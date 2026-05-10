#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ACTION_NAME="$(basename "$0")"
LOGS_DIR="${SCRIPT_DIR}/logs"
LOG_FILE=""
LOG_TO_FILE_ENABLED=0

init_script() {
  INTERACTIVE_TTY_AVAILABLE=0
  if [[ -t 0 && -t 1 && -r /dev/tty && -w /dev/tty ]]; then
    INTERACTIVE_TTY_AVAILABLE=1
  fi
  export INTERACTIVE_TTY_AVAILABLE
  if [[ "${FS_LOG_TO_FILE:-false}" == "true" ]]; then
    LOG_TO_FILE_ENABLED=1
    mkdir -p "${LOGS_DIR}"
    local timestamp
    timestamp="$(date '+%Y%m%d-%H%M%S')"
    LOG_FILE="${LOGS_DIR}/$(basename "$0" .sh)-${timestamp}.log"
    export LOG_FILE
    export LOG_TO_FILE_ENABLED
    exec > >(tee -a "${LOG_FILE}") 2>&1
  else
    export LOG_FILE
    export LOG_TO_FILE_ENABLED
  fi
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] Starting ${ACTION_NAME}"
  echo "Repository root: ${REPO_ROOT}"
  if [[ "${LOG_TO_FILE_ENABLED}" -eq 1 ]]; then
    echo "Log file: ${LOG_FILE}"
  fi
  trap 'on_err "$BASH_COMMAND" "$LINENO" "$?"' ERR
  trap 'on_exit "$?"' EXIT
}

on_err() {
  local command="$1"
  local line_number="$2"
  local exit_code="$3"
  echo
  echo "[ERROR] Command failed with exit code ${exit_code} at line ${line_number}: ${command}"
}

on_exit() {
  local exit_code="$1"
  echo
  if [[ "${exit_code}" -eq 0 ]]; then
    echo "[OK] ${ACTION_NAME} finished successfully"
  else
    echo "[FAIL] ${ACTION_NAME} finished with exit code ${exit_code}"
  fi
  if [[ "${LOG_TO_FILE_ENABLED:-0}" -eq 1 ]]; then
    echo "Log file: ${LOG_FILE}"
  fi
  pause_if_requested
}

pause_if_requested() {
  local pause_mode="${FS_PAUSE_ON_EXIT:-auto}"
  case "${pause_mode}" in
    never|false) return ;;
    auto)
      if [[ "${INTERACTIVE_TTY_AVAILABLE:-0}" -ne 1 ]]; then
        return
      fi
      ;;
  esac
  printf 'Press Enter to close...' >/dev/tty
  read -r _ </dev/tty || true
}

is_windows_bash() {
  [[ "${OSTYPE:-}" == msys* || "${OSTYPE:-}" == cygwin* || -n "${MSYSTEM:-}" ]]
}

to_wsl_path() {
  local input_path="$1"
  if [[ "${input_path}" =~ ^/([A-Za-z])/(.*)$ ]]; then
    printf '/mnt/%s/%s\n' "${BASH_REMATCH[1],,}" "${BASH_REMATCH[2]}"
    return
  fi
  if [[ "${input_path}" =~ ^([A-Za-z]):[\\/](.*)$ ]]; then
    local remainder="${BASH_REMATCH[2]}"
    remainder="${remainder//\\//}"
    printf '/mnt/%s/%s\n' "${BASH_REMATCH[1],,}" "${remainder}"
    return
  fi
  printf '%s\n' "${input_path}"
}

relaunch_in_wsl_if_needed() {
  if ! is_windows_bash || [[ "${FS_DISABLE_WSL_RELAUNCH:-false}" == "true" ]]; then
    return
  fi
  if ! command -v wsl.exe >/dev/null 2>&1; then
    echo "wsl.exe is required to relaunch this script in WSL"
    return 1
  fi
  local repo_root_wsl
  repo_root_wsl="$(to_wsl_path "${REPO_ROOT}")"
  local command_string
  command_string="cd $(printf '%q' "${repo_root_wsl}") && env FS_ALREADY_IN_WSL_RELAUNCH=1"
  local forwarded_env_var
  local forwarded_env_vars=(FS_PAUSE_ON_EXIT FS_LOG_TO_FILE K3D_BIN FS_WSL_DISTRO)
  for forwarded_env_var in "${forwarded_env_vars[@]}"; do
    if [[ -n "${!forwarded_env_var+x}" ]]; then
      command_string+=" $(printf '%q' "${forwarded_env_var}=${!forwarded_env_var}")"
    fi
  done
  command_string+=" bash $(printf '%q' "./scripts/${ACTION_NAME}")"
  local script_arg
  for script_arg in "$@"; do
    command_string+=" $(printf '%q' "${script_arg}")"
  done
  local -a wsl_launcher=(wsl.exe)
  if [[ -n "${FS_WSL_DISTRO:-}" ]]; then
    wsl_launcher+=(-d "${FS_WSL_DISTRO}")
  fi
  echo "Detected Windows bash environment, relaunching in WSL"
  echo "Repository root in WSL: ${repo_root_wsl}"
  echo "Script in WSL: ./scripts/${ACTION_NAME}"
  exec "${wsl_launcher[@]}" bash -lc "${command_string}"
}

k3d_cmd() {
  local k3d_path=""
  if [[ -n "${K3D_BIN:-}" ]]; then
    k3d_path="${K3D_BIN}"
  else
    local candidate
    for candidate in k3d k3d.exe; do
      if command -v "${candidate}" >/dev/null 2>&1; then
        k3d_path="$(command -v "${candidate}")"
        break
      fi
    done
  fi
  if [[ -z "${k3d_path}" ]]; then
    echo "Missing required command: k3d"
    if is_windows_bash; then
      echo "Hint: this script is running in Git Bash on Windows, not in WSL"
      echo "Hint: if k3d is installed only in WSL, run the script from WSL so the tool is on PATH there"
    fi
    echo "Hint: install k3d or export K3D_BIN with the full path to k3d/k3d.exe"
    return 1
  fi
  "${k3d_path}" "$@"
}

cluster_name="${K3D_CLUSTER_NAME:-finscheduler}"

relaunch_in_wsl_if_needed "$@"
init_script
if ! k3d_cmd cluster list 2>/dev/null | awk 'NR > 1 { print $1 }' | grep -Fxq "${cluster_name}"; then
  echo "k3d cluster not found: ${cluster_name}"
  exit 0
fi
echo "Deleting k3d cluster ${cluster_name}"
k3d_cmd cluster delete "${cluster_name}"
echo "Cluster is deleted"

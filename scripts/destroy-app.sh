#!/usr/bin/env bash

set -Eeuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ACTION_NAME="$(basename "$0")"
LOGS_DIR="${SCRIPT_DIR}/logs"
K8S_BASE_DIR="${REPO_ROOT}/k8s/base"
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
  local forwarded_env_vars=(FS_PAUSE_ON_EXIT FS_LOG_TO_FILE KUBECTL_BIN)
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

kubectl_cmd() {
  local kubectl_path=""
  if [[ -n "${KUBECTL_BIN:-}" ]]; then
    kubectl_path="${KUBECTL_BIN}"
  else
    local candidate
    for candidate in kubectl kubectl.exe; do
      if command -v "${candidate}" >/dev/null 2>&1; then
        kubectl_path="$(command -v "${candidate}")"
        break
      fi
    done
  fi
  if [[ -z "${kubectl_path}" ]]; then
    echo "Missing required command: kubectl"
    echo "Hint: install kubectl or export KUBECTL_BIN with the full path to kubectl/kubectl.exe"
    return 1
  fi
  "${kubectl_path}" "$@"
}

current_kube_context() {
  kubectl_cmd config current-context 2>/dev/null || true
}

cluster_name_for_context() {
  local context_name="$1"
  kubectl_cmd config view --raw -o "jsonpath={.contexts[?(@.name==\"${context_name}\")].context.cluster}" 2>/dev/null
}

server_for_cluster_entry() {
  local cluster_entry_name="$1"
  kubectl_cmd config view --raw -o "jsonpath={.clusters[?(@.name==\"${cluster_entry_name}\")].cluster.server}" 2>/dev/null
}

prepare_kubectl_access() {
  local context_name
  context_name="$(current_kube_context)"
  if [[ -z "${context_name}" ]]; then
    echo "kubectl current context is not set"
    echo "Run ./scripts/create-test-cluster.sh first or switch to a valid context"
    return 1
  fi
  local cluster_entry_name
  cluster_entry_name="$(cluster_name_for_context "${context_name}")"
  if [[ -z "${cluster_entry_name}" ]]; then
    return
  fi
  local current_server
  current_server="$(server_for_cluster_entry "${cluster_entry_name}")"
  if [[ -z "${current_server}" ]]; then
    return
  fi
  local normalized_server="${current_server//0.0.0.0/127.0.0.1}"
  normalized_server="${normalized_server//\[::\]/127.0.0.1}"
  if [[ "${normalized_server}" != "${current_server}" ]]; then
    echo "Normalizing kubeconfig server for ${context_name}"
    echo "  from: ${current_server}"
    echo "  to:   ${normalized_server}"
    kubectl_cmd config set-cluster "${cluster_entry_name}" --server="${normalized_server}" >/dev/null
  fi
}

relaunch_in_wsl_if_needed "$@"
init_script
prepare_kubectl_access
echo "Deleting app ingress"
kubectl_cmd delete -f "${K8S_BASE_DIR}/ingress.yaml" --ignore-not-found=true
echo "Deleting web manifests"
kubectl_cmd delete -k "${K8S_BASE_DIR}/web" --ignore-not-found=true
echo "Deleting API manifests"
kubectl_cmd delete -k "${K8S_BASE_DIR}/api" --ignore-not-found=true
echo "Application is deleted"

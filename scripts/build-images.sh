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
    never|false)
      return
      ;;
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
  if ! is_windows_bash; then
    return
  fi

  if [[ "${FS_DISABLE_WSL_RELAUNCH:-false}" == "true" ]]; then
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
  local forwarded_env_vars=(
    FS_PAUSE_ON_EXIT
    FS_LOG_TO_FILE
    DOCKER_BIN
  )

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

find_docker_cmd() {
  if [[ -n "${DOCKER_BIN:-}" ]]; then
    printf '%s\n' "${DOCKER_BIN}"
    return 0
  fi

  local candidate
  for candidate in docker docker.exe; do
    if command -v "${candidate}" >/dev/null 2>&1; then
      command -v "${candidate}"
      return 0
    fi
  done

  return 1
}

docker_cmd() {
  local docker_path
  if ! docker_path="$(find_docker_cmd)"; then
    echo "Missing required command: docker"
    echo "Hint: install Docker CLI or export DOCKER_BIN with the full path to docker/docker.exe"
    return 1
  fi

  "${docker_path}" "$@"
}

build_images() {
  echo "Building finscheduler-api:latest"
  docker_cmd build \
    -f "${REPO_ROOT}/FinScheduler.API/Dockerfile" \
    -t finscheduler-api:latest \
    "${REPO_ROOT}/FinScheduler.API"

  echo "Building finscheduler-web:latest"
  docker_cmd build \
    -f "${REPO_ROOT}/finscheduler-web/Dockerfile" \
    -t finscheduler-web:latest \
    "${REPO_ROOT}"

  echo "Docker images are ready"
}

relaunch_in_wsl_if_needed "$@"
init_script
build_images

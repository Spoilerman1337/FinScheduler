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

check_secret_files() {
  local secret_file
  for secret_file in "$@"; do
    if [[ ! -f "${secret_file}" ]]; then
      echo "Secret file not found: ${secret_file}"
      return 1
    fi

    if grep -Eqi 'replace-me|replace_with|change-me|fill-me|your-value-here' "${secret_file}"; then
      echo "Placeholder value detected in ${secret_file}"
      return 1
    fi
  done
}

init_script
check_secret_files \
  "${REPO_ROOT}/k8s/base/api/secret.yaml" \
  "${REPO_ROOT}/k8s/base/storage/minio/secret.yaml" \
  "${REPO_ROOT}/k8s/base/observability/grafana/admin-secret.yaml" \
  "${REPO_ROOT}/k8s/base/observability/grafana/db-secret.yaml"

echo "Kubernetes secret files look ready"

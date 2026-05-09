#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

required_commands=(
  kubectl
)

for command_name in "${required_commands[@]}"; do
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    echo "${command_name} is required"
    exit 1
  fi
done

secret_files=(
  "${REPO_ROOT}/k8s/base/api/secret.yaml"
  "${REPO_ROOT}/k8s/base/observability/grafana-db-secret.yaml"
  "${REPO_ROOT}/k8s/base/storage/minio-secret.yaml"
)

for secret_file in "${secret_files[@]}"; do
  if [[ ! -f "${secret_file}" ]]; then
    echo "Secret file not found: ${secret_file}"
    exit 1
  fi
done

for secret_file in "${secret_files[@]}"; do
  if grep -Fq "replace-with-your-" "${secret_file}"; then
    echo "Placeholder value detected in ${secret_file}"
    exit 1
  fi
done

if grep -Fq "GF_SECURITY_ADMIN_PASSWORD: admin" "${REPO_ROOT}/k8s/base/observability/grafana-admin-secret.yaml"; then
  echo "Grafana admin password is still set to the default value"
  exit 1
fi

echo "Kubernetes secrets look ready"

#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
KUSTOMIZE_ROOT="${REPO_ROOT}/k8s/base"

required_commands=(
  kubectl
)

for command_name in "${required_commands[@]}"; do
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    echo "${command_name} is required"
    exit 1
  fi
done

SKIP_BUILD="${SKIP_BUILD:-false}"
WAIT_TIMEOUT="${WAIT_TIMEOUT:-180s}"

if [[ "${SKIP_BUILD}" != "true" ]]; then
  "${SCRIPT_DIR}/build-images.sh"
fi

"${SCRIPT_DIR}/check-k8s-secrets.sh"

echo "Rendering manifests"
kubectl kustomize "${KUSTOMIZE_ROOT}" >/dev/null

echo "Applying manifests"
kubectl apply -k "${KUSTOMIZE_ROOT}"

echo "Waiting for core deployments"
kubectl rollout status deployment/finscheduler-api --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/finscheduler-web --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/minio --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/prometheus --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/mimir --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/pyroscope --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/tempo --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/loki --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/alloy-logs --timeout="${WAIT_TIMEOUT}"
kubectl rollout status deployment/grafana --timeout="${WAIT_TIMEOUT}"

echo
echo "Contour is deployed"
echo "App ingress: http://www.finscheduler.local"
echo "Grafana ingress: http://grafana.finscheduler.local"

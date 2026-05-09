#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required"
  exit 1
fi

echo "Building finscheduler-api:latest"
docker build \
  -f "${REPO_ROOT}/FinScheduler.API/Dockerfile" \
  -t finscheduler-api:latest \
  "${REPO_ROOT}/FinScheduler.API"

echo "Building finscheduler-web:latest"
docker build \
  -f "${REPO_ROOT}/finscheduler-web/Dockerfile" \
  -t finscheduler-web:latest \
  "${REPO_ROOT}"

echo "Docker images are ready"

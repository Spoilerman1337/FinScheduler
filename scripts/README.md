# Scripts

Helper bash scripts for the local Kubernetes test contour.

## Prerequisites

- `docker`
- `kubectl`
- real values in:
  - `k8s/base/api/secret.yaml`
  - `k8s/base/observability/grafana-db-secret.yaml`
  - `k8s/base/observability/grafana-admin-secret.yaml`
  - `k8s/base/storage/minio-secret.yaml`

## Main flow

Build images, validate secrets, apply `k8s/base`, and wait for the main deployments:

```bash
./scripts/deploy-test-contour.sh
```

Skip docker build if images are already available in the cluster runtime:

```bash
SKIP_BUILD=true ./scripts/deploy-test-contour.sh
```

## Other scripts

- `build-images.sh` - builds `finscheduler-api:latest` and `finscheduler-web:latest`
- `check-k8s-secrets.sh` - fails if placeholder values are still present
- `destroy-test-contour.sh` - deletes the contour from the cluster

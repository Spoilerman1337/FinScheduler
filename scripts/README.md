# Scripts

Helper bash scripts for the local Kubernetes test contour.

By default, each script logs only to the current terminal.
Set `FS_LOG_TO_FILE=true` if you also want a timestamped log file in `scripts/logs/`.
When a script runs in an interactive terminal, it waits for Enter before closing, even after a failure, so the output stays visible for investigation.
For WSL and similar setups, the scripts also normalize `k3d` kubeconfig endpoints from `0.0.0.0` or `[::]` to `127.0.0.1` before talking to the cluster.
If a script is started from Git Bash on Windows, it now relaunches itself automatically inside WSL. Set `FS_WSL_DISTRO` if you want to force a specific distro.

## Prerequisites

- `docker`
- `kubectl`
- `bash`
- `k3d` for cluster create/delete scripts
- configured values in:
  - `k8s/base/api/secret.yaml`
  - `k8s/base/observability/grafana/db-secret.yaml`
  - `k8s/base/observability/grafana/admin-secret.yaml`
  - `k8s/base/storage/minio/secret.yaml`

## Main flow

Create the local `k3d` cluster with ingress ports on `80/443`:

```bash
./scripts/create-test-cluster.sh
```

Override defaults if needed:

```bash
K3D_CLUSTER_NAME=finscheduler K3D_AGENTS_COUNT=2 K3D_API_PORT=6550 ./scripts/create-test-cluster.sh
```

Then deploy the layers independently:

```bash
./scripts/deploy-storage.sh
./scripts/deploy-observability.sh
./scripts/deploy-app.sh
```

Or run the full contour in the same order:

```bash
./scripts/deploy-test-contour.sh
```

Skip docker build when deploying the application layer or the full contour if images are already available in the cluster runtime:

```bash
SKIP_BUILD=true ./scripts/deploy-app.sh
SKIP_BUILD=true ./scripts/deploy-test-contour.sh
```

When a script runs in an interactive terminal, it now waits for Enter at the end so the log stays visible.

## Other scripts

- `create-test-cluster.sh` - creates or reuses the local `k3d` cluster and switches `kubectl` context
- `destroy-test-cluster.sh` - deletes the local `k3d` cluster
- `deploy-storage.sh` - applies the shared MinIO storage layer
- `deploy-observability.sh` - applies the Grafana/Mimir/Tempo/Loki/Pyroscope/Alloy layer
- `deploy-app.sh` - builds images if needed and applies `api`, `web`, and the application ingress
- `build-images.sh` - builds `finscheduler-api:latest` and `finscheduler-web:latest`
- `check-k8s-secrets.sh` - fails if placeholder values are still present
- `destroy-storage.sh` - deletes the storage layer
- `destroy-observability.sh` - deletes the observability layer
- `destroy-app.sh` - deletes `api`, `web`, and the application ingress
- `destroy-test-contour.sh` - deletes the contour from the cluster

## Environment variables

- `WAIT_TIMEOUT` - rollout wait timeout, default `180s`
- `SKIP_BUILD=true` - skips docker builds in `deploy-app.sh` and `deploy-test-contour.sh`
- `IMPORT_TO_K3D=auto|true|false` - controls whether built images are imported into the local `k3d` cluster
- `FS_PAUSE_ON_EXIT=auto|always|never` - controls the final Enter prompt, default `auto`
- `FS_LOG_TO_FILE=true` - also writes the terminal output to `scripts/logs/*.log`
- `FS_WSL_DISTRO` - optional WSL distro name for Windows-to-WSL relaunch
- `FS_DISABLE_WSL_RELAUNCH=true` - disables the automatic Windows-to-WSL relaunch
- `K3D_CLUSTER_NAME` - cluster name, default `finscheduler`
- `K3D_AGENTS_COUNT` - worker count for `create-test-cluster.sh`, default `1`
- `K3D_API_PORT` - API port for `create-test-cluster.sh`, default `6550`
- `K3D_HTTP_PORT` - host HTTP port for `create-test-cluster.sh`, default `80`
- `K3D_HTTPS_PORT` - host HTTPS port for `create-test-cluster.sh`, default `443`

# Go Sample API

A simple Go HTTP API built with [Echo](https://echo.labstack.com/) for testing **ArgoCD GitOps** deployments across multiple environments.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/` | Returns app info (version, hostname, environment, timestamp) |
| `GET` | `/health` | Liveness probe |
| `GET` | `/ready` | Readiness probe |

## Run Locally

```bash
go run main.go
# or with custom port/env:
PORT=3000 APP_ENV=local APP_VERSION=1.0.0 go run main.go
```

## Docker

```bash
docker build -t go-sample-api .
docker run -p 8080:8080 go-sample-api
```

## Environments

| Environment | Namespace | Replicas | Ingress Host | ArgoCD Sync |
|-------------|-----------|----------|--------------|-------------|
| Dev | `dev` | 1 | `dev-api.angang.biz.id` | Auto |
| Staging | `staging` | 2 | `staging-api.angang.biz.id` | Manual |
| Prod | `prod` | 2 | `api.angang.biz.id` | Manual |

## Deployment Flow

### Automatic (Dev)
1. Push code to `main` branch
2. GitHub Actions runs tests → builds Docker image → pushes to `ghcr.io`
3. CI updates `kubernetes/overlays/dev/kustomization.yaml` with new image tag
4. ArgoCD auto-syncs the dev environment

### Manual Promotion (Staging / Prod)
1. Go to GitHub Actions → **Promote Environment** workflow
2. Click **Run workflow**
3. Select target environment (`staging` or `prod`)
4. Enter the image tag to deploy (e.g., `abc1234`)
5. ArgoCD detects the manifest change → click **Sync** in ArgoCD UI

## Setup

1. Create this repo on GitHub: `https://github.com/trianggianggara/go-sample-api`
2. Push this code to the `main` branch
3. ArgoCD Application CRDs are pre-configured in `k8s-gcp/kubernetes/argocd/*.yaml`
4. Enable GitHub Actions in the repo settings
5. Ensure the GitHub Container Registry package visibility is set appropriately

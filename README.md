<div align="center">

# 🚀 go-web-app

**A production-grade Go web server — built simple, deployed like a pro.**

[![CI/CD Pipeline](https://github.com/DevSars24/go-web-app/actions/workflows/cicd.yml/badge.svg)](https://github.com/DevSars24/go-web-app/actions/workflows/cicd.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22.5-00ADD8?logo=go)](https://go.dev/)
[![Docker Image](https://img.shields.io/badge/DockerHub-sars2006%2Fgo--web--app-2496ED?logo=docker)](https://hub.docker.com/r/sars2006/go-web-app)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.24%2B-326CE5?logo=kubernetes)](https://kubernetes.io/)

*Designed and maintained by **[Saurabh Singh Rajput](https://github.com/DevSars24)***

</div>

---

## 📖 Table of Contents

- [What Is This Project?](#-what-is-this-project)
- [Architecture Overview](#-architecture-overview)
- [Project Structure](#-project-structure)
- [Technology Stack](#-technology-stack)
- [Quick Start — Run Locally](#-quick-start--run-locally)
- [Docker — Build & Run](#-docker--build--run)
- [CI/CD Pipeline — GitHub Actions](#️-cicd-pipeline--github-actions)
- [Kubernetes — Deploy to a Cluster](#️-kubernetes--deploy-to-a-cluster)
- [Helm — Package Manager for K8s](#️-helm--package-manager-for-k8s)
- [GitOps — ArgoCD Continuous Delivery](#-gitops--argocd-continuous-delivery)
- [File-by-File Explanation](#-file-by-file-explanation)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)

---

## 🎯 What Is This Project?

`go-web-app` is a **minimal Go HTTP web server** that serves four static HTML pages (Home, Courses, About, Contact). The application itself is intentionally simple — a single `main.go` file with no external dependencies.

**The real learning is in the infrastructure surrounding it.**

This project demonstrates a **complete, end-to-end DevOps workflow** that mirrors what you'd find in a real production environment at a technology company:

```
Developer pushes code
        │
        ▼
GitHub Actions CI/CD
 ├── Tests (go test -race)
 ├── Static analysis (go vet)
 ├── Docker image build + push
 └── GitOps manifest update
        │
        ▼
ArgoCD detects the git change
        │
        ▼
Kubernetes Cluster (EKS/GKE/AKS/Minikube)
 ├── Namespace: webapps
 ├── Deployment (2 replicas, rolling update)
 ├── Service (ClusterIP)
 └── Ingress (NGINX)
        │
        ▼
User accesses https://go-web-app.local
```

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          DEVELOPER MACHINE                          │
│   main.go  →  git push  →  GitHub                                  │
└────────────────────────────────┬────────────────────────────────────┘
                                 │ triggers
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    GITHUB ACTIONS (CI/CD)                           │
│                                                                     │
│  Job 1: test          Job 2: build-and-push    Job 3: update-manifest│
│  ─────────────────    ─────────────────────    ────────────────────  │
│  go vet ./...    →    docker build         →   sed image tag         │
│  go test -race        docker push               git commit & push    │
│                       sars2006/go-web-app                           │
│                       :sha-abc1234                                  │
└──────────────────────────────────────────────┬──────────────────────┘
                                               │ commits to
                                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      GIT REPOSITORY                                 │
│   helm/go-web-app/values.yaml  ← image tag updated                 │
└───────────────────────────────┬─────────────────────────────────────┘
                                │ ArgoCD polls every 3 minutes
                                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    KUBERNETES CLUSTER                               │
│                                                                     │
│   Namespace: webapps                                                │
│   ┌─────────────────────────────────────────┐                      │
│   │  Ingress (NGINX)                        │ ← go-web-app.local   │
│   │         │                               │                      │
│   │         ▼                               │                      │
│   │  Service: go-web-app (ClusterIP:80)     │                      │
│   │         │                               │                      │
│   │    ┌────┴────┐                          │                      │
│   │    ▼         ▼                          │                      │
│   │  Pod 1     Pod 2   (2 replicas)         │                      │
│   │  :8080     :8080                        │                      │
│   └─────────────────────────────────────────┘                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
go-web-app/
│
├── main.go                          # Application entrypoint (Go HTTP server)
├── main_test.go                     # Unit & integration tests
├── go.mod                           # Go module definition (dependency lock)
│
├── static/                          # Static HTML pages served by the app
│   ├── home.html
│   ├── courses.html
│   ├── about.html
│   └── contact.html
│
├── Dockerfile                       # Multi-stage container build definition
│
├── .github/
│   └── workflows/
│       └── cicd.yml                 # GitHub Actions CI/CD pipeline (3 jobs)
│
├── K8s/
│   └── manifests/                   # Raw Kubernetes YAML (learning / fallback)
│       ├── namespace.yaml           # Dedicated `webapps` namespace
│       ├── deployment.yaml          # App deployment (replicas, probes, limits)
│       ├── service.yaml             # Internal ClusterIP service
│       └── ingress.yaml             # External NGINX ingress routing
│
├── helm/
│   └── go-web-app/                  # Helm chart (preferred deployment method)
│       ├── Chart.yaml               # Chart metadata and versioning
│       ├── values.yaml              # All configurable defaults (documented)
│       └── template/
│           ├── deployment.yaml      # Helm-templated deployment
│           ├── service.yaml         # Helm-templated service
│           └── ingress.yaml         # Helm-templated ingress (TLS-ready)
│
├── argocd-app.yaml                  # ArgoCD Application (GitOps registration)
├── .gitignore                       # Excludes binaries, secrets, IDE files
└── README.md                        # This file
```

---

## 🛠️ Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Language** | Go 1.22.5 | HTTP web server, zero external deps |
| **Testing** | `go test`, `httptest` | Unit tests, table-driven test patterns |
| **Container** | Docker (multi-stage) | Produce a ~20MB distroless image |
| **Registry** | DockerHub | Store and version container images |
| **CI/CD** | GitHub Actions | Automate test → build → deploy |
| **Orchestration** | Kubernetes 1.24+ | Run, scale, self-heal containers |
| **Package Mgr** | Helm 3 | Parameterise K8s manifests |
| **GitOps** | ArgoCD | Sync Git state → cluster state |

---

## ⚡ Quick Start — Run Locally

**Prerequisites:** Go 1.22+ installed

```bash
# 1. Clone the repository
git clone https://github.com/DevSars24/go-web-app.git
cd go-web-app

# 2. Run tests to verify everything is working
go test -race -v ./...

# 3. Start the server
go run main.go

# 4. Visit in your browser
open http://localhost:8080/home
```

Available routes:

| Route | Page |
|-------|------|
| `GET /home` | Home page |
| `GET /courses` | Courses catalogue |
| `GET /about` | About page |
| `GET /contact` | Contact page |
| `GET /healthz` | Liveness probe (returns `ok`) |
| `GET /readyz` | Readiness probe (returns `ready`) |

---

## 🐳 Docker — Build & Run

### Build the Image

```bash
# Build using the multi-stage Dockerfile
docker build -t go-web-app:local .

# Inspect the final image size (should be ~20MB)
docker images go-web-app:local
```

### Run the Container

```bash
docker run --rm -p 8080:8080 go-web-app:local

# Visit: http://localhost:8080/home
```

### Pull from DockerHub

```bash
docker pull sars2006/go-web-app:latest
docker run --rm -p 8080:8080 sars2006/go-web-app:latest
```

### Why Multi-Stage?

| Stage | Base Image | Size | Contains |
|-------|-----------|------|---------|
| `builder` | `golang:1.22.5-alpine3.20` | ~300MB | Compiler, source, deps |
| `final` | `gcr.io/distroless/static-debian12` | ~20MB | Binary only |

The Go compiler and source code **never ship to production**. Only the compiled binary does.

---

## ⚙️ CI/CD Pipeline — GitHub Actions

The pipeline (`.github/workflows/cicd.yml`) is divided into three jobs that run in strict order:

```
[test] ──→ [build-and-push] ──→ [update-manifest]
```

### Job 1: `test` (runs on every push + PR)

```
✔ go vet ./...          — static analysis (catches bugs the compiler misses)
✔ go test -race -v ./...— unit tests with race detector enabled
```

### Job 2: `build-and-push` (runs on push to `main` only)

```
✔ Set up Docker Buildx
✔ Authenticate with DockerHub (uses GitHub Secrets)
✔ Build Docker image
✔ Push with two tags:
    sars2006/go-web-app:latest
    sars2006/go-web-app:sha-abc1234
✔ Layer cache stored in registry (saves ~2 min per build)
```

### Job 3: `update-manifest` (runs after successful build)

```
✔ Updates image tag in helm/go-web-app/values.yaml
✔ Commits with message: "chore: update image tag to sha-xxx [skip ci]"
✔ [skip ci] prevents an infinite loop of pipeline triggers
```

### Required GitHub Secrets

Set these in **GitHub → Settings → Secrets → Actions**:

| Secret | Value |
|--------|-------|
| `DOCKER_USERNAME` | Your DockerHub username |
| `DOCKER_PASSWORD` | Your DockerHub access token (not password!) |

---

## ☸️ Kubernetes — Deploy to a Cluster

### Prerequisites

- A running Kubernetes cluster (Minikube, kind, EKS, GKE, AKS)
- `kubectl` configured to point to your cluster
- NGINX Ingress Controller installed

```bash
# Install NGINX Ingress Controller (one-time)
helm upgrade --install ingress-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --namespace ingress-nginx --create-namespace
```

### Deploy with Raw Manifests

```bash
# Apply in dependency order
kubectl apply -f K8s/manifests/namespace.yaml
kubectl apply -f K8s/manifests/deployment.yaml
kubectl apply -f K8s/manifests/service.yaml
kubectl apply -f K8s/manifests/ingress.yaml
```

### Verify the Deployment

```bash
# Check pod status
kubectl get pods -n webapps

# Check pod logs
kubectl logs -n webapps -l app.kubernetes.io/name=go-web-app

# Test health endpoint directly from inside the cluster
kubectl exec -n webapps deploy/go-web-app -- wget -qO- localhost:8080/healthz
```

### Access the Application (Minikube)

```bash
# Get the Minikube IP
minikube ip

# Add to /etc/hosts (Linux/Mac) or C:\Windows\System32\drivers\etc\hosts (Windows)
echo "$(minikube ip) go-web-app.local" | sudo tee -a /etc/hosts

# Visit: http://go-web-app.local/home
```

### Useful kubectl Commands

```bash
# Watch rolling update progress
kubectl rollout status deployment/go-web-app -n webapps

# Roll back to the previous version
kubectl rollout undo deployment/go-web-app -n webapps

# Scale up/down (hot-patch, not GitOps)
kubectl scale deployment/go-web-app --replicas=4 -n webapps

# Show resource usage
kubectl top pods -n webapps
```

---

## ⛵ Helm — Package Manager for K8s

Helm is the recommended deployment method. It replaces hardcoded values with parameterised templates, making the same chart reusable across dev, staging, and production environments.

### Install with Helm

```bash
# Install (or upgrade) the chart
helm upgrade --install go-web-app ./helm/go-web-app \
  --namespace webapps \
  --create-namespace

# Install with a specific image tag (e.g., from CI)
helm upgrade --install go-web-app ./helm/go-web-app \
  --set image.tag=sha-abc1234 \
  --namespace webapps
```

### Useful Helm Commands

```bash
# Preview the rendered YAML without installing
helm template go-web-app ./helm/go-web-app

# Check chart for issues
helm lint ./helm/go-web-app

# View installed releases
helm list -n webapps

# Inspect currently deployed values
helm get values go-web-app -n webapps

# Roll back to a previous release
helm rollback go-web-app 1 -n webapps

# Uninstall (removes all K8s resources)
helm uninstall go-web-app -n webapps
```

### Environment-Specific Overrides

Create a `values-production.yaml` for production settings without touching the default `values.yaml`:

```yaml
# values-production.yaml
replicaCount: 4
image:
  tag: sha-abc1234  # specific SHA from CI
ingress:
  host: app.yourdomain.com
  tls:
    enabled: true
resources:
  limits:
    memory: "256Mi"
    cpu: "500m"
```

```bash
helm upgrade --install go-web-app ./helm/go-web-app \
  -f values-production.yaml \
  --namespace webapps
```

---

## 🔄 GitOps — ArgoCD Continuous Delivery

ArgoCD is a Kubernetes operator that watches this Git repository and automatically applies any changes to the cluster. It is the CD (Continuous Delivery) half of the pipeline.

### Install ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd \
  -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
kubectl wait --for=condition=Available deployment/argocd-server \
  -n argocd --timeout=120s
```

### Access ArgoCD UI

```bash
# Port-forward the ArgoCD server
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get the initial admin password
kubectl get secret argocd-initial-admin-secret -n argocd \
  -o jsonpath="{.data.password}" | base64 -d

# Open: https://localhost:8080  (Username: admin)
```

### Register the Application

```bash
kubectl apply -f argocd-app.yaml
```

After this, every commit to `main` will automatically deploy to the cluster within ~3 minutes. The deployment history is your git history.

### GitOps vs Traditional CD

| | Traditional (Push) | GitOps (ArgoCD) |
|--|----|----|
| Cluster credentials | Stored in CI | Never leave cluster |
| Audit trail | Pipeline logs | Git commits |
| Rollback | Re-run old pipeline | `git revert` |
| Drift detection | None | Automatic |
| Manual changes | Persist forever | Auto-reverted |

---

## 📚 File-by-File Explanation

Every file in this project is extensively commented with **the reasoning behind each decision**, not just what the code does.

| File | Key Concepts Explained |
|------|----------------------|
| [`main.go`](./main.go) | Graceful shutdown, SIGTERM handling, HTTP timeouts, health probe design |
| [`main_test.go`](./main_test.go) | Table-driven tests, httptest, race detector, test strategy |
| [`Dockerfile`](./Dockerfile) | Multi-stage builds, distroless images, non-root user, static binary flags |
| [`.gitignore`](./.gitignore) | Why compiled binaries don't belong in git, secret protection |
| [`cicd.yml`](./.github/workflows/cicd.yml) | GitOps delivery model, Docker layer caching, concurrency control |
| [`namespace.yaml`](./K8s/manifests/namespace.yaml) | Why `default` namespace is dangerous, RBAC isolation |
| [`deployment.yaml`](./K8s/manifests/deployment.yaml) | Rolling updates, resource limits, liveness vs readiness probes |
| [`service.yaml`](./K8s/manifests/service.yaml) | Service types, named ports, kube-proxy load balancing |
| [`ingress.yaml`](./K8s/manifests/ingress.yaml) | Layer-7 routing, NGINX annotations, TLS termination |
| [`Chart.yaml`](./helm/go-web-app/Chart.yaml) | Helm SemVer, chart version vs app version |
| [`values.yaml`](./helm/go-web-app/values.yaml) | Single source of truth, environment-specific overrides |
| [`argocd-app.yaml`](./argocd-app.yaml) | GitOps model, drift detection, selfHeal, cascade delete |

---

## 🔧 Troubleshooting

### Pods stuck in `Pending`
```bash
kubectl describe pod -n webapps -l app.kubernetes.io/name=go-web-app
# Look for: Insufficient CPU/memory → increase node size or reduce resource requests
# Look for: No nodes available → check node count and taints
```

### Pods stuck in `CrashLoopBackOff`
```bash
kubectl logs -n webapps -l app.kubernetes.io/name=go-web-app --previous
# Check for: port already in use, file not found (static/ missing), OOM kill
```

### `ImagePullBackOff`
```bash
# The image tag in values.yaml doesn't exist in the registry
# Check DockerHub for available tags, then:
helm upgrade go-web-app ./helm/go-web-app --set image.tag=<valid-tag> -n webapps
```

### Cannot reach the app via Ingress
```bash
# Check Ingress controller is running
kubectl get pods -n ingress-nginx

# Check Ingress has an Address assigned
kubectl get ingress -n webapps

# Ensure /etc/hosts has the correct entry
cat /etc/hosts | grep go-web-app
```

### ArgoCD shows `OutOfSync`
```bash
# Manually trigger a sync
argocd app sync go-web-app

# Or via kubectl
kubectl patch application go-web-app -n argocd \
  --type merge -p '{"operation":{"sync":{}}}'
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-improvement`
3. Make your changes (keep files thoroughly commented!)
4. Run tests: `go test -race ./...`
5. Push and open a Pull Request

---

<div align="center">

Made with ❤️ by **[Saurabh Singh Rajput](https://github.com/DevSars24)**

*"Keep it simple. Deploy it right."*

</div>

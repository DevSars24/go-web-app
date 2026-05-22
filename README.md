# 🚀 Go Web App - Complete DevOps Project

> Created by **Saurabh Singh Rajput**  
> Beginner Friendly End-to-End DevOps Project using Docker, Kubernetes, Helm, GitHub Actions, and ArgoCD.

---

# 📌 About This Project

This project is a complete beginner-friendly DevOps workflow project designed to help students and developers understand how modern DevOps pipelines work in real-world companies.

The project demonstrates:

- CI/CD Pipeline
- Docker Containerization
- Kubernetes Deployment
- Helm Charts
- GitHub Actions Automation
- GitOps using ArgoCD

This repository is made especially for beginners who want to understand DevOps practically step by step.

---

# 🛠️ Tech Stack

- Golang
- Docker
- Kubernetes
- Helm
- GitHub Actions
- ArgoCD
- YAML

---

# 📂 Project Structure

```bash
go-web-app/
│
├── .github/workflows/
│   └── cicd.yml
│
├── K8s/manifests/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
│
├── helm/go-web-app/
│   ├── Chart.yaml
│   ├── values.yaml
│   └── templates/
│
├── static/
│
├── Dockerfile
├── argocd-app.yaml
├── README.md
└── main.go


---

🔥 Complete DevOps Workflow

1️⃣ Developer Pushes Code

The developer pushes code to GitHub.

git add .
git commit -m "updated app"
git push origin main

After push, GitHub Actions automatically starts the CI/CD pipeline.


---

⚙️ GitHub Actions CI/CD

📄 File:

.github/workflows/cicd.yml

🔥 Purpose

This file automates:

Build process

Docker image creation

Docker image push

Kubernetes deployment update



---

🐳 Dockerfile

📄 File:

Dockerfile

🔥 Purpose

Dockerfile is used to create a Docker image of the Go application.

Example Workflow

FROM golang:1.22
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]


---

🧠 Docker Concepts Used

Concept	Meaning

FROM	Base image
WORKDIR	Working directory
COPY	Copy files
RUN	Execute commands
CMD	Start application



---

☸️ Kubernetes Manifests

Folder:

K8s/manifests/

This folder contains Kubernetes YAML files.


---

📄 deployment.yaml

Purpose

Deployment manages application pods.

It ensures:

Desired number of pods run

Auto restart if pod crashes

Scaling support


Important Concepts

Field	Meaning

replicas	Number of pods
containers	App container
image	Docker image
ports	Container port



---

📄 service.yaml

Purpose

Service exposes the application inside or outside the cluster.

Types of Services

Type	Meaning

ClusterIP	Internal access
NodePort	External access
LoadBalancer	Cloud external IP



---

📄 ingress.yaml

Purpose

Ingress provides domain-based routing.

Example:

myapp.com → Kubernetes Service

Ingress helps manage:

Custom domains

HTTPS

Reverse proxy routing



---

📦 Helm Charts

Folder:

helm/go-web-app/

Helm is called the package manager of Kubernetes.


---

📄 Chart.yaml

Purpose

Contains metadata about the Helm chart.

Example:

apiVersion: v2
name: go-web-app
version: 0.1.0


---

📄 values.yaml

Purpose

Stores configurable values.

Example:

replicaCount: 2
image:
  repository: your-image

This allows customization without changing templates.


---

📄 templates/

Contains reusable Kubernetes templates.

Helm dynamically generates manifests from templates.


---

🚀 ArgoCD GitOps

📄 File:

argocd-app.yaml

Purpose

ArgoCD continuously watches GitHub repository.

If changes happen:

ArgoCD automatically syncs Kubernetes cluster

Deployment updates automatically


This process is called:

🔥 GitOps


---

🧠 What is GitOps?

GitOps means:

> Git repository becomes the single source of truth.



Whenever Git changes:

Infrastructure updates automatically.



---

🔄 Complete Flow

Developer Push Code
        ↓
GitHub Actions Runs
        ↓
Docker Image Build
        ↓
Push to Docker Registry
        ↓
Update Kubernetes Manifest
        ↓
ArgoCD Detects Changes
        ↓
Kubernetes Cluster Updated


---

🎯 DevOps Concepts Covered

✅ CI/CD
✅ Docker
✅ Kubernetes
✅ Helm
✅ GitHub Actions
✅ GitOps
✅ ArgoCD
✅ YAML Configuration
✅ Containerization
✅ Deployment Automation


---

🚀 How To Run Locally

Clone Repository

git clone <repo-url>

Run Go App

go run main.go

Build Docker Image

docker build -t go-web-app .

Run Docker Container

docker run -p 8080:8080 go-web-app


---

☸️ Kubernetes Deployment

Apply manifests:

kubectl apply -f K8s/manifests/

Check pods:

kubectl get pods

Check services:

kubectl get svc


---

📦 Helm Deployment

Install Helm chart:

helm install go-web-app ./helm/go-web-app


---

🚀 ArgoCD Deployment

Apply ArgoCD app:

kubectl apply -f argocd-app.yaml


---

💡 Why This Project is Important

This project demonstrates real-world DevOps practices used in companies.

It helps beginners understand:

How CI/CD works

How containers work

How Kubernetes deployments work

How GitOps automates infrastructure



---

👨‍💻 Author

Saurabh Singh Rajput

---

⭐ Final Note

If you are a beginner learning DevOps, this repository can help you understand the complete deployment lifecycle from code push to Kubernetes production deployment.

Happy Learning 🚀



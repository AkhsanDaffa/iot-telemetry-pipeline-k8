# 🚀 Edge IoT Telemetry Pipeline: From Legacy to Kubernetes

![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Helm](https://img.shields.io/badge/Helm-0F1689?style=for-the-badge&logo=Helm&color=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)
![Raspberry Pi](https://img.shields.io/badge/-RaspberryPi-C51A4A?style=for-the-badge&logo=Raspberry-Pi)
![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)

## 📌 Project Overview
This project demonstrates the migration of a dummy IoT sensor application from a legacy Docker Compose environment to a robust **Edge Kubernetes Cluster (K3s)** running on a Raspberry Pi. The primary goal is to establish a scalable, automated, and observable infrastructure suitable for Edge Computing scenarios.

**Development Workflow:** Prototyped and tested on WSL (Windows Subsystem for Linux) before deploying to production Raspberry Pi K3s cluster.

## ✨ Key Achievements & Features

* **Multi-Architecture CI/CD Pipeline:** Engineered a GitHub Actions workflow to automatically build and push Docker images supporting both `linux/amd64` and `linux/arm64` architectures, ensuring seamless deployment across different hardware environments.
* **Infrastructure as Code (IaC) with Helm:** Transitioned from static YAML manifests to dynamic **Helm Charts**, enabling templated deployments, easy rollbacks, and scalable application management.
* **Edge Computing Deployment:** Successfully deployed and stabilized the microservices architecture on a resource-constrained ARM64 device (Raspberry Pi) using K3s.
* **Enterprise-Grade Observability:** Replaced basic monitoring with the `kube-prometheus-stack` (Prometheus Operator & Grafana) to achieve deep Kubernetes pods and cluster metrics monitoring, auto-discovery, and visual dashboards.

## 🛠️ Tech Stack
* **Containerization:** Docker
* **Orchestration:** K3s (Lightweight Kubernetes)
* **Package Manager:** Helm v3
* **CI/CD:** GitHub Actions
* **Monitoring/Observability:** Prometheus, Grafana, Node Exporter, Kube-State-Metrics
* **Hardware:** Raspberry Pi (ARM64)

## 📂 Project Structure
```
.
├── backend/                    # Go IoT sensor application
│   ├── main.go                 # Sensor logic - generates random temperature data
│   ├── go.mod                  # Go module dependencies
│   └── Dockerfile              # Multi-stage build for amd64 & arm64
├── iot-chart/                  # Helm chart for Kubernetes deployment
│   ├── Chart.yaml              # Chart metadata & version
│   ├── values.yaml             # Default configuration values
│   └── templates/
│       └── deployment.yaml     # K8s deployment manifest
├── k8s/                        # Raw Kubernetes manifests
│   └── deployment.yaml         # Basic K8s deployment (legacy)
├── .github/workflows/          # CI/CD automation
│   └── docker-build.yml        # Multi-arch Docker image build & push
└── assets/                     # Project documentation screenshots
```

## 🎯 Skills Demonstrated
| Category | Technologies |
|----------|--------------|
| **Languages** | Go, YAML |
| **Containerization** | Docker, Docker Compose |
| **Orchestration** | Kubernetes, K3s |
| **CI/CD** | GitHub Actions |
| **Monitoring** | Prometheus, Grafana |
| **Tools** | Helm, kubectl |

## 📸 Project Showcase

### 1. Multi-Arch CI/CD Pipeline (GitHub Actions)
*Successfully building and pushing images for both AMD64 and ARM64 architectures.*
![CI/CD Pipeline](assets/cicd-success.png)

### 2. Kubernetes Pods Status & Workloads
*All microservices and enterprise monitoring stacks are running smoothly inside the K3s cluster.*
![K8s Pods Running](assets/k8s-pods-running.png)

### 3. IoT Sensor Data Generation
*The deployed dummy sensor generating and transmitting real-time telemetry data.*
![IoT Sensor Logs](assets/iot-sensor-logs.png)

### 4. Infrastructure Observability (Grafana)
*Real-time cluster compute resources monitoring using the Kube-Prometheus-Stack.*
![Grafana Dashboard](assets/grafana-dashboard.png)

---
*Developed by Akhsan - Building scalable infrastructure for the future of IoT.*
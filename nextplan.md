# Next Plan - IoT-K8s DevOps Upgrade

> Last Updated: 2026-04-05

## Context
Project ini adalah Edge IoT Telemetry Pipeline yang mendemonstrasikan migrasi dari Docker Compose ke K3s. Berdasarkan review, terdapat 3 upgrade kritis untuk meningkatkan kadar "DevOps" project ini ke level production-grade.

---

## 🚀 UPGRADE 1: App Observability (Prometheus Metrics)
**Priority: HIGH** - Foundation, wajib terlebih dahulu

### A. Arsitektur Flow (Pseudocode & Flowchart)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        IOT SENSOR → PROMETHEUS FLOW                     │
└─────────────────────────────────────────────────────────────────────────┘

    ┌──────────┐      ┌──────────────┐      ┌─────────────┐      ┌───────────┐
    │  IoT     │      │   Go App     │      │  Prometheus │      │  Grafana  │
    │  Sensor  │──────│ /metrics     │──────│  Scrapes    │──────│ Dashboard │
    └──────────┘      └──────────────┘      └─────────────┘      └───────────┘
                              │                     ▲
                              │                     │
                              ▼                     │
                       ┌──────────────┐             │
                       │ Exposed :8080│─────────────┘
                       └──────────────┘
```

**Pseudocode - Go Application dengan Prometheus:**

```go
package main

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

// 1. DEKLARASI METRICS
var (
    temperatureCurrent = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "iot_temperature_current",
        Help: "Current temperature reading",
    })
    totalReadings = promauto.NewCounter(prometheus.CounterOpts{
        Name: "iot_temperature_total_readings",
        Help: "Total number of temperature readings",
    })
    sensorErrors = promauto.NewCounter(prometheus.CounterOpts{
        Name: "iot_sensor_read_errors_total",
        Help: "Total sensor read errors",
    })
)

func main() {
    // 2. SETUP HTTP SERVER + METRICS ENDPOINT
    http.Handle("/metrics", promhttp.Handler())
    http.HandleFunc("/health", healthCheck)
    
    go collectSensorData()  // Background goroutine untuk baca sensor
    
    // 3. START SERVER
    http.ListenAndServe(":8080", nil)
}

func collectSensorData() {
    for {
        temp, err := readFromSensor()
        if err != nil {
            sensorErrors.Inc()      // Increment error counter
            continue
        }
        
        temperatureCurrent.Set(temp)  // Update gauge
        totalReadings.Inc()           // Increment counter
    }
}

func readFromSensor() (float64, error) {
    // Logic baca sensor (MQTT, serial, dll)
}
```

**Flowchart - Metrics Collection:**

```
START
  │
  ▼
┌─────────────────────┐
│  Inisialisasi       │
│  Prometheus Metrics │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Start HTTP Server  │
│  Port 8080          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Goroutine:         │
│  Read Sensor Loop    │◄──────────────┐
└──────────┬──────────┘               │
           │                          │
           ▼                          │
┌─────────────────────┐               │
│  Baca suhu sensor   │               │
└──────────┬──────────┘               │
           │                          │
           ▼                          │
    ┌──────────────┐                  │
    │ Error?       │──YES──────────────┤
    └───────┬──────┘                  │
            │ NO                      │
            ▼                         │
┌─────────────────────┐               │
│  Set gauge value    │               │
│  Inc counter        │               │
└──────────┬──────────┘               │
           │                          │
           ▼                          │
┌─────────────────────┐               │
│  Sleep (interval)   │───────────────┘
└─────────────────────┘
```

**External Request Flow:**

```
┌────────────────────────────────────────────────────────────┐
│                      METRICS SCRAPING                       │
└────────────────────────────────────────────────────────────┘

  Prometheus                    Go App (:8080)               K8s Pod
      │                              │                          │
      │  ┌─ GET /metrics ────────────┼──────────────────────►  │
      │  │                           │                          │
      │  │                           │  ┌─ iot_temperature_current 0.85
      │  │                           │  ├─ iot_temperature_total_readings 1523
      │  │                           │  └─ iot_sensor_read_errors_total 3
      │  │                           │                          │
      │  ◄───────────────────────────┼──────────────────────────┘
      │                              │
      ▼                              ▼
  Prometheus                   Response in
  stores metrics               Prometheus format
```

### C. Modify Go Application (`backend/`)
- [ ] Install library `github.com/prometheus/client_golang`
- [ ] Import packages: `prometheus`, `promhttp`, `net/http`
- [ ] Buat custom Prometheus metrics:
  - [ ] `iot_temperature_current` (Gauge) - suhu saat ini
  - [ ] `iot_temperature_total_readings` (Counter) - total pembacaan
  - [ ] `iot_sensor_read_errors_total` (Counter) - error counter
- [ ] Setup HTTP server di port `8080`
- [ ] Buat endpoint `/metrics` → `promhttp.Handler()`
- [ ] Update main loop: collect metrics, bukan hanya `fmt.Println()`
- [ ] Update `go.mod` → `go mod tidy`

### D. Update Dockerfile (`backend/Dockerfile`)
- [ ] Install `ca-certificates` di Alpine stage
- [ ] Expose port `8080`
- [ ] Update CMD/ENTRYPOINT untuk HTTP server

### E. Update Helm Chart (`iot-chart/`)
- [ ] Buat `templates/service.yaml` → expose port 8080
- [ ] Update `templates/deployment.yaml`: tambahkan port container
- [ ] Buat `templates/servicemonitor.yaml` → untuk Prometheus Operator scraping
- [ ] Update `values.yaml`: 
  - [ ] Tambahkan section `service`
  - [ ] Tambahkan label `release: prometheus` di pod spec

### F. Create Grafana Dashboard
- [ ] Buat JSON dashboard: `iot-chart/templates/grafanadashboard.yaml`
- [ ] Panel 1: Suhu Real-time (Gauge)
- [ ] Panel 2: Temperature Over Time (Graph)
- [ ] Panel 3: Total Readings Counter
- [ ] Include Prometheus expression queries

### G. Update CI/CD (`.github/workflows/`)
- [ ] Update Docker tag menjadi `v4` saat push

---

## 🚀 UPGRADE 2: Horizontal Pod Autoscaler (HPA) + Load Testing
**Priority: HIGH** - Topik interview yang sering ditanyakan

> WARNING: Ini wajib dipraktekkan, bukan hanya teori. User harus bisa demonstrate live scaling.

### A. Setup Prometheus Metrics (Prerequisite)
- [ ] Lihat UPGRADE 1 - pastikan `/metrics` endpoint sudah aktif
- [ ] Verify: `curl localhost:8080/metrics | grep iot_`

### B. Enable Metrics Server di K3s
- [ ] Check if metrics-server sudah terinstall: `kubectl get apiservices`
- [ ] Jika belum: `kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml`
- [ ] Verify: `kubectl top nodes` dan `kubectl top pods`

### C. Create HPA Manifest
- [ ] Buat `templates/hpa.yaml`:
  ```yaml
  apiVersion: autoscaling/v2
  kind: HorizontalPodAutoscaler
  metadata:
    name: iot-backend-hpa
  spec:
    scaleTargetRef:
      apiVersion: apps/v1
      kind: Deployment
      name: iot-backend
    minReplicas: 1
    maxReplicas: 5
    metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  ```
- [ ] Apply: `kubectl apply -f templates/hpa.yaml`
- [ ] Verify: `kubectl get hpa`

### D. Hands-on Load Testing (WAJIB)
- [ ] Install `hey` atau `ab` (Apache Bench):
  - macOS: `brew install hey`
  - Linux: `go install github.com/rakyll/hey@latest`
  - atau: `apt install apache2-utils` (untuk `ab`)
- [ ] Trigger load test:
  ```bash
  # Use hey
  hey -z 60s -c 50 http://iot-backend:8080/metrics
  
  # OR use ab
  ab -n 10000 -c 100 http://iot-backend:8080/metrics
  ```
- [ ] Monitor scaling in real-time:
  ```bash
  watch -n 2 'kubectl get hpa && kubectl get pods'
  ```
- [ ] Verify: Pod count harus increase dari 1 ke 2-5

### E. Interview Questions to Prepare
- [ ] "Kapan HPA melakukan scaling?"
- [ ] "Bedanya HPA v1 vs v2?"
- [ ] "Kenapa pakai CPU 70% sebagai threshold?"
- [ ] "Apa itu cooldown period?"

---

## 🚀 UPGRADE 3: Ingress Nginx + Path Routing
**Priority: MEDIUM** - Basic networking yang sering ditanya

> User sering bertanya: "Kenapa pakai Ingress? Kenapa tidak NodePort saja?"
> Jawaban: Ingress menghemat IP, support path-based routing (/api, /web), TLS termination, dll.

### A. Install Ingress Nginx
- [ ] K3s sudah include Traefik, tapi Ingress Nginx lebih umum di interview:
  ```bash
  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml
  ```
- [ ] Verify: `kubectl get pods -n ingress-nginx`

### B. Create Ingress Manifest
- [ ] Buat `templates/ingress.yaml`:
  ```yaml
  apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    name: iot-ingress
    annotations:
      nginx.ingress.kubernetes.io/rewrite-target: /
  spec:
    ingressClassName: nginx
    rules:
    - host: iot.local
      http:
        paths:
        - path: /api
          pathType: Prefix
          backend:
            service:
              name: iot-backend
              port:
                number: 8080
        - path: /grafana
          pathType: Prefix
          backend:
            service:
              name: grafana
              port:
                number: 3000
  ```
- [ ] Apply: `kubectl apply -f templates/ingress.yaml`
- [ ] Add to `/etc/hosts`: `127.0.0.1 iot.local`
- [ ] Test: `curl http://iot.local/api/metrics`

### C. Bandingkan dengan Service Types (Untuk Interview)
- [ ] Pahami beda:
  | Type | Use Case | Port | IP |
  |------|----------|------|-----|
  | ClusterIP | Internal only | Cluster-internal | No external IP |
  | NodePort | Dev/Testing | 30000-32767 | Node IP |
  | LoadBalancer | Cloud prod | Any | Cloud LB |
  | Ingress | Prod + Routing | 80/443 | Single IP, multiple paths |

---

## 🚀 UPGRADE 4: GitOps with ArgoCD
**Priority: LOW** - Nice-to-have, turunkan prioritas

### A. Setup ArgoCD di K3s
- [ ] Buat namespace `argocd`
- [ ] Install ArgoCD via kubectl atau Helm
- [ ] Patch ArgoCD server service → NodePort/LoadBalancer
- [ ] Get initial admin password
- [ ] Login via CLI atau Web UI

### B. Create ArgoCD Application
- [ ] Buat `argocd-app.yaml` manifest
- [ ] Konfigurasi:
  - [ ] `source`: repository Git + path `iot-chart/`
  - [ ] `destination`: cluster K3s + namespace `default`
  - [ ] `syncPolicy`: automated with auto-sync
- [ ] Apply `argocd-app.yaml` ke cluster

### C. Setup Image Updater (Optional)
- [ ] Install ArgoCD Image Updater
- [ ] Konfigurasi untuk auto-update image tag di Helm values
- [ ] Test: push code baru → ArgoCD auto-sync → cluster updated

### D. GitOps Bootstrap Pattern (Optional)
- [ ] Buat folder `argocd/` di repository
- [ ] Buat `Application.yaml` yang di-commit di Git
- [ ] ArgoCD sync dari Git itself

---

## 📋 Post-Upgrade Verification

### Functional Tests
- [ ] `curl localhost:8080/metrics` → returns Prometheus format
- [ ] Prometheus UI → Targets → IoT sensor discovered
- [ ] Grafana → Dashboard → Temperature data visible
- [ ] **HPA Test**: `hey -z 60s -c 50 http://backend:8080/metrics` → pods scale up
- [ ] **Ingress Test**: `curl http://iot.local/api` → returns 200
- [ ] ArgoCD UI → Application → Status "Synced" & "Healthy" (optional)

### Documentation
- [ ] Update README.md dengan screenshot baru
- [ ] Dokumentasi hasil load test HPA

---

## 📊 Status Tracking

| Upgrade | Status | Catatan |
|---------|--------|---------|
| #1 App Observability | PENDING | Foundation, wajib duluan |
| #2 HPA + Load Test | PENDING | Wajib hands-on practice |
| #3 Ingress Nginx | PENDING | Basic networking |
| #4 GitOps ArgoCD | PENDING | Nice-to-have |

---

## 📝 Notes

- **Upgrade #1 adalah fondasi** - harus selesai sebelum #2, #3, #4
- **HPA (#2) adalah prioritas interview** - Wajib hands-on, bukan teori
- **Estimasi waktu total**: 8-12 jam
- **Branch strategy**: 
  - Branch `feature/prometheus-metrics` untuk Upgrade #1
  - Branch `feature/hpa-loadtest` untuk Upgrade #2
- **Interview Prep**: Setelah HPA selesai, praktikkan menjawab:
  - "Demonstrate HPA scaling dengan load test"
  - "Bedanya Ingress dengan NodePort?"
  - "Kenapa pakai 70% CPU threshold?"

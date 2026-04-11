# IoT-K8s Learning Notes

## Helm Chart Structure (`iot-chart/`)

```
iot-chart/
├── Chart.yaml              # Metadata chart
├── values.yaml             # Konfigurasi default (yang bisa di-override)
└── templates/              # Template Kubernetes manifests
    ├── deployment.yaml     # Pod/Container definition
    ├── service.yaml        # Network expose
    ├── configmap-dashboard.yaml  # Grafana dashboard config
    └── servicemonitor.yaml # Prometheus monitoring config
```

---

### 1. Chart.yaml (Metadata)
```yaml
name: iot-chart           # Nama chart
version: 0.1.0           # Versi chart (Helm versioning)
appVersion: "1.16.0"     # Versi aplikasi yang di-deploy
```
**Fungsi:** Metadata chart untuk tracking oleh Helm.

---

### 2. values.yaml (Konfigurasi Default)
```yaml
image:
  repository: jawaracode/iot-dummy
  tag: v6
service:
  type: ClusterIP
  port: 8080
```
**Fungsi:** Nilai default yang bisa di-override saat `helm install`. Membuat chart reusable untuk environment berbeda.

---

### 3. templates/deployment.yaml (Pod Specification)
```yaml
containers:
- name: dummy-sensor
  image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
  ports:
    - containerPort: 8080
```
**Fungsi:** Mendefinisikan pod - container mana yang jalan, image apa, port apa.

---

### 4. templates/service.yaml (Network)
```yaml
spec:
  type: ClusterIP
  ports:
    - port: 8080
      targetPort: metrics
```
**Fungsi:** Mengekspose deployment ke jaringan internal cluster.

---

### 5. templates/configmap-dashboard.yaml (Grafana Config)
```yaml
data:
  iot-sensor-dashboard.json: |-
    {{ .Files.Get "dashboards/dashboard.json" | indent 4 }}
```
**Fungsi:** Mount file dashboard Grafana ke namespace monitoring.

---

### 6. templates/servicemonitor.yaml (Prometheus Scraping)
```yaml
kind: ServiceMonitor
spec:
  endpoints:
  - port: metrics
    path: /metrics
```
**Fungsi:** Memberitahu Prometheus untuk scrape metrics dari endpoint `/metrics`.

---

## Service Type vs Protocol

### type: ClusterIP (Service Type)
Cara service diekspose ke jaringan:

| Type | Akses |
|------|-------|
| `ClusterIP` | Hanya dari dalam cluster |
| `NodePort` | via IP node + port (30000-32767) |
| `LoadBalancer` | via cloud LB (eksternal) |
| `Headless` | langsung ke pod (no cluster IP) |

### protocol: TCP (Protocol Layer)
Protokol komunikasi:

| Protocol | Fungsi |
|----------|--------|
| `TCP` | Reliable connection, most common |
| `UDP` | Connectionless, untuk streaming/DNS |
| `SCTP` | Stream control transmission |

---

### Hubungan ClusterIP vs TCP?

Mereka saling melengkapi tapi tidak bergantung:

```yaml
spec:
  type: ClusterIP          # Service type (eksposur jaringan)
  ports:
    - port: 8080            # Port yang diekspose service
      targetPort: metrics   # Port di container
      protocol: TCP         # Protokol komunikasi
```

- `ClusterIP` → **siapa yang bisa akses** service
- `TCP` → **bagaimana data ditransmisikan**

Tidak harus TCP - kalau aplikasi pakai UDP, bisa set `protocol: UDP`. Untuk aplikasi web/HTTP/MQTT, TCP adalah default.

---

## Helm Workflow

1. `helm template` atau `helm install` → Helm render semua template
2. Template menggunakan `{{ .Values.xxx }}` → substitusi dari `values.yaml`
3. Hasilnya → Kubernetes manifests (Deployment, Service, dll)
4. Kubernetes apply manifests → Pods running

---

## CI/CD Pipeline Notes

### CI (Continuous Integration) - docker-build.yml
- Build Docker image
- Push ke Docker Hub
- Trigger: push ke branch `main`, perubahan di `backend/**` atau `.github/workflows/**`

### CD (Continuous Deployment) - Opsi Best Practice

| Approach | Best Untuk |
|----------|-----------|
| **ArgoCD/Flux (GitOps)** | Production, multi-environment |
| **OCI Registry** | Simpel, chart versioning |
| **CI Direct Deploy** | MVP/simpel, single environment |

**Recommended: GitOps dengan ArgoCD/Flux**
- Server pull dari Git repo secara otomatis
- Tidak perlu remote execution dari CI
- Changes di-commit ke Git → ArgoCD mendeteksi & deploy

---

## Quick Commands

```bash
# Render template (preview)
helm template ./iot-chart

# Install/Upgrade
helm upgrade --install iot-sensor ./iot-chart -n monitoring

# Package chart
helm package ./iot-chart

# List releases
helm list -n monitoring

# Uninstall
helm uninstall iot-sensor -n monitoring
```

package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ===================== METRICS =====================
var (
	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_temperature_celsius", Help: "Current temperature in Celsius"},
		[]string{"sensor_id"},
	)

	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_humidity_percent", Help: "Current relative humidity percentage"},
		[]string{"sensor_id"},
	)

	batteryGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_battery_level_percent", Help: "Remaining battery level"},
		[]string{"sensor_id"},
	)

	signalGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_signal_strength_dbm", Help: "Signal strength in dBm"},
		[]string{"sensor_id"},
	)

	totalReadingsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "iot_total_readings", Help: "Total readings processed"},
		[]string{"sensor_id"},
	)

	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "iot_gateway_errors_total", Help: "Total gateway errors"},
	)
)

func init() {
	prometheus.MustRegister(
		tempGauge,
		humidityGauge,
		batteryGauge,
		signalGauge,
		totalReadingsCounter,
		errorCounter,
	)
}

// ===================== MODEL =====================
type SensorData struct {
	ID        string  `json:"sensor_id"`
	Temp      float64 `json:"temperature_c"`
	Humidity  float64 `json:"humidity_pct"`
	Battery   int     `json:"battery_pct"`
	Signal    int     `json:"signal_dbm"`
	Timestamp string  `json:"timestamp"`
}

// ===================== SENSOR LOOP =====================
func sensorLoop(ctx context.Context) {
	podName, _ := os.Hostname()
	log.Printf(`{"level":"info","msg":"IoT Gateway Started","pod":"%s"}`, podName)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println(`{"level":"info","msg":"Shutting down sensor loop"}`)
			return

		case <-ticker.C:
			// Simulasi error
			if rand.Float32() < 0.05 {
				errorCounter.Inc()
				log.Println(`{"level":"error","msg":"Simulated gateway timeout"}`)
				continue
			}

			for i := 1; i <= 5; i++ {
				sensorID := "Sensor-" + string(rune(i+'0'))

				temp := math.Round((20+rand.Float64()*15)*100) / 100
				hum := math.Round((40+rand.Float64()*50)*100) / 100
				batt := rand.Intn(101)
				sig := -100 + rand.Intn(70)

				// Update metrics
				tempGauge.WithLabelValues(sensorID).Set(temp)
				humidityGauge.WithLabelValues(sensorID).Set(hum)
				batteryGauge.WithLabelValues(sensorID).Set(float64(batt))
				signalGauge.WithLabelValues(sensorID).Set(float64(sig))
				totalReadingsCounter.WithLabelValues(sensorID).Inc()

				data := SensorData{
					ID:        sensorID,
					Temp:      temp,
					Humidity:  hum,
					Battery:   batt,
					Signal:    sig,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				jsonData, _ := json.Marshal(data)
				log.Println(string(jsonData))
			}
		}
	}
}

// ===================== HEALTH CHECK =====================
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ===================== MAIN =====================
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Handle shutdown signal (Kubernetes friendly)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println(`{"level":"info","msg":"Shutdown signal received"}`)
		cancel()
	}()

	// Start sensor loop
	go sensorLoop(ctx)

	// HTTP endpoints
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println(`{"level":"info","msg":"HTTP server started on :8080"}`)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for shutdown
	<-ctx.Done()

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	server.Shutdown(ctxTimeout)
	log.Println(`{"level":"info","msg":"Server gracefully stopped"}`)
}
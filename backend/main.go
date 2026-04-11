package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Gunakan GaugeVec untuk mendukung multiple sensor_id
	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_temperature_current",
			Help: "Current temperature reading from the IoT sensor",
		},
		[]string{"sensor_id"},
	)
	
	// Gunakan CounterVec agar kita tahu total pembacaan per masing-masing sensor
	totalReadingsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iot_temperature_total_readings",
			Help: "Total number of temperature readings processed",
		},
		[]string{"sensor_id"},
	)
	
	// Error counter cukup global saja (menandakan gateway-nya yang gagal baca)
	errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "iot_sensor_read_errors_total",
		Help: "Total number of sensor read errors on the gateway",
	})
)

func init() {
	prometheus.MustRegister(tempGauge)
	prometheus.MustRegister(totalReadingsCounter)
	prometheus.MustRegister(errorCounter)
}

type SensorData struct {
	ID        string  `json:"sensor_id"`
	Temp      float64 `json:"temperature"`
	Timestamp string  `json:"timestamp"`
}

func sensorLoop() {
	// Ambil nama Pod
	podName, err := os.Hostname()
	if err != nil {
		podName = "unknown-gateway"
	}

	fmt.Printf("Starting IoT Dummy Sensor Gateway on Pod: %s...\n", podName)

	for {
		// Simulasi error sesekali (5% kemungkinan gagal di tingkat Gateway)
		if rand.Float32() < 0.05 {
			errorCounter.Inc()
			fmt.Println("{\"error\": \"Simulated sensor read error on gateway timeout\"}")
			time.Sleep(2 * time.Second)
			continue // Lanjut ke putaran berikutnya
		}

		// 1 Pod mensimulasikan 5 sensor sekaligus
		for i := 1; i <= 5; i++ {
			// Bikin ID unik. Contoh: SENSOR-x8hlf-1
			sensorID := fmt.Sprintf("SENSOR-%s-%d", podName, i)

			// Generate suhu & bulatkan 2 desimal
			rawTemp := 20.0 + rand.Float64()*(35.0-20.0)
			fixedTemp := math.Round(rawTemp*100) / 100

			// Update Metrics
			tempGauge.WithLabelValues(sensorID).Set(fixedTemp)
			totalReadingsCounter.WithLabelValues(sensorID).Inc()

			// Bungkus ke JSON struct
			data := SensorData{
				ID:        sensorID,
				Temp:      fixedTemp,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			jsonData, _ := json.Marshal(data)
			fmt.Println(string(jsonData))
		}

		// Jeda sebelum gateway membaca data 5 sensor itu lagi
		time.Sleep(5 * time.Second)
	}
}

func main() {
	go sensorLoop()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Serving metrics on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
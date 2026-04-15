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
	// METRIK 1: Suhu
	tempGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_temperature_celsius", Help: "Current temperature in Celsius"},
		[]string{"sensor_id"},
	)
	// METRIK 2: Kelembapan
	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_humidity_percent", Help: "Current relative humidity percentage"},
		[]string{"sensor_id"},
	)
	// METRIK 3: Baterai
	batteryGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_battery_level_percent", Help: "Remaining battery level"},
		[]string{"sensor_id"},
	)
	// METRIK 4: Sinyal WiFi/LoRa (RSSI)
	signalGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: "iot_signal_strength_dbm", Help: "Signal strength in dBm"},
		[]string{"sensor_id"},
	)
	
	totalReadingsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "iot_total_readings", Help: "Total readings processed"},
		[]string{"sensor_id"},
	)
	errorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "iot_gateway_errors_total", Help: "Total gateway simulation errors"},
	)
)

func init() {
	prometheus.MustRegister(tempGauge, humidityGauge, batteryGauge, signalGauge, totalReadingsCounter, errorCounter)
}

type SensorData struct {
	ID        string  `json:"sensor_id"`
	Temp      float64 `json:"temperature_c"`
	Humidity  float64 `json:"humidity_pct"`
	Battery   int     `json:"battery_pct"`
	Signal    int     `json:"signal_dbm"`
	Timestamp string  `json:"timestamp"`
}

func sensorLoop() {
	podName, err := os.Hostname()
	if err != nil {
		podName = "unknown-gateway"
	}

	fmt.Printf(" Starting Enterprise IoT Gateway on Pod: %s...\n", podName)

	for {
		if rand.Float32() < 0.05 {
			errorCounter.Inc()
			fmt.Println("{\"level\":\"error\", \"message\":\"Simulated gateway timeout fetching node data\"}")
			time.Sleep(5 * time.Second)
			continue
		}

		for i := 1; i <= 5; i++ {
			sensorID := fmt.Sprintf("NODE-%s-%d", podName, i) // Ubah nama dari SENSOR ke NODE agar lebih industri

			// Generate Data Kompleks
			temp := math.Round((20.0+rand.Float64()*(35.0-20.0))*100) / 100
			hum := math.Round((40.0+rand.Float64()*(90.0-40.0))*100) / 100
			batt := rand.Intn(101) // 0 - 100
			sig := -100 + rand.Intn(71) // -100 sampai -30

			// Set Metrics
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
			fmt.Println(string(jsonData))
		}
		time.Sleep(30 * time.Second)
	}
}

func main() {
	go sensorLoop()
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
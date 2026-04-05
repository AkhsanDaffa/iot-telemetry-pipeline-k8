package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tempGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iot_temperature_current",
		Help: "Current temperature reading from the IoT sensor",
	})	
	totalReadingsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "iot_temperature_total_readings",
		Help: "Total number of temperature readings processed",
	})
	errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "iot_sensor_read_errors_total",
		Help: "Total number of sensor read errors",
	})
)

func init() {
	prometheus.MustRegister(tempGauge)
	prometheus.MustRegister(totalReadingsCounter)
	prometheus.MustRegister(errorCounter)
}

type SensorData struct {
	ID			string	`json:"sensor_id"`
	Temp		float64 `json:"temperature"`
	Timestamp 	string 	`json:"timestamp"`
}

func sensorLoop() {
	for {
		if rand.Float32() < 0.05 {
			errorCounter.Inc()
			fmt.Println("Simulated sensor read error")
			time.Sleep(1 * time.Second)
			continue
		}

			fmt.Println("Starting IoT Dummy Sensor...")

	for {
		rawTemp := 20.0 + rand.Float64()*(35.0-20.0)
		fixedTemp := math.Round(rawTemp*100) / 100

		tempGauge.Set(fixedTemp)
		totalReadingsCounter.Inc()

		data := SensorData{
				ID:			"SENSOR-001",
				Temp:		fixedTemp,
				Timestamp: 	time.Now().Format(time.RFC3339),
		}

		jsonData, _ := json.Marshal(data)
		fmt.Println(string(jsonData))

		time.Sleep(1 * time.Second)
	}
	}
}

func main() {
	fmt.Println("Starting IoT Dummy Sensor with Prometheus Metrics...")

	go sensorLoop()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Serving metrics on port 8080...")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
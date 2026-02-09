package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type SensorData struct {
	ID			string	`json:"sensor_id"`
	Temp		float64 `json:"temperature"`
	Timestamp 	string 	`json:"timestamp"`
}

func main() {
	fmt.Println("Starting IoT Dummy Sensor...")

	for {
		rawTemp := 20.0 + rand.Float64()*(35.0-20.0)

		fixedTemp := math.Round(rawTemp*100) / 100

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
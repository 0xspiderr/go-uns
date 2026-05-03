package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type OTData struct {
	Temperature float64 `json:"temperature"`
	IsRunning   bool    `json:"is_running"`
	MotorRPM    int     `json:"motor_rpm"`
}

func startOTPublisher(brokerURL string) {
	options := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("plc_01")
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to start OT publisher: %v", token.Error())
	}

	topic := "assembly_line/plc_01/data"
	fmt.Printf("OT publishing to %s\n", topic)

	// publish mock data every 5 seconds
	for {
		rawTemp := 45.0 + (rand.Float64() * 4.0)
		roundedTemp := float64(int(rawTemp*10)) / 10
		data := OTData{
			Temperature: roundedTemp,
			IsRunning:   true,
			MotorRPM:    1000 + rand.Intn(50),
		}
		payload, _ := json.Marshal(data)
		client.Publish(topic, 0, false, payload).Wait()
		time.Sleep(5 * time.Second)
	}
}

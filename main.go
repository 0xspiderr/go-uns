package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// unified namespace data from the IT and OT level
type UNSData struct {
	TimeStamp string `json:"timestamp"`
	Topic     string `json:"topic"`
	OT        OTData `json:"ot_data"`
	IT        ITData `json:"it_data"`
}

func mqttMsgHandler(client mqtt.Client, msg mqtt.Message) {
	var ot OTData
	json.Unmarshal(msg.Payload(), &ot)
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %v\n", ot.Temperature)
}

func startUNSListener(brokerURL string) {
	options := mqtt.NewClientOptions().AddBroker(brokerURL)
	client := mqtt.NewClient(options)
	// connect the client to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to start UNS listener: %v", token.Error())
	}

	if token := client.Subscribe("#", 0, mqttMsgHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("UNS listener started successfully")
}

func main() {
	const brokerURL string = "tcp://localhost:1883"
	go startERPServer()
	go startOTPublisher(brokerURL)
	startUNSListener(brokerURL)

	// block the main routine to let the program listen for messages
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
}

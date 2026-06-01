package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

// unified namespace data from the IT and OT level
type UNSData struct {
	TimeStamp string          `json:"timestamp"`
	Topic     string          `json:"topic"`
	OT        json.RawMessage `json:"ot_data"`
	IT        ITData          `json:"it_data"`
}

func mqttMsgHandler(client mqtt.Client, msg mqtt.Message) {
	// OLD MOCK OT DATA
	// var ot OTData
	//if err := json.Unmarshal(msg.Payload(), &ot); err != nil {
	//	log.Printf("Error decoding OT data: %v", err)
	//	return
	// }

	log.Printf("New ot message received:\n")
	log.Printf("Topic: %s\n", msg.Topic())
	log.Printf("Raw payload: %s\n", string(msg.Payload()))

	it, err := fetchITData()
	if err != nil {
		log.Printf("Error fetching IT data: %v", err)
		return
	}

	uns := UNSData{
		Topic:     msg.Topic(),
		TimeStamp: time.Now().Format(time.StampMilli), // placeholder format for now
		OT:        msg.Payload(),
		IT:        it,
	}

	insertUNSData(uns)
	uns_json, _ := json.MarshalIndent(uns, "\t", "")
	fmt.Printf("UNS payload: %s\n", string(uns_json))
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	brokerURL := os.Getenv("MQTT_BROKER_URL")

	configureDB()
	// go startITServer()
	// go startOTPublisher(brokerURL)
	startUNSListener(brokerURL)
	// block the main routine to let the program listen for messages
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
}

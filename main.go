package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// unified namespace data from the IT and OT level
type UNSData struct {
	TimeStamp string `json:"timestamp"`
	Topic     string `json:"topic"`
	OT        OTData `json:"ot_data"`
	IT        ITData `json:"it_data"`
}

func fetchITData() (ITData, error) {
	var it ITData
	response, err := http.Get("http://localhost:8089/api/orders")
	if err != nil {
		return it, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&it)
	return it, err
}

func mqttMsgHandler(client mqtt.Client, msg mqtt.Message) {
	var ot OTData
	if err := json.Unmarshal(msg.Payload(), &ot); err != nil {
		log.Printf("Error decoding OT data: %v", err)
		return
	}

	it, err := fetchITData()
	if err != nil {
		log.Printf("Error fetching IT data: %v", err)
		return
	}

	uns := UNSData{
		Topic:     msg.Topic(),
		TimeStamp: time.Now().Format(time.StampMilli), // placeholder format for now
		OT:        ot,
		IT:        it,
	}

	uns_json, _ := json.MarshalIndent(uns, "\t", "")
	fmt.Printf("UNS payload: %s\n", string(uns_json))
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %v\n", msg.Payload())
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

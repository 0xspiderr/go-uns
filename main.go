package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func mqttMsgHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	const brokerURL string = "tcp://localhost:1883"
	client_options := mqtt.NewClientOptions()
	client_options.AddBroker(brokerURL)
	client := mqtt.NewClient(client_options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := client.Subscribe("go-mqtt/test", 0, mqttMsgHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Println("subscribed succesfully")

	// block the main routine here to let the program listen for messages
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
}

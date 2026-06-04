package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/0xspiderr/go-uns/internal/db"
	"github.com/0xspiderr/go-uns/internal/mqtt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load .env file")
	}

	brokerURL := os.Getenv("MQTT_BROKER_URL")

	db.ConfigureDB()

	var wg sync.WaitGroup
	// start with two workers because there are 2 PLCs on the OT level
	mqtt.InitWorkers(2, &wg)

	mqttClient := mqtt.StartListener(brokerURL)
	// block the main routine to let the program listen for messages
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	log.Println("shutdown signal received, disconnecting mqtt client")
	mqttClient.Disconnect(1000)
	mqtt.CloseChannel()

	log.Println("waiting for workers to finish processing pending messages...")
	wg.Wait()

	log.Println("closing database...")
	db.Close()
}

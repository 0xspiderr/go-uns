package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
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

var otMessageChannel = make(chan mqtt.Message, 100)

func unsWorker(workerID int, wg *sync.WaitGroup) {
	// ensure that when the channel is closed, the loop exits and the wwaitgroup is decremented
	defer wg.Done()

	for msg := range otMessageChannel {
		log.Printf("[worker %d] processing topic: %s\n", workerID, msg.Topic())

		it, err := fetchITData()
		// log that the ERP is down, but insert the OT data anyway.
		if err != nil {
			log.Printf("[worker %d] error fetching IT data: %v\n", workerID, err)
		}

		uns := UNSData{
			Topic:     msg.Topic(),
			TimeStamp: time.Now().Format(time.RFC3339),
			OT:        msg.Payload(),
			IT:        it,
		}

		err = insertUNSData(uns)
		if err != nil {
			log.Printf("[worker %d] failed to insert UNS data: %v\n", workerID, err)
		}
	}

	log.Printf("[worker %d] shutdown\n", workerID)
}

func mqttMsgHandler(client mqtt.Client, msg mqtt.Message) {
	// OLD MOCK OT DATA
	// var ot OTData
	//if err := json.Unmarshal(msg.Payload(), &ot); err != nil {
	//	log.Printf("Error decoding OT data: %v", err)
	//	return
	// }

	//log.Printf("New ot message received:\n")
	//log.Printf("Topic: %s\n", msg.Topic())
	//log.Printf("Raw payload: %s\n", string(msg.Payload()))

	//it, err := fetchITData()
	//if err != nil {
	//	log.Printf("Error fetching IT data: %v", err)
	//	return
	//}

	//uns := UNSData{
	//	Topic:     msg.Topic(),
	//	TimeStamp: time.Now().Format(time.StampMilli), // placeholder format for now
	//	OT:        msg.Payload(),
	//	IT:        it,
	//}

	//insertUNSData(uns)
	//uns_json, _ := json.MarshalIndent(uns, "\t", "")
	//fmt.Printf("UNS payload: %s\n", string(uns_json))
	select {
	case otMessageChannel <- msg:
	default:
		log.Printf("message channel full!")

	}
}

func startUNSListener(brokerURL string) mqtt.Client {
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
	return client
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

	var wg sync.WaitGroup

	// start go routines for uns workers
	numWorkers := 2
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go unsWorker(i, &wg)
	}

	mqttClient := startUNSListener(brokerURL)
	// block the main routine to let the program listen for messages
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	log.Println("shutdown signal received, disconnecting mqtt client")
	mqttClient.Disconnect(1000)
	close(otMessageChannel)

	log.Println("waiting for workers to finish processing pending messages...")
	wg.Wait()

	log.Println("closing database...")
	if db != nil {
		db.Close()

	}
}

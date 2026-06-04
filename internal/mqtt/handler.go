package mqtt

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/0xspiderr/go-uns/internal/db"
	"github.com/0xspiderr/go-uns/internal/erp"
	"github.com/0xspiderr/go-uns/internal/models"

	mqttLib "github.com/eclipse/paho.mqtt.golang"
)

var otMessageChannel = make(chan mqttLib.Message, 100)

func InitWorkers(numWorkers int, wg *sync.WaitGroup) {
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go unsWorker(i, wg)
	}
}

func unsWorker(workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range otMessageChannel {
		it, err := erp.FetchITData()
		if err != nil {
			log.Printf("[Worker %d] WARNING: ERP unavailable, inserting without IT context. Error: %v", workerID, err)
		}

		uns := models.UNSData{
			Topic:     msg.Topic(),
			TimeStamp: time.Now().Format(time.RFC3339),
			OT:        msg.Payload(),
			IT:        it,
		}

		if err := db.InsertUNSData(uns); err != nil {
			log.Printf("[Worker %d] Failed to insert UNS data: %v", workerID, err)
		}
	}
	fmt.Printf("[Worker %d] Gracefully shut down.\n", workerID)
}

func mqttMsgHandler(client mqttLib.Client, msg mqttLib.Message) {
	select {
	case otMessageChannel <- msg:
	default:
		log.Printf("WARNING: Message channel full! Dropping payload from %s", msg.Topic())
	}
}

func StartListener(brokerURL string) mqttLib.Client {
	options := mqttLib.NewClientOptions().AddBroker(brokerURL)
	client := mqttLib.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to start UNS listener: %v", token.Error())
	}

	if token := client.Subscribe("#", 0, mqttMsgHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Println("UNS listener started successfully")
	return client
}

func CloseChannel() {
	close(otMessageChannel)
}

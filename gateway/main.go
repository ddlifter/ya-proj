package gateway

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json: "name"`
	Value string `json:"value`
}

var responseChannels = make(map[string]chan *sarama.ConsumerMessage)
var mu sync.Mutex

func main() {
	responseChannels = make(map[string]chan *sarama.ConsumerMessage)

	//Creating Kafka producer
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	//Creating Kafka consumer
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	//creating partition "pong"
	pastConsumer, err := consumer.ConsumePartition("pong", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
		defer pastConsumer.Close()
	}

	go func() {
		for {
			select {
			//Reading messages from Kafka
			case msg, ok := <-pastConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting goroutine")
				}
				responseID := string(msg.Key)
				mu.Lock()
				ch, exist := responseChannels[responseID]
				if exist {
					ch <- msg
					delete(responseChannels, responseID)
				}
				mu.Unlock()
			}
		}
	}()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		requestID := uuid.New().String()

		message := MyMessage{
			ID:    requestID,
			Name:  "Ping",
			Value: "Pong",
		}

		//Formatting to JSON message to send with Kafka later
		bytes, err := json.Marshal(message)
		if err != nil {
			c.JSON(500, gin.H{"erro": "failed to marshal JSON"})
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: "ping",
			Key:   sarama.StringEncoder(bytes),
		}

		//Sending messages to Kafka
		_, _, err = producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message to Kafka")
			c.JSON(500, gin.H{"error": "failed to send message to Kafka"})
			return
		}

		responseCh := make(chan *sarama.ConsumerMessage)
		mu.Lock()
		responseChannels[requestID] = responseCh
		mu.Unlock()

		select {
		case responseMsg := <-responseCh:
			c.JSON(200, gin.H{"message": string(responseMsg.Value)})
		case <-time.After(10 * time.Second):
			mu.Lock()
			delete(responseChannels, requestID)
			mu.Unlock()
			c.JSON(500, gin.H{"error": "timeout waiting for response"})
		}
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

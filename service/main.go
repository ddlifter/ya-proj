package service

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type MyMessage struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	// Создание продюсера Kafka
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Создание консьюмера Kafka
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Подписка на партицию "ping" в Kafka
	partConsumer, err := consumer.ConsumePartition("ping", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	for {
		select {
		// (обработка входящего сообщения и отправка ответа в Kafka)
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Println("Channel closed, exiting")
				return
			}

			// Десериализация входящего сообщения из JSON
			var receivedMessage MyMessage
			err := json.Unmarshal(msg.Value, &receivedMessage)

			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			log.Printf("Received message: %+v\n", receivedMessage)

			responseText := receivedMessage.Name + " " + receivedMessage.Value + " ( " + receivedMessage.ID + " ) "

			// Формируем ответное сообщение
			resp := &sarama.ProducerMessage{
				Topic: "pong",
				Key:   sarama.StringEncoder(receivedMessage.ID),
				Value: sarama.StringEncoder(responseText),
			}
			// Отпровляем ответ в gateway
			_, _, err = producer.SendMessage(resp)
			if err != nil {
				log.Printf("Failed to send message to Kafka: %v", err)
			}
		}
	}
}

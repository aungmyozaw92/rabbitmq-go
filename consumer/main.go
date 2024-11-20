package main

import (
	"log"
	"sync"

	"github.com/aungmyozaw92/rabbitmq-go"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "TestingQueue"
)

func main() {
	// Connect to RabbitMQ
	conn := rabbitmq.Connect(rabbitMQURL)
	defer conn.Close()

	// Create a channel
	ch := rabbitmq.CreateChannel(conn)
	defer ch.Close()

	// Declare a queue
	q := rabbitmq.DeclareQueue(ch, queueName)

	// Consume messages
	msgs, err := ch.Consume(
		q.Name,
		"",    // Consumer
		true,  // Auto-acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Arguments
	)
	rabbitmq.FailOnError(err, "Failed to register a consumer")

	// Process messages

	var processedMessages sync.Map // Thread-safe map

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	for d := range msgs {
		if _, exists := processedMessages.Load(d.MessageId); exists {
			log.Printf("duplicate message detected: %s", d.MessageId)
			continue
		}
		processedMessages.Store(d.MessageId, true)
		log.Printf("Received a message: %s", d.Body)
		log.Printf("Processing message: %s", d.Body)

	}
}

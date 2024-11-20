package main

import (
	"context"
	"log"
	"time"

	"github.com/aungmyozaw92/rabbitmq-go"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
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

	// Set a timeout context for publishing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish a message
	body := "Hello RabbitMQ!"
	messageId := uuid.New().String() // Generate unique ID

	err := ch.PublishWithContext(ctx,
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			MessageId:   messageId, // Unique identifier

		})
	rabbitmq.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n %s\n", body, messageId)
}

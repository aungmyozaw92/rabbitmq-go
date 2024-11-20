package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErr(err error, msg string){
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

const queueName = "TestingQueue"


func main(){
	// connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")
	defer ch.Close()


	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // Queue name
		false,   // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	failOnErr(err, "Failed to declare a queue")

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer
		true,   // Auto-acknowledge
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	failOnErr(err, "Failed to register a consumer")

	// Process messages
	forever := make(chan bool)
	
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
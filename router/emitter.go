package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ROUTER_INPUT_EXCHANGE = "router_input"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error connecting to RabitMQ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error opening a channel", err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		ROUTER_INPUT_EXCHANGE, // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring the input exchange", err)
	}

	data := Data{"123456789", "Hello World"}

	body, _ := json.Marshal(data)
	err = channel.Publish(
		ROUTER_INPUT_EXCHANGE, // exchange
		"",                    // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatalf("Error sending a message", err)
	}
	log.Printf(" [x] Sent %s\n", body)
}

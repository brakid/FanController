package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ROUTER_INPUT_EXCHANGE  = "router_input"
	ROUTER_OUTPUT_EXCHANGE = "router_output"
)

func handleMessage(channel *amqp.Channel, transformer *Transformer, content []byte) {
	var data Data
	err := json.Unmarshal(content, &data)
	if err != nil {
		log.Fatalf("Error deserializing data", err)
	}

	commands := (*transformer).TransformData(data)

	for _, command := range commands {
		commandData, _ := json.Marshal(command)
		err = channel.Publish(
			ROUTER_OUTPUT_EXCHANGE, // exchange
			command.TargetId,       // routing key
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        commandData,
			},
		)
		if err != nil {
			log.Fatalf("Error sending command %v", command, err)
		}
		log.Printf(" [x] Sent %s\n", commandData)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ", err)
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

	err = channel.ExchangeDeclare(
		ROUTER_OUTPUT_EXCHANGE, // name
		"direct",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring the output exchange", err)
	}

	inputQueue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Error declaring the input queue", err)
	}

	err = channel.QueueBind(
		inputQueue.Name,       // queue name
		"",                    // routing key
		ROUTER_INPUT_EXCHANGE, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error binding the input queue to the input exchange", err)
	}

	msgs, err := channel.Consume(
		inputQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Error registering a consumer", err)
	}

	var forever chan struct{}

	var transformer Transformer = DummyTransformer{}

	go func() {
		for msg := range msgs {
			log.Printf("[x] %s", msg.Body)
			handleMessage(channel, &transformer, msg.Body)
		}
	}()

	log.Printf("[*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

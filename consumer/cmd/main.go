package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	addr  = "amqp://rabbitmq:5673"
	queue = "movies"
)

func main() {
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Printf("Cannot connect to RabbitMQ: %v\n", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Cannot open a channel to RabbitMQ: %v\n", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Cannot declare a queue in RabbitMQ: %v\n", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Cannot consume messages from RabbitMQ %v\n", err)
	}

	msgChan := make(chan movie)
	go func() {
		for msg := range msgs {
			m := msg.Body
			fmt.Printf("Received message: %s\n", string(m))
			var mov movie
			err := json.Unmarshal(m, &mov)
			if err != nil {
				log.Printf("Cannot unmarshal message from queue: %v\n", err)
			}
			msgChan <- mov
		}
	}()
	log.Printf("Started listening for messages on '%s' queue", queue)

	for m := range msgChan {
		log.Printf("Received message from chan: %v\n", m)
	}
}

type movie struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

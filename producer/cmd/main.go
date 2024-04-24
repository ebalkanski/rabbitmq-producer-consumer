package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	serviceName = "metadata_producer"
	addr        = "amqp://rabbitmq:5673"
	queue       = "movies"
)

func main() {
	log.Printf("Starting service [%s]\n", serviceName)

	ctx := context.Background()

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

	type movie struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	i := 1
	for {
		msg := movie{
			//Id:   "movie" + strconv.Itoa(rand.Intn(11)),
			Id:   "movie" + strconv.Itoa(i),
			Name: "Movie title",
		}
		encoded, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Cannot marshal message: %v\n", err)
			continue
		}

		err = ch.PublishWithContext(
			ctx,
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        encoded,
			})
		if err != nil {
			log.Printf("Cannot publish a message to RabbitMQ: %v\n", err)
			return
		}
		log.Printf("Sent message: %s\n", string(encoded))
		i++
		time.Sleep(1 * time.Second)
	}
}

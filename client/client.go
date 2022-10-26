package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var defaultMessageSize int = 1

type Customer struct {
	id int
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"bytes",  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 1000; i++ {
		fmt.Printf("Customer %d going to sushi bar\n", i)
		time.Sleep(time.Second * 1)
		err = ch.PublishWithContext(ctx,
			"bytes", // exchange
			"",      // routing key
			false,   // mandatory
			false,   // immediate
			amqp.Publishing{
				Timestamp: time.Now(),
				Body:      []byte(strconv.Itoa(i)),
			})
		failOnError(err, "Failed to publish a message")
	}
}

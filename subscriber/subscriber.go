package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
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

	q, err := ch.QueueDeclare(
		"bytes", // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,  // queue name
		"",      // routing key
		"bytes", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//table
	isTableFull := false
	tableSeatsAvailable := 5
	for msg := range msgs {
		// when table reaches size 5
		// for loop should be blocked while people are eating sushi
		customerId := string(msg.Body)
		fmt.Printf("Customer %s arrived\n", customerId)
		fmt.Printf("Customer %s waiting\n", customerId)
		for isTableFull {
		}

		if tableSeatsAvailable > 0 {
			fmt.Printf("Customer %s sitting\n", customerId)
			tableSeatsAvailable--
		}

		if tableSeatsAvailable == 0 {
			isTableFull = true
			go eating(&isTableFull, &tableSeatsAvailable)
		}
	}
}

func eating(isTableFull *bool, tableSeatsAvailable *int) {
	fmt.Println("Five friends eating sushi")
	time.Sleep(time.Second * 5)
	*isTableFull = false
	*tableSeatsAvailable = 5
}

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
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

	prefetchCount, err := strconv.Atoi(os.Args[1])
	failOnError(err, "Failed to get prefetchCount")

	ch.Qos(prefetchCount, 0, false)

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
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	csvFile, err := os.Create("execs.csv")
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write([]string{"Bytes", "ExecTime"})
	msgCounter := 0
	for msg := range msgs {
		err = msg.Ack(false)
		msgCounter++
		failOnError(err, "Failed to ack")
		arrivalTime := time.Now().UnixNano()
		failOnError(err, "Failed to receive a message")
		err = csvwriter.Write([]string{fmt.Sprintf("%d", len(msg.Body)), fmt.Sprintf("%d", arrivalTime-msg.Timestamp.UnixNano())})
		if msgCounter%100 == 0 {
			csvwriter.Flush()
		}
		failOnError(err, "Could not write to csv")
	}
	csvFile.Close()
}

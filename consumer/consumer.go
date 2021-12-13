package main

import (
	"fmt"
	logs "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_username = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {
	consume()
}

func consume() {
	conn, err := amqp.Dial("amqp://" + rabbit_username + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	if err != nil {
		logs.Fatal("%s: %s", "failed to connect rabbitmq", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		logs.Fatal("%s: %s", "failed to connect channel", err)
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"publisher",
		false,
		false,
		false,
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
		)

	if err != nil {
		logs.Fatal("%s: %s", "failed to connect consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logs.Printf("message received: %s", d.Body)

			d.Ack(false)
		}
	}()

	fmt.Println("running.......")
	<-forever
}
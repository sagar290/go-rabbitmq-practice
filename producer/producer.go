package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	logs "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"net/http"
	"os"
)

var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_username = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")

func main() {

	router := httprouter.New()

	router.GET("/publish/:message", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		submit(w, r, p)
	})

	fmt.Println("running......")
	logs.Fatal(http.ListenAndServe(":80", router))
}

func submit(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	message := p.ByName("message")

	fmt.Println("Received message: " + message)

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

	if err != nil {
		logs.Fatal("%s: %s", "failed to published message", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		logs.Fatal("%s: %s", "failed to published message", err)
	}

	fmt.Println("published successfully")
}

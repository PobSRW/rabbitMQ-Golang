package main

import (
	"context"
	"log"
	"obptwentyeight/evendrivenrabbitmq/internal"
	"os"
	"strings"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

var username = "guest"
var password = "guest"
var host = "localhost:5672"
var vHost = "/"

func main() {
	conn, err := internal.ConnectRabbitMQ(username, password, host, vHost)
	if err != nil {
		panic(err)
	}

	// close connection
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}

	// close channel
	defer client.Close()

	if err := client.CreateExchange("logs", "fanout", true, false); err != nil {
		panic(err)
	}

	body := bodyFrom(os.Args)

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := client.Send(ctx, "logs", "", ampq.Publishing{
		// we need to mark our messages as persistent - by using the amqp.Persistent option amqp.Publishing takes
		DeliveryMode: ampq.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	}); err != nil {
		panic(err)
	}

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

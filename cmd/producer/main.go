package main

import (
	"context"
	"log"
	"obptwentyeight/evendrivenrabbitmq/internal"
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

	if err := client.CreateQueue("hello_world", false, false); err != nil {
		panic(err)
	}

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := client.Send(ctx, "", "hello_world", ampq.Publishing{
		ContentType: "text/plain",
		Body:        []byte(`An cool message between services`),
	}); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	log.Println(client, "client")
}

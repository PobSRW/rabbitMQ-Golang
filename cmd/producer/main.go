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

	if err := client.CreateQueue("customer_created", true, false); err != nil {
		panic(err)
	}
	if err := client.CreateQueue("customer_test", false, true); err != nil {
		panic(err)
	}

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := client.Send(ctx, "", "customer_created", ampq.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: ampq.Persistent,
		Body:         []byte(`An cool message between services`),
	}); err != nil {
		panic(err)
	}

	if err := client.Send(ctx, "", "customer_test", ampq.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: ampq.Transient,
		Body:         []byte(`An uncool durable message`),
	}); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)
	log.Println(client, "client")
}

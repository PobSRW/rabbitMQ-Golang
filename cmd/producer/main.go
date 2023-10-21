package main

import (
	"log"
	"obptwentyeight/evendrivenrabbitmq/internal"
	"time"
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

	time.Sleep(10 * time.Second)
	log.Println(client, "client")
}

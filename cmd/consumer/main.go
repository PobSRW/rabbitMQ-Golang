package main

import (
	"log"
	"obptwentyeight/evendrivenrabbitmq/internal"
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

	messageBus, err := client.Consume("hello_world", "", true)
	if err != nil {
		panic(err)
	}

	blocking := make(chan struct{})

	go func() {
		for message := range messageBus {
			log.Printf("New Message: %v\n", message)
		}
	}()

	<-blocking

}

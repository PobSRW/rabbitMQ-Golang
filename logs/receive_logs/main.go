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

	if err := client.CreateExchange("logs", "fanout", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateExclusiveQueue("", false, false); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("", "", "logs"); err != nil {
		panic(err)
	}

	messageBus, err := client.Consume("", "", true)
	if err != nil {
		panic(err)
	}

	blocking := make(chan struct{})

	go func() {
		for message := range messageBus {
			log.Printf(" [x] %s", message.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocking

}

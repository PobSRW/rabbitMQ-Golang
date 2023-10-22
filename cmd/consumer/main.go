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

	messageBus, err := client.Consume("customer_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}

	go func() {
		for message := range messageBus {
			log.Println("New Message: %v", message)

			if err := message.Ack(false); err != nil {
				log.Println("Acknowledge message failes")
				continue
			}

			log.Printf("Acknowledge message %s\n", message.MessageId)
		}
	}()

	log.Println("Consuming, to close the program press CTRL+C")

	<-blocking

}

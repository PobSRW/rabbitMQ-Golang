package main

import (
	"bytes"
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

	if err := client.CreateQueue("test_queue", true, false); err != nil {
		panic(err)
	}

	if err := client.FairDispatch(1, 0); err != nil {
		panic(err)
	}

	messageBus, err := client.Consume("test_queue", "", false)
	if err != nil {
		panic(err)
	}

	blocking := make(chan struct{})

	go func() {
		for message := range messageBus {
			log.Printf("Received a message: %s", message.Body)

			dotCount := bytes.Count(message.Body, []byte("."))

			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			message.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocking

}

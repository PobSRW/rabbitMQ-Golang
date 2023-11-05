package internal

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// The connection used by the client
	conn *amqp.Connection
	// Channel is used to process / Send messages
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vHost string) (*amqp.Connection, error) {
	url := fmt.Sprintf(`amqp://%[1]s:%[2]s@%[3]s/%[4]s`, username, password, host, vHost)
	return amqp.Dial(url)
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	// close only channel, we don't wanna close the connection
	// incase others people using the same connection

	return rc.ch.Close()
}

// CreateQueue will create a new queue based on given configurations
func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

func (rc RabbitClient) CreateExclusiveQueue(queueName string, durable, autodelete bool) error {
	// When the connection that declared it closes, the queue will be deleted because it is declared as exclusive.
	// for broadcast
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, true, false, nil)
	return err
}

// CreateBinding will bind the current channel to the given exchange using the routingkey provided
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving nowait false, having nowait set to false will make the channel return en error if its false to bind
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

// Send is used to publish payloads onto an exchange with the given routingKey
func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	return rc.ch.PublishWithContext(ctx,
		exchange,
		routingKey,
		// Mandatory is used to determine if an error should be returned upon failure
		true,
		// Immediate is deprecated, so leave it false
		false,
		options,
	)
}

// Consume is used to consume a queue
func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, // queue name
		consumer, // consumer
		autoAck,
		// exclusive: true => the one and only customer consume that queue, if you have one customer to consume set it to true
		false, // false: the server will distribute messages using a load balance
		false, // noLocal => not support => set it to false
		false, // noWait
		nil,   // args
	)
}

func (rc RabbitClient) FairDispatch(prefetchCount, prefetchSize int) error {
	// set the prefetch count with the value of 1. This tells RabbitMQ not to give more than one message to a worker at a time. Or, in other words, don't dispatch a new message to a worker until it has processed and acknowledged the previous one. Instead, it will dispatch it to the next worker that is not still busy.
	return rc.ch.Qos(prefetchCount, prefetchSize, false)
}

func (rc RabbitClient) CreateExchange(name, kind string, duration, autoDelete bool) error {
	return rc.ch.ExchangeDeclare(name,
		kind, // type
		duration, autoDelete, false, false, nil)
}

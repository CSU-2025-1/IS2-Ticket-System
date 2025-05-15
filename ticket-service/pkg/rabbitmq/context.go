package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type Context[T any] struct {
	Delivery *amqp.Delivery
	Message  T
}

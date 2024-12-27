package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Connect() error
	Publish(queue string, message interface{}, headers map[string]interface{}) error
	Consume(queue string, consumer string, handler func(message amqp.Delivery, headers map[string]interface{})) error
	RequestReply(queue string, replyTo string, message interface{}, correlationID string, headers map[string]interface{}) (*amqp.Delivery, error)
	Close() error
}

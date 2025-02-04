package broker

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
	"otus/pkg/logger"
	"sync"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	URL        string
	log        slog.Logger
	mu         sync.Mutex
}

func NewRabbitMQ() *RabbitMQ {
	dbUser := os.Getenv("RABBITMQ_USER")
	dbPassword := os.Getenv("RABBITMQ_PASSWORD")
	dbHost := os.Getenv("RABBITMQ_HOST")
	return &RabbitMQ{
		URL: fmt.Sprintf("amqp://%s:%s@%s/", dbUser, dbPassword, dbHost),
		log: *logger.GetLogger(),
	}
}

func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(r.URL)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err.Error()))
		return err
	}
	r.Connection = conn
	err = r.openChannel()
	if err != nil {
		return err
	}
	r.log.Info("Connected to RabbitMQ")
	return nil
}

func (r *RabbitMQ) Publish(queue string, message interface{}, headers map[string]interface{}) error {
	err := r.openChannel()
	if err != nil {
		return err
	}
	body, err := json.Marshal(message)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to marshal message: %v", err.Error()))
		return err
	}
	q, err := r.declareQueue(queue)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to declare queue: %v", err.Error()))
		return err
	}
	amqpHeaders := amqp.Table{}
	for key, value := range headers {
		amqpHeaders[key] = value
	}
	err = r.Channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body:        body,
			Headers:     amqpHeaders,
			ContentType: "application/json",
		},
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to publish message: %s", err.Error()))
		return err
	}
	r.log.Info(fmt.Sprintf("Message published to queue: %s", queue))
	return nil
}

func (r *RabbitMQ) Consume(queue string, consumer string, handler func(message amqp.Delivery, headers map[string]interface{})) error {
	err := r.openChannel()
	if err != nil {
		return err
	}
	_, err = r.declareQueue(queue)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to declare queue: %s", err))
		return err
	}
	msgs, err := r.Channel.Consume(queue, consumer, true, false, false, false, nil)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to consume messages: %v", err))
		return err
	}
	go func() {
		for msg := range msgs {
			headers := make(map[string]interface{})
			for key, value := range msg.Headers {
				headers[key] = value
			}
			r.log.Info(fmt.Sprintf("Message consumed from queue: %s", queue))
			ctx := context.Background()
			go func(ctx context.Context, message amqp.Delivery) {
				handler(message, headers)
			}(ctx, msg)
		}
	}()
	return nil
}

func (r *RabbitMQ) RequestReply(queue string, replyTo string, message interface{}, correlationID string, headers map[string]interface{}) (*amqp.Delivery, error) {
	ch, err := r.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	defer ch.Close()
	msgs, err := ch.Consume(replyTo, "", true, false, false, false, nil)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to consume messages: %v", err))
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to declare queue: %s", err))
		return nil, err
	}
	amqpHeaders := amqp.Table{}
	for key, value := range headers {
		amqpHeaders[key] = value
	}
	body, err := json.Marshal(message)
	if err != nil {
		r.log.Error("Failed to marshal message: %v", err)
		return nil, err
	}
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body:          body,
			Headers:       amqpHeaders,
			CorrelationId: correlationID,
			ReplyTo:       replyTo,
			ContentType:   "application/json",
		},
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to publish message: %s", err))
		return nil, err
	}
	r.log.Info(fmt.Sprintf("Message published to queue: %s", queue))

	for msg := range msgs {
		headers := make(map[string]interface{})
		for key, value := range msg.Headers {
			headers[key] = value
		}
		return &msg, nil
	}
	return nil, nil
}

func (r *RabbitMQ) Close() error {
	err := r.Channel.Close()
	if err != nil {
		return err
	}
	return r.Connection.Close()
}

func (r *RabbitMQ) openChannel() error {
	if r.Channel == nil || r.Channel.IsClosed() {
		ch, err := r.Connection.Channel()
		if err != nil {
			return fmt.Errorf("failed to create channel: %w", err)
		}
		r.Channel = ch
	}
	return nil
}

func (r *RabbitMQ) declareQueue(queue string) (amqp.Queue, error) {
	q, err := r.Channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	return q, err
}

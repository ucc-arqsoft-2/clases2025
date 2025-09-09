package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

const (
	encodingJSON = "application/json"
	encodingUTF8 = "UTF-8"
)

type RabbitMQClient struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      *amqp091.Queue
}

func NewRabbitMQClient(user string, password string, queueName string, host string, port string) *RabbitMQClient {
	connection, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port))
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %w", err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %w", err)
	}
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare a queue: %w", err)
	}
	return &RabbitMQClient{
		connection: connection,
		channel:    channel,
		queue:      &queue,
	}
}

func (r RabbitMQClient) Publish(ctx context.Context, action string, itemID string) error {
	message := map[string]interface{}{
		"action":  action,
		"item_id": itemID,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message to JSON: %w", err)
	}

	if err := r.channel.PublishWithContext(ctx, "", r.queue.Name, false, false, amqp091.Publishing{
		ContentType:     encodingJSON,
		ContentEncoding: encodingUTF8,
		DeliveryMode:    amqp091.Transient,
		MessageId:       uuid.New().String(),
		Timestamp:       time.Now().UTC(),
		AppId:           "items-api",
		Body:            bytes,
	}); err != nil {
		return fmt.Errorf("error publishing message to RabbitMQ: %w", err)
	}
	return nil
}

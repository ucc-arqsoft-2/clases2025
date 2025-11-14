package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
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

// Reintenta conectar (hasta 10 intentos) y NO hace log.Fatalf.
func NewRabbitMQClient(user, password, queueName, host, port string) (*RabbitMQClient, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)

	var (
		conn *amqp091.Connection
		ch   *amqp091.Channel
		q    amqp091.Queue
		err  error
	)

	for i := 1; i <= 10; i++ {
		conn, err = amqp091.Dial(dsn)
		if err == nil {
			ch, err = conn.Channel()
		}
		if err == nil {
			// durable=true para sobrevivir reinicios
			q, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
		}
		if err == nil {
			log.Printf("✅ RabbitMQ conectado (intento %d), queue=%s", i, queueName)
			return &RabbitMQClient{connection: conn, channel: ch, queue: &q}, nil
		}
		wait := time.Duration(i) * time.Second
		log.Printf("⏳ Rabbit no disponible (intento %d): %v. Reintento en %s...", i, err, wait)
		time.Sleep(wait)
	}
	return nil, fmt.Errorf("no pude conectar a RabbitMQ: %w", err)
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil { _ = r.channel.Close() }
	if r.connection != nil { _ = r.connection.Close() }
}

func (r *RabbitMQClient) Publish(ctx context.Context, action, itemID string) error {
	msg := map[string]any{"action": action, "item_id": itemID}
	body, err := json.Marshal(msg)
	if err != nil { return fmt.Errorf("marshal: %w", err) }

	// mandatory=true para detectar “unroutable” si la cola no existe
	return r.channel.PublishWithContext(
		ctx,
		"", r.queue.Name, true, false,
		amqp091.Publishing{
			ContentType:     encodingJSON,
			ContentEncoding: encodingUTF8,
			DeliveryMode:    amqp091.Persistent, // persistente
			MessageId:       uuid.New().String(),
			Timestamp:       time.Now().UTC(),
			AppId:           "items-api",
			Body:            body,
		},
	)
}

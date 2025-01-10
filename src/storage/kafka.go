package storage

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/shyam0507/pd-payment/src/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	storage Storage
}

// ProducePaymentReceivedEvent implements Producer.
func (k KafkaProducer) ProducePaymentReceivedEvent(key string, value types.PaymentReceivedEvent) error {

	m, _ := json.Marshal(value)
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: m,
		},
	)

	if err != nil {
		slog.Error("failed to write messages:", "Err", err)
		return err
	}

	slog.Info("Payment Received event published for Key", "Key", key)

	return nil
}

func NewKafkaProducer(topic string, brokers []string) Producer {
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: mechanism,
		},
		AllowAutoTopicCreation: true,
	}

	return KafkaProducer{writer: w}
}

func NewKafkaConsumer(topic string, brokers []string, storage Storage) Consumer {

	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     "pd-payment-group",
		StartOffset: kafka.LastOffset,
		MaxBytes:    10e6, // 10MB
		Dialer:      dialer,
	})

	return KafkaConsumer{reader: r, storage: storage}
}

// ConsumeOrderCreated implements Consumer.
//
// It will block until a message is consumed, then it will return the value of the message.
// If the consumer is closed, it will return an error.
func (k KafkaConsumer) ConsumeOrderCreated() error {
	for {
		m, err := k.reader.ReadMessage(context.Background())

		if err != nil {
			slog.Error("failed to read message:", "Err", err)
			continue
		}

		slog.Info("kafka event consumed for Key", "Key", string(m.Key))
		var order types.OrderCreatedEvent
		err = json.Unmarshal(m.Value, &order)

		if err != nil {
			slog.Error("failed to unmarshal message:", "Err", err)
			continue
		}

		slog.Info("Order created event consumed", "Order", order)

		payment := types.Payment{
			Id:         primitive.NewObjectID(),
			OrderId:    order.Data.Id.Hex(),
			Total:      order.Data.Total,
			CustomerId: order.Data.CustomerId,
			Status:     "PENDING",
		}

		if err := k.storage.CreatePayment(payment); err != nil {
			slog.Error("Error while creating the payment", "Err", err)
			continue
		}

	}

}

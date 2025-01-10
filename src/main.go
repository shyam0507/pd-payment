package main

import (
	"os"

	"github.com/shyam0507/pd-payment/src/api"
	"github.com/shyam0507/pd-payment/src/storage"
)

func main() {

	kafkaBrokers := os.Getenv("KAFKA_BROKERS") //localhost:9092
	mongoURI := os.Getenv("MONGO_URI")         //mongodb://localhost:27017

	mongoStorage := storage.NewMongoStorage(mongoURI, "pd_payments")
	consumer := storage.NewKafkaConsumer("order.created", []string{kafkaBrokers}, mongoStorage)
	producer := storage.NewKafkaProducer("payment.received", []string{kafkaBrokers})

	go consumer.ConsumeOrderCreated()

	server := api.NewServer("3006", mongoStorage, producer)
	server.Start()
}

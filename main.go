package main

import (
	"github.com/shyam0507/pd-payment/api"
	"github.com/shyam0507/pd-payment/storage"
)

func main() {

	mongoStorage := storage.NewMongoStorage("mongodb://localhost:27017", "pd_payments")
	consumer := storage.NewKafkaConsumer("order.created", []string{"localhost:9092"}, mongoStorage)
	producer := storage.NewKafkaProducer("payment.received", []string{"localhost:9092"})

	go consumer.ConsumeOrderCreated()

	server := api.NewServer("3006", mongoStorage, producer)
	server.Start()
}

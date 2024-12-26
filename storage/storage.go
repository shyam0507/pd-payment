package storage

import "github.com/shyam0507/pd-payment/types"

type Storage interface {
	FindPayment(order string) (types.Payment, error)
	UpdatePayment(id string, status string) error
	CreatePayment(types.Payment) error
}

type Producer interface {
	ProducePaymentReceivedEvent(string, types.PaymentReceivedEvent) error
}

type Consumer interface {
	ConsumeOrderCreated() error
}

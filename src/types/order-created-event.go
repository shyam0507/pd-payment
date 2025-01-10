package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderCreatedEvent struct {
	SpecVersion     string `json:"specversion"`
	Type            string `json:"type"`
	Source          string `json:"source"`
	Subject         string `json:"subject"`
	Id              string `json:"id"`
	Time            string `json:"time"`
	DataContentType string `json:"datacontenttype"`
	Data            Order  `json:"data"`
}

type Order struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Products   []any              `json:"products" bson:"products"`
	CustomerId string             `json:"customer_id" bson:"customer_id"`
	AddressId  string             `json:"address_id" bson:"address_id"`
	Status     string             `json:"status" bson:"status"`
	Total      float64            `json:"total,omitempty" bson:"total"`
}

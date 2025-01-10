package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	Id          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	CustomerId  string             `json:"customer_id" bson:"customer_id"`
	OrderId     string             `json:"order_id" bson:"order_id"`
	Status      string             `json:"status" bson:"status"`
	Total       float64            `json:"total,omitempty" bson:"total"`
	ReferenceId string             `json:"reference_id" bson:"reference_id"` //TODO use a nested struct in future
}

func (p *Payment) Validate() error {
	if p.Total > 0 {
		return fmt.Errorf("total should be greater than 0")
	}

	return nil
}

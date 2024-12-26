package storage

import (
	"context"
	"time"

	"github.com/shyam0507/pd-payment/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	client *mongo.Client
	dbName string
}

const Collection_Name = "payment"

// CreatePayment implements Storage.
func (m *MongoStorage) CreatePayment(o types.Payment) error {

	m.client.Database(m.dbName).Collection(Collection_Name).InsertOne(
		context.Background(),
		o,
	)

	return nil
}

// UpdatePayment implements Storage.
func (m *MongoStorage) UpdatePayment(id string, status string) error {
	oId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}
	_, err = m.client.Database(m.dbName).Collection(Collection_Name).UpdateOne(
		context.Background(),
		bson.M{"_id": oId},
		bson.M{"$set": bson.M{"status": status}},
	)

	return err
}

func (m *MongoStorage) FindPayment(id string) (types.Payment, error) {

	data := m.client.Database(m.dbName).Collection(Collection_Name).FindOne(
		context.Background(),
		bson.M{"order_id": id},
	)

	if data.Err() != nil {
		return types.Payment{}, data.Err()
	}

	var payment types.Payment

	if err := data.Decode(&payment); err != nil {
		return types.Payment{}, err
	}

	return payment, nil
}

func NewMongoStorage(uri string, dbName string) *MongoStorage {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return nil
	}

	return &MongoStorage{client: client, dbName: dbName}
}

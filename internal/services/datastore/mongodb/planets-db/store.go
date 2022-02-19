package planetsdb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Querier
}

type MongoDBStore struct {
	mongodbClient *mongo.Client
}

func NewStore(mongodbClient *mongo.Client) MongoDBStore {
	return MongoDBStore{
		mongodbClient: mongodbClient,
	}
}

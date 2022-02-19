package planetsdb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"testing"
)

const (
	mongoURI = "mongodb://localhost:27017"
)

var testStore MongoDBStore

func TestMain(m *testing.M) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalln("could not connect to database:", err)
	}
	testStore = NewStore(client)
	os.Exit(m.Run())
}

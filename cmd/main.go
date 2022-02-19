package main

import (
	"context"
	planetsfactory "github.com/gmaschi/b2w-sw-planets/internal/factories/planets-factory"
	planetsdb "github.com/gmaschi/b2w-sw-planets/internal/services/datastore/mongodb/planets-db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln("could not connect to database:", err)
	}
	store := planetsdb.NewStore(client)
	server, err := planetsfactory.New(&store)
	if err != nil {
		log.Fatalln("could not start server:", err)
	}

	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatalln("cannot start server", err)
	}
}

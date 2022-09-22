package data

import (
	"alex-api/utils"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetClient() *mongo.Client {
	uri := utils.GetEnvVar("MONGO_CONNECTION_URI")

	if uri == "" {
		log.Fatal("$MONGO_CONNECTION_URI must be set")
	}

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func DisconnectClient(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

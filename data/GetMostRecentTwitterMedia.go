package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMostRecentTwitterMedia() (*TwitterMedia, error) {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("takoland").Collection("twitter-media")
	var twitterMedia TwitterMedia
	var opts = &options.FindOneOptions{
		Sort: bson.D{primitive.E{Key: "created_at", Value: -1}},
	}

	filter := bson.D{}
	if err := collection.FindOne(context.Background(), filter, opts).Decode(&twitterMedia); err != nil {
		return nil, err
	}

	return &twitterMedia, nil
}

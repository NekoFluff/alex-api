package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTwitterMedia(skip *int64, limit *int64) ([]TwitterMedia, error) {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("takoland").Collection("twitter-media")
	opts := &options.FindOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  bson.D{primitive.E{Key: "created_at", Value: -1}},
	}
	cur, err := collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []TwitterMedia
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

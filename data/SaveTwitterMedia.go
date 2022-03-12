package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveTwitterMedia(twitterMedia TwitterMedia) *mongo.UpdateResult {
	client := GetClient()
	defer DisconnectClient(client)

	collection := client.Database("takoland").Collection("twitter-media")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "url", Value: twitterMedia.Url}, primitive.E{Key: "tweet_id", Value: twitterMedia.TweetId}}
	update := bson.D{primitive.E{Key: "$set", Value: twitterMedia}}

	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Update Result: %#v\n", result)
	return result
}

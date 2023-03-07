package data

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TwitterAuthor struct {
	Id       string `bson:"id"`
	Name     string `bson:"name"`
	UserName string `bson:"username"`
}

type TwitterMedia struct {
	Author            TwitterAuthor `bson:"author"`
	TweetId           string        `bson:"tweet_id"`
	Url               string        `bson:"url"`
	Updated           time.Time     `bson:"updated"`
	CreatedAt         time.Time     `bson:"created_at"`
	PossiblySensitive bool          `bson:"possibly_sensitive"`
	Width             int16         `bson:"width"`
	Height            int16         `bson:"height"`
}

func (db *DB) GetTwitterMedia(skip *int64, limit *int64) ([]TwitterMedia, error) {
	collection := db.client.Database("takoland").Collection("twitter-media")
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

func (db *DB) GetMostRecentTwitterMedia() (*TwitterMedia, error) {
	collection := db.client.Database("takoland").Collection("twitter-media")
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

func (db *DB) CreateTwitterMedia(twitterMedia TwitterMedia) *mongo.UpdateResult {
	collection := db.client.Database("takoland").Collection("twitter-media")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "url", Value: twitterMedia.Url}, primitive.E{Key: "tweet_id", Value: twitterMedia.TweetId}}
	update := bson.D{primitive.E{Key: "$set", Value: twitterMedia}}

	result, err := collection.UpdateOne(ctx, filter, update, options)

	if err != nil {
		db.log.WithFields(logrus.Fields{
			"update": update,
		}).Error("Failed to upsert twitter media:", err)
		return nil
	}

	db.log.WithFields(logrus.Fields{
		"result": result,
	}).Info("Created Twitter Media")
	return result
}

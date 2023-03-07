package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	Name             string       `bson:"Name,omitempty"`
	Action           string       `bson:"Action,omitempty"`
	Image            string       `bson:"Image,omitempty"`
	Recipe           []Ingredient `bson:"Recipe,omitempty"`
	QuantityProduced float64      `bson:"Quantity Produced,omitempty"`
	MinProduced      float64      `bson:"Min Produced,omitempty"`
	MaxProduced      float64      `bson:"Max Produced,omitempty"`
	TimeToProduce    float64      `bson:"Time to Produce,omitempty"`
	MarketData       *MarketData  `bson:"Market Data,omitempty"`
}

type MarketData struct {
	LastUpdateAttempt time.Time `bson:"Last Update Attempt,omitempty"`
	LastUpdated       time.Time `bson:"Last Updated,omitempty"`
	Price             float64   `bson:"Price,omitempty"`
	Quantity          float64   `bson:"Quantity,omitempty"`
	TotalTradeCount   float64   `bson:"Total Trade Count,omitempty"`
	Name              string    `bson:"Name,omitempty"`
}

type Ingredient struct {
	ItemName string  `bson:"Item Name,omitempty"`
	Amount   float64 `bson:"Amount,omitempty"`
}

func (db *DB) GetRecipes(skip *int64, limit *int64) ([]Recipe, error) {
	collection := db.client.Database("bdo-craft-profit").Collection("recipesWithMarketPrice")
	opts := &options.FindOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  bson.D{primitive.E{Key: "Name", Value: -1}},
	}
	cur, err := collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []Recipe
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *DB) GetPPSRecipes(skip *int64, limit *int64) ([]Recipe, error) {
	collection := db.client.Database("bdo-craft-profit").Collection("recipes")
	opts := &options.FindOptions{
		Skip:  skip,
		Limit: limit,
		Sort:  bson.D{primitive.E{Key: "Name", Value: -1}},
	}
	cur, err := collection.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []Recipe
	if err = cur.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

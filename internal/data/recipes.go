package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	Name             string      `json:"name" bson:"name"`
	Action           string      `json:"action,omitempty" bson:"action,omitempty"`
	Facility         string      `json:"facility,omitempty" bson:"facility,omitempty"`
	Image            string      `json:"image,omitempty" bson:"image,omitempty"`
	Ingredients      Ingredients `json:"ingredients" bson:"ingredients"`
	QuantityProduced float64     `json:"quantityProduced,omitempty" bson:"quantityProduced,omitempty"`
	MinProduced      float64     `json:"minProduced,omitempty" bson:"minProduced,omitempty"`
	MaxProduced      float64     `json:"maxProduced,omitempty" bson:"maxProduced,omitempty"`
	TimeToProduce    float64     `json:"timeToProduce,omitempty" bson:"timeToProduce,omitempty"`
	MarketData       *MarketData `json:"marketData,omitempty" bson:"marketData,omitempty"`
}

type Ingredients map[string]float64

type MarketData struct {
	LastUpdateAttempt time.Time `json:"lastUpdateAttempt,omitempty" bson:"lastUpdateAttempt,omitempty"`
	LastUpdated       time.Time `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
	Price             float64   `json:"price,omitempty" bson:"price,omitempty"`
	Quantity          float64   `json:"quantity,omitempty" bson:"quantity,omitempty"`
	TotalTradeCount   float64   `json:"totalTradeCount,omitempty" bson:"totalTradeCount,omitempty"`
	Name              string    `json:"name,omitempty" bson:"name,omitempty"`
}

func (db *DB) GetBDORecipes(skip *int64, limit *int64) ([]Recipe, error) {
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

func (db *DB) GetBDOPPSRecipes(skip *int64, limit *int64) ([]Recipe, error) {
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

func (db *DB) GetDSPRecipes(skip *int64, limit *int64) ([]Recipe, error) {
	collection := db.client.Database("dsp").Collection("recipes")
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

func (db *DB) CreateDSPRecipes(recipes []Recipe) error {
	collection := db.client.Database("dsp").Collection("recipes")
	iRecipes := make([]interface{}, len(recipes))
	for i := range recipes {
		iRecipes[i] = recipes[i]
	}

	_, err := collection.InsertMany(context.Background(), iRecipes)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteDSPRecipes() error {
	collection := db.client.Database("dsp").Collection("recipes")
	_, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	return nil
}

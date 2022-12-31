package data

import (
	"alex-api/utils"
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	log    *logrus.Logger
	client *mongo.Client
}

func New(l *logrus.Logger) *DB {
	db := &DB{
		log: l,
	}

	db.Connect()
	return db
}

func (db *DB) Connect() *mongo.Client {
	uri := utils.GetEnvVar("MONGO_CONNECTION_URI")

	if uri == "" {
		db.log.Fatal("$MONGO_CONNECTION_URI must be set")
	}

	// Create a new client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client
}

func (db *DB) Disconnect(client *mongo.Client) {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

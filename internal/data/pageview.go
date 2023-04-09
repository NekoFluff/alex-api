package data

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PageView struct {
	Domain       string    `json:"domain" bson:"domain"`
	Path         string    `json:"path" bson:"path"`
	TimesTracked int64     `json:"timesTracked" bson:"timesTracked"`
	LastTracked  time.Time `json:"lastTracked" bson:"lastTracked"`
}

func (pageView *PageView) Increment() *PageView {
	pageView.TimesTracked++
	pageView.LastTracked = time.Now()
	return pageView
}

func (db *DB) GetPageView(domain string, path string) (PageView, error) {
	collection := db.client.Database("analytics").Collection("page-views")
	filter := bson.D{primitive.E{Key: "domain", Value: domain}, primitive.E{Key: "path", Value: path}}

	result := collection.FindOne(context.Background(), filter, nil)
	var pageView PageView
	err := result.Decode(&pageView)
	if err != nil {
		db.log.WithFields(logrus.Fields{
			"domain": domain,
			"path":   path,
		}).Error("Failed to deocde Page View", err)
	}
	return pageView, err
}

func (db *DB) CreatePageView(pageView PageView) error {
	collection := db.client.Database("analytics").Collection("page-views")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, pageView)

	if err != nil {
		db.log.WithFields(logrus.Fields{
			"PageView": pageView,
		}).Error("Failed to insert Page View:", err)
		return err
	}

	db.log.WithFields(logrus.Fields{
		"result": result,
	}).Info("Created Page View")
	return nil
}

func (db *DB) UpdatePageView(pageView PageView) error {
	collection := db.client.Database("analytics").Collection("page-views")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	filter := bson.D{primitive.E{Key: "domain", Value: pageView.Domain}, primitive.E{Key: "path", Value: pageView.Path}}
	update := bson.D{primitive.E{Key: "$set", Value: pageView}}

	result, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		db.log.WithFields(logrus.Fields{
			"update": update,
		}).Error("Failed to update Page View:", err)
		return err
	}

	db.log.WithFields(logrus.Fields{
		"result": result,
	}).Info("Updated Page View")
	return nil
}

package main

import (
	"alex-api/internal/data"
	"alex-api/internal/dspscraper"
	"context"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	dspRecipes := dspscraper.Scrape()
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	log := logger.WithContext(ctx)

	db := data.New(log)
	_ = db.DeleteDSPRecipes()
	_ = db.CreateDSPRecipes(dspRecipes)
	cancel()
}

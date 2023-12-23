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
	err := db.DeleteDSPRecipes()
	if err != nil {
		log.WithError(err).Error("Failed to delete DSP recipes")
	}
	err = db.CreateDSPRecipes(dspRecipes)
	if err != nil {
		log.WithError(err).Error("Failed to create DSP recipes")
	}
	cancel()
}

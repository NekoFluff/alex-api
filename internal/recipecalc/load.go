package recipecalc

import (
	"alex-api/internal/data"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func LoadDSPRecipes(file string) map[string][]Recipe {
	recipeMap := map[string][]Recipe{}

	// Open up the file
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err, "failed to open recipes file")
	}
	defer jsonFile.Close()

	// Read and unmarshal the file
	byteValue, _ := io.ReadAll(jsonFile)
	var recipes []Recipe
	err = json.Unmarshal(byteValue, &recipes)
	if err != nil {
		log.Fatal(err, "failed to unmarshal recipes from file")
	}

	// Map the recipe
	for _, recipe := range recipes {
		name := strings.ToLower(string(recipe.OutputItem))
		recipeMap[name] = append(recipeMap[name], recipe)
	}

	return recipeMap
}

func LoadBDORecipes() map[string][]Recipe {
	log := logrus.New().WithContext(context.TODO())
	recipeMap := map[string][]Recipe{}

	// TEMP BDO Section
	db := data.New(log)
	defer db.Disconnect()
	dbRecipes, err := db.GetRecipes(nil, nil)
	if err != nil {
		log.Fatal("failed to load recipes", err)
	}
	log.Print("recipes", dbRecipes)

	// Map the recipe
	for _, recipe := range dbRecipes {
		name := strings.ToLower(string(recipe.Name))

		mats := Materials{}
		for _, ingredient := range recipe.Recipe {
			mats[ingredient.ItemName] = ingredient.Amount
		}

		marketData := &MarketData{}
		if recipe.MarketData != nil {
			marketData = &MarketData{
				LastUpdateAttempt: recipe.MarketData.LastUpdateAttempt,
				LastUpdated:       recipe.MarketData.LastUpdated,
				Price:             recipe.MarketData.Price,
				Quantity:          recipe.MarketData.Quantity,
				TotalTradeCount:   recipe.MarketData.TotalTradeCount,
				Name:              recipe.MarketData.Name,
			}
		}

		recipeMap[name] = append(recipeMap[name], Recipe{
			OutputItem:         recipe.Name,
			OutputItemCount:    recipe.QuantityProduced,
			MinOutputItemCount: recipe.MinProduced,
			MaxOutputItemCount: recipe.MaxProduced,
			Facility:           recipe.Action,
			Time:               recipe.TimeToProduce,
			Materials:          mats,
			Image:              recipe.Image,
			MarketData:         marketData,
		})
	}

	return recipeMap
}

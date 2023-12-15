package recipecalc

import (
	"alex-api/internal/data"
	"strings"

	"github.com/sirupsen/logrus"
)

func LoadDSPRecipes(log *logrus.Entry, db DB) map[string][]data.Recipe {
	recipeMap := map[string][]data.Recipe{}

	dbRecipes, err := db.GetDSPRecipes(nil, nil)
	if err != nil {
		log.Fatal("failed to load recipes", err)
	}
	log.Print("recipes", dbRecipes)

	// Map the recipe
	for _, recipe := range dbRecipes {
		name := strings.ToLower(string(recipe.Name))
		recipeMap[name] = append(recipeMap[name], recipe)
	}

	return recipeMap
}

func LoadBDORecipes(log *logrus.Entry, db DB) map[string][]data.Recipe {
	recipeMap := map[string][]data.Recipe{}

	dbRecipes, err := db.GetBDORecipes(nil, nil)
	if err != nil {
		log.Fatal("failed to load recipes", err)
	}
	log.Print("recipes", dbRecipes)

	// Map the recipe
	for _, recipe := range dbRecipes {
		name := strings.ToLower(string(recipe.Name))
		recipeMap[name] = append(recipeMap[name], recipe)
	}

	return recipeMap
}

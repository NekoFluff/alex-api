package recipecalc

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOptimizer_DoesNotInfinitelyLoop(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	recipeMap := map[string][]Recipe{}

	// Open up the file
	jsonFile, err := os.Open("../data/loopy_items.json")
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

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(recipeMap)

	recipe := o.GetOptimalRecipe("Loopy item", 1, "", map[string]bool{}, 1, map[string]int{})
	assert.Equal(t, []ComputedRecipe{
		{
			OutputItem:          "Loopy item",
			Facility:            "Smelting facility",
			NumFacilitiesNeeded: 1,
			ItemsConsumedPerSec: map[string]float64{
				"Loopy item": 1,
			},
			SecondsSpentPerCraft: 1,
			CraftingPerSec:       1,
			UsedFor:              "",
			Depth:                1,
		},
	}, recipe)
}

func TestOptimizer_E2E_ConveyorBeltMKII(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes("../data/items.json"))

	expectedRecipes := []ComputedRecipe{}
	f, err := os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)
	o.SortRecipes(expectedRecipes)

	recipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{})
	o.SortRecipes(recipes)
	for i, recipe := range recipes {
		recipe.Image = ""
		recipes[i] = recipe
	}

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], recipes[k])
	}
}

func TestOptimizer_E2E_ConveyorBeltMKII_Combined(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes("../data/items.json"))

	expectedRecipes := []ComputedRecipe{}
	f, err := os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2 combined.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)
	o.SortRecipes(expectedRecipes)

	recipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{})
	o.SortRecipes(recipes)
	recipes = o.CombineRecipes(recipes)
	o.SortRecipes(recipes)
	for i, recipe := range recipes {
		recipe.Image = ""
		recipes[i] = recipe
	}

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], recipes[k])
	}
}

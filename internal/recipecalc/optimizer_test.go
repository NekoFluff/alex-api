package recipecalc

import (
	"alex-api/internal/data"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOptimizer_DoesNotInfinitelyLoop(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	recipeMap := map[string][]data.Recipe{}

	// Open up the file
	jsonFile, err := os.Open("test_data/loopy_items.json")
	if err != nil {
		log.Fatal(err, "failed to open recipes file")
	}
	defer jsonFile.Close()

	// Read and unmarshal the file
	byteValue, _ := io.ReadAll(jsonFile)
	var recipes []data.Recipe
	err = json.Unmarshal(byteValue, &recipes)
	if err != nil {
		log.Fatal(err, "failed to unmarshal recipes from file")
	}

	// Map the recipe
	for _, recipe := range recipes {
		name := strings.ToLower(string(recipe.Name))
		recipeMap[name] = append(recipeMap[name], recipe)
	}

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(recipeMap)

	recipe := o.GetOptimalRecipe("Loopy item", 1, "", map[string]bool{}, 1, map[string]int{}, 2, 1, 1)
	assert.Equal(t, []ComputedRecipe{
		{
			Name:                "Loopy item",
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

	recipes := []data.Recipe{}
	f, err := os.ReadFile("../data/items.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &recipes)
	assert.Equal(t, nil, err)

	expectedRecipes := []ComputedRecipe{}
	f, err = os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)

	ctrl := gomock.NewController(t)
	dbMock := NewMockDB(ctrl)
	dbMock.EXPECT().GetDSPRecipes(nil, nil).Return(recipes, nil)

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes(log, dbMock))
	o.SortRecipes(expectedRecipes)

	optimalRecipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{}, 2, 1, 1)
	o.SortRecipes(optimalRecipes)
	for i, recipe := range optimalRecipes {
		recipe.Image = ""
		optimalRecipes[i] = recipe
	}

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], optimalRecipes[k])
	}
}

func TestOptimizer_E2E_ConveyorBeltMKII_Combined(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	recipes := []data.Recipe{}
	f, err := os.ReadFile("../data/items.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &recipes)
	assert.Equal(t, nil, err)

	expectedRecipes := []ComputedRecipe{}
	f, err = os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2_combined.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)

	ctrl := gomock.NewController(t)
	dbMock := NewMockDB(ctrl)
	dbMock.EXPECT().GetDSPRecipes(nil, nil).Return(recipes, nil)

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes(log, dbMock))
	o.SortRecipes(expectedRecipes)

	optimalRecipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{}, 2, 1, 1)
	o.SortRecipes(optimalRecipes)
	optimalRecipes = o.CombineRecipes(optimalRecipes)
	o.SortRecipes(optimalRecipes)
	for i, recipe := range optimalRecipes {
		recipe.Image = ""
		optimalRecipes[i] = recipe
	}

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], optimalRecipes[k])
	}
}

func TestOptimizer_E2E_ConveyorBeltMKII_Combined_WithLevel3Assembler(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	recipes := []data.Recipe{}
	f, err := os.ReadFile("../data/items.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &recipes)
	assert.Equal(t, nil, err)

	expectedRecipes := []ComputedRecipe{}
	f, err = os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2_combined_level_3_assembler.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)

	ctrl := gomock.NewController(t)
	dbMock := NewMockDB(ctrl)
	dbMock.EXPECT().GetDSPRecipes(nil, nil).Return(recipes, nil)

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes(log, dbMock))
	o.SortRecipes(expectedRecipes)

	optimalRecipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{}, 3, 1, 1)
	o.SortRecipes(optimalRecipes)
	optimalRecipes = o.CombineRecipes(optimalRecipes)
	o.SortRecipes(optimalRecipes)
	for i, recipe := range optimalRecipes {
		recipe.Image = ""
		optimalRecipes[i] = recipe
	}

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], optimalRecipes[k])
	}
}

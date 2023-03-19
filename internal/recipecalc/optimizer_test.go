package recipecalc

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDSP_Optimizer_DoesNotInfinitelyLoop(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes())

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

func TestDSP_Optimizer_E2E_ConveyorBeltMKII(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes())

	expectedRecipes := []ComputedRecipe{}
	f, err := os.ReadFile("test_data/computed_recipe_conveyor_belt_mk_2.json")
	assert.Equal(t, nil, err)
	err = json.Unmarshal(f, &expectedRecipes)
	assert.Equal(t, nil, err)
	o.SortRecipes(expectedRecipes)

	recipes := o.GetOptimalRecipe("Conveyor belt MK.II", 1, "", map[string]bool{}, 1, map[string]int{})
	o.SortRecipes(recipes)

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], recipes[k])
	}
}

func TestDSP_Optimizer_E2E_ConveyorBeltMKII_Combined(t *testing.T) {
	log := logrus.New().WithContext(context.TODO())

	o := NewOptimizer(log, OptimizerConfig{})
	o.SetRecipes(LoadDSPRecipes())

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

	for k := range expectedRecipes {
		assert.Equal(t, expectedRecipes[k], recipes[k])
	}
}

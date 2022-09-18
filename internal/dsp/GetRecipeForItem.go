package dsp

import (
	"addi/models"
	"math"
)

func GetRecipeForItem(itemName string, craftingSpeed float64, parentItemName string) []*models.DSPRecipeResponse {
	recipes := []*models.DSPRecipeResponse{}
	item, ok := GetItem(itemName)

	if ok {
		consumedMats := make(map[string]float64)
		numberOfFacilitiesNeeded := 0.0
		if item.Produce != nil && item.Time != nil {
			numberOfFacilitiesNeeded = *item.Time * craftingSpeed / *item.Produce
		}

		for _, material := range item.Materials {
			if material.Name != nil && item.Time != nil {
				consumedMats[*material.Name] = *material.Count * numberOfFacilitiesNeeded / *item.Time
				if math.IsNaN(consumedMats[*material.Name]) {
					consumedMats[*material.Name] = 0
				}
			}
		}

		recipe := &models.DSPRecipeResponse{
			Produce:                  item.Name,
			MadeIn:                   item.MadeIn,
			NumberOfFacilitiesNeeded: &numberOfFacilitiesNeeded,
			ConsumesPerSec:           consumedMats,
			SecondsSpendPerCrafting:  item.Time,
			CraftingPerSecond:        &craftingSpeed,
			For:                      &parentItemName,
		}
		recipes = append(recipes, recipe)

		for _, material := range item.Materials {
			targetCraftingSpeed := 0.0
			if material.Count != nil && item.Time != nil {
				targetCraftingSpeed = numberOfFacilitiesNeeded * *material.Count / *item.Time
			}
			if math.IsNaN(targetCraftingSpeed) {
				targetCraftingSpeed = 0
			}
			if material.Name != nil && item.Name != nil {
				materialRecipe := GetRecipeForItem(*material.Name, targetCraftingSpeed, *item.Name)
				recipes = append(recipes, materialRecipe...)
			}
		}
	}

	return recipes
}

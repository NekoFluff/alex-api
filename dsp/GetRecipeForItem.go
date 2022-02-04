package dsp

import "addi/models"

func GetRecipeForItem(itemName string, craftingSpeed float64, parentItemName string) []*models.Recipe {
	recipes := []*models.Recipe{}
	item, ok := GetItem(itemName)

	if ok {
		consumedMats := make(map[string]float64)
		numberOfFacilitiesNeeded := item.Time * craftingSpeed / item.Produce

		for _, material := range item.Materials {
			consumedMats[material.Name] = material.Count * numberOfFacilitiesNeeded / item.Time
		}

		recipe := &models.Recipe{
			Produce:                  item.Name,
			MadeIn:                   item.MadeIn,
			NumberOfFacilitiesNeeded: numberOfFacilitiesNeeded,
			ConsumesPerSec:           consumedMats,
			SecondsSpendPerCrafting:  item.Time,
			CraftingPerSecond:        craftingSpeed,
			For:                      parentItemName,
		}
		recipes = append(recipes, recipe)

		for _, material := range item.Materials {
			targetCraftingSpeed := numberOfFacilitiesNeeded * material.Count / item.Time
			materialRecipe := GetRecipeForItem(material.Name, targetCraftingSpeed, item.Name)
			recipes = append(recipes, materialRecipe...)
		}
	}

	return recipes
}

package handlers

import (
	"alex-api/models"
	"alex-api/restapi/operations"
	"encoding/json"
	"fmt"

	"github.com/NekoFluff/go-dsp/dsp"
	"github.com/go-openapi/runtime/middleware"
)

func GetDSPRecipeHandler(params operations.GetDSPRecipeParams) middleware.Responder {
	recipe := []*models.DSPRecipeResponse{}

	for _, v := range params.RecipeRequest {
		itemName := dsp.ItemName(*v.Name)
		computedRecipes := optimizer.GetOptimalRecipe(itemName, float32(*v.Count), "", map[dsp.ItemName]bool{})

		for _, cr := range computedRecipes {
			cps := float64(cr.CraftingPerSec)
			usedFor := string(cr.UsedFor)
			madeIn := string(cr.Facility)
			numFacilities := float64(cr.NumFacilitiesNeeded)
			product := string(cr.OutputItem)
			sspc := float64(cr.SecondsSpentPerCraft)

			recipe = append(recipe, &models.DSPRecipeResponse{
				ConsumesPerSec:           cr.ItemsConsumedPerSec,
				CraftingPerSecond:        &cps,
				For:                      &usedFor,
				MadeIn:                   &madeIn,
				NumberOfFacilitiesNeeded: &numFacilities,
				Produce:                  &product,
				SecondsSpendPerCrafting:  &sspc,
			})
		}
	}

	jsonStr, _ := json.MarshalIndent(recipe, "", "\t")
	fmt.Println("RECIPE")
	fmt.Println(string(jsonStr))

	return operations.NewGetDSPRecipeOK().WithPayload(recipe)
}

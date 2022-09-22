package handlers

import (
	"alex-api/models"
	"alex-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetDSPRecipesHandler(params operations.GetDSPRecipesParams) middleware.Responder {
	recipes := optimizer.GetRecipes()

	recipeModels := []*models.DSPRecipe{}

	for _, recipe := range recipes {
		madeIn := recipe.Facility
		materials := []*models.DSPMaterial{}

		for mat, count := range recipe.Materials {
			c := float64(count)
			m := string(mat)
			materials = append(materials, &models.DSPMaterial{
				Count: &c,
				Name:  &m,
			})
		}

		name := string(recipe.OutputItem)
		produce := float64(recipe.OutputItemCount)
		time := float64(recipe.Time)

		r := &models.DSPRecipe{
			MadeIn:    &madeIn,
			Materials: materials,
			Name:      &name,
			Produce:   &produce,
			Time:      &time,
		}

		recipeModels = append(recipeModels, r)
	}

	return operations.NewGetDSPRecipesOK().WithPayload(recipeModels)
}

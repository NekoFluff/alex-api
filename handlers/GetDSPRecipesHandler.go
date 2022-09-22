package handlers

import (
	"alex-api/models"
	"alex-api/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetDSPRecipesHandler(params operations.GetDSPRecipesParams) middleware.Responder {
	madeIn := "test"
	materials := []*models.DSPMaterial{}
	name := "name"
	produce := 1.0
	time := 1.0

	recipe := models.DSPRecipe{
		MadeIn:    &madeIn,
		Materials: materials,
		Name:      &name,
		Produce:   &produce,
		Time:      &time,
	}

	recipes := []*models.DSPRecipe{&recipe}
	return operations.NewGetDSPRecipesOK().WithPayload(recipes)
}

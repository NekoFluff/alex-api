package handlers

import (
	"addi/dsp"
	"addi/models"
	"addi/restapi/operations"
	"encoding/json"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
)

func PostDSPHandler(params operations.PostDspParams) middleware.Responder {

	// log.Println("Starting DSP Optimizer Program")

	recipe := []*models.Recipe{}

	for _, v := range params.RecipeRequest {
		recipe = append(recipe, dsp.GetRecipeForItem(*v.Name, float64(*v.Count), "")...)
	}

	jsonStr, _ := json.MarshalIndent(recipe, "", "\t")
	fmt.Println("RECIPE")
	fmt.Println(string(jsonStr))

	return operations.NewPostDspOK().WithPayload(recipe)
}

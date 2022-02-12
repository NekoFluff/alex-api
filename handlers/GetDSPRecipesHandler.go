package handlers

import (
	"addi/dsp"
	"addi/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetDSPRecipesHandler(params operations.GetDSPRecipesParams) middleware.Responder {
	return operations.NewGetDSPRecipesOK().WithPayload(dsp.GetItems())
}

package handlers

import (
	"addi/dsp"
	"addi/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func ReloadDSPRecipesHandler(params operations.ReloadDSPRecipesParams) middleware.Responder {
	dsp.ScrapeDSPItems()
	return operations.NewReloadDSPRecipesOK()
}

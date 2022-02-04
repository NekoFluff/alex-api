package handlers

import (
	"addi/dsp"
	"addi/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func GetDSPItemsHandler(params operations.GetDspItemsParams) middleware.Responder {
	return operations.NewGetDspItemsOK().WithPayload(dsp.GetItems())
}

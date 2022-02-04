package handlers

import (
	"addi/dsp"
	"addi/restapi/operations"

	"github.com/go-openapi/runtime/middleware"
)

func PostDSPItemsReloadHandler(params operations.PostDspItemsReloadParams) middleware.Responder {
	dsp.ScrapeDSPItems()
	return operations.NewPostDspItemsReloadOK()
}

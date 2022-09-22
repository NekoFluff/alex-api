package handlers

import (
	"alex-api/internal/dsphelper"
	"alex-api/restapi/operations"

	"github.com/NekoFluff/go-dsp/dsp"
	"github.com/go-openapi/runtime/middleware"
)

func ReloadDSPRecipesHandler(params operations.ReloadDSPRecipesParams) middleware.Responder {
	dsphelper.Scrape()
	optimizer = dsp.NewOptimizer(dsp.OptimizerConfig{
		DataSource: "../data/items.json",
	})
	return operations.NewReloadDSPRecipesOK()
}

package server

import (
	"net/http"
)

func (s *Server) Route() {
	// Dyson Sphere Program Endpoints
	s.router.Handle("/dsp/computedRecipes", (s.DSPComputedRecipes())).
		Methods(http.MethodPost).Name("getDSPComputedRecipe")
	s.router.Handle("/dsp/recipes", (s.DSPRecipes())).
		Methods(http.MethodGet).Name("getDSPRecipes")
	s.router.Handle("/dsp/recipes/reload", (s.DSPReloadRecipes())).
		Methods(http.MethodPost).Name("reloadDSPRecipes")

	// Black Desert Online Endpoints
	s.router.Handle("/bdo/computedRecipes", (s.BDOComputedRecipes())).
		Methods(http.MethodPost).Name("getBDOComputedRecipe")
	s.router.Handle("/bdo/recipes", (s.BDORecipes())).
		Methods(http.MethodGet).Name("getBDORecipes")

	// Twitter Endpoints
	s.router.Handle("/twitter/inArt", (s.InArt())).
		Methods(http.MethodGet).Name("getInArt")
}

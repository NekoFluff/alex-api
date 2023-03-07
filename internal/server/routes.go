package server

import (
	"alex-api/internal/middleware"
	"net/http"
)

func (s *Server) Route() {
	// Dyson Sphere Program Endpoints
	s.router.Handle("/dsp/computedRecipes", middleware.CORSMiddleware(s.DSPComputedRecipes())).
		Methods(http.MethodPost).Name("getDSPComputedRecipe")
	s.router.Handle("/dsp/recipes", middleware.CORSMiddleware(s.DSPRecipes())).
		Methods(http.MethodGet).Name("getDSPRecipes")
	s.router.Handle("/dsp/recipes/reload", middleware.CORSMiddleware(s.DSPReloadRecipes())).
		Methods(http.MethodPost).Name("reloadDSPRecipes")

	// Black Desert Online Endpoints
	s.router.Handle("/bdo/computedRecipes", middleware.CORSMiddleware(s.BDOComputedRecipes())).
		Methods(http.MethodPost).Name("getBDOComputedRecipe")
	s.router.Handle("/bdo/recipes", middleware.CORSMiddleware(s.BDORecipes())).
		Methods(http.MethodGet).Name("getBDORecipes")

	// Twitter Endpoints
	s.router.Handle("/inArt", middleware.CORSMiddleware(s.InArt())).
		Methods(http.MethodGet).Name("getInArt")
}

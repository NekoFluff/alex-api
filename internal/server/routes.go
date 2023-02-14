package server

import (
	"alex-api/internal/middleware"
	"net/http"
)

func (s *Server) Route() {
	s.router.Handle("/dsp/computedRecipe", middleware.CORSMiddleware(s.ComputedRecipe())).
		Methods(http.MethodPost).Name("getComputedRecipe")
	s.router.Handle("/dsp/recipes", middleware.CORSMiddleware(s.Recipes())).
		Methods(http.MethodGet).Name("getRecipes")
	s.router.Handle("/dsp/recipes/reload", middleware.CORSMiddleware(s.ReloadRecipes())).
		Methods(http.MethodGet).Name("reloadRecipes")
	s.router.Handle("/inArt", middleware.CORSMiddleware(s.InArt())).
		Methods(http.MethodGet).Name("getInArt")
}

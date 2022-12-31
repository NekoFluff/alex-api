package server

import (
	"net/http"
)

func (s *Server) Route() {
	s.router.Handle("/dsp/computedRecipe", s.ComputedRecipe()).
		Methods(http.MethodPost).Name("getComputedRecipe")
	s.router.Handle("/dsp/recipes", s.Recipes()).
		Methods(http.MethodGet).Name("getRecipes")
	s.router.Handle("/dsp/recipes/reload", s.ReloadRecipes()).
		Methods(http.MethodGet).Name("reloadRecipes")
	s.router.Handle("/inArt", s.InArt()).
		Methods(http.MethodGet).Name("getInArt")
}

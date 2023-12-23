package server

import (
	"alex-api/internal/dspscraper"
	"alex-api/internal/recipecalc"
	"alex-api/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *Server) DSPComputedRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body []recipecalc.ComputedRecipeRequest
		err := utils.DecodeValidate(r.Body, s.validator, &body)
		defer r.Body.Close()
		if err != nil {
			l.WithError(err).Error("failed to decode request")
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		l = l.WithField("body", body)

		var computedRecipes = []recipecalc.ComputedRecipe{}
		for _, recipeRequest := range body {
			itemName := recipeRequest.Name
			optimalRecipe := s.dspOptimizer.GetOptimalRecipe(itemName, recipeRequest.Rate, "", map[string]bool{}, 0, recipecalc.RecipeRequirements(recipeRequest.Requirements), recipeRequest.AssemberLevel)
			computedRecipes = append(computedRecipes, optimalRecipe...)
		}

		q := r.URL.Query()
		if q.Get("group") == "true" {
			s.dspOptimizer.SortRecipes(computedRecipes)
			computedRecipes = s.dspOptimizer.CombineRecipes(computedRecipes)
		}
		s.dspOptimizer.SortRecipes(computedRecipes)
		l.WithField("computedRecipes", computedRecipes).Info("Computed Recipes")

		jsonStr, err := json.MarshalIndent(computedRecipes, "", "\t")
		if err != nil {
			l.WithError(err).Error("failed to marshal computed recipes")
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) DSPReloadRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		dspscraper.Scrape()
		s.dspOptimizer.SetRecipes(recipecalc.LoadDSPRecipes(l, s.db))

		l.Info("Reloaded Recipes")
		_, _ = w.Write([]byte("Successfully reloaded recipes"))
	}

}

func (s *Server) DSPRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		recipes := s.dspOptimizer.GetRecipes()
		l.WithField("recipes", recipes).Info("recipes")

		jsonStr, _ := json.MarshalIndent(recipes, "", "\t")
		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) BDOComputedRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body []recipecalc.ComputedRecipeRequest
		err := utils.DecodeValidate(r.Body, s.validator, &body)
		defer r.Body.Close()
		if err != nil {
			l.WithError(err).Error("failed to decode request")
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		l = l.WithField("body", body)

		var computedRecipes = []recipecalc.ComputedRecipe{}
		for _, recipeRequest := range body {
			itemName := recipeRequest.Name
			optimalRecipe := s.bdoOptimizer.GetOptimalRecipe(itemName, recipeRequest.Rate, "", map[string]bool{}, 0, recipeRequest.Requirements, 2)
			computedRecipes = append(computedRecipes, optimalRecipe...)
		}

		q := r.URL.Query()
		if q.Get("group") == "true" {
			s.bdoOptimizer.SortRecipes(computedRecipes)
			computedRecipes = s.bdoOptimizer.CombineRecipes(computedRecipes)
		}
		s.bdoOptimizer.SortRecipes(computedRecipes)
		l.WithField("computedRecipes", computedRecipes).Info("Computed Recipes")

		jsonStr, err := json.MarshalIndent(computedRecipes, "", "\t")
		if err != nil {
			l.WithError(err).Error("failed to marshal computed recipes")
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) BDORecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		recipes := s.bdoOptimizer.GetRecipes()
		l.WithField("recipes", recipes).Info("recipes")

		jsonStr, _ := json.MarshalIndent(recipes, "", "\t")
		_, _ = w.Write(jsonStr)
	}
}

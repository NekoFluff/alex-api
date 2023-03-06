package server

import (
	"alex-api/internal/dspscraper"
	"alex-api/utils"
	"encoding/json"
	"net/http"

	"github.com/NekoFluff/go-dsp/dsp"
	"github.com/sirupsen/logrus"
)

func (s *Server) ComputedRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body []dspscraper.ComputedRecipeRequest
		err := utils.DecodeValidate(r.Body, s.validator, &body)
		defer r.Body.Close()
		if err != nil {
			l.WithError(err).Error("failed to decode request")
			w.Write([]byte(err.Error()))
			return
		}
		l = l.WithField("body", body)

		var computedRecipes = []dsp.ComputedRecipe{}
		for _, recipeRequest := range body {
			itemName := dsp.ItemName(recipeRequest.Name)
			optimalRecipe := optimizer.GetOptimalRecipe(itemName, recipeRequest.Rate, "", map[dsp.ItemName]bool{}, 0, recipeRequest.Requirements)
			computedRecipes = append(computedRecipes, optimalRecipe...)
		}

		q := r.URL.Query()
		if q.Get("group") == "1" {
			optimizer.SortRecipes(computedRecipes)
			computedRecipes = optimizer.CombineRecipes(computedRecipes)
		}
		optimizer.SortRecipes(computedRecipes)
		l.WithField("computedRecipes", computedRecipes).Info("Computed Recipes")

		jsonStr, err := json.MarshalIndent(computedRecipes, "", "\t")
		if err != nil {
			l.WithError(err).Error("failed to marshal computed recipes")
			w.Write([]byte(err.Error()))
			return
		}

		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) ReloadRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		dspscraper.Scrape()
		optimizer.LoadRecipes()

		l.Info("Reloaded Recipes")
		_, _ = w.Write([]byte("Successfully reloaded recipes"))
	}

}

func (s *Server) Recipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		recipes := optimizer.GetRecipes()
		l.WithField("recipes", recipes).Info("recipes")

		jsonStr, _ := json.MarshalIndent(recipes, "", "\t")
		_, _ = w.Write(jsonStr)
	}

}

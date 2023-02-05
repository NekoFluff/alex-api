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
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body dspscraper.ComputedRecipeRequest
		err := utils.DecodeValidate(r.Body, s.validator, &body)
		defer r.Body.Close()
		if err != nil {
			l.WithError(err).Error("failed to decode request")
			w.Write([]byte(err.Error()))
			return
		}

		itemName := dsp.ItemName(body.Name)
		computedRecipes := optimizer.GetOptimalRecipe(itemName, body.Rate, "", map[dsp.ItemName]bool{})

		jsonStr, _ := json.MarshalIndent(computedRecipes, "", "\t")
		l.Info("COMPUTED RECIPE")
		l.Info(string(jsonStr))

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
		optimizer = dsp.NewOptimizer(dsp.OptimizerConfig{
			DataSource: "data/items.json",
		})

		l.Info("RELOAD RECIPES")
		_, _ = w.Write([]byte("RELOAD COMPUTED RECIPES"))
	}

}

func (s *Server) Recipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		recipes := optimizer.GetRecipes()

		jsonStr, _ := json.MarshalIndent(recipes, "", "\t")
		l.Info("RECIPES")
		l.Info(string(jsonStr))
		_, _ = w.Write(jsonStr)
	}

}

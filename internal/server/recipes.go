package server

import (
	"alex-api/internal/data"
	"alex-api/internal/dspscraper"
	"alex-api/internal/utils"
	"encoding/json"
	"net/http"

	"github.com/NekoFluff/go-dsp/dsp"
	"github.com/sirupsen/logrus"
)

func (s *Server) DSPComputedRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body []ComputedRecipeRequest
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

		crList := []ComputedRecipe{}
		for _, cr := range computedRecipes {
			crList = append(crList, s.convertDSPComputedRecipeToComputedRecipe(cr))
		}

		jsonStr, err := json.MarshalIndent(crList, "", "\t")
		if err != nil {
			l.WithError(err).Error("failed to marshal computed recipes")
			w.Write([]byte(err.Error()))
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
		optimizer.LoadRecipes()

		l.Info("Reloaded Recipes")
		_, _ = w.Write([]byte("Successfully reloaded recipes"))
	}

}

func (s *Server) convertDSPComputedRecipeToComputedRecipe(recipe dsp.ComputedRecipe) ComputedRecipe {
	m := map[string]float64{}

	for k, v := range recipe.ItemsConsumedPerSec {
		m[string(k)] = v
	}

	return ComputedRecipe{
		OutputItem:           string(recipe.OutputItem),
		Facility:             recipe.Facility,
		NumFacilitiesNeeded:  recipe.NumFacilitiesNeeded,
		ItemsConsumedPerSec:  m,
		SecondsSpentPerCraft: recipe.SecondsSpentPerCraft,
		CraftingPerSec:       recipe.CraftingPerSec,
		UsedFor:              recipe.UsedFor,
		Depth:                int64(recipe.Depth),
	}
}

func (s *Server) DSPRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		recipes := optimizer.GetRecipes()
		l.WithField("recipes", recipes).Info("recipes")

		recipesResp := [][]Recipe{}

		for _, recipeList := range recipes {
			rList := []Recipe{}
			for _, recipe := range recipeList {
				rList = append(rList, s.convertDSPRecipeToRecipe(recipe))
			}
			recipesResp = append(recipesResp, rList)
		}

		jsonStr, _ := json.MarshalIndent(recipesResp, "", "\t")
		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) convertDSPRecipeToRecipe(recipe dsp.Recipe) Recipe {
	m := Materials{}

	for k, v := range recipe.Materials {
		m[string(k)] = v
	}

	return Recipe{
		OutputItem:      string(recipe.OutputItem),
		OutputItemCount: recipe.OutputItemCount,
		Facility:        recipe.Facility,
		Time:            recipe.Time,
		Materials:       m,
	}
}

func (s *Server) BDOComputedRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		var body []ComputedRecipeRequest
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

func (s *Server) BDORecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := s.logger.WithContext(ctx).WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.EscapedPath(),
		})

		db := data.New(l)
		defer db.Disconnect()
		recipes, err := db.GetRecipes(nil, nil)
		if err != nil {
			l.WithError(err).Error("failed to retrieve bdo recipes")
			w.Write([]byte(err.Error()))
			return
		}
		l.WithField("recipes", recipes).Info("recipes")

		rList := []Recipe{}
		for _, recipe := range recipes {
			rList = append(rList, s.convertBDORecipeToRecipe(recipe))
		}

		jsonStr, _ := json.MarshalIndent(rList, "", "\t")
		_, _ = w.Write(jsonStr)
	}
}

func (s *Server) convertBDORecipeToRecipe(recipe data.Recipe) Recipe {
	m := Materials{}

	for _, v := range recipe.Recipe {
		m[v.ItemName] = v.Amount
	}

	r := Recipe{
		OutputItem:         recipe.Name,
		OutputItemCount:    recipe.QuantityProduced,
		MinOutputItemCount: recipe.MinProduced,
		MaxOutputItemCount: recipe.MaxProduced,
		Facility:           recipe.Action,
		Time:               recipe.TimeToProduce,
		Materials:          m,
		Image:              recipe.Image,
	}

	if recipe.MarketData != nil {
		r.MarketData = &MarketData{
			LastUpdateAttempt: recipe.MarketData.LastUpdateAttempt,
			LastUpdated:       recipe.MarketData.LastUpdated,
			Price:             recipe.MarketData.Price,
			Quantity:          recipe.MarketData.Quantity,
			TotalTradeCount:   recipe.MarketData.TotalTradeCount,
			Name:              recipe.MarketData.Name,
		}
	}

	return r
}

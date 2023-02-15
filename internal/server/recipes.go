package server

import (
	"alex-api/internal/dspscraper"
	"alex-api/utils"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"

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

		var computedRecipes = []dsp.ComputedRecipe{}
		for _, recipeRequest := range body {
			l.Info(recipeRequest.Name, recipeRequest.Rate)
			itemName := dsp.ItemName(recipeRequest.Name)
			optimalRecipe := optimizer.GetOptimalRecipe(itemName, recipeRequest.Rate, "", map[dsp.ItemName]bool{}, 0)
			l.Info(optimalRecipe)
			computedRecipes = append(computedRecipes, optimalRecipe...)
			l.Info(computedRecipes)
		}

		q := r.URL.Query()
		if q.Get("group") == "1" {
			computedRecipes = combineRecipes(computedRecipes)
		}
		sortRecipes(computedRecipes)

		jsonStr, err := json.MarshalIndent(computedRecipes, "", "\t")
		if err != nil {
			l.WithError(err).Error("failed to marshal computed recipes")
			w.Write([]byte(err.Error()))
			return
		}
		l.Info("COMPUTED RECIPES")
		l.Info(string(jsonStr))

		_, _ = w.Write(jsonStr)
	}
}

func combineRecipes(recipes []dsp.ComputedRecipe) []dsp.ComputedRecipe {
	uniqueRecipes := make(map[dsp.ItemName]dsp.ComputedRecipe)

	for _, recipe := range recipes {
		if uRecipe, ok := uniqueRecipes[recipe.OutputItem]; ok { // combine recipe objects

			old_num := uRecipe.NumFacilitiesNeeded
			new_num := recipe.NumFacilitiesNeeded
			total_num := old_num + new_num
			for materialName, perSecConsumption := range uRecipe.ItemsConsumedPerSec {
				uRecipe.ItemsConsumedPerSec[materialName] = perSecConsumption + recipe.ItemsConsumedPerSec[materialName]
			}

			sspc := (uRecipe.SecondsSpentPerCraft*old_num + recipe.SecondsSpentPerCraft*new_num) / total_num
			if math.IsNaN(float64(sspc)) {
				sspc = 0.0
			}
			uRecipe.SecondsSpentPerCraft = sspc

			uRecipe.CraftingPerSec = uRecipe.CraftingPerSec + recipe.CraftingPerSec
			uRecipe.UsedFor = dsp.ItemName(fmt.Sprintf("%s | %s (Uses %0.2f/s)", uRecipe.UsedFor, recipe.UsedFor, recipe.CraftingPerSec))
			// uRecipe.UsedFor = uRecipe.UsedFor.filter((v, i, a) => a.indexOf(v) === i); // get unique values
			uRecipe.NumFacilitiesNeeded += recipe.NumFacilitiesNeeded
			uRecipe.Depth = max(uRecipe.Depth, recipe.Depth)
			uniqueRecipes[recipe.OutputItem] = uRecipe

		} else { // add recipe object
			recipe.UsedFor = dsp.ItemName(fmt.Sprintf("%s (Uses %0.2f/s)", recipe.UsedFor, recipe.CraftingPerSec))
			uniqueRecipes[recipe.OutputItem] = recipe
		}
	}

	v := []dsp.ComputedRecipe{}
	for _, value := range uniqueRecipes {
		v = append(v, value)
	}
	return v
}

func sortRecipes(recipes []dsp.ComputedRecipe) {
	sort.SliceStable(recipes, func(i, j int) bool {
		if recipes[i].Depth != recipes[j].Depth {
			return recipes[i].Depth < recipes[j].Depth
		} else if recipes[i].OutputItem != recipes[j].OutputItem {
			return recipes[i].OutputItem < recipes[j].OutputItem
		} else if recipes[i].UsedFor != recipes[j].UsedFor {
			return recipes[i].UsedFor < recipes[j].UsedFor
		} else {
			return recipes[i].CraftingPerSec < recipes[j].CraftingPerSec
		}
	})
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
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

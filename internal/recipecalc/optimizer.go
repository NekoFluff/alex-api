package recipecalc

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Optimizer struct {
	recipeMap map[string][]Recipe
	config    OptimizerConfig
	log       *logrus.Entry
}

type OptimizerConfig struct{}

func NewOptimizer(log *logrus.Entry, config OptimizerConfig) *Optimizer {
	o := &Optimizer{
		recipeMap: make(map[string][]Recipe),
		config:    config,
		log:       log,
	}
	return o
}

func (o *Optimizer) SetRecipes(recipes map[string][]Recipe) {
	o.recipeMap = recipes
}

func (o *Optimizer) GetRecipe(itemName string, recipeIdx int) (Recipe, bool) {
	name := strings.ToLower(string(itemName))
	recipes, ok := o.recipeMap[name]

	if !ok {
		return Recipe{}, ok
	}

	if len(recipes) > recipeIdx {
		recipe := recipes[recipeIdx]
		return recipe, true
	}

	return recipes[0], true
}

func (o *Optimizer) GetRecipes() [][]Recipe {
	recipes := [][]Recipe{}
	for _, recipe := range o.recipeMap {
		recipes = append(recipes, recipe)
	}
	return recipes
}

func (o *Optimizer) GetOptimalRecipe(itemName string, craftingSpeed float64, parentItemName string, seenRecipes map[string]bool, depth int64, recipeRequirements RecipeRequirements) []ComputedRecipe {
	computedRecipes := []ComputedRecipe{}

	if seenRecipes[itemName] {
		return computedRecipes
	}
	seenRecipes[itemName] = true

	rRequirement, ok := recipeRequirements[itemName]
	recipeIdx := 0
	if ok {
		recipeIdx = rRequirement
	}

	recipe, ok := o.GetRecipe(itemName, recipeIdx)
	if !ok {
		return computedRecipes
	}

	consumedMats := make(map[string]float64)
	numberOfFacilitiesNeeded := guardInf(float64(recipe.Time * craftingSpeed / recipe.OutputItemCount))

	for materialName, materialCount := range recipe.Materials {
		consumedMats[materialName] = guardInf(float64(materialCount * numberOfFacilitiesNeeded / recipe.Time))
	}

	computedRecipe := ComputedRecipe{
		OutputItem:           recipe.OutputItem,
		Facility:             recipe.Facility,
		NumFacilitiesNeeded:  numberOfFacilitiesNeeded,
		ItemsConsumedPerSec:  consumedMats,
		SecondsSpentPerCraft: recipe.Time,
		CraftingPerSec:       craftingSpeed,
		UsedFor:              string(parentItemName),
		Depth:                depth,
		Image:                recipe.Image,
	}
	computedRecipes = append(computedRecipes, computedRecipe)

	for materialName, materialCountPerSec := range computedRecipe.ItemsConsumedPerSec {
		targetCraftingSpeed := materialCountPerSec
		seenRecipesCopy := make(map[string]bool)
		for k, v := range seenRecipes {
			seenRecipesCopy[k] = v
		}
		cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.OutputItem, seenRecipesCopy, depth+1, recipeRequirements)
		computedRecipes = append(computedRecipes, cr...)
	}

	return computedRecipes
}

func (o *Optimizer) SortRecipes(recipes []ComputedRecipe) {
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

func (o *Optimizer) CombineRecipes(recipes []ComputedRecipe) []ComputedRecipe {
	uniqueRecipes := make(map[string]ComputedRecipe)

	for _, recipe := range recipes {
		if uRecipe, ok := uniqueRecipes[recipe.OutputItem]; ok { // combine recipe objects

			old_num := uRecipe.NumFacilitiesNeeded
			new_num := recipe.NumFacilitiesNeeded
			total_num := old_num + new_num
			for materialName, perSecConsumption := range uRecipe.ItemsConsumedPerSec {
				uRecipe.ItemsConsumedPerSec[materialName] = perSecConsumption + recipe.ItemsConsumedPerSec[materialName]
			}

			uRecipe.SecondsSpentPerCraft = guardInf(float64(uRecipe.SecondsSpentPerCraft*old_num+recipe.SecondsSpentPerCraft*new_num) / total_num)
			uRecipe.CraftingPerSec = uRecipe.CraftingPerSec + recipe.CraftingPerSec
			uRecipe.UsedFor = fmt.Sprintf("%s | %s (Uses %0.2f/s)", uRecipe.UsedFor, recipe.UsedFor, recipe.CraftingPerSec)
			uRecipe.NumFacilitiesNeeded += recipe.NumFacilitiesNeeded
			uRecipe.Depth = max(uRecipe.Depth, recipe.Depth)
			uniqueRecipes[recipe.OutputItem] = uRecipe

		} else { // add recipe object
			if recipe.UsedFor != "" {
				recipe.UsedFor = fmt.Sprintf("%s (Uses %0.2f/s)", recipe.UsedFor, recipe.CraftingPerSec)
			}
			uniqueRecipes[recipe.OutputItem] = recipe
		}
	}

	v := []ComputedRecipe{}
	for _, value := range uniqueRecipes {
		v = append(v, value)
	}
	return v
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func guardInf(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0.0
	}
	return x
}

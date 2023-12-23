package recipecalc

import (
	"alex-api/internal/data"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type Optimizer struct {
	recipeMap map[string][]data.Recipe
	config    OptimizerConfig
	log       *logrus.Entry
}

type OptimizerConfig struct{}

func NewOptimizer(log *logrus.Entry, config OptimizerConfig) *Optimizer {
	o := &Optimizer{
		recipeMap: make(map[string][]data.Recipe),
		config:    config,
		log:       log,
	}
	return o
}

func (o *Optimizer) SetRecipes(recipes map[string][]data.Recipe) {
	o.recipeMap = recipes
}

func (o *Optimizer) GetRecipe(itemName string, recipeIdx int) (data.Recipe, bool) {
	name := strings.ToLower(string(itemName))
	recipes, ok := o.recipeMap[name]

	if !ok {
		return data.Recipe{}, ok
	}

	if len(recipes) > recipeIdx {
		recipe := recipes[recipeIdx]
		return recipe, true
	}

	return recipes[0], true
}

func (o *Optimizer) GetRecipes() [][]data.Recipe {
	recipes := [][]data.Recipe{}
	for _, recipe := range o.recipeMap {
		recipes = append(recipes, recipe)
	}
	return recipes
}

func (o *Optimizer) GetOptimalRecipe(itemName string, craftingSpeed float64, parentItemName string, seenRecipes map[string]bool, depth int64, recipeRequirements RecipeRequirements, assemblerLevel, smelterLevel, chemicalPlantLevel int) []ComputedRecipe {
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
	speedModifier := 1.0
	if (recipe.Facility != "") && (strings.Contains(strings.ToLower(recipe.Facility), "assembling machine")) {
		switch assemblerLevel {
		case 1:
			speedModifier = 0.75
			recipe.Facility = "Assembling Machine Mk.I"
		case 2:
			speedModifier = 1.0
			recipe.Facility = "Assembling Machine Mk.II"
		case 3:
			speedModifier = 1.5
			recipe.Facility = "Assembling Machine Mk.III"
		default:
			speedModifier = 1.0
			recipe.Facility = "Assembling Machine Mk.II"
		}
	} else if (recipe.Facility != "") && (strings.Contains(strings.ToLower(recipe.Facility), "smelter")) {
		switch smelterLevel {
		case 1:
			speedModifier = 1.0
			recipe.Facility = "Arc Smelter"
		case 2:
			speedModifier = 2.0
			recipe.Facility = "Plane Smelter"
		default:
			speedModifier = 1.0
			recipe.Facility = "Arc Smelter"
		}
	} else if (recipe.Facility != "") && (strings.Contains(strings.ToLower(recipe.Facility), "chemical plant")) {
		switch chemicalPlantLevel {
		case 1:
			speedModifier = 1.0
			recipe.Facility = "Chemical Plant"
		case 2:
			speedModifier = 2.0
			recipe.Facility = "Quantum Chemical Plant"
		default:
			speedModifier = 1.0
			recipe.Facility = "Chemical Plant"
		}
	}

	numberOfFacilitiesNeeded := guardInf(float64(recipe.TimeToProduce * craftingSpeed / recipe.QuantityProduced))

	for materialName, materialCount := range recipe.Ingredients {
		consumedMats[materialName] = guardInf(float64(materialCount * numberOfFacilitiesNeeded / recipe.TimeToProduce))
	}

	computedRecipe := ComputedRecipe{
		Name:                 recipe.Name,
		Facility:             recipe.Facility,
		NumFacilitiesNeeded:  numberOfFacilitiesNeeded / speedModifier,
		ItemsConsumedPerSec:  consumedMats,
		SecondsSpentPerCraft: recipe.TimeToProduce,
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
		cr := o.GetOptimalRecipe(materialName, targetCraftingSpeed, recipe.Name, seenRecipesCopy, depth+1, recipeRequirements, assemblerLevel, smelterLevel, chemicalPlantLevel)
		computedRecipes = append(computedRecipes, cr...)
	}

	return computedRecipes
}

func (o *Optimizer) SortRecipes(recipes []ComputedRecipe) {
	sort.SliceStable(recipes, func(i, j int) bool {
		if recipes[i].Depth != recipes[j].Depth {
			return recipes[i].Depth < recipes[j].Depth
		} else if recipes[i].Name != recipes[j].Name {
			return recipes[i].Name < recipes[j].Name
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
		if uRecipe, ok := uniqueRecipes[recipe.Name]; ok { // combine recipe objects

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
			uniqueRecipes[recipe.Name] = uRecipe

		} else { // add recipe object
			if recipe.UsedFor != "" {
				recipe.UsedFor = fmt.Sprintf("%s (Uses %0.2f/s)", recipe.UsedFor, recipe.CraftingPerSec)
			}
			uniqueRecipes[recipe.Name] = recipe
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

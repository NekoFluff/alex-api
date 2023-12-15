//go:generate mockgen -source=optimizer.go -destination=optimizer_mock_test.go -package=server
package server

import (
	"alex-api/internal/data"
	"alex-api/internal/recipecalc"
)

type Optimizer interface {
	GetOptimalRecipe(itemName string, rate float64, recipeName string, ignore map[string]bool, depth int64, requirements recipecalc.RecipeRequirements) []recipecalc.ComputedRecipe
	SortRecipes(recipes []recipecalc.ComputedRecipe)
	CombineRecipes(recipes []recipecalc.ComputedRecipe) []recipecalc.ComputedRecipe
	SetRecipes(recipes map[string][]data.Recipe)
	GetRecipes() [][]data.Recipe
}

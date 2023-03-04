package dspscraper

import "github.com/NekoFluff/go-dsp/dsp"

type ComputedRecipeRequest struct {
	Name         string                 `json:"name" validate:"required"`
	Rate         float64                `json:"rate" validate:"required"`
	Requirements dsp.RecipeRequirements `json:"requirements"`
}

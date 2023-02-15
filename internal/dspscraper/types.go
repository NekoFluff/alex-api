package dspscraper

type ComputedRecipeRequest struct {
	Name string  `json:"name" validate:"required"`
	Rate float64 `json:"rate" validate:"required"`
}

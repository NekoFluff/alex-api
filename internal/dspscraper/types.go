package dspscraper

type ComputedRecipeRequest struct {
	Name string  `json:"name" validate:"required"`
	Rate float32 `json:"rate" validate:"required"`
}

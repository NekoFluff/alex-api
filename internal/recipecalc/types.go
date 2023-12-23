package recipecalc

import "encoding/json"

type ComputedRecipeRequest struct {
	Name          string             `json:"name" validate:"required"`
	Rate          float64            `json:"rate" validate:"required"`
	AssemberLevel int                `json:"assemblerLevel" validate:"gte=1,lte=3"`
	Requirements  RecipeRequirements `json:"requirements,omitempty"`
}

func (crr *ComputedRecipeRequest) UnmarshalJSON(b []byte) error {
	crr.AssemberLevel = 2
	type alias ComputedRecipeRequest
	return json.Unmarshal(b, (*alias)(crr))
}

type RecipeRequirements map[string]int

type ComputedRecipe struct {
	Name                 string             `json:"name" validate:"required"`
	Facility             string             `json:"facility"`
	NumFacilitiesNeeded  float64            `json:"numFacilitiesNeeded"`
	ItemsConsumedPerSec  map[string]float64 `json:"itemsConsumedPerSec"`
	SecondsSpentPerCraft float64            `json:"secondsSpentPerCraft"`
	CraftingPerSec       float64            `json:"craftingPerSec"`
	UsedFor              string             `json:"usedFor"`
	Depth                int64              `json:"depth,omitempty"`
	Image                string             `json:"image,omitempty"`
}

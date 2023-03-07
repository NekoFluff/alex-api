package server

import (
	"time"

	"github.com/NekoFluff/go-dsp/dsp"
)

type Recipe struct {
	OutputItem         string      `json:"outputItem" validate:"required"`
	OutputItemCount    float64     `json:"outputItemCount" validate:"required"`
	MinOutputItemCount float64     `json:"minOutputItemCount,omitempty"`
	MaxOutputItemCount float64     `json:"maxOutputItemCount,omitempty"`
	Facility           string      `json:"facility" validate:"required"`
	Time               float64     `json:"time" validate:"required"`
	Materials          Materials   `json:"materials" validate:"required"`
	Image              string      `json:"image,omitempty"`
	MarketData         *MarketData `json:"marketData,omitempty"`
}

type Materials map[string]float64

type ComputedRecipeRequest struct {
	Name         string                 `json:"name" validate:"required"`
	Rate         float64                `json:"rate" validate:"required"`
	Requirements dsp.RecipeRequirements `json:"requirements,omitempty"`
}

type ComputedRecipe struct {
	OutputItem           string             `json:"outputItem" validate:"required"`
	Facility             string             `json:"facility" validate:"required"`
	NumFacilitiesNeeded  float64            `json:"numFacilitiesNeeded" validate:"required"`
	ItemsConsumedPerSec  map[string]float64 `json:"itemsConsumedPerSec" validate:"required"`
	SecondsSpentPerCraft float64            `json:"secondsSpentPerCraft" validate:"required"`
	CraftingPerSec       float64            `json:"craftingPerSec" validate:"required"`
	UsedFor              string             `json:"usedFor" validate:"required"`
	Depth                int64              `json:"depth,omitempty"`
	Image                string             `json:"image,omitempty"`
}

type MarketData struct {
	LastUpdateAttempt time.Time `json:"lastUpdateAttempt,omitempty"`
	LastUpdated       time.Time `json:"lastUpdated,omitempty"`
	Price             float64   `json:"price,omitempty"`
	Quantity          float64   `json:"quantity,omitempty"`
	TotalTradeCount   float64   `json:"totalTradeCount,omitempty"`
	Name              string    `json:"name,omitempty"`
}

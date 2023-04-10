package server

import (
	"alex-api/internal/config"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	os.Setenv("PORT", "8080")
	server := New(config.Config{}, logrus.NewEntry(logrus.New()), nil, nil, nil, nil)

	expected := []string{
		"getDSPComputedRecipe",
		"getDSPRecipes",
		"reloadDSPRecipes",
		"getBDOComputedRecipe",
		"getBDORecipes",
		"getInArt",
		"pageViewed",
	}
	var received []string
	_ = server.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		received = append(received, route.GetName())
		return nil
	})

	assert.Equal(t, expected, received)
}

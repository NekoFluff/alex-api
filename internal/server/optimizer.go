package server

import (
	"alex-api/internal/recipecalc"
	"context"

	"github.com/sirupsen/logrus"
)

var dspOptimizer *recipecalc.Optimizer
var bdoOptimizer *recipecalc.Optimizer

func initializeOptimizers() {
	log := logrus.New().WithContext(context.TODO())

	dspOptimizer = recipecalc.NewOptimizer(log, recipecalc.OptimizerConfig{})
	dspOptimizer.SetRecipes(recipecalc.LoadDSPRecipes())

	bdoOptimizer = recipecalc.NewOptimizer(log, recipecalc.OptimizerConfig{})
	bdoOptimizer.SetRecipes(recipecalc.LoadBDORecipes())
}

package handlers

import "github.com/NekoFluff/go-dsp/dsp"

var optimizer *dsp.Optimizer

func init() {
	optimizer = dsp.NewOptimizer(dsp.OptimizerConfig{
		DataSource: "../data/items.json",
	})
}

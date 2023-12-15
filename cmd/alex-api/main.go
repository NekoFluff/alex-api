package main

import (
	"alex-api/internal/config"
	"alex-api/internal/data"
	"alex-api/internal/recipecalc"
	"alex-api/internal/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	log := logger.WithContext(ctx)

	cfg, err := config.New(".env")
	if err != nil {
		log.WithError(err).Fatal("failed to process config")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		log.Info("termination signaled")
		cancel()
	}()

	log.Info("Starting up Alex API")

	db := data.New(log)
	defer db.Disconnect()

	dspOptimizer := recipecalc.NewOptimizer(log, recipecalc.OptimizerConfig{})
	dspOptimizer.SetRecipes(recipecalc.LoadDSPRecipes(log, db))

	bdoOptimizer := recipecalc.NewOptimizer(log, recipecalc.OptimizerConfig{})
	bdoOptimizer.SetRecipes(recipecalc.LoadBDORecipes(log, db))

	s := server.New(cfg, log, db, dspOptimizer, bdoOptimizer)
	go func() { log.Info(s.ListenAndServe()) }()

	// cronjobs.ScheduleTwitterMediaFetch()

	<-ctx.Done()
	_ = s.Shutdown(context.Background())
	log.Info("exiting")
}

package main

import (
	"alex-api/cronjobs"
	"alex-api/data"
	"alex-api/internal/config"
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
	service := server.NewService()
	db := data.New(log)
	defer db.Disconnect()
	s := server.New(cfg, log, service, db)
	go func() { log.Info(s.ListenAndServe()) }()

	cronjobs.ScheduleTwitterMediaFetch()

	<-ctx.Done()
	_ = s.Shutdown(context.Background())
	log.Info("exiting")
}

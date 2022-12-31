package main

import (
	"alex-api/cronjobs"
	"alex-api/internal/config"
	"alex-api/internal/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	cfg, err := config.New(".env")
	if err != nil {
		log.WithError(err).Fatal("failed to process config")
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		log.Info("termination signaled")
		cancel()
	}()

	log.Info("Starting up Alex API")
	service := server.NewService()
	s := server.New(cfg, log.WithContext(ctx), service)
	go func() { log.Info(s.ListenAndServe()) }()

	cronjobs.ScheduleTwitterMediaFetch()

	<-ctx.Done()
	_ = s.Shutdown(context.Background())
	log.Info("exiting")
}

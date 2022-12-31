package server

import (
	"alex-api/internal/config"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg config.Config
	*http.Server
	router    *mux.Router
	logger    *logrus.Logger
	service   Servicer
	validator *validator.Validate
}

func New(cfg config.Config, log *logrus.Logger, service Servicer) *Server {
	timeout := 60 * time.Second
	router := mux.NewRouter()
	router.StrictSlash(true)

	v := validator.New()
	server := &Server{
		Server: &http.Server{
			Handler:        router,
			Addr:           ":80",
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 65536,
		},
		router:    router,
		logger:    log,
		cfg:       cfg,
		validator: v,
		service:   service,
	}

	server.Route()

	return server
}

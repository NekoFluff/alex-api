package server

import (
	"alex-api/internal/config"
	"alex-api/utils"
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
	logger    *logrus.Entry
	service   Servicer
	validator *validator.Validate
	db        DB
}

func New(cfg config.Config, log *logrus.Entry, service Servicer, db DB) *Server {
	timeout := 60 * time.Second
	router := mux.NewRouter()
	router.StrictSlash(true)

	v := validator.New()
	server := &Server{
		Server: &http.Server{
			Handler:        router,
			Addr:           ":" + utils.GetEnvVar("PORT"),
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 65536,
		},
		router:    router,
		logger:    log,
		cfg:       cfg,
		validator: v,
		service:   service,
		db:        db,
	}

	server.Route()

	return server
}

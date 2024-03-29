package server

import (
	"alex-api/internal/config"
	"alex-api/internal/utils"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg config.Config
	*http.Server
	router       *mux.Router
	logger       *logrus.Entry
	validator    *validator.Validate
	db           DB
	dspOptimizer Optimizer
	bdoOptimizer Optimizer
}

func New(cfg config.Config, log *logrus.Entry, db DB, dspOptimizer Optimizer, bdoOptimizer Optimizer) *Server {
	timeout := 60 * time.Second
	router := mux.NewRouter()
	router.StrictSlash(true)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://nekofluff.github.io",
			"https://bdo-craft-profit.herokuapp.com",
			"https://wahtako.herokuapp.com",
			"https://www.wahtako.com",
			"https://www.alexnou.com",
			"https://alexnou.herokuapp.com",
			"https://stg-alexnou.herokuapp.com",
		},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	v := validator.New()
	server := &Server{
		Server: &http.Server{
			Handler:        c.Handler(router),
			Addr:           ":" + utils.GetEnvVar("PORT"),
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 65536,
		},
		router:       router,
		logger:       log,
		cfg:          cfg,
		validator:    v,
		db:           db,
		dspOptimizer: dspOptimizer,
		bdoOptimizer: bdoOptimizer,
	}

	server.Route()

	return server
}

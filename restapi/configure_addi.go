// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/carbocation/interpose/adaptors"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/dre1080/recovr"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	negronilogrus "github.com/meatballhat/negroni-logrus"

	"addi/restapi/operations"
)

//go:generate swagger generate server --target ../../src --name Addi --spec ../swagger.json --principal interface{}

func configureFlags(api *operations.AddiAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.AddiAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()
	api.TxtProducer = runtime.TextProducer()

	if api.GetDspItemsHandler == nil {
		api.GetDspItemsHandler = operations.GetDspItemsHandlerFunc(func(params operations.GetDspItemsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetDspItems has not yet been implemented")
		})
	}
	if api.PostDspHandler == nil {
		api.PostDspHandler = operations.PostDspHandlerFunc(func(params operations.PostDspParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostDsp has not yet been implemented")
		})
	}
	if api.PostDspItemsReloadHandler == nil {
		api.PostDspItemsReloadHandler = operations.PostDspItemsReloadHandlerFunc(func(params operations.PostDspItemsReloadParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostDspItemsReload has not yet been implemented")
		})
	}
	if api.CheckHealthHandler == nil {
		api.CheckHealthHandler = operations.CheckHealthHandlerFunc(func(params operations.CheckHealthParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.CheckHealth has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("A request was rejected: %v", r.URL)
	})
	return tollbooth.LimitFuncHandler(lmt, handler.ServeHTTP)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	// handleCORS := cors.Default().Handler
	credentials := handlers.AllowCredentials()
	ignoreOptions := handlers.IgnoreOptions()
	// methods := handlers.AllowedMethods([]string{"POST"})
	// // ttl := handlers.MaxAge(3600)
	// origins := handlers.AllowedOrigins([]string{"www.local.com"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS", "FETCH"})
	origins := handlers.AllowedOrigins([]string{"*"})
	handleCORS := handlers.CORS(credentials, methods, origins, ignoreOptions)

	recovery := recovr.New()
	negroniMiddleware := negronilogrus.NewMiddleware()
	logViaLogrus := adaptors.FromNegroni(negroniMiddleware)
	return recovery(logViaLogrus(handleCORS(handler)))
}

package main

import (
	"log"
	"strconv"

	"github.com/go-openapi/loads"

	"addi/handlers"
	"addi/restapi"
	"addi/restapi/operations"
	"addi/utils"
)

func main() {

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewAddiAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle
			log.Fatalln(err)
		}
	}()

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(handlers.CheckHealthHandler)

	api.PostDspHandler = operations.PostDspHandlerFunc(handlers.PostDSPHandler)

	api.GetDspItemsHandler = operations.GetDspItemsHandlerFunc(handlers.GetDSPItemsHandler)

	server.ConfigureAPI()

	port := utils.GetEnvVar("PORT")

	// Translate port string into int
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}

	server.Port = portInt

	// Start server
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"log"
	"strconv"

	"github.com/go-openapi/loads"

	"alex-api/cronjobs"
	"alex-api/handlers"
	"alex-api/restapi"
	"alex-api/restapi/operations"
	"alex-api/utils"
)

func main() {

	cronjobs.ScheduleTwitterMediaFetch()

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

	api.GetDSPRecipeHandler = operations.GetDSPRecipeHandlerFunc(handlers.GetDSPRecipeHandler)

	api.GetDSPRecipesHandler = operations.GetDSPRecipesHandlerFunc(handlers.GetDSPRecipesHandler)

	api.ReloadDSPRecipesHandler = operations.ReloadDSPRecipesHandlerFunc(handlers.ReloadDSPRecipesHandler)

	api.GetInArtHandler = operations.GetInArtHandlerFunc(handlers.GetInArtHandler)

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

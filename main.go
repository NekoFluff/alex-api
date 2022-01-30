package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

	"addi/models"
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

	// parser := flags.NewParser(server, flags.Default)
	// parser.ShortDescription = "go-rest-api"
	// parser.LongDescription = "HTTP server in Go with Swagger endpoints definition."
	// server.ConfigureFlags()
	// for _, optsGroup := range api.CommandLineOptionsGroups {
	// 	_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }

	// if _, err := parser.Parse(); err != nil {
	// 	code := 1
	// 	if fe, ok := err.(*flags.Error); ok {
	// 		if fe.Type == flags.ErrHelp {
	// 			code = 0
	// 		}
	// 	}
	// 	os.Exit(code)
	// }

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)

	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	api.PostDspHandler = operations.PostDspHandlerFunc(DSP)

	api.GetDspItemsHandler = operations.GetDspItemsHandlerFunc(GetDSPItems)

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

// variavel Global
var once sync.Once
var itemMap = make(map[string]*models.DSPItem)

func GetItem(itemName string) (*models.DSPItem, bool) {

	once.Do(func() {

		// Open up the file
		jsonFile, err := os.Open("data/items_arr.json")
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		// Read and unmarshal the file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var items []*models.DSPItem
		json.Unmarshal(byteValue, &items)

		// Map the items
		for _, v := range items {
			itemMap[v.Name] = v
		}
	})

	result, ok := itemMap[itemName]
	return result, ok
}

func init() {
	GetItem("asdf")
	log.Println("Initialized data")
}

func GetItems() []*models.DSPItem {
	m := make([]*models.DSPItem, 0, len(itemMap))
	for _, val := range itemMap {
		m = append(m, val)
	}
	return m
}

func GetRecipeForItem(itemName string, craftingSpeed float64, parentItemName string) []*models.Recipe {
	recipes := []*models.Recipe{}
	item, ok := GetItem(itemName)

	if ok {
		consumedMats := make(map[string]float64)
		numberOfFacilitiesNeeded := item.Time * craftingSpeed / item.Produce

		for _, material := range item.Materials {
			consumedMats[material.Name] = material.Count * numberOfFacilitiesNeeded / item.Time
		}

		recipe := &models.Recipe{
			Produce:                  item.Name,
			MadeIn:                   item.MadeIn,
			NumberOfFacilitiesNeeded: numberOfFacilitiesNeeded,
			ConsumesPerSec:           consumedMats,
			SecondsSpendPerCrafting:  item.Time,
			CraftingPerSecond:        craftingSpeed,
			For:                      parentItemName,
		}
		recipes = append(recipes, recipe)

		for _, material := range item.Materials {
			targetCraftingSpeed := numberOfFacilitiesNeeded * material.Count / item.Time
			materialRecipe := GetRecipeForItem(material.Name, targetCraftingSpeed, item.Name)
			recipes = append(recipes, materialRecipe...)
		}
	}

	return recipes
}

func GetDSPItems(params operations.GetDspItemsParams) middleware.Responder {
	return operations.NewGetDspItemsOK().WithPayload(GetItems())
}

func DSP(params operations.PostDspParams) middleware.Responder {

	// log.Println("Starting DSP Optimizer Program")

	recipe := []*models.Recipe{}

	for _, v := range params.RecipeRequest {
		recipe = append(recipe, GetRecipeForItem(*v.Name, float64(*v.Count), "")...)
	}

	// jsonStr, _ := json.MarshalIndent(recipe, "", "\t")
	// fmt.Println(string(jsonStr))

	return operations.NewPostDspOK().WithPayload(recipe)
}

//Health route returns OK
func Health(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

//GetHelloUser returns Hello + your name
func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}

//GetGopherByName returns a gopher in png
func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {

	var URL string
	if gopher.Name != "" {
		URL = "https://github.com/scraly/gophers/raw/main/" + gopher.Name + ".png"
	} else {
		//by default we return dr who gopher
		URL = "https://github.com/scraly/gophers/raw/main/dr-who.png"
	}

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}

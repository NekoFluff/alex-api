package main

import (
	"addi/models"
	"addi/restapi"
	"addi/restapi/operations"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func main() {

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
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

	server.Port = 8080

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)

	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	api.GetDspHandler = operations.GetDspHandlerFunc(DSP)

	// Start server which listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

type Item struct {
	Name      string     `json:"Name"`
	Produce   float64    `json:"Produce"`
	MadeIn    string     `json:"MadeIn"`
	Time      float64    `json:"Time"`
	Materials []Material `json:"Materials"`
}

type Material struct {
	Name  string  `json:"Name"`
	Count float64 `json:"Count"`
}

// variavel Global
var once sync.Once
var itemMap = make(map[string]Item)

func GetItem(itemName string) (Item, bool) {

	once.Do(func() {

		// Open up the file
		jsonFile, err := os.Open("data/items_arr.json")
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		// Read and unmarshal the file
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var items []Item
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
	fmt.Println(GetItem("asdf"))
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

func DSP(params operations.GetDspParams) middleware.Responder {

	log.Println("Starting DSP Optimizer Program")

	recipe := []*models.Recipe{}

	for _, v := range params.User {
		recipe = append(recipe, GetRecipeForItem(*v.Name, float64(*v.Count), "")...)
	}

	jsonStr, _ := json.MarshalIndent(recipe, "", "\t")
	fmt.Println(string(jsonStr))

	return operations.NewGetDspOK().WithPayload(recipe)
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

package dsp

import (
	"addi/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var once sync.Once
var items []*models.DSPRecipe
var itemMap = make(map[string][]*models.DSPRecipe)

func init() {
	GetItem("")
	log.Println("Initialized data")
}

func GetItem(itemName string) (*models.DSPRecipe, bool) {
	once.Do(func() {
		loadItemsFromFile()
	})

	result, ok := itemMap[itemName]
	if ok {
		return result[0], ok
	}

	return nil, false
}

func ReloadItems() {
	itemMap = make(map[string][]*models.DSPRecipe)
	loadItemsFromFile()
}

func loadItemsFromFile() {
	// Open up the file
	jsonFile, err := os.Open("data/items.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Read and unmarshal the file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &items)
	// Map the items
	for _, recipe := range items {
		// fmt.Println(recipe)

		if recipe.Name != nil {
			itemName := *recipe.Name
			itemMap[itemName] = append(itemMap[itemName], recipe)
		} else {
			fmt.Println(recipe)
			fmt.Println("Invalid Recipe Found")
			fmt.Println(*recipe.MadeIn)
			fmt.Println(*recipe.Materials[0].Name)
			fmt.Println(*recipe.Materials[0].Count)
		}
	}
}

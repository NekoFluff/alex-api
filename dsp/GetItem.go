package dsp

import (
	"addi/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var once sync.Once
var items []*models.DSPItem
var itemMap = make(map[string][]*models.DSPItem)

func init() {
	GetItem("")
	log.Println("Initialized data")
}

func GetItem(itemName string) (*models.DSPItem, bool) {
	once.Do(func() {
		loadItemsFromFile()
	})

	result, ok := itemMap[itemName]
	item := result[0]

	return item, ok
}

func ReloadItems() {
	itemMap = make(map[string][]*models.DSPItem)
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
	for _, v := range items {
		itemMap[v.Name] = append(itemMap[v.Name], v)
	}
}

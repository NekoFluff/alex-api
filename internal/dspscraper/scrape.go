package dspscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/NekoFluff/go-dsp/dsp"
	"github.com/gocolly/colly"
)

func Scrape() {
	urls := getURLs()
	// fmt.Println("URLS",urls)

	var dspRecipes []dsp.Recipe
	for itemName, url := range urls {
		dspRecipes = append(dspRecipes, scrapeURL(itemName, url)...)
	}

	file, _ := json.MarshalIndent(dspRecipes, "", "\t")

	_ = ioutil.WriteFile("data/items.json", file, 0644)
}

func getURLs() map[string]string {
	url := "https://dsp-wiki.com/Items"
	urls := make(map[string]string)

	c := colly.NewCollector(
		colly.AllowedDomains("dsp-wiki.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("div.item_icon_container a[href]", func(e *colly.HTMLElement) {
		urls[e.Attr("title")] = e.Request.AbsoluteURL(e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
	}
	return urls
}

func scrapeURL(itemName string, url string) []dsp.Recipe {
	var dspRecipes []dsp.Recipe
	c := colly.NewCollector()

	// c.Limit(&colly.LimitRule{
	// 	RandomDelay: 5 * time.Second,
	// })

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("table.pc_table:nth-of-type(1)", func(e *colly.HTMLElement) {
		e.ForEach("tr:nth-of-type(n+1)", func(_ int, e2 *colly.HTMLElement) {
			i := dsp.Recipe{}
			i.Materials = make(map[dsp.ItemName]float64)

			// Materials
			e2.ForEach("div.tt_recipe_item", func(_ int, e3 *colly.HTMLElement) {
				count, _ := strconv.ParseFloat(e3.ChildText("div"), 32)
				name := e3.ChildAttr("a", "title")
				i.Materials[dsp.ItemName(name)] = count
			})

			// Time Taken
			secondsStr := e2.ChildText("div.tt_rec_arrow")
			r, _ := regexp.Compile(`(\d+\.*\d*)`)
			secondsStr = r.FindString(secondsStr)
			time, _ := strconv.ParseFloat(secondsStr, 32)
			i.Time = time

			// Output
			e2.ForEach("div.tt_output_item", func(_ int, e3 *colly.HTMLElement) {
				outputItemName := e3.ChildAttr("a", "title")
				if itemName == outputItemName {
					i.OutputItem = dsp.ItemName(outputItemName)
					outputItemCount, _ := strconv.ParseFloat(e3.ChildText("div"), 64)
					i.OutputItemCount = float64(outputItemCount)
				}
			})

			// Made In
			e2.ForEach("td:nth-of-type(2)", func(_ int, e3 *colly.HTMLElement) {
				facility := e3.ChildAttr("a:nth-of-type(1)", "title")
				i.Facility = facility
			})

			fmt.Printf("Item: %+v\n", i)

			if i.OutputItem != "" {
				dspRecipes = append(dspRecipes, i)
			}
		})

	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
	}
	return dspRecipes
}

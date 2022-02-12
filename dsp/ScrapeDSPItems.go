package dsp

import (
	"addi/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

func ScrapeDSPItems() {
	urls := getDSPItemUrls()
	fmt.Println(urls)

	var dspItems []*models.DSPItem
	for itemName, url := range urls {
		dspItems = append(dspItems, handleDSPItemUrl(itemName, url)...)
	}

	file, _ := json.MarshalIndent(dspItems, "", "\t")

	_ = ioutil.WriteFile("data/items.json", file, 0644)
	ReloadItems()
}

func getDSPItemUrls() map[string]string {
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

	c.Visit(url)
	return urls
}

func handleDSPItemUrl(itemName string, url string) []*models.DSPItem {
	var dspItems []*models.DSPItem
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
			i := models.DSPItem{}

			// Materials
			e2.ForEach("div.tt_recipe_item", func(_ int, e3 *colly.HTMLElement) {
				m := models.DSPMaterial{}
				count, _ := strconv.ParseFloat(e3.ChildText("div"), 64)
				m.Count = &count
				name := e3.ChildAttr("a", "title")
				i.Name = &name
				i.Materials = append(i.Materials, &m)
				fmt.Printf("%+v\n", m)
			})

			// Time Taken
			secondsStr := e2.ChildText("div.tt_rec_arrow")
			r, _ := regexp.Compile(`(\d)+`)
			secondsStr = r.FindString(secondsStr)
			time, _ := strconv.ParseFloat(secondsStr, 64)
			i.Time = &time

			// Output
			e2.ForEach("div.tt_output_item", func(_ int, e3 *colly.HTMLElement) {
				outputItemName := e3.ChildAttr("a", "title")
				if itemName == outputItemName {
					i.Name = &outputItemName
					produce, _ := strconv.ParseFloat(e3.ChildText("div"), 64)
					i.Produce = &produce
				}
			})

			// Made In
			e2.ForEach("td:nth-of-type(2)", func(_ int, e3 *colly.HTMLElement) {
				madeIn := e3.ChildAttr("a:nth-of-type(1)", "title")
				i.MadeIn = &madeIn
			})

			fmt.Printf("%+v\n", i)

			if *i.Name != "" {
				dspItems = append(dspItems, &i)
			}
		})

	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	return dspItems
}

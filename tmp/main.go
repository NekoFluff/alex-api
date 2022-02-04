package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type Weapon struct {
	WeaponType   string
	MaxAmmo      float32
	MagazineSize float32
	Damage       float32
	Precision    float32
	Stagger      float32
	ROF          float32
	ReloadTime   float32
	Range        float32
}

type DSPItem struct {

	// made in
	MadeIn string `json:"madeIn,omitempty"`

	// materials
	Materials []*DSPMaterial `json:"materials"`

	// name
	Name string `json:"name,omitempty"`

	// produce
	Produce float64 `json:"produce,omitempty"`

	// time
	Time float64 `json:"time,omitempty"`
}

type DSPMaterial struct {

	// count
	Count float64 `json:"count,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

func main() {
	// scrapeGTFP()
	scrapeDSP()
}

func scrapeDSP() {
	urls := getDSPItemUrls()
	fmt.Println(urls)

	var dspItems []*DSPItem
	for itemName, url := range urls {
		dspItems = append(dspItems, handleDSPItemUrl(itemName, url)...)
	}

	file, _ := json.MarshalIndent(dspItems, "", "\t")

	_ = ioutil.WriteFile("test.json", file, 0644)
}

func handleDSPItemUrl(itemName string, url string) []*DSPItem {
	var dspItems []*DSPItem
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
			i := DSPItem{}
			// div.tt_recipe

			// Materials
			e2.ForEach("div.tt_recipe_item", func(_ int, e3 *colly.HTMLElement) {
				m := DSPMaterial{}
				m.Count, _ = strconv.ParseFloat(e3.ChildText("div"), 64)
				m.Name = e3.ChildAttr("a", "title")
				i.Materials = append(i.Materials, &m)
				fmt.Printf("%+v\n", m)
			})

			// Time Taken
			secondsStr := e2.ChildText("div.tt_rec_arrow")
			r, _ := regexp.Compile(`(\d)+`)
			secondsStr = r.FindString(secondsStr)
			i.Time, _ = strconv.ParseFloat(secondsStr, 64)

			// Output
			e2.ForEach("div.tt_output_item", func(_ int, e3 *colly.HTMLElement) {
				outputItemName := e3.ChildAttr("a", "title")
				if itemName == outputItemName {
					i.Name = outputItemName
					i.Produce, _ = strconv.ParseFloat(e3.ChildText("div"), 64)
				}
			})

			// Made In
			e2.ForEach("td:nth-of-type(2)", func(_ int, e3 *colly.HTMLElement) {
				i.MadeIn = e3.ChildAttr("a:nth-of-type(1)", "title")
			})

			fmt.Printf("%+v\n", i)

			if i.Name != "" {
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
		// urls = append(urls, e.Attr("href"))
		// e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		urls[e.Attr("title")] = e.Request.AbsoluteURL(e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	return urls
}

func scrapeGTFO() {
	urls := getWeaponUrls()
	fmt.Println(urls)
	// url := "https://gtfo.fandom.com/wiki/Buckland_S870_Shotgun"
	// handleWeaponUrl(url)

	for _, url := range urls {
		handleWeaponUrl("https://gtfo.fandom.com/" + url)
	}
}

func getWeaponUrls() []string {
	url := "https://gtfo.fandom.com/wiki/Weapons"

	c := colly.NewCollector(
		colly.AllowedDomains("gtfo.fandom.com"),
	)
	urls := []string{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		urls = append(urls, e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
	return urls
}

func handleWeaponUrl(url string) {
	c := colly.NewCollector(
		colly.AllowedDomains("gtfo.fandom.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// c.OnResponseHeaders(func(r *colly.Response) {
	// 	fmt.Println("Visited", r.Request.URL)
	// })

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		// fmt.Println(string(r.Body))
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("table.infoboxtable tr", func(e *colly.HTMLElement) {
		fmt.Println("First column of a table row:", e.ChildText("td:nth-of-type(1)"))
		fmt.Println("Second column of a table row:", e.ChildText("td:nth-of-type(2)"))
	})

	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(url)
}

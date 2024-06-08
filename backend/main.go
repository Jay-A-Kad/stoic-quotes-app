package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type Quote struct {
	ID        int    `json:"id"`
	Paragraph string `json:"paragraph"`
}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

func main() {
	quotes := Quotes{}
	quoteID := 1

	scrapeUrl := "https://mindofastoic.com/stoic-quotes"

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("mindofastoic.com"),
	)

	// On every <ol> element call callback
	c.OnHTML("ol", func(e *colly.HTMLElement) {
		// Iterate through each <li> inside the <ol>
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			quote := el.Text
			quotes.Quotes = append(quotes.Quotes, Quote{ID: quoteID, Paragraph: quote})
			quoteID++
		})
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(quotes); err != nil {
			fmt.Println("Error encoding JSON:", err)
		}
	})

	// Start scraping on the specified URL
	err := c.Visit(scrapeUrl)
	if err != nil {
		fmt.Println("Error visiting site:", err)
	}
}

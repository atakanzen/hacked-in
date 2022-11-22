package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

func NewScraper() {
	c := colly.NewCollector(
		// colly.AllowedDomains("google.com"),
		colly.AllowURLRevisit(),
		// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Something went wrong:\nResponse:\n%s\nError:\n%s", r.Body, err)
	})

	c.OnHTML("#rso", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(i int, h *colly.HTMLElement) {
			fmt.Printf("Got this: %s", h.Text)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.Visit("https://www.google.com/search?site%3Astackoverflow.com+javascript+variables")
}

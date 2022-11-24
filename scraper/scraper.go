package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

func NewScraper() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.google.com", "google.com", "consent.google.com"),
		colly.AllowURLRevisit(),
		// colly.UserAgent("firefox"),
		colly.Debugger(&debug.LogDebugger{}),
		// colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Got response from: %s", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Something went wrong:\nResponse:\n%s\nError:\n%s", r.Body, err)
	})

	c.OnHTML("div[id=res]", func(e *colly.HTMLElement) {
		fmt.Printf("Got: %s", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.Visit("https://www.google.com/search?q=site%3Astackoverflow.com+react+server+components")
}

package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

type hackerPosts struct {
	Title string `selector:".titleline > a"`
	URL   string `selector:".titleline > a" attr:"href"`
	// Score string `selector:".subline > .score"`
	// Age   string `selector:".subline > .age > a"`
}

func NewScraper() {
	posts := []hackerPosts{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.news.ycombinator.com", "news.ycombinator.com"),
		// colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (platform; rv:geckoversion) Gecko/geckotrail Firefox/firefoxversion"),
	// colly.Async(true),
	)

	// c.SetDebugger(&debug.LogDebugger{})

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

	c.OnHTML("tr.athing", func(e *colly.HTMLElement) {
		post := &hackerPosts{}
		e.Unmarshal(post)
		posts = append(posts, *post)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.Visit("https://news.ycombinator.com/")

	fmt.Printf("\n\n%v\n\n", posts)
}

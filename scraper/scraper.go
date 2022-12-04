package scraper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

type postMetadata struct {
	Score string `selector:".subline > .score"`
	Age   string `selector:".subline > .age > a"`
}

type hackerPost struct {
	Title    string `selector:".titleline > a"`
	URL      string `selector:".titleline > a" attr:"href"`
	Metadata postMetadata
}

func setupCollectors() (*colly.Collector, *colly.Collector) {
	mc := colly.NewCollector(
		colly.AllowedDomains("www.news.ycombinator.com", "news.ycombinator.com"),
		// colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (platform; rv:geckoversion) Gecko/geckotrail Firefox/firefoxversion"),
	// colly.Async(true),
	)

	sc := colly.NewCollector(
		colly.AllowedDomains("www.news.ycombinator.com", "news.ycombinator.com"),
		// colly.AllowURLRevisit(),
		colly.UserAgent("Mozilla/5.0 (platform; rv:geckoversion) Gecko/geckotrail Firefox/firefoxversion"),
	// colly.Async(true),
	)

	mc.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	sc.Limit(&colly.LimitRule{
		// Set a delay between requests to these domains
		Delay: 1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	return mc, sc
}

func ScrapePosts() *[]hackerPost {
	posts := []hackerPost{}
	postMetas := []postMetadata{}

	// Main Collector
	mc, sc := setupCollectors()

	mc.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error on Main Collector:\nResponse:\n%s\nError:\n%s", r.Body, err)
	})

	// Post Basic
	mc.OnHTML("tr.athing", func(h *colly.HTMLElement) {
		post := &hackerPost{}
		h.Unmarshal(post)
		posts = append(posts, *post)
	})

	sc.OnHTML("tr.athing + tr", func(h2 *colly.HTMLElement) {
		postMeta := &postMetadata{}
		h2.Unmarshal(postMeta)
		postMetas = append(postMetas, *postMeta)
	})

	mc.Visit("https://news.ycombinator.com/")
	sc.Visit("https://news.ycombinator.com/")

	for i, _ := range posts {
		posts[i].Metadata = postMetas[i]
	}

	return &posts
}

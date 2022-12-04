package main

import (
	"fmt"

	"github.com/atakanzen/stacked-in/scraper"
)

func main() {
	// p := tea.NewProgram(cli.InitialModel(), tea.WithAltScreen())
	// if err := p.Start(); err != nil {
	// 	fmt.Printf("Alas, there's been an error: %v", err)
	// 	os.Exit(1)
	// }
	posts := scraper.ScrapePosts()

	fmt.Printf("\n%+v\n", posts)
}

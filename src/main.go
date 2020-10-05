package main

import (
	"github.com/ythosa/linkschecker/src/links"
	"os"
)

func main() {
	worklist := make(chan []string)

	var n int // Number of waiting to be sent to the list

	n++

	go func() {
		worklist <- os.Args[1:] // Start with cmd arguments
	}()

	seen := make(map[string]bool)
	// Concurrency scan
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++

				go func(link string) {
					worklist <- links.Crawl(link)
				}(link)
			}
		}
	}
}

package main

import "os"

func main() {
	worklist := make(chan []string)
	var n int // Number of waiting to be sent to the list

	// Start with cmd arguments
	n++
	go func() {
		worklist <- os.Args[1:]
	}()

	// Concurrency scan
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- Crawl(link)
				}(link)
			}
		}
	}
}

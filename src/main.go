package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/ythosa/linkschecker/src/links"
)

func Crawl(link string, baseURL string, worklist *chan []string) error {
    var err error

    if strings.HasPrefix(link, baseURL) {
        go func(link string) {
            parsed, err := links.Extract(link)
            if err != nil {
                *worklist <- nil
                log.Println(err)
            }
            *worklist <- parsed
        }(link)
    }

    go func(link string) {
        _, _, err := links.Check(link)
        if err != nil {
            log.Println(err)
        }
        *worklist <- nil
    }(link)

    return err
}

func main() {
    worklist := make(chan []string)

    var baseURL string

    var n int // Number of waiting to be sent to the list

    if len(os.Args) == 1 {
        fmt.Println("Pls enter sthg :(")
        os.Exit(1)
    }

    baseURL = os.Args[1]
    n++

    go func() {
        worklist <- []string{baseURL} // Start with cmd arguments
    }()

    seen := make(map[string]bool)
    // Concurrency scan
    for ; n > 0; n-- {
        list := <-worklist
        if list == nil {
            continue
        }

        for _, link := range list {
            if !seen[link] {
                seen[link] = true

                n++

                if err := Crawl(link, baseURL, &worklist); err != nil {
                    log.Println(err)
                }
            }
        }
    }
}

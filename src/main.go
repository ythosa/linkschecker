package main

import (
    "fmt"
    "github.com/ythosa/linkschecker/src/links"
    "os"
    "strings"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("Please, pass something in arguments :(")
        os.Exit(1)
    }

    baseURL := os.Args[1]
    worklist := make(chan []links.ParsingURL)
    errlist := make(chan links.BadURL, 500)

    var n int

    n++
    go func() {
        worklist <- []links.ParsingURL{links.ParsingURL(baseURL)}
    }()

    go func() {
        for {
            fmt.Println((<-errlist).Err)
        }
    }()

    seen := make(map[links.ParsingURL]bool)
    for ; n > 0; n-- {
        slist := <-worklist
        if slist == nil {
            continue
        }

        for _, link := range slist {
            if !seen[link] {
                seen[link] = true
                n++
                go func(link links.ParsingURL) {
                    res, doc, err := links.CheckURL(link)
                    //fmt.Println(err)
                    if err != nil {
                        errlist <- links.BadURL{ParsingURL: link, Err: err}
                        worklist <- nil
                        return
                    }

                    if strings.HasPrefix(string(link), baseURL) {
                        worklist <- links.Extract(res, doc)
                    } else {
                        worklist <- nil
                    }
                }(link)
            }
        }
    }
}

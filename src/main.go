package main

import (
    "fmt"
    "github.com/ythosa/linkschecker/src/links"
    "os"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("Please, pass something in arguments :(")
        os.Exit(1)
    }

    baseURL := links.ParsingURL(os.Args[1])

    for _, errLink := range links.FindBrokenLinks(baseURL) {
        fmt.Println(errLink.Err)
    }
}

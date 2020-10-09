package main

import (
    "flag"
    "log"

    "github.com/ythosa/linkschecker/src/internal/app/apiserver"
)

func main() {
    flag.Parse()

    config := apiserver.NewConfig()

    if err := apiserver.Start(config); err != nil {
        log.Fatal(err)
    }
}


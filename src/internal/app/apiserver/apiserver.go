package apiserver

import (
    "log"
    "net/http"
)

// Start ...
func Start(config *Config) error {
    srv := newServer()

    log.Printf("Server has been started on PORT=%s\n", config.BindAddr)

    return http.ListenAndServe(config.BindAddr, srv)
}

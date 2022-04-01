package apiserver

import (
	"fmt"
	"log"
	"net/http"
)

// Start ...
func Start(config *Config) error {
	srv := newServer()

	log.Printf("Server has been started on PORT=%s\n", config.BindAddr)

	if err := http.ListenAndServe(config.BindAddr, srv); err != nil {
		return fmt.Errorf("error while starting server: %w", err)
	}

	return nil
}

package apiserver

import (
    "fmt"
    "os"
)

// Config ...
type Config struct {
    BindAddr string `toml:"bind_addr"`
    LogLevel string `toml:"log_level"`
}

// NewConfig ...
func NewConfig() *Config {
    return &Config{
        BindAddr: fmt.Sprintf(":%s", os.Getenv("PORT")),
        LogLevel: os.Getenv("LOG_LEVEL"),
    }
}

package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Useragent string
	Host      string
	Port      int
}

type AppConfig struct {
	Server Server `json:"server" yaml:"server"`

	// BadPatterns is a list of patterns for http request handler
	// see https://pkg.go.dev/github.com/gorilla/mux#Route.Path
	BadPaths        []string `json:"paths" yaml:"paths"`
	BadPathPrefixes []string `json:"path_prefixes" yaml:"path_prefixes"`
}

var DefaultConfig = &AppConfig{
	Server: Server{
		Useragent: "nginx/1.29.1",
		Host:      "0.0.0.0",
		Port:      8000,
	},
	BadPaths: []string{},
	BadPathPrefixes: []string{
		"/",
	},
}

func ReadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg AppConfig

	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

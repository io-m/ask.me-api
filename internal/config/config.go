package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	ArangoDB ArangoDBConfig
}

type ArangoDBConfig struct {
	Endpoint string
	Database string
	Username string
	Password string
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	arangoEndpoint := os.Getenv("ARANGO_ENDPOINT")
	if arangoEndpoint == "" {
		arangoEndpoint = "http://localhost:8529"
	}

	arangoDB := os.Getenv("ARANGO_DATABASE")
	if arangoDB == "" {
		return nil, fmt.Errorf("ARANGO_DATABASE environment variable is required")
	}

	return &Config{
		Port: port,
		ArangoDB: ArangoDBConfig{
			Endpoint: arangoEndpoint,
			Database: arangoDB,
			Username: os.Getenv("ARANGO_USERNAME"),
			Password: os.Getenv("ARANGO_PASSWORD"),
		},
	}, nil
}

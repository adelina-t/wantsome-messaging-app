package config

import "os"

type Config struct {
	URL  string
	Port string
}

func LoadConfig() Config {
	// Get url and port from environment variables
	url := os.Getenv("SERVER_URL")
	port := os.Getenv("SERVER_PORT")

	// If url or port are empty, use default values
	if url == "" {
		url = "localhost"
	}
	if port == "" {
		port = "8080"
	}

	return Config{
		URL:  url,
		Port: port,
	}
}

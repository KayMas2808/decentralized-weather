package main

import (
	"os"
	"strconv"
)

type Config struct {
	BackendURL         string
	SubmissionInterval int
	KeysPath           string
	DeviceLocation     string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		BackendURL:         getEnvOrDefault("BACKEND_URL", "http://localhost:8080"),
		SubmissionInterval: getEnvIntOrDefault("SUBMISSION_INTERVAL", 300),
		KeysPath:           getEnvOrDefault("KEYS_PATH", "./device_keys.json"),
		DeviceLocation:     getEnvOrDefault("DEVICE_LOCATION", "Unknown"),
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

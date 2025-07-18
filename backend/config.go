package main

import (
	"os"
	"strconv"
)

type Config struct {
	EthereumRPC             string
	PrivateKey              string
	DeviceRegistryAddr      string
	WeatherDataAddr         string
	RewardManagerAddr       string
	PinataAPIKey            string
	PinataSecretKey         string
	RateLimitWindow         int
	MaxSubmissionsPerWindow int
}

func LoadConfig() (*Config, error) {
	config := &Config{
		EthereumRPC:             getEnvOrDefault("ETHEREUM_RPC", "https://sepolia-rollup.arbitrum.io/rpc"),
		PrivateKey:              getEnvOrDefault("PRIVATE_KEY", ""),
		DeviceRegistryAddr:      getEnvOrDefault("DEVICE_REGISTRY_ADDRESS", ""),
		WeatherDataAddr:         getEnvOrDefault("WEATHER_DATA_ADDRESS", ""),
		RewardManagerAddr:       getEnvOrDefault("REWARD_MANAGER_ADDRESS", ""),
		PinataAPIKey:            getEnvOrDefault("PINATA_API_KEY", ""),
		PinataSecretKey:         getEnvOrDefault("PINATA_SECRET_KEY", ""),
		RateLimitWindow:         getEnvIntOrDefault("RATE_LIMIT_WINDOW", 3600),
		MaxSubmissionsPerWindow: getEnvIntOrDefault("MAX_SUBMISSIONS_PER_WINDOW", 12),
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

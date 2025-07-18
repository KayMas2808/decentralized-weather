package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	service, err := NewWeatherService(config)
	if err != nil {
		log.Fatalf("Failed to create weather service: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/register", service.RegisterDevice)
		api.POST("/submit", service.SubmitWeatherData)
		api.GET("/data", service.GetWeatherData)
		api.GET("/data/latest", service.GetLatestData)
		api.GET("/devices", service.GetDevices)
		api.GET("/health", service.HealthCheck)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting weather backend server on port %s", port)
	log.Fatal(r.Run(":" + port))
}

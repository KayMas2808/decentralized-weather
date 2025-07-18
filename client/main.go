package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting Weather Data Client...")

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := NewWeatherClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "register" {
		err = client.RegisterDevice()
		if err != nil {
			log.Fatalf("Failed to register device: %v", err)
		}
		fmt.Println("Device registered successfully!")
		return
	}

	err = client.LoadOrCreateKeys()
	if err != nil {
		log.Fatalf("Failed to load keys: %v", err)
	}

	ticker := time.NewTicker(time.Duration(config.SubmissionInterval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Client running with device ID: %x\n", client.DeviceID)
	fmt.Printf("Submitting data every %d seconds\n", config.SubmissionInterval)

	for {
		select {
		case <-ticker.C:
			err := client.SubmitWeatherData()
			if err != nil {
				log.Printf("Failed to submit weather data: %v", err)
			} else {
				fmt.Printf("Weather data submitted successfully at %s\n", time.Now().Format(time.RFC3339))
			}
		}
	}
}

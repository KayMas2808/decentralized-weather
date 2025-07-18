package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func (c *WeatherClient) generateMockWeatherData() WeatherData {
	rand.Seed(time.Now().UnixNano())

	windDirections := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

	return WeatherData{
		DeviceID:    hex.EncodeToString(c.DeviceID),
		Location:    c.Config.DeviceLocation,
		Temperature: 15.0 + rand.Float64()*20.0,
		Humidity:    30.0 + rand.Float64()*40.0,
		Pressure:    980.0 + rand.Float64()*40.0,
		WindSpeed:   rand.Float64() * 30.0,
		WindDir:     windDirections[rand.Intn(len(windDirections))],
		Timestamp:   time.Now(),
	}
}

func (c *WeatherClient) sendToBackend(payload SubmissionPayload) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(c.Config.BackendURL+"/submit", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("backend returned status %d: %s", resp.StatusCode, string(body))
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Backend response: %v\n", response)
	return nil
}

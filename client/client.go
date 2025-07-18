package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WeatherClient struct {
	Config     *Config
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	DeviceID   []byte
}

type DeviceKeys struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	DeviceID   string `json:"device_id"`
}

type WeatherData struct {
	DeviceID    string    `json:"device_id"`
	Location    string    `json:"location"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	WindSpeed   float64   `json:"wind_speed"`
	WindDir     string    `json:"wind_direction"`
	Timestamp   time.Time `json:"timestamp"`
}

type SubmissionPayload struct {
	WeatherData WeatherData `json:"weather_data"`
	DataHash    string      `json:"data_hash"`
	Signature   string      `json:"signature"`
	PublicKey   string      `json:"public_key"`
}

func NewWeatherClient(config *Config) (*WeatherClient, error) {
	return &WeatherClient{
		Config: config,
	}, nil
}

func (c *WeatherClient) LoadOrCreateKeys() error {
	keys, err := c.loadKeys()
	if err != nil {
		fmt.Println("Creating new device keys...")
		return c.generateAndSaveKeys()
	}

	privateKeyBytes, err := hex.DecodeString(keys.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %v", err)
	}

	c.PrivateKey, err = DeserializePrivateKey(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to deserialize private key: %v", err)
	}

	c.PublicKey = &c.PrivateKey.PublicKey

	deviceIDBytes, err := hex.DecodeString(keys.DeviceID)
	if err != nil {
		return fmt.Errorf("failed to decode device ID: %v", err)
	}
	c.DeviceID = deviceIDBytes

	fmt.Printf("Loaded existing device keys, device ID: %x\n", c.DeviceID)
	return nil
}

func (c *WeatherClient) SubmitWeatherData() error {
	weatherData := c.generateMockWeatherData()

	dataBytes, err := json.Marshal(weatherData)
	if err != nil {
		return fmt.Errorf("failed to marshal weather data: %v", err)
	}

	dataHash := sha256.Sum256(dataBytes)
	signature, err := c.signData(dataHash[:])
	if err != nil {
		return fmt.Errorf("failed to sign data: %v", err)
	}

	publicKeyBytes := SerializePublicKey(c.PublicKey)

	payload := SubmissionPayload{
		WeatherData: weatherData,
		DataHash:    hex.EncodeToString(dataHash[:]),
		Signature:   hex.EncodeToString(signature),
		PublicKey:   hex.EncodeToString(publicKeyBytes),
	}

	return c.sendToBackend(payload)
}

func (c *WeatherClient) RegisterDevice() error {
	err := c.LoadOrCreateKeys()
	if err != nil {
		return fmt.Errorf("failed to load keys: %v", err)
	}

	publicKeyBytes := SerializePublicKey(c.PublicKey)

	registrationData := map[string]string{
		"device_id":  hex.EncodeToString(c.DeviceID),
		"public_key": hex.EncodeToString(publicKeyBytes),
		"location":   c.Config.DeviceLocation,
	}

	payloadBytes, err := json.Marshal(registrationData)
	if err != nil {
		return fmt.Errorf("failed to marshal registration data: %v", err)
	}

	resp, err := http.Post(c.Config.BackendURL+"/register", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

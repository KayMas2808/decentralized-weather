package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

type WeatherService struct {
	Config     *Config
	EthClient  *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	Auth       *bind.TransactOpts

	submissionCounts map[string][]time.Time
	mu               sync.RWMutex
}

type DeviceRegistration struct {
	DeviceID  string `json:"device_id"`
	PublicKey string `json:"public_key"`
	Location  string `json:"location"`
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

type PinataResponse struct {
	IpfsHash string `json:"IpfsHash"`
}

func NewWeatherService(config *Config) (*WeatherService, error) {
	client, err := ethclient.Dial(config.EthereumRPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	var privateKey *ecdsa.PrivateKey
	var auth *bind.TransactOpts

	if config.PrivateKey != "" {
		privateKey, err = crypto.HexToECDSA(config.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}

		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %v", err)
		}

		auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		if err != nil {
			return nil, fmt.Errorf("failed to create transactor: %v", err)
		}
	}

	return &WeatherService{
		Config:           config,
		EthClient:        client,
		PrivateKey:       privateKey,
		Auth:             auth,
		submissionCounts: make(map[string][]time.Time),
	}, nil
}

func (s *WeatherService) RegisterDevice(c *gin.Context) {
	var registration DeviceRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data"})
		return
	}

	if registration.DeviceID == "" || registration.PublicKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID and public key are required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Device registration received",
		"device_id": registration.DeviceID,
		"status":    "pending_blockchain_confirmation",
	})
}

func (s *WeatherService) SubmitWeatherData(c *gin.Context) {
	var payload SubmissionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	deviceID := payload.WeatherData.DeviceID
	if !s.checkRateLimit(deviceID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
		return
	}

	if !s.verifySignature(payload) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	if !s.validateWeatherData(payload.WeatherData) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid weather data"})
		return
	}

	ipfsHash, err := s.uploadToPinata(payload.WeatherData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to IPFS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Weather data submitted successfully",
		"ipfs_hash": ipfsHash,
		"device_id": deviceID,
		"timestamp": time.Now(),
		"data_hash": payload.DataHash,
	})
}

func (s *WeatherService) GetWeatherData(c *gin.Context) {
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	mockData := s.generateMockHistoricalData(limit)
	c.JSON(http.StatusOK, gin.H{
		"data":  mockData,
		"count": len(mockData),
	})
}

func (s *WeatherService) GetLatestData(c *gin.Context) {
	mockData := s.generateMockHistoricalData(10)
	c.JSON(http.StatusOK, gin.H{
		"data": mockData,
	})
}

func (s *WeatherService) GetDevices(c *gin.Context) {
	mockDevices := []map[string]interface{}{
		{
			"device_id":         "0x1a2b3c4d5e6f",
			"location":          "New York, NY",
			"last_submission":   time.Now().Add(-time.Hour * 2),
			"total_submissions": 245,
			"status":            "active",
		},
		{
			"device_id":         "0x9a8b7c6d5e4f",
			"location":          "London, UK",
			"last_submission":   time.Now().Add(-time.Minute * 30),
			"total_submissions": 189,
			"status":            "active",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": mockDevices,
		"count":   len(mockDevices),
	})
}

func (s *WeatherService) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

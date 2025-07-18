package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

func (s *WeatherService) checkRateLimit(deviceID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-time.Duration(s.Config.RateLimitWindow) * time.Second)

	submissions := s.submissionCounts[deviceID]
	validSubmissions := make([]time.Time, 0)

	for _, submission := range submissions {
		if submission.After(windowStart) {
			validSubmissions = append(validSubmissions, submission)
		}
	}

	if len(validSubmissions) >= s.Config.MaxSubmissionsPerWindow {
		return false
	}

	validSubmissions = append(validSubmissions, now)
	s.submissionCounts[deviceID] = validSubmissions

	return true
}

func (s *WeatherService) verifySignature(payload SubmissionPayload) bool {
	publicKeyBytes, err := hex.DecodeString(payload.PublicKey)
	if err != nil {
		return false
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), publicKeyBytes)
	if x == nil || y == nil {
		return false
	}

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	dataBytes, err := json.Marshal(payload.WeatherData)
	if err != nil {
		return false
	}

	dataHash := sha256.Sum256(dataBytes)
	expectedHash := hex.EncodeToString(dataHash[:])

	if expectedHash != payload.DataHash {
		return false
	}

	signatureBytes, err := hex.DecodeString(payload.Signature)
	if err != nil {
		return false
	}

	if len(signatureBytes) != 64 {
		return false
	}

	r := new(big.Int).SetBytes(signatureBytes[:32])
	sigS := new(big.Int).SetBytes(signatureBytes[32:])

	return ecdsa.Verify(publicKey, dataHash[:], r, sigS)
}

func (s *WeatherService) validateWeatherData(data WeatherData) bool {
	if data.Temperature < -100 || data.Temperature > 70 {
		return false
	}

	if data.Humidity < 0 || data.Humidity > 100 {
		return false
	}

	if data.Pressure < 800 || data.Pressure > 1200 {
		return false
	}

	if data.WindSpeed < 0 || data.WindSpeed > 200 {
		return false
	}

	validDirections := map[string]bool{
		"N": true, "NE": true, "E": true, "SE": true,
		"S": true, "SW": true, "W": true, "NW": true,
	}

	if !validDirections[data.WindDir] {
		return false
	}

	timeDiff := time.Since(data.Timestamp)
	if timeDiff > time.Hour || timeDiff < -time.Minute*5 {
		return false
	}

	return true
}

func (s *WeatherService) uploadToPinata(data WeatherData) (string, error) {
	if s.Config.PinataAPIKey == "" {
		return s.generateMockIPFSHash(), nil
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile("file", "weather_data.json")
	if err != nil {
		return "", err
	}

	_, err = fileWriter.Write(dataBytes)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.pinata.cloud/pinning/pinFileToIPFS", &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", s.Config.PinataAPIKey)
	req.Header.Set("pinata_secret_api_key", s.Config.PinataSecretKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("pinata API error: %s", string(body))
	}

	var pinataResp PinataResponse
	err = json.NewDecoder(resp.Body).Decode(&pinataResp)
	if err != nil {
		return "", err
	}

	return pinataResp.IpfsHash, nil
}

func (s *WeatherService) generateMockIPFSHash() string {
	rand.Seed(time.Now().UnixNano())
	hash := make([]byte, 32)
	rand.Read(hash)
	return "Qm" + hex.EncodeToString(hash)[:44]
}

func (s *WeatherService) generateMockHistoricalData(count int) []map[string]interface{} {
	data := make([]map[string]interface{}, count)
	locations := []string{"New York, NY", "London, UK", "Tokyo, Japan", "Sydney, Australia", "Berlin, Germany"}
	devices := []string{"0x1a2b3c4d5e6f", "0x9a8b7c6d5e4f", "0x3f2e1d0c9b8a", "0x7e6d5c4b3a29"}

	for i := 0; i < count; i++ {
		data[i] = map[string]interface{}{
			"device_id":      devices[rand.Intn(len(devices))],
			"location":       locations[rand.Intn(len(locations))],
			"temperature":    15.0 + rand.Float64()*20.0,
			"humidity":       30.0 + rand.Float64()*40.0,
			"pressure":       980.0 + rand.Float64()*40.0,
			"wind_speed":     rand.Float64() * 30.0,
			"wind_direction": []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}[rand.Intn(8)],
			"timestamp":      time.Now().Add(-time.Duration(i*5) * time.Minute),
			"ipfs_hash":      s.generateMockIPFSHash(),
		}
	}

	return data
}

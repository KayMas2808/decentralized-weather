package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

func (c *WeatherClient) generateAndSaveKeys() error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	c.PrivateKey = privateKey
	c.PublicKey = &privateKey.PublicKey

	publicKeyBytes := SerializePublicKey(c.PublicKey)
	privateKeyBytes := SerializePrivateKey(c.PrivateKey)

	hash := sha256.Sum256(publicKeyBytes)
	c.DeviceID = hash[:16]

	keys := DeviceKeys{
		PrivateKey: hex.EncodeToString(privateKeyBytes),
		PublicKey:  hex.EncodeToString(publicKeyBytes),
		DeviceID:   hex.EncodeToString(c.DeviceID),
	}

	return c.saveKeys(keys)
}

func (c *WeatherClient) loadKeys() (*DeviceKeys, error) {
	data, err := os.ReadFile(c.Config.KeysPath)
	if err != nil {
		return nil, err
	}

	var keys DeviceKeys
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, err
	}

	return &keys, nil
}

func (c *WeatherClient) saveKeys(keys DeviceKeys) error {
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.Config.KeysPath, data, 0600)
}

func (c *WeatherClient) signData(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, c.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func SerializePublicKey(pubKey *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
}

func SerializePrivateKey(privKey *ecdsa.PrivateKey) []byte {
	return privKey.D.Bytes()
}

func DeserializePrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(data)

	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
		D: d,
	}

	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	return privKey, nil
}

package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

type Config struct {
	Port               string
	BaseURL            string
	TidalClientID      string
	GoogleClientID     string
	GoogleClientSecret string
	EncryptionKey      []byte
}

func Load() (*Config, error) {
	key, err := parseEncryptionKey(getEnv("ENCRYPTION_KEY", ""))
	if err != nil {
		return nil, fmt.Errorf("invalid ENCRYPTION_KEY: %w", err)
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080"),
		TidalClientID:      getEnv("TIDAL_CLIENT_ID", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		EncryptionKey:      key,
	}, nil
}

func (c *Config) Address() string {
	return ":" + c.Port
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseEncryptionKey(hexKey string) ([]byte, error) {
	if hexKey == "" {
		fmt.Fprintln(os.Stderr, "WARNING: ENCRYPTION_KEY not set, generating random key (sessions won't survive restarts)")
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("failed to generate random key: %w", err)
		}
		return key, nil
	}
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("must be hex-encoded: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("must be 32 bytes (64 hex chars), got %d bytes", len(key))
	}
	return key, nil
}

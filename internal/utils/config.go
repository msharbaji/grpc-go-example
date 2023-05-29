package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Endpoint  string
	KeyID     string
	SecretKey string
}

func LoadConfig() (*Config, error) {
	// Load Endpoint
	endpoint, err := loadEndpoint()
	if err != nil {
		return nil, err
	}

	// Load KeyID
	keyID, err := loadKeyID()
	if err != nil {
		return nil, err
	}

	// Load SecretKey
	secretKey, err := loadSecretKey()
	if err != nil {
		return nil, err
	}

	return &Config{
		Endpoint:  endpoint,
		KeyID:     keyID,
		SecretKey: secretKey,
	}, nil
}

// LoadEndpoint loads the endpoint from the environment variable GRPC_ENDPOINT.
func loadEndpoint() (string, error) {
	endpoint := os.Getenv("GRPC_ENDPOINT")
	if endpoint == "" {
		return "", errors.New("GRPC_ENDPOINT environment variable is not set")
	}

	if err := validateEndpoint(endpoint); err != nil {
		return "", fmt.Errorf("invalid GRPC_ENDPOINT: %w", err)
	}

	return endpoint, nil
}

// loadKeyID loads the key id from the environment variable KEY_ID.
func loadKeyID() (string, error) {
	if os.Getenv("KEY_ID") != "" {
		return os.Getenv("KEY_ID"), nil
	}
	return "", nil
}

// LoadSecretKey loads the secret key from the environment variable SECRET_KEY.
func loadSecretKey() (string, error) {
	if os.Getenv("SECRET_KEY") != "" {
		return os.Getenv("SECRET_KEY"), nil
	}
	return "", nil
}

// ValidateEndpoint validates the endpoint.
func validateEndpoint(endpoint string) error {
	parts := strings.Split(endpoint, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid endpoint: %s", endpoint)
	}

	return nil
}

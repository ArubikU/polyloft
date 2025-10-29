package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Credentials stores the user's authentication information
type Credentials struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

var (
	// ErrNotAuthenticated is returned when the user is not logged in
	ErrNotAuthenticated = errors.New("not authenticated")
)

// getCredentialsPath returns the path to the credentials file
func getCredentialsPath() (string, error) {
	home, err := resolveHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".polyloft")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "credentials.json"), nil
}

func resolveHomeDir() (string, error) {
	if custom := os.Getenv("POLYLOFT_HOME"); custom != "" {
		return filepath.Clean(custom), nil
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Clean(home), nil
	}
	if profile := os.Getenv("USERPROFILE"); profile != "" {
		return filepath.Clean(profile), nil
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return dir, nil
}

// SaveCredentials saves authentication credentials to disk
func SaveCredentials(creds *Credentials) error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// LoadCredentials loads authentication credentials from disk
func LoadCredentials() (*Credentials, error) {
	path, err := getCredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotAuthenticated
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	return &creds, nil
}

// ClearCredentials removes stored credentials
func ClearCredentials() error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	return nil
}

// IsAuthenticated checks if the user is currently authenticated
func IsAuthenticated() bool {
	creds, err := LoadCredentials()
	return err == nil && creds.Token != ""
}

// GetRegistryURL returns the registry URL from environment or default
func GetRegistryURL() string {
	if url := os.Getenv("POLYLOFT_REGISTRY_URL"); url != "" {
		return url
	}
	return "https://registry.polyloft.dev"
}

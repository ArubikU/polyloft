package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadCredentials(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving credentials
	creds := &Credentials{
		Username: "testuser",
		Token:    "test-token-123",
	}

	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials failed: %v", err)
	}

	// Verify file was created
	credPath := filepath.Join(tmpDir, ".polyloft", "credentials.json")
	if _, err := os.Stat(credPath); err != nil {
		t.Fatalf("Credentials file not created: %v", err)
	}

	// Test loading credentials
	loadedCreds, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials failed: %v", err)
	}

	if loadedCreds.Username != creds.Username {
		t.Errorf("Username mismatch: got %s, want %s", loadedCreds.Username, creds.Username)
	}

	if loadedCreds.Token != creds.Token {
		t.Errorf("Token mismatch: got %s, want %s", loadedCreds.Token, creds.Token)
	}
}

func TestIsAuthenticated(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Initially should not be authenticated
	if IsAuthenticated() {
		t.Error("Expected not authenticated initially")
	}

	// Save credentials
	creds := &Credentials{
		Username: "testuser",
		Token:    "test-token",
	}
	SaveCredentials(creds)

	// Now should be authenticated
	if !IsAuthenticated() {
		t.Error("Expected authenticated after saving credentials")
	}
}

func TestClearCredentials(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Save credentials
	creds := &Credentials{
		Username: "testuser",
		Token:    "test-token",
	}
	SaveCredentials(creds)

	// Clear credentials
	if err := ClearCredentials(); err != nil {
		t.Fatalf("ClearCredentials failed: %v", err)
	}

	// Should not be authenticated anymore
	if IsAuthenticated() {
		t.Error("Expected not authenticated after clearing credentials")
	}

	// Loading should return error
	_, err := LoadCredentials()
	if err != ErrNotAuthenticated {
		t.Errorf("Expected ErrNotAuthenticated, got: %v", err)
	}
}

func TestGetRegistryURL(t *testing.T) {
	// Test default URL
	originalEnv := os.Getenv("POLYLOFT_REGISTRY_URL")
	os.Unsetenv("POLYLOFT_REGISTRY_URL")
	defer func() {
		if originalEnv != "" {
			os.Setenv("POLYLOFT_REGISTRY_URL", originalEnv)
		}
	}()

	url := GetRegistryURL()
	if url != "https://registry.polyloft.dev" {
		t.Errorf("Expected default URL, got: %s", url)
	}

	// Test custom URL
	customURL := "https://custom.registry.com"
	os.Setenv("POLYLOFT_REGISTRY_URL", customURL)
	url = GetRegistryURL()
	if url != customURL {
		t.Errorf("Expected custom URL %s, got: %s", customURL, url)
	}
}

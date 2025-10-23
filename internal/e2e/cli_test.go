package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ArubikU/polyloft/internal/auth"
	"github.com/ArubikU/polyloft/internal/config"
	"github.com/ArubikU/polyloft/internal/publisher"
)

// TestRunCurrentDirectory tests the polyloft run command without file argument
func TestRunCurrentDirectory(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	
	// Create a simple polyloft.toml
	configContent := `[project]
name = "test-project"
version = "0.1.0"
entry_point = "main.pf"
`
	configPath := filepath.Join(tmpDir, "polyloft.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}
	
	// Create a simple main.pf file
	mainContent := `println("Test passed!")`
	mainPath := filepath.Join(tmpDir, "main.pf")
	if err := os.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to create main.pf: %v", err)
	}
	
	// Load the config
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify the entry point is set correctly
	if cfg.Project.EntryPoint != "main.pf" {
		t.Errorf("Expected entry_point 'main.pf', got '%s'", cfg.Project.EntryPoint)
	}
	
	// Verify the entry point file exists
	entryPath := filepath.Join(tmpDir, cfg.Project.EntryPoint)
	if _, err := os.Stat(entryPath); err != nil {
		t.Errorf("Entry point file does not exist: %v", err)
	}
}

// TestPublishValidation tests the publisher validation
func TestPublishValidation(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	
	// Test with missing entry point
	t.Run("MissingEntryPoint", func(t *testing.T) {
		configContent := `[project]
name = "test-project"
version = "0.1.0"
entry_point = "nonexistent.pf"
`
		configPath := filepath.Join(tmpDir, "test1.toml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to create config: %v", err)
		}
		
		cfg, err := config.Load(configPath)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// Use temporary home directory for credentials
		originalHome := os.Getenv("HOME")
		os.Setenv("HOME", tmpDir)
		defer os.Setenv("HOME", originalHome)
		
		// Create fake credentials
		creds := &auth.Credentials{
			Username: "testuser",
			Token:    "test-token",
		}
		auth.SaveCredentials(creds)
		
		// Try to publish - should fail validation
		pub := publisher.New(cfg)
		err = pub.Publish()
		if err == nil {
			t.Error("Expected error when entry point doesn't exist, got nil")
		}
	})
	
	// Test with valid configuration
	t.Run("ValidConfig", func(t *testing.T) {
		// Create entry point file
		mainPath := filepath.Join(tmpDir, "valid.pf")
		if err := os.WriteFile(mainPath, []byte("println('test')"), 0644); err != nil {
			t.Fatalf("Failed to create main.pf: %v", err)
		}
		
		configContent := `[project]
name = "test-project"
version = "1.0.0"
entry_point = "valid.pf"
`
		configPath := filepath.Join(tmpDir, "test2.toml")
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to create config: %v", err)
		}
		
		cfg, err := config.Load(configPath)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// Change to tmpDir so entry point can be found
		originalWd, _ := os.Getwd()
		os.Chdir(tmpDir)
		defer os.Chdir(originalWd)
		
		// Validation should pass (but publish will fail because server doesn't exist)
		pub := publisher.New(cfg)
		err = pub.Publish()
		// We expect this to fail at the upload stage, not validation
		if err != nil && err.Error() == "entry point file not found: valid.pf" {
			t.Error("Entry point validation should have passed")
		}
	})
}

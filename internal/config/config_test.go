package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "polyloft.toml")
	
	configContent := `[project]
name = "test-project"
version = "1.0.0"
entry_point = "src/main.pf"

[[dependencies.go]]
name = "github.com/example/lib"
version = "v1.0.0"

[[dependencies.pf]]
name = "math.vector"
version = "1.0.0"
`
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	// Validate the loaded config
	if cfg.Project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", cfg.Project.Name)
	}
	
	if cfg.Project.EntryPoint != "src/main.pf" {
		t.Errorf("Expected entry point 'src/main.pf', got '%s'", cfg.Project.EntryPoint)
	}
	
	if len(cfg.Dependencies.Go) != 1 {
		t.Errorf("Expected 1 Go dependency, got %d", len(cfg.Dependencies.Go))
	}
	
	if len(cfg.Dependencies.Pf) != 1 {
		t.Errorf("Expected 1 Polyloft dependency, got %d", len(cfg.Dependencies.Pf))
	}
}

func TestLoadMissingEntryPoint(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "polyloft.toml")
	
	configContent := `[project]
name = "test-project"
version = "1.0.0"
`
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}
	
	// Load should fail without entry_point
	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for missing entry_point, got nil")
	}
}

func TestExample(t *testing.T) {
	example := Example()
	if example == "" {
		t.Error("Example should not be empty")
	}
	
	// Verify example contains required sections
	if !contains(example, "[project]") {
		t.Error("Example should contain [project] section")
	}
	
	if !contains(example, "entry_point") {
		t.Error("Example should contain entry_point field")
	}
	
	if !contains(example, "[[dependencies.go]]") {
		t.Error("Example should contain Go dependencies section")
	}
	
	if !contains(example, "[[dependencies.pf]]") {
		t.Error("Example should contain Polyloft dependencies section")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

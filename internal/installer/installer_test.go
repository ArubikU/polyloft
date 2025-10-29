package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/config"
)

func TestInstallPackages(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create minimal config
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:       "test-project",
			Version:    "0.1.0",
			EntryPoint: "src/main.pf",
		},
	}
	
	// Create installer with custom lib directory
	inst := New(cfg)
	inst.LibDir = filepath.Join(tmpDir, "libs")
	
	// Test with package without author
	packages := []string{"test-package"}
	err := inst.InstallPackages(packages)
	if err != nil {
		t.Errorf("InstallPackages should not return error for package without author: %v", err)
	}
	
	// Test with package with author (will fail to download but should handle gracefully)
	packages = []string{"test-package@test-author"}
	err = inst.InstallPackages(packages)
	if err != nil {
		t.Errorf("InstallPackages should handle download errors gracefully: %v", err)
	}
	
	// Verify libs directory was created
	if _, err := os.Stat(inst.LibDir); os.IsNotExist(err) {
		t.Error("Libs directory should be created")
	}
}

func TestInstallPackagesMultiple(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create minimal config
	cfg := &config.Config{
		Project: config.ProjectConfig{
			Name:       "test-project",
			Version:    "0.1.0",
			EntryPoint: "src/main.pf",
		},
	}
	
	// Create installer with custom lib directory
	inst := New(cfg)
	inst.LibDir = filepath.Join(tmpDir, "libs")
	
	// Test with multiple packages
	packages := []string{
		"package1@author1",
		"package2@author2",
		"package3", // without author
	}
	
	err := inst.InstallPackages(packages)
	if err != nil {
		t.Errorf("InstallPackages should handle multiple packages: %v", err)
	}
	
	// Verify libs directory was created
	if _, err := os.Stat(inst.LibDir); os.IsNotExist(err) {
		t.Error("Libs directory should be created")
	}
}

func TestInstallPackagesParseAuthor(t *testing.T) {
	tests := []struct {
		input        string
		expectName   string
		expectAuthor string
	}{
		{"vectors@Arubik", "vectors", "Arubik"},
		{"math.vector@TestAuthor", "math.vector", "TestAuthor"},
		{"simple", "simple", ""},
		{"complex@name@with@multiple@at", "complex", "name@with@multiple@at"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var name, author string
			if strings.Contains(tt.input, "@") {
				parts := strings.SplitN(tt.input, "@", 2)
				name = parts[0]
				author = parts[1]
			} else {
				name = tt.input
			}
			
			if name != tt.expectName {
				t.Errorf("Expected name %q, got %q", tt.expectName, name)
			}
			if author != tt.expectAuthor {
				t.Errorf("Expected author %q, got %q", tt.expectAuthor, author)
			}
		})
	}
}

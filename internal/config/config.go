package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the polyloft.toml configuration file structure
type Config struct {
	Project      ProjectConfig      `toml:"project"`
	Dependencies DependenciesConfig `toml:"dependencies"`
}

// ProjectConfig contains project-level settings
type ProjectConfig struct {
	EntryPoint string `toml:"entry_point"`
	Name       string `toml:"name"`
	Version    string `toml:"version"`
}

// DependenciesConfig contains both Go and Polyloft library dependencies
type DependenciesConfig struct {
	Go []GoDependency `toml:"go"`
	Pf []PfDependency `toml:"pf"`
}

// GoDependency represents a Go library dependency for compatibility
type GoDependency struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
}

// PfDependency represents a Polyloft library dependency
type PfDependency struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Source  string `toml:"source,omitempty"` // Optional: custom source URL
}

// Load reads and parses a polyloft.toml configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse TOML config: %w", err)
	}

	// Validate required fields
	if cfg.Project.EntryPoint == "" {
		return nil, fmt.Errorf("project.entry_point is required in config")
	}

	return &cfg, nil
}

// LoadDefault attempts to load polyloft.toml from the current directory
func LoadDefault() (*Config, error) {
	return Load("polyloft.toml")
}

// Example generates an example polyloft.toml configuration
func Example() string {
	return `[project]
name = "my-polyloft-project"
version = "0.1.0"
entry_point = "src/main.pf"

[[dependencies.go]]
name = "github.com/example/library"
version = "v1.0.0"

[[dependencies.pf]]
name = "math.vector"
version = "1.0.0"

[[dependencies.pf]]
name = "utils"
version = "1.0.0"
source = "https://polyloft-registry.example.com/utils"
`
}

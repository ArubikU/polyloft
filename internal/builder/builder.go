package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ArubikU/polyloft/internal/config"
)

// Builder handles the compilation of Hy source to executables
type Builder struct {
	Config     *config.Config
	OutputPath string
}

// New creates a new Builder with the given configuration
func New(cfg *config.Config, outputPath string) *Builder {
	return &Builder{
		Config:     cfg,
		OutputPath: outputPath,
	}
}

// Build compiles the Hy project to an executable
func (b *Builder) Build() error {
	fmt.Println("[build] Starting build process...")
	
	// Verify entry point exists
	if _, err := os.Stat(b.Config.Project.EntryPoint); os.IsNotExist(err) {
		return fmt.Errorf("entry point not found: %s", b.Config.Project.EntryPoint)
	}

	fmt.Printf("[build] Entry point: %s\n", b.Config.Project.EntryPoint)
	
	// Make output path absolute
	absOutput, err := filepath.Abs(b.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}
	b.OutputPath = absOutput
	fmt.Printf("[build] Output: %s\n", b.OutputPath)

	// Create a temporary directory for build artifacts
	tmpDir, err := os.MkdirTemp("", "polyloft-build-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Generate Go wrapper code
	wrapperPath := filepath.Join(tmpDir, "main.go")
	if err := b.generateGoWrapper(wrapperPath); err != nil {
		return fmt.Errorf("failed to generate Go wrapper: %w", err)
	}

	// Copy polyloft source and dependencies
	if err := b.copySourceFiles(tmpDir); err != nil {
		return fmt.Errorf("failed to copy source files: %w", err)
	}

	// Build the Go executable
	if err := b.compileGoWrapper(tmpDir); err != nil {
		return fmt.Errorf("failed to compile: %w", err)
	}

	fmt.Printf("[build] Successfully built executable: %s\n", b.OutputPath)
	return nil
}

// generateGoWrapper creates a Go main.go that embeds and runs the Hy code
func (b *Builder) generateGoWrapper(outputPath string) error {
	entryPoint := b.Config.Project.EntryPoint
	
	// Read the entry point source
	sourceData, err := os.ReadFile(entryPoint)
	if err != nil {
		return err
	}

	// Escape the source for embedding in Go string
	escapedSource := strings.ReplaceAll(string(sourceData), "`", "` + \"`\" + `")

	goCode := fmt.Sprintf(`package main

import (
	"os"
	"github.com/ArubikU/polyloft/pkg/runtime"
)

const embeddedSource = %s

func main() {
	if err := runtime.ExecuteSource(embeddedSource, "%s"); err != nil {
		os.Exit(1)
	}
}
`, "`"+escapedSource+"`", entryPoint)

	return os.WriteFile(outputPath, []byte(goCode), 0644)
}

// copySourceFiles copies necessary source files to the build directory
func (b *Builder) copySourceFiles(buildDir string) error {
	// For now, we embed the source directly in the wrapper
	// In the future, we could copy library files here
	return nil
}

// compileGoWrapper compiles the generated Go code to an executable
func (b *Builder) compileGoWrapper(buildDir string) error {
	fmt.Println("[build] Compiling Go executable...")
	
	// Try to find a local polyloft module first (for development)
	// If not found, will try to download from remote
	var modPath string
	var useLocalModule bool
	
	// Try to locate polyloft module in development environment
	exePath, err := os.Executable()
	if err == nil {
		// Navigate up to find go.mod
		searchPath := filepath.Dir(exePath)
		for i := 0; i < 5; i++ { // Limit search depth
			goModPath := filepath.Join(searchPath, "go.mod")
			if data, err := os.ReadFile(goModPath); err == nil {
				// Check if this is the polyloft module
				if strings.Contains(string(data), "module github.com/ArubikU/polyloft") {
					modPath = searchPath
					useLocalModule = true
					break
				}
			}
			parent := filepath.Dir(searchPath)
			if parent == searchPath {
				break
			}
			searchPath = parent
		}
	}

	// Initialize go.mod in the temp directory
	initCmd := exec.Command("go", "mod", "init", "polyloft-build")
	initCmd.Dir = buildDir
	if err := initCmd.Run(); err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}

	// If we found a local module, use replace directive
	// Otherwise, the module should be available via go get
	if useLocalModule {
		fmt.Printf("[build] Using local polyloft module from: %s\n", modPath)
		replaceCmd := exec.Command("go", "mod", "edit", 
			"-replace", fmt.Sprintf("github.com/ArubikU/polyloft=%s", modPath))
		replaceCmd.Dir = buildDir
		if err := replaceCmd.Run(); err != nil {
			return fmt.Errorf("go mod edit replace failed: %w", err)
		}
		
		// Use v0.0.0 for local development
		requireCmd := exec.Command("go", "mod", "edit", 
			"-require", "github.com/ArubikU/polyloft@v0.0.0")
		requireCmd.Dir = buildDir
		if err := requireCmd.Run(); err != nil {
			return fmt.Errorf("go mod edit require failed: %w", err)
		}
	} else {
		fmt.Println("[build] Using polyloft module from Go module cache")
		// For published version, use @latest or specific version
		requireCmd := exec.Command("go", "mod", "edit", 
			"-require", "github.com/ArubikU/polyloft@latest")
		requireCmd.Dir = buildDir
		if err := requireCmd.Run(); err != nil {
			return fmt.Errorf("go mod edit require failed: %w", err)
		}
	}

	// Tidy dependencies
	fmt.Println("[build] Downloading dependencies...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = buildDir
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Build the executable
	fmt.Println("[build] Building executable...")
	buildCmd := exec.Command("go", "build", "-o", b.OutputPath, ".")
	buildCmd.Dir = buildDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}

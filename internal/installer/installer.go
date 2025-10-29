package installer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ArubikU/polyloft/internal/auth"
	"github.com/ArubikU/polyloft/internal/config"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Installer handles dependency installation
type Installer struct {
	Config          *config.Config
	LibDir          string
	GlobalMode      bool
	installed       map[string]bool // Track installed packages to avoid duplicates
	dependencyChain []string        // Track dependency chain to detect cycles
}

// New creates a new Installer with the given configuration
func New(cfg *config.Config) *Installer {
	return &Installer{
		Config:          cfg,
		LibDir:          "libs", // Default library directory
		GlobalMode:      false,
		installed:       make(map[string]bool),
		dependencyChain: []string{},
	}
}

// SetGlobalMode enables global installation mode
func (i *Installer) SetGlobalMode(global bool) {
	i.GlobalMode = global
	if global {
		// Use global library directory in user's home
		homeDir, err := os.UserHomeDir()
		if err == nil {
			i.LibDir = filepath.Join(homeDir, ".polyloft", "libs")
		}
	}
}

// Install downloads and installs all dependencies
func (i *Installer) Install() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	
	fmt.Printf("\n%s Installing dependencies...\n", cyan("📦"))

	// Install Go dependencies
	if err := i.installGoDependencies(); err != nil {
		return fmt.Errorf("failed to install Go dependencies: %w", err)
	}

	// Install Polyloft dependencies
	if err := i.installHyDependencies(); err != nil {
		return fmt.Errorf("failed to install Polyloft dependencies: %w", err)
	}

	fmt.Printf("\n%s All dependencies installed successfully\n\n", green("✓"))
	return nil
}

// InstallPackages installs specific packages from command line arguments
func (i *Installer) InstallPackages(packages []string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	
	fmt.Printf("\n%s Installing %d package(s)...\n", cyan("📦"), len(packages))
	
	// Ensure libs directory exists
	if err := os.MkdirAll(i.LibDir, 0755); err != nil {
		return fmt.Errorf("failed to create libs directory: %w", err)
	}
	
	for _, pkg := range packages {
		// Parse package name for @author syntax
		var name, author string
		if strings.Contains(pkg, "@") {
			parts := strings.SplitN(pkg, "@", 2)
			name = parts[0]
			author = parts[1]
		} else {
			name = pkg
			// If no author specified, we can't download from registry
			fmt.Printf("  %s Package %s has no author specified. Use format: package@author\n", yellow("⚠"), pkg)
			fmt.Printf("  %s Skipping %s\n", yellow("→"), pkg)
			continue
		}
		
		packageKey := fmt.Sprintf("%s@%s", name, author)
		
		// Check if already installed
		if i.installed[packageKey] {
			fmt.Printf("  %s %s already processed, skipping\n", green("✓"), packageKey)
			continue
		}
		
		fmt.Printf("\n  %s Installing %s...\n", cyan("→"), packageKey)
		
		libPath := filepath.Join(i.LibDir, name)
		
		// Check if library already exists
		if _, err := os.Stat(libPath); err == nil {
			fmt.Printf("    %s %s already exists\n", green("✓"), pkg)
			i.installed[packageKey] = true
			// Still check for transitive dependencies
			if err := i.installTransitiveDependencies(libPath, packageKey); err != nil {
				fmt.Printf("    %s Warning: Failed to install transitive dependencies: %v\n", yellow("⚠"), err)
			}
			continue
		}
		
		// Download from registry with spinner
		if err := i.downloadPackageWithAnimation(name, author, "", libPath); err != nil {
			fmt.Printf("    %s Failed to download %s: %v\n", red("✗"), packageKey, err)
			continue
		}
		
		i.installed[packageKey] = true
		fmt.Printf("    %s Successfully installed %s\n", green("✓"), packageKey)
		
		// Install transitive dependencies
		if err := i.installTransitiveDependencies(libPath, packageKey); err != nil {
			fmt.Printf("    %s Warning: Failed to install transitive dependencies: %v\n", yellow("⚠"), err)
		}
	}
	
	fmt.Printf("\n%s Package installation complete\n\n", green("✓"))
	return nil
}

// installTransitiveDependencies installs dependencies of dependencies
func (i *Installer) installTransitiveDependencies(packagePath, packageKey string) error {
	// Check for cyclic dependencies
	for _, dep := range i.dependencyChain {
		if dep == packageKey {
			return fmt.Errorf("cyclic dependency detected: %s", strings.Join(append(i.dependencyChain, packageKey), " -> "))
		}
	}
	
	// Add current package to chain
	i.dependencyChain = append(i.dependencyChain, packageKey)
	defer func() {
		// Remove from chain when done
		i.dependencyChain = i.dependencyChain[:len(i.dependencyChain)-1]
	}()
	
	// Look for polyloft.toml in the package directory
	configPath := filepath.Join(packagePath, "polyloft.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file, no transitive dependencies
		return nil
	}
	
	// Load the package's config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load package config: %w", err)
	}
	
	// Check if there are Polyloft dependencies
	if len(cfg.Dependencies.Pf) == 0 {
		return nil
	}
	
	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("    %s Installing %d transitive dependencies...\n", cyan("→"), len(cfg.Dependencies.Pf))
	
	// Install each dependency
	for _, dep := range cfg.Dependencies.Pf {
		var name, author string
		if strings.Contains(dep.Name, "@") {
			parts := strings.SplitN(dep.Name, "@", 2)
			name = parts[0]
			author = parts[1]
		} else {
			name = dep.Name
			// Try to infer from package structure or skip
			continue
		}
		
		transKey := fmt.Sprintf("%s@%s", name, author)
		
		// Check if already installed
		if i.installed[transKey] {
			continue
		}
		
		libPath := filepath.Join(i.LibDir, name)
		
		// Check if library already exists
		if _, err := os.Stat(libPath); err == nil {
			i.installed[transKey] = true
			// Recursively check this package's dependencies
			if err := i.installTransitiveDependencies(libPath, transKey); err != nil {
				return err
			}
			continue
		}
		
		// Download the transitive dependency
		if err := i.downloadPackageWithAnimation(name, author, dep.Version, libPath); err != nil {
			return fmt.Errorf("failed to download transitive dependency %s: %w", transKey, err)
		}
		
		i.installed[transKey] = true
		
		// Recursively install its dependencies
		if err := i.installTransitiveDependencies(libPath, transKey); err != nil {
			return err
		}
	}
	
	return nil
}

// downloadPackageWithAnimation downloads a package with a nice spinner animation
func (i *Installer) downloadPackageWithAnimation(name, author, version, destPath string) error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Downloading %s@%s...", name, author)
	s.Start()
	defer s.Stop()
	
	err := i.downloadFromRegistry(name, author, version, destPath)
	return err
}

// installGoDependencies installs Go library dependencies
func (i *Installer) installGoDependencies() error {
	if len(i.Config.Dependencies.Go) == 0 {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("  %s No Go dependencies to install\n", cyan("ℹ"))
		return nil
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n  %s Installing %d Go dependencies...\n", cyan("→"), len(i.Config.Dependencies.Go))

	for _, dep := range i.Config.Dependencies.Go {
		fmt.Printf("    %s %s@%s\n", cyan("→"), dep.Name, dep.Version)
		
		// Use go get to install the dependency
		cmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", dep.Name, dep.Version))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", dep.Name, err)
		}
		
		fmt.Printf("    %s Installed %s\n", green("✓"), dep.Name)
	}

	// Run go mod tidy to clean up
	fmt.Printf("    %s Running go mod tidy...\n", cyan("→"))
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	return nil
}

// installHyDependencies installs Polyloft library dependencies
func (i *Installer) installHyDependencies() error {
	if len(i.Config.Dependencies.Pf) == 0 {
		cyan := color.New(color.FgCyan).SprintFunc()
		fmt.Printf("  %s No Polyloft dependencies to install\n", cyan("ℹ"))
		return nil
	}

	cyan := color.New(color.FgCyan).SprintFunc()
	fmt.Printf("\n  %s Installing %d Polyloft dependencies...\n", cyan("→"), len(i.Config.Dependencies.Pf))

	// Ensure libs directory exists
	if err := os.MkdirAll(i.LibDir, 0755); err != nil {
		return fmt.Errorf("failed to create libs directory: %w", err)
	}

	for _, dep := range i.Config.Dependencies.Pf {
		if err := i.installPfDependency(dep); err != nil {
			return fmt.Errorf("failed to install %s: %w", dep.Name, err)
		}
	}

	return nil
}

// installPfDependency installs a single Polyloft library dependency
func (i *Installer) installPfDependency(dep config.PfDependency) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	
	// Parse package name for @author syntax
	var name, author string
	if strings.Contains(dep.Name, "@") {
		parts := strings.SplitN(dep.Name, "@", 2)
		name = parts[0]
		author = parts[1]
	} else {
		name = dep.Name
	}
	
	packageKey := fmt.Sprintf("%s@%s", name, author)
	
	libPath := filepath.Join(i.LibDir, name)
	
	// Check if library already exists
	if _, err := os.Stat(libPath); err == nil {
		if !i.installed[packageKey] {
			fmt.Printf("    %s %s already exists\n", green("✓"), dep.Name)
			i.installed[packageKey] = true
			// Check for transitive dependencies
			if err := i.installTransitiveDependencies(libPath, packageKey); err != nil {
				fmt.Printf("    %s Warning: %v\n", yellow("⚠"), err)
			}
		}
		return nil
	}

	// Try to download from registry if author is specified
	if author != "" {
		fmt.Printf("    %s Downloading %s...\n", color.CyanString("→"), packageKey)
		if err := i.downloadPackageWithAnimation(name, author, dep.Version, libPath); err != nil {
			fmt.Printf("    %s Warning: Failed to download from registry: %v\n", yellow("⚠"), err)
			return nil // Don't fail the install, just warn
		}
		i.installed[packageKey] = true
		fmt.Printf("    %s Successfully installed %s\n", green("✓"), packageKey)
		
		// Install transitive dependencies
		if err := i.installTransitiveDependencies(libPath, packageKey); err != nil {
			fmt.Printf("    %s Warning: %v\n", yellow("⚠"), err)
		}
		return nil
	}

	// If source is specified, download from there
	if dep.Source != "" {
		fmt.Printf("    %s Source: %s\n", color.CyanString("→"), dep.Source)
		if err := i.downloadFromSource(dep.Source, libPath); err != nil {
			fmt.Printf("    %s Warning: Failed to download from source: %v\n", yellow("⚠"), err)
			return nil
		}
		return nil
	}

	// For local development, just verify the library exists somewhere
	fmt.Printf("    %s Warning: Library %s not found locally. Ensure it exists in %s/\n", yellow("⚠"), dep.Name, i.LibDir)
	
	return nil
}

// downloadFromRegistry downloads a package from the Polyloft registry
func (i *Installer) downloadFromRegistry(name, author, version, destPath string) error {
	registryURL := auth.GetRegistryURL()
	
	// Construct download URL
	var downloadURL string
	if version != "" {
		downloadURL = fmt.Sprintf("%s/api/download/%s/%s/%s", registryURL, author, name, version)
	} else {
		downloadURL = fmt.Sprintf("%s/api/download/%s/%s", registryURL, author, name)
	}
	
	// Download package archive
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download package: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	// Read archive data
	archiveData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read package data: %w", err)
	}
	
	// Extract archive
	if err := i.extractArchive(archiveData, destPath); err != nil {
		return fmt.Errorf("failed to extract package: %w", err)
	}
	
	return nil
}

// extractArchive extracts a tar.gz archive to the destination path
func (i *Installer) extractArchive(archiveData []byte, destPath string) error {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Create gzip reader
	gzReader, err := gzip.NewReader(bytes.NewReader(archiveData))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()
	
	// Create tar reader
	tarReader := tar.NewReader(gzReader)
	
	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}
		
		// Construct target path
		targetPath := filepath.Join(destPath, header.Name)
		
		// Handle directories
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}
		
		// Handle files
		if header.Typeflag == tar.TypeReg {
			// Create parent directories
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}
			
			// Create file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			
			// Copy file contents
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()
		}
	}
	
	return nil
}

// downloadFromSource downloads a package from a custom source URL
func (i *Installer) downloadFromSource(source, destPath string) error {
	// TODO: Implement downloading from custom sources (git, http, etc.)
	fmt.Printf("[install]     Custom source download not yet implemented\n")
	return nil
}

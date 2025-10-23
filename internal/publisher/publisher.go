package publisher

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArubikU/polyloft/internal/auth"
	"github.com/ArubikU/polyloft/internal/config"
)

// Publisher handles publishing packages to the registry
type Publisher struct {
	cfg         *config.Config
	registryURL string
}

// New creates a new publisher
func New(cfg *config.Config) *Publisher {
	return &Publisher{
		cfg:         cfg,
		registryURL: auth.GetRegistryURL(),
	}
}

// Publish publishes the current package to the registry
func (p *Publisher) Publish() error {
	// Check authentication
	creds, err := auth.LoadCredentials()
	if err != nil {
		return fmt.Errorf("not authenticated. Please run 'polyloft login' first")
	}

	// Validate configuration
	if err := p.validateConfig(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Println("📦 Packaging files...")
	
	// Create package archive
	archiveData, checksum, err := p.createArchive()
	if err != nil {
		return fmt.Errorf("failed to create package archive: %w", err)
	}
	
	fmt.Printf("   Archive size: %d bytes\n", len(archiveData))
	fmt.Printf("   Checksum: %s\n", checksum)

	// Prepare package metadata
	metadata := map[string]interface{}{
		"name":        p.cfg.Project.Name,
		"version":     p.cfg.Project.Version,
		"entry_point": p.cfg.Project.EntryPoint,
		"author":      creds.Username,
		"checksum":    checksum,
		"data":        base64.StdEncoding.EncodeToString(archiveData),
	}

	fmt.Println("🚀 Uploading to registry...")
	
	// Send to registry
	return p.uploadPackage(metadata, creds.Token)
}

// createArchive creates a tar.gz archive of the package files
func (p *Publisher) createArchive() ([]byte, string, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)
	
	// Get the directory containing the entry point
	baseDir := filepath.Dir(p.cfg.Project.EntryPoint)
	if baseDir == "." {
		baseDir = ""
	}
	
	// Collect files to include
	filesToInclude := []string{
		p.cfg.Project.EntryPoint,
		"polyloft.toml",
	}
	
	// Add all .pf files in the project directory
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and non-.pf files (except already included)
		if info.IsDir() {
			return nil
		}
		
		// Include .pf files
		if strings.HasSuffix(path, ".pf") {
			// Avoid duplicates
			isDuplicate := false
			for _, f := range filesToInclude {
				if f == path {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				filesToInclude = append(filesToInclude, path)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, "", fmt.Errorf("failed to walk directory: %w", err)
	}
	
	// Add files to archive
	for _, filePath := range filesToInclude {
		if err := p.addFileToArchive(tarWriter, filePath); err != nil {
			// If file doesn't exist, skip it (except for required files)
			if filePath == p.cfg.Project.EntryPoint || filePath == "polyloft.toml" {
				return nil, "", fmt.Errorf("required file not found: %s", filePath)
			}
		}
	}
	
	// Close writers
	if err := tarWriter.Close(); err != nil {
		return nil, "", err
	}
	if err := gzWriter.Close(); err != nil {
		return nil, "", err
	}
	
	// Calculate checksum
	archiveData := buf.Bytes()
	hash := sha256.Sum256(archiveData)
	checksum := hex.EncodeToString(hash[:])
	
	return archiveData, checksum, nil
}

// addFileToArchive adds a single file to the tar archive
func (p *Publisher) addFileToArchive(tw *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		return err
	}
	
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	
	// Use relative path in archive
	header.Name = filePath
	
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	
	_, err = io.Copy(tw, file)
	return err
}

// validateConfig validates the package configuration before publishing
func (p *Publisher) validateConfig() error {
	if p.cfg.Project.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if p.cfg.Project.Version == "" {
		return fmt.Errorf("project version is required")
	}
	if p.cfg.Project.EntryPoint == "" {
		return fmt.Errorf("project entry_point is required")
	}

	// Validate version format (basic semver check)
	// TODO: More robust semver validation
	if len(p.cfg.Project.Version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}

	// Check if entry point file exists
	if _, err := os.Stat(p.cfg.Project.EntryPoint); err != nil {
		return fmt.Errorf("entry point file not found: %s", p.cfg.Project.EntryPoint)
	}

	return nil
}

// uploadPackage uploads the package to the registry
func (p *Publisher) uploadPackage(metadata map[string]interface{}, token string) error {
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	url := fmt.Sprintf("%s/api/packages", p.registryURL)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("publish failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✓ Package published successfully")
	return nil
}

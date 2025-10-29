package searcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ArubikU/polyloft/internal/auth"
)

// Searcher handles package searching
type Searcher struct {
	registryURL string
}

// PackageResult represents a search result
type PackageResult struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// New creates a new searcher
func New() *Searcher {
	return &Searcher{
		registryURL: auth.GetRegistryURL(),
	}
}

// Search searches for packages matching the query
func (s *Searcher) Search(query string) ([]PackageResult, error) {
	// Build URL with query parameter
	searchURL := fmt.Sprintf("%s/api/search?q=%s", s.registryURL, url.QueryEscape(query))
	
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search packages: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response struct {
		Results []PackageResult `json:"results"`
		Count   int             `json:"count"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Results, nil
}

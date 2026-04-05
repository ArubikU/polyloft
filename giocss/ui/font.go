package ui

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func LooksLikeFontBinary(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	lowerHead := strings.ToLower(string(data[:imin(128, len(data))]))
	if strings.Contains(lowerHead, "<!doctype html") || strings.Contains(lowerHead, "<html") {
		return false
	}
	if bytes.HasPrefix(data, []byte("OTTO")) {
		return true
	}
	if bytes.HasPrefix(data, []byte("ttcf")) {
		return true
	}
	if bytes.HasPrefix(data, []byte("wOFF")) || bytes.HasPrefix(data, []byte("wOF2")) {
		return true
	}
	if data[0] == 0x00 && data[1] == 0x01 && data[2] == 0x00 && data[3] == 0x00 {
		return true
	}
	return false
}

func FontHasSpaceGlyph(data []byte) bool {
	parsedFont, err := sfnt.Parse(data)
	if err != nil {
		return false
	}
	var buf sfnt.Buffer
	idx, err := parsedFont.GlyphIndex(&buf, ' ')
	if err != nil {
		return false
	}
	adv, err := parsedFont.GlyphAdvance(&buf, idx, fixed.I(12), font.HintingNone)
	if err != nil || adv <= 0 {
		return false
	}
	bounds, _, err := parsedFont.GlyphBounds(&buf, idx, fixed.I(12), font.HintingNone)
	if err != nil {
		return false
	}
	if bounds.Min.X != 0 || bounds.Min.Y != 0 || bounds.Max.X != 0 || bounds.Max.Y != 0 {
		return false
	}
	return true
}

func LoadFontResource(pathOrURL string) ([]byte, error) {
	trimmed := strings.TrimSpace(pathOrURL)
	if trimmed == "" {
		return nil, fmt.Errorf("font path/url is empty")
	}
	if strings.HasPrefix(strings.ToLower(trimmed), "http://") || strings.HasPrefix(strings.ToLower(trimmed), "https://") {
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(trimmed)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("font request failed: %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if !LooksLikeFontBinary(body) {
			return nil, fmt.Errorf("remote resource is not a valid font binary")
		}
		if !FontHasSpaceGlyph(body) {
			return nil, fmt.Errorf("remote font has no valid space glyph")
		}
		return body, nil
	}
	bytes, err := os.ReadFile(trimmed)
	if err != nil {
		return nil, err
	}
	if !LooksLikeFontBinary(bytes) {
		return nil, fmt.Errorf("local resource is not a valid font binary")
	}
	if !FontHasSpaceGlyph(bytes) {
		return nil, fmt.Errorf("local font has no valid space glyph")
	}
	return bytes, nil
}

func imin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

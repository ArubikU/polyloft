package ui

import (
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type textMeasureKey struct {
	text          string
	fontSizeMilli int
	letterSpacing string
	lineHeight    string
	mode          byte
}

type textMeasureValue struct {
	w int
	h int
}

var (
	textMeasureCacheMu sync.RWMutex
	textMeasureCache   = map[textMeasureKey]textMeasureValue{}
)

func textMeasureKeyFor(text string, fontSize float32, css map[string]string, mode byte) (textMeasureKey, bool) {
	if len(text) > 256 {
		return textMeasureKey{}, false
	}
	return textMeasureKey{
		text:          text,
		fontSizeMilli: int(fontSize * 1000),
		letterSpacing: strings.TrimSpace(css["letter-spacing"]),
		lineHeight:    strings.TrimSpace(css["line-height"]),
		mode:          mode,
	}, true
}

func getCachedTextMeasure(key textMeasureKey) (textMeasureValue, bool) {
	textMeasureCacheMu.RLock()
	v, ok := textMeasureCache[key]
	textMeasureCacheMu.RUnlock()
	return v, ok
}

func putCachedTextMeasure(key textMeasureKey, v textMeasureValue) {
	textMeasureCacheMu.Lock()
	if len(textMeasureCache) > 2048 {
		textMeasureCache = map[textMeasureKey]textMeasureValue{}
	}
	textMeasureCache[key] = v
	textMeasureCacheMu.Unlock()
}

func CSSTextTransform(text string, css map[string]string) string {
	switch strings.ToLower(strings.TrimSpace(css["text-transform"])) {
	case "uppercase":
		return strings.ToUpper(text)
	case "lowercase":
		return strings.ToLower(text)
	case "capitalize":
		parts := strings.Fields(text)
		for i, p := range parts {
			if len(p) == 0 {
				continue
			}
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
		return strings.Join(parts, " ")
	default:
		return text
	}
}

func CSSApplyLetterSpacing(text string, css map[string]string) string {
	v := strings.TrimSpace(css["letter-spacing"])
	if v == "" || text == "" {
		return text
	}
	if strings.Contains(v, "0.") {
		return text
	}
	steps := CSSLengthValue(v, 0, 0, 0, 0)
	if steps <= 0 {
		return text
	}
	pad := strings.Repeat(" ", steps)
	runes := []rune(text)
	if len(runes) < 2 {
		return text
	}
	var b strings.Builder
	for i, r := range runes {
		b.WriteRune(r)
		if i < len(runes)-1 {
			b.WriteString(pad)
		}
	}
	return b.String()
}

func CSSLineHeightPx(css map[string]string, fontSize float32) int {
	v := strings.TrimSpace(css["line-height"])
	if v == "" {
		return max(1, int(float32(fontSize)*1.35))
	}
	if strings.HasSuffix(v, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
		if err == nil && pct > 0 {
			return max(1, int((pct/100.0)*float64(fontSize)))
		}
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
		if !strings.Contains(v, "px") && !strings.Contains(v, "vh") && !strings.Contains(v, "vw") && !strings.Contains(v, "%") {
			return max(1, int(f*float64(fontSize)))
		}
	}
	return max(1, CSSLengthValue(v, int(float32(fontSize)*1.35), 0, 0, 0))
}

func EstimateTextBox(text string, fontSize float32, css map[string]string) (int, int) {
	if key, ok := textMeasureKeyFor(text, fontSize, css, 0); ok {
		if cached, hit := getCachedTextMeasure(key); hit {
			return cached.w, cached.h
		}
		w, h := estimateTextBoxUncached(text, fontSize, css)
		putCachedTextMeasure(key, textMeasureValue{w: w, h: h})
		return w, h
	}
	return estimateTextBoxUncached(text, fontSize, css)
}

func estimateTextBoxUncached(text string, fontSize float32, css map[string]string) (int, int) {
	if text == "" {
		return max(8, int(fontSize)*2/3), CSSLineHeightPx(css, fontSize)
	}
	lines := strings.Split(text, "\n")
	maxChars := 0
	hasWideSymbol := false
	for _, line := range lines {
		count := len([]rune(line))
		if count > maxChars {
			maxChars = count
		}
		for _, r := range line {
			switch r {
			case '%', '&', '@', '#', '$', 'W', 'M':
				hasWideSymbol = true
			}
		}
	}
	charWidth := float64(fontSize) * 0.78
	width := int(float64(maxChars) * charWidth)
	width += CSSLengthValue(css["letter-spacing"], 0, 0, 0, 0) * max(0, maxChars-1)
	width += max(6, int(float64(fontSize)*0.5))
	if hasWideSymbol {
		width += max(2, int(float64(fontSize)*0.25))
	}
	height := len(lines) * CSSLineHeightPx(css, fontSize)
	height += max(1, int(float64(fontSize)*0.06))
	return max(8, width), max(1, height)
}

func EstimateTextLayoutBox(text string, fontSize float32, css map[string]string) (int, int) {
	if key, ok := textMeasureKeyFor(text, fontSize, css, 1); ok {
		if cached, hit := getCachedTextMeasure(key); hit {
			return cached.w, cached.h
		}
		w, h := estimateTextLayoutBoxUncached(text, fontSize, css)
		putCachedTextMeasure(key, textMeasureValue{w: w, h: h})
		return w, h
	}
	return estimateTextLayoutBoxUncached(text, fontSize, css)
}

func estimateTextLayoutBoxUncached(text string, fontSize float32, css map[string]string) (int, int) {
	if text == "" {
		return max(4, int(fontSize)*2/3), CSSLineHeightPx(css, fontSize)
	}
	lines := strings.Split(text, "\n")
	maxChars := 0
	for _, line := range lines {
		count := len([]rune(line))
		if count > maxChars {
			maxChars = count
		}
	}
	charWidth := float64(fontSize) * 0.72
	width := int(float64(maxChars) * charWidth)
	width += CSSLengthValue(css["letter-spacing"], 0, 0, 0, 0) * max(0, maxChars-1)
	width += max(4, int(float64(fontSize)*0.28))
	height := len(lines)*CSSLineHeightPx(css, fontSize) + 1
	return max(4, width), max(1, height)
}

func FitTextToWidth(text string, fontSize float32, css map[string]string, width int) string {
	if width <= 0 || text == "" {
		return text
	}
	availableChars := int(float64(width) / max(1.0, float64(fontSize)*0.62))
	if availableChars <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= availableChars {
		return text
	}
	if availableChars <= 1 {
		return string(runes[:1])
	}
	if strings.ToLower(strings.TrimSpace(css["text-overflow"])) != "ellipsis" {
		return string(runes[:availableChars])
	}
	return string(runes[:availableChars-1]) + "…"
}

func TextCharsForWidth(fontSize float32, width int) int {
	if width <= 0 {
		return 0
	}
	chars := int(float64(width) / max(1.0, float64(fontSize)*0.62))
	if chars < 1 {
		return 1
	}
	return chars
}

func WrapTextToWidth(text string, fontSize float32, width int) string {
	if text == "" || width <= 0 {
		return text
	}
	maxChars := TextCharsForWidth(fontSize, width)
	if maxChars <= 1 {
		return text
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	lines := make([]string, 0, 4)
	line := ""
	for _, word := range words {
		if line == "" {
			line = word
			continue
		}
		candidate := line + " " + word
		if len([]rune(candidate)) <= maxChars {
			line = candidate
			continue
		}
		lines = append(lines, line)
		line = word
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func SanitizeRenderText(text string) string {
	if text == "" {
		return text
	}
	var b strings.Builder
	lastWasSpace := false
	for _, r := range text {
		switch r {
		case '\uFFFD':
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		case '\t', '\r':
			if !lastWasSpace {
				b.WriteRune(' ')
				lastWasSpace = true
			}
		case '\n':
			b.WriteRune('\n')
			lastWasSpace = false
		default:
			if r < 32 {
				continue
			}
			if unicode.IsSpace(r) || unicode.Is(unicode.Cf, r) {
				if !lastWasSpace {
					b.WriteRune(' ')
					lastWasSpace = true
				}
				continue
			}
			if !unicode.IsPrint(r) {
				continue
			}
			b.WriteRune(r)
			lastWasSpace = false
		}
	}
	return b.String()
}

func CSSAllowsWrap(css map[string]string) bool {
	whiteSpace := strings.ToLower(strings.TrimSpace(css["white-space"]))
	if whiteSpace == "nowrap" {
		return false
	}
	overflow := strings.ToLower(strings.TrimSpace(css["overflow"]))
	if overflow == "hidden" && strings.ToLower(strings.TrimSpace(css["text-overflow"])) == "ellipsis" {
		return false
	}
	return true
}

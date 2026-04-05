package ui

import (
	"image/color"
	"strconv"
	"strings"
	"time"
)

func SplitCommaOutsideParens(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := make([]string, 0, 4)
	depth := 0
	start := 0
	for i, r := range input {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				chunk := strings.TrimSpace(input[start:i])
				if chunk != "" {
					parts = append(parts, chunk)
				}
				start = i + 1
			}
		}
	}
	last := strings.TrimSpace(input[start:])
	if last != "" {
		parts = append(parts, last)
	}
	return parts
}

func ParseTransition(css map[string]string) (time.Duration, []string) {
	raw := strings.TrimSpace(css["transition"])
	if raw == "" {
		return 0, nil
	}
	entries := SplitCommaOutsideParens(raw)
	if len(entries) == 0 {
		return 0, nil
	}
	props := make([]string, 0, len(entries))
	maxDuration := time.Duration(0)
	for _, entry := range entries {
		parts := strings.Fields(strings.ToLower(strings.TrimSpace(entry)))
		if len(parts) == 0 {
			continue
		}
		prop := strings.TrimRight(parts[0], ",")
		if prop == "" {
			continue
		}
		props = append(props, prop)
		for _, tok := range parts[1:] {
			if strings.HasSuffix(tok, "ms") {
				if n, err := strconv.Atoi(strings.TrimSuffix(tok, "ms")); err == nil {
					if d := time.Duration(n) * time.Millisecond; d > maxDuration {
						maxDuration = d
					}
					break
				}
			}
			if strings.HasSuffix(tok, "s") {
				if f, err := strconv.ParseFloat(strings.TrimSuffix(tok, "s"), 64); err == nil {
					if d := time.Duration(f * float64(time.Second)); d > maxDuration {
						maxDuration = d
					}
					break
				}
			}
		}
	}
	if maxDuration <= 0 {
		maxDuration = 180 * time.Millisecond
	}
	if len(props) == 0 {
		props = []string{"all"}
	}
	return maxDuration, props
}

func HasTransitionProp(props []string, name string) bool {
	for _, prop := range props {
		if prop == name {
			return true
		}
	}
	return false
}

func CSSScale(css map[string]string) float64 {
	if raw := strings.TrimSpace(css["scale"]); raw != "" {
		if f, err := strconv.ParseFloat(raw, 64); err == nil && f > 0 {
			return f
		}
	}
	if raw := strings.TrimSpace(css["transform"]); strings.Contains(raw, "scale(") {
		start := strings.Index(raw, "scale(")
		end := strings.Index(raw[start:], ")")
		if start >= 0 && end > 0 {
			inner := raw[start+len("scale(") : start+end]
			if f, err := strconv.ParseFloat(strings.TrimSpace(inner), 64); err == nil && f > 0 {
				return f
			}
		}
	}
	return 1.0
}

func ParseTransformTranslate(css map[string]string) (translateX, translateY float64) {
	raw := strings.TrimSpace(css["transform"])
	if raw == "" {
		return 0, 0
	}
	extractFunc := func(funcName string) float64 {
		start := strings.Index(raw, funcName+"(")
		if start < 0 {
			return 0
		}
		start += len(funcName)
		end := strings.Index(raw[start:], ")")
		if end < 0 {
			return 0
		}
		inner := strings.TrimSpace(raw[start+1 : start+end])
		inner = strings.TrimSpace(strings.Map(func(r rune) rune {
			if (r >= '0' && r <= '9') || r == '.' || r == '-' {
				return r
			}
			return -1
		}, strings.TrimSpace(inner)))
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v
		}
		return 0
	}

	if strings.Contains(raw, "translate(") {
		start := strings.Index(raw, "translate(")
		if start >= 0 {
			end := strings.Index(raw[start:], ")")
			if end > 0 {
				inner := strings.TrimSpace(raw[start+len("translate(") : start+end])
				parts := strings.Split(inner, ",")
				if len(parts) >= 1 {
					xStr := strings.TrimSpace(strings.Map(func(r rune) rune {
						if (r >= '0' && r <= '9') || r == '.' || r == '-' {
							return r
						}
						return -1
					}, strings.TrimSpace(parts[0])))
					if v, err := strconv.ParseFloat(xStr, 64); err == nil {
						translateX = v
					}
				}
				if len(parts) >= 2 {
					yStr := strings.TrimSpace(strings.Map(func(r rune) rune {
						if (r >= '0' && r <= '9') || r == '.' || r == '-' {
							return r
						}
						return -1
					}, strings.TrimSpace(parts[1])))
					if v, err := strconv.ParseFloat(yStr, 64); err == nil {
						translateY = v
					}
				}
			}
		}
	} else {
		translateX = extractFunc("translateX")
		translateY = extractFunc("translateY")
	}

	return translateX, translateY
}

func ParseFilterBrightness(css map[string]string) float64 {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return 1.0
	}
	start := strings.Index(raw, "brightness(")
	if start < 0 {
		return 1.0
	}
	start += len("brightness(")
	end := strings.Index(raw[start:], ")")
	if end < 0 {
		return 1.0
	}
	inner := strings.TrimSpace(raw[start : start+end])
	if strings.HasSuffix(inner, "%") {
		inner = strings.TrimSuffix(inner, "%")
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v / 100.0
		}
	}
	if v, err := strconv.ParseFloat(inner, 64); err == nil {
		return v
	}
	return 1.0
}

func ParseFilterContrast(css map[string]string) float64 {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return 1.0
	}
	start := strings.Index(raw, "contrast(")
	if start < 0 {
		return 1.0
	}
	start += len("contrast(")
	end := strings.Index(raw[start:], ")")
	if end < 0 {
		return 1.0
	}
	inner := strings.TrimSpace(raw[start : start+end])
	if strings.HasSuffix(inner, "%") {
		inner = strings.TrimSuffix(inner, "%")
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v / 100.0
		}
	}
	if v, err := strconv.ParseFloat(inner, 64); err == nil {
		return v
	}
	return 1.0
}

func ParseFilterSaturate(css map[string]string) float64 {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return 1.0
	}
	start := strings.Index(raw, "saturate(")
	if start < 0 {
		return 1.0
	}
	start += len("saturate(")
	end := strings.Index(raw[start:], ")")
	if end < 0 {
		return 1.0
	}
	inner := strings.TrimSpace(raw[start : start+end])
	if strings.HasSuffix(inner, "%") {
		inner = strings.TrimSuffix(inner, "%")
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v / 100.0
		}
	}
	if v, err := strconv.ParseFloat(inner, 64); err == nil {
		return v
	}
	return 1.0
}

func ParseFilterGrayscale(css map[string]string) float64 {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return 0.0
	}
	start := strings.Index(raw, "grayscale(")
	if start < 0 {
		return 0.0
	}
	start += len("grayscale(")
	end := strings.Index(raw[start:], ")")
	if end < 0 {
		return 0.0
	}
	inner := strings.TrimSpace(raw[start : start+end])
	if strings.HasSuffix(inner, "%") {
		inner = strings.TrimSuffix(inner, "%")
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v / 100.0
		}
	}
	if v, err := strconv.ParseFloat(inner, 64); err == nil {
		return v
	}
	return 0.0
}

func ParseFilterInvert(css map[string]string) float64 {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return 0.0
	}
	start := strings.Index(raw, "invert(")
	if start < 0 {
		return 0.0
	}
	start += len("invert(")
	end := strings.Index(raw[start:], ")")
	if end < 0 {
		return 0.0
	}
	inner := strings.TrimSpace(raw[start : start+end])
	if strings.HasSuffix(inner, "%") {
		inner = strings.TrimSuffix(inner, "%")
		if v, err := strconv.ParseFloat(inner, 64); err == nil {
			return v / 100.0
		}
	}
	if v, err := strconv.ParseFloat(inner, 64); err == nil {
		return v
	}
	return 0.0
}

type PositionInfo struct {
	Position  string
	OffsetX   float64
	OffsetY   float64
	IsLeft    bool
	IsTop     bool
	HasLeft   bool
	HasRight  bool
	HasTop    bool
	HasBottom bool
}

func ParsePositionAndOffset(css map[string]string, elemW, elemH int) PositionInfo {
	info := PositionInfo{Position: "static"}
	posRaw := css["position"]
	if posRaw == "" && css["top"] == "" && css["bottom"] == "" && css["left"] == "" && css["right"] == "" {
		return info
	}
	pos := strings.ToLower(strings.TrimSpace(posRaw))
	if pos == "relative" || pos == "absolute" || pos == "fixed" || pos == "sticky" {
		info.Position = pos
	}
	if info.Position != "static" {
		topVal := strings.TrimSpace(css["top"])
		bottomVal := strings.TrimSpace(css["bottom"])
		if topVal != "" && topVal != "auto" {
			info.OffsetY = float64(CSSLengthValue(topVal, 0, elemH, elemW, elemH))
			info.IsTop = true
			info.HasTop = true
		} else if bottomVal != "" && bottomVal != "auto" {
			info.OffsetY = float64(CSSLengthValue(bottomVal, 0, elemH, elemW, elemH))
			info.IsTop = false
			info.HasBottom = true
		}
		leftVal := strings.TrimSpace(css["left"])
		rightVal := strings.TrimSpace(css["right"])
		if leftVal != "" && leftVal != "auto" {
			info.OffsetX = float64(CSSLengthValue(leftVal, 0, elemW, elemW, elemH))
			info.IsLeft = true
			info.HasLeft = true
		} else if rightVal != "" && rightVal != "auto" {
			info.OffsetX = float64(CSSLengthValue(rightVal, 0, elemW, elemW, elemH))
			info.IsLeft = false
			info.HasRight = true
		}
	}
	return info
}

func ParseZIndex(css map[string]string) int {
	raw := css["z-index"]
	if raw == "" {
		return 0
	}
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "auto") {
		return 0
	}
	if v, err := strconv.Atoi(raw); err == nil {
		return v
	}
	return 0
}

func ParseWidthHeightCSS(css map[string]string, elemW, elemH int) (int, int, bool, bool) {
	width := elemW
	height := elemH
	hasWidth := false
	hasHeight := false
	if css["width"] == "" && css["height"] == "" && css["min-width"] == "" && css["max-width"] == "" && css["min-height"] == "" && css["max-height"] == "" {
		return width, height, hasWidth, hasHeight
	}
	if raw := strings.TrimSpace(css["width"]); raw != "" && !strings.EqualFold(raw, "auto") {
		width = CSSLengthValue(raw, elemW, elemW, elemW, elemH)
		hasWidth = true
	}
	if raw := strings.TrimSpace(css["height"]); raw != "" && !strings.EqualFold(raw, "auto") {
		height = CSSLengthValue(raw, elemH, elemH, elemW, elemH)
		hasHeight = true
	}
	if raw := strings.TrimSpace(css["min-width"]); raw != "" {
		minW := CSSLengthValue(raw, 0, elemW, elemW, elemH)
		if width < minW {
			width = minW
		}
	}
	if raw := strings.TrimSpace(css["max-width"]); raw != "" {
		maxW := CSSLengthValue(raw, elemW, elemW, elemW, elemH)
		if width > maxW {
			width = maxW
		}
	}
	if raw := strings.TrimSpace(css["min-height"]); raw != "" {
		minH := CSSLengthValue(raw, 0, elemH, elemW, elemH)
		if height < minH {
			height = minH
		}
	}
	if raw := strings.TrimSpace(css["max-height"]); raw != "" {
		maxH := CSSLengthValue(raw, elemH, elemH, elemW, elemH)
		if height > maxH {
			height = maxH
		}
	}
	return width, height, hasWidth, hasHeight
}

type ShadowParams struct {
	OffsetX float64
	OffsetY float64
	Blur    float64
	Color   string
}

func ParseFilterDropShadow(css map[string]string) ShadowParams {
	raw := strings.TrimSpace(css["filter"])
	if raw == "" {
		return ShadowParams{}
	}
	start := strings.Index(raw, "drop-shadow(")
	if start < 0 {
		return ShadowParams{}
	}
	start += len("drop-shadow(")
	depth := 1
	end := start
	for end < len(raw) && depth > 0 {
		if raw[end] == '(' {
			depth++
		} else if raw[end] == ')' {
			depth--
		}
		end++
	}
	if depth != 0 {
		return ShadowParams{}
	}
	inner := strings.TrimSpace(raw[start : end-1])
	parts := strings.Fields(inner)
	if len(parts) < 3 {
		return ShadowParams{}
	}
	shadow := ShadowParams{}
	if v, err := strconv.ParseFloat(strings.TrimRight(strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			return r
		}
		return -1
	}, parts[0]), "-"), 64); err == nil {
		shadow.OffsetX = v
	}
	if v, err := strconv.ParseFloat(strings.TrimRight(strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			return r
		}
		return -1
	}, parts[1]), "-"), 64); err == nil {
		shadow.OffsetY = v
	}
	if v, err := strconv.ParseFloat(strings.TrimRight(strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			return r
		}
		return -1
	}, parts[2]), "-"), 64); err == nil {
		shadow.Blur = v
	}
	if len(parts) >= 4 {
		shadow.Color = strings.Join(parts[3:], " ")
	}
	return shadow
}

type TextDecorationInfo struct {
	Line      string
	Style     string
	Color     string
	Thickness string
}

func ParseTextDecoration(css map[string]string) TextDecorationInfo {
	info := TextDecorationInfo{Line: "none"}
	if raw := strings.TrimSpace(css["text-decoration"]); raw != "" {
		parts := strings.Fields(raw)
		if len(parts) > 0 {
			info.Line = strings.ToLower(parts[0])
		}
		if len(parts) > 1 {
			info.Style = strings.ToLower(parts[1])
		}
		if len(parts) > 2 {
			info.Color = parts[2]
		}
	}
	if raw := strings.TrimSpace(css["text-decoration-line"]); raw != "" {
		info.Line = strings.ToLower(raw)
	}
	if raw := strings.TrimSpace(css["text-decoration-style"]); raw != "" {
		info.Style = strings.ToLower(raw)
	}
	if raw := strings.TrimSpace(css["text-decoration-color"]); raw != "" {
		info.Color = raw
	}
	if raw := strings.TrimSpace(css["text-decoration-thickness"]); raw != "" {
		info.Thickness = raw
	}
	return info
}

func ParseCursor(css map[string]string) string {
	if raw := strings.TrimSpace(css["cursor"]); raw != "" {
		return strings.ToLower(raw)
	}
	return "default"
}

func ParseAspectRatio(css map[string]string) (float64, bool) {
	raw := strings.TrimSpace(css["aspect-ratio"])
	if raw == "" || strings.EqualFold(raw, "auto") {
		return 0, false
	}
	parts := strings.Split(raw, "/")
	if len(parts) >= 2 {
		widthStr := strings.TrimSpace(parts[0])
		heightStr := strings.TrimSpace(parts[1])
		if w, err1 := strconv.ParseFloat(widthStr, 64); err1 == nil {
			if h, err2 := strconv.ParseFloat(heightStr, 64); err2 == nil && h != 0 {
				return w / h, true
			}
		}
	} else if len(parts) == 1 {
		if ratio, err := strconv.ParseFloat(parts[0], 64); err == nil && ratio > 0 {
			return ratio, true
		}
	}
	return 0, false
}

func MixColor(from color.Color, to color.Color, t float64) color.NRGBA {
	fr, fg, fb, fa := from.RGBA()
	tr, tg, tb, ta := to.RGBA()
	lerp := func(a uint32, b uint32) uint8 {
		v := float64(a>>8) + (float64(b>>8)-float64(a>>8))*t
		if v < 0 {
			v = 0
		}
		if v > 255 {
			v = 255
		}
		return uint8(v)
	}
	return color.NRGBA{R: lerp(fr, tr), G: lerp(fg, tg), B: lerp(fb, tb), A: lerp(fa, ta)}
}

func ClampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v + 0.5)
}

func ApplyBrightnessToColor(col color.NRGBA, brightness float64) color.NRGBA {
	if brightness < 0 {
		brightness = 0
	}
	return color.NRGBA{
		R: ClampUint8(float64(col.R) * brightness),
		G: ClampUint8(float64(col.G) * brightness),
		B: ClampUint8(float64(col.B) * brightness),
		A: col.A,
	}
}

func ApplyContrastToColor(col color.NRGBA, contrast float64) color.NRGBA {
	if contrast < 0 {
		contrast = 0
	}
	mid := 128.0
	return color.NRGBA{
		R: ClampUint8(mid + (float64(col.R)-mid)*contrast),
		G: ClampUint8(mid + (float64(col.G)-mid)*contrast),
		B: ClampUint8(mid + (float64(col.B)-mid)*contrast),
		A: col.A,
	}
}

func ApplySaturationToColor(col color.NRGBA, saturation float64) color.NRGBA {
	if saturation < 0 {
		saturation = 0
	}
	r := float64(col.R)
	g := float64(col.G)
	b := float64(col.B)
	r2 := (0.213+0.787*saturation)*r + (0.715-0.715*saturation)*g + (0.072-0.072*saturation)*b
	g2 := (0.213-0.213*saturation)*r + (0.715+0.285*saturation)*g + (0.072-0.072*saturation)*b
	b2 := (0.213-0.213*saturation)*r + (0.715-0.715*saturation)*g + (0.072+0.928*saturation)*b
	return color.NRGBA{
		R: ClampUint8(r2),
		G: ClampUint8(g2),
		B: ClampUint8(b2),
		A: col.A,
	}
}

func ApplyGrayscaleToColor(col color.NRGBA, grayscale float64) color.NRGBA {
	if grayscale < 0 {
		grayscale = 0
	}
	if grayscale > 1.0 {
		grayscale = 1.0
	}
	r := float64(col.R)
	g := float64(col.G)
	b := float64(col.B)
	s := 1.0 - grayscale
	r2 := (0.2126+0.7874*s)*r + (0.7152-0.7152*s)*g + (0.0722-0.0722*s)*b
	g2 := (0.2126-0.2126*s)*r + (0.7152+0.2848*s)*g + (0.0722-0.0722*s)*b
	b2 := (0.2126-0.2126*s)*r + (0.7152-0.7152*s)*g + (0.0722+0.9278*s)*b
	return color.NRGBA{
		R: ClampUint8(r2),
		G: ClampUint8(g2),
		B: ClampUint8(b2),
		A: col.A,
	}
}

func ApplyInvertToColor(col color.NRGBA, invert float64) color.NRGBA {
	if invert < 0 {
		invert = 0
	}
	if invert > 1.0 {
		invert = 1.0
	}
	inverted := color.NRGBA{R: 255 - col.R, G: 255 - col.G, B: 255 - col.B, A: col.A}
	return color.NRGBA{
		R: ClampUint8(float64(col.R)*(1.0-invert) + float64(inverted.R)*invert),
		G: ClampUint8(float64(col.G)*(1.0-invert) + float64(inverted.G)*invert),
		B: ClampUint8(float64(col.B)*(1.0-invert) + float64(inverted.B)*invert),
		A: col.A,
	}
}

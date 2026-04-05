package ui

import (
	"image/color"
	"math"
	"strconv"
	"strings"
)

func ParseHexColor(input string, fallback color.Color) color.Color {
	text := strings.TrimSpace(strings.TrimPrefix(input, "#"))
	if len(text) == 0 {
		return fallback
	}
	lowerInput := strings.ToLower(strings.TrimSpace(input))
	if strings.HasPrefix(lowerInput, "linear-gradient(") || strings.HasPrefix(lowerInput, "radial-gradient(") {
		start := strings.Index(lowerInput, "(")
		end := strings.LastIndex(lowerInput, ")")
		if start >= 0 && end > start {
			inner := input[start+1 : end]
			parts := SplitColorArgs(inner)
			for _, p := range parts {
				cand := strings.TrimSpace(p)
				if cand == "" {
					continue
				}
				if strings.Contains(cand, "deg") || strings.HasPrefix(strings.ToLower(cand), "to ") || strings.HasPrefix(strings.ToLower(cand), "circle") || strings.HasPrefix(strings.ToLower(cand), "ellipse") {
					continue
				}
				if c := ParseHexColor(cand, nil); c != nil {
					return c
				}
				fields := strings.Fields(cand)
				if len(fields) > 0 {
					if c := ParseHexColor(fields[0], nil); c != nil {
						return c
					}
				}
			}
		}
		return fallback
	}
	if strings.HasPrefix(lowerInput, "rgb(") || strings.HasPrefix(lowerInput, "rgba(") {
		return ParseRGBColor(input, fallback)
	}
	if strings.HasPrefix(lowerInput, "hsl(") || strings.HasPrefix(lowerInput, "hsla(") {
		return ParseHSLColor(input, fallback)
	}
	if strings.HasPrefix(lowerInput, "cmyk(") {
		return ParseCMYKColor(input, fallback)
	}
	if named := ParseNamedColor(input); named != nil {
		return named
	}
	if len(text) == 4 {
		text = string([]byte{text[0], text[0], text[1], text[1], text[2], text[2], text[3], text[3]})
	}
	if len(text) == 3 {
		text = string([]byte{text[0], text[0], text[1], text[1], text[2], text[2]})
	}
	if len(text) == 8 {
		r, errR := strconv.ParseUint(text[0:2], 16, 8)
		g, errG := strconv.ParseUint(text[2:4], 16, 8)
		b, errB := strconv.ParseUint(text[4:6], 16, 8)
		a, errA := strconv.ParseUint(text[6:8], 16, 8)
		if errR != nil || errG != nil || errB != nil || errA != nil {
			return fallback
		}
		return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}
	if len(text) != 6 {
		return fallback
	}
	r, errR := strconv.ParseUint(text[0:2], 16, 8)
	g, errG := strconv.ParseUint(text[2:4], 16, 8)
	b, errB := strconv.ParseUint(text[4:6], 16, 8)
	if errR != nil || errG != nil || errB != nil {
		return fallback
	}
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xFF}
}

func SplitColorArgs(input string) []string {
	text := strings.TrimSpace(input)
	if text == "" {
		return nil
	}
	if strings.Contains(text, ",") {
		parts := make([]string, 0, 6)
		depth := 0
		start := 0
		for i, r := range text {
			switch r {
			case '(':
				depth++
			case ')':
				if depth > 0 {
					depth--
				}
			case ',':
				if depth == 0 {
					chunk := strings.TrimSpace(text[start:i])
					if chunk != "" {
						parts = append(parts, chunk)
					}
					start = i + 1
				}
			}
		}
		last := strings.TrimSpace(text[start:])
		if last != "" {
			parts = append(parts, last)
		}
		if len(parts) > 0 {
			return parts
		}
	}
	text = strings.ReplaceAll(text, "/", " ")
	return strings.Fields(text)
}

func ParseRGBColor(input string, fallback color.Color) color.Color {
	text := strings.ToLower(strings.TrimSpace(input))
	start := strings.Index(text, "(")
	end := strings.LastIndex(text, ")")
	if start < 0 || end <= start {
		return fallback
	}
	parts := SplitColorArgs(text[start+1 : end])
	if len(parts) < 3 {
		return fallback
	}
	parseChannel := func(raw string) (uint8, bool) {
		v := strings.TrimSpace(raw)
		if strings.HasSuffix(v, "%") {
			f, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
			if err != nil {
				return 0, false
			}
			if f < 0 {
				f = 0
			}
			if f > 100 {
				f = 100
			}
			return uint8((f / 100.0) * 255.0), true
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		if n < 0 {
			n = 0
		}
		if n > 255 {
			n = 255
		}
		return uint8(n), true
	}
	r, okR := parseChannel(parts[0])
	g, okG := parseChannel(parts[1])
	b, okB := parseChannel(parts[2])
	if !okR || !okG || !okB {
		return fallback
	}
	a := uint8(255)
	if len(parts) >= 4 {
		af, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil {
			return fallback
		}
		if af < 0 {
			af = 0
		}
		if af > 1 {
			af = 1
		}
		a = uint8(af * 255.0)
	}
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func ParseHSLColor(input string, fallback color.Color) color.Color {
	text := strings.ToLower(strings.TrimSpace(input))
	start := strings.Index(text, "(")
	end := strings.LastIndex(text, ")")
	if start < 0 || end <= start {
		return fallback
	}
	parts := SplitColorArgs(text[start+1 : end])
	if len(parts) < 3 {
		return fallback
	}
	parseHue := func(raw string) (float64, bool) {
		v := strings.TrimSpace(strings.TrimSuffix(raw, "deg"))
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		f = math.Mod(f, 360.0)
		if f < 0 {
			f += 360.0
		}
		return f, true
	}
	parsePct := func(raw string) (float64, bool) {
		v := strings.TrimSpace(raw)
		if !strings.HasSuffix(v, "%") {
			return 0, false
		}
		f, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
		if err != nil {
			return 0, false
		}
		if f < 0 {
			f = 0
		}
		if f > 100 {
			f = 100
		}
		return f / 100.0, true
	}
	h, okH := parseHue(parts[0])
	s, okS := parsePct(parts[1])
	l, okL := parsePct(parts[2])
	if !okH || !okS || !okL {
		return fallback
	}
	a := 1.0
	if len(parts) >= 4 {
		v := strings.TrimSpace(parts[3])
		if strings.HasSuffix(v, "%") {
			f, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
			if err != nil {
				return fallback
			}
			a = f / 100.0
		} else {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fallback
			}
			a = f
		}
		if a < 0 {
			a = 0
		}
		if a > 1 {
			a = 1
		}
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2
	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	to8 := func(v float64) uint8 {
		p := (v + m) * 255.0
		if p < 0 {
			p = 0
		}
		if p > 255 {
			p = 255
		}
		return uint8(p + 0.5)
	}
	return color.NRGBA{R: to8(r1), G: to8(g1), B: to8(b1), A: uint8(a*255 + 0.5)}
}

func ParseCMYKColor(input string, fallback color.Color) color.Color {
	text := strings.ToLower(strings.TrimSpace(input))
	start := strings.Index(text, "(")
	end := strings.LastIndex(text, ")")
	if start < 0 || end <= start {
		return fallback
	}
	parts := SplitColorArgs(text[start+1 : end])
	if len(parts) < 4 {
		return fallback
	}
	parseC := func(raw string) (float64, bool) {
		v := strings.TrimSpace(raw)
		if strings.HasSuffix(v, "%") {
			f, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64)
			if err != nil {
				return 0, false
			}
			if f < 0 {
				f = 0
			}
			if f > 100 {
				f = 100
			}
			return f / 100.0, true
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		if f < 0 {
			f = 0
		}
		if f > 1 {
			f = 1
		}
		return f, true
	}
	c, okC := parseC(parts[0])
	m, okM := parseC(parts[1])
	y, okY := parseC(parts[2])
	k, okK := parseC(parts[3])
	if !okC || !okM || !okY || !okK {
		return fallback
	}
	a := 1.0
	if len(parts) >= 5 {
		f, err := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
		if err == nil {
			a = f
			if a < 0 {
				a = 0
			}
			if a > 1 {
				a = 1
			}
		}
	}
	r := 255.0 * (1 - c) * (1 - k)
	g := 255.0 * (1 - m) * (1 - k)
	b := 255.0 * (1 - y) * (1 - k)
	clamp := func(v float64) uint8 {
		if v < 0 {
			v = 0
		}
		if v > 255 {
			v = 255
		}
		return uint8(v + 0.5)
	}
	return color.NRGBA{R: clamp(r), G: clamp(g), B: clamp(b), A: uint8(a*255 + 0.5)}
}

func ParseNamedColor(input string) color.Color {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "white":
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case "black":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case "red":
		return color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	case "green":
		return color.NRGBA{R: 0, G: 128, B: 0, A: 255}
	case "blue":
		return color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	case "yellow":
		return color.NRGBA{R: 255, G: 255, B: 0, A: 255}
	case "gray", "grey":
		return color.NRGBA{R: 128, G: 128, B: 128, A: 255}
	case "silver":
		return color.NRGBA{R: 192, G: 192, B: 192, A: 255}
	case "maroon":
		return color.NRGBA{R: 128, G: 0, B: 0, A: 255}
	case "navy":
		return color.NRGBA{R: 0, G: 0, B: 128, A: 255}
	case "teal":
		return color.NRGBA{R: 0, G: 128, B: 128, A: 255}
	case "olive":
		return color.NRGBA{R: 128, G: 128, B: 0, A: 255}
	case "lime":
		return color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	case "aqua", "cyan":
		return color.NRGBA{R: 0, G: 255, B: 255, A: 255}
	case "magenta", "fuchsia":
		return color.NRGBA{R: 255, G: 0, B: 255, A: 255}
	case "orange":
		return color.NRGBA{R: 255, G: 165, B: 0, A: 255}
	case "purple":
		return color.NRGBA{R: 128, G: 0, B: 128, A: 255}
	case "transparent":
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	return nil
}

func ApplyCSSOpacity(c color.Color, css map[string]string) color.Color {
	v := strings.TrimSpace(css["opacity"])
	if v == "" {
		return c
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return c
	}
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	r, g, b, a := c.RGBA()
	n := color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
	n.A = uint8(float64(n.A) * f)
	return n
}

package render

import (
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

func ParseBorderRadii(css map[string]string, w int, h int) CornerRadii {
	basis := max(w, h)
	parse := func(raw string, fallback int) int {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return fallback
		}
		if idx := strings.Index(raw, "/"); idx >= 0 {
			raw = strings.TrimSpace(raw[:idx])
		}
		parts := strings.Fields(raw)
		if len(parts) == 0 {
			return fallback
		}
		return max(0, uicore.CSSLengthValue(parts[0], fallback, basis, w, h))
	}
	r := CornerRadii{}
	if raw := strings.TrimSpace(css["border-radius"]); raw != "" {
		main := raw
		if idx := strings.Index(main, "/"); idx >= 0 {
			main = strings.TrimSpace(main[:idx])
		}
		parts := strings.Fields(main)
		if len(parts) == 1 {
			v := parse(parts[0], 0)
			r = CornerRadii{NW: v, NE: v, SE: v, SW: v}
		} else if len(parts) == 2 {
			v1 := parse(parts[0], 0)
			v2 := parse(parts[1], 0)
			r = CornerRadii{NW: v1, NE: v2, SE: v1, SW: v2}
		} else if len(parts) == 3 {
			v1 := parse(parts[0], 0)
			v2 := parse(parts[1], 0)
			v3 := parse(parts[2], 0)
			r = CornerRadii{NW: v1, NE: v2, SE: v3, SW: v2}
		} else if len(parts) >= 4 {
			r = CornerRadii{
				NW: parse(parts[0], 0),
				NE: parse(parts[1], 0),
				SE: parse(parts[2], 0),
				SW: parse(parts[3], 0),
			}
		}
	}
	r.NW = parse(css["border-top-left-radius"], r.NW)
	r.NE = parse(css["border-top-right-radius"], r.NE)
	r.SE = parse(css["border-bottom-right-radius"], r.SE)
	r.SW = parse(css["border-bottom-left-radius"], r.SW)
	maxR := min(w, h) / 2
	r.NW = min(r.NW, maxR)
	r.NE = min(r.NE, maxR)
	r.SE = min(r.SE, maxR)
	r.SW = min(r.SW, maxR)
	return r
}

func BorderRadiusValue(css map[string]string, w int, h int) int {
	r := ParseBorderRadii(css, w, h)
	if r.NW != r.NE || r.NE != r.SE || r.SE != r.SW {
		return r.Max()
	}
	return r.NW
}

func ResolveRoundedFromProps(props map[string]any, css map[string]string, w int, h int) int {
	radius := BorderRadiusValue(css, w, h)
	if radius > 0 {
		return radius
	}
	if props == nil {
		return 0
	}
	for _, key := range []string{"rounded", "radius"} {
		v, ok := props[key]
		if !ok {
			continue
		}
		s := anyToString(v, "")
		if s != "" {
			return max(0, uicore.CSSLengthValue(s, 0, max(w, h), w, h))
		}
		switch n := v.(type) {
		case int:
			return max(0, n)
		case int64:
			return max(0, int(n))
		case float64:
			return max(0, int(n))
		}
	}
	return 0
}

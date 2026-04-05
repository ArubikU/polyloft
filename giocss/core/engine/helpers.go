package engine

import (
	"image"
	"math"
	"sort"
	"strconv"
	"strings"
)

func RefreshBoolMap(dst map[string]bool, src map[string]bool) map[string]bool {
	if len(src) == 0 {
		if dst != nil {
			clear(dst)
		}
		return dst
	}
	if dst == nil {
		dst = make(map[string]bool, len(src))
	} else {
		clear(dst)
	}
	for k, v := range src {
		if strings.TrimSpace(k) == "" {
			continue
		}
		dst[k] = v
	}
	return dst
}

func ChildPath(parent string, idx int) string {
	return parent + "/" + strconv.Itoa(idx)
}

func LowerASCIIIfNeeded(s string) string {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return strings.ToLower(s)
		}
	}
	return s
}

func InheritedTextCSSSignature(parent map[string]string) string {
	if len(parent) == 0 {
		return ""
	}
	keys := []string{"color", "font-family", "font-size", "font-style", "font-weight", "line-height", "letter-spacing", "text-align", "text-transform", "white-space"}
	var b strings.Builder
	b.Grow(128)
	for _, k := range keys {
		v := strings.TrimSpace(parent[k])
		if v == "" {
			continue
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
		b.WriteByte(';')
	}
	return b.String()
}

func CanonicalCSSSignature(css map[string]string) string {
	if len(css) == 0 {
		return ""
	}
	keys := make([]string, 0, len(css))
	for k, v := range css {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(strings.TrimSpace(css[k]))
		b.WriteByte(';')
	}
	return b.String()
}

func ExpandedViewportRect(viewW, viewH, margin int) image.Rectangle {
	return image.Rect(-margin, -margin, viewW+margin, viewH+margin)
}

func NodeRenderRect(x, y, w, h int, tx, ty float64, scale float64) image.Rectangle {
	r := image.Rect(x+int(tx), y+int(ty), x+w+int(tx), y+h+int(ty))
	if scale <= 0 || math.Abs(scale-1.0) < 0.001 {
		return r
	}
	cx := float64(r.Min.X+r.Max.X) / 2.0
	cy := float64(r.Min.Y+r.Max.Y) / 2.0
	hw := float64(r.Dx()) * scale / 2.0
	hh := float64(r.Dy()) * scale / 2.0
	return image.Rect(int(cx-hw), int(cy-hh), int(cx+hw), int(cy+hh))
}

func PointInRect(p image.Point, r image.Rectangle) bool {
	return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y
}

func NodeHasExplicitZIndex(props map[string]any) bool {
	if props == nil {
		return false
	}
	for _, key := range []string{"z-index", "zIndex"} {
		if raw, ok := props[key]; ok {
			v := strings.TrimSpace(anyToString(raw, ""))
			if v != "" && !strings.EqualFold(v, "auto") && v != "0" {
				return true
			}
		}
	}
	styleRaw := strings.ToLower(strings.TrimSpace(anyToString(props["style"], "")))
	return strings.Contains(styleRaw, "z-index")
}

func ParseZIndexFromPropsFast(props map[string]any) int {
	if len(props) == 0 {
		return 0
	}
	for _, key := range []string{"z-index", "zIndex"} {
		if raw, ok := props[key]; ok {
			v := strings.TrimSpace(anyToString(raw, ""))
			if v == "" || strings.EqualFold(v, "auto") {
				continue
			}
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	style := strings.ToLower(anyToString(props["style"], ""))
	if idx := strings.Index(style, "z-index:"); idx >= 0 {
		rest := style[idx+len("z-index:"):]
		if semi := strings.Index(rest, ";"); semi >= 0 {
			rest = rest[:semi]
		}
		rest = strings.TrimSpace(rest)
		if n, err := strconv.Atoi(rest); err == nil {
			return n
		}
	}
	return 0
}

func PropsMayBeOutOfFlow(props map[string]any) bool {
	if len(props) == 0 {
		return false
	}
	pos := strings.ToLower(strings.TrimSpace(anyToString(props["position"], "")))
	if pos == "absolute" || pos == "fixed" || pos == "sticky" {
		return true
	}
	for _, k := range []string{"top", "left", "right", "bottom"} {
		if _, ok := props[k]; ok {
			return true
		}
	}
	style := strings.ToLower(anyToString(props["style"], ""))
	return strings.Contains(style, "position:") || strings.Contains(style, "transform:")
}

// anyToString and cloneStringMap are unexported helpers used throughout core.
// They were previously defined in style_stylesheet.go (now ui/stylesheet.go).

func anyToString(candidate any, fallback string) string {
	if typed, ok := candidate.(string); ok {
		normalized := strings.ToValidUTF8(typed, "")
		trimmed := strings.TrimSpace(normalized)
		if trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func CloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneStringMap(in map[string]string) map[string]string {
	return CloneStringMap(in)
}

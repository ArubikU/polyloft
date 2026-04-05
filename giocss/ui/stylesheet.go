package ui

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// StyleSheet stores CSS-like class rules.
// Supports both .css file loading and programmatic rule setting via SetRule().
//
// Two ways to use from Polyloft script:
//
//	// Method 1 – CSS text / file
//	var ss = UI.StyleSheet()
//	ss.loadFile("styles.css")        // .title { color: #fff; font-size: 20px }
//	app.attachStylesheet(ss)
//
//	// Method 2 – programmatic
//	var ss = UI.StyleSheet()
//	ss.set(".title", "color", "#e8e8e8")
//	ss.set(".title", "text-align", "center")
//	app.attachStylesheet(ss)
//
//	// Applying to a node
//	title.setClass("title")           // or: title.set("class", "title header")
//	title.style("font-size", "24px")  // inline override (highest priority)
type StyleSheet struct {
	mu      sync.RWMutex
	classes map[string]map[string]string // normalizedClassName → property → value
	rev     int64
}

type resolveStyleCacheKey struct {
	rev       int64
	viewportW int
	className string
	tag       string
	inlineSig string
}

var (
	resolveStyleCacheMu sync.RWMutex
	resolveStyleCache   = map[resolveStyleCacheKey]map[string]string{}
)

func NewStyleSheet() *StyleSheet {
	return &StyleSheet{classes: make(map[string]map[string]string)}
}

// normalizeCSSClass strips a leading dot and trims whitespace.
func normalizeCSSClass(name string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(name), "."))
}

// SetRule stores a CSS declaration for className.
// className may optionally start with ".".
func (ss *StyleSheet) SetRule(className, property, val string) {
	cls := normalizeCSSClass(className)
	if cls == "" {
		return
	}
	prop := strings.ToLower(strings.TrimSpace(property))
	v := strings.TrimSpace(val)
	if prop == "" {
		return
	}
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if ss.classes[cls] == nil {
		ss.classes[cls] = make(map[string]string, 4)
	}
	ss.classes[cls][prop] = v
	ss.rev++
}

// GetRule retrieves a single CSS property for a class.
func (ss *StyleSheet) GetRule(className, property string) (string, bool) {
	cls := normalizeCSSClass(className)
	prop := strings.ToLower(strings.TrimSpace(property))
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if m, ok := ss.classes[cls]; ok {
		if v, ok2 := m[prop]; ok2 {
			return v, true
		}
	}
	return "", false
}

// GetClassProps returns a snapshot copy of all declarations for className.
func (ss *StyleSheet) GetClassProps(className string) map[string]string {
	cls := normalizeCSSClass(className)
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if m, ok := ss.classes[cls]; ok {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	return nil
}

// ParseCSSText parses a CSS string into the StyleSheet and returns the number
// of declarations successfully loaded.
//
// Supported syntax:
//
//	.classname { property: value; property: value }
//	.a, .b   { property: value }
//	/* block comments */
//	// line comments
func (ss *StyleSheet) ParseCSSText(text string) int {
	// Strip /* ... */ block comments.
	for {
		si := strings.Index(text, "/*")
		if si < 0 {
			break
		}
		ei := strings.Index(text[si:], "*/")
		if ei < 0 {
			text = text[:si]
			break
		}
		text = text[:si] + " " + text[si+ei+2:]
	}
	// Strip // line comments.
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		if ci := strings.Index(l, "//"); ci >= 0 {
			lines[i] = l[:ci]
		}
	}
	text = strings.Join(lines, "\n")

	findMatchingBrace := func(src string, open int) int {
		depth := 0
		for i := open; i < len(src); i++ {
			switch src[i] {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					return i
				}
			}
		}
		return -1
	}

	count := 0
	rest := text
	for {
		openIdx := strings.IndexByte(rest, '{')
		if openIdx < 0 {
			break
		}
		selectorRaw := strings.TrimSpace(rest[:openIdx])
		closeIdx := findMatchingBrace(rest, openIdx)
		if closeIdx < 0 {
			break
		}
		body := rest[openIdx+1 : closeIdx]
		rest = rest[closeIdx+1:]

		if selectorRaw == "" {
			continue
		}

		if strings.HasPrefix(strings.ToLower(selectorRaw), "@media") {
			// Media queries are not evaluated yet against viewport in this parser.
			// Ignore their blocks for now instead of applying them unconditionally.
			continue
		}

		for _, sel := range strings.Split(selectorRaw, ",") {
			sel = strings.TrimSpace(sel)
			if sel == "" {
				continue
			}
			cls := normalizeCSSClass(sel)
			if cls == "" {
				continue
			}
			for _, decl := range strings.Split(body, ";") {
				parts := strings.SplitN(strings.TrimSpace(decl), ":", 2)
				if len(parts) != 2 {
					continue
				}
				prop := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				if prop == "" {
					continue
				}
				ss.SetRule(cls, prop, val)
				count++
			}
		}
	}
	return count
}

// LoadFile reads a CSS file from disk and parses it.
func (ss *StyleSheet) LoadFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return ss.ParseCSSText(string(data)), nil
}

// ─── Style resolution ─────────────────────────────────────────────────────────

// ResolveStyle builds a flat map of canonical CSS property → value string from:
//  1. Class-level rules (left-to-right from the "class" prop, later overrides earlier).
//  2. Inline style.* node props (highest priority, override class rules).
//
// The returned map uses standard CSS property names:
// "color", "background-color", "background", "font-size", "font-weight",
// "font-style", "text-align", "padding", "padding-top", etc.
func ResolveStyle(props map[string]any, appSS *StyleSheet, viewportW int) map[string]string {
	cacheKey, cacheable := resolveStyleCacheKeyForProps(props, appSS, viewportW)
	if cacheable {
		resolveStyleCacheMu.RLock()
		if cached, ok := resolveStyleCache[cacheKey]; ok {
			resolveStyleCacheMu.RUnlock()
			return cloneStringMap(cached)
		}
		resolveStyleCacheMu.RUnlock()
	}

	resolved := make(map[string]string, 12)
	applyRule := func(m map[string]string) {
		if m == nil {
			return
		}
		for k, v := range m {
			resolved[CanonicalName(k)] = v
		}
	}

	// Global selectors baseline.
	if appSS != nil {
		if root := appSS.GetClassProps(":root"); root != nil {
			for k, v := range root {
				canonical := CanonicalName(k)
				if strings.HasPrefix(strings.TrimSpace(canonical), "--") {
					resolved[canonical] = v
				}
			}
		}
		applyRule(appSS.GetClassProps("*"))
		// `body` should not be applied wholesale to every node; that causes backgrounds,
		// margins and layout properties to leak into children. Keep only inheritable text
		// defaults as a practical approximation.
		if body := appSS.GetClassProps("body"); body != nil {
			for _, key := range []string{"color", "font-family", "font-size", "font-weight", "font-style", "line-height", "text-align"} {
				if v := strings.TrimSpace(body[key]); v != "" {
					if resolved[CanonicalName(key)] == "" {
						resolved[CanonicalName(key)] = v
					}
				}
			}
		}
	}

	// HTML-like default (UA) styles by semantic tag.
	// These defaults are low-precedence and can be overridden by class/inline styles.
	cssApplyTagDefaults(resolved, props)

	// Author tag selectors (e.g. h1 { ... }, p { ... }) should override UA defaults
	// and be overridden by class/inline styles, matching browser precedence intuition.
	if appSS != nil {
		tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
		if tag != "" {
			applyRule(appSS.GetClassProps(tag))
		}
	}

	// 1. Class styles (e.g. node.set("class", "title hero"))
	// Property names are normalized through CanonicalName so that CSS shorthand
	// aliases like `flex-direction`, `justify-content`, `align-items`, `bg` etc. are
	// stored under the canonical keys the layout/render engine reads.
	if appSS != nil {
		classStr := anyToString(props["class"], "")
		if classStr != "" {
			for _, cls := range strings.Fields(strings.ReplaceAll(classStr, ",", " ")) {
				cls = strings.TrimPrefix(strings.TrimSpace(cls), ".")
				if cls == "" {
					continue
				}
				applyRule(appSS.GetClassProps(cls))
			}
		}
	}

	// 2. Inline style.* overrides (e.g. node.style("color", "#fff"))
	for k, v := range props {
		if !strings.HasPrefix(k, "style.") {
			continue
		}
		prop := CanonicalName(strings.ToLower(strings.TrimPrefix(k, "style.")))
		if s := anyToString(v, ""); s != "" {
			resolved[prop] = s
		}
	}

	CSSExpandBoxShorthand(resolved, "padding")
	CSSExpandBoxShorthand(resolved, "margin")
	CSSExpandBorderShorthand(resolved)
	CSSResolveVariables(resolved)
	_ = viewportW
	if cacheable {
		resolveStyleCacheMu.Lock()
		if len(resolveStyleCache) > 4096 {
			resolveStyleCache = map[resolveStyleCacheKey]map[string]string{}
		}
		resolveStyleCache[cacheKey] = cloneStringMap(resolved)
		resolveStyleCacheMu.Unlock()
	}

	return resolved
}

func resolveStyleCacheKeyForProps(props map[string]any, appSS *StyleSheet, viewportW int) (resolveStyleCacheKey, bool) {
	if len(props) == 0 {
		return resolveStyleCacheKey{}, false
	}
	rev := int64(0)
	if appSS != nil {
		appSS.mu.RLock()
		rev = appSS.rev
		appSS.mu.RUnlock()
	}
	className := anyToString(props["class"], "")
	tag := anyToString(props["tag"], "")

	inline := make([]string, 0, 6)
	for k, v := range props {
		if !strings.HasPrefix(k, "style.") {
			continue
		}
		sv := strings.TrimSpace(anyToString(v, ""))
		if sv == "" {
			continue
		}
		inline = append(inline, k+"="+sv)
	}
	sort.Strings(inline)
	inlineSig := strings.Join(inline, ";")
	if len(className) > 240 || len(tag) > 80 || len(inlineSig) > 1200 {
		return resolveStyleCacheKey{}, false
	}
	return resolveStyleCacheKey{
		rev:       rev,
		viewportW: viewportW,
		className: className,
		tag:       tag,
		inlineSig: inlineSig,
	}, true
}

func cssApplyTagDefaults(resolved map[string]string, props map[string]any) {
	if resolved == nil {
		return
	}
	tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
	if tag == "" {
		return
	}
	setIfEmpty := func(k, v string) {
		k = CanonicalName(k)
		if strings.TrimSpace(resolved[k]) == "" {
			resolved[k] = v
		}
	}
	switch tag {
	case "h1":
		setIfEmpty("font-size", "2em")
		setIfEmpty("font-weight", "700")
		setIfEmpty("margin-top", "0.67em")
		setIfEmpty("margin-bottom", "0.67em")
	case "h2":
		setIfEmpty("font-size", "1.5em")
		setIfEmpty("font-weight", "700")
		setIfEmpty("margin-top", "0.83em")
		setIfEmpty("margin-bottom", "0.83em")
	case "h3":
		setIfEmpty("font-size", "1.17em")
		setIfEmpty("font-weight", "700")
		setIfEmpty("margin-top", "1em")
		setIfEmpty("margin-bottom", "1em")
	case "h4":
		setIfEmpty("font-size", "1em")
		setIfEmpty("font-weight", "700")
		setIfEmpty("margin-top", "1.33em")
		setIfEmpty("margin-bottom", "1.33em")
	case "p":
		setIfEmpty("margin-top", "1em")
		setIfEmpty("margin-bottom", "1em")
	case "label":
		setIfEmpty("font-weight", "500")
	case "span":
		// Span intentionally keeps current Text behavior (no extra defaults).
	}
}

func CSSResolveVariables(css map[string]string) {
	if len(css) == 0 {
		return
	}
	vars := make(map[string]string, 8)
	for k, v := range css {
		if strings.HasPrefix(strings.TrimSpace(k), "--") {
			vars[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	replaceVar := func(value string) string {
		out := value
		for {
			idx := strings.Index(out, "var(")
			if idx < 0 {
				break
			}
			start := idx + len("var(")
			depth := 1
			end := -1
			for i := start; i < len(out); i++ {
				switch out[i] {
				case '(':
					depth++
				case ')':
					depth--
					if depth == 0 {
						end = i
						i = len(out)
					}
				}
			}
			if end < 0 {
				break
			}
			inner := strings.TrimSpace(out[start:end])
			parts := strings.SplitN(inner, ",", 2)
			name := strings.TrimSpace(parts[0])
			repl := ""
			if v, ok := vars[name]; ok {
				repl = v
			} else if len(parts) == 2 {
				repl = strings.TrimSpace(parts[1])
			}
			out = out[:idx] + repl + out[end+1:]
		}
		return out
	}

	for i := 0; i < 4; i++ {
		changed := false
		for k, v := range css {
			nv := replaceVar(v)
			if nv != v {
				css[k] = nv
				changed = true
			}
			if strings.HasPrefix(strings.TrimSpace(k), "--") {
				vars[strings.TrimSpace(k)] = css[k]
			}
		}
		if !changed {
			break
		}
	}
}

func CanonicalName(name string) string {
	trimmed := strings.ToLower(strings.TrimSpace(name))
	switch trimmed {
	case "inline-size":
		return "width"
	case "block-size":
		return "height"
	case "min-inline-size":
		return "min-width"
	case "max-inline-size":
		return "max-width"
	case "min-block-size":
		return "min-height"
	case "max-block-size":
		return "max-height"
	case "flex-direction":
		return "direction"
	case "flex-grow":
		return "flex"
	case "justify-content":
		return "justify"
	case "align-items":
		return "align"
	case "overflow-inline":
		return "overflow-x"
	case "overflow-block":
		return "overflow-y"
	case "padding-inline-start":
		return "padding-left"
	case "padding-inline-end":
		return "padding-right"
	case "padding-block-start":
		return "padding-top"
	case "padding-block-end":
		return "padding-bottom"
	case "margin-inline-start":
		return "margin-left"
	case "margin-inline-end":
		return "margin-right"
	case "margin-block-start":
		return "margin-top"
	case "margin-block-end":
		return "margin-bottom"
	case "minw":
		return "min-width"
	case "maxw":
		return "max-width"
	case "minh":
		return "min-height"
	case "maxh":
		return "max-height"
	case "radius":
		return "border-radius"
	case "rounded", "round", "corner-radius":
		return "border-radius"
	case "border-start-start-radius":
		return "border-top-left-radius"
	case "border-start-end-radius":
		return "border-top-right-radius"
	case "border-end-end-radius":
		return "border-bottom-right-radius"
	case "border-end-start-radius":
		return "border-bottom-left-radius"
	case "bg-image", "background-image":
		return "background-image"
	case "bg-position", "background-position":
		return "background-position"
	case "bg-size", "background-size":
		return "background-size"
	case "bg-repeat", "background-repeat":
		return "background-repeat"
	case "scrollbar-thumb":
		return "scrollbar-thumb-color"
	case "scrollbar-track":
		return "scrollbar-track-color"
	case "bg":
		return "background"
	default:
		return trimmed
	}
}

func cssExtractFunctionalColor(value string) string {
	lower := strings.ToLower(value)
	for _, fn := range []string{"rgba(", "rgb(", "hsla(", "hsl(", "cmyk("} {
		if idx := strings.Index(lower, fn); idx >= 0 {
			if end := strings.Index(value[idx:], ")"); end > 0 {
				return strings.TrimSpace(value[idx : idx+end+1])
			}
		}
	}
	return ""
}

func CSSExpandBoxShorthand(css map[string]string, prop string) {
	value := strings.TrimSpace(css[prop])
	if value == "" {
		return
	}
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return
	}
	top, right, bottom, left := parts[0], parts[0], parts[0], parts[0]
	if len(parts) == 2 {
		top, bottom = parts[0], parts[0]
		right, left = parts[1], parts[1]
	} else if len(parts) == 3 {
		top = parts[0]
		right, left = parts[1], parts[1]
		bottom = parts[2]
	} else if len(parts) >= 4 {
		top, right, bottom, left = parts[0], parts[1], parts[2], parts[3]
	}
	if css[prop+"-top"] == "" {
		css[prop+"-top"] = top
	}
	if css[prop+"-right"] == "" {
		css[prop+"-right"] = right
	}
	if css[prop+"-bottom"] == "" {
		css[prop+"-bottom"] = bottom
	}
	if css[prop+"-left"] == "" {
		css[prop+"-left"] = left
	}
}

func CSSExpandBorderShorthand(css map[string]string) {
	// Expand directional values (1..4 items) into side-specific keys.
	expandDirectional := func(prop string) {
		value := strings.TrimSpace(css[prop])
		if value == "" {
			return
		}
		parts := strings.Fields(value)
		if len(parts) == 0 {
			return
		}
		top, right, bottom, left := parts[0], parts[0], parts[0], parts[0]
		if len(parts) == 2 {
			top, bottom = parts[0], parts[0]
			right, left = parts[1], parts[1]
		} else if len(parts) == 3 {
			top = parts[0]
			right, left = parts[1], parts[1]
			bottom = parts[2]
		} else if len(parts) >= 4 {
			top, right, bottom, left = parts[0], parts[1], parts[2], parts[3]
		}
		if css[prop+"-top"] == "" {
			css[prop+"-top"] = top
		}
		if css[prop+"-right"] == "" {
			css[prop+"-right"] = right
		}
		if css[prop+"-bottom"] == "" {
			css[prop+"-bottom"] = bottom
		}
		if css[prop+"-left"] == "" {
			css[prop+"-left"] = left
		}
	}

	expandDirectional("border-width")
	expandDirectional("border-style")
	expandDirectional("border-color")

	// Helper: parse a "1px solid #color" value into its components.
	parseBV := func(val string) (widthStr, styleStr, colorStr string) {
		if fnCol := cssExtractFunctionalColor(val); fnCol != "" {
			colorStr = fnCol
			val = strings.Replace(val, fnCol, " ", 1)
		}
		for _, part := range strings.Fields(val) {
			lower := strings.ToLower(part)
			if widthStr == "" && (strings.HasSuffix(lower, "px") || strings.HasSuffix(lower, "em") || strings.HasSuffix(lower, "rem") || strings.HasSuffix(lower, "pt") || strings.HasSuffix(lower, "pc") || strings.HasSuffix(lower, "in") || strings.HasSuffix(lower, "cm") || strings.HasSuffix(lower, "mm") || strings.HasSuffix(lower, "%")) {
				widthStr = part
				continue
			}
			if styleStr == "" && (lower == "none" || lower == "solid" || lower == "dashed" || lower == "dotted" || lower == "double" || lower == "groove" || lower == "ridge") {
				styleStr = lower
				continue
			}
			if colorStr == "" {
				colorStr = part
			}
		}
		return
	}

	// 1. Process the `border` shorthand — sets global border-* values.
	if v := strings.TrimSpace(css["border"]); v != "" {
		w, s, c := parseBV(v)
		if css["border-width"] == "" && w != "" {
			css["border-width"] = w
		}
		if css["border-style"] == "" && s != "" {
			css["border-style"] = s
		}
		if css["border-color"] == "" && c != "" {
			css["border-color"] = c
		}
	}

	// 2. Process individual side shorthands: border-top, border-right, border-bottom, border-left.
	//    These set per-side values and do NOT touch the global border-width/border-color
	//    (which are used as fallback for sides that have NO per-side override).
	for _, side := range []string{"top", "right", "bottom", "left"} {
		sv := strings.TrimSpace(css["border-"+side])
		if sv == "" {
			continue
		}
		w, s, c := parseBV(sv)
		key := "border-" + side + "-"
		if css[key+"width"] == "" && w != "" {
			css[key+"width"] = w
		}
		if css[key+"style"] == "" && s != "" {
			css[key+"style"] = s
		}
		if css[key+"color"] == "" && c != "" {
			css[key+"color"] = c
		}
	}

	// If directional `border-width/style/color` were set, materialize side-specific
	// keys used by the renderer.
	for _, side := range []string{"top", "right", "bottom", "left"} {
		if css["border-"+side+"-width"] == "" && css["border-width-"+side] != "" {
			css["border-"+side+"-width"] = css["border-width-"+side]
		}
		if css["border-"+side+"-style"] == "" && css["border-style-"+side] != "" {
			css["border-"+side+"-style"] = css["border-style-"+side]
		}
		if css["border-"+side+"-color"] == "" && css["border-color-"+side] != "" {
			css["border-"+side+"-color"] = css["border-color-"+side]
		}
	}
}

// ─── CSS property accessors ───────────────────────────────────────────────────

// CSSGetColor returns the CSS property value for prop, or fallback if absent/empty.
func CSSGetColor(css map[string]string, prop, fallback string) string {
	if v := css[prop]; v != "" {
		return v
	}
	return fallback
}

// CSSBackground returns the background color from "background-color" or "background".
func CSSBackground(css map[string]string) string {
	if v := css["background"]; v != "" {
		return v
	}
	if v := css["background-image"]; v != "" {
		return v
	}
	return css["background-color"]
}

// CSSFontSize parses "font-size" from resolved map (supports px suffix).
func CSSFontSize(css map[string]string, fallback float32) float32 {
	v := css["font-size"]
	if v == "" {
		return fallback
	}
	trimmed := strings.ToLower(strings.TrimSpace(v))
	if strings.HasSuffix(trimmed, "rem") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(trimmed, "rem"), 32); err == nil && f > 0 {
			return float32(f * 16.0)
		}
	}
	if strings.HasSuffix(trimmed, "em") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(trimmed, "em"), 32); err == nil && f > 0 {
			return float32(f * 16.0)
		}
	}
	if strings.HasSuffix(trimmed, "pt") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(trimmed, "pt"), 32); err == nil && f > 0 {
			return float32(f * (96.0 / 72.0))
		}
	}
	s := strings.TrimSuffix(trimmed, "px")
	if f, err := strconv.ParseFloat(s, 32); err == nil && f > 0 {
		return float32(f)
	}
	return fallback
}

// CSSLength parses a CSS length property (bare number or px/em suffix) to int.
func CSSLength(css map[string]string, prop string, fallback int) int {
	return CSSLengthValue(css[prop], fallback, 0, 0, 0)
}

// CSSValueInt reads an integer length value from a resolved CSS map.
func CSSValueInt(css map[string]string, prop string, fallback int) int {
	return CSSLength(css, prop, fallback)
}

func CSSLengthValue(raw string, fallback int, basis int, viewportW int, viewportH int) int {
	v := strings.ToLower(strings.TrimSpace(raw))
	if v == "" {
		return fallback
	}
	splitArgs := func(input string) []string {
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
	if strings.HasPrefix(v, "min(") && strings.HasSuffix(v, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(v, "min("), ")")
		args := splitArgs(inner)
		if len(args) > 0 {
			best := CSSLengthValue(args[0], fallback, basis, viewportW, viewportH)
			for i := 1; i < len(args); i++ {
				cand := CSSLengthValue(args[i], best, basis, viewportW, viewportH)
				if cand < best {
					best = cand
				}
			}
			return best
		}
	}
	if strings.HasPrefix(v, "max(") && strings.HasSuffix(v, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(v, "max("), ")")
		args := splitArgs(inner)
		if len(args) > 0 {
			best := CSSLengthValue(args[0], fallback, basis, viewportW, viewportH)
			for i := 1; i < len(args); i++ {
				cand := CSSLengthValue(args[i], best, basis, viewportW, viewportH)
				if cand > best {
					best = cand
				}
			}
			return best
		}
	}
	if strings.HasPrefix(v, "clamp(") && strings.HasSuffix(v, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(v, "clamp("), ")")
		args := splitArgs(inner)
		if len(args) == 3 {
			minV := CSSLengthValue(args[0], fallback, basis, viewportW, viewportH)
			prefV := CSSLengthValue(args[1], minV, basis, viewportW, viewportH)
			maxV := CSSLengthValue(args[2], prefV, basis, viewportW, viewportH)
			if prefV < minV {
				return minV
			}
			if prefV > maxV {
				return maxV
			}
			return prefV
		}
	}
	round := func(f float64) int {
		if f < 0 {
			return int(f - 0.5)
		}
		return int(f + 0.5)
	}
	parseNumeric := func(text string) (float64, bool) {
		if f, err := strconv.ParseFloat(strings.TrimSpace(text), 64); err == nil {
			return f, true
		}
		return 0, false
	}
	if strings.HasSuffix(v, "vw") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "vw")); ok && viewportW > 0 {
			return round((f / 100.0) * float64(viewportW))
		}
		return fallback
	}
	if strings.HasSuffix(v, "vh") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "vh")); ok && viewportH > 0 {
			return round((f / 100.0) * float64(viewportH))
		}
		return fallback
	}
	if strings.HasSuffix(v, "%") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "%")); ok && basis > 0 {
			return round((f / 100.0) * float64(basis))
		}
		return fallback
	}
	if strings.HasSuffix(v, "rem") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "rem")); ok {
			return round(f * 16.0)
		}
		return fallback
	}
	if strings.HasSuffix(v, "em") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "em")); ok {
			// Without font-size context in this helper, use a stable baseline.
			return round(f * 16.0)
		}
		return fallback
	}
	if strings.HasSuffix(v, "pt") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "pt")); ok {
			return round(f * (96.0 / 72.0))
		}
		return fallback
	}
	if strings.HasSuffix(v, "pc") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "pc")); ok {
			return round(f * 16.0)
		}
		return fallback
	}
	if strings.HasSuffix(v, "in") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "in")); ok {
			return round(f * 96.0)
		}
		return fallback
	}
	if strings.HasSuffix(v, "cm") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "cm")); ok {
			return round(f * (96.0 / 2.54))
		}
		return fallback
	}
	if strings.HasSuffix(v, "mm") {
		if f, ok := parseNumeric(strings.TrimSuffix(v, "mm")); ok {
			return round(f * (96.0 / 25.4))
		}
		return fallback
	}
	s := strings.TrimSuffix(v, "px")
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return round(f)
	}
	return fallback
}

func CSSFloatValue(raw string, fallback float64) float64 {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

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

// CSSBold returns true when font-weight indicates bold.
func CSSBold(css map[string]string) bool {
	v := strings.ToLower(strings.TrimSpace(css["font-weight"]))
	return v == "bold" || v == "700" || v == "800" || v == "900" || v == "bolder"
}

// CSSItalic returns true when font-style indicates italic/oblique.
func CSSItalic(css map[string]string) bool {
	v := strings.ToLower(strings.TrimSpace(css["font-style"]))
	return strings.Contains(v, "italic") || strings.Contains(v, "oblique")
}

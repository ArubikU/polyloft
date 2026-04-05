package layout

import (
	"strconv"
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

// anyToString is a local helper to avoid external dep for simple conversion.
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

// Node is a lightweight pseudo-HTML element model owned by giocss.
type Node struct {
	Tag      string
	Text     string
	Props    map[string]any
	Children []*Node
}

func NewNode(tag string) *Node {
	return &Node{Tag: strings.ToLower(strings.TrimSpace(tag)), Props: map[string]any{}}
}

func (n *Node) SetProp(name string, value any) {
	if n == nil {
		return
	}
	if n.Props == nil {
		n.Props = map[string]any{}
	}
	n.Props[name] = value
}

func (n *Node) GetProp(name string) any {
	if n == nil || n.Props == nil {
		return nil
	}
	return n.Props[name]
}

func (n *Node) AddChild(child *Node) {
	if n == nil || child == nil {
		return
	}
	n.Children = append(n.Children, child)
}

func (n *Node) AddClass(className string) {
	if n == nil {
		return
	}
	existing := strings.TrimSpace(anyToString(n.GetProp("class"), ""))
	className = strings.TrimSpace(className)
	if className == "" {
		return
	}
	if existing == "" {
		n.SetProp("class", className)
		return
	}
	n.SetProp("class", existing+" "+className)
}

func ResolveNodeStyle(node *Node, ss *uicore.StyleSheet, viewportW int) map[string]string {
	if node == nil {
		return map[string]string{}
	}
	props := make(map[string]any, len(node.Props)+1)
	for k, v := range node.Props {
		props[k] = v
	}
	if strings.TrimSpace(node.Tag) != "" {
		props["tag"] = strings.ToLower(strings.TrimSpace(node.Tag))
	}
	return uicore.ResolveStyle(props, ss, viewportW)
}

func MergeInheritedTextCSS(css map[string]string, parent map[string]string) map[string]string {
	if css == nil {
		css = map[string]string{}
	}
	if parent == nil {
		return css
	}
	inheritable := []string{
		"color",
		"font-family",
		"font-size",
		"font-style",
		"font-weight",
		"line-height",
		"letter-spacing",
		"text-align",
		"text-transform",
		"white-space",
	}
	for _, k := range inheritable {
		v := strings.TrimSpace(css[k])
		if v == "" || strings.EqualFold(v, "inherit") {
			if pv := strings.TrimSpace(parent[k]); pv != "" {
				css[k] = pv
			}
		}
	}
	return css
}

func SplitSpaceOutsideParens(input string) []string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil
	}
	parts := make([]string, 0, 8)
	depth := 0
	start := -1
	for i, r := range input {
		switch r {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
		if start < 0 && !strings.ContainsRune(" \t\n\r", r) {
			start = i
		}
		if start >= 0 && depth == 0 && strings.ContainsRune(" \t\n\r", r) {
			part := strings.TrimSpace(input[start:i])
			if part != "" {
				parts = append(parts, part)
			}
			start = -1
		}
	}
	if start >= 0 {
		part := strings.TrimSpace(input[start:])
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func NodeLayoutString(node *Node, css map[string]string, propName string, cssName string, fallback string) string {
	if node != nil {
		if value := strings.TrimSpace(anyToString(node.GetProp(propName), "")); value != "" {
			return strings.ToLower(value)
		}
	}
	if css != nil {
		if value := strings.TrimSpace(css[cssName]); value != "" {
			return strings.ToLower(value)
		}
	}
	return fallback
}

func NodeLayoutLength(node *Node, css map[string]string, propName string, cssName string, basis int, viewportW int, viewportH int, fallback int) int {
	if node != nil {
		if candidate := node.GetProp(propName); candidate != nil {
			switch typed := candidate.(type) {
			case int:
				return typed
			case int64:
				return int(typed)
			case float64:
				return int(typed)
			case float32:
				return int(typed)
			case string:
				return uicore.CSSLengthValue(typed, fallback, basis, viewportW, viewportH)
			}
		}
	}
	if css != nil {
		if value := css[cssName]; value != "" {
			return uicore.CSSLengthValue(value, fallback, basis, viewportW, viewportH)
		}
	}
	return fallback
}

func NodeLayoutFloat(node *Node, css map[string]string, propName string, cssName string, fallback float64) float64 {
	if node != nil {
		if candidate := node.GetProp(propName); candidate != nil {
			switch typed := candidate.(type) {
			case int:
				return float64(typed)
			case int64:
				return float64(typed)
			case float64:
				return typed
			case float32:
				return float64(typed)
			case string:
				return uicore.CSSFloatValue(typed, fallback)
			}
		}
	}
	if css != nil {
		if value := css[cssName]; value != "" {
			return uicore.CSSFloatValue(value, fallback)
		}
	}
	return fallback
}

func ParseGridTrackSpec(spec string, total int, viewportW int, viewportH int) []int {
	trimmed := strings.TrimSpace(strings.ToLower(spec))
	if trimmed == "" || total <= 0 {
		return []int{imax(1, total)}
	}
	if strings.HasPrefix(trimmed, "repeat(") && strings.HasSuffix(trimmed, ")") {
		inner := strings.TrimSuffix(strings.TrimPrefix(trimmed, "repeat("), ")")
		parts := strings.SplitN(inner, ",", 2)
		if len(parts) == 2 {
			count, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err == nil && count > 0 {
				unit := strings.TrimSpace(parts[1])
				items := make([]string, count)
				for i := 0; i < count; i++ {
					items[i] = unit
				}
				trimmed = strings.Join(items, " ")
			}
		}
	}
	tokens := SplitSpaceOutsideParens(trimmed)
	if len(tokens) == 0 {
		return []int{imax(1, total)}
	}
	tracks := make([]int, len(tokens))
	fixed := 0
	frIndexes := make([]int, 0, len(tokens))
	frTotal := 0.0
	for i, token := range tokens {
		if strings.HasSuffix(token, "fr") {
			f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(token, "fr")), 64)
			if err != nil || f <= 0 {
				f = 1
			}
			frIndexes = append(frIndexes, i)
			frTotal += f
			continue
		}
		basis := total
		if basis <= 0 {
			basis = imax(viewportW, viewportH)
		}
		v := uicore.CSSLengthValue(token, 0, basis, viewportW, viewportH)
		if v <= 0 {
			v = 1
		}
		tracks[i] = v
		fixed += v
	}
	remaining := total - fixed
	if remaining < len(frIndexes) {
		remaining = len(frIndexes)
	}
	if len(frIndexes) > 0 {
		if frTotal <= 0 {
			frTotal = float64(len(frIndexes))
		}
		for _, idx := range frIndexes {
			token := tokens[idx]
			f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(token, "fr")), 64)
			if err != nil || f <= 0 {
				f = 1
			}
			tracks[idx] = imax(1, int((float64(remaining)*f)/frTotal))
		}
	}
	return tracks
}

func imax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

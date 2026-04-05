package ui

import "strings"

func UIClassTokens(raw string) []string {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		cls := strings.TrimSpace(strings.TrimPrefix(part, "."))
		if cls == "" {
			continue
		}
		if _, ok := seen[cls]; ok {
			continue
		}
		seen[cls] = struct{}{}
		out = append(out, cls)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func ApplyHoverStyles(css map[string]string, props map[string]any, ss *StyleSheet, hovered bool) {
	if !hovered || ss == nil {
		return
	}
	for _, cls := range UIClassTokens(anyToString(props["class"], "")) {
		if m := ss.GetClassProps(cls + ":hover"); m != nil {
			for k, v := range m {
				css[CanonicalName(k)] = v
			}
		}
	}
	if strings.TrimSpace(anyToString(props["tag"], "")) != "" {
		tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
		if m := ss.GetClassProps(tag + ":hover"); m != nil {
			for k, v := range m {
				css[CanonicalName(k)] = v
			}
		}
	}
	CSSExpandBoxShorthand(css, "padding")
	CSSExpandBoxShorthand(css, "margin")
	CSSExpandBorderShorthand(css)
	CSSResolveVariables(css)
}

func ApplyActiveStyles(css map[string]string, props map[string]any, ss *StyleSheet, active bool) {
	if !active || ss == nil {
		return
	}
	for _, cls := range UIClassTokens(anyToString(props["class"], "")) {
		if m := ss.GetClassProps(cls + ":active"); m != nil {
			for k, v := range m {
				css[CanonicalName(k)] = v
			}
		}
	}
	if strings.TrimSpace(anyToString(props["tag"], "")) != "" {
		tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
		if m := ss.GetClassProps(tag + ":active"); m != nil {
			for k, v := range m {
				css[CanonicalName(k)] = v
			}
		}
	}
	CSSExpandBoxShorthand(css, "padding")
	CSSExpandBoxShorthand(css, "margin")
	CSSExpandBorderShorthand(css)
	CSSResolveVariables(css)
}

func HasHoverSelector(props map[string]any, ss *StyleSheet) bool {
	if ss == nil {
		return false
	}
	for _, cls := range UIClassTokens(anyToString(props["class"], "")) {
		if m := ss.GetClassProps(cls + ":hover"); len(m) > 0 {
			return true
		}
	}
	tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
	if tag != "" {
		if m := ss.GetClassProps(tag + ":hover"); len(m) > 0 {
			return true
		}
	}
	return false
}

func HasActiveSelector(props map[string]any, ss *StyleSheet) bool {
	if ss == nil {
		return false
	}
	for _, cls := range UIClassTokens(anyToString(props["class"], "")) {
		if m := ss.GetClassProps(cls + ":active"); len(m) > 0 {
			return true
		}
	}
	tag := strings.ToLower(strings.TrimSpace(anyToString(props["tag"], "")))
	if tag != "" {
		if m := ss.GetClassProps(tag + ":active"); len(m) > 0 {
			return true
		}
	}
	return false
}

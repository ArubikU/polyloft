package layout

import (
	"strconv"
	"strings"

	pkgcore "github.com/ArubikU/giocss/pkg"
	uicore "github.com/ArubikU/giocss/ui"
)

func IntrinsicNodeSize(node *Node, ss *uicore.StyleSheet, viewportW int, viewportH int, parentCSS map[string]string) (int, int) {
	if node == nil {
		return 0, 0
	}
	css := ResolveNodeStyle(node, ss, viewportW)
	css = MergeInheritedTextCSS(css, parentCSS)
	if strings.EqualFold(strings.TrimSpace(css["display"]), "none") {
		return 0, 0
	}
	hasExplicitWidth := strings.TrimSpace(css["width"]) != ""
	hasExplicitHeight := strings.TrimSpace(css["height"]) != ""
	width := NodeLayoutLength(node, css, "width", "width", viewportW, viewportW, viewportH, -1)
	height := NodeLayoutLength(node, css, "height", "height", viewportH, viewportW, viewportH, -1)
	paddingTop := NodeLayoutLength(node, css, "paddingTop", "padding-top", viewportH, viewportW, viewportH, nodePropInt(node, "pady", nodePropInt(node, "padding", 0)))
	paddingRight := NodeLayoutLength(node, css, "paddingRight", "padding-right", viewportW, viewportW, viewportH, nodePropInt(node, "padx", nodePropInt(node, "padding", 0)))
	paddingBottom := NodeLayoutLength(node, css, "paddingBottom", "padding-bottom", viewportH, viewportW, viewportH, nodePropInt(node, "pady", nodePropInt(node, "padding", 0)))
	paddingLeft := NodeLayoutLength(node, css, "paddingLeft", "padding-left", viewportW, viewportW, viewportH, nodePropInt(node, "padx", nodePropInt(node, "padding", 0)))
	gap := NodeLayoutLength(node, css, "gap", "gap", max(viewportW, viewportH), viewportW, viewportH, 0)
	kind := strings.ToLower(strings.TrimSpace(node.Tag))

	if width >= 0 && height >= 0 {
		return width, height
	}

	switch kind {
	case "text", "label", "span", "p", "h1", "h2", "h3", "h4", "h5", "h6":
		fontSize := uicore.CSSFontSize(css, 14)
		text := uicore.CSSApplyLetterSpacing(uicore.CSSTextTransform(nodePropString(node, "text", ""), css), css)
		if width > 0 && uicore.CSSAllowsWrap(css) && !pkgcore.ShouldClipText(css) {
			inner := max(1, width-paddingLeft-paddingRight)
			text = uicore.WrapTextToWidth(text, fontSize, inner)
		}
		w, h := uicore.EstimateTextLayoutBox(text, fontSize, css)
		if width < 0 {
			width = w + paddingLeft + paddingRight
		}
		if height < 0 {
			height = h + paddingTop + paddingBottom
		}
	case "button":
		fontSize := uicore.CSSFontSize(css, 13)
		text := uicore.CSSApplyLetterSpacing(uicore.CSSTextTransform(nodePropString(node, "text", ""), css), css)
		w, h := uicore.EstimateTextLayoutBox(text, fontSize, css)
		iconExtra := 0
		if nodePropString(node, "icon", "") != "" {
			iconExtra = max(18, int(fontSize)+4)
		}
		if width < 0 {
			width = w + iconExtra + paddingLeft + paddingRight + 28
		}
		if height < 0 {
			height = max(30, h+paddingTop+paddingBottom+10)
		}
	case "input":
		inputType := strings.ToLower(nodePropString(node, "type", "text"))
		if width < 0 {
			if inputType == "checkbox" || inputType == "check" {
				width = 120
			} else if inputType == "range" || inputType == "slider" {
				width = 220
			} else {
				width = 180
			}
		}
		if height < 0 {
			if inputType == "textarea" || inputType == "multiline" {
				height = 92
			} else {
				height = 36
			}
		}
	case "native":
		component := strings.ToLower(nodePropString(node, "component", nodePropString(node, "native", "label")))
		if width < 0 {
			switch component {
			case "image", "img", "svg":
				width = 180
			case "icon":
				width = 18
			case "slider":
				width = 220
			case "progress", "progressbar":
				width = 220
			case "select", "dropdown":
				width = 180
			default:
				width = 120
			}
		}
		if height < 0 {
			switch component {
			case "image", "img", "svg":
				height = 120
			case "icon":
				height = 18
			case "progress", "progressbar":
				height = 14
			default:
				height = 36
			}
		}
	default:
		direction := NodeLayoutString(node, css, "direction", "direction", strings.ToLower(nodePropString(node, "layout", "column")))
		if direction != "row" {
			direction = "column"
		}
		contentW := 0
		contentH := 0
		for _, child := range node.Children {
			cw, ch := IntrinsicNodeSize(child, ss, viewportW, viewportH, css)
			childCSS := ResolveNodeStyle(child, ss, viewportW)
			childCSS = MergeInheritedTextCSS(childCSS, css)
			mw := NodeLayoutLength(child, childCSS, "marginLeft", "margin-left", viewportW, viewportW, viewportH, 0) + NodeLayoutLength(child, childCSS, "marginRight", "margin-right", viewportW, viewportW, viewportH, 0)
			mh := NodeLayoutLength(child, childCSS, "marginTop", "margin-top", viewportH, viewportW, viewportH, 0) + NodeLayoutLength(child, childCSS, "marginBottom", "margin-bottom", viewportH, viewportW, viewportH, 0)
			if direction == "row" {
				contentW += cw + mw
				contentH = max(contentH, ch+mh)
			} else {
				contentH += ch + mh
				contentW = max(contentW, cw+mw)
			}
		}
		visibleChildren := len(node.Children)
		if visibleChildren > 1 {
			if direction == "row" {
				contentW += gap * (visibleChildren - 1)
			} else {
				contentH += gap * (visibleChildren - 1)
			}
		}
		if width < 0 {
			width = contentW + paddingLeft + paddingRight
		}
		if height < 0 {
			height = contentH + paddingTop + paddingBottom
		}
	}

	if ar, ok := uicore.ParseAspectRatio(css); ok && ar > 0.001 {
		switch {
		case width > 0 && !hasExplicitHeight:
			height = max(1, int(float64(width)/ar))
		case height > 0 && !hasExplicitWidth:
			width = max(1, int(float64(height)*ar))
		case width <= 0 && height <= 0:
			baseW := max(96, paddingLeft+paddingRight+96)
			width = baseW
			height = max(1, int(float64(baseW)/ar))
		}
	}

	return max(1, width), max(1, height)
}

func nodePropString(node *Node, name string, fallback string) string {
	if node == nil {
		return fallback
	}
	if strings.EqualFold(strings.TrimSpace(name), "text") {
		if strings.TrimSpace(node.Text) != "" {
			return node.Text
		}
	}
	v := anyToString(node.GetProp(name), "")
	if v == "" {
		return fallback
	}
	return v
}

func nodePropInt(node *Node, name string, fallback int) int {
	if node == nil {
		return fallback
	}
	switch typed := node.GetProp(name).(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case float32:
		return int(typed)
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(typed)); err == nil {
			return parsed
		}
	}
	return fallback
}

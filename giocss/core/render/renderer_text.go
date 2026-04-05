package render

import (
	"strings"

	pkgcore "github.com/ArubikU/giocss/pkg"
	uicore "github.com/ArubikU/giocss/ui"
)

type TextContentPlan struct {
	Text        string
	FontSize    float32
	WrapAllowed bool
	MaxLines    int
	UseMono     bool
	Bold        bool
	Italic      bool
}

func DefaultTextFallbackSize(kind string) float32 {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "h1":
		return 32
	case "h2":
		return 26
	case "h3":
		return 22
	case "h4":
		return 18
	case "h5":
		return 16
	case "h6":
		return 14
	case "p":
		return 14
	case "span", "label", "text":
		return 13
	default:
		return 14
	}
}

func BuildTextContentPlan(props map[string]any, css map[string]string, width int, fallbackSize float32) TextContentPlan {
	textValue := uicore.CSSTextTransform(anyToString(props["text"], ""), css)
	textValue = uicore.CSSApplyLetterSpacing(textValue, css)
	textValue = uicore.SanitizeRenderText(textValue)
	fontSize := uicore.CSSFontSize(css, fallbackSize)
	wrapAllowed := uicore.CSSAllowsWrap(css)
	if pkgcore.ShouldClipText(css) {
		textValue = uicore.FitTextToWidth(textValue, fontSize, css, width)
		wrapAllowed = false
	}
	fontFamily := strings.ToLower(strings.TrimSpace(css["font-family"]))
	maxLines := 0
	if !wrapAllowed {
		maxLines = 1
	}
	return TextContentPlan{
		Text:        textValue,
		FontSize:    fontSize,
		WrapAllowed: wrapAllowed,
		MaxLines:    maxLines,
		UseMono:     strings.Contains(fontFamily, "mono"),
		Bold:        uicore.CSSBold(css),
		Italic:      uicore.CSSItalic(css),
	}
}

type ButtonContentPlan struct {
	Label    string
	FontSize float32
	TextX    int
	TextY    int
	TextW    int
	TextH    int
	Bold     bool
	Italic   bool
}

func BuildButtonContentPlan(props map[string]any, css map[string]string, x, y, w, h int, fallbackSize float32) ButtonContentPlan {
	label := uicore.SanitizeRenderText(anyToString(props["text"], "Button"))
	label = uicore.CSSTextTransform(label, css)
	fontSize := uicore.CSSFontSize(css, fallbackSize)
	padL := uicore.CSSLengthValue(css["padding-left"], 0, w, w, h)
	padR := uicore.CSSLengthValue(css["padding-right"], 0, w, w, h)
	padT := uicore.CSSLengthValue(css["padding-top"], 0, h, w, h)
	padB := uicore.CSSLengthValue(css["padding-bottom"], 0, h, w, h)

	// Guard against double-accounting CSS padding between layout and paint plans.
	// If the effective text box becomes too small, fallback to a centered full box.
	maxInsetX := max(0, w/4)
	maxInsetY := max(0, h/4)
	effPadL := min(padL, maxInsetX)
	effPadR := min(padR, maxInsetX)
	effPadT := min(padT, maxInsetY)
	effPadB := min(padB, maxInsetY)

	textX := x + effPadL
	textY := y + effPadT
	textW := max(1, w-effPadL-effPadR)
	textH := max(1, h-effPadT-effPadB)
	minTextW := max(24, int(fontSize)*2)
	if textW < minTextW {
		textX = x
		textY = y
		textW = max(1, w)
		textH = max(1, h)
	}

	return ButtonContentPlan{
		Label:    label,
		FontSize: fontSize,
		TextX:    textX,
		TextY:    textY,
		TextW:    textW,
		TextH:    textH,
		Bold:     uicore.CSSBold(css),
		Italic:   uicore.CSSItalic(css),
	}
}

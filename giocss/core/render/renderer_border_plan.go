package render

import (
	"image/color"
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

type BorderSidePlan struct {
	Width int
	Color color.NRGBA
	Style string
}

type ElementBorderPlan struct {
	Top          BorderSidePlan
	Right        BorderSidePlan
	Bottom       BorderSidePlan
	Left         BorderSidePlan
	HasAny       bool
	Uniform      bool
	UniformW     int
	UniformCol   color.NRGBA
	UniformStyle string
}

func BuildElementBorderPlan(css map[string]string, w int, h int) ElementBorderPlan {
	getSide := func(name string) BorderSidePlan {
		wStr := css["border-"+name+"-width"]
		if wStr == "" {
			wStr = css["border-width"]
		}
		cStr := css["border-"+name+"-color"]
		if cStr == "" {
			cStr = css["border-color"]
		}
		sStyle := strings.ToLower(strings.TrimSpace(css["border-"+name+"-style"]))
		if sStyle == "" {
			sStyle = strings.ToLower(strings.TrimSpace(css["border-style"]))
		}
		bw := uicore.CSSLengthValue(wStr, 0, max(w, h), w, h)
		bc := toNRGBA(uicore.ParseHexColor(cStr, color.NRGBA{A: 255}))
		if sStyle == "none" || sStyle == "hidden" {
			bw = 0
		}
		return BorderSidePlan{Width: bw, Color: bc, Style: sStyle}
	}

	top := getSide("top")
	right := getSide("right")
	bottom := getSide("bottom")
	left := getSide("left")

	hasAny := top.Width > 0 || right.Width > 0 || bottom.Width > 0 || left.Width > 0
	uniform := top.Width == right.Width && right.Width == bottom.Width && bottom.Width == left.Width &&
		top.Color == right.Color && right.Color == bottom.Color && bottom.Color == left.Color &&
		top.Style == right.Style && right.Style == bottom.Style && bottom.Style == left.Style

	return ElementBorderPlan{
		Top:          top,
		Right:        right,
		Bottom:       bottom,
		Left:         left,
		HasAny:       hasAny,
		Uniform:      uniform,
		UniformW:     top.Width,
		UniformCol:   top.Color,
		UniformStyle: strings.ToLower(strings.TrimSpace(top.Style)),
	}
}

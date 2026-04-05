package render

import (
	"image"
	"image/color"
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

type PaintRect struct {
	X       int
	Y       int
	W       int
	H       int
	Radius  int
	Rounded bool
}

type TextDecorationPlan struct {
	Color color.NRGBA
	Rects []PaintRect
}

type ImagePaintPlan struct {
	Clip   image.Rectangle
	Offset image.Point
	ScaleX float32
	ScaleY float32
}

func BuildDecorationLineRects(x int, y int, w int, thickness int, style string) []PaintRect {
	if w <= 0 || thickness <= 0 {
		return nil
	}
	style = strings.ToLower(strings.TrimSpace(style))
	rects := make([]PaintRect, 0, 16)
	switch style {
	case "dotted":
		dot := max(1, thickness)
		gap := max(1, thickness)
		for cx := x; cx < x+w; cx += dot + gap {
			dw := min(dot, x+w-cx)
			rects = append(rects, PaintRect{X: cx, Y: y, W: dw, H: thickness, Radius: thickness / 2, Rounded: true})
		}
	case "dashed":
		dash := max(2, thickness*3)
		gap := max(1, thickness*2)
		for cx := x; cx < x+w; cx += dash + gap {
			dw := min(dash, x+w-cx)
			rects = append(rects, PaintRect{X: cx, Y: y, W: dw, H: thickness})
		}
	default:
		rects = append(rects, PaintRect{X: x, Y: y, W: w, H: thickness})
	}
	return rects
}

func BuildDashedBorderRects(x int, y int, w int, h int, borderW int, dotted bool) []PaintRect {
	if borderW <= 0 || w <= 0 || h <= 0 {
		return nil
	}
	dashLen := max(1, borderW*3)
	gapLen := max(1, borderW*2)
	if dotted {
		dashLen = max(1, borderW)
		gapLen = max(1, borderW*2)
	}
	rects := make([]PaintRect, 0, 64)
	for cx := x; cx < x+w; cx += dashLen + gapLen {
		dw := min(dashLen, x+w-cx)
		rects = append(rects, PaintRect{X: cx, Y: y, W: dw, H: min(borderW, h)})
		if h > borderW {
			rects = append(rects, PaintRect{X: cx, Y: y + h - borderW, W: dw, H: min(borderW, h)})
		}
	}
	for cy := y; cy < y+h; cy += dashLen + gapLen {
		dh := min(dashLen, y+h-cy)
		rects = append(rects, PaintRect{X: x, Y: cy, W: min(borderW, w), H: dh})
		if w > borderW {
			rects = append(rects, PaintRect{X: x + w - borderW, Y: cy, W: min(borderW, w), H: dh})
		}
	}
	return rects
}

func BuildImagePaintPlan(x int, y int, w int, h int, imgW int, imgH int) (ImagePaintPlan, bool) {
	if w <= 0 || h <= 0 || imgW <= 0 || imgH <= 0 {
		return ImagePaintPlan{}, false
	}
	return ImagePaintPlan{
		Clip:   image.Rect(x, y, x+w, y+h),
		Offset: image.Pt(x, y),
		ScaleX: float32(w) / float32(imgW),
		ScaleY: float32(h) / float32(imgH),
	}, true
}

func BuildDebugOutlineRects(x int, y int, w int, h int, stroke int) []PaintRect {
	if stroke <= 0 || w <= 0 || h <= 0 {
		return nil
	}
	rects := make([]PaintRect, 0, 4)
	rects = append(rects, PaintRect{X: x, Y: y, W: w, H: min(stroke, h)})
	if h > stroke {
		rects = append(rects, PaintRect{X: x, Y: y + h - stroke, W: w, H: min(stroke, h)})
	}
	innerH := h - 2*stroke
	if innerH > 0 {
		rects = append(rects, PaintRect{X: x, Y: y + stroke, W: min(stroke, w), H: innerH})
		if w > stroke {
			rects = append(rects, PaintRect{X: x + w - stroke, Y: y + stroke, W: min(stroke, w), H: innerH})
		}
	}
	return rects
}

func BuildTextDecorationPlan(x int, y int, w int, h int, fontSize float32, align string, txt string, base color.NRGBA, deco uicore.TextDecorationInfo) TextDecorationPlan {
	line := strings.ToLower(strings.TrimSpace(deco.Line))
	if line == "" || line == "none" {
		return TextDecorationPlan{}
	}

	thickness := uicore.CSSLengthValue(strings.TrimSpace(deco.Thickness), max(1, int(fontSize*0.08)+1), max(w, h), w, h)
	if thickness <= 0 {
		thickness = 1
	}
	lineColor := base
	if strings.TrimSpace(deco.Color) != "" {
		lineColor = toNRGBA(uicore.ParseHexColor(deco.Color, base))
	}
	lineStyle := strings.ToLower(strings.TrimSpace(deco.Style))

	firstLine := txt
	if idx := strings.Index(firstLine, "\n"); idx >= 0 {
		firstLine = firstLine[:idx]
	}
	firstLine = strings.TrimSpace(firstLine)
	if firstLine == "" {
		return TextDecorationPlan{}
	}
	runeCount := len([]rune(firstLine))
	if runeCount <= 0 {
		return TextDecorationPlan{}
	}

	avgCharW := float64(fontSize) * 0.56
	lineW := int(float64(runeCount)*avgCharW) + max(2, int(float64(fontSize)*0.22))
	lineW = min(w, max(1, lineW))
	lineX := x
	switch strings.ToLower(strings.TrimSpace(align)) {
	case "middle", "center":
		lineX = x + max(0, (w-lineW)/2)
	case "end", "right":
		lineX = x + max(0, w-lineW)
	}

	lineH := max(1, int(float64(fontSize)*1.28))
	if lineH > h {
		lineH = h
	}
	lineTop := y + max(0, (h-lineH)/2)

	plan := TextDecorationPlan{Color: lineColor, Rects: make([]PaintRect, 0, 24)}
	parts := strings.Fields(strings.ReplaceAll(line, ",", " "))
	for _, p := range parts {
		switch p {
		case "underline":
			uy := min(y+h-thickness, lineTop+lineH-max(1, int(fontSize*0.10)))
			plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, uy, lineW, thickness, lineStyle)...)
			if lineStyle == "double" {
				plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, min(y+h-thickness, uy+thickness+1), lineW, thickness, "solid")...)
			}
		case "overline":
			oy := lineTop + max(0, int(fontSize*0.06))
			plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, oy, lineW, thickness, lineStyle)...)
			if lineStyle == "double" {
				plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, oy+thickness+1, lineW, thickness, "solid")...)
			}
		case "line-through":
			sy := min(y+h-thickness, lineTop+int(float32(lineH)*0.50))
			plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, sy, lineW, thickness, lineStyle)...)
			if lineStyle == "double" {
				plan.Rects = append(plan.Rects, BuildDecorationLineRects(lineX, sy+thickness+1, lineW, thickness, "solid")...)
			}
		}
	}
	if len(plan.Rects) == 0 {
		return TextDecorationPlan{}
	}
	return plan
}

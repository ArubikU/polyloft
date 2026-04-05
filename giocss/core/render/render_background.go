package render

import (
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	coreengine "github.com/ArubikU/giocss/core/engine"
	uicore "github.com/ArubikU/giocss/ui"
)

// drawGioBackground draws the CSS background (background-color + border-radius).
func drawGioBackground(ops *op.Ops, x, y, w, h int, css map[string]string) {
	bgRaw := strings.TrimSpace(uicore.CSSBackground(css))
	if bgRaw == "" {
		return
	}
	radii := cssBorderRadiiValues(css, w, h)
	layers := backgroundLayers(bgRaw)
	filters := ParseBackgroundFilterFactors(css)
	// Paint from back to front: last layer first, first layer on top.
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		if layer == "" {
			continue
		}
		drawSingleBackgroundLayer(ops, x, y, w, h, radii, layer, css, filters)
	}
}

func DrawGioBackground(rctx *coreengine.RenderContext, x, y, w, h int, css map[string]string) {
	drawGioBackground(rctx.Gtx.Ops, x, y, w, h, css)
}

func drawSingleBackgroundLayer(ops *op.Ops, x, y, w, h int, radii cornerRadii, layer string, css map[string]string, filters BackgroundFilterFactors) {
	initRenderConfig()
	if renderLowPower {
		lower := strings.ToLower(strings.TrimSpace(layer))
		if strings.Contains(lower, "gradient(") {
			fallback := toNRGBA(uicore.ParseHexColor(layer, color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0xFF}))
			drawGioRRectCorners(ops, x, y, w, h, radii, fallback)
			return
		}
	}

	lower := strings.ToLower(strings.TrimSpace(layer))
	isLinear := strings.HasPrefix(lower, "linear-gradient(") || strings.HasPrefix(lower, "repeating-linear-gradient(")
	isRadial := strings.HasPrefix(lower, "radial-gradient(") || strings.HasPrefix(lower, "repeating-radial-gradient(")
	isRepeating := strings.HasPrefix(lower, "repeating-linear-gradient(") || strings.HasPrefix(lower, "repeating-radial-gradient(")
	if isLinear && isRepeating {
		if pxStops, ok := parseRepeatingLinearPxStops(layer); ok {
			drawGioRepeatingLinearGradientPx(ops, x, y, w, h, radii, layer, pxStops)
			return
		}
	}
	if isLinear || isRadial {
		stops := gradientColorStops(layer)
		if len(stops) >= 2 {
			// Apply all filters to gradient stops.
			if filters.Brightness != 1.0 || filters.Contrast != 1.0 || filters.Saturate != 1.0 || filters.Grayscale > 0.001 || filters.Invert > 0.001 {
				for i := range stops {
					col := stops[i].Color
					if filters.Brightness != 1.0 {
						col = uicore.ApplyBrightnessToColor(col, filters.Brightness)
					}
					if filters.Contrast != 1.0 {
						col = uicore.ApplyContrastToColor(col, filters.Contrast)
					}
					if filters.Saturate != 1.0 {
						col = uicore.ApplySaturationToColor(col, filters.Saturate)
					}
					if filters.Grayscale > 0.001 {
						col = uicore.ApplyGrayscaleToColor(col, filters.Grayscale)
					}
					if filters.Invert > 0.001 {
						col = uicore.ApplyInvertToColor(col, filters.Invert)
					}
					stops[i].Color = col
				}
			}
			if isRadial {
				drawGioRadialGradientApprox(ops, x, y, w, h, radii, layer, stops, isRepeating)
			} else {
				drawGioLinearGradientApprox(ops, x, y, w, h, radii, layer, stops, isRepeating)
			}
			return
		}
	}
	raw := uicore.ParseHexColor(layer, color.Transparent)
	col := toNRGBA(uicore.ApplyCSSOpacity(raw, css))
	if col.A == 0 {
		return
	}
	// Apply all filters to solid color.
	if filters.Brightness != 1.0 || filters.Contrast != 1.0 || filters.Saturate != 1.0 || filters.Grayscale > 0.001 || filters.Invert > 0.001 {
		if filters.Brightness != 1.0 {
			col = uicore.ApplyBrightnessToColor(col, filters.Brightness)
		}
		if filters.Contrast != 1.0 {
			col = uicore.ApplyContrastToColor(col, filters.Contrast)
		}
		if filters.Saturate != 1.0 {
			col = uicore.ApplySaturationToColor(col, filters.Saturate)
		}
		if filters.Grayscale > 0.001 {
			col = uicore.ApplyGrayscaleToColor(col, filters.Grayscale)
		}
		if filters.Invert > 0.001 {
			col = uicore.ApplyInvertToColor(col, filters.Invert)
		}
	}
	drawGioRRectCorners(ops, x, y, w, h, radii, col)
}

type gradientStop = GradientStop
type gradientStopPx = GradientStopPx
type gradientRasterKey = GradientRasterKey

func gradientStopsSignature(stops []gradientStop) string {
	return GradientStopsSignature(stops)
}

func gradientStopsPxSignature(stops []gradientStopPx) string {
	return GradientStopsPxSignature(stops)
}

func getCachedGradientRaster(key gradientRasterKey, build func() *image.NRGBA) *image.NRGBA {
	return GetCachedGradientRaster(key, build)
}

func gradientColorStops(raw string) []gradientStop {
	return GradientColorStops(raw)
}

func parseLinearGradientDirection(raw string) (float64, float64) {
	return ParseLinearGradientDirection(raw)
}

func parseRadialGradientCenter(raw string, x, y, w, h int) (int, int) {
	return ParseRadialGradientCenter(raw, x, y, w, h)
}

func gradientPaletteColor(stops []gradientStop, t float64) color.NRGBA {
	return GradientPaletteColor(stops, t)
}

func gradientPaletteColorRepeating(stops []gradientStop, t float64) color.NRGBA {
	return GradientPaletteColorRepeating(stops, t)
}

func parseRepeatingLinearPxStops(raw string) ([]gradientStopPx, bool) {
	return ParseRepeatingLinearPxStops(raw)
}

func gradientPaletteColorPx(stops []gradientStopPx, t float64) color.NRGBA {
	return GradientPaletteColorPx(stops, t)
}

func drawGioRepeatingLinearGradientPx(ops *op.Ops, x, y, w, h int, radii cornerRadii, raw string, stops []gradientStopPx) {
	if w <= 0 || h <= 0 || len(stops) < 2 {
		return
	}
	var cs clip.Stack
	if radii.max() > 0 {
		maxR := min(w, h) / 2
		cs = clip.RRect{
			Rect: image.Rect(x, y, x+w, y+h),
			NW:   min(max(0, radii.nw), maxR),
			NE:   min(max(0, radii.ne), maxR),
			SE:   min(max(0, radii.se), maxR),
			SW:   min(max(0, radii.sw), maxR),
		}.Push(ops)
	} else {
		cs = clip.Rect(image.Rect(x, y, x+w, y+h)).Push(ops)
	}

	dx, dy := parseLinearGradientDirection(raw)
	start := stops[0].Pos
	end := stops[len(stops)-1].Pos
	span := end - start
	if span <= 0 {
		span = 1
	}
	key := gradientRasterKey{
		Kind:     "repeating-linear-px",
		W:        w,
		H:        h,
		NW:       radii.nw,
		NE:       radii.ne,
		SE:       radii.se,
		SW:       radii.sw,
		Raw:      raw,
		StopsSig: gradientStopsPxSignature(stops),
	}
	img := getCachedGradientRaster(key, func() *image.NRGBA {
		out := image.NewNRGBA(image.Rect(0, 0, w, h))
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				p := float64(px)*dx + float64(py)*dy
				wrapped := start + math.Mod(p-start, span)
				if wrapped < start {
					wrapped += span
				}
				out.SetNRGBA(px, py, gradientPaletteColorPx(stops, wrapped))
			}
		}
		return out
	})
	if img == nil {
		cs.Pop()
		return
	}
	tr := op.Offset(image.Pt(x, y)).Push(ops)
	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{}.Add(ops)
	tr.Pop()
	cs.Pop()
}

func drawGioLinearGradientApprox(ops *op.Ops, x, y, w, h int, radii cornerRadii, raw string, stops []gradientStop, repeating bool) {
	if w <= 0 || h <= 0 {
		return
	}
	var cs clip.Stack
	if radii.max() > 0 {
		maxR := min(w, h) / 2
		cs = clip.RRect{
			Rect: image.Rect(x, y, x+w, y+h),
			NW:   min(max(0, radii.nw), maxR),
			NE:   min(max(0, radii.ne), maxR),
			SE:   min(max(0, radii.se), maxR),
			SW:   min(max(0, radii.sw), maxR),
		}.Push(ops)
	} else {
		cs = clip.Rect(image.Rect(x, y, x+w, y+h)).Push(ops)
	}

	dx, dy := parseLinearGradientDirection(raw)
	proj := func(px, py float64) float64 {
		return px*dx + py*dy
	}
	minP := proj(0, 0)
	maxP := minP
	for _, pt := range [][2]float64{{float64(w - 1), 0}, {0, float64(h - 1)}, {float64(w - 1), float64(h - 1)}} {
		p := proj(pt[0], pt[1])
		if p < minP {
			minP = p
		}
		if p > maxP {
			maxP = p
		}
	}
	den := maxP - minP
	if den <= 0 {
		drawGioRect(ops, x, y, w, h, stops[0].Color)
		cs.Pop()
		return
	}

	key := gradientRasterKey{
		Kind:     "linear",
		W:        w,
		H:        h,
		NW:       radii.nw,
		NE:       radii.ne,
		SE:       radii.se,
		SW:       radii.sw,
		Raw:      raw + "|rep:" + strconv.FormatBool(repeating),
		StopsSig: gradientStopsSignature(stops),
	}
	img := getCachedGradientRaster(key, func() *image.NRGBA {
		// Render 2D to preserve diagonal angles (e.g. 45deg) and multi-stop richness.
		out := image.NewNRGBA(image.Rect(0, 0, w, h))
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				p := proj(float64(px), float64(py))
				t := (p - minP) / den
				col := gradientPaletteColor(stops, t)
				if repeating {
					col = gradientPaletteColorRepeating(stops, t)
				}
				out.SetNRGBA(px, py, col)
			}
		}
		return out
	})
	if img == nil {
		cs.Pop()
		return
	}

	tr := op.Offset(image.Pt(x, y)).Push(ops)
	imgOp := paint.NewImageOp(img)
	imgOp.Add(ops)
	paint.PaintOp{}.Add(ops)
	tr.Pop()
	cs.Pop()
}

func drawGioRadialGradientApprox(ops *op.Ops, x, y, w, h int, radii cornerRadii, raw string, stops []gradientStop, repeating bool) {
	if w <= 0 || h <= 0 {
		return
	}
	var cs clip.Stack
	if radii.max() > 0 {
		maxR := min(w, h) / 2
		cs = clip.RRect{
			Rect: image.Rect(x, y, x+w, y+h),
			NW:   min(max(0, radii.nw), maxR),
			NE:   min(max(0, radii.ne), maxR),
			SE:   min(max(0, radii.se), maxR),
			SW:   min(max(0, radii.sw), maxR),
		}.Push(ops)
	} else {
		cs = clip.Rect(image.Rect(x, y, x+w, y+h)).Push(ops)
	}
	cx, cy := parseRadialGradientCenter(raw, 0, 0, w, h)
	maxRadius := math.Max(
		math.Hypot(float64(cx), float64(cy)),
		math.Max(
			math.Hypot(float64(cx-w), float64(cy)),
			math.Max(
				math.Hypot(float64(cx), float64(cy-h)),
				math.Hypot(float64(cx-w), float64(cy-h)),
			),
		),
	)
	if maxRadius <= 0 {
		cs.Pop()
		return
	}
	key := gradientRasterKey{
		Kind:     "radial",
		W:        w,
		H:        h,
		NW:       radii.nw,
		NE:       radii.ne,
		SE:       radii.se,
		SW:       radii.sw,
		Raw:      raw + "|rep:" + strconv.FormatBool(repeating),
		StopsSig: gradientStopsSignature(stops),
	}
	img := getCachedGradientRaster(key, func() *image.NRGBA {
		out := image.NewNRGBA(image.Rect(0, 0, w, h))
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				d := math.Hypot(float64(px-cx), float64(py-cy))
				t := d / maxRadius
				col := gradientPaletteColor(stops, t)
				if repeating {
					col = gradientPaletteColorRepeating(stops, t)
				}
				out.SetNRGBA(px, py, col)
			}
		}
		return out
	})
	if img == nil {
		cs.Pop()
		return
	}
	tr := op.Offset(image.Pt(x, y)).Push(ops)
	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{}.Add(ops)
	tr.Pop()
	cs.Pop()
}

func cssBorderRadiiValues(css map[string]string, w int, h int) cornerRadii {
	r := ParseBorderRadii(css, w, h)
	return cornerRadii{nw: r.NW, ne: r.NE, se: r.SE, sw: r.SW}
}

func cssBorderRadiusValue(css map[string]string, w int, h int) int {
	return BorderRadiusValue(css, w, h)
}

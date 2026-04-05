package render

import (
	"image"
	"image/color"
	"math"
	"strings"
	"sync"

	"gioui.org/f32"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	uicore "github.com/ArubikU/giocss/ui"
)

type BoxShadowLayer struct {
	OffsetX int
	OffsetY int
	Blur    int
	Spread  int
	Color   color.NRGBA
	Inset   bool
}

type boxShadowLayersKey struct {
	raw string
	w   int
	h   int
}

var (
	boxShadowLayersCacheMu sync.RWMutex
	boxShadowLayersCache   = map[boxShadowLayersKey][]BoxShadowLayer{}
)

func cloneBoxShadowLayers(in []BoxShadowLayer) []BoxShadowLayer {
	if len(in) == 0 {
		return nil
	}
	out := make([]BoxShadowLayer, len(in))
	copy(out, in)
	return out
}

func ParseBoxShadowLayersCached(raw string, basisW int, basisH int) []BoxShadowLayer {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "none") {
		return nil
	}
	key := boxShadowLayersKey{raw: raw, w: basisW, h: basisH}
	boxShadowLayersCacheMu.RLock()
	if cached, ok := boxShadowLayersCache[key]; ok {
		boxShadowLayersCacheMu.RUnlock()
		return cloneBoxShadowLayers(cached)
	}
	boxShadowLayersCacheMu.RUnlock()

	parsed := parseBoxShadowLayers(raw, basisW, basisH)
	boxShadowLayersCacheMu.Lock()
	if len(boxShadowLayersCache) > 384 {
		boxShadowLayersCache = map[boxShadowLayersKey][]BoxShadowLayer{}
	}
	boxShadowLayersCache[key] = cloneBoxShadowLayers(parsed)
	boxShadowLayersCacheMu.Unlock()
	return parsed
}

func parseBoxShadowLayers(raw string, basisW int, basisH int) []BoxShadowLayer {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "none") {
		return nil
	}
	items := uicore.SplitCommaOutsideParens(raw)
	if len(items) == 0 {
		return nil
	}
	out := make([]BoxShadowLayer, 0, len(items))
	for _, item := range items {
		cleanItem := item
		colorToken := ""
		for _, fn := range []string{"rgba(", "rgb(", "hsla(", "hsl(", "cmyk("} {
			if idx := strings.Index(strings.ToLower(cleanItem), fn); idx >= 0 {
				end := strings.Index(cleanItem[idx:], ")")
				if end > 0 {
					colorToken = strings.TrimSpace(cleanItem[idx : idx+end+1])
					cleanItem = strings.TrimSpace(cleanItem[:idx] + " " + cleanItem[idx+end+1:])
					break
				}
			}
		}
		tokens := strings.Fields(cleanItem)
		if len(tokens) < 2 {
			continue
		}
		layer := BoxShadowLayer{}
		vals := make([]int, 0, 4)
		for _, tok := range tokens {
			ltok := strings.ToLower(strings.TrimSpace(tok))
			if ltok == "inset" {
				layer.Inset = true
				continue
			}
			if strings.HasPrefix(ltok, "#") || strings.HasPrefix(ltok, "rgb(") || strings.HasPrefix(ltok, "rgba(") || strings.HasPrefix(ltok, "hsl(") || strings.HasPrefix(ltok, "hsla(") || strings.HasPrefix(ltok, "cmyk(") || ltok == "transparent" || ltok == "black" || ltok == "white" || ltok == "gray" || ltok == "grey" || ltok == "red" || ltok == "green" || ltok == "blue" || ltok == "yellow" || ltok == "orange" || ltok == "purple" {
				colorToken = tok
				continue
			}
			v := uicore.CSSLengthValue(tok, 0, max(basisW, basisH), basisW, basisH)
			vals = append(vals, v)
		}
		if len(vals) < 2 {
			continue
		}
		layer.OffsetX = vals[0]
		layer.OffsetY = vals[1]
		if len(vals) >= 3 {
			layer.Blur = max(0, vals[2])
		}
		if len(vals) >= 4 {
			layer.Spread = vals[3]
		}
		if colorToken == "" {
			colorToken = "rgba(0,0,0,0.35)"
		}
		layer.Color = toNRGBA(uicore.ParseHexColor(colorToken, color.NRGBA{R: 0, G: 0, B: 0, A: 90}))
		out = append(out, layer)
	}
	return out
}

type shadowRasterKey struct {
	w       int
	h       int
	radius  int
	offsetX int
	offsetY int
	blur    int
	spread  int
	color   uint32
	inset   bool
}

type ShadowRaster struct {
	Img *image.NRGBA
	DX  int
	DY  int
}

type shadowTemplateKey struct {
	radius  int
	offsetX int
	offsetY int
	blur    int
	spread  int
	color   uint32
}

type ShadowTemplate struct {
	Raster ShadowRaster
	CoreW  int
	CoreH  int
}

type ShadowSlice struct {
	Src image.Rectangle
	Dst image.Rectangle
}

type ShadowNineSlicePlan struct {
	Img    *image.NRGBA
	Slices [9]ShadowSlice
}

var (
	shadowRasterCacheMu sync.RWMutex
	shadowRasterCache   = map[shadowRasterKey]ShadowRaster{}
	shadowTemplateMu    sync.RWMutex
	shadowTemplateCache = map[shadowTemplateKey]ShadowTemplate{}
)

func shadowColorKey(c color.NRGBA) uint32 {
	return (uint32(c.R) << 24) | (uint32(c.G) << 16) | (uint32(c.B) << 8) | uint32(c.A)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func signedDistanceRoundedRect(px, py, left, top, width, height, radius float64) float64 {
	if width <= 0 || height <= 0 {
		return 1e9
	}
	hx := width / 2
	hy := height / 2
	r := radius
	if r < 0 {
		r = 0
	}
	if r > hx {
		r = hx
	}
	if r > hy {
		r = hy
	}
	cx := left + hx
	cy := top + hy
	qx := math.Abs(px-cx) - (hx - r)
	qy := math.Abs(py-cy) - (hy - r)
	ax := math.Max(qx, 0)
	ay := math.Max(qy, 0)
	outside := math.Hypot(ax, ay)
	inside := math.Min(math.Max(qx, qy), 0)
	return outside + inside - r
}

func shadowAlphaOuter(d float64, blur int, baseA float64) float64 {
	if blur <= 0 {
		if d <= 0 {
			return baseA
		}
		return 0
	}
	if d <= 0 {
		return baseA
	}
	sigma := math.Max(0.75, float64(blur)*0.6)
	a := baseA * math.Exp(-(d*d)/(2*sigma*sigma))
	if a < 0.5 {
		return 0
	}
	return a
}

func shadowAlphaInset(d float64, blur int, baseA float64) float64 {
	if blur <= 0 {
		if d >= 0 {
			return baseA
		}
		return 0
	}
	if d >= 0 {
		return baseA
	}
	dist := -d
	sigma := math.Max(0.75, float64(blur)*0.6)
	a := baseA * math.Exp(-(dist*dist)/(2*sigma*sigma))
	if a < 0.5 {
		return 0
	}
	return a
}

func buildOuterShadowRaster(layer BoxShadowLayer, w int, h int, radius int) (ShadowRaster, bool) {
	if w <= 0 || h <= 0 || layer.Color.A == 0 {
		return ShadowRaster{}, false
	}
	pad := max(2, layer.Blur*3+max(0, layer.Spread)+absInt(layer.OffsetX)+absInt(layer.OffsetY)+2)
	imgW := w + 2*pad
	imgH := h + 2*pad
	if imgW <= 0 || imgH <= 0 {
		return ShadowRaster{}, false
	}
	img := image.NewNRGBA(image.Rect(0, 0, imgW, imgH))
	rectLeft := float64(-layer.Spread)
	rectTop := float64(-layer.Spread)
	rectW := float64(w + 2*layer.Spread)
	rectH := float64(h + 2*layer.Spread)
	if rectW <= 0 || rectH <= 0 {
		return ShadowRaster{}, false
	}
	rad := float64(max(0, radius+layer.Spread))
	baseA := float64(layer.Color.A)
	for py := 0; py < imgH; py++ {
		ry := float64(py-pad) - float64(layer.OffsetY)
		for px := 0; px < imgW; px++ {
			rx := float64(px-pad) - float64(layer.OffsetX)
			d := signedDistanceRoundedRect(rx+0.5, ry+0.5, rectLeft, rectTop, rectW, rectH, rad)
			a := shadowAlphaOuter(d, layer.Blur, baseA)
			if a <= 0 {
				continue
			}
			img.SetNRGBA(px, py, color.NRGBA{R: layer.Color.R, G: layer.Color.G, B: layer.Color.B, A: uicore.ClampUint8(a)})
		}
	}
	return ShadowRaster{Img: img, DX: -pad, DY: -pad}, true
}

func buildInsetShadowRaster(layer BoxShadowLayer, w int, h int, radius int) (ShadowRaster, bool) {
	if w <= 0 || h <= 0 || layer.Color.A == 0 {
		return ShadowRaster{}, false
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	baseA := float64(layer.Color.A)
	for py := 0; py < h; py++ {
		y := float64(py) + 0.5
		for px := 0; px < w; px++ {
			x := float64(px) + 0.5
			d0 := signedDistanceRoundedRect(x, y, 0, 0, float64(w), float64(h), float64(radius))
			if d0 > 0 {
				continue
			}
			sx := x + float64(layer.OffsetX)
			sy := y + float64(layer.OffsetY)
			d := signedDistanceRoundedRect(sx, sy, 0, 0, float64(w), float64(h), float64(radius))
			d -= float64(layer.Spread)
			a := shadowAlphaInset(d, layer.Blur, baseA)
			if a <= 0 {
				continue
			}
			img.SetNRGBA(px, py, color.NRGBA{R: layer.Color.R, G: layer.Color.G, B: layer.Color.B, A: uicore.ClampUint8(a)})
		}
	}
	return ShadowRaster{Img: img, DX: 0, DY: 0}, true
}

func ShadowRasterForLayer(layer BoxShadowLayer, w int, h int, radius int) (ShadowRaster, bool) {
	key := shadowRasterKey{
		w:       w,
		h:       h,
		radius:  radius,
		offsetX: layer.OffsetX,
		offsetY: layer.OffsetY,
		blur:    layer.Blur,
		spread:  layer.Spread,
		color:   shadowColorKey(layer.Color),
		inset:   layer.Inset,
	}
	shadowRasterCacheMu.RLock()
	if cached, ok := shadowRasterCache[key]; ok {
		shadowRasterCacheMu.RUnlock()
		return cached, true
	}
	shadowRasterCacheMu.RUnlock()

	var raster ShadowRaster
	var ok bool
	if layer.Inset {
		raster, ok = buildInsetShadowRaster(layer, w, h, radius)
	} else {
		raster, ok = buildOuterShadowRaster(layer, w, h, radius)
	}
	if !ok {
		return ShadowRaster{}, false
	}

	shadowRasterCacheMu.Lock()
	if len(shadowRasterCache) > 384 {
		shadowRasterCache = map[shadowRasterKey]ShadowRaster{}
	}
	shadowRasterCache[key] = raster
	shadowRasterCacheMu.Unlock()
	return raster, true
}

func ShadowTemplateForLayer(layer BoxShadowLayer, radius int) (ShadowTemplate, bool) {
	if layer.Inset {
		return ShadowTemplate{}, false
	}
	key := shadowTemplateKey{
		radius:  radius,
		offsetX: layer.OffsetX,
		offsetY: layer.OffsetY,
		blur:    layer.Blur,
		spread:  layer.Spread,
		color:   shadowColorKey(layer.Color),
	}
	shadowTemplateMu.RLock()
	if cached, ok := shadowTemplateCache[key]; ok {
		shadowTemplateMu.RUnlock()
		return cached, true
	}
	shadowTemplateMu.RUnlock()

	coreW := max(4, 2*max(1, radius)+4)
	coreH := max(4, 2*max(1, radius)+4)
	raster, ok := ShadowRasterForLayer(layer, coreW, coreH, radius)
	if !ok || raster.Img == nil {
		return ShadowTemplate{}, false
	}
	tpl := ShadowTemplate{Raster: raster, CoreW: coreW, CoreH: coreH}

	shadowTemplateMu.Lock()
	if len(shadowTemplateCache) > 256 {
		shadowTemplateCache = map[shadowTemplateKey]ShadowTemplate{}
	}
	shadowTemplateCache[key] = tpl
	shadowTemplateMu.Unlock()
	return tpl, true
}

func BuildShadowTemplateNineSlicePlan(x int, y int, w int, h int, tpl ShadowTemplate) (ShadowNineSlicePlan, bool) {
	img := tpl.Raster.Img
	if img == nil || w <= 0 || h <= 0 {
		return ShadowNineSlicePlan{}, false
	}
	full := img.Bounds()
	left := max(0, -tpl.Raster.DX)
	top := max(0, -tpl.Raster.DY)
	right := min(full.Max.X, left+tpl.CoreW)
	bottom := min(full.Max.Y, top+tpl.CoreH)
	if left <= full.Min.X || top <= full.Min.Y || right >= full.Max.X || bottom >= full.Max.Y {
		return ShadowNineSlicePlan{}, false
	}

	outL := left - full.Min.X
	outT := top - full.Min.Y
	outR := full.Max.X - right
	outB := full.Max.Y - bottom

	x0 := x - outL
	x1 := x
	x2 := x + w
	x3 := x + w + outR
	y0 := y - outT
	y1 := y
	y2 := y + h
	y3 := y + h + outB

	if x1 < x0 || x2 < x1 || x3 < x2 || y1 < y0 || y2 < y1 || y3 < y2 {
		return ShadowNineSlicePlan{}, false
	}

	sx := [4]int{full.Min.X, left, right, full.Max.X}
	sy := [4]int{full.Min.Y, top, bottom, full.Max.Y}
	dx := [4]int{x0, x1, x2, x3}
	dy := [4]int{y0, y1, y2, y3}

	plan := ShadowNineSlicePlan{Img: img}
	idx := 0
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			plan.Slices[idx] = ShadowSlice{
				Src: image.Rect(sx[col], sy[row], sx[col+1], sy[row+1]),
				Dst: image.Rect(dx[col], dy[row], dx[col+1], dy[row+1]),
			}
			idx++
		}
	}
	return plan, true
}

func drawShadowRaster(ops *op.Ops, x int, y int, raster ShadowRaster) {
	if raster.Img == nil {
		return
	}
	tr := op.Offset(image.Pt(x+raster.DX, y+raster.DY)).Push(ops)
	imgOp := paint.NewImageOp(raster.Img)
	imgOp.Add(ops)
	paint.PaintOp{}.Add(ops)
	tr.Pop()
}

func drawShadowRasterSlice(ops *op.Ops, img *image.NRGBA, src image.Rectangle, dst image.Rectangle) {
	if img == nil || src.Empty() || dst.Empty() {
		return
	}
	if src.Dx() <= 0 || src.Dy() <= 0 || dst.Dx() <= 0 || dst.Dy() <= 0 {
		return
	}
	part := img.SubImage(src)
	if part == nil {
		return
	}
	clipStack := clip.Rect(dst).Push(ops)
	tr := op.Offset(image.Pt(dst.Min.X, dst.Min.Y)).Push(ops)
	sx := float32(dst.Dx()) / float32(src.Dx())
	sy := float32(dst.Dy()) / float32(src.Dy())
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(sx, sy))).Add(ops)
	paint.NewImageOp(part).Add(ops)
	paint.PaintOp{}.Add(ops)
	tr.Pop()
	clipStack.Pop()
}

func drawShadowTemplateNineSlice(ops *op.Ops, x int, y int, w int, h int, tpl ShadowTemplate) bool {
	plan, ok := BuildShadowTemplateNineSlicePlan(x, y, w, h, tpl)
	if !ok {
		return false
	}
	for _, s := range plan.Slices {
		drawShadowRasterSlice(ops, plan.Img, s.Src, s.Dst)
	}
	return true
}

func drawGioBoxShadow(ops *op.Ops, x int, y int, w int, h int, radius int, css map[string]string) {
	initRenderConfig()
	if renderLowPower {
		return
	}
	shadowRaw := strings.TrimSpace(css["box-shadow"])
	if shadowRaw == "" || strings.EqualFold(shadowRaw, "none") {
		return
	}
	layers := ParseBoxShadowLayersCached(shadowRaw, w, h)
	if len(layers) == 0 {
		return
	}
	for _, layer := range layers {
		if layer.Inset {
			continue
		}
		if tpl, ok := ShadowTemplateForLayer(layer, radius); ok {
			if drawShadowTemplateNineSlice(ops, x, y, w, h, tpl) {
				continue
			}
		}
		raster, ok := ShadowRasterForLayer(layer, w, h, radius)
		if !ok {
			continue
		}
		drawShadowRaster(ops, x, y, raster)
	}
}

func drawGioInsetBoxShadow(ops *op.Ops, x int, y int, w int, h int, radius int, css map[string]string) {
	initRenderConfig()
	if renderLowPower {
		return
	}
	shadowRaw := strings.TrimSpace(css["box-shadow"])
	if shadowRaw == "" || strings.EqualFold(shadowRaw, "none") {
		return
	}
	layers := ParseBoxShadowLayersCached(shadowRaw, w, h)
	if len(layers) == 0 {
		return
	}
	for _, layer := range layers {
		if !layer.Inset {
			continue
		}
		raster, ok := ShadowRasterForLayer(layer, w, h, radius)
		if !ok {
			continue
		}
		drawShadowRaster(ops, x, y, raster)
	}
}

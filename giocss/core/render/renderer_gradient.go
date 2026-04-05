package render

import (
	"image"
	"image/color"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"

	uicore "github.com/ArubikU/giocss/ui"
)

type BackgroundFilterFactors struct {
	Brightness float64
	Contrast   float64
	Saturate   float64
	Grayscale  float64
	Invert     float64
}

func ParseBackgroundFilterFactors(css map[string]string) BackgroundFilterFactors {
	ff := BackgroundFilterFactors{Brightness: 1.0, Contrast: 1.0, Saturate: 1.0, Grayscale: 0.0, Invert: 0.0}
	if rawBright := strings.TrimSpace(css["__brightness_factor"]); rawBright != "" {
		if f, err := strconv.ParseFloat(rawBright, 64); err == nil {
			ff.Brightness = f
		}
	}
	if rawContrast := strings.TrimSpace(css["__contrast_factor"]); rawContrast != "" {
		if f, err := strconv.ParseFloat(rawContrast, 64); err == nil {
			ff.Contrast = f
		}
	}
	if rawSaturate := strings.TrimSpace(css["__saturate_factor"]); rawSaturate != "" {
		if f, err := strconv.ParseFloat(rawSaturate, 64); err == nil {
			ff.Saturate = f
		}
	}
	if rawGrayscale := strings.TrimSpace(css["__grayscale_factor"]); rawGrayscale != "" {
		if f, err := strconv.ParseFloat(rawGrayscale, 64); err == nil {
			ff.Grayscale = f
		}
	}
	if rawInvert := strings.TrimSpace(css["__invert_factor"]); rawInvert != "" {
		if f, err := strconv.ParseFloat(rawInvert, 64); err == nil {
			ff.Invert = f
		}
	}
	return ff
}

var (
	backgroundLayersCacheMu sync.RWMutex
	backgroundLayersCache   = map[string][]string{}
)

func BackgroundLayers(raw string) []string {
	backgroundLayersCacheMu.RLock()
	if cached, ok := backgroundLayersCache[raw]; ok {
		backgroundLayersCacheMu.RUnlock()
		out := make([]string, len(cached))
		copy(out, cached)
		return out
	}
	backgroundLayersCacheMu.RUnlock()

	layers := uicore.SplitCommaOutsideParens(raw)
	if len(layers) == 0 {
		layers = []string{raw}
	}
	clean := make([]string, 0, len(layers))
	for _, l := range layers {
		lt := strings.TrimSpace(l)
		if lt != "" {
			clean = append(clean, lt)
		}
	}

	backgroundLayersCacheMu.Lock()
	if len(backgroundLayersCache) > 1024 {
		backgroundLayersCache = map[string][]string{}
	}
	copyClean := make([]string, len(clean))
	copy(copyClean, clean)
	backgroundLayersCache[raw] = copyClean
	backgroundLayersCacheMu.Unlock()

	out := make([]string, len(clean))
	copy(out, clean)
	return out
}

type GradientStop struct {
	Color color.NRGBA
	Pos   float64
}

type GradientStopPx struct {
	Color color.NRGBA
	Pos   float64
}

type GradientRasterKey struct {
	Kind     string
	W        int
	H        int
	NW       int
	NE       int
	SE       int
	SW       int
	Raw      string
	StopsSig string
}

var (
	gradientStopsCacheMu    sync.RWMutex
	gradientStopsCache      = map[string][]GradientStop{}
	repeatingPxStopsCacheMu sync.RWMutex
	repeatingPxStopsCache   = map[string]struct {
		stops []GradientStopPx
		ok    bool
	}{}
	gradientRasterCacheMu sync.RWMutex
	gradientRasterCache   = map[GradientRasterKey]*image.NRGBA{}
)

func cloneGradientStops(in []GradientStop) []GradientStop {
	if len(in) == 0 {
		return nil
	}
	out := make([]GradientStop, len(in))
	copy(out, in)
	return out
}

func cloneGradientStopsPx(in []GradientStopPx) []GradientStopPx {
	if len(in) == 0 {
		return nil
	}
	out := make([]GradientStopPx, len(in))
	copy(out, in)
	return out
}

func GradientStopsSignature(stops []GradientStop) string {
	if len(stops) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(stops) * 32)
	for _, s := range stops {
		b.WriteString(strconv.FormatFloat(s.Pos, 'f', 4, 64))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(int(s.Color.R)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.G)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.B)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.A)))
		b.WriteByte(';')
	}
	return b.String()
}

func GradientStopsPxSignature(stops []GradientStopPx) string {
	if len(stops) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(stops) * 32)
	for _, s := range stops {
		b.WriteString(strconv.FormatFloat(s.Pos, 'f', 2, 64))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(int(s.Color.R)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.G)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.B)))
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(int(s.Color.A)))
		b.WriteByte(';')
	}
	return b.String()
}

func GetCachedGradientRaster(key GradientRasterKey, build func() *image.NRGBA) *image.NRGBA {
	gradientRasterCacheMu.RLock()
	if img, ok := gradientRasterCache[key]; ok && img != nil {
		gradientRasterCacheMu.RUnlock()
		return img
	}
	gradientRasterCacheMu.RUnlock()

	img := build()
	if img == nil {
		return nil
	}

	gradientRasterCacheMu.Lock()
	if len(gradientRasterCache) > 512 {
		gradientRasterCache = map[GradientRasterKey]*image.NRGBA{}
	}
	gradientRasterCache[key] = img
	gradientRasterCacheMu.Unlock()
	return img
}

func GradientColorStops(raw string) []GradientStop {
	gradientStopsCacheMu.RLock()
	if cached, ok := gradientStopsCache[raw]; ok {
		gradientStopsCacheMu.RUnlock()
		return cloneGradientStops(cached)
	}
	gradientStopsCacheMu.RUnlock()

	start := strings.Index(raw, "(")
	end := strings.LastIndex(raw, ")")
	if start < 0 || end <= start {
		return nil
	}
	inner := raw[start+1 : end]
	parts := uicore.SplitCommaOutsideParens(inner)
	type rawStop struct {
		col    color.NRGBA
		pos    float64
		hasPos bool
	}
	st := make([]rawStop, 0, len(parts))
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		lp := strings.ToLower(p)
		if strings.Contains(lp, "deg") || strings.HasPrefix(lp, "to ") || strings.HasPrefix(lp, "circle") || strings.HasPrefix(lp, "ellipse") {
			continue
		}
		c := uicore.ParseHexColor(p, nil)
		hasPos := false
		pos := 0.0
		if c == nil {
			fields := strings.Fields(p)
			if len(fields) > 0 {
				c = uicore.ParseHexColor(fields[0], nil)
				if len(fields) > 1 {
					last := strings.TrimSpace(fields[len(fields)-1])
					if strings.HasSuffix(last, "%") {
						if f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(last, "%")), 64); err == nil {
							hasPos = true
							pos = f / 100.0
						}
					}
				}
			}
		} else {
			fields := strings.Fields(p)
			if len(fields) > 1 {
				last := strings.TrimSpace(fields[len(fields)-1])
				if strings.HasSuffix(last, "%") {
					if f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(last, "%")), 64); err == nil {
						hasPos = true
						pos = f / 100.0
					}
				}
			}
		}
		if c != nil {
			st = append(st, rawStop{col: toNRGBA(c), pos: pos, hasPos: hasPos})
		}
	}
	if len(st) == 0 {
		return nil
	}
	if !st[0].hasPos {
		st[0].pos = 0
		st[0].hasPos = true
	}
	if !st[len(st)-1].hasPos {
		st[len(st)-1].pos = 1
		st[len(st)-1].hasPos = true
	}
	for i := 1; i < len(st)-1; {
		if st[i].hasPos {
			i++
			continue
		}
		j := i
		for j < len(st)-1 && !st[j].hasPos {
			j++
		}
		left := st[i-1].pos
		right := st[j].pos
		count := j - i + 1
		for k := i; k < j; k++ {
			st[k].pos = left + (right-left)*float64(k-i+1)/float64(count)
			st[k].hasPos = true
		}
		i = j + 1
	}
	out := make([]GradientStop, 0, len(st))
	prev := -1.0
	for _, s := range st {
		p := s.pos
		if p < 0 {
			p = 0
		}
		if p > 1 {
			p = 1
		}
		if p < prev {
			p = prev
		}
		out = append(out, GradientStop{Color: s.col, Pos: p})
		prev = p
	}
	gradientStopsCacheMu.Lock()
	if len(gradientStopsCache) > 1024 {
		gradientStopsCache = map[string][]GradientStop{}
	}
	gradientStopsCache[raw] = cloneGradientStops(out)
	gradientStopsCacheMu.Unlock()
	return out
}

func GradientParts(raw string) []string {
	start := strings.Index(raw, "(")
	end := strings.LastIndex(raw, ")")
	if start < 0 || end <= start {
		return nil
	}
	inner := raw[start+1 : end]
	return uicore.SplitCommaOutsideParens(inner)
}

func ParseLinearGradientDirection(raw string) (float64, float64) {
	dx, dy := 0.0, 1.0
	parts := GradientParts(raw)
	if len(parts) == 0 {
		return dx, dy
	}
	tok := strings.ToLower(strings.TrimSpace(parts[0]))
	if tok == "" {
		return dx, dy
	}
	if c := uicore.ParseHexColor(tok, nil); c != nil {
		return dx, dy
	}
	if strings.HasPrefix(tok, "rgb(") || strings.HasPrefix(tok, "rgba(") || strings.HasPrefix(tok, "hsl(") || strings.HasPrefix(tok, "hsla(") || strings.HasPrefix(tok, "cmyk(") {
		return dx, dy
	}
	if strings.HasPrefix(tok, "to ") {
		dx, dy = 0, 0
		for _, f := range strings.Fields(strings.TrimPrefix(tok, "to ")) {
			switch f {
			case "left":
				dx -= 1
			case "right":
				dx += 1
			case "top":
				dy -= 1
			case "bottom":
				dy += 1
			}
		}
		if dx == 0 && dy == 0 {
			return 0, 1
		}
		m := math.Hypot(dx, dy)
		if m <= 0 {
			return 0, 1
		}
		return dx / m, dy / m
	}
	if strings.Contains(tok, "deg") || strings.Contains(tok, "grad") || strings.Contains(tok, "rad") || strings.Contains(tok, "turn") {
		if ang, ok := ParseCSSAngleToDegrees(tok); ok {
			rad := ang * math.Pi / 180.0
			dx = math.Sin(rad)
			dy = -math.Cos(rad)
			m := math.Hypot(dx, dy)
			if m > 0 {
				dx /= m
				dy /= m
			}
			return dx, dy
		}
	}
	return dx, dy
}

func ParseCSSAngleToDegrees(raw string) (float64, bool) {
	v := strings.ToLower(strings.TrimSpace(raw))
	if v == "" {
		return 0, false
	}
	if strings.HasSuffix(v, "deg") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "deg")), 64)
		return f, err == nil
	}
	if strings.HasSuffix(v, "grad") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "grad")), 64)
		if err != nil {
			return 0, false
		}
		return f * 0.9, true
	}
	if strings.HasSuffix(v, "turn") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "turn")), 64)
		if err != nil {
			return 0, false
		}
		return f * 360.0, true
	}
	if strings.HasSuffix(v, "rad") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "rad")), 64)
		if err != nil {
			return 0, false
		}
		return f * 180.0 / math.Pi, true
	}
	return 0, false
}

func ParseRadialGradientCenter(raw string, x, y, w, h int) (int, int) {
	cx, cy := x+w/2, y+h/2
	parts := GradientParts(raw)
	if len(parts) == 0 {
		return cx, cy
	}
	tok := strings.ToLower(strings.TrimSpace(parts[0]))
	idx := strings.Index(tok, " at ")
	if idx < 0 {
		return cx, cy
	}
	pos := strings.TrimSpace(tok[idx+4:])
	if pos == "" {
		return cx, cy
	}
	parseAxis := func(v string, basis int, fallback int) int {
		switch strings.TrimSpace(v) {
		case "left", "top":
			return 0
		case "right", "bottom":
			return basis
		case "center":
			return basis / 2
		default:
			return uicore.CSSLengthValue(v, fallback, basis, w, h)
		}
	}
	fields := strings.Fields(pos)
	if len(fields) == 1 {
		f := fields[0]
		switch f {
		case "top", "bottom":
			cy = y + parseAxis(f, h, h/2)
		default:
			cx = x + parseAxis(f, w, w/2)
		}
		return cx, cy
	}
	cx = x + parseAxis(fields[0], w, w/2)
	cy = y + parseAxis(fields[1], h, h/2)
	return cx, cy
}

func MixNRGBA(a color.NRGBA, b color.NRGBA, t float64) color.NRGBA {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	lerp := func(x, y float64) float64 {
		return x + (y-x)*t
	}
	clamp255 := func(v float64) uint8 {
		if v < 0 {
			v = 0
		}
		if v > 255 {
			v = 255
		}
		return uint8(v + 0.5)
	}

	aA := float64(a.A) / 255.0
	bA := float64(b.A) / 255.0
	outA := lerp(aA, bA)
	if outA <= 0 {
		return color.NRGBA{}
	}

	aR := (float64(a.R) / 255.0) * aA
	aG := (float64(a.G) / 255.0) * aA
	aB := (float64(a.B) / 255.0) * aA
	bR := (float64(b.R) / 255.0) * bA
	bG := (float64(b.G) / 255.0) * bA
	bB := (float64(b.B) / 255.0) * bA

	outR := lerp(aR, bR) / outA
	outG := lerp(aG, bG) / outA
	outB := lerp(aB, bB) / outA

	return color.NRGBA{
		R: clamp255(outR * 255.0),
		G: clamp255(outG * 255.0),
		B: clamp255(outB * 255.0),
		A: clamp255(outA * 255.0),
	}
}

func GradientPaletteColor(stops []GradientStop, t float64) color.NRGBA {
	if len(stops) == 0 {
		return color.NRGBA{}
	}
	if len(stops) == 1 {
		return stops[0].Color
	}
	if t <= stops[0].Pos {
		return stops[0].Color
	}
	if t >= stops[len(stops)-1].Pos {
		return stops[len(stops)-1].Color
	}
	for i := 0; i < len(stops)-1; i++ {
		l := stops[i]
		r := stops[i+1]
		if t < l.Pos || t > r.Pos {
			continue
		}
		if r.Pos <= l.Pos {
			return r.Color
		}
		local := (t - l.Pos) / (r.Pos - l.Pos)
		return MixNRGBA(l.Color, r.Color, local)
	}
	return stops[len(stops)-1].Color
}

func GradientPaletteColorRepeating(stops []GradientStop, t float64) color.NRGBA {
	if len(stops) == 0 {
		return color.NRGBA{}
	}
	if len(stops) == 1 {
		return stops[0].Color
	}
	start := stops[0].Pos
	end := stops[len(stops)-1].Pos
	span := end - start
	if span <= 0 {
		span = 1
		start = 0
	}
	wrapped := start + math.Mod(t-start, span)
	if wrapped < start {
		wrapped += span
	}
	return GradientPaletteColor(stops, wrapped)
}

func ParseRepeatingLinearPxStops(raw string) ([]GradientStopPx, bool) {
	repeatingPxStopsCacheMu.RLock()
	if cached, ok := repeatingPxStopsCache[raw]; ok {
		repeatingPxStopsCacheMu.RUnlock()
		return cloneGradientStopsPx(cached.stops), cached.ok
	}
	repeatingPxStopsCacheMu.RUnlock()

	parts := GradientParts(raw)
	if len(parts) == 0 {
		repeatingPxStopsCacheMu.Lock()
		repeatingPxStopsCache[raw] = struct {
			stops []GradientStopPx
			ok    bool
		}{stops: nil, ok: false}
		repeatingPxStopsCacheMu.Unlock()
		return nil, false
	}
	out := make([]GradientStopPx, 0, len(parts))
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" {
			continue
		}
		lp := strings.ToLower(p)
		if strings.Contains(lp, "deg") || strings.HasPrefix(lp, "to ") {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) < 2 {
			continue
		}
		c := uicore.ParseHexColor(fields[0], nil)
		if c == nil {
			continue
		}
		last := strings.ToLower(strings.TrimSpace(fields[len(fields)-1]))
		if !strings.HasSuffix(last, "px") {
			continue
		}
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(last, "px")), 64)
		if err != nil {
			continue
		}
		out = append(out, GradientStopPx{Color: toNRGBA(c), Pos: f})
	}
	if len(out) < 2 {
		repeatingPxStopsCacheMu.Lock()
		repeatingPxStopsCache[raw] = struct {
			stops []GradientStopPx
			ok    bool
		}{stops: nil, ok: false}
		repeatingPxStopsCacheMu.Unlock()
		return nil, false
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Pos < out[j].Pos })
	repeatingPxStopsCacheMu.Lock()
	if len(repeatingPxStopsCache) > 1024 {
		repeatingPxStopsCache = map[string]struct {
			stops []GradientStopPx
			ok    bool
		}{}
	}
	repeatingPxStopsCache[raw] = struct {
		stops []GradientStopPx
		ok    bool
	}{stops: cloneGradientStopsPx(out), ok: true}
	repeatingPxStopsCacheMu.Unlock()
	return out, true
}

func GradientPaletteColorPx(stops []GradientStopPx, t float64) color.NRGBA {
	if len(stops) == 0 {
		return color.NRGBA{}
	}
	if len(stops) == 1 {
		return stops[0].Color
	}
	if t <= stops[0].Pos {
		return stops[0].Color
	}
	if t >= stops[len(stops)-1].Pos {
		return stops[len(stops)-1].Color
	}
	for i := 0; i < len(stops)-1; i++ {
		l := stops[i]
		r := stops[i+1]
		if t < l.Pos || t > r.Pos {
			continue
		}
		if r.Pos <= l.Pos {
			return r.Color
		}
		local := (t - l.Pos) / (r.Pos - l.Pos)
		return MixNRGBA(l.Color, r.Color, local)
	}
	return stops[len(stops)-1].Color
}

package layout

import (
	"image"
	"math"
	"strconv"
	"strings"
)

type PositionInfo struct {
	Position  string
	OffsetX   float64
	OffsetY   float64
	IsLeft    bool
	IsTop     bool
	HasLeft   bool
	HasRight  bool
	HasTop    bool
	HasBottom bool
}

type PositionLayoutResult struct {
	X                  int
	Y                  int
	NodeRect           image.Rectangle
	ChildContainingBox image.Rectangle
}

func ApplyPositionLayout(x, y, w, h int, pos PositionInfo, containingBox image.Rectangle, viewport image.Rectangle) PositionLayoutResult {
	outX := x
	outY := y
	switch pos.Position {
	case "relative", "sticky":
		if pos.HasLeft {
			outX += int(pos.OffsetX)
		} else if pos.HasRight {
			outX -= int(pos.OffsetX)
		}
		if pos.HasTop {
			outY += int(pos.OffsetY)
		} else if pos.HasBottom {
			outY -= int(pos.OffsetY)
		}
	case "absolute":
		if pos.HasLeft {
			outX = containingBox.Min.X + int(pos.OffsetX)
		} else if pos.HasRight {
			outX = containingBox.Max.X - w - int(pos.OffsetX)
		}
		if pos.HasTop {
			outY = containingBox.Min.Y + int(pos.OffsetY)
		} else if pos.HasBottom {
			outY = containingBox.Max.Y - h - int(pos.OffsetY)
		}
	case "fixed":
		if viewport.Empty() {
			viewport = containingBox
		}
		if pos.HasLeft {
			outX = viewport.Min.X + int(pos.OffsetX)
		} else if pos.HasRight {
			outX = viewport.Max.X - w - int(pos.OffsetX)
		}
		if pos.HasTop {
			outY = viewport.Min.Y + int(pos.OffsetY)
		} else if pos.HasBottom {
			outY = viewport.Max.Y - h - int(pos.OffsetY)
		}
	}

	nodeRect := image.Rect(outX, outY, outX+w, outY+h)
	childBox := containingBox
	if pos.Position != "static" {
		childBox = nodeRect
	}
	return PositionLayoutResult{X: outX, Y: outY, NodeRect: nodeRect, ChildContainingBox: childBox}
}

type OverflowLayoutResult struct {
	OverflowX    string
	OverflowY    string
	ClipChildren bool
}

type RenderVisibilityResult struct {
	ViewportRect image.Rectangle
	RenderRect   image.Rectangle
	IsOnScreen   bool
}

type ChildTraversalPlan struct {
	ChildShiftX          int
	ChildShiftY          int
	EnableClip           bool
	ApplyScrollTransform bool
}

func ShouldQuickCullNode(path string, x, y, w, h, inheritedShiftX, inheritedShiftY int, frameViewport image.Rectangle, viewW, viewH, childCount int, kind string) bool {
	quickRect := image.Rect(x+inheritedShiftX, y+inheritedShiftY, x+w+inheritedShiftX, y+h+inheritedShiftY)
	quickViewport := frameViewport
	if quickViewport.Empty() {
		quickViewport = ExpandedViewportRect(viewW, viewH, 96)
	}
	if path != "root" && !quickRect.Overlaps(quickViewport) && childCount == 0 && kind != "input" && kind != "button" {
		return true
	}
	return false
}

func ShouldProfileNode(profileComponents, profileFull bool, path string, profileSampleFrame bool) bool {
	if !profileComponents {
		return false
	}
	return profileFull || (path != "root" && profileSampleFrame)
}

func ResolveOverflowLayout(css map[string]string) OverflowLayoutResult {
	overflowMode := LowerASCIIIfNeeded(strings.TrimSpace(css["overflow"]))
	overflowX := LowerASCIIIfNeeded(strings.TrimSpace(css["overflow-x"]))
	overflowY := LowerASCIIIfNeeded(strings.TrimSpace(css["overflow-y"]))
	if overflowX == "" {
		overflowX = overflowMode
	}
	if overflowY == "" {
		overflowY = overflowMode
	}
	clipX := overflowX == "hidden" || overflowX == "auto" || overflowX == "scroll" || overflowX == "clip"
	clipY := overflowY == "hidden" || overflowY == "auto" || overflowY == "scroll" || overflowY == "clip"
	return OverflowLayoutResult{OverflowX: overflowX, OverflowY: overflowY, ClipChildren: clipX || clipY}
}

func ComputeRenderVisibility(path string, x, y, w, h int, inheritedShiftX, inheritedShiftY int, translateX, translateY, scale float64, frameViewport image.Rectangle, viewW, viewH int) RenderVisibilityResult {
	viewport := frameViewport
	if viewport.Empty() {
		viewport = ExpandedViewportRect(viewW, viewH, 48)
	}
	renderRect := NodeRenderRect(x+inheritedShiftX, y+inheritedShiftY, w, h, translateX, translateY, scale)
	isOnScreen := renderRect.Overlaps(viewport) || path == "root"
	return RenderVisibilityResult{
		ViewportRect: viewport,
		RenderRect:   renderRect,
		IsOnScreen:   isOnScreen,
	}
}

func ShouldRegisterContainerEvents(kind string, hasHover, hasActive bool, overflowX, overflowY string) bool {
	if kind == "button" || kind == "input" {
		return false
	}
	if hasHover || hasActive {
		return true
	}
	return overflowX == "auto" || overflowX == "scroll" || overflowY == "auto" || overflowY == "scroll"
}

func ShouldApplyNodeCursor(cursorRaw string) bool {
	v := strings.TrimSpace(cursorRaw)
	if v == "" {
		return false
	}
	return !strings.EqualFold(v, "default") && !strings.EqualFold(v, "auto")
}

func ShouldPromoteFrameCursor(path, currentPath string, pointerKnown bool, pointerPos image.Point, renderRect image.Rectangle) bool {
	if !pointerKnown {
		return false
	}
	if !PointInRect(pointerPos, renderRect) {
		return false
	}
	return len(path) >= len(currentPath)
}

func BuildChildTraversalPlan(visible, isOnScreen, clipChildren bool, inheritedShiftX, inheritedShiftY int, translateX, translateY float64, scroll image.Point) ChildTraversalPlan {
	plan := ChildTraversalPlan{
		ChildShiftX: inheritedShiftX + int(translateX),
		ChildShiftY: inheritedShiftY + int(translateY),
	}
	if !(visible && isOnScreen && clipChildren) {
		return plan
	}
	plan.EnableClip = true
	if scroll.X != 0 || scroll.Y != 0 {
		plan.ApplyScrollTransform = true
		plan.ChildShiftX -= scroll.X
		plan.ChildShiftY -= scroll.Y
	}
	return plan
}

func IsDisplayNone(displayValue string) bool {
	v := strings.TrimSpace(displayValue)
	return strings.EqualFold(v, "none")
}

func IsVisibilityHidden(visibilityValue string) bool {
	v := strings.TrimSpace(visibilityValue)
	return strings.EqualFold(v, "hidden")
}

func IsFilterEnabled(filterRaw string) bool {
	v := strings.TrimSpace(filterRaw)
	if v == "" {
		return false
	}
	return !strings.EqualFold(v, "none")
}

type CSSFilterParams struct {
	Grayscale    float64
	Invert       float64
	Sepia        float64
	Saturate     float64
	Brightness   float64
	Contrast     float64
	Opacity      float64
	HueRotateDeg float64
	BlurPx       int
}

func BuildFilterFactorVars(filterRaw string) map[string]string {
	params := ParseCSSFilterChain(filterRaw)
	vars := make(map[string]string, 5)
	if math.Abs(params.Brightness-1.0) > 0.001 {
		vars["__brightness_factor"] = strconv.FormatFloat(params.Brightness, 'f', 3, 64)
	}
	if math.Abs(params.Contrast-1.0) > 0.001 {
		vars["__contrast_factor"] = strconv.FormatFloat(params.Contrast, 'f', 3, 64)
	}
	if math.Abs(params.Saturate-1.0) > 0.001 {
		vars["__saturate_factor"] = strconv.FormatFloat(params.Saturate, 'f', 3, 64)
	}
	if params.Grayscale > 0.001 {
		vars["__grayscale_factor"] = strconv.FormatFloat(params.Grayscale, 'f', 3, 64)
	}
	if params.Invert > 0.001 {
		vars["__invert_factor"] = strconv.FormatFloat(params.Invert, 'f', 3, 64)
	}
	if len(vars) == 0 {
		return nil
	}
	return vars
}

func ParseCSSFilterChain(filter string) CSSFilterParams {
	p := CSSFilterParams{Saturate: 1.0, Brightness: 1.0, Contrast: 1.0, Opacity: 1.0}
	chain := strings.TrimSpace(filter)
	if chain == "" {
		return p
	}
	parts := splitFilterFunctions(chain)
	for _, fn := range parts {
		lower := strings.ToLower(strings.TrimSpace(fn))
		s := strings.Index(lower, "(")
		e := strings.LastIndex(lower, ")")
		if s < 0 || e <= s {
			continue
		}
		name := strings.TrimSpace(lower[:s])
		arg := strings.TrimSpace(lower[s+1 : e])
		switch name {
		case "grayscale":
			p.Grayscale = clamp01(parseFilterAmount(arg, 1.0))
		case "invert":
			p.Invert = clamp01(parseFilterAmount(arg, 1.0))
		case "sepia":
			p.Sepia = clamp01(parseFilterAmount(arg, 1.0))
		case "saturate":
			p.Saturate = maxFloat(0.0, parseFilterAmount(arg, 1.0))
		case "brightness":
			p.Brightness = maxFloat(0.0, parseFilterAmount(arg, 1.0))
		case "contrast":
			p.Contrast = maxFloat(0.0, parseFilterAmount(arg, 1.0))
		case "opacity":
			p.Opacity = clamp01(parseFilterAmount(arg, 1.0))
		case "hue-rotate":
			p.HueRotateDeg = parseAngleToDeg(arg)
		case "blur":
			px := parseLengthPx(arg)
			if px > 0 {
				p.BlurPx = px
			}
		}
	}
	return p
}

func splitFilterFunctions(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	out := make([]string, 0, 8)
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
		case ' ', '\t', '\n':
			if depth == 0 {
				part := strings.TrimSpace(input[start:i])
				if part != "" {
					out = append(out, part)
				}
				start = i + 1
			}
		}
	}
	last := strings.TrimSpace(input[start:])
	if last != "" {
		out = append(out, last)
	}
	return out
}

func parseFilterAmount(raw string, fallback float64) float64 {
	v := strings.TrimSpace(raw)
	if v == "" {
		return fallback
	}
	if strings.HasSuffix(v, "%") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "%")), 64)
		if err != nil {
			return fallback
		}
		return f / 100.0
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	return f
}

func parseLengthPx(raw string) int {
	v := strings.TrimSpace(raw)
	if v == "" {
		return 0
	}
	if strings.HasSuffix(v, "px") {
		f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "px")), 64)
		if err == nil {
			return max(0, int(math.Round(f)))
		}
		return 0
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return max(0, int(math.Round(f)))
}

func parseAngleToDeg(raw string) float64 {
	v := strings.TrimSpace(strings.ToLower(raw))
	if v == "" {
		return 0
	}
	if strings.HasSuffix(v, "deg") {
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "deg")), 64)
		return f
	}
	if strings.HasSuffix(v, "grad") {
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "grad")), 64)
		return f * 0.9
	}
	if strings.HasSuffix(v, "turn") {
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "turn")), 64)
		return f * 360.0
	}
	if strings.HasSuffix(v, "rad") {
		f, _ := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, "rad")), 64)
		return f * 180.0 / math.Pi
	}
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

func LowerASCIIIfNeeded(s string) string {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return strings.ToLower(s)
		}
	}
	return s
}

func ExpandedViewportRect(viewW, viewH, margin int) image.Rectangle {
	return image.Rect(-margin, -margin, viewW+margin, viewH+margin)
}

func NodeRenderRect(x, y, w, h int, tx, ty float64, scale float64) image.Rectangle {
	r := image.Rect(x+int(tx), y+int(ty), x+w+int(tx), y+h+int(ty))
	if scale <= 0 || math.Abs(scale-1.0) < 0.001 {
		return r
	}
	cx := float64(r.Min.X+r.Max.X) / 2.0
	cy := float64(r.Min.Y+r.Max.Y) / 2.0
	hw := float64(r.Dx()) * scale / 2.0
	hh := float64(r.Dy()) * scale / 2.0
	return image.Rect(int(cx-hw), int(cy-hh), int(cx+hw), int(cy+hh))
}

func PointInRect(p image.Point, r image.Rectangle) bool {
	return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

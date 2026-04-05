package input

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gioui.org/io/pointer"
)

func InputValueString(candidate any) string {
	switch v := candidate.(type) {
	case string:
		return strings.ToValidUTF8(v, "")
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return ""
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		fv := float64(v)
		if math.IsNaN(fv) || math.IsInf(fv, 0) {
			return ""
		}
		return strconv.FormatFloat(fv, 'f', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func InputExternalValue(props map[string]any) string {
	if props == nil {
		return ""
	}
	if raw, ok := props["text"]; ok {
		if s := strings.TrimSpace(InputValueString(raw)); s != "" {
			return s
		}
	}
	if raw, ok := props["value"]; ok {
		return strings.TrimSpace(InputValueString(raw))
	}
	return ""
}

func InputPropFloat(props map[string]any, key string) (float64, bool) {
	if props == nil {
		return 0, false
	}
	raw, ok := props[key]
	if !ok {
		return 0, false
	}
	switch v := raw.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, false
		}
		return v, true
	case float32:
		fv := float64(v)
		if math.IsNaN(fv) || math.IsInf(fv, 0) {
			return 0, false
		}
		return fv, true
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return 0, false
		}
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func NumberStepFromProps(props map[string]any) (float64, bool) {
	if props == nil {
		return 0, false
	}
	raw, ok := props["step"]
	if !ok {
		return 0, false
	}
	if s, ok := raw.(string); ok {
		s = strings.TrimSpace(strings.ToLower(s))
		if s == "" || s == "any" {
			return 0, false
		}
	}
	step, has := InputPropFloat(props, "step")
	if !has || step <= 0 {
		return 0, false
	}
	return step, true
}

func SanitizeNumberLive(raw string) string {
	raw = strings.TrimSpace(strings.ToValidUTF8(raw, ""))
	if raw == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(raw))
	hasSign := false
	hasDot := false
	for i, r := range raw {
		switch {
		case unicode.IsDigit(r):
			b.WriteRune(r)
		case (r == '+' || r == '-') && i == 0 && !hasSign:
			hasSign = true
			b.WriteRune(r)
		case r == '.' && !hasDot:
			hasDot = true
			b.WriteRune(r)
		}
	}
	return b.String()
}

func NormalizeNumberInput(raw string, props map[string]any, finalize bool) string {
	clean := SanitizeNumberLive(raw)
	if clean == "" {
		return ""
	}
	if !finalize {
		switch clean {
		case "+", "-", ".", "+.", "-.":
			return clean
		}
		return clean
	}
	v, err := strconv.ParseFloat(clean, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return ""
	}
	if minV, ok := InputPropFloat(props, "min"); ok && v < minV {
		v = minV
	}
	if maxV, ok := InputPropFloat(props, "max"); ok && v > maxV {
		v = maxV
	}
	if step, ok := NumberStepFromProps(props); ok {
		base := 0.0
		if minV, hasMin := InputPropFloat(props, "min"); hasMin {
			base = minV
		}
		steps := math.Round((v - base) / step)
		v = base + steps*step
		if minV, ok := InputPropFloat(props, "min"); ok && v < minV {
			v = minV
		}
		if maxV, ok := InputPropFloat(props, "max"); ok && v > maxV {
			v = maxV
		}
	}
	if math.Abs(v) < 1e-12 {
		v = 0
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func SanitizeDateLive(raw string) string {
	raw = strings.TrimSpace(strings.ToValidUTF8(raw, ""))
	if raw == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(min(len(raw), 10))
	for _, r := range raw {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		} else if r == '/' || r == '-' {
			b.WriteRune('-')
		}
		if b.Len() >= 10 {
			break
		}
	}
	return b.String()
}

func NormalizeDateInput(raw string, finalize bool) string {
	clean := SanitizeDateLive(raw)
	if clean == "" {
		return ""
	}
	if !finalize {
		return clean
	}
	for _, layout := range []string{"2006-01-02", "2006/01/02", "02-01-2006", "02/01/2006"} {
		if parsed, err := time.Parse(layout, clean); err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return ""
}

func SanitizeTimeLive(raw string) string {
	raw = strings.TrimSpace(strings.ToValidUTF8(raw, ""))
	if raw == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(min(len(raw), 8))
	for _, r := range raw {
		if unicode.IsDigit(r) || r == ':' {
			b.WriteRune(r)
		}
		if b.Len() >= 8 {
			break
		}
	}
	return b.String()
}

func NormalizeTimeInput(raw string, props map[string]any, finalize bool) string {
	clean := SanitizeTimeLive(raw)
	if clean == "" {
		return ""
	}
	if !finalize {
		return clean
	}
	step, hasStep := NumberStepFromProps(props)
	wantSeconds := strings.Count(clean, ":") == 2 || (hasStep && step < 60)
	for _, layout := range []string{"15:04", "15:04:05"} {
		if parsed, err := time.Parse(layout, clean); err == nil {
			if wantSeconds {
				return parsed.Format("15:04:05")
			}
			return parsed.Format("15:04")
		}
	}
	return ""
}

func NormalizeTypedInputValue(inputType, raw string, props map[string]any, finalize bool) string {
	switch strings.ToLower(strings.TrimSpace(inputType)) {
	case "number":
		return NormalizeNumberInput(raw, props, finalize)
	case "date":
		return NormalizeDateInput(raw, finalize)
	case "time":
		return NormalizeTimeInput(raw, props, finalize)
	default:
		return strings.TrimSpace(strings.ToValidUTF8(raw, ""))
	}
}

func StepNumberInputValue(current string, props map[string]any, delta int) string {
	if delta == 0 {
		return NormalizeNumberInput(current, props, true)
	}
	step := 1.0
	if customStep, ok := NumberStepFromProps(props); ok {
		step = customStep
	}
	cur := 0.0
	normalizedCurrent := NormalizeNumberInput(current, props, true)
	if normalizedCurrent != "" {
		if parsed, err := strconv.ParseFloat(normalizedCurrent, 64); err == nil && !math.IsNaN(parsed) && !math.IsInf(parsed, 0) {
			cur = parsed
		}
	} else if minV, ok := InputPropFloat(props, "min"); ok {
		cur = minV
	}
	cur += float64(delta) * step
	return NormalizeNumberInput(strconv.FormatFloat(cur, 'f', -1, 64), props, true)
}

func CSSCursorToGio(raw string) pointer.Cursor {
	v := strings.ToLower(strings.TrimSpace(raw))
	switch v {
	case "", "auto", "default", "initial", "inherit", "unset":
		return pointer.CursorDefault
	case "none":
		return pointer.CursorNone
	case "pointer":
		return pointer.CursorPointer
	case "text":
		return pointer.CursorText
	case "vertical-text":
		return pointer.CursorVerticalText
	case "crosshair":
		return pointer.CursorCrosshair
	case "move", "all-scroll":
		return pointer.CursorAllScroll
	case "grab":
		return pointer.CursorGrab
	case "grabbing":
		return pointer.CursorGrabbing
	case "not-allowed", "no-drop":
		return pointer.CursorNotAllowed
	case "wait":
		return pointer.CursorWait
	case "progress":
		return pointer.CursorProgress
	case "col-resize":
		return pointer.CursorColResize
	case "row-resize":
		return pointer.CursorRowResize
	case "n-resize":
		return pointer.CursorNorthResize
	case "s-resize":
		return pointer.CursorSouthResize
	case "e-resize":
		return pointer.CursorEastResize
	case "w-resize":
		return pointer.CursorWestResize
	case "ns-resize":
		return pointer.CursorNorthSouthResize
	case "ew-resize":
		return pointer.CursorEastWestResize
	case "ne-resize":
		return pointer.CursorNorthEastResize
	case "nw-resize":
		return pointer.CursorNorthWestResize
	case "se-resize":
		return pointer.CursorSouthEastResize
	case "sw-resize":
		return pointer.CursorSouthWestResize
	case "nesw-resize":
		return pointer.CursorNorthEastSouthWestResize
	case "nwse-resize":
		return pointer.CursorNorthWestSouthEastResize
	default:
		return pointer.CursorDefault
	}
}



// ResolveSliderValueFromPointer resolves a slider value from the pointer X position.
func ResolveSliderValueFromPointer(pointerX, boundsMinX, boundsWidth int, minV, maxV float64) (float64, bool) {
if boundsWidth <= 0 {
return 0, false
}
if maxV <= minV {
maxV = minV + 100
}
pct := float64(pointerX-boundsMinX) / float64(boundsWidth)
if pct < 0 {
pct = 0
}
if pct > 1 {
pct = 1
}
return minV + pct*(maxV-minV), true
}

// NextSelectIndex advances select index by 1 wrapping around.
func NextSelectIndex(current, total int) int {
if total <= 0 {
return 0
}
if current < 0 {
current = 0
}
return (current + 1) % total
}


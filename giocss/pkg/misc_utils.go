package pkg

import (
	"strconv"
	"strings"
	"time"

	"gioui.org/text"
)

// CSSGridSpan parses a CSS grid-column/row span value and returns the integer span count.
func CSSGridSpan(value string) int {
	v := strings.ToLower(strings.TrimSpace(value))
	if v == "" {
		return 1
	}
	if strings.HasPrefix(v, "span") {
		parts := strings.Fields(v)
		if len(parts) >= 2 {
			if n, err := strconv.Atoi(parts[1]); err == nil && n > 0 {
				return n
			}
		}
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return 1
}

// ShouldClipText returns true when CSS properties indicate text should be clipped.
func ShouldClipText(css map[string]string) bool {
	overflow := strings.ToLower(strings.TrimSpace(css["overflow"]))
	textOverflow := strings.ToLower(strings.TrimSpace(css["text-overflow"]))
	whiteSpace := strings.ToLower(strings.TrimSpace(css["white-space"]))
	if textOverflow == "ellipsis" {
		return true
	}
	if overflow == "hidden" && whiteSpace == "nowrap" {
		return true
	}
	return false
}

// CSSTextAlign converts the CSS text-align property to a Gio text.Alignment value.
func CSSTextAlign(css map[string]string) text.Alignment {
	switch strings.ToLower(strings.TrimSpace(css["text-align"])) {
	case "center":
		return text.Middle
	case "right", "end":
		return text.End
	default:
		return text.Start
	}
}

// MathAbs returns the absolute value of v.
func MathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// AnimateFrames runs frame callbacks from t=0..1 over the given duration in a goroutine.
func AnimateFrames(duration time.Duration, frame func(t float64)) {
	go func() {
		steps := 10
		if duration < 100*time.Millisecond {
			steps = 6
		}
		for i := 0; i <= steps; i++ {
			frame(float64(i) / float64(steps))
			time.Sleep(duration / time.Duration(steps))
		}
	}()
}

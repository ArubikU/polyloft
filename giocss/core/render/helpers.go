package render

import (
	"fmt"
	"image/color"
)

func anyToString(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		s := fmt.Sprintf("%v", v)
		if s == "<nil>" {
			return fallback
		}
		return s
	}
}

func toNRGBA(c color.Color) color.NRGBA {
	if c == nil {
		return color.NRGBA{}
	}
	r, g, b, a := c.RGBA()
	if a == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8(r * 255 / a),
		G: uint8(g * 255 / a),
		B: uint8(b * 255 / a),
		A: uint8(a >> 8),
	}
}

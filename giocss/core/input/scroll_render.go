package input

import (
	"image"
	"image/color"
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

type ScrollbarRenderData struct {
	Hit       ScrollHitInfo
	Thickness int
	Radius    int
	ThumbCol  color.NRGBA
	TrackCol  color.NRGBA
}

func ScrollbarThickness(css map[string]string, w int, h int) int {
	v := strings.ToLower(strings.TrimSpace(css["scrollbar-width"]))
	switch v {
	case "none":
		return 0
	case "thin":
		return 8
	case "auto", "":
		return 12
	default:
		return max(4, uicore.CSSLengthValue(v, 12, max(w, h), w, h))
	}
}

func ScrollbarColors(css map[string]string) (color.NRGBA, color.NRGBA) {
	thumb := color.NRGBA{R: 0x8E, G: 0x9C, B: 0xC9, A: 0xF0}
	track := color.NRGBA{R: 0x2B, G: 0x33, B: 0x47, A: 0xD8}
	if raw := strings.TrimSpace(css["scrollbar-color"]); raw != "" {
		parts := strings.Fields(raw)
		if len(parts) >= 1 {
			thumb = toNRGBA(uicore.ParseHexColor(parts[0], thumb))
		}
		if len(parts) >= 2 {
			track = toNRGBA(uicore.ParseHexColor(parts[1], track))
		}
	}
	if raw := strings.TrimSpace(css["scrollbar-thumb-color"]); raw != "" {
		thumb = toNRGBA(uicore.ParseHexColor(raw, thumb))
	}
	if raw := strings.TrimSpace(css["scrollbar-track-color"]); raw != "" {
		track = toNRGBA(uicore.ParseHexColor(raw, track))
	}
	return thumb, track
}

func ComputeScrollbarRenderData(x int, y int, w int, h int, contentRight int, contentBottom int, css map[string]string, scrollX int, scrollY int) ScrollbarRenderData {
	data := ScrollbarRenderData{Hit: ScrollHitInfo{Rect: image.Rect(x, y, x+w, y+h)}}
	overflow := strings.ToLower(strings.TrimSpace(css["overflow"]))
	overflowX := strings.ToLower(strings.TrimSpace(css["overflow-x"]))
	overflowY := strings.ToLower(strings.TrimSpace(css["overflow-y"]))
	if overflowX == "" {
		overflowX = overflow
	}
	if overflowY == "" {
		overflowY = overflow
	}
	if (overflowX != "auto" && overflowX != "scroll") && (overflowY != "auto" && overflowY != "scroll") {
		return data
	}
	th := ScrollbarThickness(css, w, h)
	if th <= 0 {
		return data
	}
	data.Thickness = th
	thumbCol, trackCol := ScrollbarColors(css)
	data.ThumbCol = thumbCol
	data.TrackCol = trackCol
	sr := max(2, uicore.CSSLengthValue(css["scrollbar-radius"], th/2, max(w, h), w, h))
	if sr > th {
		sr = th
	}
	data.Radius = sr

	contentW := max(w, contentRight-x)
	contentH := max(h, contentBottom-y)

	showV := overflowY == "scroll" || (overflowY == "auto" && contentH > h)
	showH := overflowX == "scroll" || (overflowX == "auto" && contentW > w)
	data.Hit.HasV = showV
	data.Hit.HasH = showH
	data.Hit.MaxX = max(0, contentW-w)
	data.Hit.MaxY = max(0, contentH-h)

	if showV {
		trackX := x + w - th - 1
		trackY := y + 1
		trackH := h - 2
		if showH {
			trackH -= th
		}
		if trackH > 0 {
			data.Hit.TrackV = image.Rect(trackX, trackY, trackX+th, trackY+trackH)
			ratio := float64(h) / float64(max(h, contentH))
			thumbH := max(th, int(float64(trackH)*ratio))
			thumbY := trackY
			if contentH > h {
				thumbY = trackY + int(float64(trackH-thumbH)*float64(scrollY)/float64(max(1, contentH-h)))
			}
			data.Hit.ThumbV = image.Rect(trackX, thumbY, trackX+th, thumbY+thumbH)
		}
	}
	if showH {
		trackX := x + 1
		trackY := y + h - th - 1
		trackW := w - 2
		if showV {
			trackW -= th
		}
		if trackW > 0 {
			data.Hit.TrackH = image.Rect(trackX, trackY, trackX+trackW, trackY+th)
			ratio := float64(w) / float64(max(w, contentW))
			thumbW := max(th, int(float64(trackW)*ratio))
			thumbX := trackX
			if contentW > w {
				thumbX = trackX + int(float64(trackW-thumbW)*float64(scrollX)/float64(max(1, contentW-w)))
			}
			data.Hit.ThumbH = image.Rect(thumbX, trackY, thumbX+thumbW, trackY+th)
		}
	}
	return data
}

func toNRGBA(c color.Color) color.NRGBA {
	if c == nil {
		return color.NRGBA{}
	}
	if n, ok := c.(color.NRGBA); ok {
		return n
	}
	r, g, b, a := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
}

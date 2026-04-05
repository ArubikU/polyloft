package input

import (
	"image"
	"math"

	"gioui.org/io/key"
)

func normalizePointerScrollDelta(delta float32) float32 {
	if delta == 0 {
		return 0
	}
	absDelta := math.Abs(float64(delta))
	switch {
	case absDelta < 1:
		return delta * 28
	case absDelta < 4:
		return delta * 12
	default:
		return delta
	}
}

type ScrollHitInfo struct {
	Rect   image.Rectangle
	HasV   bool
	HasH   bool
	TrackV image.Rectangle
	ThumbV image.Rectangle
	TrackH image.Rectangle
	ThumbH image.Rectangle
	MaxX   int
	MaxY   int
}

type ScrollbarPressOutcome struct {
	Handled  bool
	Changed  bool
	DragAxis string
	DragGrab int
}

func SelectScrollTarget(hits map[string]ScrollHitInfo, preferredPath string, p image.Point) string {
	if hit, ok := hits[preferredPath]; ok && (hit.HasV || hit.HasH) {
		return preferredPath
	}
	bestPath := ""
	bestArea := int(^uint(0) >> 1)
	for path, hit := range hits {
		if !(hit.HasV || hit.HasH) {
			continue
		}
		if !pointInRect(p, hit.Rect) {
			continue
		}
		area := hit.Rect.Dx() * hit.Rect.Dy()
		if area < bestArea {
			bestArea = area
			bestPath = path
		}
	}
	return bestPath
}

func ApplyPointerScroll(cur image.Point, carryX, carryY float32, scrollX, scrollY float32) (image.Point, float32, float32) {
	dx := carryX + normalizePointerScrollDelta(scrollX)
	dy := carryY + normalizePointerScrollDelta(scrollY)

	stepX := 0
	for dx >= 1 {
		stepX++
		dx -= 1
	}
	for dx <= -1 {
		stepX--
		dx += 1
	}

	stepY := 0
	for dy >= 1 {
		stepY++
		dy -= 1
	}
	for dy <= -1 {
		stepY--
		dy += 1
	}

	cur.X += stepX
	cur.Y += stepY
	return cur, dx, dy
}

func ClampScrollOffset(cur image.Point, hit ScrollHitInfo) (image.Point, bool) {
	old := cur
	if cur.X < 0 {
		cur.X = 0
	}
	if cur.Y < 0 {
		cur.Y = 0
	}
	if cur.X > hit.MaxX {
		cur.X = hit.MaxX
	}
	if cur.Y > hit.MaxY {
		cur.Y = hit.MaxY
	}
	return cur, cur != old
}

func ApplyKeyboardScroll(cur image.Point, hit ScrollHitInfo, keyName key.Name) (image.Point, bool) {
	old := cur
	lineY := 36
	lineX := 36
	pageY := max(1, hit.Rect.Dy()-24)

	switch keyName {
	case key.NameDownArrow:
		cur.Y += lineY
	case key.NameUpArrow:
		cur.Y -= lineY
	case key.NameRightArrow:
		cur.X += lineX
	case key.NameLeftArrow:
		cur.X -= lineX
	case key.NamePageDown:
		cur.Y += pageY
	case key.NamePageUp:
		cur.Y -= pageY
	case key.NameHome:
		cur.Y = 0
		cur.X = 0
	case key.NameEnd:
		cur.Y = hit.MaxY
	}

	next, clamped := ClampScrollOffset(cur, hit)
	if clamped {
		return next, true
	}
	return next, next != old
}

func UpdateScrollByThumb(cur image.Point, hit ScrollHitInfo, p image.Point, axis string, scrollDragGrab int) (image.Point, bool) {
	changed := false
	if axis == "y" && hit.HasV {
		track := hit.TrackV
		thumb := hit.ThumbV
		usable := max(1, track.Dy()-thumb.Dy())
		pos := p.Y - track.Min.Y - scrollDragGrab
		if pos < 0 {
			pos = 0
		}
		if pos > usable {
			pos = usable
		}
		next := 0
		if hit.MaxY > 0 {
			next = int(float64(pos) * float64(hit.MaxY) / float64(usable))
		}
		if next != cur.Y {
			cur.Y = next
			changed = true
		}
	}
	if axis == "x" && hit.HasH {
		track := hit.TrackH
		thumb := hit.ThumbH
		usable := max(1, track.Dx()-thumb.Dx())
		pos := p.X - track.Min.X - scrollDragGrab
		if pos < 0 {
			pos = 0
		}
		if pos > usable {
			pos = usable
		}
		next := 0
		if hit.MaxX > 0 {
			next = int(float64(pos) * float64(hit.MaxX) / float64(usable))
		}
		if next != cur.X {
			cur.X = next
			changed = true
		}
	}
	return cur, changed
}

func HandleScrollbarPress(cur image.Point, hit ScrollHitInfo, p image.Point) (ScrollbarPressOutcome, image.Point) {
	if hit.HasV && pointInRect(p, hit.ThumbV) {
		return ScrollbarPressOutcome{Handled: true, DragAxis: "y", DragGrab: p.Y - hit.ThumbV.Min.Y}, cur
	}
	if hit.HasH && pointInRect(p, hit.ThumbH) {
		return ScrollbarPressOutcome{Handled: true, DragAxis: "x", DragGrab: p.X - hit.ThumbH.Min.X}, cur
	}
	if hit.HasV && pointInRect(p, hit.TrackV) {
		out := ScrollbarPressOutcome{Handled: true, DragAxis: "y", DragGrab: hit.ThumbV.Dy() / 2}
		next, changed := UpdateScrollByThumb(cur, hit, p, "y", out.DragGrab)
		out.Changed = changed
		return out, next
	}
	if hit.HasH && pointInRect(p, hit.TrackH) {
		out := ScrollbarPressOutcome{Handled: true, DragAxis: "x", DragGrab: hit.ThumbH.Dx() / 2}
		next, changed := UpdateScrollByThumb(cur, hit, p, "x", out.DragGrab)
		out.Changed = changed
		return out, next
	}
	return ScrollbarPressOutcome{}, cur
}

func ClampScrollToContent(cur image.Point, x, y, w, h, contentRight, contentBottom int) image.Point {
	maxScrollX := max(0, contentRight-(x+w))
	maxScrollY := max(0, contentBottom-(y+h))
	if cur.X > maxScrollX {
		cur.X = maxScrollX
	}
	if cur.Y > maxScrollY {
		cur.Y = maxScrollY
	}
	if cur.X < 0 {
		cur.X = 0
	}
	if cur.Y < 0 {
		cur.Y = 0
	}
	return cur
}

func pointInRect(p image.Point, r image.Rectangle) bool {
	return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y
}




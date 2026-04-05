package render

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	coreengine "github.com/ArubikU/giocss/core/engine"
)

// drawGioRect draws a solid colored rectangle at absolute position (x,y) with size (w,h).
func drawGioRect(ops *op.Ops, x, y, w, h int, col color.NRGBA) {
	if col.A == 0 || w <= 0 || h <= 0 {
		return
	}
	defer clip.Rect(image.Rect(x, y, x+w, y+h)).Push(ops).Pop()
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func DrawGioRect(ops *op.Ops, x, y, w, h int, col color.NRGBA) {
	drawGioRect(ops, x, y, w, h, col)
}

type cornerRadii struct {
	nw int
	ne int
	se int
	sw int
}

func (r cornerRadii) max() int {
	m := r.nw
	if r.ne > m {
		m = r.ne
	}
	if r.se > m {
		m = r.se
	}
	if r.sw > m {
		m = r.sw
	}
	return m
}

func (r cornerRadii) allEqual() bool {
	return r.nw == r.ne && r.ne == r.se && r.se == r.sw
}

func drawGioRRectCorners(ops *op.Ops, x, y, w, h int, radii cornerRadii, col color.NRGBA) {
	if col.A == 0 || w <= 0 || h <= 0 {
		return
	}
	if radii.max() <= 0 {
		drawGioRect(ops, x, y, w, h, col)
		return
	}
	maxR := min(w, h) / 2
	nw := min(max(0, radii.nw), maxR)
	ne := min(max(0, radii.ne), maxR)
	se := min(max(0, radii.se), maxR)
	sw := min(max(0, radii.sw), maxR)
	defer clip.RRect{
		Rect: image.Rect(x, y, x+w, y+h),
		NW:   nw, NE: ne, SE: se, SW: sw,
	}.Push(ops).Pop()
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// drawGioRRect draws a solid rounded rectangle.
func drawGioRRect(ops *op.Ops, x, y, w, h, radius int, col color.NRGBA) {
	drawGioRRectCorners(ops, x, y, w, h, cornerRadii{nw: radius, ne: radius, se: radius, sw: radius}, col)
}

func DrawGioRRect(rctx *coreengine.RenderContext, x, y, w, h, radius int, col color.NRGBA) {
	drawGioRRect(rctx.Gtx.Ops, x, y, w, h, radius, col)
}

// drawGioBorder draws a border outline. Supports border-radius via a stroked rounded-rect path.
func drawGioBorder(ops *op.Ops, x, y, w, h, radius, borderW int, col color.NRGBA) {
	if col.A == 0 || borderW <= 0 || w <= 0 || h <= 0 {
		return
	}
	if radius <= 0 {
		drawGioRect(ops, x, y, w, borderW, col)
		drawGioRect(ops, x, y+h-borderW, w, borderW, col)
		if h > 2*borderW {
			drawGioRect(ops, x, y+borderW, borderW, h-2*borderW, col)
			drawGioRect(ops, x+w-borderW, y+borderW, borderW, h-2*borderW, col)
		}
		return
	}
	maxR := float32(min(w, h)) / 2
	r := float32(radius)
	if r > maxR {
		r = maxR
	}
	const k = float32(0.5523)
	hw := float32(borderW) / 2
	x0 := float32(x) + hw
	y0 := float32(y) + hw
	w0 := float32(w) - 2*hw
	h0 := float32(h) - 2*hw
	if r > w0/2 {
		r = w0 / 2
	}
	if r > h0/2 {
		r = h0 / 2
	}
	var p clip.Path
	p.Begin(ops)
	p.Move(f32.Pt(x0+r, y0))
	p.Line(f32.Pt(w0-2*r, 0))
	p.Cube(f32.Pt(k*r, 0), f32.Pt(r, r-k*r), f32.Pt(r, r))
	p.Line(f32.Pt(0, h0-2*r))
	p.Cube(f32.Pt(0, k*r), f32.Pt(-(r-k*r), r), f32.Pt(-r, r))
	p.Line(f32.Pt(-(w0 - 2*r), 0))
	p.Cube(f32.Pt(-k*r, 0), f32.Pt(-r, -(r-k*r)), f32.Pt(-r, -r))
	p.Line(f32.Pt(0, -(h0 - 2*r)))
	p.Cube(f32.Pt(0, -k*r), f32.Pt(r-k*r, -r), f32.Pt(r, -r))
	p.Close()
	cs := clip.Stroke{Path: p.End(), Width: float32(borderW)}.Op().Push(ops)
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
	cs.Pop()
}

func DrawGioBorder(rctx *coreengine.RenderContext, x, y, w, h, radius, borderW int, col color.NRGBA) {
	drawGioBorder(rctx.Gtx.Ops, x, y, w, h, radius, borderW, col)
}

func drawGioDashedBorder(ops *op.Ops, x, y, w, h, radius, borderW int, col color.NRGBA, dotted bool) {
	if col.A == 0 || borderW <= 0 || w <= 0 || h <= 0 {
		return
	}
	var cs clip.Stack
	if radius > 0 {
		r := min(radius, min(w, h)/2)
		cs = clip.RRect{Rect: image.Rect(x, y, x+w, y+h), NW: r, NE: r, SE: r, SW: r}.Push(ops)
	} else {
		cs = clip.Rect(image.Rect(x, y, x+w, y+h)).Push(ops)
	}

	for _, r := range BuildDashedBorderRects(x, y, w, h, borderW, dotted) {
		drawGioRect(ops, r.X, r.Y, r.W, r.H, col)
	}

	cs.Pop()
}

// drawGioElementBorder renders borders respecting per-side declarations.
func drawGioElementBorder(ops *op.Ops, x, y, w, h int, radii cornerRadii, css map[string]string) {
	plan := BuildElementBorderPlan(css, w, h)
	top := plan.Top
	right := plan.Right
	bottom := plan.Bottom
	left := plan.Left

	if !plan.HasAny {
		return
	}

	if plan.Uniform {
		style := plan.UniformStyle
		if style == "dashed" || style == "dotted" {
			drawGioDashedBorder(ops, x, y, w, h, radii.max(), plan.UniformW, plan.UniformCol, style == "dotted")
		} else {
			if radii.allEqual() {
				drawGioBorder(ops, x, y, w, h, radii.nw, plan.UniformW, plan.UniformCol)
			} else {
				var cs clip.Stack
				cs = clip.RRect{Rect: image.Rect(x, y, x+w, y+h), NW: radii.nw, NE: radii.ne, SE: radii.se, SW: radii.sw}.Push(ops)
				drawGioRect(ops, x, y, w, top.Width, top.Color)
				drawGioRect(ops, x, y+h-bottom.Width, w, bottom.Width, bottom.Color)
				innerTop := top.Width
				innerBot := bottom.Width
				innerH := h - innerTop - innerBot
				if innerH > 0 {
					drawGioRect(ops, x, y+innerTop, left.Width, innerH, left.Color)
					drawGioRect(ops, x+w-right.Width, y+innerTop, right.Width, innerH, right.Color)
				}
				cs.Pop()
			}
		}
		return
	}

	var cs clip.Stack
	useClip := radii.max() > 0
	if useClip {
		cs = clip.RRect{Rect: image.Rect(x, y, x+w, y+h), NW: radii.nw, NE: radii.ne, SE: radii.se, SW: radii.sw}.Push(ops)
	}
	if top.Width > 0 && top.Color.A > 0 {
		drawGioRect(ops, x, y, w, top.Width, top.Color)
	}
	if bottom.Width > 0 && bottom.Color.A > 0 {
		drawGioRect(ops, x, y+h-bottom.Width, w, bottom.Width, bottom.Color)
	}
	innerTop := top.Width
	innerBot := bottom.Width
	innerH := h - innerTop - innerBot
	if innerH > 0 {
		if left.Width > 0 && left.Color.A > 0 {
			drawGioRect(ops, x, y+innerTop, left.Width, innerH, left.Color)
		}
		if right.Width > 0 && right.Color.A > 0 {
			drawGioRect(ops, x+w-right.Width, y+innerTop, right.Width, innerH, right.Color)
		}
	}
	if useClip {
		cs.Pop()
	}
}

func drawDebugOutline(ops *op.Ops, x, y, w, h, stroke int, col color.NRGBA) {
	if col.A == 0 {
		return
	}
	for _, r := range BuildDebugOutlineRects(x, y, w, h, stroke) {
		drawGioRect(ops, r.X, r.Y, r.W, r.H, col)
	}
}

package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #f8fafc;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 32px;
}

.palette-title {
  font-size: 22px;
  font-weight: bold;
  color: #1e293b;
  margin-bottom: 24px;
}

.swatches-row {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16px;
}

.swatch {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 80px;
}

.swatch-block {
  width: 80px;
  height: 56px;
  border-radius: 10px;
  margin-bottom: 6px;
}

.swatch-hex {
  font-size: 11px;
  color: #475569;
  font-family: monospace;
  text-align: center;
}

.swatch-name {
  font-size: 11px;
  color: #0f172a;
  font-weight: bold;
  text-align: center;
}
`

type swatchData struct {
	hex  string
	name string
}

var palette = []swatchData{
	{"#ef4444", "Red"},
	{"#f97316", "Orange"},
	{"#eab308", "Yellow"},
	{"#22c55e", "Green"},
	{"#14b8a6", "Teal"},
	{"#3b82f6", "Blue"},
	{"#8b5cf6", "Violet"},
	{"#ec4899", "Pink"},
	{"#64748b", "Slate"},
	{"#1e293b", "Dark"},
	{"#f1f5f9", "Light"},
	{"#ffffff", "White"},
}

func buildSwatch(d swatchData) *giocss.Node {
	block := giocss.NewNode("div")
	block.AddClass("swatch-block")
	block.SetProp("style.background-color", d.hex)

	swatch := components.Column(
		block,
		components.Span(d.hex, "swatch-hex"),
		components.Span(d.name, "swatch-name"),
	)
	swatch.AddClass("swatch")
	return swatch
}

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	shell := components.Column()
	shell.AddClass("palette-shell")
	shell.AddChild(components.Heading(2, "Color Palette", "palette-title"))

	row := giocss.NewNode("div")
	row.AddClass("swatches-row")
	for _, s := range palette {
		row.AddChild(buildSwatch(s))
	}
	shell.AddChild(row)
	root.AddChild(shell)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 07 – Color Swatches", Width: 760, Height: 360},
		giocss.WindowRuntimeHooks{
			Snapshot: func(size image.Point) giocss.WindowRuntimeSnapshot {
				root := buildUI()
				return giocss.WindowRuntimeSnapshot{
					RootLayout:   giocss.LayoutNodeToNative(root, size.X, size.Y, ss),
					RootCSS:      giocss.ResolveNodeStyle(root, ss, size.X),
					StyleSheet:   ss,
					ScreenWidth:  size.X,
					ScreenHeight: size.Y,
				}
			},
		},
	)
	rt.Run()
}

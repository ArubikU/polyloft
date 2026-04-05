package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #f1f5f9;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
	align-items: flex-start;
	padding: 24px 0;
	overflow: auto;
}

.cards-grid {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 20px;
  padding: 24px;
  justify-content: center;
	align-content: flex-start;
	align-items: flex-start;
	width: 720px;
	max-width: 100%;
	height: 100%;
	overflow: auto;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 24px;
  width: 280px;
	min-height: 280px;
  display: flex;
  flex-direction: column;
}

.card:hover {
  background-color: #f8fafc;
}

.card-thumb {
  background-color: #e2e8f0;
  border-radius: 8px;
  height: 120px;
  width: 100%;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-thumb-label {
  color: #94a3b8;
  font-size: 13px;
}

.card-title {
  font-size: 17px;
  font-weight: bold;
  color: #1e293b;
  margin-bottom: 6px;
}

.card-body {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 16px;
	line-height: 1.45;
	max-height: 38px;
	overflow: hidden;
}

.card-tag {
  font-size: 11px;
  font-weight: bold;
  color: #3b82f6;
  background-color: #eff6ff;
  border-radius: 999px;
  padding: 3px 10px;
  align-self: flex-start;
}
`

type cardData struct {
	title, body, tag string
}

func buildCard(d cardData) *giocss.Node {
	thumb := components.Div([]string{"card-thumb"},
		components.Span("Image placeholder", "card-thumb-label"),
	)
	body := giocss.NewNode("p")
	body.Text = d.body
	body.AddClass("card-body")
	tag := components.Span(d.tag, "card-tag")
	return components.Card(
		thumb,
		components.Heading(3, d.title, "card-title"),
		body,
		tag,
	)
}

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	items := []cardData{
		{"Mountain Retreat", "Escape the city and breathe fresh air in the highlands.", "Travel"},
		{"TypeScript Tips", "10 lesser-known tricks that will improve your codebase.", "Dev"},
		{"Minimalist Design", "Why less is more — lessons from Swiss graphic design.", "Design"},
		{"Home Cooking", "Simple weeknight dinners that taste like restaurant food.", "Food"},
	}

	grid := giocss.NewNode("div")
	grid.AddClass("cards-grid")
	for _, item := range items {
		grid.AddChild(buildCard(item))
	}
	root.AddChild(grid)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 04 – Cards", Width: 720, Height: 520},
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

package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #f0f4f8;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.btn-row {
  display: flex;
  flex-direction: row;
  padding: 32px 40px;
  background-color: #ffffff;
  border-radius: 12px;
  gap: 16px;
  align-items: center;
}

.btn {
  padding: 10px 28px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: bold;
  cursor: pointer;
}

.btn-primary {
  background-color: #3b82f6;
  color: #ffffff;
}

.btn-primary:hover {
  background-color: #2563eb;
}

.btn-primary:active {
  background-color: #1d4ed8;
}

.btn-secondary {
  background-color: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover {
  background-color: #d1d5db;
}

.btn-secondary:active {
  background-color: #9ca3af;
}

.btn-danger {
  background-color: #ef4444;
  color: #ffffff;
}

.btn-danger:hover {
  background-color: #dc2626;
}

.btn-danger:active {
  background-color: #b91c1c;
}

.btn-ghost {
  background-color: transparent;
  color: #6b7280;
  border: 2px solid #d1d5db;
}

.btn-ghost:hover {
  background-color: #f9fafb;
  color: #374151;
}
`

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	row := components.Row(
		components.Button("Primary", "btn", "btn-primary"),
		components.Button("Secondary", "btn", "btn-secondary"),
		components.Button("Danger", "btn", "btn-danger"),
		components.Button("Ghost", "btn", "btn-ghost"),
	)
	row.AddClass("btn-row")
	root.AddChild(row)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 02 – Buttons", Width: 700, Height: 220},
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

package main

import (
	"image"
	"path/filepath"
	"runtime"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #1a1a2e;
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background-color: #16213e;
  border-radius: 16px;
  width: 480px;
  height: 260px;
}

.hero-title {
  font-size: 36px;
  font-weight: bold;
  color: #e94560;
  text-align: center;
  margin-bottom: 12px;
}

.hero-subtitle {
  font-size: 16px;
  color: #a8b2d8;
  text-align: center;
}
`

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	hero := components.Column(
		components.Heading(1, "Hello, giocss!", "hero-title"),
		components.Text("A CSS-driven UI toolkit for Gio.", "hero-subtitle"),
	)
	hero.AddClass("hero")
	root.AddChild(hero)

	return root
}

func sampleCSSPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "style.css"
	}
	return filepath.Join(filepath.Dir(file), "style.css")
}

func main() {
	ss := giocss.NewStyleSheet()
	if _, err := ss.LoadFile(sampleCSSPath()); err != nil {
		// Keep embedded CSS fallback for portability in non-standard run contexts.
		ss.ParseCSSText(css)
	}

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 01 – Hello World", Width: 700, Height: 500},
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

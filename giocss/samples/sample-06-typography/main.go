package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #ffffff;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 40px 60px;
  overflow: scroll;
}

h1 { font-size: 40px; font-weight: bold; color: #111827; margin-bottom: 6px; }
h2 { font-size: 32px; font-weight: bold; color: #1f2937; margin-bottom: 6px; }
h3 { font-size: 24px; font-weight: bold; color: #374151; margin-bottom: 6px; }
h4 { font-size: 20px; font-weight: bold; color: #4b5563; margin-bottom: 6px; }
h5 { font-size: 16px; font-weight: bold; color: #6b7280; margin-bottom: 6px; }
h6 { font-size: 13px; font-weight: bold; color: #9ca3af; margin-bottom: 12px; }

.divider {
  height: 1px;
  background-color: #e5e7eb;
  margin: 20px 0;
  width: 100%;
}

.body-text {
  font-size: 16px;
  color: #374151;
  line-height: 1.6;
  margin-bottom: 16px;
  max-width: 620px;
}

.text-muted  { color: #9ca3af; }
.text-accent { color: #7c3aed; font-weight: bold; }
.text-small  { font-size: 12px; color: #6b7280; }
.text-upper  { text-transform: uppercase; letter-spacing: 2px; font-size: 12px; color: #6b7280; }
.text-mono   { font-family: monospace; background-color: #f3f4f6; padding: 2px 6px;
               border-radius: 4px; font-size: 13px; color: #374151; }
`

func divider() *giocss.Node {
	d := giocss.NewNode("div")
	d.AddClass("divider")
	return d
}

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	root.AddChild(components.Heading(1, "The quick brown fox", "t-h1"))
	root.AddChild(components.Heading(2, "Jumps over the lazy dog", "t-h2"))
	root.AddChild(components.Heading(3, "Typography showcase", "t-h3"))
	root.AddChild(components.Heading(4, "Heading level four", "t-h4"))
	root.AddChild(components.Heading(5, "Heading level five", "t-h5"))
	root.AddChild(components.Heading(6, "Heading level six", "t-h6"))
	root.AddChild(divider())

	root.AddChild(components.Text(
		"Body text at 16px with a comfortable line-height for reading long-form content. "+
			"giocss resolves font, color, and spacing via CSS cascade — programmatic or from a file.",
		"body-text",
	))

	root.AddChild(components.Text("Muted secondary text", "body-text", "text-muted"))
	root.AddChild(components.Text("Accented bold highlight", "body-text", "text-accent"))
	root.AddChild(divider())

	root.AddChild(components.Text("UPPERCASE TRACKING LABEL", "text-upper"))
	root.AddChild(components.Text("Small helper caption at 12px", "text-small"))
	root.AddChild(components.Text("NewWindowRuntime(opts, hooks)", "text-mono"))

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 06 – Typography", Width: 760, Height: 620},
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

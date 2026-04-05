package main

import (
	"image"
	"strings"
	"sync"

	giocss "github.com/ArubikU/giocss"
)

const css = `
body {
  background-color: #f8fafc;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: 0;
}

h1 {
  font-size: 32px;
  color: #1e293b;
  margin-bottom: 32px;
}

.open-button {
  padding: 12px 24px;
  background-color: #3b82f6;
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
}
.open-button:hover { background-color: #2563eb; }

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  padding: 32px;
  border-radius: 12px;
  width: 400px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: 0 10px 15px -3px rgba(0,0,0,0.1);
}

.modal-title {
  font-size: 24px;
  color: #0f172a;
  margin: 0;
  font-weight: bold;
}

.modal-body {
  color: #475569;
  line-height: 1.5;
  margin: 0;
}

.close-button {
  margin-top: 16px;
  padding: 8px 16px;
  background-color: #ef4444;
  color: white;
  border-radius: 6px;
  align-self: flex-end;
  cursor: pointer;
  font-weight: bold;
}
.close-button:hover { background-color: #dc2626; }
`

type Modal struct {
	IsOpen bool
	Title  string
	Body   string
}

func (m *Modal) Open(title, body string) {
	m.IsOpen = true
	m.Title = title
	m.Body = body
}

func (m *Modal) Close() {
	m.IsOpen = false
}

type appState struct {
	mu    sync.Mutex
	modal *Modal
}

func newAppState() *appState {
	return &appState{modal: &Modal{}}
}

func (s *appState) onEvent(eventName string, payload map[string]any) error {
	eventName = strings.TrimSpace(eventName)
	
	s.mu.Lock()
	defer s.mu.Unlock()

	switch eventName {
	case "modal.open":
		s.modal.Open("Hello, Method Component!", "This modal is dynamically injected via component methods. It uses position: absolute and z-index to overlay the main content.")
	case "modal.close":
		s.modal.Close()
	}
	return nil
}

func (s *appState) snapshot() (bool, string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.modal.IsOpen, s.modal.Title, s.modal.Body
}

func buildUI(isOpen bool, title string, body string) *giocss.Node {
	root := giocss.NewNode("body")

	h1 := giocss.NewNode("h1")
	h1.Text = "Modal Method Demo"
	root.AddChild(h1)

	btn := giocss.NewNode("button")
	btn.AddClass("open-button")
	btn.Text = "Open Modal"
	btn.SetProp("onclick", "modal.open")
	root.AddChild(btn)

	if isOpen {
		overlay := giocss.NewNode("div")
		overlay.AddClass("modal-overlay")
		overlay.SetProp("onclick", "modal.close") // Close when clicking background

		content := giocss.NewNode("div")
		content.AddClass("modal-content")
		content.SetProp("onclick", "no-op") // Swallow clicks inside modal

		h2 := giocss.NewNode("h2")
		h2.AddClass("modal-title")
		h2.Text = title
		
		p := giocss.NewNode("p")
		p.AddClass("modal-body")
		p.Text = body
		
		closeBtn := giocss.NewNode("button")
		closeBtn.AddClass("close-button")
		closeBtn.Text = "Close"
		closeBtn.SetProp("onclick", "modal.close")

		content.AddChild(h2)
		content.AddChild(p)
		content.AddChild(closeBtn)
		
		overlay.AddChild(content)
		root.AddChild(overlay)
	}

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)
	app := newAppState()

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 11 - Modals via Absolute/Z-Index", Width: 800, Height: 600},
		giocss.WindowRuntimeHooks{
			DispatchEvent: app.onEvent,
			Snapshot: func(size image.Point) giocss.WindowRuntimeSnapshot {
				isOpen, title, body := app.snapshot()
				root := buildUI(isOpen, title, body)
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

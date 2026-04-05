package main

import (
	"fmt"
	"image"
	"strings"
	"sync"

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
	align-items: stretch;
	overflow: hidden;
}

.navbar {
	background-color: #0f172a;
  width: 100%;
	height: 60px;
  display: flex;
  flex-direction: row;
  align-items: center;
	justify-content: space-between;
	padding: 0 18px;
}

.navbar-brand {
	font-size: 16px;
  font-weight: bold;
	color: #e2e8f0;
}

.navbar-items {
  display: flex;
  flex-direction: row;
  align-items: center;
	gap: 10px;
	min-width: 0;
	overflow-x: auto;
}

.navbar-item {
  font-size: 14px;
	font-weight: bold;
	color: #cbd5e1;
	padding: 8px 14px;
  border-radius: 6px;
	border: 1px solid #1e293b;
	background-color: #111827;
  cursor: pointer;
}

.navbar-item:hover {
  color: #f1f5f9;
	background-color: #1e293b;
}

.navbar-item-active {
	color: #0f172a;
	background-color: #38bdf8;
	border-color: #38bdf8;
}

.navbar-badge {
	font-size: 10px;
	font-weight: bold;
	color: #ffffff;
	background-color: #f43f5e;
	border-radius: 999px;
	padding: 2px 7px;
}

.page-content {
	width: 100%;
	flex: 1;
	min-height: 0;
	overflow: auto;
	padding: 22px;
  display: flex;
  flex-direction: column;
}

.page-title {
	font-size: 30px;
  font-weight: bold;
	color: #0f172a;
	margin-bottom: 8px;
}

.page-desc {
  font-size: 15px;
	color: #475569;
	margin-bottom: 22px;
}

.section {
	background-color: #ffffff;
	border-radius: 12px;
	padding: 18px 20px;
	margin-bottom: 12px;
	display: flex;
	flex-direction: column;
}

.section-title {
	font-size: 16px;
  font-weight: bold;
	color: #1e293b;
	margin-bottom: 6px;
}

.section-body {
	font-size: 13px;
	color: #64748b;
	line-height: 1.45;
}
`

type pageSection struct {
	title string
	body  string
}

type pageModel struct {
	Title       string
	Description string
	Sections    []pageSection
}

type appState struct {
	mu         sync.Mutex
	activePage string
}

func newAppState() *appState {
	return &appState{activePage: "home"}
}

func (s *appState) onEvent(eventName string, _ map[string]any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch strings.TrimSpace(eventName) {
	case "nav.home":
		s.activePage = "home"
	case "nav.components":
		s.activePage = "components"
	case "nav.docs":
		s.activePage = "docs"
	case "nav.changelog":
		s.activePage = "changelog"
	}
	return nil
}

func (s *appState) snapshot() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activePage
}

func buildPageModel(activePage string) pageModel {
	makeSections := func(prefix string, detail string, count int) []pageSection {
		sections := make([]pageSection, 0, count)
		for i := 1; i <= count; i++ {
			sections = append(sections, pageSection{
				title: fmt.Sprintf("%s %d", prefix, i),
				body:  detail,
			})
		}
		return sections
	}

	switch activePage {
	case "components":
		return pageModel{
			Title:       "Components",
			Description: "Reusable primitives, cards, forms, and interactive controls available in the sample set.",
			Sections:    makeSections("Component", "Use this page to validate scrolling while switching between component-focused sections.", 16),
		}
	case "docs":
		return pageModel{
			Title:       "Docs",
			Description: "Layout, style resolution, rendering and runtime notes presented as long-form content blocks.",
			Sections:    makeSections("Guide", "This content simulates documentation sections so the page can exercise vertical scrolling reliably.", 18),
		}
	case "changelog":
		return pageModel{
			Title:       "Changelog",
			Description: "Recent updates across samples and runtime behavior, keeping the same visual shell while content changes.",
			Sections:    makeSections("Release", "Each release card stands in for a changelog entry and gives the page enough content to scroll.", 15),
		}
	default:
		return pageModel{
			Title:       "Navigation",
			Description: "Top nav with visible labels and a scrollable content area.",
			Sections:    makeSections("Section", "This section exists to validate vertical scrolling and readable text rendering in navigation sample.", 14),
		}
	}
}

func buildUI(activePage string) *giocss.Node {
	root := giocss.NewNode("body")
	page := buildPageModel(activePage)

	// Navbar
	nav := giocss.NewNode("nav")
	nav.AddClass("navbar")

	brand := components.Span("giocss", "navbar-brand")
	nav.AddChild(brand)

	navItems := components.Row()
	navItems.AddClass("navbar-items")

	items := []struct {
		label  string
		key    string
		event  string
		badge  string
	}{
		{"Home", "home", "nav.home", ""},
		{"Components", "components", "nav.components", ""},
		{"Docs", "docs", "nav.docs", ""},
		{"Changelog", "changelog", "nav.changelog", "New"},
	}
	for _, item := range items {
		a := components.Button(item.label, "navbar-item")
		a.SetProp("event", item.event)
		if item.key == activePage {
			a.AddClass("navbar-item-active")
		}
		navItems.AddChild(a)
		if item.badge != "" {
			navItems.AddChild(components.Badge(item.badge, "navbar-badge"))
		}
	}
	nav.AddChild(navItems)
	root.AddChild(nav)

	// Page content below navbar (long enough to require scroll)
	content := components.Column(
		components.Heading(1, page.Title, "page-title"),
		components.Span(page.Description, "page-desc"),
	)
	content.AddClass("page-content")

	for i, sectionData := range page.Sections {
		section := components.Column(
			components.Heading(3, sectionData.title, "section-title"),
			components.Span(sectionData.body, "section-body"),
		)
		section.SetProp("id", fmt.Sprintf("section-%d", i+1))
		section.AddClass("section")
		content.AddChild(section)
	}

	root.AddChild(content)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)
	app := newAppState()

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 05 – Navigation", Width: 800, Height: 400},
		giocss.WindowRuntimeHooks{
			DispatchEvent: app.onEvent,
			Snapshot: func(size image.Point) giocss.WindowRuntimeSnapshot {
				root := buildUI(app.snapshot())
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

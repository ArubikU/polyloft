package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #0f172a;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: row;
}

.sidebar {
  width: 220px;
  height: 100%;
  background-color: #1e293b;
  display: flex;
  flex-direction: column;
  padding: 24px 0;
}

.sidebar-logo {
  font-size: 16px;
  font-weight: bold;
  color: #38bdf8;
  padding: 0 20px;
  margin-bottom: 28px;
}

.sidebar-section {
  font-size: 10px;
  font-weight: bold;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 2px;
  padding: 0 20px;
  margin-bottom: 8px;
  margin-top: 16px;
}

.sidebar-item {
  font-size: 14px;
  color: #94a3b8;
  padding: 9px 20px;
  cursor: pointer;
}

.sidebar-item:hover {
  color: #e2e8f0;
  background-color: #0f172a;
}

.sidebar-item-active {
  color: #f8fafc;
  background-color: #0f172a;
  border-left: 3px solid #38bdf8;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 32px 36px;
  overflow: scroll;
}

.topbar {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-bottom: 28px;
}

.topbar-title {
  font-size: 20px;
  font-weight: bold;
  color: #f1f5f9;
  flex: 1;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  background-color: #38bdf8;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-initials {
  font-size: 13px;
  font-weight: bold;
  color: #0f172a;
}

.cards-row {
  display: flex;
  flex-direction: row;
  gap: 16px;
  margin-bottom: 28px;
}

.dark-card {
  background-color: #1e293b;
  border-radius: 12px;
  padding: 20px 22px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.dark-card:hover {
  background-color: #243247;
}

.dark-card-value {
  font-size: 30px;
  font-weight: bold;
  color: #38bdf8;
  margin-bottom: 6px;
}

.dark-card-label {
  font-size: 13px;
  color: #64748b;
}

.section-title {
  font-size: 15px;
  font-weight: bold;
  color: #cbd5e1;
  margin-bottom: 12px;
}

.activity-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid #1e293b;
}

.activity-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background-color: #38bdf8;
  margin-right: 12px;
}

.activity-text {
  font-size: 14px;
  color: #94a3b8;
  flex: 1;
}

.activity-time {
  font-size: 12px;
  color: #475569;
}
`

func sidebarItem(label string, active bool) *giocss.Node {
	n := components.Text(label, "sidebar-item")
	n.SetProp("style.color", "#94a3b8")
	n.AddClass("sidebar-item")
	if active {
		n.AddClass("sidebar-item-active")
		n.SetProp("style.color", "#f8fafc")
	}
	return n
}

func themedText(content, className, color string) *giocss.Node {
	n := components.Text(content, className)
	if color != "" {
		n.SetProp("style.color", color)
	}
	return n
}

func darkCard(value, label string) *giocss.Node {
	c := components.Column(
		themedText(value, "dark-card-value", "#38bdf8"),
		themedText(label, "dark-card-label", "#64748b"),
	)
	c.AddClass("dark-card")
	return c
}

func activityItem(text, time string) *giocss.Node {
	dot := giocss.NewNode("div")
	dot.AddClass("activity-dot")

	row := components.Row(
		dot,
		themedText(text, "activity-text", "#94a3b8"),
		themedText(time, "activity-time", "#475569"),
	)
	row.AddClass("activity-item")
	return row
}

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")

	// ── Sidebar ──────────────────────────────────────────────────────────
	sidebar := giocss.NewNode("div")
	sidebar.AddClass("sidebar")

	logo := themedText("giocss", "sidebar-logo", "#38bdf8")
	sidebar.AddChild(logo)

	for _, s := range []struct {
		section string
		items   []struct {
			label  string
			active bool
		}
	}{
		{"General", []struct {
			label  string
			active bool
		}{
			{"Dashboard", true},
			{"Analytics", false},
			{"Reports", false},
		}},
		{"Settings", []struct {
			label  string
			active bool
		}{
			{"Profile", false},
			{"Preferences", false},
			{"Billing", false},
		}},
	} {
		sec := themedText(s.section, "sidebar-section", "#475569")
		sidebar.AddChild(sec)
		for _, it := range s.items {
			sidebar.AddChild(sidebarItem(it.label, it.active))
		}
	}

	root.AddChild(sidebar)

	// ── Main ─────────────────────────────────────────────────────────────
	main := giocss.NewNode("div")
	main.AddClass("main")

	// Topbar
	avatar := giocss.NewNode("div")
	avatar.AddClass("avatar")
	avatar.AddChild(components.Text("AW", "avatar-initials"))

	topbar := components.Row(
		themedText("Overview", "topbar-title", "#f1f5f9"),
		avatar,
	)
	topbar.AddClass("topbar")
	main.AddChild(topbar)

	// Stat cards
	cardsRow := giocss.NewNode("div")
	cardsRow.AddClass("cards-row")
	for _, d := range []struct{ v, l string }{
		{"2,481", "Active Users"},
		{"$12.4k", "MRR"},
		{"99.8%", "Uptime"},
	} {
		cardsRow.AddChild(darkCard(d.v, d.l))
	}
	main.AddChild(cardsRow)

	// Activity
	main.AddChild(themedText("Recent Activity", "section-title", "#cbd5e1"))
	for _, a := range []struct{ text, time string }{
		{"Alice pushed a commit to main", "2m ago"},
		{"Bob opened PR #142: fix layout overflow", "14m ago"},
		{"Deployment to production succeeded", "1h ago"},
		{"Carol commented on issue #87", "3h ago"},
		{"Nightly build passed (1m 42s)", "8h ago"},
	} {
		main.AddChild(activityItem(a.text, a.time))
	}

	root.AddChild(main)
	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 10 – Dark Theme", Width: 900, Height: 620},
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

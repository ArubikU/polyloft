package main

import (
	"image"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
body {
  background-color: #f0f4ff;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 32px 40px;
}

.dash-title {
  font-size: 22px;
  font-weight: bold;
  color: #1e293b;
  margin-bottom: 24px;
}

.stats-row {
  display: flex;
  flex-direction: row;
  gap: 16px;
  margin-bottom: 32px;
}

.stat-card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 20px 24px;
  width: 180px;
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #4f46e5;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #64748b;
}

.table-container {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 0;
  overflow: hidden;
  width: 640px;
}

.table-header {
  display: flex;
  flex-direction: row;
  background-color: #f1f5f9;
  padding: 12px 20px;
}

.table-row {
  display: flex;
  flex-direction: row;
  padding: 12px 20px;
  border-top: 1px solid #f1f5f9;
}

.table-row:hover {
  background-color: #fafbff;
}

.col-name  { width: 180px; font-size: 14px; color: #1e293b; font-weight: bold; }
.col-role  { width: 180px; font-size: 14px; color: #475569; }
.col-status{ width: 120px; font-size: 14px; }
.col-score { width: 80px; font-size: 14px; color: #4f46e5; font-weight: bold; text-align: right; }

.th { font-size: 12px; font-weight: bold; color: #94a3b8; text-transform: uppercase; letter-spacing: 1px; }

.badge { font-size: 11px; font-weight: bold; border-radius: 999px; padding: 2px 8px; }
.badge-active   { background-color: #dcfce7; color: #16a34a; }
.badge-inactive { background-color: #fee2e2; color: #dc2626; }
.badge-trial    { background-color: #fef3c7; color: #d97706; }
`

type statsItem struct {
	value, label string
}

type tableRow struct {
	name, role, status, score string
}

var stats = []statsItem{
	{"1,248", "Total Users"},
	{"94.2%", "Uptime"},
	{"$8,340", "Revenue"},
}

var rows = []tableRow{
	{"Alice Webb", "Frontend Dev", "active", "98"},
	{"Bob Torres", "Backend Dev", "active", "91"},
	{"Carol Lin", "Designer", "trial", "87"},
	{"David Park", "DevOps", "inactive", "74"},
	{"Eva Mller", "Product", "active", "95"},
}

func statCard(s statsItem) *giocss.Node {
	return components.Column(
		components.Text(s.value, "stat-value"),
		components.Text(s.label, "stat-label"),
	)
}

func tableRowNode(r tableRow) *giocss.Node {
	badgeClass := "badge-active"
	if r.status == "inactive" {
		badgeClass = "badge-inactive"
	} else if r.status == "trial" {
		badgeClass = "badge-trial"
	}

	row := giocss.NewNode("div")
	row.AddClass("table-row")
	row.AddChild(components.Text(r.name, "col-name"))
	row.AddChild(components.Text(r.role, "col-role"))
	statusCell := giocss.NewNode("div")
	statusCell.AddClass("col-status")
	statusCell.AddChild(components.Badge(r.status, badgeClass))
	row.AddChild(statusCell)
	row.AddChild(components.Text(r.score, "col-score"))
	return row
}

func buildUI() *giocss.Node {
	root := giocss.NewNode("body")
	root.AddChild(components.Heading(2, "Dashboard", "dash-title"))

	// Stat cards
	statsRow := giocss.NewNode("div")
	statsRow.AddClass("stats-row")
	for _, s := range stats {
		c := statCard(s)
		c.AddClass("stat-card")
		statsRow.AddChild(c)
	}
	root.AddChild(statsRow)

	// Table
	table := giocss.NewNode("div")
	table.AddClass("table-container")

	header := giocss.NewNode("div")
	header.AddClass("table-header")
	for _, h := range []string{"Name", "Role", "Status", "Score"} {
		cell := giocss.NewNode("span")
		cell.Text = h
		classes := map[string]string{
			"Name": "col-name th", "Role": "col-role th",
			"Status": "col-status th", "Score": "col-score th",
		}
		for _, cls := range []string{classes[h]} {
			cell.AddClass(cls)
		}
		header.AddChild(cell)
	}
	table.AddChild(header)

	for _, r := range rows {
		table.AddChild(tableRowNode(r))
	}
	root.AddChild(table)

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 09 – Dashboard", Width: 760, Height: 560},
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

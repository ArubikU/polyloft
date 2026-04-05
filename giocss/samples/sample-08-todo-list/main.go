package main

import (
	"image"
	"strconv"
	"strings"
	"sync"

	giocss "github.com/ArubikU/giocss"
	"github.com/ArubikU/giocss/components"
)

const css = `
.app-shell {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
	padding: 24px 30px;
  background-color: #f8fafc;
	color: #0f172a;
	overflow: auto;
}

.app-shell-dark {
  background-color: #0f172a;
	color: #e2e8f0;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
  color: #0f172a;
}

.title-dark {
  color: #e2e8f0;
}

.subtitle {
  font-size: 14px;
  margin-bottom: 18px;
  color: #475569;
}

.subtitle-dark {
  color: #94a3b8;
}

.controls {
  display: flex;
  flex-direction: row;
  gap: 10px;
  margin-bottom: 16px;
}

.action-btn {
  padding: 9px 16px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: bold;
  cursor: pointer;
}

.btn-inc { background-color: #22c55e; color: #052e16; }
.btn-dec { background-color: #f97316; color: #431407; }
.btn-reset { background-color: #64748b; color: #f8fafc; }
.btn-theme { background-color: #3b82f6; color: #eff6ff; }

.stats-card {
	width: 680px;
  background-color: rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 16px 18px;
  margin-bottom: 18px;
}

.stats-card-dark {
	background-color: rgba(15, 23, 42, 0.34);
}

.counter-label {
  font-size: 12px;
  letter-spacing: 1px;
  text-transform: uppercase;
  margin-bottom: 2px;
	color: #64748b;
}

.counter-label-dark { color: #94a3b8; }

.counter-value {
  font-size: 34px;
  font-weight: bold;
  margin-bottom: 8px;
	color: #0f172a;
}

.counter-value-dark { color: #e2e8f0; }

.counter-positive { color: #16a34a; }
.counter-negative { color: #dc2626; }

.mut-line {
  font-size: 13px;
	color: #334155;
}

.mut-line-dark { color: #cbd5e1; }

.tasks-title {
  font-size: 15px;
  font-weight: bold;
  margin-bottom: 8px;
	color: #1e293b;
}

.tasks-title-dark { color: #e2e8f0; }

.task-row {
	width: 680px;
  display: flex;
  flex-direction: row;
  align-items: center;
	gap: 10px;
	padding: 10px 12px;
  border-radius: 10px;
  margin-bottom: 8px;
	background-color: #ffffff;
}

.task-row-dark { background-color: #1e293b; }

.task-check {
	width: 22px;
	height: 22px;
}

.task-text {
  flex: 1;
  font-size: 14px;
	color: #1f2937;
}

.task-text-dark { color: #e2e8f0; }

.task-done {
  text-decoration: line-through;
  color: #94a3b8;
}
`

type task struct {
	label string
	done  bool
}

type appState struct {
	mu        sync.Mutex
	counter   int
	mutations int
	lastEvent string
	darkTheme bool
	tasks     []task
}

func newAppState() *appState {
	return &appState{
		counter:   0,
		mutations: 0,
		lastEvent: "boot",
		darkTheme: false,
		tasks: []task{
			{label: "Design state model", done: true},
			{label: "Wire click handlers", done: true},
			{label: "Exercise runtime mutation", done: false},
			{label: "Verify rerender on events", done: false},
		},
	}
}

type viewState struct {
	counter   int
	mutations int
	lastEvent string
	darkTheme bool
	tasks     []task
}

func (s *appState) snapshot() viewState {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := viewState{
		counter:   s.counter,
		mutations: s.mutations,
		lastEvent: s.lastEvent,
		darkTheme: s.darkTheme,
		tasks:     make([]task, len(s.tasks)),
	}
	copy(out.tasks, s.tasks)
	return out
}

func (s *appState) apply(eventName string, payload map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.mutations++
	s.lastEvent = eventName

	source := strings.TrimSpace(anyString(payload["source"]))
	if source != "" {
		s.lastEvent = eventName + " @ " + source
	}

	switch eventName {
	case "counter.increment":
		s.counter++
	case "counter.decrement":
		s.counter--
	case "counter.reset":
		s.counter = 0
	case "theme.toggle":
		s.darkTheme = !s.darkTheme
	default:
		if strings.HasPrefix(eventName, "task.toggle.") {
			idxRaw := strings.TrimPrefix(eventName, "task.toggle.")
			idx, err := strconv.Atoi(idxRaw)
			if err == nil && idx >= 0 && idx < len(s.tasks) {
				s.tasks[idx].done = !s.tasks[idx].done
			}
		}
	}
}

func anyString(v any) string {
	s, _ := v.(string)
	return s
}

func buildTaskRow(index int, t task, dark bool) *giocss.Node {
	check := components.Input("task-"+strconv.Itoa(index), "", "checkbox")
	check.AddClass("task-check")
	check.SetProp("event", "task.toggle."+strconv.Itoa(index))
	if t.done {
		check.SetProp("checked", "true")
	} else {
		check.SetProp("checked", "false")
	}

	label := components.Text(t.label, "task-text")
	if dark {
		label.AddClass("task-text-dark")
	}
	if t.done {
		label.AddClass("task-done")
	}

	row := components.Row(check, label)
	row.AddClass("task-row")
	if dark {
		row.AddClass("task-row-dark")
	}
	return row
}

func buildUI(v viewState) *giocss.Node {
	root := giocss.NewNode("body")
	root.AddClass("app-shell")
	if v.darkTheme {
		root.AddClass("app-shell-dark")
	}

	title := components.Heading(2, "State Mutation Lab", "title")
	subtitle := components.Text("Interactive sample: events mutate app state and Snapshot rerenders UI.", "subtitle")
	if v.darkTheme {
		title.AddClass("title-dark")
		subtitle.AddClass("subtitle-dark")
	}
	root.AddChild(title)
	root.AddChild(subtitle)

	controls := components.Row(
		components.Button("+1", "action-btn", "btn-inc"),
		components.Button("-1", "action-btn", "btn-dec"),
		components.Button("Reset", "action-btn", "btn-reset"),
		components.Button("Toggle Theme", "action-btn", "btn-theme"),
	)
	controls.AddClass("controls")
	controls.Children[0].SetProp("event", "counter.increment")
	controls.Children[1].SetProp("event", "counter.decrement")
	controls.Children[2].SetProp("event", "counter.reset")
	controls.Children[3].SetProp("event", "theme.toggle")
	root.AddChild(controls)

	stats := components.Column(
		components.Text("Counter", "counter-label"),
		components.Text(strconv.Itoa(v.counter), "counter-value"),
		components.Text("Mutations: "+strconv.Itoa(v.mutations), "mut-line"),
		components.Text("Last event: "+v.lastEvent, "mut-line"),
	)
	stats.AddClass("stats-card")
	if v.darkTheme {
		stats.AddClass("stats-card-dark")
		stats.Children[0].AddClass("counter-label-dark")
		stats.Children[1].AddClass("counter-value-dark")
		stats.Children[2].AddClass("mut-line-dark")
		stats.Children[3].AddClass("mut-line-dark")
	}
	if v.counter > 0 {
		stats.Children[1].AddClass("counter-positive")
	} else if v.counter < 0 {
		stats.Children[1].AddClass("counter-negative")
	}
	root.AddChild(stats)

	tasksTitle := components.Text("Mutable tasks", "tasks-title")
	if v.darkTheme {
		tasksTitle.AddClass("tasks-title-dark")
	}
	root.AddChild(tasksTitle)
	for i, t := range v.tasks {
		root.AddChild(buildTaskRow(i, t, v.darkTheme))
	}

	return root
}

func main() {
	ss := giocss.NewStyleSheet()
	ss.ParseCSSText(css)
	state := newAppState()

	rt := giocss.NewWindowRuntime(
		giocss.WindowOptions{Title: "Sample 08 – State Mutation", Width: 760, Height: 620},
		giocss.WindowRuntimeHooks{
			DispatchEvent: func(eventName string, payload map[string]any) error {
				state.apply(eventName, payload)
				return nil
			},
			Snapshot: func(size image.Point) giocss.WindowRuntimeSnapshot {
				root := buildUI(state.snapshot())
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

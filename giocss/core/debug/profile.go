package debug

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// anyToString is a local helper to avoid importing core/engine for a single function.
func anyToString(candidate any, fallback string) string {
	if typed, ok := candidate.(string); ok {
		normalized := strings.ToValidUTF8(typed, "")
		trimmed := strings.TrimSpace(normalized)
		if trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

type DebugComponentStat struct {
	Path          string
	Kind          string
	ID            string
	ClassName     string
	Component     string
	Tag           string
	Role          string
	StyleHint     string
	Attrs         map[string]string
	Samples       int64
	TotalMS       float64
	SelfTotalMS   float64
	MaxMS         float64
	SelfMaxMS     float64
	EstMemBytes   int64
	MaxMemBytes   int64
	LastRenderMS  float64
	LastSeenFrame int64
}

type DebugComponentUpdate struct {
	Path      string
	Kind      string
	ID        string
	ClassName string
	Component string
	Tag       string
	Role      string
	StyleHint string
	Attrs     map[string]string
	TotalMS   float64
	SelfMS    float64
	EstBytes  int64
	Frame     int64
}

type DebugComponentSampleInput struct {
	Path    string
	Kind    string
	Props   map[string]any
	Width   int
	Height  int
	TotalMS float64
	SelfMS  float64
	Frame   int64
}

func BuildDebugComponentUpdate(in DebugComponentSampleInput) (DebugComponentUpdate, bool) {
	if strings.TrimSpace(in.Path) == "" {
		return DebugComponentUpdate{}, false
	}
	w := in.Width
	if w < 1 {
		w = 1
	}
	h := in.Height
	if h < 1 {
		h = 1
	}
	s := SummarizeComponentProps(in.Kind, in.Props)
	selfMS := in.SelfMS
	if selfMS < 0 {
		selfMS = 0
	}
	return DebugComponentUpdate{
		Path:      strings.TrimSpace(in.Path),
		Kind:      in.Kind,
		ID:        s.ID,
		ClassName: s.ClassName,
		Component: s.Component,
		Tag:       s.Tag,
		Role:      s.Role,
		StyleHint: s.StyleHint,
		Attrs:     s.Attrs,
		TotalMS:   in.TotalMS,
		SelfMS:    selfMS,
		EstBytes:  int64(w) * int64(h) * 4,
		Frame:     in.Frame,
	}, true
}

func SnapshotDebugComponentStats(stats map[string]*DebugComponentStat) map[string]*DebugComponentStat {
	if len(stats) == 0 {
		return nil
	}
	out := make(map[string]*DebugComponentStat, len(stats))
	for k, v := range stats {
		if v == nil {
			continue
		}
		copyV := *v
		out[k] = &copyV
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func EnabledProfileFlags(profile map[string]bool) []string {
	if len(profile) == 0 {
		return nil
	}
	flags := make([]string, 0, len(profile))
	for k, enabled := range profile {
		if enabled {
			flags = append(flags, k)
		}
	}
	if len(flags) == 0 {
		return nil
	}
	return flags
}

func NormalizeProfileFlags(profile map[string]bool) map[string]bool {
	if len(profile) == 0 {
		return map[string]bool{}
	}
	flags := make(map[string]bool, len(profile))
	for k, v := range profile {
		flags[strings.ToLower(strings.TrimSpace(k))] = v
	}
	return flags
}

func ParseDebugProfileConfig(config map[string]any) (map[string]bool, string) {
	profile := make(map[string]bool, len(config))
	profilerPath := ""
	for key, raw := range config {
		name := strings.ToLower(strings.TrimSpace(key))
		if name == "" {
			continue
		}
		switch name {
		case "profiler_path", "profilerpath", "profile_path", "dump_path":
			if text := strings.TrimSpace(anyToString(raw, "")); text != "" {
				profilerPath = text
			}
		default:
			profile[name] = anyIsTruthy(raw)
		}
	}
	return profile, profilerPath
}

func anyIsTruthy(raw any) bool {
	switch typed := raw.(type) {
	case nil:
		return false
	case bool:
		return typed
	case int:
		return typed != 0
	case int64:
		return typed != 0
	case float64:
		return typed != 0
	case float32:
		return typed != 0
	case string:
		text := strings.ToLower(strings.TrimSpace(typed))
		return text != "" && text != "0" && text != "false" && text != "no" && text != "off"
	default:
		return true
	}
}

type ProfileDumpInput struct {
	Timestamp    string
	Frames       int64
	FPS          float64
	RenderMS     float64
	RenderAvgMS  float64
	RenderMaxMS  float64
	SlowFrames   int64
	HeapMB       float64
	GPUEnabled   bool
	LowPower     bool
	ProfileFlags []string
	WindowTitle  string
	WindowWidth  int
	WindowHeight int
	TopLimit     int
}

func BuildProfileDumpData(in ProfileDumpInput, componentStats map[string]*DebugComponentStat) map[string]any {
	timestamp := strings.TrimSpace(in.Timestamp)
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339Nano)
	}
	flags := make([]string, 0, len(in.ProfileFlags))
	for _, flag := range in.ProfileFlags {
		if strings.TrimSpace(flag) == "" {
			continue
		}
		flags = append(flags, flag)
	}
	data := map[string]any{
		"timestamp":     timestamp,
		"frames":        in.Frames,
		"fps":           in.FPS,
		"render_ms":     in.RenderMS,
		"render_avg_ms": in.RenderAvgMS,
		"render_max_ms": in.RenderMaxMS,
		"slow_frames":   in.SlowFrames,
		"heap_mb":       in.HeapMB,
		"gpu_enabled":   in.GPUEnabled,
		"low_power":     in.LowPower,
		"profile_flags": flags,
		"window_title":  in.WindowTitle,
		"window_width":  in.WindowWidth,
		"window_height": in.WindowHeight,
	}
	if len(componentStats) == 0 {
		return data
	}
	limit := in.TopLimit
	if limit <= 0 {
		limit = 12
	}
	data["components_profiled"] = len(componentStats)
	data["top_components_cpu"] = BuildTopComponentStats(componentStats, false, limit)
	data["top_components_mem"] = BuildTopComponentStats(componentStats, true, limit)
	return data
}

func UpsertDebugComponentStat(stats map[string]*DebugComponentStat, u DebugComponentUpdate) {
	if stats == nil || u.Path == "" {
		return
	}
	st := stats[u.Path]
	if st == nil {
		st = &DebugComponentStat{Path: u.Path}
		stats[u.Path] = st
	}
	st.Kind = u.Kind
	if u.ID != "" {
		st.ID = u.ID
	}
	if u.ClassName != "" {
		st.ClassName = u.ClassName
	}
	if u.Component != "" {
		st.Component = u.Component
	}
	if u.Tag != "" {
		st.Tag = u.Tag
	}
	if u.Role != "" {
		st.Role = u.Role
	}
	if u.StyleHint != "" {
		st.StyleHint = u.StyleHint
	}
	if u.Attrs != nil {
		st.Attrs = CloneStringMap(u.Attrs)
	}
	st.Samples++
	st.TotalMS += u.TotalMS
	st.SelfTotalMS += u.SelfMS
	if u.TotalMS > st.MaxMS {
		st.MaxMS = u.TotalMS
	}
	if u.SelfMS > st.SelfMaxMS {
		st.SelfMaxMS = u.SelfMS
	}
	st.LastRenderMS = u.SelfMS
	st.EstMemBytes += u.EstBytes
	if u.EstBytes > st.MaxMemBytes {
		st.MaxMemBytes = u.EstBytes
	}
	st.LastSeenFrame = u.Frame
}

func BuildTopComponentStats(stats map[string]*DebugComponentStat, byMem bool, limit int) []map[string]any {
	if len(stats) == 0 || limit <= 0 {
		return nil
	}
	items := make([]*DebugComponentStat, 0, len(stats))
	for _, st := range stats {
		if st == nil {
			continue
		}
		items = append(items, st)
	}
	if len(items) == 0 {
		return nil
	}
	sort.Slice(items, func(i, j int) bool {
		a := items[i]
		b := items[j]
		if byMem {
			if a.EstMemBytes == b.EstMemBytes {
				return a.SelfTotalMS > b.SelfTotalMS
			}
			return a.EstMemBytes > b.EstMemBytes
		}
		if a.SelfTotalMS == b.SelfTotalMS {
			return a.SelfMaxMS > b.SelfMaxMS
		}
		return a.SelfTotalMS > b.SelfTotalMS
	})
	if len(items) > limit {
		items = items[:limit]
	}
	out := make([]map[string]any, 0, len(items))
	for _, st := range items {
		avg := 0.0
		selfAvg := 0.0
		if st.Samples > 0 {
			avg = st.TotalMS / float64(st.Samples)
			selfAvg = st.SelfTotalMS / float64(st.Samples)
		}
		out = append(out, map[string]any{
			"path":              st.Path,
			"kind":              st.Kind,
			"id":                st.ID,
			"class_name":        st.ClassName,
			"component":         st.Component,
			"tag":               st.Tag,
			"role":              st.Role,
			"style_hint":        st.StyleHint,
			"attrs":             CloneStringMap(st.Attrs),
			"samples":           st.Samples,
			"total_ms":          st.TotalMS,
			"self_total_ms":     st.SelfTotalMS,
			"avg_ms":            avg,
			"self_avg_ms":       selfAvg,
			"max_ms":            st.MaxMS,
			"self_max_ms":       st.SelfMaxMS,
			"estimated_mem_mb":  float64(st.EstMemBytes) / (1024.0 * 1024.0),
			"estimated_peak_mb": float64(st.MaxMemBytes) / (1024.0 * 1024.0),
		})
	}
	return out
}

func WriteProfileDump(path string, data map[string]any) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("profiler dump mkdir failed: %w", err)
		}
	}
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("profiler dump marshal failed: %w", err)
	}
	if err := os.WriteFile(path, buf, 0o644); err != nil {
		return fmt.Errorf("profiler dump write failed: %w", err)
	}
	return nil
}

type ProfilerMetrics struct {
	Frames         int64
	FPS            float64
	RenderMS       float64
	RenderAvgMS    float64
	RenderMaxMS    float64
	SlowFrames     int64
	HeapMB         float64
	ComponentCount int
	GPUEnabled     bool
	LowPower       bool
}

type DebugFrameMetricsState struct {
	Frames      int64
	FPSFrames   int64
	FPSLastAt   time.Time
	FPS         float64
	RenderMS    float64
	RenderAvgMS float64
	RenderMaxMS float64
	SlowFrames  int64
}

func UpdateDebugFrameMetrics(state DebugFrameMetricsState, frameTime time.Time) DebugFrameMetricsState {
	state.Frames++
	state.FPSFrames++
	if state.FPSLastAt.IsZero() {
		state.FPSLastAt = frameTime
	}
	dt := frameTime.Sub(state.FPSLastAt)
	if dt >= time.Second {
		if dt > 0 {
			state.FPS = float64(state.FPSFrames) / dt.Seconds()
		}
		state.FPSFrames = 0
		state.FPSLastAt = frameTime
	}
	return state
}

func UpdateDebugRenderMetrics(state DebugFrameMetricsState, renderMS float64) DebugFrameMetricsState {
	state.RenderMS = renderMS
	if state.Frames <= 1 {
		state.RenderAvgMS = renderMS
	} else {
		n := float64(state.Frames)
		state.RenderAvgMS += (renderMS - state.RenderAvgMS) / n
	}
	if renderMS > state.RenderMaxMS {
		state.RenderMaxMS = renderMS
	}
	if renderMS > 16.7 {
		state.SlowFrames++
	}
	return state
}

type DebugOverlayInput struct {
	Path          string
	Kind          string
	NativeKind    string
	Active        bool
	Hovered       bool
	HasHandler    bool
	ClipChildren  bool
	X             int
	Y             int
	W             int
	H             int
	ContentRight  int
	ContentBottom int
}

type DebugOverlayModel struct {
	OutlineColor   color.NRGBA
	Interactive    bool
	ClipChildren   bool
	HasOverflow    bool
	Label          string
	BaseColor      color.NRGBA
	InteractiveCol color.NRGBA
	ScrollCol      color.NRGBA
}

func BuildDebugOverlayModel(in DebugOverlayInput) DebugOverlayModel {
	baseCol := color.NRGBA{R: 255, G: 173, B: 51, A: 190}
	interactiveCol := color.NRGBA{R: 255, G: 64, B: 129, A: 190}
	scrollCol := color.NRGBA{R: 0, G: 220, B: 255, A: 200}
	hoverCol := color.NRGBA{R: 140, G: 255, B: 140, A: 190}
	activeCol := color.NRGBA{R: 255, G: 90, B: 90, A: 200}

	outlineCol := baseCol
	if in.Active {
		outlineCol = activeCol
	} else if in.Hovered {
		outlineCol = hoverCol
	}

	interactive := false
	if in.Kind == "button" || in.Kind == "input" {
		interactive = true
	}
	if in.NativeKind == "slider" || in.NativeKind == "checkbox" || in.NativeKind == "radio" || in.NativeKind == "select" || in.NativeKind == "dropdown" {
		interactive = true
	}

	label := in.Kind
	if label == "native" && in.NativeKind != "" {
		label = label + ":" + in.NativeKind
	}

	return DebugOverlayModel{
		OutlineColor:   outlineCol,
		Interactive:    interactive && in.HasHandler,
		ClipChildren:   in.ClipChildren,
		HasOverflow:    in.ContentRight > in.X+in.W || in.ContentBottom > in.Y+in.H,
		Label:          label,
		BaseColor:      baseCol,
		InteractiveCol: interactiveCol,
		ScrollCol:      scrollCol,
	}
}

func BuildProfilerOverlayLines(flags map[string]bool, m ProfilerMetrics) []string {
	showFrames := flags["frames"]
	showFPS := flags["fpscounter"] || flags["fps"]
	showRender := flags["render"] || flags["render_ms"]
	showCPU := flags["cpu"]
	showGPU := flags["gpu"]
	showMem := flags["mem"] || flags["memory"]
	showSlow := flags["slow"] || flags["slowframes"]
	showComponents := flags["components"] || flags["component_cpu"] || flags["component_mem"]
	if !showFrames && !showFPS && !showRender && !showCPU && !showGPU && !showMem && !showSlow && !showComponents {
		return nil
	}
	lines := make([]string, 0, 8)
	if showFrames {
		lines = append(lines, fmt.Sprintf("frames: %d", m.Frames))
	}
	if showFPS {
		lines = append(lines, fmt.Sprintf("fps: %.1f", m.FPS))
	}
	if showRender {
		lines = append(lines, fmt.Sprintf("render ms: %.2f (avg %.2f / max %.2f)", m.RenderMS, m.RenderAvgMS, m.RenderMaxMS))
	}
	if showCPU {
		cpuApprox := (m.RenderAvgMS / 16.7) * 100.0
		if cpuApprox < 0 {
			cpuApprox = 0
		}
		lines = append(lines, fmt.Sprintf("cpu est: %.0f%% of 60fps budget", cpuApprox))
	}
	if showGPU {
		mode := "on"
		if !m.GPUEnabled {
			mode = "off"
		}
		power := "normal"
		if m.LowPower {
			power = "low"
		}
		lines = append(lines, fmt.Sprintf("gpu: %s | power: %s", mode, power))
	}
	if showMem {
		lines = append(lines, fmt.Sprintf("heap: %.1f MB", m.HeapMB))
	}
	if showSlow {
		lines = append(lines, fmt.Sprintf("slow frames (>16.7ms): %d", m.SlowFrames))
	}
	if showComponents {
		lines = append(lines, fmt.Sprintf("components tracked: %d", m.ComponentCount))
	}
	return lines
}

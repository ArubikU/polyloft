package core

import (
	"image"
	"image/color"
	stdruntime "runtime"
	"strings"
	"sync"
	"time"
)

type WindowRuntimeSnapshot struct {
	RootLayout   map[string]any
	RootCSS      map[string]string
	StyleSheet   *StyleSheet
	ScreenWidth  int
	ScreenHeight int
}

type WindowRuntimeHooks struct {
	Snapshot         func(size image.Point) WindowRuntimeSnapshot
	DispatchEvent    func(eventName string, payload map[string]any) error
	EmitRuntimeError func(err error)
	OnClose          func()
}

type WindowRuntime struct {
	mu              sync.Mutex
	win             *Window
	title           string
	width           int
	height          int
	hooks           WindowRuntimeHooks
	debug           bool
	profile         map[string]bool
	profileDumpPath string
	frameMetrics    DebugFrameMetricsState
	heapMB          float64
	lastHeapSample  time.Time
	components      map[string]*DebugComponentStat
	hoverState      map[string]bool
	activeState     map[string]bool
	renderStore     *RenderStore
}

func NewWindowRuntime(opts WindowOptions, hooks WindowRuntimeHooks) *WindowRuntime {
	return &WindowRuntime{
		win:         NewWindow(opts),
		title:       opts.Title,
		width:       opts.Width,
		height:      opts.Height,
		hooks:       hooks,
		profile:     make(map[string]bool),
		components:  make(map[string]*DebugComponentStat),
		hoverState:  make(map[string]bool),
		activeState: make(map[string]bool),
		renderStore: NewRenderStore(),
	}
}

func (rt *WindowRuntime) Window() *Window {
	if rt == nil {
		return nil
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	return rt.win
}

func (rt *WindowRuntime) Invalidate() {
	if win := rt.Window(); win != nil {
		win.Invalidate()
	}
}

func (rt *WindowRuntime) Close() {
	if win := rt.Window(); win != nil {
		win.Close()
	}
}

func (rt *WindowRuntime) SetTitle(title string) {
	rt.mu.Lock()
	rt.title = title
	rt.mu.Unlock()
	if win := rt.Window(); win != nil {
		win.SetTitle(title)
	}
}

func (rt *WindowRuntime) SetDebug(enabled bool) {
	if rt == nil {
		return
	}
	rt.mu.Lock()
	rt.debug = enabled
	rt.mu.Unlock()
}

func (rt *WindowRuntime) SetDebugProfile(profile map[string]bool, profilerPath string) {
	if rt == nil {
		return
	}
	rt.mu.Lock()
	rt.profile = NormalizeProfileFlags(profile)
	rt.profileDumpPath = strings.TrimSpace(profilerPath)
	rt.mu.Unlock()
}

func (rt *WindowRuntime) emitRuntimeError(err error) {
	if err == nil || rt == nil || rt.hooks.EmitRuntimeError == nil {
		return
	}
	rt.hooks.EmitRuntimeError(err)
}

func (rt *WindowRuntime) dumpProfile() {
	if rt == nil {
		return
	}

	rt.mu.Lock()
	path := strings.TrimSpace(rt.profileDumpPath)
	if path == "" {
		rt.mu.Unlock()
		return
	}
	flags := EnabledProfileFlags(rt.profile)
	componentStats := SnapshotDebugComponentStats(rt.components)
	metrics := rt.frameMetrics
	heapMB := rt.heapMB
	windowTitle := rt.title
	windowWidth := rt.width
	windowHeight := rt.height
	rt.mu.Unlock()

	data := BuildProfileDumpData(ProfileDumpInput{
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		Frames:       metrics.Frames,
		FPS:          metrics.FPS,
		RenderMS:     metrics.RenderMS,
		RenderAvgMS:  metrics.RenderAvgMS,
		RenderMaxMS:  metrics.RenderMaxMS,
		SlowFrames:   metrics.SlowFrames,
		HeapMB:       heapMB,
		GPUEnabled:   RenderGPUEnabled(),
		LowPower:     RenderLowPowerEnabled(),
		ProfileFlags: flags,
		WindowTitle:  windowTitle,
		WindowWidth:  windowWidth,
		WindowHeight: windowHeight,
		TopLimit:     12,
	}, componentStats)
	if err := WriteProfileDump(path, data); err != nil {
		rt.emitRuntimeError(err)
	}
}

func (rt *WindowRuntime) buildProfilerOverlay(rctx *RenderContext, state *GioWindowState, screenW int, screenH int) {
	if rt == nil {
		return
	}

	rt.mu.Lock()
	flags := NormalizeProfileFlags(rt.profile)
	metrics := ProfilerMetrics{
		Frames:         rt.frameMetrics.Frames,
		FPS:            rt.frameMetrics.FPS,
		RenderMS:       rt.frameMetrics.RenderMS,
		RenderAvgMS:    rt.frameMetrics.RenderAvgMS,
		RenderMaxMS:    rt.frameMetrics.RenderMaxMS,
		SlowFrames:     rt.frameMetrics.SlowFrames,
		HeapMB:         rt.heapMB,
		ComponentCount: len(rt.components),
		GPUEnabled:     RenderGPUEnabled(),
		LowPower:       RenderLowPowerEnabled(),
	}
	rt.mu.Unlock()

	lines := BuildProfilerOverlayLines(flags, metrics)
	if len(lines) == 0 {
		return
	}

	pad := 8
	lineH := 16
	boxW := 340
	boxH := pad*2 + lineH*len(lines)
	boxX := max(8, screenW-boxW-12)
	boxY := 12

	DrawGioRRect(rctx, boxX, boxY, boxW, boxH, 8, color.NRGBA{R: 12, G: 18, B: 28, A: 220})
	DrawGioBorder(rctx, boxX, boxY, boxW, boxH, 8, 1, color.NRGBA{R: 71, G: 85, B: 105, A: 230})

	for i, line := range lines {
		ty := boxY + pad + i*lineH
		DrawGioText(rctx, state, boxX+10, ty, boxW-20, lineH, line, color.NRGBA{R: 226, G: 232, B: 240, A: 255}, 12, false, false, false, CSSTextAlign(nil), 1, false)
	}
	_ = screenH
}

func (rt *WindowRuntime) Run() {
	if rt == nil {
		return
	}
	win := rt.Window()
	if win == nil {
		return
	}

	InitRenderConfig()
	RunWindowRenderer(win, WindowRunnerHooks{
		OnDestroy: func() {
			rt.dumpProfile()
			if rt.hooks.OnClose != nil {
				rt.hooks.OnClose()
			}
			rt.mu.Lock()
			rt.win = nil
			rt.mu.Unlock()
		},
		BeforeFrame: func(size image.Point, frameTime time.Time) WindowFrameSnapshot {
			rt.mu.Lock()
			if size.X > 0 {
				rt.width = size.X
			}
			if size.Y > 0 {
				rt.height = size.Y
			}
			rt.frameMetrics = UpdateDebugFrameMetrics(rt.frameMetrics, frameTime)
			if rt.hoverState == nil {
				rt.hoverState = make(map[string]bool)
			}
			if rt.activeState == nil {
				rt.activeState = make(map[string]bool)
			}
			debug := rt.debug
			profile := NormalizeProfileFlags(rt.profile)
			hoverState := rt.hoverState
			activeState := rt.activeState
			renderStore := rt.renderStore
			frameNumber := rt.frameMetrics.Frames
			rt.mu.Unlock()

			snapshot := WindowRuntimeSnapshot{}
			if rt.hooks.Snapshot != nil {
				snapshot = rt.hooks.Snapshot(size)
			}
			screenW := snapshot.ScreenWidth
			screenH := snapshot.ScreenHeight
			if screenW <= 0 {
				screenW = size.X
			}
			if screenH <= 0 {
				screenH = size.Y
			}

			return WindowFrameSnapshot{
				HoverState:          hoverState,
				ActiveState:         activeState,
				StyleSheet:          snapshot.StyleSheet,
				RenderStore:         renderStore,
				CanUseWindowEffects: true,
				FrameNumber:         frameNumber,
				Debug:               debug,
				ProfileComponents:   profile["components"] || profile["component_cpu"] || profile["component_mem"],
				ProfileFull:         profile["components_full"],
				RootLayout:          snapshot.RootLayout,
				RootCSS:             snapshot.RootCSS,
				ScreenW:             screenW,
				ScreenH:             screenH,
			}
		},
		AfterProcess: func(hoverState map[string]bool, activeState map[string]bool) {
			rt.mu.Lock()
			rt.hoverState = hoverState
			rt.activeState = activeState
			rt.mu.Unlock()
		},
		AfterRender: func(renderStore *RenderStore, renderDur time.Duration, frameTime time.Time) {
			rt.mu.Lock()
			rt.frameMetrics = UpdateDebugRenderMetrics(rt.frameMetrics, float64(renderDur)/float64(time.Millisecond))
			if rt.lastHeapSample.IsZero() || frameTime.Sub(rt.lastHeapSample) >= time.Second {
				var mem stdruntime.MemStats
				stdruntime.ReadMemStats(&mem)
				rt.heapMB = float64(mem.Alloc) / (1024.0 * 1024.0)
				rt.lastHeapSample = frameTime
			}
			rt.renderStore = renderStore
			if rt.renderStore != nil {
				rt.renderStore.Finalize()
			}
			rt.mu.Unlock()
		},
		DrawOverlay: rt.buildProfilerOverlay,
		DispatchEvent: func(eventName string, payload map[string]any) error {
			if rt.hooks.DispatchEvent == nil {
				return nil
			}
			return rt.hooks.DispatchEvent(eventName, payload)
		},
		EmitRuntimeError: rt.emitRuntimeError,
		RecordDebugSample: func(path string, kind string, props map[string]any, w int, h int, totalDur time.Duration, selfDur time.Duration) {
			rt.mu.Lock()
			frame := rt.frameMetrics.Frames
			rt.mu.Unlock()
			update, ok := BuildDebugComponentUpdate(DebugComponentSampleInput{
				Path:    path,
				Kind:    kind,
				Props:   props,
				Width:   w,
				Height:  h,
				TotalMS: float64(totalDur) / float64(time.Millisecond),
				SelfMS:  float64(selfDur) / float64(time.Millisecond),
				Frame:   frame,
			})
			if !ok {
				return
			}
			rt.mu.Lock()
			if rt.components == nil {
				rt.components = make(map[string]*DebugComponentStat)
			}
			UpsertDebugComponentStat(rt.components, update)
			rt.mu.Unlock()
		},
	})
}


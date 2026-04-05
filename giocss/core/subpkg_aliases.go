package core

// subpkg_aliases.go â€” bridges core sub-packages into package core so that
// window_renderer.go, window_runtime.go, and exports.go can use bare names.

import (
	"image"
	"image/color"
	"time"

	"gioui.org/text"

	coredebug "github.com/ArubikU/giocss/core/debug"
	coreengine "github.com/ArubikU/giocss/core/engine"
	coreinput "github.com/ArubikU/giocss/core/input"
	corelayout "github.com/ArubikU/giocss/core/layout"
	corerender "github.com/ArubikU/giocss/core/render"
)

// â”€â”€ Layout â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type Node = corelayout.Node

func NewNode(tag string) *Node { return corelayout.NewNode(tag) }

func ReconcileTrees(oldRoot *Node, newRoot *Node) ([]map[string]any, []map[string]any) {
	return corelayout.ReconcileTrees(oldRoot, newRoot)
}

func LayoutNodeToNative(node *Node, width int, height int, ss *StyleSheet) map[string]any {
	return corelayout.LayoutNodeToNative(node, width, height, ss)
}

func ResolveNodeStyle(node *Node, ss *StyleSheet, viewportW int) map[string]string {
	return corelayout.ResolveNodeStyle(node, ss, viewportW)
}

// â”€â”€ Engine â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type Window = coreengine.Window
type WindowOptions = coreengine.WindowOptions
type RenderContext = coreengine.RenderContext
type GioWindowLoopHooks = coreengine.GioWindowLoopHooks

func NewWindow(opts WindowOptions) *Window              { return coreengine.NewWindow(opts) }
func RunApp()                                           { coreengine.RunApp() }
func RunGioWindowLoop(gw *Window, h GioWindowLoopHooks) { coreengine.RunGioWindowLoop(gw, h) }

// â”€â”€ Debug â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type DebugFrameMetricsState = coredebug.DebugFrameMetricsState
type DebugComponentStat = coredebug.DebugComponentStat
type DebugComponentUpdate = coredebug.DebugComponentUpdate
type DebugComponentSampleInput = coredebug.DebugComponentSampleInput
type ProfilerMetrics = coredebug.ProfilerMetrics
type ProfileDumpInput = coredebug.ProfileDumpInput

func ParseDebugProfileConfig(config map[string]any) (map[string]bool, string) {
	return coredebug.ParseDebugProfileConfig(config)
}
func NormalizeProfileFlags(profile map[string]bool) map[string]bool {
	return coredebug.NormalizeProfileFlags(profile)
}
func EnabledProfileFlags(profile map[string]bool) []string {
	return coredebug.EnabledProfileFlags(profile)
}
func SnapshotDebugComponentStats(stats map[string]*DebugComponentStat) map[string]*DebugComponentStat {
	return coredebug.SnapshotDebugComponentStats(stats)
}
func BuildProfileDumpData(in ProfileDumpInput, componentStats map[string]*DebugComponentStat) map[string]any {
	return coredebug.BuildProfileDumpData(in, componentStats)
}
func WriteProfileDump(path string, data map[string]any) error {
	return coredebug.WriteProfileDump(path, data)
}
func BuildProfilerOverlayLines(flags map[string]bool, metrics ProfilerMetrics) []string {
	return coredebug.BuildProfilerOverlayLines(flags, metrics)
}
func UpdateDebugFrameMetrics(state DebugFrameMetricsState, frameTime time.Time) DebugFrameMetricsState {
	return coredebug.UpdateDebugFrameMetrics(state, frameTime)
}
func UpdateDebugRenderMetrics(state DebugFrameMetricsState, renderMS float64) DebugFrameMetricsState {
	return coredebug.UpdateDebugRenderMetrics(state, renderMS)
}
func BuildDebugComponentUpdate(in DebugComponentSampleInput) (DebugComponentUpdate, bool) {
	return coredebug.BuildDebugComponentUpdate(in)
}
func UpsertDebugComponentStat(stats map[string]*DebugComponentStat, u DebugComponentUpdate) {
	coredebug.UpsertDebugComponentStat(stats, u)
}

// â”€â”€ Input â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func InputExternalValue(props map[string]any) string { return coreinput.InputExternalValue(props) }
func NormalizeNumberInput(raw string, props map[string]any, finalize bool) string {
	return coreinput.NormalizeNumberInput(raw, props, finalize)
}
func NormalizeDateInput(raw string, finalize bool) string {
	return coreinput.NormalizeDateInput(raw, finalize)
}
func NormalizeTimeInput(raw string, props map[string]any, finalize bool) string {
	return coreinput.NormalizeTimeInput(raw, props, finalize)
}

// â”€â”€ Render â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type GioWindowState = corerender.GioWindowState
type GioRenderHost = corerender.GioRenderHost

func NewGioWindowState() *GioWindowState { return corerender.NewGioWindowState() }
func InitRenderConfig()                  { corerender.InitRenderConfig() }
func RenderGPUEnabled() bool             { return corerender.RenderGPUEnabled() }
func RenderLowPowerEnabled() bool        { return corerender.RenderLowPowerEnabled() }

func ProcessGioEvents(rctx *RenderContext, host *GioRenderHost, state *GioWindowState, frameTime time.Time) bool {
	return corerender.ProcessGioEvents(rctx, host, state, frameTime)
}
func DrawGioBackground(rctx *RenderContext, x, y, w, h int, css map[string]string) {
	corerender.DrawGioBackground(rctx, x, y, w, h, css)
}
func DrawGioTree(rctx *RenderContext, node map[string]any, host *GioRenderHost, path string, state *GioWindowState, parentClass string, containingBlock image.Rectangle, parentCSS map[string]string, inheritedShiftX int, inheritedShiftY int, frameTime time.Time) {
	corerender.DrawGioTree(rctx, node, host, path, state, parentClass, containingBlock, parentCSS, inheritedShiftX, inheritedShiftY, frameTime)
}
func DrawPickerModal(rctx *RenderContext, state *GioWindowState, screenW, screenH int) {
	corerender.DrawPickerModal(rctx, state, screenW, screenH)
}
func ApplyFrameCursorOverride(rctx *RenderContext, state *GioWindowState) {
	corerender.ApplyFrameCursorOverride(rctx, state)
}
func DrawGioRRect(rctx *RenderContext, x, y, w, h, radius int, col color.NRGBA) {
	corerender.DrawGioRRect(rctx, x, y, w, h, radius, col)
}
func DrawGioBorder(rctx *RenderContext, x, y, w, h, radius, borderW int, col color.NRGBA) {
	corerender.DrawGioBorder(rctx, x, y, w, h, radius, borderW, col)
}
func DrawGioText(rctx *RenderContext, state *GioWindowState, x, y, w, h int, txt string, col color.NRGBA, fontSize float32, bold, italic, mono bool, align text.Alignment, maxLines int, wrapText bool) {
	corerender.DrawGioText(rctx, state, x, y, w, h, txt, col, fontSize, bold, italic, mono, align, maxLines, wrapText)
}

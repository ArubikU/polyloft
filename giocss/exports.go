package giocss

import (
	"image/color"

	core "github.com/ArubikU/giocss/core"
)

// Temporary compatibility facade. New code should import domain subpackages directly.
type Node = core.Node

type StyleSheet = core.StyleSheet

type WindowOptions = core.WindowOptions

type WindowRuntimeSnapshot = core.WindowRuntimeSnapshot

type WindowRuntimeHooks = core.WindowRuntimeHooks

type WindowRuntime = core.WindowRuntime

func NewNode(tag string) *Node { return core.NewNode(tag) }

func ReconcileTrees(oldRoot *Node, newRoot *Node) ([]map[string]any, []map[string]any) {
	return core.ReconcileTrees(oldRoot, newRoot)
}

func LayoutNodeToNative(node *Node, width int, height int, ss *StyleSheet) map[string]any {
	return core.LayoutNodeToNative(node, width, height, ss)
}

func ResolveNodeStyle(node *Node, ss *StyleSheet, viewportW int) map[string]string {
	return core.ResolveNodeStyle(node, ss, viewportW)
}

func NewStyleSheet() *StyleSheet { return core.NewStyleSheet() }

func ResolveStyle(props map[string]any, appSS *StyleSheet, viewportW int) map[string]string {
	return core.ResolveStyle(props, appSS, viewportW)
}

func NewWindowRuntime(opts WindowOptions, hooks WindowRuntimeHooks) *WindowRuntime {
	return core.NewWindowRuntime(opts, hooks)
}

func ParseDebugProfileConfig(config map[string]any) (map[string]bool, string) {
	return core.ParseDebugProfileConfig(config)
}

func RunApp() { core.RunApp() }

func ParseHexColor(s string, fallback color.Color) color.Color {
	return core.ParseHexColor(s, fallback)
}

func ParseRGBColor(s string, fallback color.Color) color.Color {
	return core.ParseRGBColor(s, fallback)
}

func ParseHSLColor(s string, fallback color.Color) color.Color {
	return core.ParseHSLColor(s, fallback)
}

func CSSTextTransform(s string, css map[string]string) string {
	return core.CSSTextTransform(s, css)
}

func InputExternalValue(props map[string]any) string { return core.InputExternalValue(props) }

func NormalizeNumberInput(raw string, props map[string]any, finalize bool) string {
	return core.NormalizeNumberInput(raw, props, finalize)
}

func NormalizeDateInput(raw string, finalize bool) string {
	return core.NormalizeDateInput(raw, finalize)
}

func NormalizeTimeInput(raw string, props map[string]any, finalize bool) string {
	return core.NormalizeTimeInput(raw, props, finalize)
}

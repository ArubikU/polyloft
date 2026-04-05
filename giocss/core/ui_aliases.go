package core

// ui_aliases.go â€” bridge from core to the /ui and /pkg packages.
// All types and functions previously scattered across thin-wrapper files
// (render_gio_state.go, render_store.go, render_font.go, style_color_helpers.go,
// style_text_helpers.go, style_transition.go) are consolidated here.

import (
	"image/color"
	"time"

	"gioui.org/text"

	pkgcore "github.com/ArubikU/giocss/pkg"
	uicore "github.com/ArubikU/giocss/ui"
)

// â”€â”€ Render / GIO state â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type ResolvedCSSCacheEntry = uicore.ResolvedCSSCacheEntry
type ZChildrenHintCacheEntry = uicore.ZChildrenHintCacheEntry

var RenderConfigOnce = &uicore.RenderConfigOnce

func NewGioShaper() *text.Shaper {
	return uicore.NewGioShaper()
}

func PurgeStaleRenderCaches(resolvedCSS map[string]ResolvedCSSCacheEntry, zChildrenHint map[string]ZChildrenHintCacheEntry, frame int64) {
	uicore.PurgeStaleRenderCaches(resolvedCSS, zChildrenHint, frame)
}

// â”€â”€ Render store â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type RenderSnapshot = uicore.RenderSnapshot
type RenderStore = uicore.RenderStore

func NewRenderSnapshot(css map[string]string, x int, y int, w int, h int) RenderSnapshot {
	return uicore.NewRenderSnapshot(css, x, y, w, h)
}

func NewRenderStore() *RenderStore {
	return uicore.NewRenderStore()
}

// â”€â”€ Font utilities â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func LooksLikeFontBinary(data []byte) bool {
	return uicore.LooksLikeFontBinary(data)
}

func FontHasSpaceGlyph(data []byte) bool {
	return uicore.FontHasSpaceGlyph(data)
}

func LoadFontResource(pathOrURL string) ([]byte, error) {
	return uicore.LoadFontResource(pathOrURL)
}

// â”€â”€ Color helpers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func ParseHexColor(input string, fallback color.Color) color.Color {
	return uicore.ParseHexColor(input, fallback)
}

func SplitColorArgs(input string) []string {
	return uicore.SplitColorArgs(input)
}

func ParseRGBColor(input string, fallback color.Color) color.Color {
	return uicore.ParseRGBColor(input, fallback)
}

func ParseHSLColor(input string, fallback color.Color) color.Color {
	return uicore.ParseHSLColor(input, fallback)
}

func ParseCMYKColor(input string, fallback color.Color) color.Color {
	return uicore.ParseCMYKColor(input, fallback)
}

func ParseNamedColor(input string) color.Color {
	return uicore.ParseNamedColor(input)
}

func ApplyCSSOpacity(c color.Color, css map[string]string) color.Color {
	return uicore.ApplyCSSOpacity(c, css)
}

// â”€â”€ Text helpers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func CSSTextTransform(text string, css map[string]string) string {
	return uicore.CSSTextTransform(text, css)
}

func CSSApplyLetterSpacing(text string, css map[string]string) string {
	return uicore.CSSApplyLetterSpacing(text, css)
}

func CSSLineHeightPx(css map[string]string, fontSize float32) int {
	return uicore.CSSLineHeightPx(css, fontSize)
}

func EstimateTextBox(text string, fontSize float32, css map[string]string) (int, int) {
	return uicore.EstimateTextBox(text, fontSize, css)
}

func EstimateTextLayoutBox(text string, fontSize float32, css map[string]string) (int, int) {
	return uicore.EstimateTextLayoutBox(text, fontSize, css)
}

func FitTextToWidth(text string, fontSize float32, css map[string]string, width int) string {
	return uicore.FitTextToWidth(text, fontSize, css, width)
}

func TextCharsForWidth(fontSize float32, width int) int {
	return uicore.TextCharsForWidth(fontSize, width)
}

func WrapTextToWidth(text string, fontSize float32, width int) string {
	return uicore.WrapTextToWidth(text, fontSize, width)
}

func SanitizeRenderText(text string) string {
	return uicore.SanitizeRenderText(text)
}

func CSSAllowsWrap(css map[string]string) bool {
	return uicore.CSSAllowsWrap(css)
}

// â”€â”€ Transition / transform / filter helpers â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type PositionInfo = uicore.PositionInfo
type ShadowParams = uicore.ShadowParams
type TextDecorationInfo = uicore.TextDecorationInfo

func SplitCommaOutsideParens(input string) []string {
	return uicore.SplitCommaOutsideParens(input)
}

func ParseTransition(css map[string]string) (time.Duration, []string) {
	return uicore.ParseTransition(css)
}

func HasTransitionProp(props []string, name string) bool {
	return uicore.HasTransitionProp(props, name)
}

func CSSScale(css map[string]string) float64 {
	return uicore.CSSScale(css)
}

func ParseTransformTranslate(css map[string]string) (translateX, translateY float64) {
	return uicore.ParseTransformTranslate(css)
}

func ParseFilterBrightness(css map[string]string) float64 {
	return uicore.ParseFilterBrightness(css)
}

func ParseFilterContrast(css map[string]string) float64 {
	return uicore.ParseFilterContrast(css)
}

func ParseFilterSaturate(css map[string]string) float64 {
	return uicore.ParseFilterSaturate(css)
}

func ParseFilterGrayscale(css map[string]string) float64 {
	return uicore.ParseFilterGrayscale(css)
}

func ParseFilterInvert(css map[string]string) float64 {
	return uicore.ParseFilterInvert(css)
}

func ParsePositionAndOffset(css map[string]string, elemW, elemH int) PositionInfo {
	return uicore.ParsePositionAndOffset(css, elemW, elemH)
}

func ParseZIndex(css map[string]string) int {
	return uicore.ParseZIndex(css)
}

func ParseWidthHeightCSS(css map[string]string, elemW, elemH int) (int, int, bool, bool) {
	return uicore.ParseWidthHeightCSS(css, elemW, elemH)
}

func ParseFilterDropShadow(css map[string]string) ShadowParams {
	return uicore.ParseFilterDropShadow(css)
}

func ParseTextDecoration(css map[string]string) TextDecorationInfo {
	return uicore.ParseTextDecoration(css)
}

func ParseCursor(css map[string]string) string {
	return uicore.ParseCursor(css)
}

func ParseAspectRatio(css map[string]string) (float64, bool) {
	return uicore.ParseAspectRatio(css)
}

func MixColor(from color.Color, to color.Color, t float64) color.NRGBA {
	return uicore.MixColor(from, to, t)
}

func ClampUint8(v float64) uint8 {
	return uicore.ClampUint8(v)
}

func ApplyBrightnessToColor(col color.NRGBA, brightness float64) color.NRGBA {
	return uicore.ApplyBrightnessToColor(col, brightness)
}

func ApplyContrastToColor(col color.NRGBA, contrast float64) color.NRGBA {
	return uicore.ApplyContrastToColor(col, contrast)
}

func ApplySaturationToColor(col color.NRGBA, saturation float64) color.NRGBA {
	return uicore.ApplySaturationToColor(col, saturation)
}

func ApplyGrayscaleToColor(col color.NRGBA, grayscale float64) color.NRGBA {
	return uicore.ApplyGrayscaleToColor(col, grayscale)
}

func ApplyInvertToColor(col color.NRGBA, invert float64) color.NRGBA {
	return uicore.ApplyInvertToColor(col, invert)
}

// â”€â”€ Stylesheet (StyleSheet type now lives in /ui) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

type StyleSheet = uicore.StyleSheet

func NewStyleSheet() *StyleSheet {
	return uicore.NewStyleSheet()
}

func ResolveStyle(props map[string]any, appSS *StyleSheet, viewportW int) map[string]string {
	return uicore.ResolveStyle(props, appSS, viewportW)
}

func CSSResolveVariables(css map[string]string) {
	uicore.CSSResolveVariables(css)
}

func CanonicalName(name string) string {
	return uicore.CanonicalName(name)
}

func CSSExpandBoxShorthand(css map[string]string, prop string) {
	uicore.CSSExpandBoxShorthand(css, prop)
}

func CSSExpandBorderShorthand(css map[string]string) {
	uicore.CSSExpandBorderShorthand(css)
}

func CSSGetColor(css map[string]string, prop, fallback string) string {
	return uicore.CSSGetColor(css, prop, fallback)
}

func CSSBackground(css map[string]string) string {
	return uicore.CSSBackground(css)
}

func CSSFontSize(css map[string]string, fallback float32) float32 {
	return uicore.CSSFontSize(css, fallback)
}

func CSSLength(css map[string]string, prop string, fallback int) int {
	return uicore.CSSLength(css, prop, fallback)
}

func CSSValueInt(css map[string]string, prop string, fallback int) int {
	return uicore.CSSValueInt(css, prop, fallback)
}

func CSSLengthValue(raw string, fallback int, basis int, viewportW int, viewportH int) int {
	return uicore.CSSLengthValue(raw, fallback, basis, viewportW, viewportH)
}

func CSSFloatValue(raw string, fallback float64) float64 {
	return uicore.CSSFloatValue(raw, fallback)
}

func CSSBold(css map[string]string) bool {
	return uicore.CSSBold(css)
}

func CSSItalic(css map[string]string) bool {
	return uicore.CSSItalic(css)
}

// â”€â”€ State / selector helpers (previously in style/state_helpers.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func UIClassTokens(raw string) []string {
	return uicore.UIClassTokens(raw)
}

func ApplyHoverStyles(css map[string]string, props map[string]any, ss *StyleSheet, hovered bool) {
	uicore.ApplyHoverStyles(css, props, ss, hovered)
}

func ApplyActiveStyles(css map[string]string, props map[string]any, ss *StyleSheet, active bool) {
	uicore.ApplyActiveStyles(css, props, ss, active)
}

func HasHoverSelector(props map[string]any, ss *StyleSheet) bool {
	return uicore.HasHoverSelector(props, ss)
}

func HasActiveSelector(props map[string]any, ss *StyleSheet) bool {
	return uicore.HasActiveSelector(props, ss)
}

// â”€â”€ Common / misc utilities (previously in common/misc_utils.go) â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€

func CSSGridSpan(value string) int {
	return pkgcore.CSSGridSpan(value)
}

func ShouldClipText(css map[string]string) bool {
	return pkgcore.ShouldClipText(css)
}

func CSSTextAlign(css map[string]string) text.Alignment {
	return pkgcore.CSSTextAlign(css)
}

func MathAbs(v float64) float64 {
	return pkgcore.MathAbs(v)
}

func AnimateFrames(duration time.Duration, frame func(t float64)) {
	pkgcore.AnimateFrames(duration, frame)
}

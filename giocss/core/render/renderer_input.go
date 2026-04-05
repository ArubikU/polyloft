package render

import (
	"strconv"
	"strings"

	uicore "github.com/ArubikU/giocss/ui"
)

type InputBoxMetrics struct {
	FontSize    float32
	PadL        int
	PadR        int
	PadT        int
	PadB        int
	SpinnerW    int
	PickerW     int
	ControlW    int
	ContentPadR int
}

type InputEditorConfig struct {
	SingleLine       bool
	UseHeuristicWrap bool
	MaxLen           int
	Mask             rune
}

type SingleLineInputLayoutPlan struct {
	ScrollX      int
	FullWidth    int
	ContentWidth int
}

type InputScrollIndicatorPlan struct {
	Visible  bool
	TrackX   int
	TrackY   int
	TrackW   int
	TrackH   int
	ThumbX   int
	ThumbW   int
	ThumbMin int
}

type SpinnerLayoutPlan struct {
	Visible bool
	X       int
	Y       int
	W       int
	H       int
}

type PickerLayoutPlan struct {
	Visible bool
	X       int
	Y       int
	W       int
	H       int
}

type ProgressVisualPlan struct {
	Value float64
	FillW int
}

type BinaryControlLayout struct {
	BoxSize int
	BoxX    int
	BoxY    int
	Inner   int
}

type ColorSwatchLayout struct {
	X      int
	Y      int
	W      int
	H      int
	Radius int
}

type TextareaLayoutPlan struct {
	EditorW int
	EditorH int
	OffsetX int
	OffsetY int
}

func BuildInputEditorConfig(inputType string) InputEditorConfig {
	isTextarea := inputType == "textarea"
	cfg := InputEditorConfig{
		SingleLine:       !isTextarea,
		UseHeuristicWrap: isTextarea,
		MaxLen:           0,
		Mask:             0,
	}
	if inputType == "date" {
		cfg.MaxLen = 10
	} else if inputType == "time" {
		cfg.MaxLen = 8
	}
	if inputType == "password" {
		cfg.Mask = '\u2022'
	}
	return cfg
}

func ShouldSyncExternalInput(inputExternal map[string]string, path, externalVal string) bool {
	prevExternal, seen := inputExternal[path]
	return !seen || externalVal != prevExternal
}

func BuildSingleLineInputLayoutPlan(text string, fontSize float32, css map[string]string, caretCol int, visibleW int, currentScrollX int) SingleLineInputLayoutPlan {
	if visibleW < 1 {
		visibleW = 1
	}
	fullW, _ := uicore.EstimateTextBox(text, fontSize, css)
	if fullW < visibleW {
		fullW = visibleW
	}

	runes := []rune(text)
	if caretCol < 0 {
		caretCol = 0
	}
	if caretCol > len(runes) {
		caretCol = len(runes)
	}
	caretPrefix := string(runes[:caretCol])
	caretW, _ := uicore.EstimateTextBox(caretPrefix, fontSize, css)

	sx := currentScrollX
	rightPad := 8
	if caretW-sx > visibleW-rightPad {
		sx = caretW - (visibleW - rightPad)
	}
	if caretW-sx < 0 {
		sx = caretW - 2
	}

	maxSX := riMaxInt(0, fullW-visibleW)
	if sx < 0 {
		sx = 0
	}
	if sx > maxSX {
		sx = maxSX
	}

	return SingleLineInputLayoutPlan{
		ScrollX:      sx,
		FullWidth:    fullW,
		ContentWidth: riMaxInt(visibleW, fullW+2),
	}
}

func BuildInputScrollIndicatorPlan(x, y, h, padL, padB, visibleW, totalW, scrollX int) InputScrollIndicatorPlan {
	trackW := riMaxInt(1, visibleW)
	trackH := 2
	plan := InputScrollIndicatorPlan{
		Visible:  totalW > visibleW,
		TrackX:   x + padL,
		TrackW:   trackW,
		TrackH:   trackH,
		TrackY:   y + h - riMaxInt(2, padB/2) - trackH,
		ThumbMin: 12,
	}
	if !plan.Visible {
		return plan
	}
	thumbW := riMaxInt(plan.ThumbMin, int(float64(trackW)*float64(trackW)/float64(totalW)))
	maxThumbX := riMaxInt(0, trackW-thumbW)
	thumbX := 0
	if totalW > visibleW {
		thumbX = int(float64(maxThumbX) * (float64(scrollX) / float64(totalW-visibleW)))
	}
	if thumbX < 0 {
		thumbX = 0
	}
	if thumbX > maxThumbX {
		thumbX = maxThumbX
	}
	plan.ThumbX = thumbX
	plan.ThumbW = thumbW
	return plan
}

func BuildSpinnerLayoutPlan(x, y, w, h, spinnerW int, focused bool) SpinnerLayoutPlan {
	return SpinnerLayoutPlan{
		Visible: spinnerW > 0 && focused,
		X:       x + w - spinnerW,
		Y:       y,
		W:       spinnerW,
		H:       h,
	}
}

func BuildPickerLayoutPlan(x, y, w, h, pickerW int, direction string, focused bool) PickerLayoutPlan {
	return PickerLayoutPlan{
		Visible: pickerW > 0 && focused,
		X:       ResolvePickerX(x, w, pickerW, direction),
		Y:       y,
		W:       pickerW,
		H:       h,
	}
}

func BuildProgressVisualPlan(rawValue any, width int, intsArePercent bool) ProgressVisualPlan {
	value := 0.0
	switch v := rawValue.(type) {
	case float64:
		value = v
	case float32:
		value = float64(v)
	case int:
		value = float64(v)
		if intsArePercent {
			value = value / 100.0
		}
	case int64:
		value = float64(v)
		if intsArePercent {
			value = value / 100.0
		}
	}
	if value < 0 {
		value = 0
	}
	if value > 1 {
		value = 1
	}
	return ProgressVisualPlan{Value: value, FillW: int(float64(width) * value)}
}

func BuildBinaryControlLayout(x, y, w, h, inset int) BinaryControlLayout {
	boxSize := riMinInt(w, h) - inset
	if boxSize < 1 {
		boxSize = 1
	}
	return BinaryControlLayout{
		BoxSize: boxSize,
		BoxX:    x + (w-boxSize)/2,
		BoxY:    y + (h-boxSize)/2,
		Inner:   boxSize / 2,
	}
}

func BuildColorSwatchLayout(x, y, w, h, inset, minSize, radius int) ColorSwatchLayout {
	sw := w - inset*2
	sh := h - inset*2
	if sw < minSize {
		sw = minSize
	}
	if sh < minSize {
		sh = minSize
	}
	return ColorSwatchLayout{X: x + inset, Y: y + inset, W: sw, H: sh, Radius: radius}
}

func BuildTextareaLayoutPlan(x, y, w, h, insetX, insetY int) TextareaLayoutPlan {
	return TextareaLayoutPlan{
		EditorW: riMaxInt(1, w-insetX*2),
		EditorH: riMaxInt(1, h-insetY*2),
		OffsetX: x + insetX,
		OffsetY: y + insetY,
	}
}

func BuildInputBoxMetrics(css map[string]string, inputType string, isTextarea bool, w, h int, padX, padY int, fallbackFont float32) InputBoxMetrics {
	fontSize := uicore.CSSFontSize(css, fallbackFont)
	padL := maxInt(0, uicore.CSSLengthValue(css["padding-left"], padX, w, w, h))
	padR := maxInt(0, uicore.CSSLengthValue(css["padding-right"], padX, w, w, h))
	padT := maxInt(0, uicore.CSSLengthValue(css["padding-top"], padY, h, w, h))
	padB := maxInt(0, uicore.CSSLengthValue(css["padding-bottom"], padY, h, w, h))
	spinnerW := 0
	if inputType == "number" && !isTextarea {
		spinnerW = riMinInt(24, riMaxInt(16, w/5))
	}
	pickerW := 0
	if (inputType == "date" || inputType == "time") && !isTextarea {
		pickerW = riMinInt(24, riMaxInt(16, w/5))
	}
	controlW := riMaxInt(spinnerW, pickerW)
	return InputBoxMetrics{
		FontSize:    fontSize,
		PadL:        padL,
		PadR:        padR,
		PadT:        padT,
		PadB:        padB,
		SpinnerW:    spinnerW,
		PickerW:     pickerW,
		ControlW:    controlW,
		ContentPadR: padR + controlW,
	}
}

func ResolvePickerX(x, w, pickerW int, direction string) int {
	pickerX := x + w - pickerW
	if direction == "left" {
		pickerX = x
	}
	return pickerX
}

type SliderVisualPlan struct {
	MinV   float64
	MaxV   float64
	Value  float64
	Pct    float64
	FillW  int
	ThumbX int
	ThumbR int
}

func BuildSliderVisualPlan(minV, maxV, value float64, w, h int, minThumb int) SliderVisualPlan {
	if maxV <= minV {
		maxV = minV + 100
	}
	if value < minV {
		value = minV
	}
	if value > maxV {
		value = maxV
	}
	pct := (value - minV) / (maxV - minV)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	fillW := int(float64(w) * pct)
	thumbR := riMinInt(h/3, 10)
	if thumbR < minThumb {
		thumbR = minThumb
	}
	return SliderVisualPlan{
		MinV:   minV,
		MaxV:   maxV,
		Value:  value,
		Pct:    pct,
		FillW:  fillW,
		ThumbX: fillW,
		ThumbR: thumbR,
	}
}

type SliderStateModel struct {
	MinV  float64
	MaxV  float64
	Value float64
}

func ResolveSliderState(path string, props map[string]any, sliderValues map[string]float64, defaultMin, defaultMax float64) SliderStateModel {
	minV := parseSliderPropFloat(props["min"], defaultMin)
	maxV := parseSliderPropFloat(props["max"], defaultMax)
	if maxV <= minV {
		maxV = minV + 100
	}

	value, seen := sliderValues[path]
	if !seen {
		if v, ok := parseOptionalSliderValue(props["value"]); ok {
			value = v
		} else {
			value = minV + (maxV-minV)/2
		}
	}

	if value < minV {
		value = minV
	}
	if value > maxV {
		value = maxV
	}
	sliderValues[path] = value

	return SliderStateModel{MinV: minV, MaxV: maxV, Value: value}
}

func ResolveSliderValueFromPointer(pointerX, boundsMinX, boundsWidth int, minV, maxV float64) (float64, bool) {
	if boundsWidth <= 0 {
		return 0, false
	}
	if maxV <= minV {
		maxV = minV + 100
	}
	pct := float64(pointerX-boundsMinX) / float64(boundsWidth)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	return minV + pct*(maxV-minV), true
}

type RadioGroupModel struct {
	GroupKey string
	Value    string
	Selected bool
}

func ResolveCheckboxChecked(path string, props map[string]any, boolValues map[string]bool, inputValues map[string]string) bool {
	checked := boolValues[path]
	if _, seen := inputValues[path+"__init"]; !seen {
		if anyToString(props["checked"], "") == "true" {
			checked = true
			boolValues[path] = true
		}
		inputValues[path+"__init"] = "1"
	}
	return checked
}

func ResolveRadioGroupModel(path string, props map[string]any, inputValues map[string]string) RadioGroupModel {
	groupName := anyToString(props["name"], path)
	groupKey := "radio:" + groupName
	myValue := anyToString(props["value"], path)
	if _, seen := inputValues[path+"__init"]; !seen {
		if anyToString(props["checked"], "") == "true" {
			inputValues[groupKey] = myValue
		}
		inputValues[path+"__init"] = "1"
	}
	return RadioGroupModel{
		GroupKey: groupKey,
		Value:    myValue,
		Selected: inputValues[groupKey] == myValue,
	}
}

type SelectModel struct {
	Key           string
	Options       []string
	Index         int
	SelectedLabel string
}

func ResolveSelectModel(path string, props map[string]any, inputValues map[string]string) SelectModel {
	options := parseSelectOptions(props["options"])
	if len(options) == 0 {
		options = []string{"Option 1", "Option 2", "Option 3"}
	}

	selKey := "sel:" + path
	if _, seen := inputValues[path+"__init"]; !seen {
		initSel := anyToString(props["value"], anyToString(props["selected"], ""))
		for i, o := range options {
			if o == initSel {
				inputValues[selKey] = strconv.Itoa(i)
				break
			}
		}
		inputValues[path+"__init"] = "1"
	}

	idx := 0
	if idxStr := inputValues[selKey]; idxStr != "" {
		if parsed, err := strconv.Atoi(idxStr); err == nil {
			idx = parsed
		}
	}
	if idx < 0 {
		idx = 0
	}
	if idx >= len(options) {
		idx = len(options) - 1
	}

	return SelectModel{
		Key:           selKey,
		Options:       options,
		Index:         idx,
		SelectedLabel: options[idx],
	}
}

func NextSelectIndex(current, total int) int {
	if total <= 0 {
		return 0
	}
	if current < 0 {
		current = 0
	}
	return (current + 1) % total
}

func parseSelectOptions(raw any) []string {
	var out []string
	switch opts := raw.(type) {
	case []any:
		for _, o := range opts {
			out = append(out, anyToString(o, ""))
		}
	case string:
		for _, o := range strings.Split(opts, ",") {
			out = append(out, strings.TrimSpace(o))
		}
	}
	return out
}

func parseSliderPropFloat(raw any, fallback float64) float64 {
	if v, ok := parseOptionalSliderValue(raw); ok {
		return v
	}
	return fallback
}

func parseOptionalSliderValue(raw any) (float64, bool) {
	switch v := raw.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func riMinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func riMaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

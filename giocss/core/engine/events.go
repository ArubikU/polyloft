package engine

import "strings"

func NormalizeFocusableInputPath(path string) string {
	base := path
	if strings.HasSuffix(base, "__spin") {
		base = strings.TrimSuffix(base, "__spin")
	}
	if strings.HasSuffix(base, "__picker") {
		base = strings.TrimSuffix(base, "__picker")
	}
	return base
}

func IsInputLikeProps(props map[string]any) bool {
	inputType := strings.ToLower(anyToString(props["type"], anyToString(props["inputtype"], "")))
	if inputType != "" {
		return true
	}
	return anyToString(props["tag"], "") == "input"
}

func ResolveInputComponentType(props map[string]any, fallback string) string {
	component := strings.ToLower(anyToString(props["type"], anyToString(props["inputtype"], fallback)))
	if component == "" {
		return fallback
	}
	return component
}

func ResolveOnFocusEvent(props map[string]any) string {
	return strings.TrimSpace(anyToString(props["onfocus"], ""))
}

func ResolveOnBlurEvent(props map[string]any) string {
	return strings.TrimSpace(anyToString(props["onblur"], ""))
}

func ResolveOnInputEvent(props map[string]any) string {
	return anyToString(props["oninput"], anyToString(props["event"], "input"))
}

func ResolveComponentEventName(props map[string]any, fallback string) string {
	return anyToString(props["event"], fallback)
}

func ResolveOnSubmitEvent(props map[string]any) string {
	return anyToString(props["onsubmit"], anyToString(props["event"], "submit"))
}

func ResolveOnChangeEvent(props map[string]any) string {
	return strings.TrimSpace(anyToString(props["onchange"], ""))
}

func ResolveOnShowPickerEvent(props map[string]any) string {
	eventName := strings.TrimSpace(anyToString(props["onshowpicker"], ""))
	if eventName == "" {
		return "showpicker"
	}
	return eventName
}

func ResolvePickerDirection(props map[string]any, fallback string) string {
	direction := strings.TrimSpace(strings.ToLower(anyToString(props["direction"], fallback)))
	if direction == "" {
		return fallback
	}
	return direction
}

func ResolveSpinnerDelta(pointerY, spinnerY, spinnerH int) int {
	delta := 1
	if pointerY >= spinnerY+spinnerH/2 {
		delta = -1
	}
	return delta
}

func BuildInputSpinnerPath(path string) string {
	return path + "__spin"
}

func BuildInputPickerPath(path string) string {
	return path + "__picker"
}

func BuildButtonEventPayload(source, component, text string, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": component,
		"text":      text,
		"timestamp": timestamp,
	}
}

func BuildFormSubmitEventPayload(source, formPath, formID string, values map[string]any, timestamp int64) map[string]any {
	if values == nil {
		values = map[string]any{}
	}
	return map[string]any{
		"source":    source,
		"component": "form",
		"formPath":  formPath,
		"formID":    formID,
		"values":    values,
		"timestamp": timestamp,
	}
}

func BuildInputValueEventPayload(source, component, value string, focused bool, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": component,
		"value":     value,
		"focused":   focused,
		"timestamp": timestamp,
	}
}

func BuildCheckboxEventPayload(source string, checked bool, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": "checkbox",
		"checked":   checked,
		"timestamp": timestamp,
	}
}

func BuildRadioEventPayload(source, radioValue string, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": "radio",
		"value":     radioValue,
		"timestamp": timestamp,
	}
}

func BuildSliderEventPayload(source, component string, sliderValue float64, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": component,
		"value":     sliderValue,
		"timestamp": timestamp,
	}
}

func BuildSelectEventPayload(source, selectedValue string, selectedIndex int, timestamp int64) map[string]any {
	return map[string]any{
		"source":    source,
		"component": "select",
		"value":     selectedValue,
		"index":     selectedIndex,
		"timestamp": timestamp,
	}
}

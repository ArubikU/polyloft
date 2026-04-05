package engine

import (
	"image"
	"strconv"

	coreinput "github.com/ArubikU/giocss/core/input"
)

func ToggleCheckboxState(path string, boolValues map[string]bool, timestamp int64) map[string]any {
	checked := !boolValues[path]
	boolValues[path] = checked
	return BuildCheckboxEventPayload(path, checked, timestamp)
}

func SelectRadioState(path string, groupKey string, radioValue string, inputValues map[string]string, timestamp int64) map[string]any {
	inputValues[groupKey] = radioValue
	return BuildRadioEventPayload(path, radioValue, timestamp)
}

func ResolveSliderPointerState(path string, boundsForPath map[string]image.Rectangle, lastPointer map[string]image.Point, sliderValues map[string]float64, minV float64, maxV float64, component string, timestamp int64) (map[string]any, bool) {
	bounds, ok := boundsForPath[path]
	if !ok {
		return nil, false
	}
	bw := bounds.Max.X - bounds.Min.X
	px := lastPointer[path].X
	newVal, ok := coreinput.ResolveSliderValueFromPointer(px, bounds.Min.X, bw, minV, maxV)
	if !ok {
		return nil, false
	}
	sliderValues[path] = newVal
	return BuildSliderEventPayload(path, component, newVal, timestamp), true
}

func CycleSelectState(path string, key string, inputValues map[string]string, options []string, currentIndex int, timestamp int64) (map[string]any, bool) {
	if len(options) == 0 {
		return nil, false
	}
	next := coreinput.NextSelectIndex(currentIndex, len(options))
	if idxStr := inputValues[key]; idxStr != "" {
		if n, err := strconv.Atoi(idxStr); err == nil {
			next = coreinput.NextSelectIndex(n, len(options))
		}
	}
	inputValues[key] = strconv.Itoa(next)
	return BuildSelectEventPayload(path, options[next], next, timestamp), true
}

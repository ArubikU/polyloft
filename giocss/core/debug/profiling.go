package debug

import (
	"reflect"
	"strconv"
	"strings"
)

func ChildSliceSignature(children []map[string]any) uint64 {
	if len(children) == 0 {
		return 0
	}
	var sig uint64 = 1469598103934665603
	const prime uint64 = 1099511628211
	for _, cm := range children {
		if len(cm) == 0 {
			sig ^= 0x9e3779b97f4a7c15
			sig *= prime
			continue
		}
		ptr := uint64(reflect.ValueOf(cm).Pointer())
		sig ^= ptr + 0x9e3779b97f4a7c15 + (sig << 6) + (sig >> 2)
		sig *= prime
	}
	return sig
}

func CloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

type ComponentSummary struct {
	ID        string
	ClassName string
	Component string
	Tag       string
	Role      string
	StyleHint string
	Attrs     map[string]string
}

func SummarizeComponentProps(kind string, props map[string]any) ComponentSummary {
	if len(props) == 0 {
		return ComponentSummary{}
	}
	id := strings.TrimSpace(profileAnyToString(props["id"], ""))
	className := strings.TrimSpace(profileAnyToString(props["class"], profileAnyToString(props["className"], "")))
	component := strings.TrimSpace(profileAnyToString(props["component"], profileAnyToString(props["native"], "")))
	tag := strings.TrimSpace(profileAnyToString(props["tag"], ""))
	if tag == "" {
		tag = kind
	}
	role := strings.TrimSpace(profileAnyToString(props["role"], ""))
	styleHint := strings.TrimSpace(profileAnyToString(props["style"], ""))
	if len(styleHint) > 220 {
		styleHint = styleHint[:220]
	}

	keys := []string{"name", "type", "event", "oninput", "onchange", "onsubmit", "src", "href", "direction", "placeholder", "value"}
	attrs := make(map[string]string)
	for _, key := range keys {
		if raw, ok := props[key]; ok {
			val := strings.TrimSpace(profileAnyToString(raw, ""))
			if val == "" {
				continue
			}
			if len(val) > 120 {
				val = val[:120]
			}
			attrs[key] = val
		}
	}
	if len(attrs) == 0 {
		attrs = nil
	}

	return ComponentSummary{ID: id, ClassName: className, Component: component, Tag: tag, Role: role, StyleHint: styleHint, Attrs: attrs}
}

func profileAnyToString(v any, fallback string) string {
	switch vv := v.(type) {
	case nil:
		return fallback
	case string:
		return vv
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(vv)
	case int64:
		return strconv.FormatInt(vv, 10)
	case float64:
		return strconv.FormatFloat(vv, 'g', -1, 64)
	default:
		return fallback
	}
}



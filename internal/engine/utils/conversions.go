package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ArubikU/polyloft/internal/common"
)

// Common string constants to reduce allocations
const (
	strNil   = "nil"
	strTrue  = "true"
	strFalse = "false"
)

// ToString converts a value to its string representation.
// This handles all Polyloft types including primitive wrappers, class instances, etc.
func ToString(v any) string {
	switch t := v.(type) {
	case nil:
		return strNil
	case string:
		return t
	case bool:
		if t {
			return strTrue
		}
		return strFalse
	case int:
		return fmt.Sprintf("%d", t)
	case float64:
		return fmt.Sprintf("%g", t)
	case float32:
		return fmt.Sprintf("%g", float64(t))
	case *common.EnumConstructor:
		// Return formatted enum type
		return common.GetTypeName(t)
	case *common.ClassConstructor:
		// Return formatted class type
		return common.GetTypeName(t)
	case *common.ClassInstance:
		// Handle primitive wrapper classes specially
		switch t.ClassName {
		case "String":
			if val, ok := t.Fields["_value"].(string); ok {
				return val
			}
		case "Int":
			if val, ok := t.Fields["_value"].(int); ok {
				return fmt.Sprintf("%d", val)
			}
		case "Float":
			if val, ok := t.Fields["_value"].(float64); ok {
				return fmt.Sprintf("%g", val)
			}
		case "Bool":
			if val, ok := t.Fields["_value"].(bool); ok {
				if val {
					return "true"
				}
				return "false"
			}
		}
		// Try to call the toString method if it exists
		if toStringMethod, exists := t.Methods["toString"]; exists {
			// toStringMethod is already a common.Func, no need to cast
			if result, err := toStringMethod(&common.Env{Vars: map[string]any{"this": t}, Consts: map[string]bool{}}, []any{}); err == nil {
				// Handle the result - it might be a String instance
				return ToString(result)
			}
		}
		// Fallback to default representation
		return fmt.Sprintf("%s@%p", t.ClassName, t)
	case *common.ClassDefinition:
		return fmt.Sprintf("class %s", t.Name)
	case *common.EnumValueInstance:
		if method, ok := t.Methods["toString"]; ok {
			if result, err := method(&common.Env{Vars: map[string]any{"this": t}, Consts: map[string]bool{}}, []any{}); err == nil {
				if str, ok := result.(string); ok {
					return str
				}
			}
		}
		if t.Definition != nil {
			return fmt.Sprintf("%s.%s", t.Definition.Name, t.Name)
		}
		return t.Name
	case *common.RecordInstance:
		if method, ok := t.Methods["toString"]; ok {
			if result, err := method(&common.Env{Vars: map[string]any{"this": t}, Consts: map[string]bool{}}, []any{}); err == nil {
				if str, ok := result.(string); ok {
					return str
				}
			}
		}
		if t.Definition != nil {
			parts := make([]string, 0, len(t.Definition.Components))
			for _, component := range t.Definition.Components {
				parts = append(parts, fmt.Sprintf("%s=%s", component.Name, ToString(t.Values[component.Name])))
			}
			return fmt.Sprintf("%s(%s)", t.Definition.Name, strings.Join(parts, ", "))
		}
		return "record"
	case []any:
		s := "["
		for i, e := range t {
			if i > 0 {
				s += ", "
			}
			s += ToString(e)
		}
		return s + "]"
	case map[string]any:
		// Check if it has a toString method (for Lambda type wrappers)
		if toStringFunc, hasToString := t["toString"]; hasToString {
			if fn, ok := common.ExtractFunc(toStringFunc); ok {
				if result, err := fn(nil, []any{}); err == nil {
					if str, ok := result.(string); ok {
						return str
					}
				}
			}
		}
		s := "{"
		i := 0
		for k, e := range t {
			if i > 0 {
				s += ", "
			}
			s += k + ": " + ToString(e)
			i++
		}
		return s + "}"
	default:
		// Fallback: try to use toString method for ClassInstance
		if ci, ok := v.(*common.ClassInstance); ok {
			if toStringMethod, exists := ci.Methods["toString"]; exists {
				if result, err := toStringMethod(&common.Env{Vars: map[string]any{"this": ci}, Consts: map[string]bool{}}, []any{}); err == nil {
					if str, ok := result.(string); ok {
						return str
					}
				}
			}
			return fmt.Sprintf("%s@%p", ci.ClassName, ci)
		}
		return fmt.Sprintf("%v", v)
	}
}

// AsFloat converts a value to float64 if possible.
// Returns the float value and a boolean indicating success.
func AsFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case *common.ClassInstance:
		// Handle Float and Int class instances
		if t.ClassName == "Float" {
			if val, ok := t.Fields["_value"].(float64); ok {
				return val, true
			}
		} else if t.ClassName == "Int" {
			if val, ok := t.Fields["_value"].(int); ok {
				return float64(val), true
			}
		}
		return 0, false
	case string:
		// Try to parse string as float
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// AsInt converts a value to int if possible.
// Returns the int value and a boolean indicating success.
func AsInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case float32:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case *common.ClassInstance:
		// Handle Int and Float class instances
		if t.ClassName == "Int" {
			if val, ok := t.Fields["_value"].(int); ok {
				return val, true
			}
		} else if t.ClassName == "Float" {
			if val, ok := t.Fields["_value"].(float64); ok {
				return int(val), true
			}
		}
		return 0, false
	case string:
		// Try to parse string as int
		if i, err := strconv.Atoi(t); err == nil {
			return i, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// AsBool converts a value to bool.
// Returns the boolean value.
func AsBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case *common.ClassInstance:
		if t.ClassName == "Bool" {
			if val, ok := t.Fields["_value"].(bool); ok {
				return val
			}
		}
		return true // Non-nil objects are truthy
	case nil:
		return false
	case string:
		return t != ""
	case int:
		return t != 0
	case float64:
		return t != 0
	case []any:
		return len(t) > 0
	case map[string]any:
		return len(t) > 0
	default:
		return true
	}
}

// Truthy returns whether a value is considered truthy in Polyloft.
// This is similar to AsBool but specifically for control flow conditions.
func Truthy(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case string:
		return t != ""
	case int:
		return t != 0
	case float64:
		return t != 0
	case *common.ClassInstance:
		if t.ClassName == "Bool" {
			if val, ok := t.Fields["_value"].(bool); ok {
				return val
			}
		}
		// Other class instances are truthy
		return true
	case []any:
		return len(t) > 0
	case map[string]any:
		return len(t) > 0
	default:
		return true
	}
}

// AsFloatArg extracts the i-th argument and converts it to float64.
// Returns 0, false if the argument doesn't exist or can't be converted.
func AsFloatArg(args []any, i int) (float64, bool) {
	if i >= len(args) {
		return 0, false
	}
	return AsFloat(args[i])
}

// AsIntArg extracts the i-th argument and converts it to int.
// Returns 0, false if the argument doesn't exist or can't be converted.
func AsIntArg(args []any, i int) (int, bool) {
	if i >= len(args) {
		return 0, false
	}
	return AsInt(args[i])
}

// AsStringArg extracts the i-th argument and converts it to string.
// Returns empty string if the argument doesn't exist.
func AsStringArg(args []any, i int) string {
	if i >= len(args) {
		return ""
	}
	return ToString(args[i])
}

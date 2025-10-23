// Package types defines built-in types for Polyloft
// This file contains common definitions shared across all built-in types
package types

import (
	"github.com/ArubikU/polyloft/internal/common"
)

// Import the shared Func type from common package
type Func = common.Func

// Import the shared TypeError from common package
type TypeError = common.TypeError

// toString converts any value to its string representation
func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		if v == float64(int64(v)) {
			return intToString(int64(v))
		}
		return floatToString(v)
	case int:
		return intToString(int64(v))
	case int64:
		return intToString(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "nil"
	case []any:
		if len(v) == 0 {
			return "[]"
		}
		result := "["
		for i, elem := range v {
			if i > 0 {
				result += ", "
			}
			result += toString(elem)
		}
		result += "]"
		return result
	case map[string]any:
		if len(v) == 0 {
			return "{}"
		}
		result := "{"
		first := true
		for key, val := range v {
			if !first {
				result += ", "
			}
			first = false
			result += key + ": " + toString(val)
		}
		result += "}"
		return result
	default:
		return "object"
	}
}

// toInt converts a value to an integer
func toInt(value any) (int, bool) {
	switch v := value.(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// toFloat converts a value to a float
func toFloat(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

// deepEqual performs deep equality comparison
func deepEqual(a, b any) bool {
	switch va := a.(type) {
	case []any:
		vb, ok := b.([]any)
		if !ok || len(va) != len(vb) {
			return false
		}
		for i := range va {
			if !deepEqual(va[i], vb[i]) {
				return false
			}
		}
		return true
	case map[string]any:
		vb, ok := b.(map[string]any)
		if !ok || len(va) != len(vb) {
			return false
		}
		for key, val := range va {
			otherVal, exists := vb[key]
			if !exists || !deepEqual(val, otherVal) {
				return false
			}
		}
		return true
	case string:
		vb, ok := b.(string)
		return ok && va == vb
	case float64:
		vb, ok := b.(float64)
		return ok && va == vb
	case int:
		vb, ok := b.(int)
		return ok && va == vb
	case bool:
		vb, ok := b.(bool)
		return ok && va == vb
	case nil:
		return b == nil
	default:
		return a == b
	}
}

// isTruthy returns whether a value is considered true in a boolean context
func isTruthy(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		return len(v) > 0
	case []any:
		return len(v) > 0
	case map[string]any:
		return len(v) > 0
	default:
		return true
	}
}

// simpleHash returns a simple hash code for a value
func simpleHash(value any) int {
	switch v := value.(type) {
	case nil:
		return 0
	case bool:
		if v {
			return 1
		}
		return 0
	case float64:
		return int(v * 1000)
	case int:
		return v
	case string:
		hash := 0
		for _, char := range v {
			hash = hash*31 + int(char)
		}
		return hash
	case []any:
		hash := 0
		for _, elem := range v {
			hash = hash*31 + simpleHash(elem)
		}
		return hash
	case map[string]any:
		hash := 0
		for key, val := range v {
			keyHash := simpleStringHash(key)
			valueHash := simpleHash(val)
			hash += keyHash + valueHash
		}
		return hash
	default:
		return 42 // Default hash
	}
}

// simpleStringHash returns a simple hash code for a string
func simpleStringHash(str string) int {
	hash := 0
	for _, char := range str {
		hash = hash*31 + int(char)
	}
	return hash
}

// Helper functions for number formatting
func intToString(value int64) string {
	if value >= 0 {
		return positiveIntToString(value)
	}
	return "-" + positiveIntToString(-value)
}

func positiveIntToString(value int64) string {
	if value == 0 {
		return "0"
	}

	var digits []byte
	for value > 0 {
		digits = append([]byte{'0' + byte(value%10)}, digits...)
		value /= 10
	}

	return string(digits)
}

func floatToString(value float64) string {
	// Simple float to string conversion
	// In a real implementation, you'd want more sophisticated formatting
	if value == float64(int64(value)) {
		return intToString(int64(value))
	}

	// For simplicity, we'll use a basic approach
	intPart := int64(value)
	fracPart := value - float64(intPart)

	result := intToString(intPart)

	if fracPart != 0 {
		result += "."
		// Simple fractional part handling (limited precision)
		fracPart = fracPart * 1000000 // 6 decimal places
		if fracPart < 0 {
			fracPart = -fracPart
		}
		fracStr := positiveIntToString(int64(fracPart))
		// Remove trailing zeros
		for len(fracStr) > 1 && fracStr[len(fracStr)-1] == '0' {
			fracStr = fracStr[:len(fracStr)-1]
		}
		result += fracStr
	}

	return result
}

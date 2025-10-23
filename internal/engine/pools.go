package engine

import (
	"sync"
)

// String pool for common strings to reduce allocations
var stringPool = sync.Pool{
	New: func() interface{} {
		return new(string)
	},
}

// Common string cache for frequently used strings
var commonStrings = map[string]string{
	"":        "",
	"nil":     "nil",
	"true":    "true",
	"false":   "false",
	"null":    "null",
	"0":       "0",
	"1":       "1",
	"-1":      "-1",
	"[]":      "[]",
	"{}":      "{}",
	"string":  "string",
	"String":  "String",
	"int":     "int",
	"Int":     "Int",
	"float":   "float",
	"Float":   "Float",
	"bool":    "bool",
	"Bool":    "Bool",
	"array":   "array",
	"Array":   "Array",
	"map":     "map",
	"Map":     "Map",
	"object":  "object",
	"Object":  "Object",
	"number":  "number",
	"Number":  "Number",
	"void":    "void",
	"Void":    "Void",
	"any":     "any",
	"Any":     "Any",
	"this":    "this",
	"super":   "super",
	"length":  "length",
	"size":    "size",
	"name":    "name",
	"value":   "value",
	"_value":  "_value",
	"_items":  "_items",
	"_keys":   "_keys",
	"_values": "_values",
}

// GetCachedString returns a cached string if available, otherwise returns the input
func GetCachedString(s string) string {
	if cached, ok := commonStrings[s]; ok {
		return cached
	}
	return s
}

// Small integer cache for frequently used integers (-128 to 255)
const (
	intCacheMin = -128
	intCacheMax = 255
)

var intCache [intCacheMax - intCacheMin + 1]int

func init() {
	// Pre-initialize integer cache
	for i := range intCache {
		intCache[i] = i + intCacheMin
	}
}

// GetCachedInt returns a cached integer for small values
func GetCachedInt(n int) int {
	if n >= intCacheMin && n <= intCacheMax {
		return intCache[n-intCacheMin]
	}
	return n
}

// Float cache for common float values
var commonFloats = map[int]float64{
	0:  0.0,
	1:  1.0,
	-1: -1.0,
	2:  2.0,
	10: 10.0,
}

// GetCachedFloat returns a cached float for common values
func GetCachedFloat(n float64) float64 {
	intVal := int(n)
	if float64(intVal) == n {
		if cached, ok := commonFloats[intVal]; ok {
			return cached
		}
	}
	return n
}

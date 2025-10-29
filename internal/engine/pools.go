package engine

import (
	"sync"

	"github.com/ArubikU/polyloft/internal/common"
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

// Environment pool for reusing environment allocations
var envPool = sync.Pool{
	New: func() interface{} {
		return &common.Env{
			Vars:             make(map[string]any, 16), // Pre-allocate reasonable size
			Consts:           make(map[string]bool, 4),
			Finals:           make(map[string]bool, 4),
			ImportedClasses:  make(map[string]string, 8),
			ImportedPackages: make(map[string]struct{}, 4),
			Defers:           make([]func() error, 0, 4),
		}
	},
}

// GetPooledEnv gets an environment from the pool and initializes it
func GetPooledEnv(parent *common.Env) *common.Env {
	env := envPool.Get().(*common.Env)
	env.Parent = parent
	
	// Clear maps (Go 1.21+ has clear() but we'll do it manually for compatibility)
	for k := range env.Vars {
		delete(env.Vars, k)
	}
	for k := range env.Consts {
		delete(env.Consts, k)
	}
	for k := range env.Finals {
		delete(env.Finals, k)
	}
	for k := range env.ImportedClasses {
		delete(env.ImportedClasses, k)
	}
	for k := range env.ImportedPackages {
		delete(env.ImportedPackages, k)
	}
	
	// Reset slices
	env.Defers = env.Defers[:0]
	env.PositionStack = env.PositionStack[:0]
	env.CodeContext = env.CodeContext[:0]
	env.SourceLines = env.SourceLines[:0]
	
	// Reset other fields
	env.FileName = ""
	env.PackageName = ""
	env.CurrentLine = 0
	env.CurrentColumn = 0
	
	return env
}

// ReleaseEnv returns an environment to the pool
func ReleaseEnv(env *common.Env) {
	if env == nil {
		return
	}
	env.Parent = nil
	envPool.Put(env)
}

// FastIntOperation performs optimized integer arithmetic
func FastIntOperation(op int, a, b int) int {
	switch op {
	case 1: // OpPlus
		return a + b
	case 2: // OpMinus
		return a - b
	case 3: // OpMul
		return a * b
	case 4: // OpDiv
		if b == 0 {
			return 0 // Handled by caller
		}
		return a / b
	case 5: // OpMod
		if b == 0 {
			return 0 // Handled by caller
		}
		return a % b
	default:
		return 0
	}
}

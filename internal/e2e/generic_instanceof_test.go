package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// TestArrayInstanceOf tests instanceof for arrays (arrays are not generic, they're type Any)
func TestArrayInstanceOf(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Array instanceof Array",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array"))
			`,
			expected: "true",
		},
		{
			name: "Array with mixed types",
			code: `
				let arr = ["hello", 1, "world", 2]
				println(Sys.instanceof(arr, "Array"))
			`,
			expected: "true",
		},
		{
			name: "Empty array",
			code: `
				let arr = []
				println(Sys.instanceof(arr, "Array"))
			`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeInstanceOf(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContains(t, output, tt.expected)
		})
	}
}

// TestListInstanceOfGeneric tests instanceof for List with generic types
func TestListInstanceOfGeneric(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "List<String> instanceof List",
			code: `
				let list = List<String>("hello", "world")
				println(Sys.instanceof(list, "List"))
			`,
			expected: "true",
		},
		{
			name: "List<String> instanceof List<String>",
			code: `
				let list = List<String>("hello", "world")
				println(Sys.instanceof(list, "List<String>"))
			`,
			expected: "true",
		},
		{
			name: "List<Int> instanceof List<Int>",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.instanceof(list, "List<Int>"))
			`,
			expected: "true",
		},
		{
			name: "List<Int> instanceof List<Number>",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.instanceof(list, "List<Number>"))
			`,
			expected: "true",
		},
		{
			name: "List<Int> not instanceof List<String>",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.instanceof(list, "List<String>"))
			`,
			expected: "false",
		},
		{
			name: "List<Any> with mixed types",
			code: `
				let list = List<Any>(1, "hello", 3.14)
				println(Sys.instanceof(list, "List"))
			`,
			expected: "true",
		},
		{
			name: "Empty List instanceof List",
			code: `
				let list = List<Int>()
				println(Sys.instanceof(list, "List"))
			`,
			expected: "true",
		},
		{
			name: "List<Int> instanceof List<? extends Number>",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.instanceof(list, "List<? extends Number>"))
			`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeInstanceOf(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContains(t, output, tt.expected)
		})
	}
}

// TestSysTypeWithGenerics tests Sys.type() for generic collections
// Note: Arrays are not generic (they're type Any), only List/Set/Map have generic types
func TestSysTypeWithGenerics(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Array shows Array (not generic)",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.type(arr))
			`,
			expected: "Array",
		},
		{
			name: "List<String> shows type with generics",
			code: `
				let list = List<String>("hello", "world")
				println(Sys.type(list))
			`,
			expected: "List<String>",
		},
		{
			name: "List<Int> shows type with generics",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.type(list))
			`,
			expected: "List<Int>",
		},
		{
			name: "Set<String> shows type with generics",
			code: `
				let set = Set<String>("a", "b", "c")
				println(Sys.type(set))
			`,
			expected: "Set<String>",
		},
		{
			name: "Map<String, Int> shows type with generics",
			code: `
				let map = Map<String, Int>()
				println(Sys.type(map))
			`,
			expected: "Map<String, Int>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeInstanceOf(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContains(t, output, tt.expected)
		})
	}
}

// TestUnionTypesInGenerics tests union types in generic parameters
// Note: Arrays are not generic, so union tests focus on List/Set/Map
func TestUnionTypesInGenerics(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "List with union type parameter",
			code: `
				let list = List<Any>(1, "hello", 3.14)
				println(Sys.instanceof(list, "List<Int | String | Float>"))
			`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeInstanceOf(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContains(t, output, tt.expected)
		})
	}
}

// TestWildcardTypes tests wildcard type parameters
// Note: Arrays are not generic, so wildcard tests focus on List/Set/Map
func TestWildcardTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "List<? extends Number>",
			code: `
				let list = List<Int>(1, 2, 3)
				println(Sys.instanceof(list, "List<? extends Number>"))
			`,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeInstanceOf(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContains(t, output, tt.expected)
		})
	}
}

// Helper function to run code and capture output
func runCodeInstanceOf(code string) (string, error) {
	// Reset global registries to ensure test isolation
	engine.ResetGlobalRegistries()
	
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(code))
	par := parser.New(items)
	prog, err := par.Parse()
	if err != nil {
		return "", err
	}

	// Capture stdout using a buffer
	buf := &bytes.Buffer{}
	_, err = engine.Eval(prog, engine.Options{Stdout: buf})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// Helper function to check if output contains expected string
func assertOutputContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but got: %q", expected, output)
	}
}

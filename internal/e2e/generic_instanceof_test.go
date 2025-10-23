package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// TestArrayInstanceOfGeneric tests instanceof for arrays with generic types
func TestArrayInstanceOfGeneric(t *testing.T) {
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
			name: "Array<Int> - all ints",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<Int>"))
			`,
			expected: "true",
		},
		{
			name: "Array<Number> - all ints",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<Number>"))
			`,
			expected: "true",
		},
		{
			name: "Array<Integer> - all ints",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<Integer>"))
			`,
			expected: "true",
		},
		{
			name: "Array<Float> - should be false for ints",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<Float>"))
			`,
			expected: "false",
		},
		{
			name: "Array<String | Int> - union type",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<String | Int>"))
			`,
			expected: "true",
		},
		{
			name: "Array with mixed types - String | Int",
			code: `
				let arr = ["hello", 1, "world", 2]
				println(Sys.instanceof(arr, "Array<String | Int>"))
			`,
			expected: "true",
		},
		{
			name: "Array<? extends Number> - wildcard",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.instanceof(arr, "Array<? extends Number>"))
			`,
			expected: "true",
		},
		{
			name: "Array<Float> with floats",
			code: `
				let arr = [1.5, 2.5, 3.5]
				println(Sys.instanceof(arr, "Array<Float>"))
			`,
			expected: "true",
		},
		{
			name: "Array<String> with strings",
			code: `
				let arr = ["a", "b", "c"]
				println(Sys.instanceof(arr, "Array<String>"))
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
func TestSysTypeWithGenerics(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Array of ints shows Array<Int>",
			code: `
				let arr = [1, 2, 3, 4]
				println(Sys.type(arr))
			`,
			expected: "Array<Int>",
		},
		{
			name: "Array of floats shows Array<Float>",
			code: `
				let arr = [1.5, 2.5, 3.5]
				println(Sys.type(arr))
			`,
			expected: "Array<Float>",
		},
		{
			name: "Array of strings shows Array<String>",
			code: `
				let arr = ["a", "b", "c"]
				println(Sys.type(arr))
			`,
			expected: "Array<String>",
		},
		{
			name: "Mixed int/float array shows Array<Number>",
			code: `
				let arr = [1, 2.5, 3]
				println(Sys.type(arr))
			`,
			expected: "Array<Number>",
		},
		{
			name: "Empty array shows Array without type parameter",
			code: `
				let arr = []
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
func TestUnionTypesInGenerics(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Array with String | Int union",
			code: `
				let arr = ["hello", 1, "world", 2]
				println(Sys.instanceof(arr, "Array<String | Int>"))
			`,
			expected: "true",
		},
		{
			name: "Array with Int | Float union",
			code: `
				let arr = [1, 2.5, 3]
				println(Sys.instanceof(arr, "Array<Int | Float>"))
			`,
			expected: "true",
		},
		{
			name: "Array with String | Int | Float union",
			code: `
				let arr = ["hello", 1, 2.5]
				println(Sys.instanceof(arr, "Array<String | Int | Float>"))
			`,
			expected: "true",
		},
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
func TestWildcardTypes(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Array<? extends Number> with ints",
			code: `
				let arr = [1, 2, 3]
				println(Sys.instanceof(arr, "Array<? extends Number>"))
			`,
			expected: "true",
		},
		{
			name: "Array<? extends Number> with floats",
			code: `
				let arr = [1.5, 2.5, 3.5]
				println(Sys.instanceof(arr, "Array<? extends Number>"))
			`,
			expected: "true",
		},
		{
			name: "Array<? extends Number> with mixed numbers",
			code: `
				let arr = [1, 2.5, 3]
				println(Sys.instanceof(arr, "Array<? extends Number>"))
			`,
			expected: "true",
		},
		{
			name: "Array<?> unbounded wildcard",
			code: `
				let arr = ["hello", 1, true]
				println(Sys.instanceof(arr, "Array<?>"))
			`,
			expected: "true",
		},
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

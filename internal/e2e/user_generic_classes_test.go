package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// TestUserDefinedGenericClasses tests user-defined generic classes behave like built-in collections
func TestUserDefinedGenericClasses(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Pair class with explicit types shows generic parameters",
			code: `
				class PairA<K, V>:
					private var key: K
					private var value: V
					PairA(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairA<String, Int>("name", 42)
				println(Sys.type(pair))
			`,
			expected: "PairA<String, Int>",
		},
		{
			name: "Box class with single type parameter",
			code: `
				class BoxA<T>:
					private var content: T
					BoxA(item: T):
						this.content = item
					end
				end
				let box = BoxA<String>("hello")
				println(Sys.type(box))
			`,
			expected: "BoxA<String>",
		},
		{
			name: "instanceof works with user-defined generic class",
			code: `
				class PairB<K, V>:
					private var key: K
					private var value: V
					PairB(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairB<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairB"))
			`,
			expected: "true",
		},
		{
			name: "instanceof with exact generic type parameters",
			code: `
				class PairC<K, V>:
					private var key: K
					private var value: V
					PairC(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairC<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairC<String, Int>"))
			`,
			expected: "true",
		},
		{
			name: "instanceof false when type parameters don't match",
			code: `
				class PairD<K, V>:
					private var key: K
					private var value: V
					PairD(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairD<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairD<Int, String>"))
			`,
			expected: "false",
		},
		{
			name: "instanceof with wildcard on first parameter",
			code: `
				class PairE<K, V>:
					private var key: K
					private var value: V
					PairE(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairE<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairE<?, Int>"))
			`,
			expected: "true",
		},
		{
			name: "instanceof with wildcard on second parameter",
			code: `
				class PairF<K, V>:
					private var key: K
					private var value: V
					PairF(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairF<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairF<String, ?>"))
			`,
			expected: "true",
		},
		{
			name: "instanceof with all wildcards",
			code: `
				class PairG<K, V>:
					private var key: K
					private var value: V
					PairG(k: K, v: V):
						this.key = k
						this.value = v
					end
				end
				let pair = PairG<String, Int>("name", 42)
				println(Sys.instanceof(pair, "PairG<?, ?>"))
			`,
			expected: "true",
		},
		{
			name: "Multiple instances with different type parameters",
			code: `
				class BoxB<T>:
					private var content: T
					BoxB(item: T):
						this.content = item
					end
				end
				let box1 = BoxB<String>("hello")
				let box2 = BoxB<Int>(42)
				println(Sys.type(box1))
				println(Sys.type(box2))
			`,
			expected: "BoxB<String>\nBoxB<Int>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCodeUserGeneric(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertOutputContainsUserGeneric(t, output, tt.expected)
		})
	}
}

// Helper function to run code and capture output
func runCodeUserGeneric(code string) (string, error) {
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
func assertOutputContainsUserGeneric(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain %q, but got: %q", expected, output)
	}
}

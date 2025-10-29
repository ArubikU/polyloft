package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// TestGenericTypeChecking tests runtime type checking for user-defined generic classes
func TestGenericTypeChecking(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Box<String> accepts string parameter",
			code: `
				class TypeCheckBoxA<T>:
					private var content: T
					TypeCheckBoxA(item: T):
						this.content = item
					end
					def set(newContent: T) -> Void:
						this.content = newContent
					end
				end
				let box = TypeCheckBoxA<String>("hello")
				box.set("world")
				println("OK")
			`,
			shouldError: false,
		},
		{
			name: "Box<Int> accepts int parameter",
			code: `
				class TypeCheckBoxB<T>:
					private var content: T
					TypeCheckBoxB(item: T):
						this.content = item
					end
					def set(newContent: T) -> Void:
						this.content = newContent
					end
				end
				let box = TypeCheckBoxB<Int>(42)
				box.set(100)
				println("OK")
			`,
			shouldError: false,
		},
		{
			name: "Box<String> rejects int parameter",
			code: `
				class TypeCheckBoxC<T>:
					private var content: T
					TypeCheckBoxC(item: T):
						this.content = item
					end
					def set(newContent: T) -> Void:
						this.content = newContent
					end
				end
				let box = TypeCheckBoxC<String>("hello")
				box.set(42)
			`,
			shouldError: true,
			errorMsg:    "String",
		},
		{
			name: "Box<Int> rejects string parameter",
			code: `
				class TypeCheckBoxD<T>:
					private var content: T
					TypeCheckBoxD(item: T):
						this.content = item
					end
					def set(newContent: T) -> Void:
						this.content = newContent
					end
				end
				let box = TypeCheckBoxD<Int>(42)
				box.set("hello")
			`,
			shouldError: true,
			errorMsg:    "Int",
		},
		{
			name: "Pair<String, Int> accepts correct types",
			code: `
				class TypeCheckPairA<K, V>:
					private var key: K
					private var value: V
					TypeCheckPairA(k: K, v: V):
						this.key = k
						this.value = v
					end
					def setKey(newKey: K) -> Void:
						this.key = newKey
					end
					def setValue(newValue: V) -> Void:
						this.value = newValue
					end
				end
				let pair = TypeCheckPairA<String, Int>("name", 42)
				pair.setKey("newName")
				pair.setValue(100)
				println("OK")
			`,
			shouldError: false,
		},
		{
			name: "Pair<String, Int> rejects wrong key type",
			code: `
				class TypeCheckPairB<K, V>:
					private var key: K
					private var value: V
					TypeCheckPairB(k: K, v: V):
						this.key = k
						this.value = v
					end
					def setKey(newKey: K) -> Void:
						this.key = newKey
					end
				end
				let pair = TypeCheckPairB<String, Int>("name", 42)
				pair.setKey(999)
			`,
			shouldError: true,
			errorMsg:    "String",
		},
		{
			name: "Pair<String, Int> rejects wrong value type",
			code: `
				class TypeCheckPairC<K, V>:
					private var key: K
					private var value: V
					TypeCheckPairC(k: K, v: V):
						this.key = k
						this.value = v
					end
					def setValue(newValue: V) -> Void:
						this.value = newValue
					end
				end
				let pair = TypeCheckPairC<String, Int>("name", 42)
				pair.setValue("wrong")
			`,
			shouldError: true,
			errorMsg:    "Int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runCodeTypeCheck(tt.code)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error containing %q, but got no error", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, but got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// Helper function to run code and capture output
func runCodeTypeCheck(code string) (string, error) {
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

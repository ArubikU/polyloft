package e2e

import (
	"bytes"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// runCodeWithOutput runs code and returns the printed output
func runCodeWithOutput(code string) (string, error) {
	// Reset global registries to ensure test isolation
	engine.ResetGlobalRegistries()
	
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(code))
	par := parser.New(items)
	prog, err := par.Parse()
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	_, err = engine.Eval(prog, engine.Options{Stdout: buf})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TestMath_StaticMethods tests Math module static methods
func TestMath_StaticMethods(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Math.sqrt",
			code:     "let x = Math.sqrt(9)\nprintln(x)",
			expected: "3\n",
		},
		{
			name:     "Math.pow",
			code:     "let x = Math.pow(2, 3)\nprintln(x)",
			expected: "8\n",
		},
		{
			name:     "Math.PI",
			code:     "println(Math.PI)",
			expected: "3.141592653589793\n",
		},
		{
			name:     "Math.E",
			code:     "println(Math.E)",
			expected: "2.718281828459045\n",
		},
		{
			name:     "Math.abs",
			code:     "let x = Math.abs(-5.5)\nprintln(x)",
			expected: "5.5\n",
		},
		{
			name:     "Math.floor",
			code:     "let x = Math.floor(3.7)\nprintln(x)",
			expected: "3\n",
		},
		{
			name:     "Math.ceil",
			code:     "let x = Math.ceil(3.2)\nprintln(x)",
			expected: "4\n",
		},
		{
			name:     "Math.min",
			code:     "let x = Math.min(5, 3)\nprintln(x)",
			expected: "3\n",
		},
		{
			name:     "Math.max",
			code:     "let x = Math.max(5, 3)\nprintln(x)",
			expected: "5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCodeWithOutput(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSys_StaticMethods tests Sys module static methods
func TestSys_StaticMethods(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Sys.type",
			code:     "let x = 42\nprintln(Sys.type(x))",
			expected: "int\n",
		},
		{
			name:     "Sys.type_string",
			code:     "let x = \"hello\"\nprintln(Sys.type(x))",
			expected: "string\n",
		},
		{
			name:     "Sys.type_float",
			code:     "let x = 3.14\nprintln(Sys.type(x))",
			expected: "float\n",
		},
		{
			name:     "Sys.format",
			code:     "let x = Sys.format(\"Hello %s, you are %d\", \"World\", 42)\nprintln(x)",
			expected: "Hello World, you are 42\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runCodeWithOutput(tt.code)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestStaticModules_NotConstructable tests that static modules cannot be instantiated
func TestStaticModules_NotConstructable(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "Cannot instantiate Math",
			code:    "let m = Math()",
			wantErr: true,
		},
		{
			name:    "Cannot instantiate Sys",
			code:    "let s = Sys()",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runCodeWithOutput(tt.code)
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

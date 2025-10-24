package e2e

import (
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// Phase 3: Type Bounds tests

func TestBounds_WildcardUpperBound(t *testing.T) {
	code := `
let list = List<? extends Number>(10, 20, 30)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := result.(int)
	if !ok || size != 3 {
		t.Fatalf("Expected size 3, got %v", result)
	}
}

func TestBounds_WildcardLowerBound(t *testing.T) {
	code := `
let list = List<? super Integer>(10, 20, 30)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := result.(int)
	if !ok || size != 3 {
		t.Fatalf("Expected size 3, got %v", result)
	}
}

func TestBounds_SetWithBounds(t *testing.T) {
	code := `
let set = Set<? extends String>("a", "b", "c")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	iresult, eri := utils.AsInt(result)
	if !eri {
		t.Fatalf("Expected int result, got %T", result)
	}

	if iresult != 3 {
		t.Fatalf("Expected size 3, got %v", result)
	}
}

func TestBounds_MapWithBounds(t *testing.T) {
	code := `
let map = Map<? extends String, ? extends Number>()
map.put("key", 42)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	iresult, eri := utils.AsInt(result)
	if !eri {
		t.Fatalf("Expected int result, got %T", result)
	}
	if iresult != 1 {
		t.Fatalf("Expected size 1, got %v", result)
	}
}

func TestBounds_ToString_UpperBound(t *testing.T) {
	code := `
let list = List<? extends Number>()
return list.toString()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str == "" {
		t.Fatalf("Expected string, got %T", result)
	}

	if !strings.Contains(str, "? extends Number") {
		t.Fatalf("Expected type string to contain '? extends Number', got %s", str)
	}
}

func TestBounds_ToString_LowerBound(t *testing.T) {
	code := `
let list = List<? super Integer>()
return list.toString()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	if !strings.Contains(str, "? super Integer") {
		t.Fatalf("Expected type string to contain '? super Integer', got %s", str)
	}
}

func TestBounds_VarianceWithBounds(t *testing.T) {
	code := `
let list = List<out ? extends Number>(10, 20)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := result.(int)
	if !ok || size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func runCodeBounds(code string) (any, error) {
	// Reset global registries to ensure test isolation
	engine.ResetGlobalRegistries()

	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(code))
	par := parser.New(items)
	prog, err := par.Parse()

	if err != nil {
		return nil, err
	}

	return engine.Eval(prog, engine.Options{})
}

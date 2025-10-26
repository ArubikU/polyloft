package e2e

import (
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Phase 2: Variance tests - covariance (out), contravariance (in), invariant

func TestVariance_Covariant_List(t *testing.T) {
	code := `
let list = List<out Number>()
list.add(42)
list.add(3.14)
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

func TestVariance_Contravariant_List(t *testing.T) {
	code := `
let list = List<in Integer>()
list.add(10)
list.add(20)
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

func TestVariance_Invariant_List(t *testing.T) {
	code := `
let list = List<Int>()
list.add(1)
list.add(2)
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

func TestVariance_Covariant_Set(t *testing.T) {
	code := `
let set = Set<out String>()
set.add("hello")
set.add("world")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := utils.AsInt(result)
	if !ok || size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestVariance_Contravariant_Set(t *testing.T) {
	code := `
let set = Set<in String>()
set.add("a")
set.add("b")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := utils.AsInt(result)
	if !ok || size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestVariance_Mixed_Map(t *testing.T) {
	code := `
let map = Map<in String, out Integer>()
map.put("key", 42)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := result.(int)
	if !ok || size != 1 {
		t.Fatalf("Expected size 1, got %v", result)
	}
}

func TestVariance_ToString_Covariant(t *testing.T) {
	code := `
let list = List<out Number>()
return list.toString()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)

	if !strings.Contains(str, "out Number") {
		t.Fatalf("Expected 'out Number' in toString, got %s", str)
	}
}

func TestVariance_ToString_Contravariant(t *testing.T) {
	code := `
let list = List<in Integer>()
return list.toString()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)

	if !strings.Contains(str, "in Integer") {
		t.Fatalf("Expected 'in Integer' in toString, got %s", str)
	}
}

func TestVariance_ToString_Invariant(t *testing.T) {
	code := `
let list = List<String>()
return list.toString()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)

	// Invariant should not have in/out prefix
	if strings.Contains(str, "in ") || strings.Contains(str, "out ") {
		t.Fatalf("Expected no variance annotation in invariant toString, got %s", str)
	}
}

func TestVariance_WithWildcard_Out(t *testing.T) {
	code := `
let list = List<out ?>()
list.add(1)
list.add("text")
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

func TestVariance_WithWildcard_In(t *testing.T) {
	code := `
let list = List<in ?>()
list.add(42)
list.add(true)
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

func TestVariance_WithBoundedWildcard_Out(t *testing.T) {
	code := `
let list = List<out ? extends Number>()
list.add(10)
list.add(3.14)
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

func TestVariance_WithBoundedWildcard_In(t *testing.T) {
	code := `
let list = List<in ? super Integer>()
list.add(42)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, ok := result.(int)
	if !ok || size != 1 {
		t.Fatalf("Expected size 1, got %v", result)
	}
}

func TestVariance_Producer_Out(t *testing.T) {
	// Covariant (out) - producer pattern
	// List<out T> can produce T values (read-only)
	code := `
let list = List<out String>("hello", "world")
let first = list.get(0)
return first
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str != "hello" {
		t.Fatalf("Expected 'hello', got %v", result)
	}
}

func TestVariance_Consumer_In(t *testing.T) {
	// Contravariant (in) - consumer pattern
	// List<in T> can consume T values (write-only)
	code := `
let list = List<in Integer>()
list.add(10)
list.add(20)
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

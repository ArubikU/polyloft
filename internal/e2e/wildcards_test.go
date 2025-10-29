package e2e

import (
	"strings"
	"testing"
	
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Phase 1: Wildcard tests - unbounded, upper bound, lower bound

func TestWildcard_Unbounded_List(t *testing.T) {
	code := `
let list = List<?>()
list.add(1)
list.add("string")
list.add(true)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 3 {
		t.Fatalf("Expected size 3, got %v", result)
	}
}

func TestWildcard_Unbounded_Set(t *testing.T) {
	code := `
let set = Set<?>()
set.add(1)
set.add("hello")
set.add(true)
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 3 {
		t.Fatalf("Expected size 3, got %v", result)
	}
}

func TestWildcard_Unbounded_Map(t *testing.T) {
	code := `
let map = Map<?, ?>()
map.put(1, "one")
map.put("two", 2)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_UpperBound_List(t *testing.T) {
	code := `
let list = List<? extends Number>()
list.add(1)
list.add(2.5)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_UpperBound_Set(t *testing.T) {
	code := `
let set = Set<? extends String>()
set.add("hello")
set.add("world")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_UpperBound_Map(t *testing.T) {
	code := `
let map = Map<String, ? extends Number>()
map.put("one", 1)
map.put("two", 2.5)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_LowerBound_List(t *testing.T) {
	code := `
let list = List<? super Integer>()
list.add(10)
list.add(20)
return list.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_LowerBound_Set(t *testing.T) {
	code := `
let set = Set<? super String>()
set.add("hello")
set.add("world")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_LowerBound_Map(t *testing.T) {
	code := `
let map = Map<? super String, Integer>()
map.put("one", 1)
map.put("two", 2)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2, got %v", result)
	}
}

func TestWildcard_Mixed_Map(t *testing.T) {
	code := `
let map = Map<? extends String, ? super Integer>()
map.put("key", 42)
return map.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 1 {
		t.Fatalf("Expected size 1, got %v", result)
	}
}

func TestWildcard_ToString_Unbounded(t *testing.T) {
	code := `
let list = List<?>()
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

	if !strings.Contains(str, "?") {
		t.Fatalf("Expected wildcard '?' in toString, got %s", str)
	}
}

func TestWildcard_ToString_UpperBound(t *testing.T) {
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
		t.Fatalf("Expected '? extends Number' in toString, got %s", str)
	}
}

func TestWildcard_ToString_LowerBound(t *testing.T) {
	code := `
let list = List<? super Integer>()
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

	if !strings.Contains(str, "? super Integer") {
		t.Fatalf("Expected '? super Integer' in toString, got %s", str)
	}
}

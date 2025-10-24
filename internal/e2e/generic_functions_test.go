package e2e

import (
	"testing"

	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func TestGenericFunction_Identity(t *testing.T) {
	code := `
def identity<T>(value: T) -> T:
    return value
end

let result1 = identity(42)
let result2 = identity("hello")
let result3 = identity(true)

return [result1, result2, result3]
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val, ok := result.(*common.ClassInstance)
	if !ok {
		t.Fatalf("Expected ClassInstance got %T", result)
	}
	arr, err := engine.ClassInstanceToArray(val)
	if err != nil {
		t.Fatalf("Failed to convert to array: %v", err)
	}

	if len(arr) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(arr))
	}

	// Use utils functions to handle both native types and ClassInstance
	intVal, ok := utils.AsInt(arr[0])
	if !ok || intVal != 42 {
		t.Errorf("Expected 42, got %v", arr[0])
	}
	if utils.ToString(arr[1]) != "hello" {
		t.Errorf("Expected 'hello', got %v", arr[1])
	}
	if utils.AsBool(arr[2]) != true {
		t.Errorf("Expected true, got %v", arr[2])
	}
}

func TestGenericFunction_MakePair(t *testing.T) {
	code := `
def make_pair<K, V>(key: K, value: V) -> Any:
    return [key, value]
end

let pair = make_pair("name", 42)
return pair
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val, ok := result.(*common.ClassInstance)
	if !ok {
		t.Fatalf("Expected ClassInstance got %T", result)
	}
	arr, err := engine.ClassInstanceToArray(val)
	if err != nil {
		t.Fatalf("Failed to convert to array: %v", err)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(arr))
	}

	// Use utils functions to handle both native types and ClassInstance
	if utils.ToString(arr[0]) != "name" {
		t.Errorf("Expected 'name', got %v", arr[0])
	}
	intVal, ok := utils.AsInt(arr[1])
	if !ok || intVal != 42 {
		t.Errorf("Expected 42, got %v", arr[1])
	}
}

func TestGenericFunction_Top(t *testing.T) {
	code := `
def top<T>(items: Any) -> Any:
    return items[0]
end

let result = top([1, 2, 3, 4, 5])
return result
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.AsInt to handle both native ints and ClassInstance
	value, ok := utils.AsInt(result)
	if !ok {
		t.Fatalf("Expected int, got %T", result)
	}

	if value != 1 {
		t.Errorf("Expected 1, got %v", value)
	}
}

func TestGenericFunction_Swap(t *testing.T) {
	code := `
def swap<T, U>(a: T, b: U) -> Any:
    return [b, a]
end

let result = swap("hello", 123)
return result
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val, ok := result.(*common.ClassInstance)
	if !ok {
		t.Fatalf("Expected ClassInstance got %T", result)
	}
	arr, err := engine.ClassInstanceToArray(val)
	if err != nil {
		t.Fatalf("Failed to convert to array: %v", err)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(arr))
	}

	// Use utils functions to handle both native types and ClassInstance
	intVal, ok := utils.AsInt(arr[0])
	if !ok || intVal != 123 {
		t.Errorf("Expected 123, got %v", arr[0])
	}
	if utils.ToString(arr[1]) != "hello" {
		t.Errorf("Expected 'hello', got %v", arr[1])
	}
}

func TestGenericFunction_MaxOf(t *testing.T) {
	code := `
def max_of<T>(a: T, b: T) -> T:
    if a > b:
        return a
    end
    return b
end

let result1 = max_of(10, 20)
let result2 = max_of(3.5, 2.1)
return [result1, result2]
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val, ok := result.(*common.ClassInstance)
	if !ok {
		t.Fatalf("Expected ClassInstance got %T", result)
	}
	arr, err := engine.ClassInstanceToArray(val)
	if err != nil {
		t.Fatalf("Failed to convert to array: %v", err)
	}

	if len(arr) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(arr))
	}

	if arr[0].(int) != 20 {
		t.Errorf("Expected 20, got %v", arr[0])
	}
	if arr[1].(float64) != 3.5 {
		t.Errorf("Expected 3.5, got %v", arr[1])
	}
}

package e2e

import (
	"testing"

	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func TestGenericClass_Builder(t *testing.T) {
	code := `
class Builder<T>:
    private var value: T
    
    Builder(initial: T):
        this.value = initial
    end
    
    def build() -> Any:
        return this.value
    end
    
    def with_value(new_value: T) -> Any:
        this.value = new_value
        return this
    end
end

let builder = Builder("initial")
builder.with_value("updated")
let result = builder.build()
return result
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "updated" {
		t.Errorf("Expected 'updated', got %v", str)
	}
}

func TestGenericClass_Pair(t *testing.T) {
	code := `
class Pair<K, V>:
    private var key: K
    private var value: V
    
    Pair(k: K, v: V):
        this.key = k
        this.value = v
    end
    
    def getKey() -> Any:
        return this.key
    end
    
    def getValue() -> Any:
        return this.value
    end
end

let pair = Pair("name", 42)
return [pair.getKey(), pair.getValue()]
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

	if arr[0] == nil {
		t.Fatalf("First element is nil")
	}

	if arr[1] == nil {
		t.Fatalf("Second element is nil")
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	if utils.ToString(arr[0]) != "name" {
		t.Errorf("Expected 'name', got %v", arr[0])
	}
	// Use utils.AsInt to handle both native ints and ClassInstance
	intVal, ok := utils.AsInt(arr[1])
	if !ok || intVal != 42 {
		t.Errorf("Expected 42, got %v", arr[1])
	}
}

func TestGenericClass_Box(t *testing.T) {
	code := `
class Box<T>:
    private var content: T
    
    Box(item: T):
        this.content = item
    end
    
    def get() -> Any:
        return this.content
    end
    
    def set(item: T):
        this.content = item
    end
end

let box = Box(100)
box.set(200)
let result = box.get()
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

	if value != 200 {
		t.Errorf("Expected 200, got %v", value)
	}
}

func TestGenericClass_Container(t *testing.T) {
	code := `
class Container<T>:
    private var item: T
    
    Container(value: T):
        this.item = value
    end
    
    def unwrap() -> Any:
        return this.item
    end
    
    def wrap(new_item: T):
        this.item = new_item
    end
end

let container1 = Container("hello")
let result1 = container1.unwrap()

let container2 = Container(42)
container2.wrap(99)
let result2 = container2.unwrap()

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

	// Use utils functions to handle both native types and ClassInstance
	if utils.ToString(arr[0]) != "hello" {
		t.Errorf("Expected 'hello', got %v", arr[0])
	}
	intVal, ok := utils.AsInt(arr[1])
	if !ok || intVal != 99 {
		t.Errorf("Expected 99, got %v", arr[1])
	}
}

func TestGenericClass_TripleHolder(t *testing.T) {
	code := `
class Triple<A, B, C>:
    private var first: A
    private var second: B
    private var third: C
    
    Triple(a: A, b: B, c: C):
        this.first = a
        this.second = b
        this.third = c
    end
    
    def getFirst() -> Any:
        return this.first
    end
    
    def getSecond() -> Any:
        return this.second
    end
    
    def getThird() -> Any:
        return this.third
    end
end

let triple = Triple("hello", 42, true)
return [triple.getFirst(), triple.getSecond(), triple.getThird()]
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
		t.Fatalf("Expected 3 elements, got %d", len(arr))
	}

	// Use utils functions to handle both native types and ClassInstance
	if utils.ToString(arr[0]) != "hello" {
		t.Errorf("Expected 'hello', got %v", arr[0])
	}
	intVal, ok := utils.AsInt(arr[1])
	if !ok || intVal != 42 {
		t.Errorf("Expected 42, got %v", arr[1])
	}
	boolVal := utils.AsBool(arr[2])
	if boolVal != true {
		t.Errorf("Expected true, got %v", arr[2])
	}
}

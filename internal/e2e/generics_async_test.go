package e2e

import (
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

func TestGenerics_ListWithType(t *testing.T) {
	code := `
let list = List<Int>()
list.add(1)
list.add(2)
list.add(3)
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

func TestGenerics_SetWithType(t *testing.T) {
	code := `
let set = Set<String>("hello", "world", "hello")
return set.size()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	size, _ := utils.AsInt(result)
	if size != 2 {
		t.Fatalf("Expected size 2 (duplicates removed), got %d", size)
	}
}

func TestGenerics_MapWithTypes(t *testing.T) {
	code := `
let map = Map<String, Int>()
map.put("one", 1)
map.put("two", 2)
return map.size()
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

func TestAsyncAwait_BasicPromise(t *testing.T) {
	code := `
let promise = async(() => 42)
return promise.await()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, ok := result.(int)
	if !ok || value != 42 {
		t.Fatalf("Expected 42, got %v", result)
	}
}

func TestAsyncAwait_PromiseWithThen(t *testing.T) {
	code := `
let promise = async(() => 10)
let result = nil
promise.then((value) => do
    result = value * 2
    return result
end)
Sys.sleep(100)
return result
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, ok := result.(float64)
	if !ok || value != 20 {
		t.Fatalf("Expected 20, got %v", result)
	}
}

func TestAsyncAwait_PromiseWithCatch(t *testing.T) {
	code := `
let promise = async(() => do
    throw RuntimeError("error occurred")
    return nil
end)

let errorMsg = nil
promise.catch((err) => do
    errorMsg = err
    return nil
end)

Sys.sleep(100)
return errorMsg
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("Expected error message, got nil")
	}

	msg, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string error message, got %T", result)
	}

	if !strings.Contains(msg, "error") {
		t.Fatalf("Expected error message to contain 'error', got %s", msg)
	}
}

func TestAsyncAwait_CompletableFuture(t *testing.T) {
	code := `
let future = CompletableFuture()

thread spawn do
    Sys.sleep(50)
    future.complete(100)
end

return future.get()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, ok := result.(int)
	if !ok || value != 100 {
		t.Fatalf("Expected 100, got %v", result)
	}
}

func TestAsyncAwait_CompletableFutureTimeout(t *testing.T) {
	code := `
let future = CompletableFuture()

thread spawn do
    Sys.sleep(200)
    future.complete(42)
end

try
    return future.getTimeout(50)
catch e
    return "timeout"
end
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	msg, ok := result.(string)
	if !ok || msg != "timeout" {
		t.Fatalf("Expected 'timeout', got %v", result)
	}
}

func TestAsyncAwait_PromiseChaining(t *testing.T) {
	code := `
let result = nil
async(() => 5)
    .then((val) => val * 2)
    .then((val) => do
        result = val + 10
        return result
    end)

Sys.sleep(100)
return result
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, ok := result.(float64)
	if !ok || value != 20 {
		t.Fatalf("Expected 20 (5*2+10), got %v", result)
	}
}

func runCode(code string) (any, error) {
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

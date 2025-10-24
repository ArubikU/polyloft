package e2e

import (
	"strings"
	"testing"
	
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func TestFunctionType_WithExplicitTypes(t *testing.T) {
	code := `
def integerToString(x: Int) -> String:
    return "a"
end
return Sys.type(integerToString)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Int,String>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestFunctionType_InferredReturnType(t *testing.T) {
	code := `
def stringToN(x: String):
    return 1
end
return Sys.type(stringToN)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<String,Int>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestFunctionType_UntypedParam(t *testing.T) {
	code := `
def aToN(x):
    return 1
end
return Sys.type(aToN)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Any,Int>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestFunctionType_TooManyReturnTypes(t *testing.T) {
	code := `
def complexFunc(a: Int, b):
    if a > 0:
        return true
    elif a < 0:
        return 1
    elif a == 0:
        return "zero"
    elif a == 2:
        return 2.0
    elif a == 3:
        return [1,2,3]
    else:
        return {"key": "value"}
    end
end
return Sys.type(complexFunc)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Int,Any,Any>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestLambdaType_BlockBody(t *testing.T) {
	code := `
let intToString = (x: Int) => do 
    return "a"
end
return Sys.type(intToString)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Int,String>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestLambdaType_ExpressionBody(t *testing.T) {
	code := `
let simpleLambda = (x: Int) => x * 2
return Sys.type(simpleLambda)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	// Expression returns a number (could be int or float after multiplication)
	if !strings.Contains(typeStr, "Function<Int,") {
		t.Fatalf("Expected Function<Int,...>, got %s", typeStr)
	}
}

func TestFunctionType_NoParams(t *testing.T) {
	code := `
def noParams():
    return "hello"
end
return Sys.type(noParams)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<String>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestFunctionType_MultipleParams(t *testing.T) {
	code := `
def multiParams(a: Int, b: String):
    return true
end
return Sys.type(multiParams)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Int,String,Bool>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestLambdaType_CanBePassedAsArgument(t *testing.T) {
	code := `
let list = List(1, 2, 3)
let sum = 0
list.forEach((x) => do
    sum = sum + x
end)
return sum
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, ok := asInt(result)
	if !ok || value != 6 {
		t.Fatalf("Expected 6, got %v", result)
	}
}

func TestFunctionType_VariadicParam(t *testing.T) {
	code := `
def sum(numbers: Int...):
    let total = 0
    for n in numbers:
        total = total + n
    end
    return total
end
return Sys.type(sum)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<Int...,Any>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestFunctionType_VariadicWithOtherParams(t *testing.T) {
	code := `
def formatMessage(prefix: String, parts: String...):
    return prefix
end
return Sys.type(formatMessage)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<String,String...,Any>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

func TestLambdaType_Variadic(t *testing.T) {
	code := `
let printer = (items: String...) => do
    return nil
end
return Sys.type(printer)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	typeStr, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string, got %T", result)
	}

	expected := "Function<String...,Nil>"
	if typeStr != expected {
		t.Fatalf("Expected %s, got %s", expected, typeStr)
	}
}

// Helper function to convert value to int
func asInt(v any) (int, bool) {
	// Use utils.AsInt which handles ClassInstances
	return utils.AsInt(v)
}

// runCode is defined in generics_async_test.go

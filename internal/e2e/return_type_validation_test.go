package e2e

import (
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func TestReturnTypeValidation_ConcreteType(t *testing.T) {
	code := `
class TestClass:
    TestClass():
    end
    
    def getString() -> String:
        return "hello"
    end
    
    def getInt() -> Int:
        return 42
    end
end

let obj = TestClass()
return [obj.getString(), obj.getInt()]
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

	str0 := utils.ToString(arr[0])
	int1, _ := utils.AsInt(arr[1])
	if str0 != "hello" {
		t.Errorf("Expected 'hello', got %v", str0)
	}
	if int1 != 42 {
		t.Errorf("Expected 42, got %v", int1)
	}
}

func TestReturnTypeValidation_WrongType(t *testing.T) {
	code := `
class TestClassWrongType:
    TestClassWrongType():
    end
    
    def getString() -> String:
        return 42  // Wrong type!
    end
end

let obj = TestClassWrongType()
return obj.getString()
`
	_, err := runCode(code)
	if err == nil {
		t.Fatalf("Expected type error, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "type") && !strings.Contains(errMsg, "Type") && !strings.Contains(errMsg, "expected") {
		t.Errorf("Expected type error, got: %v", err)
	}
}

func TestReturnTypeValidation_VoidMethod(t *testing.T) {
	code := `
class TestClassVoid:
    private var value: Int
    
    TestClassVoid():
        this.value = 0
    end
    
    def setValue(v: Int):
        this.value = v
    end
    
    def getValue() -> Int:
        return this.value
    end
end

let obj = TestClassVoid()
obj.setValue(100)
return obj.getValue()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value, _ := utils.AsInt(result)
	if value != 100 {
		t.Errorf("Expected 100, got %v", value)
	}
}

func TestReturnTypeValidation_GenericType(t *testing.T) {
	code := `
class BoxGeneric<T>:
    private var content: T
    
    BoxGeneric(item: T):
        this.content = item
    end
    
    def get() -> T:
        return this.content
    end
    
    def set(item: T):
        this.content = item
    end
end

let box1 = BoxGeneric("hello")
let box2 = BoxGeneric(42)

return [box1.get(), box2.get()]
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

	str0 := utils.ToString(arr[0])
	int1, _ := utils.AsInt(arr[1])
	if str0 != "hello" {
		t.Errorf("Expected 'hello', got %v", str0)
	}
	if int1 != 42 {
		t.Errorf("Expected 42, got %v", int1)
	}
}

func TestReturnTypeValidation_AnyType(t *testing.T) {
	code := `
class TestClassAny:
    TestClassAny():
    end
    
    def getAny(flag: Any) -> Any:
        if flag:
            return "string"
        end
        return 42
    end
end

let obj = TestClassAny()
return [obj.getAny(true), obj.getAny(false)]
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

	str0 := utils.ToString(arr[0])
	int1, _ := utils.AsInt(arr[1])
	if str0 != "string" {
		t.Errorf("Expected 'string', got %v", str0)
	}
	if int1 != 42 {
		t.Errorf("Expected 42, got %v", int1)
	}
}

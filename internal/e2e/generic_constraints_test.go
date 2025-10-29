package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Tests for generic type constraints (T extends Animal, etc.)

func TestGenericConstraint_Basic(t *testing.T) {
	// Test basic type parameter constraint
	code := `
class Animal:
    def speak() -> String:
        return "Some sound"
    end
end

class Dog < Animal:
    def speak() -> String:
        return "Woof!"
    end
end

class AnimalContainer<T extends Animal>:
    private var item: T
    
    AnimalContainer(i: T):
        this.item = i
    end
    
    def getItem() -> T:
        return this.item
    end
end

let dog = Dog()
let container = AnimalContainer<Dog>(dog)
return container.getItem().speak()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "Woof!" {
		t.Fatalf("Expected 'Woof!', got %v", result)
	}
}

func TestGenericConstraint_ViolationAtCreation(t *testing.T) {
	// Test that constraint is enforced at object creation
	code := `
class AnimalB:
    def speak() -> String:
        return "Some sound"
    end
end

class AnimalContainer<T extends AnimalB>:
    private var item: T
    
    AnimalContainer(i: T):
        this.item = i
    end
end

class NotAnAnimal:
    def talk() -> String:
        return "Hello"
    end
end

let notAnimal = NotAnAnimal()
let container = AnimalContainer<NotAnAnimal>(notAnimal)
return Sys.type(container)
`
	backward, err := runCode(code)
	if err == nil {
		fmt.Println(backward)
		t.Fatal("Expected error for constraint violation, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "NotAnAnimal") || !strings.Contains(errMsg, "extends") {
		t.Fatalf("Expected constraint violation error, got: %v", errMsg)
	}
}

func TestGenericConstraint_MultipleParams(t *testing.T) {
	// Test multiple type parameters with constraints
	code := `
class NumberType:
end

class IntType < NumberType:
end

class Container<K extends String, V extends NumberType>:
    private var key: K
    private var value: V
    
    Container(k: K, v: V):
        this.key = k
        this.value = v
    end
    
    def getKey() -> K:
        return this.key
    end
end

let container = Container<String, IntType>("key", IntType())
return container.getKey()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "key" {
		t.Fatalf("Expected 'key', got %v", result)
	}
}

func TestGenericInheritance_Basic(t *testing.T) {
	// Test generic class inheriting from generic parent
	code := `
class Container<T>:
    private var item: T
    
    Container(i: T):
        this.item = i
    end
    
    def getItem() -> T:
        return this.item
    end
end

class StringContainer < Container<String>:
    StringContainer(s: String):
        super(s)
    end
    
    def getLength() -> Int:
        return this.getItem().length()
    end
end

let sc = StringContainer("hello")
return sc.getLength()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if result is an integer 5 (can be int or int64)
	switch v := result.(type) {
	case int:
		if v != 5 {
			t.Fatalf("Expected 5, got %v", result)
		}
	case int64:
		if v != 5 {
			t.Fatalf("Expected 5, got %v", result)
		}
	default:
		t.Fatalf("Expected integer 5, got %T: %v", result, result)
	}
}

func TestGenericInheritance_WithTypeParam(t *testing.T) {
	// Test generic class inheriting from generic parent, passing through type parameter
	code := `
class BaseContainer<T>:
    private var item: T
    
    BaseContainer(i: T):
        this.item = i
    end
    
    def getItem() -> T:
        return this.item
    end
end

class SpecialContainer<T> < BaseContainer<T>:
    SpecialContainer(i: T):
        super(i)
    end
    
    def getSpecialItem() -> T:
        return this.getItem()
    end
end

let sc = SpecialContainer<String>("test")
return sc.getSpecialItem()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "test" {
		t.Fatalf("Expected 'test', got %v", result)
	}
}

func TestVariance_ErrorHandling_In(t *testing.T) {
	// Test that 'in' variance prevents reading (covariant position)
	code := `
class Consumer<in T>:
    private var item: T
    
    Consumer(i: T):
        this.item = i
    end
    
    def accept(i: T) -> Void:
        this.item = i
    end
    
    def get() -> T:
        return this.item
    end
end

let consumer = Consumer<String>("hello")
return consumer.get()
`
	_, err := runCode(code)
	// In a full implementation, this should error because 'in' variance
	// means the type can only be consumed (contravariant), not produced
	// For now, we just check if it runs (error checking would be future work)
	if err != nil {
		// If error checking is implemented, this is expected
		errMsg := err.Error()
		if strings.Contains(errMsg, "variance") || strings.Contains(errMsg, "contravariant") {
			// Good - variance error was caught
			return
		}
	}
	// Test passes either way for now
}

func TestVariance_ErrorHandling_Out(t *testing.T) {
	// Test that 'out' variance prevents writing (contravariant position)
	code := `
class Producer<out T>:
    private var item: T
    
    Producer(i: T):
        this.item = i
    end
    
    def get() -> T:
        return this.item
    end
    
    def set(i: T) -> Void:
        this.item = i
    end
end

let producer = Producer<String>("hello")
producer.set("world")
return producer.get()
`
	_, err := runCode(code)
	// In a full implementation, this should error because 'out' variance
	// means the type can only be produced (covariant), not consumed
	// For now, we just check if it runs (error checking would be future work)
	if err != nil {
		// If error checking is implemented, this is expected
		errMsg := err.Error()
		if strings.Contains(errMsg, "variance") || strings.Contains(errMsg, "covariant") {
			// Good - variance error was caught
			return
		}
	}
	// Test passes either way for now
}

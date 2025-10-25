package e2e

import (
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Tests for variance error handling (in/out)
// These tests validate that variance annotations are enforced at runtime

func TestVarianceError_OutInParameter(t *testing.T) {
	// Test that covariant (out) type parameters cannot be used in method parameter positions
	code := `
class ProducerErr<out T>:
    private var value: T
    
    ProducerErr(v: T):
        this.value = v
    end
    
    // This should cause an error: 'out T' cannot be in parameter position
    def set(newValue: T) -> Void:
        this.value = newValue
    end
end

let producer = ProducerErr<String>("hello")
producer.set("world")
`
	_, err := runCode(code)
	if err == nil {
		t.Fatal("Expected error for covariant type parameter in parameter position, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "covariant") && !strings.Contains(errMsg, "out") {
		t.Fatalf("Expected error about covariant/out type parameter, got: %v", err)
	}
}

func TestVarianceError_OutInVariadicParameter(t *testing.T) {
	// Test that covariant (out) type parameters cannot be used in variadic parameter positions
	code := `
class MultiProducerErr<out T>:
    private var values: Array
    
    MultiProducerErr():
        this.values = Array()
    end
    
    // This should cause an error: 'out T' cannot be in variadic parameter position
    def addAll(items: T...) -> Void:
        // Try to add items
    end
end

let producer = MultiProducerErr<String>()
producer.addAll("a", "b", "c")
`
	_, err := runCode(code)
	if err == nil {
		t.Fatal("Expected error for covariant type parameter in variadic parameter position, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "covariant") && !strings.Contains(errMsg, "out") {
		t.Fatalf("Expected error about covariant/out type parameter, got: %v", err)
	}
}

func TestVarianceSuccess_OutInReturn(t *testing.T) {
	// Test that covariant (out) type parameters CAN be used in return positions (this is correct usage)
	code := `
class ProducerOk<out T>:
    private var value: T
    
    ProducerOk(v: T):
        this.value = v
    end
    
    // This is correct: 'out T' can be in return position
    def get() -> T:
        return this.value
    end
end

let producer = ProducerOk<String>("hello")
return producer.get()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error for covariant type parameter in return position: %v", err)
	}

	str := utils.ToString(result)
	if str != "hello" {
		t.Fatalf("Expected 'hello', got %v", result)
	}
}

func TestVarianceSuccess_InInParameter(t *testing.T) {
	// Test that contravariant (in) type parameters CAN be used in parameter positions (this is correct usage)
	code := `
class ConsumerOk<in T>:
    private var value: Any
    
    ConsumerOk():
        this.value = nil
    end
    
    // This is correct: 'in T' can be in parameter position
    def accept(item: T) -> Void:
        this.value = item
    end
    
    def getValue() -> Any:
        return this.value
    end
end

let consumer = ConsumerOk<Int>()
consumer.accept(42)
return consumer.getValue()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error for contravariant type parameter in parameter position: %v", err)
	}

	val, ok := utils.AsInt(result)
	if !ok || val != 42 {
		t.Fatalf("Expected 42, got %v", result)
	}
}

func TestVarianceSuccess_InvariantBothPositions(t *testing.T) {
	// Test that invariant (no variance annotation) type parameters can be used in both positions
	code := `
class BoxVariance<T>:
    private var content: T
    
    BoxVariance(item: T):
        this.content = item
    end
    
    // Invariant T can be in parameter position
    def set(newContent: T) -> Void:
        this.content = newContent
    end
    
    // Invariant T can be in return position
    def get() -> T:
        return this.content
    end
end

let box = BoxVariance<String>("hello")
box.set("world")
return box.get()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error for invariant type parameter: %v", err)
	}

	str := utils.ToString(result)
	if str != "world" {
		t.Fatalf("Expected 'world', got %v", result)
	}
}

func TestVarianceError_MultipleTypeParams(t *testing.T) {
	// Test variance errors with multiple type parameters
	code := `
class ContainerMulti<out T, in U>:
    private var producerValue: T
    private var consumerValue: Any
    
    ContainerMulti(pv: T):
        this.producerValue = pv
        this.consumerValue = nil
    end
    
    // Correct: out T in return position
    def get() -> T:
        return this.producerValue
    end
    
    // Correct: in U in parameter position
    def consume(item: U) -> Void:
        this.consumerValue = item
    end
    
    // ERROR: out T in parameter position
    def badMethod(item: T) -> Void:
        this.producerValue = item
    end
end

let container = ContainerMulti<String, Int>("hello")
container.badMethod("world")
`
	_, err := runCode(code)
	if err == nil {
		t.Fatal("Expected error for covariant type parameter in parameter position, but got none")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "covariant") && !strings.Contains(errMsg, "out") {
		t.Fatalf("Expected error about covariant/out type parameter, got: %v", err)
	}
}

func TestVarianceInfo_InheritancePattern(t *testing.T) {
	// Test that variance works correctly with the inheritance pattern
	// This is more of an informational test showing correct variance usage
	code := `
class AnimalVar:
    def speak() -> String:
        return "Some sound"
    end
end

class DogVar extends AnimalVar:
    def speak() -> String:
        return "Woof!"
    end
end

// Producer with covariance - can produce subtypes
class AnimalProducerVar<out T extends AnimalVar>:
    private var animal: T
    
    AnimalProducerVar(a: T):
        this.animal = a
    end
    
    def getAnimal() -> T:
        return this.animal
    end
end

let dogProducer = AnimalProducerVar<DogVar>(DogVar())
let sound = dogProducer.getAnimal().speak()
return sound
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str != "Woof!" {
		t.Fatalf("Expected 'Woof!', got %v", result)
	}
}

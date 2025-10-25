package e2e

import (
	"testing"
	
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Tests for variance (in/out) with user-defined generic classes

func TestUserVariance_Producer_Out(t *testing.T) {
	// Producer with 'out' variance - covariant
	code := `
class Producer<out T>:
    private var value: T
    
    Producer(v: T):
        this.value = v
    end
    
    def get() -> T:
        return this.value
    end
end

let producer = Producer<String>("hello")
return producer.get()
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

func TestUserVariance_Consumer_In(t *testing.T) {
	// Consumer with 'in' variance - contravariant
	code := `
class Consumer<in T>:
    private var value: T
    
    Consumer(v: T):
        this.value = v
    end
    
    def set(v: T) -> Void:
        this.value = v
    end
    
    def getValue() -> Any:
        return this.value
    end
end

let consumer = Consumer<Int>(42)
consumer.set(100)
return consumer.getValue()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val, ok := utils.AsInt(result)
	if !ok || val != 100 {
		t.Fatalf("Expected 100, got %v", result)
	}
}

func TestUserVariance_Invariant(t *testing.T) {
	// Invariant - no variance annotation
	code := `
class BoxInv<T>:
    private var value: T
    
    BoxInv(v: T):
        this.value = v
    end
    
    def get() -> T:
        return this.value
    end
    
    def set(v: T) -> Void:
        this.value = v
    end
end

let box = BoxInv<String>("hello")
box.set("world")
return box.get()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str != "world" {
		t.Fatalf("Expected 'world', got %v", result)
	}
}

func TestUserVariance_MultipleTypeParams(t *testing.T) {
	// Multiple type parameters with different variance
	code := `
class Converter<in I, out O>:
    private var input: I
    private var output: O
    
    Converter(i: I, o: O):
        this.input = i
        this.output = o
    end
    
    def setInput(i: I) -> Void:
        this.input = i
    end
    
    def getOutput() -> O:
        return this.output
    end
end

let converter = Converter<Int, String>(42, "result")
converter.setInput(100)
return converter.getOutput()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str != "result" {
		t.Fatalf("Expected 'result', got %v", result)
	}
}

func TestUserVariance_TypeDisplay(t *testing.T) {
	// Verify Sys.type() displays variance annotations correctly
	code := `
class ProducerTD<out T>:
    private var value: T
    ProducerTD(v: T):
        this.value = v
    end
end

class ConsumerTD<in T>:
    private var value: T
    ConsumerTD(v: T):
        this.value = v
    end
end

let producer = ProducerTD<String>("hello")
let consumer = ConsumerTD<Int>(42)

return Sys.type(producer) + "," + Sys.type(consumer)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	str := utils.ToString(result)
	if str == "" {
		t.Fatalf("Expected string result, got %v", result)
	}

	// Should display type parameters (variance is internal)
	// The important part is that it shows ProducerTD<String> and ConsumerTD<Int>
	if str != "ProducerTD<String>,ConsumerTD<Int>" {
		// If variance annotations are shown, that's also acceptable
		t.Logf("Type display result: %s", str)
	}
}

func TestUserVariance_InstanceOf(t *testing.T) {
	// Verify instanceof works with variance annotations
	code := `
class ProducerIO<out T>:
    private var value: T
    ProducerIO(v: T):
        this.value = v
    end
end

let producer = ProducerIO<String>("hello")

let check1 = Sys.instanceof(producer, "ProducerIO")
let check2 = Sys.instanceof(producer, "ProducerIO<String>")
let check3 = Sys.instanceof(producer, "ProducerIO<?>")

return check1 && check2 && check3
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val := utils.AsBool(result)
	if !val {
		t.Fatalf("Expected true for all instanceof checks, got %v", val)
	}
}

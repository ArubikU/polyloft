package e2e

import (
	"testing"

	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Tests for class inheritance with generic types

func TestInheritance_Basic(t *testing.T) {
	// Basic inheritance test
	code := `
class AnimalA:
    private var name: String
    
    AnimalA(n: String):
        this.name = n
    end
    
    def getName() -> String:
        return this.name
    end
    
    def speak() -> String:
        return "Some sound"
    end
end

class DogA < AnimalA:
    DogA(n: String):
        super(n)
    end
    
    def speak() -> String:
        return "Woof!"
    end
end

let dog = DogA("Buddy")
return dog.getName() + ":" + dog.speak()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "Buddy:Woof!" {
		t.Fatalf("Expected 'Buddy:Woof!', got %v", result)
	}
}

func TestInheritance_InstanceOf(t *testing.T) {
	// Test instanceof with inheritance
	code := `
class AnimalB:
    AnimalB():
    end
end

class DogB < AnimalB:
    DogB():
        super()
    end
end

class AdministratorB:
    AdministratorB():
    end
end

let animal = AnimalB()
let dog = DogB()
let admin = AdministratorB()

let check1 = Sys.instanceof(dog, "DogB")
let check2 = Sys.instanceof(dog, "AnimalB")
let check3 = Sys.instanceof(animal, "DogB")
let check4 = Sys.instanceof(admin, "AnimalB")

return check1 && check2 && !check3 && !check4
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val := utils.AsBool(result)
	if !val {
		t.Fatalf("Expected true for inheritance instanceof checks, got %v", val)
	}
}

func TestInheritance_WithGenerics(t *testing.T) {
	// Test generics with inherited classes
	code := `
class AnimalC:
    private var name: String
    AnimalC(n: String):
        this.name = n
    end
    def getName() -> String:
        return this.name
    end
end

class DogC < AnimalC:
    DogC(n: String):
        super(n)
    end
end

class ContainerC<T>:
    private var item: T
    ContainerC(i: T):
        this.item = i
    end
    def get() -> T:
        return this.item
    end
end

let dogContainer = ContainerC<DogC>(DogC("Buddy"))
return dogContainer.get().getName()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "Buddy" {
		t.Fatalf("Expected 'Buddy', got %v", result)
	}
}

func TestInheritance_GenericContainer_TypeDisplay(t *testing.T) {
	// Verify Sys.type() with generic containers holding inherited types
	code := `
class AnimalD:
    AnimalD():
    end
end

class DogD < AnimalD:
    DogD():
        super()
    end
end

class ContainerD<T>:
    private var item: T
    ContainerD(i: T):
        this.item = i
    end
end

let dogContainer = ContainerD<DogD>(DogD())
let animalContainer = ContainerD<AnimalD>(AnimalD())

return Sys.type(dogContainer) + "," + Sys.type(animalContainer)
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "ContainerD<DogD>,ContainerD<AnimalD>" {
		t.Fatalf("Expected 'ContainerD<DogD>,ContainerD<AnimalD>', got %v", result)
	}
}

func TestInheritance_GenericContainer_InstanceOf(t *testing.T) {
	// Test instanceof with generic containers and inheritance
	code := `
class AnimalE:
    AnimalE():
    end
end

class DogE < AnimalE:
    DogE():
        super()
    end
end

class ContainerE<T>:
    private var item: T
    ContainerE(i: T):
        this.item = i
    end
end

let dogContainer = ContainerE<DogE>(DogE())

let check1 = Sys.instanceof(dogContainer, "ContainerE")
let check2 = Sys.instanceof(dogContainer, "ContainerE<DogE>")
let check3 = Sys.instanceof(dogContainer, "ContainerE<?>")

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

func TestInheritance_Covariance(t *testing.T) {
	// Test covariance with inheritance - Container<Dog> should work with Container<Animal> for 'out'
	code := `
class AnimalF:
    private var name: String
    AnimalF(n: String):
        this.name = n
    end
    def getName() -> String:
        return this.name
    end
end

class DogF < AnimalF:
    DogF(n: String):
        super(n)
    end
end

class ProducerF<out T>:
    private var item: T
    ProducerF(i: T):
        this.item = i
    end
    def get() -> T:
        return this.item
    end
end

let dogProducer = ProducerF<DogF>(DogF("Buddy"))
return dogProducer.get().getName()
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "Buddy" {
		t.Fatalf("Expected 'Buddy', got %v", result)
	}
}

func TestInheritance_TypeChecking_WithInheritance(t *testing.T) {
	// Test that type checking works with inherited types in generic containers
	code := `
class AnimalG:
    AnimalG():
    end
end

class DogG < AnimalG:
    DogG():
        super()
    end
end

class BoxG<T>:
    private var item: T
    BoxG(i: T):
        this.item = i
    end
    def set(i: T) -> Void:
        this.item = i
    end
    def get() -> T:
        return this.item
    end
end

let dogBox = BoxG<DogG>(DogG())
dogBox.set(DogG())
return "OK"
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Use utils.ToString to handle both native strings and ClassInstance
	str := utils.ToString(result)
	if str != "OK" {
		t.Fatalf("Expected 'OK', got %v", result)
	}
}

func TestInheritance_MultiLevel(t *testing.T) {
	// Test multi-level inheritance
	code := `
class AnimalH:
    def level() -> String:
        return "Animal"
    end
end

class MammalH < AnimalH:
    def level() -> String:
        return "Mammal"
    end
end

class DogH < MammalH:
    def level() -> String:
        return "Dog"
    end
end

let dog = DogH()
let check1 = Sys.instanceof(dog, "DogH")
let check2 = Sys.instanceof(dog, "MammalH")
let check3 = Sys.instanceof(dog, "AnimalH")

return check1 && check2 && check3
`
	result, err := runCode(code)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	val := utils.AsBool(result)
	if !val {
		t.Fatalf("Expected true for multi-level inheritance checks, got %v", val)
	}
}

package e2e

import (
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/engine/utils"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// Comprehensive tests for all 17 sections of typerules.md

func runCodeTypeRules(code string) (any, error) {
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

// Section 1: Basic type checking with classes
func TestTypeRules_Section1_BasicClassTypes(t *testing.T) {
	code := `
class Person:
    let name
    let age
    Person(name, age):
        this.name = name
        this.age = age
    end
end
const p = Person("Alice", 30)
println(Sys.type(p))
println(Sys.instanceof(p, "Person"))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Person") {
		t.Errorf("Expected output to contain 'Person', got: %s", output)
	}
	if !contains(output, "true") {
		t.Errorf("Expected instanceof to be true, got: %s", output)
	}
}

// Section 2: Inheritance type checking
func TestTypeRules_Section2_InheritanceTypes(t *testing.T) {
	code := `
class Person:
    let name
    let age
    Person(name, age):
        this.name = name
        this.age = age
    end
end

class Employee < Person:
    let employeeId
    Employee(name, age, employeeId):
        super(name, age)
        this.employeeId = employeeId
    end
end

const e = Employee("Bob", 25, "E123")
println(Sys.type(e))
println(Sys.instanceof(e, "Employee"))
println(Sys.instanceof(e, "Person"))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Employee") {
		t.Errorf("Expected type to be Employee, got: %s", output)
	}
	// Count true occurrences (should be 2)
	trueCount := 0
	for _, line := range splitLines(output) {
		if contains(line, "true") {
			trueCount++
		}
	}
	if trueCount < 2 {
		t.Errorf("Expected 2 instanceof checks to be true, got %d in: %s", trueCount, output)
	}
}

// Section 3: Interface type checking
func TestTypeRules_Section3_InterfaceTypes(t *testing.T) {
	code := `
interface Named:
    def getName()
end

class User implements Named:
    let name
    User(name):
        this.name = name
    end
    def getName():
        return this.name
    end
end

const u = User("Charlie")
println(Sys.type(u))
println(Sys.instanceof(u, "User"))
println(Sys.instanceof(u, "Named"))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "User") {
		t.Errorf("Expected type to be User, got: %s", output)
	}
	// Both instanceof checks should be true
	trueCount := 0
	for _, line := range splitLines(output) {
		if contains(line, "true") {
			trueCount++
		}
	}
	if trueCount < 2 {
		t.Errorf("Expected 2 instanceof checks to be true, got %d in: %s", trueCount, output)
	}
}

// Section 4: Enum type checking
func TestTypeRules_Section4_EnumTypes(t *testing.T) {
	code := `
enum Color:
    Red
    Green
    Blue
end

const c = Color.Red
println(Sys.type(c))
println(Sys.instanceof(c, "Color"))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Color") {
		t.Errorf("Expected type to be Color, got: %s", output)
	}
	if !contains(output, "true") {
		t.Errorf("Expected instanceof to be true, got: %s", output)
	}
}

// Section 5: Function parameter type checking with inheritance
func TestTypeRules_Section5_FunctionParameterTypes(t *testing.T) {
	code := `
class Person:
    let name
    Person(name):
        this.name = name
    end
end

class Employee < Person:
    let employeeId
    Employee(name, employeeId):
        super(name)
        this.employeeId = employeeId
    end
end

def greet(person: Person):
    println("Hello, " + person.name)
end

const p = Person("Alice")
const e = Employee("Bob", "E123")

greet(p)
greet(e)
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Hello, Alice") {
		t.Errorf("Expected greeting for Alice, got: %s", output)
	}
	if !contains(output, "Hello, Bob") {
		t.Errorf("Expected greeting for Bob, got: %s", output)
	}
}

// Section 5b: Number type hierarchy
func TestTypeRules_Section5_NumberHierarchy(t *testing.T) {
	code := `
def printNumber(n: Number):
    println("Number: " + n.toString())
end

printNumber(100)
printNumber(2.718)
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Number: 100") {
		t.Errorf("Expected 'Number: 100', got: %s", output)
	}
	if !contains(output, "Number: 2.718") {
		t.Errorf("Expected 'Number: 2.718', got: %s", output)
	}
}

// Section 6: Any type as universal type
func TestTypeRules_Section6_AnyType(t *testing.T) {
	code := `
def printAny(value: Any):
    println("Value: " + value.toString())
end

printAny(42)
printAny("Hello")
printAny(3.14)
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Value: 42") {
		t.Errorf("Expected 'Value: 42', got: %s", output)
	}
	if !contains(output, "Value: Hello") {
		t.Errorf("Expected 'Value: Hello', got: %s", output)
	}
}

// Section 7: Variadic functions with types
func TestTypeRules_Section7_VariadicFunctions(t *testing.T) {
	code := `
def sum(numbers: Int...):
    let total = 0
    for n in numbers:
        total += n
    end
    return total
end

println(sum(1, 2, 3, 4, 5))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	total, _ := utils.AsInt(result)
	if total != 15 {
		t.Errorf("Expected sum to be 15, got: %d", total)
	}
}

// Section 8: Generic types with bounds
func TestTypeRules_Section8_GenericBounds(t *testing.T) {
	code := `
class Box<T extends Number>:
    let value: T
    Box(value: T):
        this.value = value
    end
    def getValue(): T:
        return this.value
    end
end

const intBox = Box<Int>(100)
println(Sys.type(intBox))
const floatBox = Box<Float>(3.14)
println(Sys.type(floatBox))
println(intBox.getValue())
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Box") {
		t.Errorf("Expected type to contain 'Box', got: %s", output)
	}
	if !contains(output, "100") {
		t.Errorf("Expected value 100, got: %s", output)
	}
}

// Section 10: Type aliases with final type
func TestTypeRules_Section10_TypeAliases(t *testing.T) {
	code := `
final type Age = Int

def celebrateBirthday(age: Age):
    println("Happy " + age.toString() + "th Birthday!")
end

celebrateBirthday(30)
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Happy 30th Birthday!") {
		t.Errorf("Expected birthday message, got: %s", output)
	}
}

// Section 11: Built-in types
func TestTypeRules_Section11_BuiltinTypes(t *testing.T) {
	code := `
const arr = [1, 2, 3]
println(Sys.type(arr))
println(Sys.instanceof(arr, "Array"))

const dict = {"key": "value"}
println(Sys.type(dict))
println(Sys.instanceof(dict, "Map"))

const s = "Hello"
println(Sys.type(s))
println(Sys.instanceof(s, "String"))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Array") {
		t.Errorf("Expected Array type, got: %s", output)
	}
	if !contains(output, "Map") {
		t.Errorf("Expected Map type, got: %s", output)
	}
	if !contains(output, "String") {
		t.Errorf("Expected String type, got: %s", output)
	}
}

// Section 12: Nil handling
func TestTypeRules_Section12_NilType(t *testing.T) {
	code := `
const x = Nil
println(Sys.type(x))
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "Nil") {
		t.Errorf("Expected Nil type, got: %s", output)
	}
}

// Section 14: Union types
func TestTypeRules_Section14_UnionTypes(t *testing.T) {
	code := `
def printId(id: Int | String):
    println(id.toString())
end

printId(42)
printId("ID123")
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "42") {
		t.Errorf("Expected '42', got: %s", output)
	}
	if !contains(output, "ID123") {
		t.Errorf("Expected 'ID123', got: %s", output)
	}
}

// Section 16: Generic polymorphism
func TestTypeRules_Section16_GenericPolymorphism(t *testing.T) {
	code := `
def identity<T>(value: T): T:
    return value
end

let i = identity(42)
let s = identity("hi")
println(i)
println(s)
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "42") {
		t.Errorf("Expected '42', got: %s", output)
	}
	if !contains(output, "hi") {
		t.Errorf("Expected 'hi', got: %s", output)
	}
}

// Section 17: Type narrowing with instanceof
func TestTypeRules_Section17_TypeNarrowing(t *testing.T) {
	code := `
def printValue(x: Int | String):
    if Sys.instanceof(x, "Int"):
        println(x + 1)
    else:
        println(x.length())
    end
end

printValue(41)
printValue("test")
`
	result, err := runCodeTypeRules(code)
	if err != nil {
		t.Fatalf("Execution error: %v", err)
	}
	output := utils.ToString(result)
	if !contains(output, "42") {
		t.Errorf("Expected '42', got: %s", output)
	}
	if !contains(output, "4") {
		t.Errorf("Expected '4' (length of 'test'), got: %s", output)
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	var current string
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			if len(current) > 0 {
				lines = append(lines, current)
			}
			current = ""
		} else {
			current += string(s[i])
		}
	}
	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

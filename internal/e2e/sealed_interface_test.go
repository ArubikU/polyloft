package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

func TestEval_SealedInterface(t *testing.T) {
	src := `
// Basic sealed interface with no permit list (same package only)
sealed interface Repository
    static var defaultTimeout: Int = 5000
    
    def save(data: any) -> Bool
    def load(id: String) -> any
end

// Class implementing sealed interface in same package - should work
class LocalRepository implements Repository
    def save(data: any) -> Bool:
        return true
    end
    
    def load(id: String) -> any:
        return "Data"
    end
end

// Access static fields
println(Repository.defaultTimeout)

// Modify static field
Repository.defaultTimeout = 10000
println(Repository.defaultTimeout)

let repo = LocalRepository()
println(repo.save("test"))
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	buf := &bytes.Buffer{}
	_, err = engine.Eval(prog, engine.Options{Stdout: buf})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	got := buf.String()

	// Check for expected outputs
	if !strings.Contains(got, "5000") {
		t.Errorf("expected to see initial timeout value '5000', got: %s", got)
	}
	if !strings.Contains(got, "10000") {
		t.Errorf("expected to see modified timeout value '10000', got: %s", got)
	}
	if !strings.Contains(got, "true") {
		t.Errorf("expected to see 'true' from save(), got: %s", got)
	}
}

func TestEval_SealedInterfaceWithPermitList(t *testing.T) {
	src := `
// Sealed interface with permit list
sealed interface Drawable(Circle, Rectangle)
    static var defaultColor: String = "black"
    static var renderCount: Int = 0
    
    def draw() -> Void
    def area() -> Float
end

// Permitted class - should work
class Circle implements Drawable
    var radius: Float
    
    Circle(r: Float):
        this.radius = r
    end
    
    def draw() -> Void:
        Drawable.renderCount = Drawable.renderCount + 1
        println("Drawing circle")
    end
    
    def area() -> Float:
        return 3.14 * this.radius * this.radius
    end
end

// Permitted class - should work
class Rectangle implements Drawable
    var width: Float
    var height: Float
    
    Rectangle(w: Float, h: Float):
        this.width = w
        this.height = h
    end
    
    def draw() -> Void:
        Drawable.renderCount = Drawable.renderCount + 1
        println("Drawing rectangle")
    end
    
    def area() -> Float:
        return this.width * this.height
    end
end

println(Drawable.defaultColor)
println(Drawable.renderCount)

let circle = Circle(5.0)
circle.draw()

let rect = Rectangle(10.0, 20.0)
rect.draw()

println(Drawable.renderCount)
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	buf := &bytes.Buffer{}
	_, err = engine.Eval(prog, engine.Options{Stdout: buf})
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	got := buf.String()

	// Check for expected outputs
	if !strings.Contains(got, "black") {
		t.Errorf("expected to see 'black' default color, got: %s", got)
	}
	if !strings.Contains(got, "Drawing circle") {
		t.Errorf("expected to see 'Drawing circle', got: %s", got)
	}
	if !strings.Contains(got, "Drawing rectangle") {
		t.Errorf("expected to see 'Drawing rectangle', got: %s", got)
	}
	if !strings.Contains(got, "2") {
		t.Errorf("expected to see final render count of '2', got: %s", got)
	}
}

func TestEval_SealedInterfaceViolation(t *testing.T) {
	src := `
// Sealed interface with permit list
sealed interface Lockable(SafeBox)
    def lock() -> Void
    def unlock() -> Void
end

// Permitted class - should work
class SafeBox implements Lockable
    def lock() -> Void:
        println("Locked")
    end
    
    def unlock() -> Void:
        println("Unlocked")
    end
end

// Non-permitted class - should fail
class Vault implements Lockable
    def lock() -> Void:
        println("Vault locked")
    end
    
    def unlock() -> Void:
        println("Vault unlocked")
    end
end
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	buf := &bytes.Buffer{}
	_, err = engine.Eval(prog, engine.Options{Stdout: buf})

	// This should fail because Vault is not in the permit list
	if err == nil {
		t.Fatalf("expected error when non-permitted class implements sealed interface, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "sealed") && !strings.Contains(errMsg, "permit") {
		t.Errorf("expected error message about sealed interface violation, got: %s", errMsg)
	}
}

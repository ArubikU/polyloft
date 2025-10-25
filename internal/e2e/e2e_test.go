package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

func TestEval_Basics(t *testing.T) {
	src := `
let x = 2 + 3 * 4
let y = [1,2,3][1]
let m = { foo: "bar" }
println(x, y, m.foo)
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
	if got == "" {
		t.Fatalf("expected output, got empty")
	}
}

func TestEval_EnumAndRecord(t *testing.T) {
	src := `
enum Color
    RED
    GREEN
    BLUE
end

record Point(x: Int, y: Int)
    def sum():
        return this.x + this.y
    end
end

let green = Color.valueOf("GREEN")
println(green.name, green.ordinal)

let p = Point(2, 5)
println(p.x, p.y, p.sum(), Sys.type(p), p instanceof Point)

enum Planet
	MERCURY(3.7)
	MARS(3.71)

	var gravity: Float
	Planet(g: Float):
		this.gravity = g
	end

	def weight(mass: Float):
		println("Weight",this.gravity * mass)
		return mass * this.gravity
	end
end

println(Planet.MARS.gravity, Planet.MARS.weight(10.0))
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
	got := strings.TrimSpace(buf.String())
	expected := "GREEN 1\n2 5 7 Point true\n3.71 37.1"
	if got != expected {
		t.Fatalf("unexpected output\nwant %q\n got %q", expected, got)
	}
}

func TestEval_SealedClassPermits(t *testing.T) {
	src := `
sealed class Animal(Dog)
end

class Dog < Animal
end

class Cat < Animal
end
`

	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	_, err = engine.Eval(prog, engine.Options{})
	if err == nil {
		t.Fatalf("expected sealed inheritance error, got nil")
	}
	if !strings.Contains(err.Error(), "does not permit Cat") {
		t.Fatalf("expected sealed permit error mentioning Cat, got %v", err)
	}
}

func TestEval_SealedEnumHelpers(t *testing.T) {
	src := `
sealed enum Mode
    OFF
    ON
end

let names = Mode.names()
let values = Mode.values()
println(names[0], values[1].name, Mode.size())
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
	got := strings.TrimSpace(buf.String())
	const expected = "OFF ON 2"
	if got != expected {
		t.Fatalf("unexpected output for sealed enum helpers\nwant %q\n got %q", expected, got)
	}
}

func TestEval_ClassTypeAndInstanceOf(t *testing.T) {
	src := `
class Creature
end

class Canine < Creature
end

// Test Sys.type on class constructor shows "Class {Name}@{Package}"
println(Sys.type(Creature))

// Test instanceof with class as second parameter (direct)
let creature = Creature()
println(Sys.instanceof(creature, Creature))

// Test instanceof with class as second parameter (inheritance)
let canine = Canine()
println(Sys.instanceof(canine, Canine))
println(Sys.instanceof(canine, Creature))

// Test instanceof still works with string parameter
println(Sys.instanceof(canine, "Canine"))
println(Sys.instanceof(canine, "Creature"))
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

	// Check that output contains "Class Creature@" (package name may vary)
	if !strings.Contains(got, "Class Creature@") {
		t.Fatalf("expected Sys.type to return 'Class Creature@...', got: %q", got)
	}

	// Check that all instanceof calls return true (on separate lines)
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 6 {
		t.Fatalf("expected at least 6 lines of output, got: %q", got)
	}

	// First line should be class type
	if !strings.HasPrefix(lines[0], "Class Creature@") {
		t.Fatalf("first line should be class type, got: %q", lines[0])
	}

	// Next 5 lines should all be "true"
	for i := 1; i <= 5; i++ {
		if lines[i] != "true" {
			t.Fatalf("line %d should be 'true', got: %q", i+1, lines[i])
		}
	}
}

func TestEval_EnumTypeFormatting(t *testing.T) {
	src := `
enum Status
    ACTIVE
    INACTIVE
end

// Test Sys.type on enum shows "Enum {Name}@{Package}"
println(Sys.type(Status))

// Test enum functionality still works
println(Status.ACTIVE)
println(Status.size())
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

	// Check that output contains "Enum Status@" (package name may vary)
	if !strings.Contains(got, "Enum Status@") {
		t.Fatalf("expected Sys.type to return 'Enum Status@...', got: %q", got)
	}

	// Check enum functionality works
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines of output, got: %q", got)
	}

	// Second line should be enum value
	if lines[1] != "Status.ACTIVE" {
		t.Fatalf("expected 'Status.ACTIVE', got: %q", lines[1])
	}

	// Third line should be size
	if lines[2] != "2" {
		t.Fatalf("expected '2', got: %q", lines[2])
	}
}

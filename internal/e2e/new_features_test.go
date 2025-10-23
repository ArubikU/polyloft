package e2e

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

func TestForWhere_BasicFiltering(t *testing.T) {
	src := `
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let sum = 0

for n in numbers where n > 5:
    sum = sum + n
end

println("Sum: " + sum.toString())
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
	// Sum of 6+7+8+9+10 = 40
	if !strings.Contains(got, "Sum: 40") {
		t.Errorf("expected sum to be 40, got: %s", got)
	}
}

func TestForWhere_WithModuloCondition(t *testing.T) {
	src := `
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evenSum = 0

for n in numbers where n % 2 == 0:
    evenSum = evenSum + n
end

println("Even sum: " + evenSum.toString())
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
	// Sum of 2+4+6+8+10 = 30
	if !strings.Contains(got, "Even sum: 30") {
		t.Errorf("expected even sum to be 30, got: %s", got)
	}
}

func TestForWhere_WithDestructuring(t *testing.T) {
	src := `
let map = { "a": 10, "b": 5, "c": 20, "d": 3 }
let sum = 0

for k, v in map where v > 5:
    sum = sum + v
end

println("Filtered sum: " + sum.toString())
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
	// Sum of values > 5: 10 + 20 = 30
	if !strings.Contains(got, "Filtered sum: 30") {
		t.Errorf("expected filtered sum to be 30, got: %s", got)
	}
}

func TestForWhere_ComplexCondition(t *testing.T) {
	src := `
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let count = 0

for n in numbers where (n > 3) && (n < 8):
    count = count + 1
end

println("Count: " + count.toString())
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
	// Numbers 4, 5, 6, 7 = 4 items
	if !strings.Contains(got, "Count: 4") {
		t.Errorf("expected count to be 4, got: %s", got)
	}
}

func TestMultiLineLambda_BasicExecution(t *testing.T) {
	src := `
let process = (x, y) => do
    let sum = x + y
    let product = x * y
    return sum + product
end

let result = process(3, 4)
println("Result: " + result.toString())
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
	// (3+4) + (3*4) = 7 + 12 = 19
	if !strings.Contains(got, "Result: 19") {
		t.Errorf("expected result to be 19, got: %s", got)
	}
}

func TestMultiLineLambda_WithLocalVariables(t *testing.T) {
	src := `
let calculator = (a, b, c) => do
    let sum = a + b + c
    let avg = sum / 3
    let floatd = avg * 2
    return floatd
end

let result = calculator(3, 6, 9)
println("Result: " + result.toString())
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
	// (3+6+9)/3 * 2 = 18/3 * 2 = 6 * 2 = 12
	if !strings.Contains(got, "Result: 12") {
		t.Errorf("expected result to be 12, got: %s", got)
	}
}

func TestMultiLineLambda_WithDefer(t *testing.T) {
	src := `
let withCleanup = (x) => do
    defer println("Cleanup")
    println("Processing: " + x.toString())
    return x * 2
end

let result = withCleanup(5)
println("Result: " + result.toString())
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
	if !strings.Contains(got, "Cleanup") {
		t.Errorf("defer should work in multi-line lambda, got: %s", got)
	}
	if !strings.Contains(got, "Result: 10") {
		t.Errorf("expected result to be 10, got: %s", got)
	}
}

func TestList_BasicOperations(t *testing.T) {
	src := `
let list = List(1, 2, 3)
list.add(4)
list.add(5)

println("Size: " + list.size().toString())
println("Get 0: " + list.get(0).toString())
println("Contains 3: " + list.contains(3).toString())
println("Contains 10: " + list.contains(10).toString())
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
	if !strings.Contains(got, "Size: 5") {
		t.Errorf("expected size 5, got: %s", got)
	}
	if !strings.Contains(got, "Get 0: 1") {
		t.Errorf("expected first element to be 1, got: %s", got)
	}
	if !strings.Contains(got, "Contains 3: true") {
		t.Errorf("expected contains 3 to be true, got: %s", got)
	}
	if !strings.Contains(got, "Contains 10: false") {
		t.Errorf("expected contains 10 to be false, got: %s", got)
	}
}

func TestList_HigherOrderFunctions(t *testing.T) {
	src := `
let list = List(1, 2, 3, 4, 5)

let sum = 0
list.forEach((x) => do
    sum = sum + x
end)
println("Sum: " + sum.toString())

let count = 0
list.forEach((x) => do
    count = count + 1
end)
println("Count: " + count.toString())
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
	if !strings.Contains(got, "Sum: 15") {
		t.Errorf("expected sum to be 15, got: %s", got)
	}
	if !strings.Contains(got, "Count: 5") {
		t.Errorf("expected count to be 5, got: %s", got)
	}
}

func TestSet_Uniqueness(t *testing.T) {
	src := `
let set = Set(1, 2, 2, 3, 3, 3, 4)
println("Size: " + set.size().toString())
println("Contains 3: " + set.contains(3).toString())

set.add(5)
set.add(5)
println("After adds, size: " + set.size().toString())

set.remove(3)
println("After remove, contains 3: " + set.contains(3).toString())
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
	if !strings.Contains(got, "Size: 4") {
		t.Errorf("expected initial size 4 (unique values), got: %s", got)
	}
	if !strings.Contains(got, "After adds, size: 5") {
		t.Errorf("expected size 5 after adding unique value, got: %s", got)
	}
	if !strings.Contains(got, "After remove, contains 3: false") {
		t.Errorf("expected contains 3 to be false after removal, got: %s", got)
	}
}

func TestDeque_BasicOperations(t *testing.T) {
	src := `
let deque = Deque()
deque.pushBack(1)
deque.pushBack(2)
deque.pushFront(0)
deque.pushFront(-1)

println("Size: " + deque.size().toString())
println("Front: " + deque.popFront().toString())
println("Back: " + deque.popBack().toString())
println("Peek front: " + deque.peekFront().toString())
println("Peek back: " + deque.peekBack().toString())
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
	if !strings.Contains(got, "Size: 4") {
		t.Errorf("expected size 4, got: %s", got)
	}
	if !strings.Contains(got, "Front: -1") {
		t.Errorf("expected front to be -1, got: %s", got)
	}
	if !strings.Contains(got, "Back: 2") {
		t.Errorf("expected back to be 2, got: %s", got)
	}
	// After pops: [-1, 0, 1, 2] -> [0, 1]
	if !strings.Contains(got, "Peek front: 0") {
		t.Errorf("expected peek front to be 0, got: %s", got)
	}
	if !strings.Contains(got, "Peek back: 1") {
		t.Errorf("expected peek back to be 1, got: %s", got)
	}
}

func TestCombined_ForWhereWithCollections(t *testing.T) {
	src := `
let list = List(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
let arr = list.toArray()

let evenSet = Set()
for n in arr where n % 2 == 0:
    evenSet.add(n)
end

println("Even set size: " + evenSet.size().toString())
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
	// Even numbers: 2, 4, 6, 8, 10 = 5
	if !strings.Contains(got, "Even set size: 5") {
		t.Errorf("expected even set size to be 5, got: %s", got)
	}
}

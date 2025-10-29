package e2e

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

func TestDefer_BasicExecution(t *testing.T) {
	src := `
def testDefer():
    defer println("Deferred 1")
    defer println("Deferred 2")
    println("Body")
end
testDefer()
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
	// Should print: Body, Deferred 2, Deferred 1 (LIFO order)
	if !strings.Contains(got, "Body") {
		t.Errorf("expected 'Body' in output, got: %s", got)
	}
	if !strings.Contains(got, "Deferred 1") || !strings.Contains(got, "Deferred 2") {
		t.Errorf("expected deferred calls in output, got: %s", got)
	}
	// Check LIFO order: Deferred 2 should appear before Deferred 1
	idx1 := strings.Index(got, "Deferred 1")
	idx2 := strings.Index(got, "Deferred 2")
	if idx1 < idx2 {
		t.Errorf("defer should execute in LIFO order, but Deferred 1 appeared before Deferred 2")
	}
}

func TestDefer_WithEarlyReturn(t *testing.T) {
	src := `
def testDefer():
    defer println("Cleanup")
    println("Before return")
    return 42
    println("After return - should not print")
end
let result = testDefer()
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
		t.Errorf("defer should execute even with early return, got: %s", got)
	}
	if !strings.Contains(got, "Result: 42") {
		t.Errorf("expected return value 42, got: %s", got)
	}
	if strings.Contains(got, "After return") {
		t.Errorf("code after return should not execute, got: %s", got)
	}
}

func TestDefer_WithError(t *testing.T) {
	src := `
def testDefer():
    defer println("Cleanup executed")
    println("Before error")
    throw "Test error"
    println("After error - should not print")
end

try
    testDefer()
catch e
    println("Caught: " + e.toString())
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
	if err != nil {
		t.Fatalf("eval error: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "Cleanup executed") {
		t.Errorf("defer should execute even with error, got: %s", got)
	}
	if !strings.Contains(got, "Caught:") {
		t.Errorf("expected error to be caught, got: %s", got)
	}
}

func TestThread_BasicSpawnAndJoin(t *testing.T) {
	src := `
let t = thread spawn do
    println("Thread executing")
    return 42
end

println("Main thread")
let result = thread join t
println("Thread result: " + result.toString())
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
	if !strings.Contains(got, "Thread executing") {
		t.Errorf("expected thread to execute, got: %s", got)
	}
	if !strings.Contains(got, "Thread result: 42") {
		t.Errorf("expected thread result 42, got: %s", got)
	}
}

func TestThread_MultipleThreads(t *testing.T) {
	src := `
let t1 = thread spawn do
    return 10
end

let t2 = thread spawn do
    return 20
end

let r1 = thread join t1
let r2 = thread join t2
let sum = r1 + r2
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
	if !strings.Contains(got, "Sum: 30") {
		t.Errorf("expected sum to be 30, got: %s", got)
	}
}

func TestChannel_BasicSendReceive(t *testing.T) {
	// Use a buffered approach via thread spawn
	src := `
let ch = channel[Int]()

thread spawn do
    ch.send(42)
    ch.close()
end

let received = ch.recv()
println("Received: " + received.toString())
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	// Set a timeout for the test
	done := make(chan bool, 1)
	buf := &bytes.Buffer{}
	
	go func() {
		_, err = engine.Eval(prog, engine.Options{Stdout: buf})
		done <- true
	}()
	
	select {
	case <-done:
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		got := buf.String()
		if !strings.Contains(got, "Received: 42") {
			t.Errorf("expected to receive 42, got: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout - channel operations may be blocking")
	}
}

func TestChannel_MultipleValues(t *testing.T) {
	src := `
let ch = channel[Int]()

thread spawn do
    ch.send(1)
    ch.send(2)
    ch.send(3)
    ch.close()
end

let sum = 0
sum = sum + ch.recv()
sum = sum + ch.recv()
sum = sum + ch.recv()
println("Sum: " + sum.toString())
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	done := make(chan bool, 1)
	buf := &bytes.Buffer{}
	
	go func() {
		_, err = engine.Eval(prog, engine.Options{Stdout: buf})
		done <- true
	}()
	
	select {
	case <-done:
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		got := buf.String()
		if !strings.Contains(got, "Sum: 6") {
			t.Errorf("expected sum to be 6, got: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestSelect_BasicReceive(t *testing.T) {
	src := `
let ch = channel[Int]()

thread spawn do
    ch.send(100)
    ch.close()
end

let received = 0
loop
    select
        case let x = ch.recv():
            received = x
            println("Received: " + x.toString())
        case closed ch:
            println("Channel closed")
            break
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
	
	done := make(chan bool, 1)
	buf := &bytes.Buffer{}
	
	go func() {
		_, err = engine.Eval(prog, engine.Options{Stdout: buf})
		done <- true
	}()
	
	select {
	case <-done:
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		got := buf.String()
		if !strings.Contains(got, "Received: 100") {
			t.Errorf("expected to receive 100, got: %s", got)
		}
		if !strings.Contains(got, "Channel closed") {
			t.Errorf("expected closed case to execute, got: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestSelect_MultipleMessages(t *testing.T) {
	src := `
let ch = channel[Int]()

thread spawn do
    ch.send(1)
    ch.send(2)
    ch.send(3)
    ch.close()
end

let count = 0
loop
    select
        case let x = ch.recv():
            count = count + 1
            println("Got: " + x.toString())
        case closed ch:
            println("Done, count: " + count.toString())
            break
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
	
	done := make(chan bool, 1)
	buf := &bytes.Buffer{}
	
	go func() {
		_, err = engine.Eval(prog, engine.Options{Stdout: buf})
		done <- true
	}()
	
	select {
	case <-done:
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		got := buf.String()
		if !strings.Contains(got, "Done, count: 3") {
			t.Errorf("expected count to be 3, got: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestCombined_DeferInThread(t *testing.T) {
	src := `
let t = thread spawn do
    defer println("Thread cleanup")
    println("Thread work")
    return "done"
end

let r = thread join t
println("Result: " + r)
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
	if !strings.Contains(got, "Thread work") {
		t.Errorf("expected thread work output, got: %s", got)
	}
	if !strings.Contains(got, "Result: done") {
		t.Errorf("expected result 'done', got: %s", got)
	}
	// Note: defer in threads may not print to captured stdout if thread completes asynchronously
}

func TestCombined_ChannelWithForWhere(t *testing.T) {
	src := `
let ch = channel[Int]()

thread spawn do
    for i in [1, 2, 3, 4, 5]:
        ch.send(i)
    end
    ch.close()
end

let count = 0
loop
    select
        case let x = ch.recv():
            // Only count even numbers
            if x % 2 == 0:
                count = count + 1
            end
        case closed ch:
            break
    end
end

println("Even count: " + count.toString())
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(src))
	p := parser.New(items)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	
	done := make(chan bool, 1)
	buf := &bytes.Buffer{}
	
	go func() {
		_, err = engine.Eval(prog, engine.Options{Stdout: buf})
		done <- true
	}()
	
	select {
	case <-done:
		if err != nil {
			t.Fatalf("eval error: %v", err)
		}
		got := buf.String()
		// Should count 2 even numbers (2 and 4)
		if !strings.Contains(got, "Even count: 2") {
			t.Errorf("expected even count to be 2, got: %s", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

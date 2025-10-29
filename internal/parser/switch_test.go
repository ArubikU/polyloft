package parser

import (
	"testing"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/lexer"
)

func TestParseSwitchValueMatching(t *testing.T) {
	input := `
switch x:
    case 1:
        println("one")
    case 2, 3:
        println("two or three")
    default:
        println("other")
end
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(input))
	p := New(items)
	
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	switchStmt, ok := prog.Stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Fatalf("Expected SwitchStmt, got %T", prog.Stmts[0])
	}

	if switchStmt.Expr == nil {
		t.Fatalf("Expected switch expression to be non-nil")
	}

	if len(switchStmt.Cases) != 2 {
		t.Fatalf("Expected 2 cases, got %d", len(switchStmt.Cases))
	}

	// Check first case
	if len(switchStmt.Cases[0].Values) != 1 {
		t.Fatalf("Expected 1 value in first case, got %d", len(switchStmt.Cases[0].Values))
	}

	// Check second case (multiple values)
	if len(switchStmt.Cases[1].Values) != 2 {
		t.Fatalf("Expected 2 values in second case, got %d", len(switchStmt.Cases[1].Values))
	}

	// Check default case
	if len(switchStmt.Default) == 0 {
		t.Fatalf("Expected default case to have statements")
	}
}

func TestParseSwitchTypeMatching(t *testing.T) {
	input := `
switch value:
    case (x: Int):
        println(x)
    case (s: String):
        println(s)
end
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(input))
	p := New(items)
	
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	switchStmt, ok := prog.Stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Fatalf("Expected SwitchStmt, got %T", prog.Stmts[0])
	}

	if len(switchStmt.Cases) != 2 {
		t.Fatalf("Expected 2 cases, got %d", len(switchStmt.Cases))
	}

	// Check first case - type matching
	if switchStmt.Cases[0].TypeName != "Int" {
		t.Fatalf("Expected type name 'Int', got '%s'", switchStmt.Cases[0].TypeName)
	}
	if switchStmt.Cases[0].VarName != "x" {
		t.Fatalf("Expected var name 'x', got '%s'", switchStmt.Cases[0].VarName)
	}

	// Check second case - type matching
	if switchStmt.Cases[1].TypeName != "String" {
		t.Fatalf("Expected type name 'String', got '%s'", switchStmt.Cases[1].TypeName)
	}
	if switchStmt.Cases[1].VarName != "s" {
		t.Fatalf("Expected var name 's', got '%s'", switchStmt.Cases[1].VarName)
	}
}

func TestParseSwitchWithoutExpression(t *testing.T) {
	// Switch without expression is not supported in the current implementation
	// The colon after switch is treated as unexpected
	// This test verifies that the parser correctly rejects this syntax
	input := `
switch:
    case true:
        println("condition 1")
    case false:
        println("condition 2")
end
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(input))
	p := New(items)
	
	_, err := p.Parse()
	if err == nil {
		t.Fatalf("Expected parse error for switch without expression, got nil")
	}
	// Error is expected
}

func TestParseSwitchDefaultOnly(t *testing.T) {
	input := `
switch x:
    default:
        println("default")
end
`
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(input))
	p := New(items)
	
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(prog.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(prog.Stmts))
	}

	switchStmt, ok := prog.Stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Fatalf("Expected SwitchStmt, got %T", prog.Stmts[0])
	}

	if len(switchStmt.Cases) != 0 {
		t.Fatalf("Expected 0 cases, got %d", len(switchStmt.Cases))
	}

	if len(switchStmt.Default) == 0 {
		t.Fatalf("Expected default case to have statements")
	}
}

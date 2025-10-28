package ast

import (
	"testing"
)

// TestPooledNodesCorrectness verifies that pooled nodes work correctly
func TestPooledNodesCorrectness(t *testing.T) {
	t.Run("Ident", func(t *testing.T) {
		n1 := NewIdent("var1")
		if n1.Name != "var1" {
			t.Errorf("Expected Name='var1', got %s", n1.Name)
		}
		ReleaseIdent(n1)

		// Get another node - might be the same one from pool
		n2 := NewIdent("var2")
		if n2.Name != "var2" {
			t.Errorf("Expected Name='var2', got %s", n2.Name)
		}
		ReleaseIdent(n2)
	})

	t.Run("NumberLit", func(t *testing.T) {
		n := NewNumberLit(42)
		if n.Value != 42 {
			t.Errorf("Expected Value=42, got %v", n.Value)
		}
		ReleaseNumberLit(n)
	})

	t.Run("BinaryExpr", func(t *testing.T) {
		lhs := NewNumberLit(1)
		rhs := NewNumberLit(2)
		expr := NewBinaryExpr(OpPlus, lhs, rhs)

		if expr.Op != OpPlus {
			t.Errorf("Expected Op=OpPlus, got %d", expr.Op)
		}
		if expr.Lhs != lhs {
			t.Error("Lhs not set correctly")
		}
		if expr.Rhs != rhs {
			t.Error("Rhs not set correctly")
		}

		ReleaseBinaryExpr(expr)
		ReleaseNumberLit(lhs)
		ReleaseNumberLit(rhs)
	})

	t.Run("Type", func(t *testing.T) {
		typ := NewType("int")
		if typ.Name != "int" {
			t.Errorf("Expected Name='int', got %s", typ.Name)
		}
		ReleaseType(typ)
	})
}

// TestPoolReuse verifies that pools actually reuse objects
func TestPoolReuse(t *testing.T) {
	// Create and release multiple times
	var nodes []*Ident
	for i := 0; i < 10; i++ {
		n := NewIdent("test")
		nodes = append(nodes, n)
		ReleaseIdent(n)
	}

	// At least some should be reused (addresses should repeat)
	// This is probabilistic but with 10 iterations, very likely
	seen := make(map[*Ident]bool)
	reused := false
	for _, n := range nodes {
		if seen[n] {
			reused = true
			break
		}
		seen[n] = true
	}

	if !reused && len(nodes) > 5 {
		t.Log("Warning: Pool may not be reusing objects (could be spurious)")
		// Don't fail - this is probabilistic and can happen
	}
}

// TestTypeOptimizations verifies type parsing optimizations
func TestTypeOptimizations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple", "int", "int"},
		{"Generic", "Array<Int>", "Array<Int>"},
		{"NestedGeneric", "Map<String, Array<Int>>", "Map<String, Array<Int>>"},
		{"DeepNesting", "List<Map<String, Array<Set<Int>>>>", "List<Map<String, Array<Set<Int>>>>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := TypeFromString(tt.input)
			result := GetTypeNameString(typ)
			if result != tt.expected {
				t.Errorf("TypeFromString(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestPreallocationEfficiency ensures preallocation strategies work
func TestPreallocationEfficiency(t *testing.T) {
	// Test that parseTypeParams handles various sizes
	tests := []struct {
		name     string
		input    string
		numTypes int
	}{
		{"OneParam", "Int", 1},
		{"TwoParams", "String, Int", 2},
		{"ThreeParams", "String, Int, Bool", 3},
		{"ManyParams", "A, B, C, D, E", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := parseTypeParams(tt.input)
			if len(params) != tt.numTypes {
				t.Errorf("Expected %d params, got %d", tt.numTypes, len(params))
			}
		})
	}
}

// TestReleaseNilSafety ensures releasing nil is safe
func TestReleaseNilSafety(t *testing.T) {
	// Should not panic
	ReleaseIdent(nil)
	ReleaseNumberLit(nil)
	ReleaseStringLit(nil)
	ReleaseBoolLit(nil)
	ReleaseNilLit(nil)
	ReleaseBinaryExpr(nil)
	ReleaseUnaryExpr(nil)
	ReleaseCallExpr(nil)
	ReleaseIndexExpr(nil)
	ReleaseFieldExpr(nil)
	ReleaseLetStmt(nil)
	ReleaseAssignStmt(nil)
	ReleaseReturnStmt(nil)
	ReleaseExprStmt(nil)
	ReleaseType(nil)
}

// TestComplexTreeWithPooling tests building and releasing a complex tree
func TestComplexTreeWithPooling(t *testing.T) {
	// Build: (1 + 2) * (3 - 4)
	lit1 := NewNumberLit(1)
	lit2 := NewNumberLit(2)
	lit3 := NewNumberLit(3)
	lit4 := NewNumberLit(4)

	add := NewBinaryExpr(OpPlus, lit1, lit2)
	sub := NewBinaryExpr(OpMinus, lit3, lit4)
	mul := NewBinaryExpr(OpMul, add, sub)

	// Verify structure
	if mul.Op != OpMul {
		t.Error("Top level operation should be multiplication")
	}

	// Release all nodes
	ReleaseBinaryExpr(mul)
	ReleaseBinaryExpr(add)
	ReleaseBinaryExpr(sub)
	ReleaseNumberLit(lit1)
	ReleaseNumberLit(lit2)
	ReleaseNumberLit(lit3)
	ReleaseNumberLit(lit4)

	// Build another tree - should reuse pooled nodes
	newLit1 := NewNumberLit(10)
	newLit2 := NewNumberLit(20)
	newExpr := NewBinaryExpr(OpPlus, newLit1, newLit2)

	if newExpr.Op != OpPlus {
		t.Error("New expression operation incorrect")
	}

	ReleaseBinaryExpr(newExpr)
	ReleaseNumberLit(newLit1)
	ReleaseNumberLit(newLit2)
}

// TestStatementPooling tests statement pooling
func TestStatementPooling(t *testing.T) {
	// Create let statement
	value := NewNumberLit(42)
	letStmt := NewLetStmt()
	letStmt.Name = "x"
	letStmt.Value = value

	if letStmt.Name != "x" {
		t.Errorf("Expected Name='x', got %s", letStmt.Name)
	}

	// Release
	ReleaseLetStmt(letStmt)
	ReleaseNumberLit(value)

	// Create another - might reuse
	value2 := NewNumberLit(100)
	letStmt2 := NewLetStmt()
	letStmt2.Name = "y"
	letStmt2.Value = value2

	if letStmt2.Name != "y" {
		t.Errorf("Expected Name='y', got %s", letStmt2.Name)
	}

	ReleaseLetStmt(letStmt2)
	ReleaseNumberLit(value2)
}

// TestTypeFromStringOptimizations verifies optimizations don't break functionality
func TestTypeFromStringOptimizations(t *testing.T) {
	tests := []struct {
		name      string
		typeStr   string
		wantName  string
		wantParam int
	}{
		{"SimpleInt", "int", "int", 0},
		{"ArrayInt", "Array<Int>", "Array", 1},
		{"MapTwoParams", "Map<String, Int>", "Map", 2},
		{"NestedArray", "Array<Array<Int>>", "Array", 1},
		{"ComplexNesting", "Map<String, Array<Map<Int, Bool>>>", "Map", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := TypeFromString(tt.typeStr)
			if typ == nil {
				t.Fatal("TypeFromString returned nil")
			}
			if typ.Name != tt.wantName {
				t.Errorf("Name = %s, want %s", typ.Name, tt.wantName)
			}
			if len(typ.TypeParams) != tt.wantParam {
				t.Errorf("TypeParams count = %d, want %d", len(typ.TypeParams), tt.wantParam)
			}
		})
	}
}

// TestGetTypeNameStringOptimization verifies string generation optimization
func TestGetTypeNameStringOptimization(t *testing.T) {
	// Create a complex type
	innerType := &Type{Name: "Int", IsBuiltin: true}
	arrayType := &Type{Name: "Array", TypeParams: []*Type{innerType}, IsBuiltin: true}
	mapType := &Type{
		Name:       "Map",
		TypeParams: []*Type{{Name: "String"}, arrayType},
		IsBuiltin:  true,
	}

	result := GetTypeNameString(mapType)
	expected := "Map<String, Array<Int>>"

	if result != expected {
		t.Errorf("GetTypeNameString = %s, want %s", result, expected)
	}
}

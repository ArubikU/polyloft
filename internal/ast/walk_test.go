package ast

import (
	"testing"
)

func TestIterativeWalk(t *testing.T) {
	// Build a test tree: let x = 1 + 2
	tree := &Program{
		Stmts: []Stmt{
			&LetStmt{
				Name: "x",
				Value: &BinaryExpr{
					Op:  OpPlus,
					Lhs: &NumberLit{Value: 1},
					Rhs: &NumberLit{Value: 2},
				},
			},
		},
	}

	t.Run("VisitsAllNodes", func(t *testing.T) {
		visited := 0
		IterativeWalk(tree, func(n Node) bool {
			visited++
			return true
		})

		// Should visit: Program, LetStmt, BinaryExpr, 2 NumberLits = 5 nodes
		expected := 5
		if visited != expected {
			t.Errorf("Expected %d nodes visited, got %d", expected, visited)
		}
	})

	t.Run("CanStopEarly", func(t *testing.T) {
		visited := 0
		IterativeWalk(tree, func(n Node) bool {
			visited++
			return visited < 3 // Stop after 3 nodes
		})

		if visited != 3 {
			t.Errorf("Expected to visit exactly 3 nodes, got %d", visited)
		}
	})

	t.Run("HandlesNilRoot", func(t *testing.T) {
		visited := 0
		IterativeWalk(nil, func(n Node) bool {
			visited++
			return true
		})

		if visited != 0 {
			t.Errorf("Expected 0 visits for nil root, got %d", visited)
		}
	})
}

func TestCountNodes(t *testing.T) {
	tests := []struct {
		name     string
		tree     Node
		expected int
	}{
		{
			name:     "SingleNode",
			tree:     &NumberLit{Value: 42},
			expected: 1,
		},
		{
			name: "BinaryExpr",
			tree: &BinaryExpr{
				Op:  OpPlus,
				Lhs: &NumberLit{Value: 1},
				Rhs: &NumberLit{Value: 2},
			},
			expected: 3, // BinaryExpr + 2 NumberLits
		},
		{
			name: "Program",
			tree: &Program{
				Stmts: []Stmt{
					&ExprStmt{X: &NumberLit{Value: 1}},
					&ExprStmt{X: &NumberLit{Value: 2}},
				},
			},
			expected: 5, // Program + 2 ExprStmts + 2 NumberLits
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CountNodes(tt.tree)
			if count != tt.expected {
				t.Errorf("Expected %d nodes, got %d", tt.expected, count)
			}
		})
	}
}

func TestFindNodes(t *testing.T) {
	tree := &Program{
		Stmts: []Stmt{
			&LetStmt{Name: "x", Value: &NumberLit{Value: 1}},
			&LetStmt{Name: "y", Value: &NumberLit{Value: 2}},
			&ExprStmt{X: &BinaryExpr{
				Op:  OpPlus,
				Lhs: &Ident{Name: "x"},
				Rhs: &Ident{Name: "y"},
			}},
		},
	}

	t.Run("FindAllNumberLits", func(t *testing.T) {
		numbers := FindNodes(tree, func(n Node) bool {
			_, ok := n.(*NumberLit)
			return ok
		})

		if len(numbers) != 2 {
			t.Errorf("Expected 2 NumberLits, got %d", len(numbers))
		}
	})

	t.Run("FindAllIdents", func(t *testing.T) {
		idents := FindNodes(tree, func(n Node) bool {
			_, ok := n.(*Ident)
			return ok
		})

		if len(idents) != 2 {
			t.Errorf("Expected 2 Idents, got %d", len(idents))
		}
	})

	t.Run("FindAllLetStmts", func(t *testing.T) {
		lets := FindNodes(tree, func(n Node) bool {
			_, ok := n.(*LetStmt)
			return ok
		})

		if len(lets) != 2 {
			t.Errorf("Expected 2 LetStmts, got %d", len(lets))
		}
	})
}

func TestFindFirstNode(t *testing.T) {
	tree := &Program{
		Stmts: []Stmt{
			&LetStmt{Name: "x", Value: &NumberLit{Value: 1}},
			&LetStmt{Name: "y", Value: &NumberLit{Value: 2}},
			&ExprStmt{X: &NumberLit{Value: 3}},
		},
	}

	t.Run("FindsFirstMatch", func(t *testing.T) {
		first := FindFirstNode(tree, func(n Node) bool {
			num, ok := n.(*NumberLit)
			return ok && num.Value == 1
		})

		if first == nil {
			t.Fatal("Expected to find a node, got nil")
		}

		num, ok := first.(*NumberLit)
		if !ok {
			t.Fatal("Expected NumberLit")
		}
		if num.Value != 1 {
			t.Errorf("Expected value 1, got %v", num.Value)
		}
	})

	t.Run("ReturnsNilWhenNotFound", func(t *testing.T) {
		first := FindFirstNode(tree, func(n Node) bool {
			_, ok := n.(*ChannelExpr)
			return ok
		})

		if first != nil {
			t.Errorf("Expected nil, got %v", first)
		}
	})
}

func TestIterativeWalkComplexTree(t *testing.T) {
	// Build a more complex tree with nested structures
	tree := &Program{
		Stmts: []Stmt{
			&IfStmt{
				Clauses: []IfClause{
					{
						Cond: &BoolLit{Value: true},
						Body: []Stmt{
							&ExprStmt{X: &NumberLit{Value: 1}},
							&ExprStmt{X: &NumberLit{Value: 2}},
						},
					},
				},
				Else: []Stmt{
					&ExprStmt{X: &NumberLit{Value: 3}},
				},
			},
			&ForInStmt{
				Name:     "item",
				Iterable: &ArrayLit{Elems: []Expr{&NumberLit{Value: 1}, &NumberLit{Value: 2}}},
				Body: []Stmt{
					&ExprStmt{X: &Ident{Name: "item"}},
				},
			},
		},
	}

	count := CountNodes(tree)
	// Let's count actual nodes (types that implement Node interface):
	// Program (1) + IfStmt (1) + BoolLit (1) + 2*ExprStmt (2) + 2*NumberLit (2) +
	// ExprStmt (1) + NumberLit (1) + ForInStmt (1) + ArrayLit (1) + 2*NumberLit (2) +
	// ExprStmt (1) + Ident (1)
	// = 1 + 1 + 1 + 2 + 2 + 1 + 1 + 1 + 1 + 2 + 1 + 1 = 15 nodes
	// (IfClause is NOT a Node, so it's not counted)
	expected := 15
	if count != expected {
		t.Errorf("Expected %d nodes, got %d", expected, count)
	}
}

func TestIterativeWalkWithClasses(t *testing.T) {
	tree := &Program{
		Stmts: []Stmt{
			&ClassDecl{
				Name: "MyClass",
				Fields: []FieldDecl{
					{Name: "x", InitValue: &NumberLit{Value: 1}},
				},
				Methods: []MethodDecl{
					{
						Name: "method1",
						Body: []Stmt{
							&ReturnStmt{Value: &NumberLit{Value: 42}},
						},
					},
				},
			},
		},
	}

	// Count FieldDecls
	fields := FindNodes(tree, func(n Node) bool {
		_, ok := n.(*FieldDecl)
		return ok
	})

	if len(fields) != 1 {
		t.Errorf("Expected 1 FieldDecl, got %d", len(fields))
	}

	// Count all nodes
	// Note: MethodDecl is not a Node, so only its body is visited
	total := CountNodes(tree)
	// Program + ClassDecl + FieldDecl + NumberLit + ReturnStmt + NumberLit
	// = 6 nodes (MethodDecl is not a Node, so it's not counted)
	expected := 6
	if total != expected {
		t.Errorf("Expected %d nodes, got %d", expected, total)
	}
}

// Benchmark iterative vs recursive traversal
func BenchmarkIterativeWalk(b *testing.B) {
	// Build a moderately deep tree
	tree := &Program{
		Stmts: []Stmt{
			&LetStmt{Name: "x", Value: &NumberLit{Value: 1}},
			&LetStmt{Name: "y", Value: &NumberLit{Value: 2}},
			&ExprStmt{X: &BinaryExpr{
				Op: OpPlus,
				Lhs: &BinaryExpr{
					Op:  OpMul,
					Lhs: &NumberLit{Value: 3},
					Rhs: &NumberLit{Value: 4},
				},
				Rhs: &BinaryExpr{
					Op:  OpMinus,
					Lhs: &NumberLit{Value: 5},
					Rhs: &NumberLit{Value: 6},
				},
			}},
		},
	}

	b.Run("IterativeTraversal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			count := 0
			IterativeWalk(tree, func(n Node) bool {
				count++
				return true
			})
		}
	})

	b.Run("CountNodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = CountNodes(tree)
		}
	})

	b.Run("FindNodes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindNodes(tree, func(n Node) bool {
				_, ok := n.(*NumberLit)
				return ok
			})
		}
	})
}

// Test that iterative walk doesn't overflow stack on deep trees
func TestIterativeWalkDeepTree(b *testing.T) {
	// Build a very deep tree that would overflow with recursion
	var expr Expr = &NumberLit{Value: 0}
	for i := 0; i < 10000; i++ {
		expr = &BinaryExpr{
			Op:  OpPlus,
			Lhs: expr,
			Rhs: &NumberLit{Value: i},
		}
	}

	tree := &Program{
		Stmts: []Stmt{
			&ExprStmt{X: expr},
		},
	}

	// This should not stack overflow
	count := CountNodes(tree)
	
	// Should count: Program + ExprStmt + 10000 BinaryExprs + 10001 NumberLits = 20003
	expected := 20003
	if count != expected {
		b.Errorf("Expected %d nodes, got %d", expected, count)
	}
}

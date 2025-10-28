package ast

import (
	"testing"
)

// Benchmark AST node creation
func BenchmarkNodeCreation(b *testing.B) {
	b.Run("CreateIdent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &Ident{Name: "testVar"}
		}
	})

	b.Run("CreateNumberLit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &NumberLit{Value: 42}
		}
	})

	b.Run("CreateStringLit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &StringLit{Value: "test"}
		}
	})

	b.Run("CreateBinaryExpr", func(b *testing.B) {
		lhs := &NumberLit{Value: 1}
		rhs := &NumberLit{Value: 2}
		for i := 0; i < b.N; i++ {
			_ = &BinaryExpr{Op: OpPlus, Lhs: lhs, Rhs: rhs}
		}
	})

	b.Run("CreateCallExpr", func(b *testing.B) {
		callee := &Ident{Name: "func"}
		args := []Expr{&NumberLit{Value: 1}, &NumberLit{Value: 2}}
		for i := 0; i < b.N; i++ {
			_ = &CallExpr{Callee: callee, Args: args}
		}
	})
}

// Benchmark AST tree building
func BenchmarkTreeBuilding(b *testing.B) {
	b.Run("SmallTree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Build a simple expression tree: (1 + 2) * 3
			expr1 := &BinaryExpr{
				Op:  OpPlus,
				Lhs: &NumberLit{Value: 1},
				Rhs: &NumberLit{Value: 2},
			}
			_ = &BinaryExpr{
				Op:  OpMul,
				Lhs: expr1,
				Rhs: &NumberLit{Value: 3},
			}
		}
	})

	b.Run("MediumTree", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Build a more complex tree: ((1 + 2) * (3 - 4)) / 5
			expr1 := &BinaryExpr{
				Op:  OpPlus,
				Lhs: &NumberLit{Value: 1},
				Rhs: &NumberLit{Value: 2},
			}
			expr2 := &BinaryExpr{
				Op:  OpMinus,
				Lhs: &NumberLit{Value: 3},
				Rhs: &NumberLit{Value: 4},
			}
			expr3 := &BinaryExpr{
				Op:  OpMul,
				Lhs: expr1,
				Rhs: expr2,
			}
			_ = &BinaryExpr{
				Op:  OpDiv,
				Lhs: expr3,
				Rhs: &NumberLit{Value: 5},
			}
		}
	})

	b.Run("ArrayLiteral", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			elems := make([]Expr, 10)
			for j := 0; j < 10; j++ {
				elems[j] = &NumberLit{Value: j}
			}
			_ = &ArrayLit{Elems: elems}
		}
	})
}

// Benchmark type operations
func BenchmarkTypeOperations(b *testing.B) {
	b.Run("CreateSimpleType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &Type{Name: "int", IsBuiltin: true}
		}
	})

	b.Run("CreateGenericType", func(b *testing.B) {
		elemType := &Type{Name: "Int", IsBuiltin: true}
		for i := 0; i < b.N; i++ {
			_ = &Type{
				Name:       "Array",
				TypeParams: []*Type{elemType},
				IsBuiltin:  true,
			}
		}
	})

	b.Run("TypeFromString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = TypeFromString("int")
		}
	})

	b.Run("TypeFromStringGeneric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = TypeFromString("Array<Int>")
		}
	})

	b.Run("TypeFromStringNestedGeneric", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = TypeFromString("Map<String, Array<Int>>")
		}
	})

	b.Run("MatchesType", func(b *testing.B) {
		typ := &Type{Name: "bool", Aliases: []string{"boolean", "Bool"}, IsBuiltin: true}
		for i := 0; i < b.N; i++ {
			_ = typ.MatchesType("boolean")
		}
	})

	b.Run("GetTypeNameString", func(b *testing.B) {
		typ := TypeFromString("Map<String, Array<Int>>")
		for i := 0; i < b.N; i++ {
			_ = GetTypeNameString(typ)
		}
	})
}

// Benchmark AST traversal patterns
func BenchmarkTraversal(b *testing.B) {
	// Build a test tree once
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

	b.Run("TypeSwitch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, stmt := range tree.Stmts {
				switch s := stmt.(type) {
				case *LetStmt:
					_ = s.Name
				case *ExprStmt:
					_ = s.X
				}
			}
		}
	})

	b.Run("InterfaceCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, stmt := range tree.Stmts {
				if letStmt, ok := stmt.(*LetStmt); ok {
					_ = letStmt.Name
				} else if exprStmt, ok := stmt.(*ExprStmt); ok {
					_ = exprStmt.X
				}
			}
		}
	})
}

// Benchmark statement creation
func BenchmarkStmtCreation(b *testing.B) {
	b.Run("LetStmt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &LetStmt{
				Name:  "x",
				Value: &NumberLit{Value: 42},
			}
		}
	})

	b.Run("IfStmt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &IfStmt{
				Clauses: []IfClause{
					{
						Cond: &BoolLit{Value: true},
						Body: []Stmt{
							&ExprStmt{X: &NumberLit{Value: 1}},
						},
					},
				},
			}
		}
	})

	b.Run("ForInStmt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &ForInStmt{
				Name:     "item",
				Iterable: &Ident{Name: "items"},
				Body: []Stmt{
					&ExprStmt{X: &Ident{Name: "item"}},
				},
			}
		}
	})
}

// Benchmark class/interface declarations
func BenchmarkDeclCreation(b *testing.B) {
	b.Run("ClassDecl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &ClassDecl{
				Name: "MyClass",
				Fields: []FieldDecl{
					{Name: "field1", Type: &Type{Name: "int"}},
					{Name: "field2", Type: &Type{Name: "string"}},
				},
				Methods: []MethodDecl{
					{
						Name: "method1",
						Params: []Parameter{
							{Name: "param1", Type: &Type{Name: "int"}},
						},
						ReturnType: &Type{Name: "void"},
					},
				},
			}
		}
	})

	b.Run("InterfaceDecl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &InterfaceDecl{
				Name: "MyInterface",
				Methods: []MethodSignature{
					{
						Name: "method1",
						Params: []Parameter{
							{Name: "param1", Type: &Type{Name: "int"}},
						},
						ReturnType: &Type{Name: "string"},
					},
				},
			}
		}
	})
}

// Benchmark pooled node creation
func BenchmarkPooledNodeCreation(b *testing.B) {
	b.Run("PooledIdent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := NewIdent("testVar")
			ReleaseIdent(n)
		}
	})

	b.Run("PooledNumberLit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := NewNumberLit(42)
			ReleaseNumberLit(n)
		}
	})

	b.Run("PooledStringLit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := NewStringLit("test")
			ReleaseStringLit(n)
		}
	})

	b.Run("PooledBinaryExpr", func(b *testing.B) {
		lhs := NewNumberLit(1)
		rhs := NewNumberLit(2)
		for i := 0; i < b.N; i++ {
			n := NewBinaryExpr(OpPlus, lhs, rhs)
			ReleaseBinaryExpr(n)
		}
		ReleaseNumberLit(lhs)
		ReleaseNumberLit(rhs)
	})

	b.Run("PooledType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t := NewType("int")
			ReleaseType(t)
		}
	})
}

// Compare pooled vs non-pooled creation
func BenchmarkPooledVsNonPooled(b *testing.B) {
	b.Run("NonPooled_Ident", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &Ident{Name: "testVar"}
		}
	})

	b.Run("Pooled_Ident", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			n := NewIdent("testVar")
			ReleaseIdent(n)
		}
	})

	b.Run("NonPooled_BinaryExpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &BinaryExpr{
				Op:  OpPlus,
				Lhs: &NumberLit{Value: 1},
				Rhs: &NumberLit{Value: 2},
			}
		}
	})

	b.Run("Pooled_BinaryExpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			lhs := NewNumberLit(1)
			rhs := NewNumberLit(2)
			n := NewBinaryExpr(OpPlus, lhs, rhs)
			ReleaseBinaryExpr(n)
			ReleaseNumberLit(lhs)
			ReleaseNumberLit(rhs)
		}
	})
}

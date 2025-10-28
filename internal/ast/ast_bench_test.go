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

// Benchmark composite operations that combine multiple functions
func BenchmarkCompositeOperations(b *testing.B) {
	b.Run("TypeParsing_And_NameGeneration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typ := TypeFromString("Map<String, Array<Int>>")
			_ = GetTypeNameString(typ)
		}
	})

	b.Run("MultipleTypeOperations", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t1 := TypeFromString("Int")
			t2 := TypeFromString("String")
			t3 := TypeFromString("Array<Int>")
			_ = t1.MatchesType("Int")
			_ = t2.MatchesType("String")
			_ = GetTypeNameString(t3)
		}
	})

	b.Run("GenericTypeCreation", func(b *testing.B) {
		baseType := TypeFromString("Array")
		paramType := TypeFromString("Int")
		for i := 0; i < b.N; i++ {
			_ = GenericType(baseType, paramType)
		}
	})
}

// Benchmark real-world patterns similar to Python execution
func BenchmarkRealWorldPatterns(b *testing.B) {
	b.Run("ForLoopWithAssignment", func(b *testing.B) {
		// Simulates: for a in range: b = a * a
		for i := 0; i < b.N; i++ {
			prog := &Program{
				Stmts: []Stmt{
					&ForInStmt{
						Name: "a",
						Iterable: &RangeExpr{
							Start:     &NumberLit{Value: 1},
							End:       &NumberLit{Value: 100},
							Inclusive: true,
						},
						Body: []Stmt{
							&AssignStmt{
								Target: &Ident{Name: "b"},
								Value: &BinaryExpr{
									Op:  OpMul,
									Lhs: &Ident{Name: "a"},
									Rhs: &Ident{Name: "a"},
								},
							},
						},
					},
				},
			}
			_ = prog
		}
	})

	b.Run("FieldAccessPattern", func(b *testing.B) {
		// Simulates: obj.field
		for i := 0; i < b.N; i++ {
			expr := &FieldExpr{
				X:    &Ident{Name: "obj"},
				Name: "field",
			}
			_ = expr
		}
	})

	b.Run("FunctionCallPattern", func(b *testing.B) {
		// Simulates: func(arg1, arg2)
		for i := 0; i < b.N; i++ {
			call := &CallExpr{
				Callee: &Ident{Name: "func"},
				Args: []Expr{
					&Ident{Name: "arg1"},
					&Ident{Name: "arg2"},
				},
			}
			_ = call
		}
	})

	b.Run("ComplexExpression", func(b *testing.B) {
		// Simulates: (a + b) * (c - d)
		for i := 0; i < b.N; i++ {
			expr := &BinaryExpr{
				Op: OpMul,
				Lhs: &BinaryExpr{
					Op:  OpPlus,
					Lhs: &Ident{Name: "a"},
					Rhs: &Ident{Name: "b"},
				},
				Rhs: &BinaryExpr{
					Op:  OpMinus,
					Lhs: &Ident{Name: "c"},
					Rhs: &Ident{Name: "d"},
				},
			}
			_ = expr
		}
	})

	b.Run("IfStatementPattern", func(b *testing.B) {
		// Simulates: if condition: body
		for i := 0; i < b.N; i++ {
			stmt := &IfStmt{
				Clauses: []IfClause{
					{
						Cond: &BinaryExpr{
							Op:  OpGt,
							Lhs: &Ident{Name: "x"},
							Rhs: &NumberLit{Value: 0},
						},
						Body: []Stmt{
							&ExprStmt{X: &Ident{Name: "x"}},
						},
					},
				},
			}
			_ = stmt
		}
	})
}

// Benchmark constant usage
func BenchmarkConstantUsage(b *testing.B) {
b.Run("CommonNumberLit_Allocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = &NumberLit{Value: 1}
}
})

b.Run("CommonNumberLit_PreAllocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = GetCommonNumberLit(1)
}
})

b.Run("BoolLit_Allocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = &BoolLit{Value: true}
}
})

b.Run("BoolLit_PreAllocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = GetCommonBoolLit(true)
}
})

b.Run("NilLit_Allocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = &NilLit{}
}
})

b.Run("NilLit_PreAllocated", func(b *testing.B) {
for i := 0; i < b.N; i++ {
_ = NilValue
}
})
}

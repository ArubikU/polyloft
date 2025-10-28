package ast

import "sync"

// Node pools for memory reuse
// These pools reduce GC pressure by reusing frequently allocated AST nodes

var (
	identPool = sync.Pool{
		New: func() interface{} {
			return &Ident{}
		},
	}

	numberLitPool = sync.Pool{
		New: func() interface{} {
			return &NumberLit{}
		},
	}

	stringLitPool = sync.Pool{
		New: func() interface{} {
			return &StringLit{}
		},
	}

	boolLitPool = sync.Pool{
		New: func() interface{} {
			return &BoolLit{}
		},
	}

	nilLitPool = sync.Pool{
		New: func() interface{} {
			return &NilLit{}
		},
	}

	binaryExprPool = sync.Pool{
		New: func() interface{} {
			return &BinaryExpr{}
		},
	}

	unaryExprPool = sync.Pool{
		New: func() interface{} {
			return &UnaryExpr{}
		},
	}

	callExprPool = sync.Pool{
		New: func() interface{} {
			return &CallExpr{}
		},
	}

	indexExprPool = sync.Pool{
		New: func() interface{} {
			return &IndexExpr{}
		},
	}

	fieldExprPool = sync.Pool{
		New: func() interface{} {
			return &FieldExpr{}
		},
	}

	letStmtPool = sync.Pool{
		New: func() interface{} {
			return &LetStmt{}
		},
	}

	assignStmtPool = sync.Pool{
		New: func() interface{} {
			return &AssignStmt{}
		},
	}

	returnStmtPool = sync.Pool{
		New: func() interface{} {
			return &ReturnStmt{}
		},
	}

	exprStmtPool = sync.Pool{
		New: func() interface{} {
			return &ExprStmt{}
		},
	}

	typePool = sync.Pool{
		New: func() interface{} {
			return &Type{}
		},
	}
)

// NewIdent gets an Ident from the pool
func NewIdent(name string) *Ident {
	n := identPool.Get().(*Ident)
	n.Name = name
	return n
}

// ReleaseIdent returns an Ident to the pool
func ReleaseIdent(n *Ident) {
	if n == nil {
		return
	}
	n.Name = ""
	identPool.Put(n)
}

// NewNumberLit gets a NumberLit from the pool
func NewNumberLit(value any) *NumberLit {
	n := numberLitPool.Get().(*NumberLit)
	n.Value = value
	return n
}

// ReleaseNumberLit returns a NumberLit to the pool
func ReleaseNumberLit(n *NumberLit) {
	if n == nil {
		return
	}
	n.Value = nil
	numberLitPool.Put(n)
}

// NewStringLit gets a StringLit from the pool
func NewStringLit(value string) *StringLit {
	n := stringLitPool.Get().(*StringLit)
	n.Value = value
	return n
}

// ReleaseStringLit returns a StringLit to the pool
func ReleaseStringLit(n *StringLit) {
	if n == nil {
		return
	}
	n.Value = ""
	stringLitPool.Put(n)
}

// NewBoolLit gets a BoolLit from the pool
func NewBoolLit(value bool) *BoolLit {
	n := boolLitPool.Get().(*BoolLit)
	n.Value = value
	return n
}

// ReleaseBoolLit returns a BoolLit to the pool
func ReleaseBoolLit(n *BoolLit) {
	if n == nil {
		return
	}
	n.Value = false
	boolLitPool.Put(n)
}

// NewNilLit gets a NilLit from the pool
func NewNilLit() *NilLit {
	return nilLitPool.Get().(*NilLit)
}

// ReleaseNilLit returns a NilLit to the pool
func ReleaseNilLit(n *NilLit) {
	if n == nil {
		return
	}
	nilLitPool.Put(n)
}

// NewBinaryExpr gets a BinaryExpr from the pool
func NewBinaryExpr(op int, lhs, rhs Expr) *BinaryExpr {
	n := binaryExprPool.Get().(*BinaryExpr)
	n.Op = op
	n.Lhs = lhs
	n.Rhs = rhs
	return n
}

// ReleaseBinaryExpr returns a BinaryExpr to the pool
func ReleaseBinaryExpr(n *BinaryExpr) {
	if n == nil {
		return
	}
	n.Op = 0
	n.Lhs = nil
	n.Rhs = nil
	binaryExprPool.Put(n)
}

// NewUnaryExpr gets a UnaryExpr from the pool
func NewUnaryExpr(op int, x Expr) *UnaryExpr {
	n := unaryExprPool.Get().(*UnaryExpr)
	n.Op = op
	n.X = x
	return n
}

// ReleaseUnaryExpr returns a UnaryExpr to the pool
func ReleaseUnaryExpr(n *UnaryExpr) {
	if n == nil {
		return
	}
	n.Op = 0
	n.X = nil
	unaryExprPool.Put(n)
}

// NewCallExpr gets a CallExpr from the pool
func NewCallExpr(callee Expr, args []Expr) *CallExpr {
	n := callExprPool.Get().(*CallExpr)
	n.Callee = callee
	n.Args = args
	return n
}

// ReleaseCallExpr returns a CallExpr to the pool
func ReleaseCallExpr(n *CallExpr) {
	if n == nil {
		return
	}
	n.Callee = nil
	n.Args = nil
	callExprPool.Put(n)
}

// NewIndexExpr gets an IndexExpr from the pool
func NewIndexExpr(x, index Expr) *IndexExpr {
	n := indexExprPool.Get().(*IndexExpr)
	n.X = x
	n.Index = index
	return n
}

// ReleaseIndexExpr returns an IndexExpr to the pool
func ReleaseIndexExpr(n *IndexExpr) {
	if n == nil {
		return
	}
	n.X = nil
	n.Index = nil
	indexExprPool.Put(n)
}

// NewFieldExpr gets a FieldExpr from the pool
func NewFieldExpr(x Expr, name string) *FieldExpr {
	n := fieldExprPool.Get().(*FieldExpr)
	n.X = x
	n.Name = name
	return n
}

// ReleaseFieldExpr returns a FieldExpr to the pool
func ReleaseFieldExpr(n *FieldExpr) {
	if n == nil {
		return
	}
	n.X = nil
	n.Name = ""
	fieldExprPool.Put(n)
}

// NewLetStmt gets a LetStmt from the pool
func NewLetStmt() *LetStmt {
	return letStmtPool.Get().(*LetStmt)
}

// ReleaseLetStmt returns a LetStmt to the pool
func ReleaseLetStmt(n *LetStmt) {
	if n == nil {
		return
	}
	n.Name = ""
	n.Names = nil
	n.Value = nil
	n.Type = nil
	n.Modifiers = nil
	n.Kind = ""
	n.Inferred = false
	letStmtPool.Put(n)
}

// NewAssignStmt gets an AssignStmt from the pool
func NewAssignStmt(target, value Expr) *AssignStmt {
	n := assignStmtPool.Get().(*AssignStmt)
	n.Target = target
	n.Value = value
	return n
}

// ReleaseAssignStmt returns an AssignStmt to the pool
func ReleaseAssignStmt(n *AssignStmt) {
	if n == nil {
		return
	}
	n.Target = nil
	n.Value = nil
	n.Pos = Position{}
	assignStmtPool.Put(n)
}

// NewReturnStmt gets a ReturnStmt from the pool
func NewReturnStmt(value Expr) *ReturnStmt {
	n := returnStmtPool.Get().(*ReturnStmt)
	n.Value = value
	return n
}

// ReleaseReturnStmt returns a ReturnStmt to the pool
func ReleaseReturnStmt(n *ReturnStmt) {
	if n == nil {
		return
	}
	n.Value = nil
	returnStmtPool.Put(n)
}

// NewExprStmt gets an ExprStmt from the pool
func NewExprStmt(x Expr) *ExprStmt {
	n := exprStmtPool.Get().(*ExprStmt)
	n.X = x
	return n
}

// ReleaseExprStmt returns an ExprStmt to the pool
func ReleaseExprStmt(n *ExprStmt) {
	if n == nil {
		return
	}
	n.X = nil
	exprStmtPool.Put(n)
}

// NewType gets a Type from the pool
func NewType(name string) *Type {
	t := typePool.Get().(*Type)
	t.Name = name
	return t
}

// ReleaseType returns a Type to the pool
func ReleaseType(t *Type) {
	if t == nil {
		return
	}
	t.Name = ""
	t.Aliases = nil
	t.TypeParams = nil
	t.UnionTypes = nil
	t.GoParallel = false
	t.IsBuiltin = false
	t.IsClass = false
	t.IsInterface = false
	t.IsEnum = false
	t.IsRecord = false
	t.IsUnion = false
	typePool.Put(t)
}

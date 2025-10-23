package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/typecheck"
)

// evalInstanceOfExpr handles instanceof expressions
func evalInstanceOfExpr(env *Env, expr *ast.InstanceOfExpr) (any, error) {
	// Evaluate the expression to check
	obj, err := evalExpr(env, expr.Expr)
	if err != nil {
		return nil, err
	}

	// Check if obj is instance of TypeName
	result := typecheck.IsInstanceOf(obj, expr.TypeName)

	// If variable assignment is specified, assign the object if instanceof is true
	if expr.Variable != "" && result {
		env.Define(expr.Variable, obj, expr.Modifier)
	}

	return result, nil
}

// evalTypeExpr handles Sys.type(obj) expressions
func evalTypeExpr(env *Env, expr *ast.TypeExpr) (any, error) {
	// Evaluate the expression to get type of
	obj, err := evalExpr(env, expr.Expr)
	if err != nil {
		return nil, err
	}

	// Return the type name
	return common.GetTypeName(obj), nil
}

// Enhanced evalExpr with support for new expression types
func evalExprEnhanced(env *Env, ex ast.Expr) (any, error) {
	switch e := ex.(type) {
	case *ast.InstanceOfExpr:
		return evalInstanceOfExpr(env, e)
	case *ast.TypeExpr:
		return evalTypeExpr(env, e)
	default:
		// Fall back to the original evalExpr for other expression types
		return evalExpr(env, ex)
	}
}

package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InitFunctionInterfaces initializes the Function and BiFunction interfaces
func InitFunctionInterfaces(env *Env) error {
	// Create Function<P, R> interface
	functionInterface := &common.InterfaceDefinition{
		Name:         "Function",
		Type:         &ast.Type{Name: "Function", IsBuiltin: true, IsInterface: true},
		Methods:      make(map[string][]common.MethodSignature),
		StaticFields: make(map[string]any),
		AccessLevel:  "public",
		FileName:     "builtin",
		PackageName:  "polyloft.lang",
	}

	// Add apply method: apply(param: P) -> R
	functionInterface.Methods["apply"] = []common.MethodSignature{{
		Name: "apply",
		Params: []ast.Parameter{
			{Name: "param", Type: nil}, // Generic parameter P
		},
		ReturnType: nil, // Generic return type R
		HasDefault: false,
	}}

	// Register the Function interface
	interfaceRegistry["Function"] = functionInterface
	env.Set("Function", functionInterface)

	// Create BiFunction<P1, P2, R> interface
	biFunctionInterface := &common.InterfaceDefinition{
		Name:         "BiFunction",
		Type:         &ast.Type{Name: "BiFunction", IsBuiltin: true, IsInterface: true},
		Methods:      make(map[string][]common.MethodSignature),
		StaticFields: make(map[string]any),
		AccessLevel:  "public",
		FileName:     "builtin",
		PackageName:  "polyloft.lang",
	}

	// Add apply method: apply(param1: P1, param2: P2) -> R
	biFunctionInterface.Methods["apply"] = []common.MethodSignature{{
		Name: "apply",
		Params: []ast.Parameter{
			{Name: "param1", Type: nil}, // Generic parameter P1
			{Name: "param2", Type: nil}, // Generic parameter P2
		},
		ReturnType: nil, // Generic return type R
		HasDefault: false,
	}}

	// Register the BiFunction interface
	interfaceRegistry["BiFunction"] = biFunctionInterface
	env.Set("BiFunction", biFunctionInterface)

	return nil
}

// WrapLambdaAsFunction wraps a lambda expression as a Function interface implementation
func WrapLambdaAsFunction(lambda *ast.LambdaExpr, env *Env) (common.Func, error) {
	return func(callEnv *Env, args []any) (any, error) {
		// Create a new environment for the lambda
		lambdaEnv := (*Env)(common.NewEnv())

		// Bind parameters
		if len(args) != len(lambda.Params) {
			return nil, ThrowArityError(env, len(lambda.Params), len(args))
		}

		for i, param := range lambda.Params {
			lambdaEnv.Set(param.Name, args[i])
		}

		// Execute lambda body
		if lambda.IsBlock {
			// Block lambda with statements
			var result any
			for _, stmt := range lambda.BlockBody {
				val, hasReturn, err := evalStmt(lambdaEnv, stmt)
				if err != nil {
					return nil, err
				}
				if hasReturn {
					return val, nil
				}
				result = val
			}
			return result, nil
		} else {
			// Expression lambda
			return evalExpr(lambdaEnv, lambda.Body)
		}
	}, nil
}

// CreateFunctionInstance creates an instance of Function interface from a lambda
func CreateFunctionInstance(lambda *ast.LambdaExpr, env *Env) (map[string]any, error) {
	funcImpl, err := WrapLambdaAsFunction(lambda, env)
	if err != nil {
		return nil, err
	}

	instance := map[string]any{
		"__interface": "Function",
		"apply":       funcImpl,
	}

	return instance, nil
}

// CreateBiFunctionInstance creates an instance of BiFunction interface from a lambda
func CreateBiFunctionInstance(lambda *ast.LambdaExpr, env *Env) (map[string]any, error) {
	if len(lambda.Params) != 2 {
		return nil, ThrowArityError(env, 2, len(lambda.Params))
	}

	funcImpl, err := WrapLambdaAsFunction(lambda, env)
	if err != nil {
		return nil, err
	}

	instance := map[string]any{
		"__interface": "BiFunction",
		"apply":       funcImpl,
	}

	return instance, nil
}

// IsFunction checks if a value implements the Function interface
func IsFunction(val any) bool {
	if m, ok := val.(map[string]any); ok {
		if iface, exists := m["__interface"]; exists {
			return iface == "Function"
		}
	}
	return false
}

// IsBiFunction checks if a value implements the BiFunction interface
func IsBiFunction(val any) bool {
	if m, ok := val.(map[string]any); ok {
		if iface, exists := m["__interface"]; exists {
			return iface == "BiFunction"
		}
	}
	return false
}

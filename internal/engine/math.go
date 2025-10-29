package engine

import (
	"math"
	"math/rand"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallMathModule installs the Math module with all mathematical functions
func InstallMathModule(env *Env) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	mathClass := NewClassBuilder("Math").
		AddStaticField("PI", math.Pi).
		AddStaticField("E", math.E).
		AddStaticMethod("abs", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Abs(a), nil
		})).
		AddStaticMethod("floor", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Floor(a), nil
		})).
		AddStaticMethod("ceil", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Ceil(a), nil
		})).
		AddStaticMethod("round", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Round(a), nil
		})).
		AddStaticMethod("sqrt", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Sqrt(a), nil
		})).
		AddStaticMethod("sin", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Sin(a), nil
		})).
		AddStaticMethod("cos", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Cos(a), nil
		})).
		AddStaticMethod("tan", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 1 {
				return nil, ThrowArityError(env, 1, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			return math.Tan(a), nil
		})).
		AddStaticMethod("pow", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "base", Type: ast.TypeFromString("Number")},
			{Name: "exponent", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 2 {
				return nil, ThrowArityError(env, 2, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			b, _ := utils.AsFloat(args[1])
			return math.Pow(a, b), nil
		})).
		AddStaticMethod("min", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "a", Type: ast.TypeFromString("Number")},
			{Name: "b", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 2 {
				return nil, ThrowArityError(env, 2, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			b, _ := utils.AsFloat(args[1])
			if a < b {
				return a, nil
			}
			return b, nil
		})).
		AddStaticMethod("max", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "a", Type: ast.TypeFromString("Number")},
			{Name: "b", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 2 {
				return nil, ThrowArityError(env, 2, len(args))
			}
			a, _ := utils.AsFloat(args[0])
			b, _ := utils.AsFloat(args[1])
			if a > b {
				return a, nil
			}
			return b, nil
		})).
		AddStaticMethod("clamp", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
			{Name: "min", Type: ast.TypeFromString("Number")},
			{Name: "max", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, args []any) (any, error) {
			if len(args) < 3 {
				return nil, ThrowArityError(env, 3, len(args))
			}
			x, _ := utils.AsFloat(args[0])
			lo, _ := utils.AsFloat(args[1])
			hi, _ := utils.AsFloat(args[2])
			if x < lo {
				x = lo
			}
			if x > hi {
				x = hi
			}
			return x, nil
		})).
		AddStaticMethod("random", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, Func(func(_ *Env, _ []any) (any, error) {
			return rnd.Float64(), nil
		}))

	_, err := mathClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
}

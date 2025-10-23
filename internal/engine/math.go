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
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Abs(a), nil
		})).
		AddStaticMethod("floor", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Floor(a), nil
		})).
		AddStaticMethod("ceil", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Ceil(a), nil
		})).
		AddStaticMethod("round", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Round(a), nil
		})).
		AddStaticMethod("sqrt", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Sqrt(a), nil
		})).
		AddStaticMethod("sin", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Sin(a), nil
		})).
		AddStaticMethod("cos", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Cos(a), nil
		})).
		AddStaticMethod("tan", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			x, _ := env.Get("x")
			a, _ := utils.AsFloat(x)
			return math.Tan(a), nil
		})).
		AddStaticMethod("pow", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "base", Type: ast.TypeFromString("Number")},
			{Name: "exponent", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			base, _ := env.Get("base")
			exponent, _ := env.Get("exponent")
			a, _ := utils.AsFloat(base)
			b, _ := utils.AsFloat(exponent)
			return math.Pow(a, b), nil
		})).
		AddStaticMethod("min", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "a", Type: ast.TypeFromString("Number")},
			{Name: "b", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			aVal, _ := env.Get("a")
			bVal, _ := env.Get("b")
			a, _ := utils.AsFloat(aVal)
			b, _ := utils.AsFloat(bVal)
			if a < b {
				return a, nil
			}
			return b, nil
		})).
		AddStaticMethod("max", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "a", Type: ast.TypeFromString("Number")},
			{Name: "b", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			aVal, _ := env.Get("a")
			bVal, _ := env.Get("b")
			a, _ := utils.AsFloat(aVal)
			b, _ := utils.AsFloat(bVal)
			if a > b {
				return a, nil
			}
			return b, nil
		})).
		AddStaticMethod("clamp", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "x", Type: ast.TypeFromString("Number")},
			{Name: "min", Type: ast.TypeFromString("Number")},
			{Name: "max", Type: ast.TypeFromString("Number")},
		}, Func(func(env *Env, _ []any) (any, error) {
			xVal, _ := env.Get("x")
			minVal, _ := env.Get("min")
			maxVal, _ := env.Get("max")
			x, _ := utils.AsFloat(xVal)
			lo, _ := utils.AsFloat(minVal)
			hi, _ := utils.AsFloat(maxVal)
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

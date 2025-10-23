package engine

import (
	"math"
	"strconv"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallNumberBuiltin installs Int and Float builtin types as classes
func InstallNumberBuiltin(env *Env) error {
	// Install Int class
	if err := installIntClass(env); err != nil {
		return err
	}

	// Install Float class
	if err := installFloatClass(env); err != nil {
		return err
	}

	return nil
}

// installIntClass installs the Int builtin type as a class
func installIntClass(env *Env) error {
	intClass := NewClassBuilder("Int").
		AddField("_value", &ast.Type{Name: "int", IsBuiltin: true}, []string{"private"})

	// utils.ToString() -> String
	intClass.AddBuiltinMethod("toString", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(int)
		return CreateStringInstance((*Env)(callEnv), strconv.Itoa(num))
	}, []string{})

	// abs() -> Int
	intClass.AddBuiltinMethod("abs", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(int)
		if num < 0 {
			return CreateIntInstance((*Env)(callEnv), -num)
		}
		return CreateIntInstance((*Env)(callEnv), num)
	}, []string{})

	// toFloat() -> Float
	intClass.AddBuiltinMethod("toFloat", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(int)
		return CreateFloatInstance((*Env)(callEnv), float64(num))
	}, []string{})

	// serialize() -> String
	intClass.AddBuiltinMethod("serialize", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(int)
		return CreateStringInstance((*Env)(callEnv), strconv.Itoa(num))
	}, []string{})

	// Build the class
	_, err := intClass.Build(env)
	return err
}

// installFloatClass installs the Float builtin type as a class
func installFloatClass(env *Env) error {
	floatClass := NewClassBuilder("Float").
		AddField("_value", &ast.Type{Name: "float", IsBuiltin: true}, []string{"private"})

	// utils.ToString() -> String
	floatClass.AddBuiltinMethod("toString", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateStringInstance((*Env)(callEnv), strconv.FormatFloat(num, 'f', -1, 64))
	}, []string{})

	// abs() -> Float
	floatClass.AddBuiltinMethod("abs", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateFloatInstance((*Env)(callEnv), math.Abs(num))
	}, []string{})

	// toInt() -> Int
	floatClass.AddBuiltinMethod("toInt", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateIntInstance((*Env)(callEnv), int(num))
	}, []string{})

	// floor() -> Float
	floatClass.AddBuiltinMethod("floor", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateFloatInstance((*Env)(callEnv), math.Floor(num))
	}, []string{})

	// ceil() -> Float
	floatClass.AddBuiltinMethod("ceil", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateFloatInstance((*Env)(callEnv), math.Ceil(num))
	}, []string{})

	// round() -> Float
	floatClass.AddBuiltinMethod("round", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateFloatInstance((*Env)(callEnv), math.Round(num))
	}, []string{})

	// sqrt() -> Float
	floatClass.AddBuiltinMethod("sqrt", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateFloatInstance((*Env)(callEnv), math.Sqrt(num))
	}, []string{})

	// serialize() -> String
	floatClass.AddBuiltinMethod("serialize", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		num := instance.Fields["_value"].(float64)
		return CreateStringInstance((*Env)(callEnv), strconv.FormatFloat(num, 'f', -1, 64))
	}, []string{})

	// Build the class
	_, err := floatClass.Build(env)
	return err
}

// CreateIntInstance creates an Int instance from a Go int
// This is used when evaluating integer literals
func CreateIntInstance(env *Env, value int) (*ClassInstance, error) {
	intClassVal, ok := env.Get("__IntClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Int class")
	}

	intClass := intClassVal.(*ClassDefinition)

	// Create instance
	instance, err := createClassInstance(intClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_value"] = value

	return classInstance, nil
}

// CreateFloatInstance creates a Float instance from a Go float64
// This is used when evaluating float literals
func CreateFloatInstance(env *Env, value float64) (*ClassInstance, error) {
	floatClassVal, ok := env.Get("__FloatClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Float class")
	}

	floatClass := floatClassVal.(*ClassDefinition)

	// Create instance
	instance, err := createClassInstance(floatClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_value"] = value

	return classInstance, nil
}

// IntValue extracts the Go int value from an Int instance or converts a value to int
func IntValue(v any) (int, bool) {
	switch val := v.(type) {
	case *ClassInstance:
		if val.ClassName == "Int" {
			if intVal, ok := val.Fields["_value"].(int); ok {
				return intVal, true
			}
		}
		return 0, false
	case int:
		return val, true
	case float64:
		return int(val), true
	default:
		return 0, false
	}
}

// FloatValue extracts the Go float64 value from a Float instance or converts a value to float64
func FloatValue(v any) (float64, bool) {
	switch val := v.(type) {
	case *ClassInstance:
		if val.ClassName == "Float" {
			if floatVal, ok := val.Fields["_value"].(float64); ok {
				return floatVal, true
			}
		}
		return 0, false
	case float64:
		return val, true
	case int:
		return float64(val), true
	default:
		return 0, false
	}
}

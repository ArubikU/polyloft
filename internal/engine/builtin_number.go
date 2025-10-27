package engine

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallNumberBuiltin installs Int and Float builtin types as classes
func InstallNumberBuiltin(env *Env) error {
	// Install Number interface first
	if err := installNumberInterface(env); err != nil {
		return err
	}

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

// installNumberInterface creates a sealed Number interface that Int and Float must implement
func installNumberInterface(env *Env) error {
	numberInterface := &common.InterfaceDefinition{
		Name:         "Number",
		Type:         &ast.Type{Name: "Number", IsBuiltin: true, IsInterface: true},
		Methods:      make(map[string][]common.MethodSignature),
		StaticFields: make(map[string]any),
		AccessLevel:  "public",
		IsSealed:     true, // Sealed - only Integer and Float can implement
		FileName:     "builtin",
		PackageName:  "polyloft.lang",
	}

	// Add toFloat method signature - all numbers must be convertible to float
	numberInterface.Methods["toFloat"] = []common.MethodSignature{{
		Name:       "toFloat",
		Params:     []ast.Parameter{},
		ReturnType: &ast.Type{Name: "float", IsBuiltin: true},
		HasDefault: false,
	}}

	// Register the Number interface
	interfaceRegistry["Number"] = numberInterface
	env.Set("Number", numberInterface)
	env.Set("__NumberInterface__", numberInterface)

	return nil
}

// installIntClass installs the Integer builtin type as a class
func installIntClass(env *Env) error {
	// Get Number interface
	numberInterface, ok := interfaceRegistry["Number"]
	if !ok {
		return fmt.Errorf("Number interface not found")
	}

	intClass := NewClassBuilder("Integer").
		AddAlias("Int").
		AddAlias("int").
		AddAlias("integer").
		AddInterface(numberInterface).
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

	// Constructor: Integer() - no args, initialize to 0
	intClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_value"] = 0
		return nil, nil
	})

	// Constructor: Integer(value: int)
	intClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "value", Type: &ast.Type{Name: "int", IsBuiltin: true}},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		// Get int value from argument
		intVal, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "Integer", args[0])
		}
		instance.Fields["_value"] = intVal
		return nil, nil
	})

	// Build the class
	_, err := intClass.Build(env)
	return err
}

// installFloatClass installs the Float builtin type as a class
func installFloatClass(env *Env) error {
	// Get Number interface
	numberInterface, ok := interfaceRegistry["Number"]
	if !ok {
		return fmt.Errorf("Number interface not found")
	}

	floatClass := NewClassBuilder("Float").
		AddAlias("float").
		AddInterface(numberInterface).
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

	// toFloat() -> Float (required by Number interface)
	floatClass.AddBuiltinMethod("toFloat", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		return thisVal, nil // Float already is a float
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

	// Constructor: Float() - no args, initialize to 0.0
	floatClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_value"] = 0.0
		return nil, nil
	})

	// Constructor: Float(value: float)
	floatClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "value", Type: &ast.Type{Name: "float", IsBuiltin: true}},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		// Get float value from argument
		floatVal, ok := utils.AsFloat(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "Float", args[0])
		}
		instance.Fields["_value"] = floatVal
		return nil, nil
	})

	// Build the class
	_, err := floatClass.Build(env)
	return err
}

// CreateIntInstance creates an Integer instance from a Go int
// This is used when evaluating integer literals
func CreateIntInstance(env *Env, value int) (*ClassInstance, error) {
	intClassVal, ok := env.Get("__IntegerClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Integer class")
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

// IntValue extracts the Go int value from an Integer instance or converts a value to int
func IntValue(v any) (int, bool) {
	switch val := v.(type) {
	case *ClassInstance:
		// Check for both "Integer" (canonical) and "Int" (legacy)
		if val.ClassName == "Integer" || val.ClassName == "Int" {
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

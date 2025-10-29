package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallBoolBuiltin installs Bool builtin type as a class
func InstallBoolBuiltin(env *Env) error {
	// Step 1: Create basic class structure first with fields
	boolClass := NewClassBuilder("Bool")

	// Get type reference for native bool field (not a registered type)
	nativeBoolType := &ast.Type{Name: "bool", IsBuiltin: true}
	boolClass.AddField("_value", nativeBoolType, []string{"private"})

	// Step 2: Now get type references for method signatures
	boolType := boolClass.GetType()
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)

	// utils.ToString() -> String
	boolClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		val := instance.Fields["_value"].(bool)
		if val {
			return CreateStringInstance((*Env)(callEnv), "true")
		}
		return CreateStringInstance((*Env)(callEnv), "false")
	}, []string{})

	// not() -> Bool
	boolClass.AddBuiltinMethod("not", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		val := instance.Fields["_value"].(bool)
		return CreateBoolInstance((*Env)(callEnv), !val)
	}, []string{})

	// serialize() -> String
	boolClass.AddBuiltinMethod("serialize", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		val := instance.Fields["_value"].(bool)
		if val {
			return CreateStringInstance((*Env)(callEnv), "true")
		}
		return CreateStringInstance((*Env)(callEnv), "false")
	}, []string{})

	// Build the class
	_, err := boolClass.Build(env)
	return err
}

// CreateBoolInstance creates a Bool instance from a Go bool
// This is used when evaluating boolean literals
func CreateBoolInstance(env *Env, value bool) (*ClassInstance, error) {
	boolClass := common.BuiltinTypeBool.GetClassDefinition(env)
	if boolClass == nil {
		return nil, ThrowInitializationError(env, "Bool class")
	}

	// Create instance
	instance, err := createClassInstance(boolClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_value"] = value

	return classInstance, nil
}

// BoolValue extracts the Go bool value from a Bool instance or converts a value to bool
func BoolValue(v any) (bool, bool) {
	switch val := v.(type) {
	case *ClassInstance:
		if val.ClassName == "Bool" {
			if boolVal, ok := val.Fields["_value"].(bool); ok {
				return boolVal, true
			}
		}
		return false, false
	case bool:
		return val, true
	default:
		return false, false
	}
}

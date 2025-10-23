package engine

import (
	"fmt"
	"strconv"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallGenericBuiltin installs the Generic builtin type
// Generic wraps any native Go value and provides conversion methods
func InstallGenericBuiltin(env *Env) error {
	genericClass := NewClassBuilder("Generic").
		AddField("_value", ast.ANY, []string{"private"})

	// Get type references
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	floatType := common.BuiltinTypeFloat.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)

	// toString() -> String - Convert to string
	genericClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		value := instance.Fields["_value"]
		return fmt.Sprintf("%v", value), nil
	}, []string{})

	// toInt() -> Int - Convert to int
	genericClass.AddBuiltinMethod("toInt", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		value := instance.Fields["_value"]

		switch v := value.(type) {
		case int:
			return v, nil
		case int32:
			return int(v), nil
		case int64:
			return int(v), nil
		case float32:
			return int(v), nil
		case float64:
			return int(v), nil
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i, nil
			}
			return 0, ThrowConversionError((*Env)(callEnv), v, "int")
		case bool:
			if v {
				return 1, nil
			}
			return 0, nil
		default:
			return 0, ThrowConversionError((*Env)(callEnv), v, "int")
		}
	}, []string{})

	// toFloat() -> Float - Convert to float
	genericClass.AddBuiltinMethod("toFloat", floatType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		value := instance.Fields["_value"]

		switch v := value.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int32:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, nil
			}
			return 0.0, ThrowConversionError((*Env)(callEnv), v, "float")
		case bool:
			if v {
				return 1.0, nil
			}
			return 0.0, nil
		default:
			return 0.0, ThrowConversionError((*Env)(callEnv), v, "float")
		}
	}, []string{})

	// toBool() -> Bool - Convert to bool
	genericClass.AddBuiltinMethod("toBool", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		value := instance.Fields["_value"]

		switch v := value.(type) {
		case bool:
			return v, nil
		case int, int32, int64:
			return v != 0, nil
		case float32, float64:
			return v != 0.0, nil
		case string:
			return v != "" && v != "false" && v != "0", nil
		default:
			return value != nil, nil
		}
	}, []string{})

	// getValue() -> Any - Get the wrapped value
	genericClass.AddBuiltinMethod("getValue", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		return instance.Fields["_value"], nil
	}, []string{})

	// type() -> String - Get the type of the wrapped value
	genericClass.AddBuiltinMethod("type", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		value := instance.Fields["_value"]

		switch value.(type) {
		case int, int32, int64:
			return "int", nil
		case float32, float64:
			return "float", nil
		case string:
			return "string", nil
		case bool:
			return "bool", nil
		case []any:
			return "array", nil
		case map[string]any:
			return "map", nil
		case nil:
			return "nil", nil
		case common.Func, *common.FunctionDefinition, *common.LambdaDefinition:
			return "function", nil
		default:
			return "unknown", nil
		}
	}, []string{})

	// Store the class definition for use when creating Generic instances
	_, err := genericClass.Build(env)
	if err != nil {
		return err
	}

	// Store it so it can be retrieved
	genericDef, _ := env.Get("Generic")
	env.Set("__GenericClass__", genericDef)

	return nil
}

// CreateGenericInstance creates a Generic instance wrapping a native Go value
func CreateGenericInstance(env *Env, value any) (any, error) {
	genericClass := common.BuiltinTypeGeneric.GetClassDefinition(env)
	if genericClass == nil {
		return nil, ThrowRuntimeError(env, "Generic class not found")
	}

	instance := &ClassInstance{
		ClassName: "Generic",
		Fields: map[string]any{
			"_value": value,
		},
		Methods:     make(map[string]common.Func),
		ParentClass: genericClass,
	}

	// Add methods from the class definition
	for methodName, methodInfos := range genericClass.Methods {
		if len(methodInfos) > 0 {
			methodInfo := methodInfos[0] // Take first overload
			if methodInfo.BuiltinImpl != nil {
				instance.Methods[methodName] = methodInfo.BuiltinImpl
			}
		}
	}

	return instance, nil
}

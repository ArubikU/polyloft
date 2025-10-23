package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallTupleClass installs the Tuple builtin class
// Tuple is an immutable sequence that implements Unstructured interface
func InstallTupleClass(env *Env) error {
	// Get interface definitions
	unstructuredInterface := common.BuiltinInterfaceUnstructured.GetInterfaceDefinition(env)

	// Step 1: Create basic class structure (name, type params, interfaces, fields)
	tupleBuilder := NewClassBuilder("Tuple")
	tupleBuilder.AddTypeParameter("T", []string{}, false)
	tupleBuilder.AddInterface(unstructuredInterface)
	
	// Add field for storing tuple elements (native Go slice)
	tupleBuilder.AddField("_elements", &ast.Type{Name: "[]any", IsBuiltin: true}, []string{"private"})

	// Step 2: Get type references for method signatures
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	anyType := ast.ANY
	tupleType := tupleBuilder.GetType()

	// Add constructor - Tuple accepts variadic arguments
	// No formal constructor needed as we'll handle it in Build

	// size() -> Int - returns number of elements
	tupleBuilder.AddBuiltinMethod("size", intType, []ast.Parameter{},
		common.Func(func(env *common.Env, args []any) (any, error) {
			thisVal, _ := env.Get("this")
			inst := thisVal.(*common.ClassInstance)
			elements := inst.Fields["_elements"].([]any)
			return len(elements), nil
		}), []string{})

	// get(index: Int) -> T - get element at index
	tupleBuilder.AddBuiltinMethod("get", anyType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		thisVal, _ := env.Get("this")
		inst := thisVal.(*common.ClassInstance)
		elements := inst.Fields["_elements"].([]any)
		
		index := args[0].(int)
		if index < 0 || index >= len(elements) {
			return nil, ThrowRuntimeError((*Env)(env), fmt.Sprintf("Tuple index out of bounds: %d (size: %d)", index, len(elements)))
		}
		
		return elements[index], nil
	}), []string{})

	// pieces() -> Int - Unstructured interface method
	tupleBuilder.AddBuiltinMethod("pieces", intType, []ast.Parameter{},
		common.Func(func(env *common.Env, args []any) (any, error) {
			thisVal, _ := env.Get("this")
			inst := thisVal.(*common.ClassInstance)
			elements := inst.Fields["_elements"].([]any)
			return len(elements), nil
		}), []string{})

	// getPiece(index: Int) -> Any - Unstructured interface method
	tupleBuilder.AddBuiltinMethod("getPiece", anyType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		thisVal, _ := env.Get("this")
		inst := thisVal.(*common.ClassInstance)
		elements := inst.Fields["_elements"].([]any)
		
		index := args[0].(int)
		if index < 0 || index >= len(elements) {
			return nil, ThrowRuntimeError((*Env)(env), fmt.Sprintf("Tuple index out of bounds: %d (size: %d)", index, len(elements)))
		}
		
		return elements[index], nil
	}), []string{})

	// toArray() -> Array - convert tuple to array
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)
	tupleBuilder.AddBuiltinMethod("toArray", arrayType, []ast.Parameter{},
		common.Func(func(env *common.Env, args []any) (any, error) {
			thisVal, _ := env.Get("this")
			inst := thisVal.(*common.ClassInstance)
			elements := inst.Fields["_elements"].([]any)
			
			// Create a copy
			result := make([]any, len(elements))
			copy(result, elements)
			return result, nil
		}), []string{})

	// toString() -> String - string representation
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	tupleBuilder.AddBuiltinMethod("toString", stringType, []ast.Parameter{},
		common.Func(func(env *common.Env, args []any) (any, error) {
			thisVal, _ := env.Get("this")
			inst := thisVal.(*common.ClassInstance)
			elements := inst.Fields["_elements"].([]any)
			
			var parts []string
			for _, elem := range elements {
				parts = append(parts, fmt.Sprintf("%v", elem))
			}
			return "(" + strings.Join(parts, ", ") + ")", nil
		}), []string{})

	// Build and install the class
	tupleBuilder.Build(env)

	// Store class reference
	env.Define("__TupleClass__", tupleType, "final")

	return nil
}

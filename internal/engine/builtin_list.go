package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallListBuiltin installs the List<T> builtin class
func InstallListBuiltin(env *Env) error {
	// Helper function to create array from items
	createArrayFromItems := func(itemsPtr *[]any, callEnv *common.Env) (any, error) {
		result := make([]any, len(*itemsPtr))
		copy(result, *itemsPtr)
		return CreateArrayInstance(callEnv, result)
	}

	// Step 1: Create basic class structure first with interfaces and fields
	listClass := NewClassBuilder("List").
		AddTypeParameters(common.TBound.AsGenericType().AsArray())

	// Get interface references
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
	collectionInterface := common.BuiltinInterfaceCollection.GetInterfaceDefinition(env)
	unstructuredInterface := common.BuiltinInterfaceUnstructured.GetInterfaceDefinition(env)

	// Add interfaces
	listClass.AddInterface(iterableInterface)
	listClass.AddInterface(collectionInterface)
	listClass.AddInterface(unstructuredInterface)

	// Get type references for fields
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)

	// Add fields
	listClass.AddField("_items", arrayType, []string{"private"})
	listClass.AddField("_currentIndex", intType, []string{"private"})

	// Step 2: Now get types for method signatures
	tType := &ast.Type{Name: "T"} // Generic type parameter
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)

	// Constructor: List() - empty list
	listClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// Constructor: List(array) - from array
	listClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "array", Type: arrayType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		if arr, ok := args[0].([]any); ok {
			items := make([]any, len(arr))
			copy(items, arr)
			instance.Fields["_items"] = &items
		} else {
			items := make([]any, 0)
			instance.Fields["_items"] = &items
		}
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// Constructor: List(items...) - variadic
	listClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "items", Type: ast.ANY, IsVariadic: true},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := make([]any, len(args))
		copy(items, args)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// size() -> Int
	listClass.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return len(*itemsPtr), nil
	}, []string{})

	// isEmpty() -> Bool
	listClass.AddBuiltinMethod("isEmpty", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return len(*itemsPtr) == 0, nil
	}, []string{})

	// contains(item: T) -> Bool
	listClass.AddBuiltinMethod("contains", boolType, []ast.Parameter{
		{Name: "item", Type: tType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		for _, item := range *itemsPtr {
			if equal(item, args[0]) {
				return true, nil
			}
		}
		return false, nil
	}, []string{})

	// add(item: T) -> Void
	listClass.AddBuiltinMethod("add", ast.ANY, []ast.Parameter{
		{Name: "item", Type: tType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		*itemsPtr = append(*itemsPtr, args[0])
		return nil, nil
	}, []string{})

	// get(index: Int) -> T
	listClass.AddBuiltinMethod("get", tType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if idx < 0 || idx >= len(*itemsPtr) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(*itemsPtr), "List")
		}
		return (*itemsPtr)[idx], nil
	}, []string{})

	// set(index: Int, value: T) -> Void
	listClass.AddBuiltinMethod("set", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
		{Name: "value", Type: tType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if idx < 0 || idx >= len(*itemsPtr) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(*itemsPtr), "List")
		}
		(*itemsPtr)[idx] = args[1]
		return nil, nil
	}, []string{})

	// remove(index: Int) -> T
	listClass.AddBuiltinMethod("remove", tType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if idx < 0 || idx >= len(*itemsPtr) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(*itemsPtr), "List")
		}
		removed := (*itemsPtr)[idx]
		*itemsPtr = append((*itemsPtr)[:idx], (*itemsPtr)[idx+1:]...)
		return removed, nil
	}, []string{})

	// clear() -> Void
	listClass.AddBuiltinMethod("clear", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	}, []string{})

	// toArray() -> Array
	listClass.AddBuiltinMethod("toArray", arrayType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return createArrayFromItems(itemsPtr, callEnv)
	}, []string{})

	// asArray() -> Array - alias for toArray
	listClass.AddBuiltinMethod("asArray", arrayType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return createArrayFromItems(itemsPtr, callEnv)
	}, []string{})

	// Iterable interface methods
	// hasNext() -> Bool
	listClass.AddBuiltinMethod("hasNext", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		currentIndex := instance.Fields["_currentIndex"].(int)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return currentIndex < len(*itemsPtr), nil
	}, []string{})

	// next() -> T
	listClass.AddBuiltinMethod("next", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		currentIndex := instance.Fields["_currentIndex"].(int)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if currentIndex >= len(*itemsPtr) {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Iterator exhausted")
		}
		result := (*itemsPtr)[currentIndex]
		instance.Fields["_currentIndex"] = currentIndex + 1
		return result, nil
	}, []string{})

	// Unstructured interface methods
	// pieces() -> Int - returns number of elements
	listClass.AddBuiltinMethod("__pieces", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return len(*itemsPtr), nil
	}, []string{})

	// getPiece(index: Int) -> Any
	listClass.AddBuiltinMethod("__get_piece", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if idx < 0 || idx >= len(*itemsPtr) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(*itemsPtr), "List")
		}
		return (*itemsPtr)[idx], nil
	}, []string{})

	// forEach(callback: Function) -> Void
	listClass.AddBuiltinMethod("forEach", ast.TypeFromString("Void"), []ast.Parameter{
		{Name: "callback", Type: ast.TypeFromString("Function")},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)

		// Extract the callback function
		callback, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		// Call the callback for each item
		for _, item := range *itemsPtr {
			_, err := callback(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}, []string{})

	// toString() -> String
	listClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)

		strs := make([]string, len(*itemsPtr))
		for i, item := range *itemsPtr {
			strs[i] = utils.ToString(item)
		}

		// Build type string with generic parameters if present
		typeStr := GetTypeName(thisVal)

		return CreateStringInstance(callEnv, fmt.Sprintf("%s(%s)", typeStr, strings.Join(strs, ", ")))
	}, []string{})

	// Build and register
	_, err := listClass.Build(env)
	return err
}

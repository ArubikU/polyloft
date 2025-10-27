package engine

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallArrayBuiltin installs the Array builtin type using ClassBuilder
// Array is now a minimal base class with only basic operations
func InstallArrayBuiltin(env *Env) error {
	// Step 1: Create basic class structure first with interfaces and fields
	arrayClass := NewClassBuilder("Array").
		AddAlias("array").
		AddTypeParameters(common.TBound.AsGenericType().AsArray())

	// Get interface references
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
	unstructuredInterface := common.BuiltinInterfaceUnstructured.GetInterfaceDefinition(env)
	sliceableInterface := common.BuiltinSliceableInterface.GetInterfaceDefinition(env)
	indexableInterface := common.BuiltinIndexableInterface.GetInterfaceDefinition(env)
	arrayClass.AddInterface(iterableInterface)
	arrayClass.AddInterface(unstructuredInterface)
	arrayClass.AddInterface(sliceableInterface)
	arrayClass.AddInterface(indexableInterface)

	// Step 2: Now get type references for method signatures
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	arrayType := arrayClass.GetType()

	arrayClass.AddField("_items", arrayType, []string{"private"})

	//Iterable interface methods , __length, __get, __get_step
	arrayClass.AddBuiltinMethod("__length", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		return len(items), nil
	}, []string{})
	arrayClass.AddBuiltinMethod("__get", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		if idx < 0 || idx >= len(items) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(items), "Array")
		}
		return items[idx], nil
	}, []string{})

	// __set(index: int, value: T) -> Void (Indexable interface)
	arrayClass.AddBuiltinMethod("__set", ast.NIL, []ast.Parameter{
		{Name: "index", Type: intType},
		{Name: "value", Type: &ast.Type{Name: "T"}},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		if idx < 0 || idx >= len(items) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(items), "Array")
		}
		items[idx] = args[1]
		return nil, nil
	}, []string{})

	// __contains(index: int) -> Bool (Indexable interface)
	arrayClass.AddBuiltinMethod("__contains", boolType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		idx, ok := utils.AsInt(args[0])
		if !ok {
			return false, nil
		}

		return idx >= 0 && idx < len(items), nil
	}, []string{})

	// length() -> Int
	arrayClass.AddBuiltinMethod("length", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		return len(items), nil
	}, []string{})

	//Sliceable interface method __slice
	arrayClass.AddBuiltinMethod("__slice", ast.ANY, []ast.Parameter{
		{Name: "start", Type: ast.ANY},
		{Name: "end", Type: ast.ANY},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		startIdx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		endIdx, ok := utils.AsInt(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[1])
		}
		if startIdx < 0 || endIdx > len(items) || startIdx > endIdx {
			return nil, ThrowIndexError((*Env)(callEnv), startIdx, len(items), "Array")
		}
		return items[startIdx:endIdx], nil
	}, []string{})

	// isEmpty() -> Bool
	arrayClass.AddBuiltinMethod("isEmpty", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		return len(items) == 0, nil
	}, []string{})

	// add(item: any) -> Void
	arrayClass.AddBuiltinMethod("add", ast.NIL, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		instance.Fields["_items"] = append(items, args[0])
		return nil, nil
	}, []string{})

	// push(item: any) -> Void (alias for add)
	arrayClass.AddBuiltinMethod("push", ast.NIL, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		instance.Fields["_items"] = append(items, args[0])
		return nil, nil
	}, []string{})

	// pop() -> any
	arrayClass.AddBuiltinMethod("pop", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		if len(items) == 0 {
			return nil, nil
		}
		last := items[len(items)-1]
		instance.Fields["_items"] = items[:len(items)-1]
		return last, nil
	}, []string{})

	// shift() -> any
	arrayClass.AddBuiltinMethod("shift", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		if len(items) == 0 {
			return nil, nil
		}
		first := items[0]
		instance.Fields["_items"] = items[1:]
		return first, nil
	}, []string{})

	// unshift(item: any) -> Void
	arrayClass.AddBuiltinMethod("unshift", ast.NIL, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		instance.Fields["_items"] = append([]any{args[0]}, items...)
		return nil, nil
	}, []string{})

	// get(index: int) -> any
	arrayClass.AddBuiltinMethod("get", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		if idx < 0 || idx >= len(items) {
			return nil, nil
		}
		return items[idx], nil
	}, []string{})

	// set(index: int, value: any) -> Void
	arrayClass.AddBuiltinMethod("set", ast.NIL, []ast.Parameter{
		{Name: "index", Type: intType},
		{Name: "value", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		if idx < 0 || idx >= len(items) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(items), "Array")
		}
		items[idx] = args[1]
		return nil, nil
	}, []string{})

	// indexOf(item: any) -> int
	arrayClass.AddBuiltinMethod("indexOf", intType, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		for i, item := range items {
			if equals(item, args[0]) {
				return i, nil
			}
		}
		return -1, nil
	}, []string{})

	// contains(item: any) -> bool
	arrayClass.AddBuiltinMethod("contains", boolType, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		for _, item := range items {
			if equals(item, args[0]) {
				return true, nil
			}
		}
		return false, nil
	}, []string{})

	// slice(start: int, end: int) -> Array
	arrayClass.AddBuiltinMethod("slice", arrayClass.GetType(), []ast.Parameter{
		{Name: "start", Type: intType},
		{Name: "end", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		start, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		end, ok := utils.AsInt(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[1])
		}

		if start < 0 {
			start = 0
		}
		if end > len(items) {
			end = len(items)
		}
		if start > end {
			start = end
		}

		sliced := items[start:end]
		return CreateArrayInstance((*Env)(callEnv), sliced)
	}, []string{})

	// join(separator: string) -> string
	arrayClass.AddBuiltinMethod("join", stringType, []ast.Parameter{
		{Name: "separator", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		sep := utils.ToString(args[0])
		parts := make([]string, len(items))
		for i, item := range items {
			parts[i] = utils.ToString(item)
		}
		return strings.Join(parts, sep), nil
	}, []string{})

	// reverse() -> Void
	arrayClass.AddBuiltinMethod("reverse", ast.NIL, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
		return nil, nil
	}, []string{})

	// sort() -> Void
	arrayClass.AddBuiltinMethod("sort", ast.NIL, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		sort.Slice(items, func(i, j int) bool {
			return utils.ToString(items[i]) < utils.ToString(items[j])
		})
		return nil, nil
	}, []string{})

	// filter(fn: Function) -> Array
	arrayClass.AddBuiltinMethod("filter", arrayClass.GetType(), []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		result := []any{}
		for _, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			if utils.AsBool(val) {
				result = append(result, item)
			}
		}
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// map(fn: Function) -> Array
	arrayClass.AddBuiltinMethod("map", arrayClass.GetType(), []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		result := make([]any, len(items))
		for i, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// forEach(fn: Function) -> Void
	arrayClass.AddBuiltinMethod("forEach", ast.NIL, []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		for _, item := range items {
			_, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}, []string{})

	// reduce(fn: Function, initial: any) -> any
	arrayClass.AddBuiltinMethod("reduce", ast.ANY, []ast.Parameter{
		{Name: "fn", Type: nil},
		{Name: "initial", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		accumulator := args[1]
		for _, item := range items {
			val, err := fn(callEnv, []any{accumulator, item})
			if err != nil {
				return nil, err
			}
			accumulator = val
		}
		return accumulator, nil
	}, []string{})

	// find(fn: Function) -> any
	arrayClass.AddBuiltinMethod("find", ast.ANY, []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		for _, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			if utils.AsBool(val) {
				return item, nil
			}
		}
		return nil, nil
	}, []string{})

	// findIndex(fn: Function) -> int
	arrayClass.AddBuiltinMethod("findIndex", intType, []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		for i, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			if utils.AsBool(val) {
				return i, nil
			}
		}
		return -1, nil
	}, []string{})

	// every(fn: Function) -> bool
	arrayClass.AddBuiltinMethod("every", boolType, []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		for _, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			if !utils.AsBool(val) {
				return false, nil
			}
		}
		return true, nil
	}, []string{})

	// some(fn: Function) -> bool
	arrayClass.AddBuiltinMethod("some", boolType, []ast.Parameter{
		{Name: "fn", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		for _, item := range items {
			val, err := fn(callEnv, []any{item})
			if err != nil {
				return nil, err
			}
			if utils.AsBool(val) {
				return true, nil
			}
		}
		return false, nil
	}, []string{})

	// concat(other: Array) -> Array
	arrayClass.AddBuiltinMethod("concat", arrayClass.GetType(), []ast.Parameter{
		{Name: "other", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		var otherItems []any
		if otherInstance, ok := args[0].(*ClassInstance); ok && otherInstance.ClassName == "Array" {
			otherItems = otherInstance.Fields["_items"].([]any)
		} else if arr, ok := args[0].([]any); ok {
			otherItems = arr
		} else {
			return nil, ThrowTypeError((*Env)(callEnv), "Array", args[0])
		}

		result := make([]any, len(items)+len(otherItems))
		copy(result, items)
		copy(result[len(items):], otherItems)
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// clear() -> Void
	arrayClass.AddBuiltinMethod("clear", ast.NIL, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_items"] = []any{}
		return nil, nil
	}, []string{})

	// utils.ToString() -> String
	arrayClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		parts := make([]string, len(items))
		for i, item := range items {
			parts[i] = utils.ToString(item)
		}
		return CreateStringInstance(callEnv, "["+strings.Join(parts, ", ")+"]")
	}, []string{})

	// serialize() -> String
	arrayClass.AddBuiltinMethod("serialize", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		arr, err := ArrayToSlice(instance)
		if err != nil {
			return nil, err
		}

		jsonBytes, err := json.Marshal(arr)
		if err != nil {
			return nil, err
		}

		return string(jsonBytes), nil
	}, []string{})

	// pieces() -> Int (Unstructured interface method)
	arrayClass.AddBuiltinMethod("__pieces", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)
		return len(items), nil
	}, []string{})

	// getPiece(index: Int) -> Any (Unstructured interface method)
	arrayClass.AddBuiltinMethod("__get_piece", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		items := instance.Fields["_items"].([]any)

		index, ok := args[0].(int)
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		if index < 0 || index >= len(items) {
			return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("Array index out of bounds: %d (size: %d)", index, len(items)))
		}

		return items[index], nil
	}, []string{})

	// Static methods
	arrayClass.AddStaticMethod("deserialize", arrayType, []ast.Parameter{
		{Name: "data", Type: stringType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		jsonStr := utils.ToString(args[0])

		var data []any
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			return nil, err
		}

		return CreateArrayInstance(env, data)
	}))

	// Build the class
	arrayDef, err := arrayClass.Build(env)
	if err != nil {
		return err
	}

	// Store the definition for use when creating [] literals
	env.Set("__ArrayClass__", arrayDef)

	return nil
}

// CreateArrayInstance creates an Array instance from a []any
// ConvertToClassInstance converts native Go types to their ClassInstance equivalents recursively
func ConvertToClassInstance(env *Env, value any) any {
	if value == nil {
		return nil
	}

	// If already a ClassInstance, return as is
	if inst, ok := value.(*ClassInstance); ok {
		return inst
	}

	// Try to determine the type and use the instance creator
	switch v := value.(type) {
	case string:
		if inst, err := CreateInstanceFor(env, "String", v); err == nil {
			return inst
		}
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		intVal, _ := utils.AsInt(v)
		if inst, err := CreateInstanceFor(env, "Int", intVal); err == nil {
			return inst
		}
		return v
	case float32, float64:
		floatVal, _ := utils.AsFloat(v)
		if inst, err := CreateInstanceFor(env, "Float", floatVal); err == nil {
			return inst
		}
		return v
	case bool:
		if inst, err := CreateInstanceFor(env, "Bool", v); err == nil {
			return inst
		}
		return v
	case []byte:
		if inst, err := CreateInstanceFor(env, "Bytes", v); err == nil {
			return inst
		}
		return v
	case []any:
		// Convert slice items recursively
		convertedSlice := make([]any, len(v))
		for i, item := range v {
			convertedSlice[i] = ConvertToClassInstance(env, item)
		}
		if inst, err := createArrayInstanceDirect(env, convertedSlice); err == nil {
			return inst
		}
		return v
	case map[string]any, map[any]any:
		if inst, err := CreateInstanceFor(env, "Map", v); err == nil {
			return inst
		}
		return v
	default:
		// For other types, keep as is
		return v
	}
}

// createArrayInstanceDirect creates an Array instance without recursive conversion (used internally)
func createArrayInstanceDirect(env *Env, items []any) (*ClassInstance, error) {
	arrayClass := common.BuiltinTypeArray.GetClassDefinition(env)
	if arrayClass == nil {
		return nil, ThrowInitializationError(env, "Array class")
	}

	instance, err := createClassInstance(arrayClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_items"] = items

	return classInstance, nil
}

func CreateArrayInstance(env *Env, items []any) (*ClassInstance, error) {
	arrayClass := common.BuiltinTypeArray.GetClassDefinition(env)
	if arrayClass == nil {
		return nil, ThrowInitializationError(env, "Array class")
	}

	// Create instance
	instance, err := createClassInstance(arrayClass, env, []any{})
	if err != nil {
		return nil, err
	}

	// Convert all items recursively to ClassInstance if they are native types
	convertedItems := make([]any, len(items))
	for i, item := range items {
		convertedItems[i] = ConvertToClassInstance(env, item)
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_items"] = convertedItems

	return classInstance, nil
}

// ArrayToSlice converts an Array instance to a Go []any for JSON serialization
func ArrayToSlice(arrayInstance *ClassInstance) ([]any, error) {
	itemsField, ok := arrayInstance.Fields["_items"]
	if !ok {
		return nil, ThrowAttributeError(nil, "_items", "Array")
	}

	items, ok := itemsField.([]any)
	if !ok {
		return nil, ThrowTypeError(nil, "slice", itemsField)
	}

	// Convert nested Arrays recursively
	result := make([]any, len(items))
	for i, item := range items {
		if nestedArray, ok := item.(*ClassInstance); ok && nestedArray.ClassName == "Array" {
			nestedSlice, err := ArrayToSlice(nestedArray)
			if err != nil {
				return nil, err
			}
			result[i] = nestedSlice
		} else if nestedMap, ok := item.(*ClassInstance); ok && nestedMap.ClassName == "Map" {
			nestedObj, err := MapToObject(nestedMap)
			if err != nil {
				return nil, err
			}
			result[i] = nestedObj
		} else {
			result[i] = item
		}
	}

	return result, nil
}

func ClassInstanceToArray(instance *ClassInstance) ([]any, error) {
	if instance.ClassName != "Array" {
		return nil, ThrowTypeError(nil, "Array", instance.ClassName)
	}

	itemsField, ok := instance.Fields["_items"]
	if !ok {
		return nil, ThrowAttributeError(nil, "_items", "Array")
	}

	items, ok := itemsField.([]any)
	if !ok {
		return nil, ThrowTypeError(nil, "slice", itemsField)
	}

	return items, nil
}
func ClassInstanceToArrayTypeString(instance *ClassInstance) (string, error) {
	if instance.ClassName != "Array" {
		return "", ThrowTypeError(nil, "Array", instance.ClassName)
	}

	items, err := ClassInstanceToArray(instance)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "Array<Any>", nil
	}

	typeSet := make(map[string]struct{})
	hasInt := false
	hasFloat := false

	for _, item := range items {
		typeName := GetTypeName(item)
		switch typeName {
		case "Int":
			hasInt = true
		case "Float":
			hasFloat = true
		default:
			typeSet[typeName] = struct{}{}
		}
	}

	// Combinar Int y Float
	if hasInt && hasFloat {
		typeSet["Number"] = struct{}{}
	} else if hasInt {
		typeSet["Int"] = struct{}{}
	} else if hasFloat {
		typeSet["Float"] = struct{}{}
	}

	// Convertir map a slice
	types := make([]string, 0, len(typeSet))
	for t := range typeSet {
		types = append(types, t)
	}

	// Ordenar alfabéticamente
	sort.Strings(types)

	typeString := strings.Join(types, " | ")
	return fmt.Sprintf("Array<%s>", typeString), nil
}

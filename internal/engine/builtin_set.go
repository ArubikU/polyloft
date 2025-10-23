package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallSetBuiltin installs the Set<T> builtin class
func InstallSetBuiltin(env *Env) error {
	// Get interface and type references
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	mapType := &ast.Type{Name: "map", IsBuiltin: true}
	voidType := &ast.Type{Name: "void", IsBuiltin: true}

	setClass := NewClassBuilder("Set").
		AddTypeParameter("T", []string{}, false).
		AddInterface(iterableInterface).
		AddField("_items", mapType, []string{"private"}). // Using map for O(1) lookups
		AddField("_keys", &ast.Type{Name: "array", IsBuiltin: true}, []string{"private"}). // Track insertion order
		AddField("_currentIndex", intType, []string{"private"})

	// Constructor: Set() - empty set
	setClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		items := make(map[string]bool)
		keys := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_keys"] = &keys
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// Constructor: Set(items...) - variadic
	setClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "items", Type: nil, IsVariadic: true},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		
		items := make(map[string]bool)
		keys := make([]any, 0)
		
		for _, item := range args {
			key := fmt.Sprintf("%v", item)
			if !items[key] {
				items[key] = true
				keys = append(keys, item)
			}
		}
		
		instance.Fields["_items"] = &items
		instance.Fields["_keys"] = &keys
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// size() -> Int
	setClass.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*map[string]bool)
		return len(*itemsPtr), nil
	}, []string{})

	// add(item: T) -> Bool - returns true if item was added (wasn't already present)
	setClass.AddBuiltinMethod("add", boolType, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*map[string]bool)
		keysPtr := instance.Fields["_keys"].(*[]any)
		
		key := fmt.Sprintf("%v", args[0])
		if (*itemsPtr)[key] {
			return false, nil // Already exists
		}
		
		(*itemsPtr)[key] = true
		*keysPtr = append(*keysPtr, args[0])
		return true, nil
	}, []string{})

	// contains(item: T) -> Bool
	setClass.AddBuiltinMethod("contains", boolType, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*map[string]bool)
		
		key := fmt.Sprintf("%v", args[0])
		return (*itemsPtr)[key], nil
	}, []string{})

	// remove(item: T) -> Bool - returns true if item was removed
	setClass.AddBuiltinMethod("remove", boolType, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*map[string]bool)
		keysPtr := instance.Fields["_keys"].(*[]any)
		
		key := fmt.Sprintf("%v", args[0])
		if !(*itemsPtr)[key] {
			return false, nil // Doesn't exist
		}
		
		delete(*itemsPtr, key)
		
		// Remove from keys array
		for i, k := range *keysPtr {
			if fmt.Sprintf("%v", k) == key {
				*keysPtr = append((*keysPtr)[:i], (*keysPtr)[i+1:]...)
				break
			}
		}
		
		return true, nil
	}, []string{})

	// clear() -> Void
	setClass.AddBuiltinMethod("clear", voidType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		items := make(map[string]bool)
		keys := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_keys"] = &keys
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	}, []string{})

	// toArray() -> Array
	setClass.AddBuiltinMethod("toArray", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		keysPtr := instance.Fields["_keys"].(*[]any)
		
		result := make([]any, len(*keysPtr))
		copy(result, *keysPtr)
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// Iterable interface methods
	// hasNext() -> Bool
	setClass.AddBuiltinMethod("hasNext", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		currentIndex := instance.Fields["_currentIndex"].(int)
		keysPtr := instance.Fields["_keys"].(*[]any)
		return currentIndex < len(*keysPtr), nil
	}, []string{})

	// next() -> T
	setClass.AddBuiltinMethod("next", &ast.Type{Name: "T"}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		currentIndex := instance.Fields["_currentIndex"].(int)
		keysPtr := instance.Fields["_keys"].(*[]any)
		if currentIndex >= len(*keysPtr) {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Iterator exhausted")
		}
		result := (*keysPtr)[currentIndex]
		instance.Fields["_currentIndex"] = currentIndex + 1
		return result, nil
	}, []string{})

	// toString() -> String
	stringType := &ast.Type{Name: "string", IsBuiltin: true}
	setClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		keysPtr := instance.Fields["_keys"].(*[]any)
		
		strs := make([]string, len(*keysPtr))
		for i, item := range *keysPtr {
			strs[i] = fmt.Sprintf("%v", item)
		}
		return fmt.Sprintf("{%s}", strings.Join(strs, ", ")), nil
	}, []string{})

	// Build and register
	setDef, err := setClass.Build(env)
	if err != nil {
		return err
	}

	env.Set("__SetClass__", setDef)
	return nil
}

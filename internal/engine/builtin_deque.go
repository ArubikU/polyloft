package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallDequeBuiltin installs the Deque<T> (double-ended queue) builtin class
func InstallDequeBuiltin(env *Env) error {
	// Step 1: Create basic class structure first with interfaces and fields
	dequeClass := NewClassBuilder("Deque").
		AddTypeParameter("T", []string{}, false)
	
	// Get interface references
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
	dequeClass.AddInterface(iterableInterface)
	
	// Get type references for fields
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	nativeArrayType := &ast.Type{Name: "array", IsBuiltin: true}
	
	// Add fields
	dequeClass.AddField("_items", nativeArrayType, []string{"private"})
	dequeClass.AddField("_currentIndex", intType, []string{"private"})
	
	// Step 2: Now get type references for method signatures
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)
	tType := &ast.Type{Name: "T"} // Generic type parameter

	// Constructor: Deque() - empty deque
	dequeClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		items := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// Constructor: Deque(items...) - variadic
	dequeClass.AddBuiltinConstructor([]ast.Parameter{
		{Name: "items", Type: nil, IsVariadic: true},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		
		items := make([]any, len(args))
		copy(items, args)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	})

	// size() -> Int
	dequeClass.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return len(*itemsPtr), nil
	}, []string{})

	// isEmpty() -> Bool
	dequeClass.AddBuiltinMethod("isEmpty", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return len(*itemsPtr) == 0, nil
	}, []string{})

	// addFirst(item: T) -> Void - add to front
	dequeClass.AddBuiltinMethod("addFirst", ast.NIL, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		*itemsPtr = append([]any{args[0]}, *itemsPtr...)
		return nil, nil
	}, []string{})

	// addLast(item: T) -> Void - add to back
	dequeClass.AddBuiltinMethod("addLast", ast.NIL, []ast.Parameter{
		{Name: "item", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		*itemsPtr = append(*itemsPtr, args[0])
		return nil, nil
	}, []string{})

	// removeFirst() -> T - remove from front
	dequeClass.AddBuiltinMethod("removeFirst", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if len(*itemsPtr) == 0 {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Deque is empty")
		}
		result := (*itemsPtr)[0]
		*itemsPtr = (*itemsPtr)[1:]
		return result, nil
	}, []string{})

	// removeLast() -> T - remove from back
	dequeClass.AddBuiltinMethod("removeLast", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if len(*itemsPtr) == 0 {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Deque is empty")
		}
		result := (*itemsPtr)[len(*itemsPtr)-1]
		*itemsPtr = (*itemsPtr)[:len(*itemsPtr)-1]
		return result, nil
	}, []string{})

	// peekFirst() -> T - get front without removing
	dequeClass.AddBuiltinMethod("peekFirst", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if len(*itemsPtr) == 0 {
			return nil, nil
		}
		return (*itemsPtr)[0], nil
	}, []string{})

	// peekLast() -> T - get back without removing
	dequeClass.AddBuiltinMethod("peekLast", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if len(*itemsPtr) == 0 {
			return nil, nil
		}
		return (*itemsPtr)[len(*itemsPtr)-1], nil
	}, []string{})

	// get(index: Int) -> T
	dequeClass.AddBuiltinMethod("get", tType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		if idx < 0 || idx >= len(*itemsPtr) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(*itemsPtr), "Deque")
		}
		return (*itemsPtr)[idx], nil
	}, []string{})

	// clear() -> Void
	dequeClass.AddBuiltinMethod("clear", ast.NIL, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		items := make([]any, 0)
		instance.Fields["_items"] = &items
		instance.Fields["_currentIndex"] = 0
		return nil, nil
	}, []string{})

	// toArray() -> Array
	dequeClass.AddBuiltinMethod("toArray", arrayType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		result := make([]any, len(*itemsPtr))
		copy(result, *itemsPtr)
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// Iterable interface methods
	// hasNext() -> Bool
	dequeClass.AddBuiltinMethod("hasNext", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		currentIndex := instance.Fields["_currentIndex"].(int)
		itemsPtr := instance.Fields["_items"].(*[]any)
		return currentIndex < len(*itemsPtr), nil
	}, []string{})

	// next() -> T
	dequeClass.AddBuiltinMethod("next", tType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
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

	// toString() -> String
	dequeClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		itemsPtr := instance.Fields["_items"].(*[]any)
		
		strs := make([]string, len(*itemsPtr))
		for i, item := range *itemsPtr {
			strs[i] = fmt.Sprintf("%v", item)
		}
		return fmt.Sprintf("Deque[%s]", strings.Join(strs, ", ")), nil
	}, []string{})

	// Build and register
	dequeDef, err := dequeClass.Build(env)
	if err != nil {
		return err
	}

	env.Set("__DequeClass__", dequeDef)
	return nil
}

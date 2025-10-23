package engine

import (
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallRangeBuiltin installs the Range builtin type
// Range is an iterable that doesn't store items in memory, just a counter
func InstallRangeBuiltin(env *Env) error {
	// Get type references from already-installed builtin types
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	stringType := &ast.Type{Name: "string", IsBuiltin: true}

	// Get the Iterable interface definition
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)

	rangeClass := NewClassBuilder("Range").
		AddInterface(iterableInterface).
		AddField("_start", intType, []string{"private"}).
		AddField("_end", intType, []string{"private"}).
		AddField("_step", intType, []string{"private"}).
		AddField("_current", intType, []string{"private"})

	// hasNext() -> Bool - Check if there are more elements
	rangeClass.AddBuiltinMethod("hasNext", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		current := instance.Fields["_current"].(int)
		end := instance.Fields["_end"].(int)
		step := instance.Fields["_step"].(int)

		if step > 0 {
			return current <= end, nil
		}
		return current >= end, nil
	}, []string{})

	// next() -> Int - Get the next element and advance
	rangeClass.AddBuiltinMethod("next", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		current := instance.Fields["_current"].(int)
		step := instance.Fields["_step"].(int)

		// Store current value to return
		result := current

		// Advance to next
		instance.Fields["_current"] = current + step

		return result, nil
	}, []string{})

	// toArray() -> Array - Convert range to array
	rangeClass.AddBuiltinMethod("toArray", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		start := instance.Fields["_start"].(int)
		end := instance.Fields["_end"].(int)
		step := instance.Fields["_step"].(int)

		// Calculate size
		var size int
		if step > 0 {
			size = (end-start)/step + 1
		} else {
			size = (start-end)/(-step) + 1
		}

		if size < 0 {
			size = 0
		}

		// Build array
		items := make([]any, 0, size)
		for i := start; ; i += step {
			if step > 0 && i > end {
				break
			}
			if step < 0 && i < end {
				break
			}
			items = append(items, i)
			if i == end {
				break
			}
		}

		return CreateArrayInstance((*Env)(callEnv), items)
	}, []string{})

	// size() -> Int - Get the number of elements in the range
	rangeClass.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		start := instance.Fields["_start"].(int)
		end := instance.Fields["_end"].(int)
		step := instance.Fields["_step"].(int)

		var size int
		if step > 0 {
			if end >= start {
				size = (end-start)/step + 1
			} else {
				size = 0
			}
		} else {
			if start >= end {
				size = (start-end)/(-step) + 1
			} else {
				size = 0
			}
		}

		return size, nil
	}, []string{})

	// toString() -> String
	rangeClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		start := instance.Fields["_start"].(int)
		end := instance.Fields["_end"].(int)
		step := instance.Fields["_step"].(int)

		if step == 1 {
			return fmt.Sprintf("Range(%d..%d)", start, end), nil
		}
		return fmt.Sprintf("Range(%d..%d step %d)", start, end, step), nil
	}, []string{})

	// Build the class
	rangeDef, err := rangeClass.Build(env)
	if err != nil {
		return err
	}

	// Store the definition for use when creating range literals
	env.Set("__RangeClass__", rangeDef)

	return nil
}

// CreateRangeInstance creates a Range instance from start, end, and optional step
func CreateRangeInstance(env *Env, start, end int, step ...int) (*ClassInstance, error) {
	rangeClassVal, ok := env.Get("__RangeClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Range class")
	}

	rangeClass := rangeClassVal.(*common.ClassDefinition)

	// Default step is 1 if not provided
	stepValue := 1
	if len(step) > 0 {
		stepValue = step[0]
	}

	// If step is not provided and start > end, use -1
	if len(step) == 0 && start > end {
		stepValue = -1
	}

	// Create instance
	instance, err := createClassInstance(rangeClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_start"] = start
	classInstance.Fields["_end"] = end
	classInstance.Fields["_step"] = stepValue
	classInstance.Fields["_current"] = start

	return classInstance, nil
}

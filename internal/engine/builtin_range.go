package engine

import (
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallRangeBuiltin installs the Range builtin type
// Range is an iterable that doesn't store items in memory, just a counter
func InstallRangeBuiltin(env *Env) error {
	// Get type references from already-installed builtin types
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	stringType := &ast.Type{Name: "string", IsBuiltin: true}

	// Get the Iterable interface definition
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)

	rangeClass := NewClassBuilder("Range").
		AddInterface(iterableInterface).
		AddField("_start", ast.ANY, []string{"private"}).
		AddField("_end", ast.ANY, []string{"private"}).
		AddField("_step", ast.ANY, []string{"private"}).

		// __length() -> Int - Get the length of the range
		AddBuiltinMethod("__length", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
			thisVal, _ := callEnv.Get("this")
			instance := thisVal.(*ClassInstance)
			start, _ := utils.AsInt(instance.Fields["_start"])
			end, _ := utils.AsInt(instance.Fields["_end"])
			step, _ := utils.AsInt(instance.Fields["_step"])
			// Calculate length
			length := end - start
			if step != 0 {
				length = length/step + 1
			} else {
				return 0, fmt.Errorf("step cannot be zero")
			}

			return CreateIntInstance(env, length)
		}, []string{})
		// __get(index: Int) -> Int - Get the value at the given index
	rangeClass.AddBuiltinMethod("__get", ast.ANY, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)

		index, _ := utils.AsInt(args[0])
		start, _ := utils.AsInt(instance.Fields["_start"])
		end, _ := utils.AsInt(instance.Fields["_end"])
		step, _ := utils.AsInt(instance.Fields["_step"])

		// Calculate the value at the given index
		var value int
		if step > 0 {
			value = start + index*step
			if value > end {
				return 0, fmt.Errorf("index out of bounds")
			}
		} else {
			value = start + index*step
			if value < end {
				return 0, fmt.Errorf("index out of bounds")
			}
		}

		return CreateIntInstance(env, value)
	}, []string{})

	// toArray() -> Array - Convert range to array
	rangeClass.AddBuiltinMethod("toArray", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		start, _ := utils.AsInt(instance.Fields["_start"])
		end, _ := utils.AsInt(instance.Fields["_end"])
		step, _ := utils.AsInt(instance.Fields["_step"])

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
			e, _ := CreateIntInstance(env, i)
			items = append(items, e)
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
		start, _ := utils.AsInt(instance.Fields["_start"])
		end, _ := utils.AsInt(instance.Fields["_end"])
		step, _ := utils.AsInt(instance.Fields["_step"])

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

		return CreateIntInstance((*Env)(callEnv), size)
	}, []string{})

	// toString() -> String
	rangeClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		start, _ := utils.AsInt(instance.Fields["_start"])
		end, _ := utils.AsInt(instance.Fields["_end"])
		step, _ := utils.AsInt(instance.Fields["_step"])

		if step == 1 {
			return CreateStringInstance(env, fmt.Sprintf("Range(%d..%d)", start, end))
		}
		return CreateStringInstance(env, fmt.Sprintf("Range(%d..%d step %d)", start, end, step))
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

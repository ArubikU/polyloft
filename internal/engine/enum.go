package engine

import (
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Global registry for enums
var enumRegistry = make(map[string]*common.EnumDefinition)

// evalEnumDecl handles enum declaration evaluation
func evalEnumDecl(env *Env, decl *ast.EnumDecl) (any, error) {
	if _, exists := enumRegistry[decl.Name]; exists {
		return nil, ThrowRuntimeError(env, fmt.Sprintf("enum '%s' is already defined", decl.Name))
	}

	definition := &common.EnumDefinition{
		Name:         decl.Name,
		Type:         &ast.Type{Name: decl.Name, IsBuiltin: false, IsEnum: true},
		AccessLevel:  decl.AccessLevel,
		FileName:     env.GetFileName(),
		PackageName:  env.GetPackageName(),
		Methods:      make(map[string][]common.MethodInfo),
		Values:       make(map[string]*common.EnumValueInstance),
		Fields:       make(map[string]common.FieldInfo),
		Constructors: []common.ConstructorInfo{},
		IsSealed:     decl.IsSealed,
		Permits:      nil, // Not resolved here; validated at usage time
	}

	// Collect enum methods
	for _, method := range decl.Methods {
		info := common.MethodInfo{
			Name:        method.Name,
			ReturnType:  method.ReturnType,
			Body:        method.Body,
			Modifiers:   method.Modifiers,
			IsAbstract:  method.IsAbstract,
			IsStatic:    contains(method.Modifiers, "static"),
			IsPrivate:   contains(method.Modifiers, "private"),
			BuiltinImpl: nil,
		}
		info.Params = append(info.Params, method.Params...)
		definition.Methods[method.Name] = append(definition.Methods[method.Name], info)
	}

	// Collect enum fields (instance-level)
	for _, field := range decl.Fields {
		if contains(field.Modifiers, "static") {
			return nil, ThrowNotImplementedError(env, fmt.Sprintf("static fields in enum '%s'", decl.Name))
		}

		fieldInfo := common.FieldInfo{
			Name:      field.Name,
			Type:      field.Type,
			Modifiers: field.Modifiers,
			IsStatic:  false,
			IsPrivate: contains(field.Modifiers, "private"),
		}

		if field.InitValue != nil {
			val, err := evalExpr(env, field.InitValue)
			if err != nil {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("error evaluating field '%s' in enum '%s': %v", field.Name, decl.Name, err))
			}
			fieldInfo.InitValue = val
		}

		definition.Fields[field.Name] = fieldInfo
	}

	// Enum constructor (optional)
	if decl.Constructor != nil {
		constructorInfo := common.ConstructorInfo{
			Body: decl.Constructor.Body,
		}
		constructorInfo.Params = append(constructorInfo.Params, decl.Constructor.Params...)
		definition.Constructors = append(definition.Constructors, constructorInfo)
	}

	enumObject := map[string]any{
		"name": decl.Name,
	}
	valuesSlice := make([]any, 0, len(decl.Values))
	nameSlice := make([]any, 0, len(decl.Values))

	for idx, value := range decl.Values {
		evaluatedArgs := make([]any, len(value.Args))
		for i, argExpr := range value.Args {
			v, err := evalExpr(env, argExpr)
			if err != nil {
				return nil, err
			}
			evaluatedArgs[i] = v
		}

		if len(definition.Constructors) == 0 && len(evaluatedArgs) > 0 {
			return nil, ThrowTypeError(env, "enum value without arguments", fmt.Sprintf("%s.%s with %d arguments", decl.Name, value.Name, len(evaluatedArgs)))
		}

		instance := &common.EnumValueInstance{
			Definition: definition,
			Name:       value.Name,
			Ordinal:    idx,
			Fields: map[string]any{
				"name":    value.Name,
				"ordinal": idx,
			},
			Methods: make(map[string]common.Func),
		}

		for fieldName, fieldInfo := range definition.Fields {
			instance.Fields[fieldName] = fieldInfo.InitValue
		}

		if len(evaluatedArgs) > 0 {
			instance.Fields["args"] = evaluatedArgs
			for argIndex, argVal := range evaluatedArgs {
				instance.Fields[fmt.Sprintf("arg%d", argIndex)] = argVal
			}
		}

		bindEnumInstanceMethods(instance)

		if len(definition.Constructors) > 0 {
			// Select appropriate constructor based on argument count
			constructor := common.SelectConstructorOverload(definition.Constructors, len(evaluatedArgs))
			if constructor == nil {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("no constructor found for enum %s with %d arguments", decl.Name, len(evaluatedArgs)))
			}

			ctorEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
			ctorEnv.Set("this", instance)

			if err := bindParametersWithVariadic(ctorEnv, constructor.Params, evaluatedArgs); err != nil {
				return nil, err
			}

			for _, stmt := range constructor.Body {
				if stmt == nil {
					continue
				}
				_, ret, err := evalStmt(ctorEnv, stmt)
				if err != nil {
					return nil, err
				}
				if ret {
					break
				}
			}
		}

		definition.Values[value.Name] = instance
		enumObject[value.Name] = instance
		valuesSlice = append(valuesSlice, instance)
		nameSlice = append(nameSlice, value.Name)
	}

	bindEnumStaticMethods(enumObject, definition, valuesSlice, nameSlice)

	enumRegistry[definition.Name] = definition

	// Wrap enum object with EnumConstructor for proper type identification
	enumConstructor := &common.EnumConstructor{
		Definition: definition,
		EnumObject: enumObject,
	}

	env.Set(definition.Name, enumConstructor)

	return enumConstructor, nil
}

// bindEnumInstanceMethods attaches instance methods to an enum value
func bindEnumInstanceMethods(instance *common.EnumValueInstance) {
	def := instance.Definition
	if def == nil {
		return
	}

	for name, methodOverloads := range def.Methods {
		overloads := methodOverloads // Copy for closure
		instance.Methods[name] = Func(func(callEnv *Env, args []any) (any, error) {
			// Select appropriate method based on argument count
			method := common.SelectMethodOverload(overloads, len(args))
			if method == nil {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("no overload found for %s.%s with %d arguments", def.Name, name, len(args)))
			}

			if method.IsStatic || method.IsAbstract {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("cannot call static or abstract method %s via instance", name))
			}

			methodEnv := callEnv.Child()
			methodEnv.Set("this", instance)

			if err := bindParametersWithVariadic(methodEnv, method.Params, args); err != nil {
				return nil, err
			}

			var last any
			for _, stmt := range method.Body {
				val, ret, err := evalStmt(methodEnv, stmt)
				if err != nil {
					return nil, err
				}
				if ret {
					return val, nil
				}
				last = val
			}
			return last, nil
		})
	}

	if _, ok := instance.Methods["toString"]; !ok {
		instance.Methods["toString"] = Func(func(_ *Env, _ []any) (any, error) {
			if instance.Definition != nil {
				return fmt.Sprintf("%s.%s", instance.Definition.Name, instance.Name), nil
			}
			return instance.Name, nil
		})
	}
}

// bindEnumStaticMethods attaches static methods to the enum object
func bindEnumStaticMethods(enumObject map[string]any, def *common.EnumDefinition, valuesSlice []any, nameSlice []any) {
	for name, methodOverloads := range def.Methods {
		overloads := methodOverloads // Copy for closure
		enumObject[name] = Func(func(callEnv *Env, args []any) (any, error) {
			// Select appropriate method based on argument count
			method := common.SelectMethodOverload(overloads, len(args))
			if method == nil {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("no overload found for static %s.%s with %d arguments", def.Name, name, len(args)))
			}

			if !method.IsStatic || method.IsAbstract {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("method %s.%s is not a static method", def.Name, name))
			}

			methodEnv := callEnv.Child()

			if err := bindParametersWithVariadic(methodEnv, method.Params, args); err != nil {
				return nil, err
			}

			var last any
			for _, stmt := range method.Body {
				val, ret, err := evalStmt(methodEnv, stmt)
				if err != nil {
					return nil, err
				}
				if ret {
					return val, nil
				}
				last = val
			}
			return last, nil
		})
	}

	if _, ok := enumObject["valueOf"]; !ok {
		enumObject["valueOf"] = Func(func(e *Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ThrowArityError(e, 1, len(args))
			}
			// Use utils.ToString to handle both string and ClassInstance
			name := utils.ToString(args[0])
			if value, exists := def.Values[name]; exists {
				return value, nil
			}
			return nil, ThrowValueError(e, fmt.Sprintf("enum value '%s' not found in enum '%s'", name, def.Name))
		})
	}

	// attach helper static methods
	if _, ok := enumObject["values"]; !ok {
		enumObject["values"] = Func(func(env *Env, _ []any) (any, error) {
			// Wrap native slice in builtin Array class
			return CreateArrayInstance(env, valuesSlice)
		})
	}
	if _, ok := enumObject["size"]; !ok {
		enumObject["size"] = Func(func(_ *Env, _ []any) (any, error) {
			return len(valuesSlice), nil
		})
	}
	if _, ok := enumObject["names"]; !ok {
		enumObject["names"] = Func(func(env *Env, _ []any) (any, error) {
			// Wrap native slice in builtin Array class
			return CreateArrayInstance(env, nameSlice)
		})
	}
}

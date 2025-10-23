package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// Global registry for records
var recordRegistry = make(map[string]*common.RecordDefinition)

// evalRecordDecl handles record declaration evaluation
func evalRecordDecl(env *Env, decl *ast.RecordDecl) (any, error) {
	if _, exists := recordRegistry[decl.Name]; exists {
		return nil, ThrowRuntimeError(env, fmt.Sprintf("record '%s' is already defined", decl.Name))
	}

	definition := &common.RecordDefinition{
		Name:        decl.Name,
		Type:        &ast.Type{Name: decl.Name, IsBuiltin: false, IsRecord: true},
		AccessLevel: decl.AccessLevel,
		FileName:    env.GetFileName(),
		PackageName: env.GetPackageName(),
		Components:  decl.Components,
		Methods:     make(map[string][]common.MethodInfo),
	}

	for _, method := range decl.Methods {
		if contains(method.Modifiers, "static") {
			return nil, ThrowNotImplementedError(env, fmt.Sprintf("static methods in record '%s'", decl.Name))
		}
		info := common.MethodInfo{
			Name:        method.Name,
			ReturnType:  method.ReturnType,
			Body:        method.Body,
			Modifiers:   method.Modifiers,
			IsAbstract:  method.IsAbstract,
			IsStatic:    false,
			IsPrivate:   contains(method.Modifiers, "private"),
			BuiltinImpl: nil,
		}
		info.Params = append(info.Params, method.Params...)
		definition.Methods[method.Name] = append(definition.Methods[method.Name], info)
	}

	recordRegistry[definition.Name] = definition

	constructor := Func(func(callEnv *Env, args []any) (any, error) {
		if len(args) != len(definition.Components) {
			return nil, ThrowArityError(callEnv, len(definition.Components), len(args))
		}

		instance := &common.RecordInstance{
			Definition: definition,
			Values:     make(map[string]any, len(definition.Components)),
			Methods:    make(map[string]common.Func),
		}

		for idx, component := range definition.Components {
			value := args[idx]
			componentTypeName := ast.GetTypeNameString(component.Type)
			if componentTypeName != "" {
				if err := ValidateArgumentType(value, componentTypeName); err != nil {
					return nil, err
				}
			}
			instance.Values[component.Name] = value
		}

		bindRecordInstanceMethods(instance)

		return instance, nil
	})

	env.Set(definition.Name, constructor)
	return constructor, nil
}

// bindRecordInstanceMethods attaches instance methods to a record instance
func bindRecordInstanceMethods(instance *common.RecordInstance) {
	def := instance.Definition
	if def == nil {
		return
	}

	for name, methodOverloads := range def.Methods {
		overloads := methodOverloads // Copy for closure
		instance.Methods[name] = Func(func(callEnv *Env, args []any) (any, error) {
			// Select appropriate method based on argument count
			method := utils.SelectMethodOverload(overloads, len(args))
			if method == nil {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("no overload found for %s.%s with %d arguments", def.Name, name, len(args)))
			}

			if method.IsAbstract {
				return nil, ThrowRuntimeError(callEnv, fmt.Sprintf("cannot call abstract method %s", name))
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
			if def == nil {
				return "record", nil
			}
			parts := make([]string, 0, len(def.Components))
			for _, component := range def.Components {
				parts = append(parts, fmt.Sprintf("%s=%s", component.Name, utils.ToString(instance.Values[component.Name])))
			}
			return fmt.Sprintf("%s(%s)", def.Name, strings.Join(parts, ", ")), nil
		})
	}
}

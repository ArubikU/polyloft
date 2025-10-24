package engine

import (
	"strings"

	"github.com/ArubikU/polyloft/internal/common"
)

// GenericType represents a generic type parameter (like T, E, K, V in Java)
type GenericType = common.GenericType

// GenericInstance represents an instantiated generic type
type GenericInstance struct {
	BaseType string
	TypeArgs []string
	FullName string // e.g., "List<Int>"
}

// InstallGenerics sets up generic type support in the environment
// NOTE: List, Set, Map, and Deque are now installed via InstallGenericCollections using builder pattern
// This function now only installs Lambda (legacy support)
func InstallGenerics(env *common.Env) {

	// Lambda generic type
	env.Set("Lambda", common.Func(func(env *common.Env, args []any) (any, error) {
		if len(args) < 1 {
			return nil, ThrowArityError((*Env)(env), 1, len(args))
		}

		// Last arg is return type, others are param types
		paramTypes := make([]string, len(args)-1)
		for i := 0; i < len(args)-1; i++ {
			if t, ok := args[i].(string); ok {
				paramTypes[i] = t
			} else {
				return nil, ThrowTypeError((*Env)(env), "string", args[i])
			}
		}

		var returnType string
		if t, ok := args[len(args)-1].(string); ok {
			returnType = t
		} else {
			return nil, ThrowTypeError((*Env)(env), "string", args[len(args)-1])
		}

		instance := &LambdaTypeInstance{
			paramTypes: paramTypes,
			returnType: returnType,
		}
		return wrapLambdaType(instance), nil
	}))
}

// LambdaTypeInstance represents a Lambda<P1, P2, ..., R>
type LambdaTypeInstance struct {
	paramTypes []string
	returnType string
}

func wrapLambdaType(lt *LambdaTypeInstance) any {
	typeStr := "Lambda<" + strings.Join(append(lt.paramTypes, lt.returnType), ", ") + ">"
	return map[string]any(map[string]any{
		"toString": common.Func(func(env *common.Env, args []any) (any, error) {
			return typeStr, nil
		}),
		"getParamTypes": common.Func(func(env *common.Env, args []any) (any, error) {
			result := make([]any, len(lt.paramTypes))
			for i, t := range lt.paramTypes {
				result[i] = t
			}
			return result, nil
		}),
		"getReturnType": common.Func(func(env *common.Env, args []any) (any, error) {
			return lt.returnType, nil
		}),
	})
}

// isTypeName checks if a string looks like a type name (starts with uppercase)
// This helps distinguish between type parameters and actual values
func isTypeName(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Common type names start with uppercase
	firstChar := s[0]
	return firstChar >= 'A' && firstChar <= 'Z'
}

// formatWildcard formats a wildcard type for display
func formatWildcard(kind string, bound string) string {
	switch kind {
	case "unbounded":
		return "?"
	case "extends":
		return "? extends " + bound
	case "super":
		return "? super " + bound
	default:
		return "?"
	}
}

// parseWildcardInfo safely extracts kind, bound, and variance from a wildcard info map
func parseWildcardInfo(wildcardInfo map[string]any) (kind, bound, variance string) {
	kind, kindOk := wildcardInfo["kind"].(string)
	bound, boundOk := wildcardInfo["bound"].(string)
	variance, _ = wildcardInfo["variance"].(string)
	if !kindOk {
		kind = "unbounded"
	}
	if !boundOk {
		bound = ""
	}
	return kind, bound, variance
}

// extractTypeArg extracts a type argument from an argument value
// It handles both string type names, variance-annotated types, and wildcard info maps
func extractTypeArg(arg any) string {
	if typeInfo, ok := arg.(map[string]any); ok {
		if isWildcard, ok := typeInfo["isWildcard"].(bool); ok {
			if isWildcard {
				// Wildcard type
				kind, bound, variance := parseWildcardInfo(typeInfo)
				result := formatWildcard(kind, bound)
				if variance != "" {
					result = variance + " " + result
				}
				return result
			} else {
				// Variance-annotated type
				name, _ := typeInfo["name"].(string)
				variance, _ := typeInfo["variance"].(string)
				if variance != "" && name != "" {
					return variance + " " + name
				}
				return name
			}
		}
	}
	if typeStr, ok := arg.(string); ok {
		return typeStr
	}
	return "Any"
}

package engine

import (
	"reflect"

	"github.com/ArubikU/polyloft/internal/ast"
)

// evalSwitchStmt evaluates a switch statement
// Supports:
// - Value matching: switch x case 1, 2: ... case 3: ...
// - Type matching: switch Sys.type(x) case (val: Int): ...
// - Enum matching: switch enumVar case Color.RED: ...
// - Default case: default: ...
func evalSwitchStmt(env *Env, stmt *ast.SwitchStmt) (val any, returned bool, err error) {
	// Evaluate the switch expression
	var switchValue any
	if stmt.Expr != nil {
		switchValue, err = evalExpr(env, stmt.Expr)
		if err != nil {
			return nil, false, err
		}
	}

	// Try each case in order
	for _, c := range stmt.Cases {
		matched := false

		// Type matching case: case (varName: TypeName):
		if c.TypeName != "" {
			// Get the type of the switch value
			typeName := getTypeName(switchValue)
			
			// Check if types match (case-insensitive comparison for built-in types)
			if matchesTypeNameSwitch(typeName, c.TypeName) {
				matched = true
				
				// If a variable name is provided, bind the value to that variable
				if c.VarName != "" {
					env.Set(c.VarName, switchValue)
				}
			}
		} else {
			// Value matching case: case value1, value2:
			for _, caseValue := range c.Values {
				caseVal, err := evalExpr(env, caseValue)
				if err != nil {
					return nil, false, err
				}

				// Check if the values match
				if valuesEqual(switchValue, caseVal) {
					matched = true
					break
				}
			}
		}

		// If this case matched, execute its body
		if matched {
			_, _, ret, val, err := runBlock(env, c.Body)
			if err != nil {
				return nil, false, err
			}
			// For switch, break exits the switch (but we don't propagate it up)
			// Continue doesn't make sense in switch context
			// Return should propagate up
			if ret {
				return val, true, nil
			}
			// After executing a matching case, exit the switch (no fall-through)
			return nil, false, nil
		}
	}

	// If no case matched and there's a default case, execute it
	if len(stmt.Default) > 0 {
		_, _, ret, val, err := runBlock(env, stmt.Default)
		if err != nil {
			return nil, false, err
		}
		if ret {
			return val, true, nil
		}
	}

	return nil, false, nil
}

// getTypeName returns the type name of a value
func getTypeName(val any) string {
	if val == nil {
		return "nil"
	}

	// Check for Polyloft types
	switch v := val.(type) {
	case bool:
		return "bool"
	case int, int8, int16, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "map"
	case *ClassInstance:
		return v.ClassName
	case *EnumValue:
		return v.EnumName
	default:
		// Use reflection as fallback
		return reflect.TypeOf(val).String()
	}
}

// matchesTypeNameSwitch checks if a type name matches the expected type name
// Handles case-insensitive matching for built-in types and their aliases
func matchesTypeNameSwitch(actual, expected string) bool {
	// Direct match
	if actual == expected {
		return true
	}

	// Normalize and compare for built-in types
	actualNorm := normalizeTypeNameSwitch(actual)
	expectedNorm := normalizeTypeNameSwitch(expected)
	
	return actualNorm == expectedNorm
}

// normalizeTypeNameSwitch normalizes type names for comparison
func normalizeTypeNameSwitch(typeName string) string {
	switch typeName {
	case "bool", "Bool", "boolean", "Boolean":
		return "bool"
	case "int", "Int", "integer", "Integer":
		return "int"
	case "float", "Float":
		return "float"
	case "number", "Number":
		return "number"
	case "string", "String":
		return "string"
	case "array", "Array":
		return "array"
	case "map", "Map", "object", "Object":
		return "map"
	case "nil", "Nil", "null", "Null":
		return "nil"
	default:
		return typeName
	}
}

// valuesEqual checks if two values are equal
func valuesEqual(a, b any) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Use reflection for deep comparison
	return reflect.DeepEqual(a, b)
}

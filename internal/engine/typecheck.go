package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// GetTypeName returns the type name for any value
// This is the implementation for Sys.type() function
// Returns lowercase names for primitives, and formatted names for classes/enums
func GetTypeName(val any) string {
	switch v := val.(type) {
	case *common.ClassConstructor:
		// Format as "Class {Name}@{Package}"
		pkg := v.Definition.PackageName
		if pkg == "" {
			pkg = "default"
		}
		return fmt.Sprintf("Class %s@%s", v.Definition.Name, pkg)
	case *common.EnumConstructor:
		// Format as "Enum {Name}@{Package}"
		pkg := v.Definition.PackageName
		if pkg == "" {
			pkg = "default"
		}
		return fmt.Sprintf("Enum %s@%s", v.Definition.Name, pkg)
	case *common.ClassInstance:
		// Return lowercase for primitive wrapper classes to match type system
		switch v.ClassName {
		case "String":
			return "string"
		case "Int":
			return "int"
		case "Float":
			return "float"
		case "Bool":
			return "bool"
		default:
			var typeArgs []string
			typeArgsGeneric := v.GenericTypes
			if len(v.ParentClass.TypeParams) > 0 {
				if len(v.GenericTypes) == 0 {
					if v.ParentClass.ImplementsInterface(common.BuiltinInterfaceCollection.InterfaceDef) {
						methods := v.ParentClass.GetMethods("asArray")
						method := common.SelectMethodOverload(methods, 0)

						if method != nil && method.ReturnType != nil {
							if method.ReturnType == common.BuiltinTypeArray.GetTypeDefinition(nil) && len(method.ReturnType.TypeParams) == 1 {
								arrayStr, ok := CallInstanceMethod(v, *method, nil, []any{})
								if ok != nil {
									inferredType := common.InferCollectionType(arrayStr.([]*common.ClassInstance))
									if inferredType != nil && len(inferredType.Bounds) > 0 {
										typeArgsGeneric = append(typeArgsGeneric, *inferredType)
									}
								}

							}
						}
					}
					// TODO: Implement InferMapType for Map instances
					// if v.ParentClass == common.BuiltinTypeMap.ClassDef {
					//	 common.InferMapType(...)
					// }
				}

				if len(v.GenericTypes) > 0 {
					for _, gt := range v.GenericTypes {
						gTypeArg := ""
						for _, bound := range gt.Bounds {
							// Use the first bound's name for type argument
							typeArg := ""
							if bound.Variance != "" {
								typeArg += bound.Variance + " "
							}
							typeArg += bound.Name.Name
							if bound.Extends != nil {
								typeArg += " extends " + bound.Extends.Name
							}
							if bound.Implements != nil {
								typeArg += " implements " + bound.Implements.Name
							}
							if bound.IsVariadic {
								typeArg += "..."
							}
							if gTypeArg == "" {
								gTypeArg = typeArg
							} else {
								gTypeArg += " | " + typeArg
							}
						}
						typeArgs = append(typeArgs, gTypeArg)
					}
					return fmt.Sprintf("%s<%s>", v.ClassName, strings.Join(typeArgs, ", "))
				}
			}
			return v.ClassName
		}

	case *common.EnumValueInstance:
		if v.Definition != nil {
			return v.Definition.Name
		}
		return "enum"
	case *common.RecordInstance:
		if v.Definition != nil {
			return v.Definition.Name
		}
		return "record"
	case int, int32, int64, float32, float64, string, bool, []any, map[string]any:
		// Native Go types should be wrapped in Generic builtin
		return "Generic"
	case nil:
		return ast.NIL.Name
	case common.Func:
		return "function"
	case *common.FunctionDefinition:
		return common.FormatFunctionType(v.Params, v.ReturnType)
	case *common.LambdaDefinition:
		return common.FormatFunctionType(v.Params, v.ReturnType)
	default:
		fmt.Println("Unknown type for GetTypeName:", v)
		return v.(fmt.Stringer).String()
	}
}

// matchesTypeName checks if a type name matches the expected name, considering aliases
// matchesTypeName checks if a base type name matches a given type name
// Handles aliases like: Integer=Int, Boolean=Bool, etc.
// Both parameters are normalized to lowercase for comparison
func matchesTypeName(baseName, typeName string) bool {
	// Normalize comparison
	typeName = strings.ToLower(strings.TrimSpace(typeName))
	baseName = strings.ToLower(baseName)

	if baseName == typeName {
		return true
	}

	// Check common type aliases
	// This allows Integer to match Int, Boolean to match Bool, etc.
	aliases := map[string][]string{
		"int":      {"integer", "int32", "int64"},
		"float":    {"double", "float32", "float64"},
		"string":   {"str"},
		"bool":     {"boolean"},
		"array":    {"list"},
		"map":      {"object", "dict", "dictionary"},
		"function": {"func", "fn"},
	}

	if aliasList, ok := aliases[baseName]; ok {
		for _, alias := range aliasList {
			if alias == typeName {
				return true
			}
		}
	}

	return false
}

// IsInstanceOf checks if a value is an instance of the given type name
// This is the core type checking function used throughout the system
// Supports: basic types, generic types (Array<Int>), union types (Int | String), wildcards (? extends Number)
func IsInstanceOf(value any, typeName string) bool {
	// Resolve type aliases
	// Check if typeName is a type alias and resolve to base type
	if resolvedType := resolveTypeAlias(typeName); resolvedType != typeName {
		return IsInstanceOf(value, resolvedType)
	}
	
	// Parse the type name to check for generic parameters
	if strings.Contains(typeName, "<") && strings.Contains(typeName, ">") {
		return isInstanceOfGenericType(value, typeName)
	}

	// Handle union types in typeName (e.g., "String | Int")
	if strings.Contains(typeName, "|") {
		return isInstanceOfUnionType(value, typeName)
	}

	switch v := value.(type) {
	case *common.ClassInstance:
		return isClassInstanceOf(v, typeName)
	case *common.EnumValueInstance:
		if v.Definition == nil {
			return false
		}
		return v.Definition.Name == typeName
	case *common.RecordInstance:
		if v.Definition == nil {
			return false
		}
		return v.Definition.Name == typeName
	case int, int32, int64:
		return matchesTypeName("int", typeName) || matchesTypeName("number", typeName)
	case float32, float64:
		return matchesTypeName("float", typeName) || matchesTypeName("number", typeName)
	case string:
		return matchesTypeName("string", typeName)
	case bool:
		return matchesTypeName("bool", typeName)
	case []any:
		return matchesTypeName("array", typeName)
	case map[string]any:
		return matchesTypeName("map", typeName)
	case nil:
		return ast.NIL.MatchesType(typeName)
	case common.Func:
		return matchesTypeName("function", typeName)
	case *common.FunctionDefinition:
		return matchesTypeName("function", typeName)
	case *common.LambdaDefinition:
		return matchesTypeName("function", typeName)
	default:
		return false
	}
}

// isInstanceOfUnionType checks if a value matches any type in a union type
func isInstanceOfUnionType(value any, typeName string) bool {
	// Split by | and trim spaces
	types := strings.Split(typeName, "|")
	for _, t := range types {
		t = strings.TrimSpace(t)
		if IsInstanceOf(value, t) {
			return true
		}
	}
	return false
}

// isInstanceOfGenericType checks if a value is an instance of a generic type
// Examples: Array<Int>, List<String>, Map<String, Int>, Array<? extends Number>
func isInstanceOfGenericType(value any, typeName string) bool {
	// Parse the generic type
	openBracket := strings.Index(typeName, "<")
	if openBracket == -1 {
		return false
	}

	baseName := strings.TrimSpace(typeName[:openBracket])
	closeBracket := strings.LastIndex(typeName, ">")
	if closeBracket == -1 || closeBracket <= openBracket {
		return false
	}

	typeParamsStr := typeName[openBracket+1 : closeBracket]
	typeParams := parseTypeParameters(typeParamsStr)

	// Check base type first
	switch v := value.(type) {
	case []any:
		// Array<T> or Array<? extends T>
		if !matchesTypeName("array", baseName) {
			return false
		}
		// If no type parameters specified, just check if it's an array
		if len(typeParams) == 0 {
			return true
		}
		// Check if all elements match the type parameter
		return allElementsMatchType(v, typeParams[0])

	case *common.ClassInstance:
		// Check for List, Set, Map, Deque, and Array
		if !isClassInstanceOf(v, baseName) {
			return false
		}

		// Get the stored type arguments from the instance
		if typeArgs, ok := v.Fields["__type_args__"].([]string); ok {
			// If stored type is Any, fall through to check elements directly
			if len(typeArgs) > 0 && normalizeTypeName(typeArgs[0]) != "any" {
				// Check if the stored type arguments are compatible with the requested type
				return areTypeArgsCompatible(typeArgs, typeParams)
			}
		}

		// If no type args stored or type args is Any, check the elements directly
		switch baseName {
		case "Array":
			// Arrays store elements in _items field as []any
			if items, ok := v.Fields["_items"].([]any); ok {
				return allElementsMatchType(items, typeParams[0])
			}
			// Also check for pointer to slice (used by some collections)
			if itemsPtr, ok := v.Fields["_items"].(*[]any); ok {
				return allElementsMatchType(*itemsPtr, typeParams[0])
			}
			return true // Empty array or array without elements field
		case "List":
			if itemsPtr, ok := v.Fields["_items"].(*[]any); ok {
				return allElementsMatchType(*itemsPtr, typeParams[0])
			}
		case "Set":
			if itemsMap, ok := v.Fields["_items"].(map[uint64]any); ok {
				items := make([]any, 0, len(itemsMap))
				for _, item := range itemsMap {
					items = append(items, item)
				}
				return allElementsMatchType(items, typeParams[0])
			}
		case "Map":
			if len(typeParams) >= 2 {
				// Map<K, V> - would need to check keys and values
				// For now, just check if it's a Map
				return true
			}
		case "Deque":
			if itemsPtr, ok := v.Fields["_items"].(*[]any); ok {
				return allElementsMatchType(*itemsPtr, typeParams[0])
			}
		}
		return true
	case *common.FunctionDefinition:
		// Function<T1, T2, ..., TRet>
		if !matchesTypeName("function", baseName) {
			return false
		}
		//check last parameter first
		if len(typeParams) == 0 {
			return true
		}
		//get last of typeParams as return type
		returnType := typeParams[len(typeParams)-1]
		if !areTypeArgsCompatible([]string{v.ReturnType.Name}, []string{returnType}) {
			return false
		}
		paramsTypes := typeParams[:len(typeParams)-1]
		funcParamTypes := extractParameterTypeNames(v.Params)
		if !areTypeArgsCompatible(funcParamTypes, paramsTypes) {
			return false
		}
		return true
	case *common.LambdaDefinition:
		// Lambda<T1, T2, ..., TRet>
		if !matchesTypeName("function", baseName) {
			return false
		}
		//check last parameter first
		if len(typeParams) == 0 {
			return true
		}
		//get last of typeParams as return type
		returnType := typeParams[len(typeParams)-1]
		if !areTypeArgsCompatible([]string{v.ReturnType.Name}, []string{returnType}) {
			return false
		}
		paramsTypes := typeParams[:len(typeParams)-1]
		lambdaParamTypes := extractParameterTypeNames(v.Params)
		if !areTypeArgsCompatible(lambdaParamTypes, paramsTypes) {
			return false
		}
		return true
	default:
		fmt.Println("Unsupported type for generic instanceof:", GetTypeName(value))
		return false
	}
}

// extractParameterTypeNames extracts type names from a slice of ast.Parameter
func extractParameterTypeNames(params []ast.Parameter) []string {
	typeNames := make([]string, len(params))
	for i, param := range params {
		if param.Type != nil {
			typeNames[i] = param.Type.Name
		} else {
			typeNames[i] = "Any"
		}
	}
	return typeNames
}

// parseTypeParameters parses comma-separated type parameters, handling nested generics
func parseTypeParameters(paramsStr string) []string {
	if paramsStr == "" {
		return nil
	}

	var params []string
	var currentParam strings.Builder
	depth := 0

	for _, ch := range paramsStr {
		switch ch {
		case '<':
			depth++
			currentParam.WriteRune(ch)
		case '>':
			depth--
			currentParam.WriteRune(ch)
		case ',':
			if depth == 0 {
				paramStr := strings.TrimSpace(currentParam.String())
				if paramStr != "" {
					params = append(params, paramStr)
				}
				currentParam.Reset()
			} else {
				currentParam.WriteRune(ch)
			}
		default:
			currentParam.WriteRune(ch)
		}
	}

	// Add the last parameter
	paramStr := strings.TrimSpace(currentParam.String())
	if paramStr != "" {
		params = append(params, paramStr)
	}

	return params
}

// allElementsMatchType checks if all elements in an array match the given type
func allElementsMatchType(elements []any, typeName string) bool {
	if len(elements) == 0 {
		// Empty array matches any type
		return true
	}

	// Handle wildcard types like "? extends Number" or "? super Int"
	if strings.HasPrefix(typeName, "?") {
		return matchesWildcardType(elements, typeName)
	}

	// Handle union types like "String | Int"
	if strings.Contains(typeName, "|") {
		for _, elem := range elements {
			if !isInstanceOfUnionType(elem, typeName) {
				return false
			}
		}
		return true
	}

	// Check if all elements match the type
	for _, elem := range elements {
		if !IsInstanceOf(elem, typeName) {
			return false
		}
	}
	return true
}

// matchesWildcardType checks if elements match a wildcard type like "? extends Number"
func matchesWildcardType(elements []any, wildcardType string) bool {
	wildcardType = strings.TrimSpace(wildcardType)

	// Parse wildcard: "?", "? extends T", "? super T"
	if wildcardType == "?" || wildcardType == "? extends Any" {
		// Unbounded wildcard - matches everything
		return true
	}

	if strings.HasPrefix(wildcardType, "? extends ") {
		boundType := strings.TrimSpace(wildcardType[len("? extends "):])
		// All elements must be instances of or subtypes of the bound type
		for _, elem := range elements {
			if !isSubtypeOf(elem, boundType) {
				return false
			}
		}
		return true
	}

	if strings.HasPrefix(wildcardType, "? super ") {
		boundType := strings.TrimSpace(wildcardType[len("? super "):])
		// All elements must be supertypes of the bound type
		// This is harder to check at runtime, so we just check if it's assignable
		for _, elem := range elements {
			if !isSupertypeOf(elem, boundType) {
				return false
			}
		}
		return true
	}

	return false
}

// isSubtypeOf checks if a value is a subtype of the given type
func isSubtypeOf(value any, typeName string) bool {
	// Direct match
	if IsInstanceOf(value, typeName) {
		return true
	}

	// Number hierarchy: Int and Float are subtypes of Number
	if typeName == "Number" || typeName == "number" {
		switch value.(type) {
		case int, int32, int64, float32, float64:
			return true
		}
	}

	// Integer is subtype of Int (but NOT Float)
	if typeName == "Int" || typeName == "int" || typeName == "Integer" || typeName == "integer" {
		switch value.(type) {
		case int, int32, int64:
			return true
		}
		return false
	}

	// Float type checking - integers are NOT subtypes of Float
	if typeName == "Float" || typeName == "float" {
		switch value.(type) {
		case float32, float64:
			return true
		}
		return false
	}

	return false
}

// isSupertypeOf checks if a value is a supertype of the given type
func isSupertypeOf(value any, typeName string) bool {
	// This is the inverse of isSubtypeOf
	// If the bound type can be assigned from the value's type
	valueType := GetTypeName(value)

	// Number is a supertype of Int and Float
	if typeName == "Int" || typeName == "int" {
		switch value.(type) {
		case int, int32, int64:
			return true
		case float32, float64:
			// Float is not a supertype of Int
			return false
		}
		// Number is a supertype of Int
		if valueType == "number" || valueType == "Number" {
			return true
		}
	}

	return IsInstanceOf(value, typeName)
}

// areTypeArgsCompatible checks if stored type arguments are compatible with requested types
func areTypeArgsCompatible(storedArgs []string, requestedTypes []string) bool {
	if len(storedArgs) != len(requestedTypes) {
		return false
	}

	for i, requested := range requestedTypes {
		stored := storedArgs[i]

		// Any matches everything (both directions)
		if normalizeTypeName(stored) == "any" || normalizeTypeName(requested) == "any" {
			continue
		}

		// Handle wildcards
		if strings.HasPrefix(requested, "?") {
			if requested == "?" || requested == "? extends Any" {
				continue // Unbounded wildcard matches everything
			}
			if strings.HasPrefix(requested, "? extends ") {
				boundType := strings.TrimSpace(requested[len("? extends "):])
				// Check if stored type is a subtype of bound
				if !isTypeSubtypeOf(stored, boundType) {
					return false
				}
			}
			continue
		}

		// Handle union types in requested
		if strings.Contains(requested, "|") {
			unionTypes := strings.Split(requested, "|")
			matched := false
			for _, unionType := range unionTypes {
				if normalizeTypeName(stored) == normalizeTypeName(unionType) || isTypeSubtypeOf(stored, unionType) {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
			continue
		}

		// Direct comparison or subtype check
		if normalizeTypeName(stored) != normalizeTypeName(requested) && !isTypeSubtypeOf(stored, requested) {
			return false
		}
	}

	return true
}

// isTypeSubtypeOf checks if one type name is a subtype of another
func isTypeSubtypeOf(subtype, supertype string) bool {
	subtype = normalizeTypeName(subtype)
	supertype = normalizeTypeName(supertype)

	if subtype == supertype {
		return true
	}

	// Number hierarchy
	if supertype == "number" {
		return subtype == "int" || subtype == "float" || subtype == "integer"
	}

	return false
}

// normalizeTypeName converts type names to lowercase for comparison
func normalizeTypeName(typeName string) string {
	typeName = strings.TrimSpace(typeName)
	lower := strings.ToLower(typeName)

	// Map common variations
	switch lower {
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	default:
		return lower
	}
}

// isClassInstanceOf checks if a class instance is of a given type.
// This handles inheritance chains, interfaces, and primitive wrapper classes.
func isClassInstanceOf(instance *common.ClassInstance, typeName string) bool {
	// Check direct class name
	if instance.ClassName == typeName {
		return true
	}

	// Handle primitive wrapper classes (String, Int, Float, Bool)
	// They should match both capitalized and lowercase type names
	switch instance.ClassName {
	case "String":
		if typeName == "string" || typeName == "String" {
			return true
		}
	case "Int":
		if typeName == "int" || typeName == "Int" || typeName == "Integer" || typeName == "integer" || typeName == "number" || typeName == "Number" {
			return true
		}
	case "Float":
		if typeName == "float" || typeName == "Float" || typeName == "number" || typeName == "Number" {
			return true
		}
	case "Bool":
		if typeName == "bool" || typeName == "Bool" {
			return true
		}
	}

	// Check aliases if available
	if instance.ParentClass != nil && instance.ParentClass.Aliases != nil {
		for _, alias := range instance.ParentClass.Aliases {
			if alias == typeName {
				return true
			}
		}
	}

	// Check inheritance chain using ParentClass field
	currentClass := instance.ParentClass
	for currentClass != nil {
		if currentClass.Name == typeName {
			return true
		}
		// Check aliases at each level of the inheritance chain
		if currentClass.Aliases != nil {
			for _, alias := range currentClass.Aliases {
				if alias == typeName {
					return true
				}
			}
		}
		currentClass = currentClass.Parent
	}

	// Check implemented interfaces using the parent class definition
	if instance.ParentClass != nil {
		for _, interfaceDef := range instance.ParentClass.Implements {
			if interfaceDef.Name == typeName {
				return true
			}
		}
	}

	return false
}

// IsClassInstanceOfDefinition checks if a class instance matches a class definition.
// This is used when comparing against a ClassConstructor or ClassDefinition directly.
func IsClassInstanceOfDefinition(instance *common.ClassInstance, class *common.ClassDefinition) bool {
	if instance == nil || class == nil {
		return false
	}

	if instance.ParentClass == class {
		return true
	}

	// Check inheritance chain
	currentClass := instance.ParentClass
	for currentClass != nil {
		if currentClass == class {
			return true
		}
		currentClass = currentClass.Parent
	}

	// Check implemented interfaces
	if instance.ParentClass != nil {
		for _, interfaceDef := range instance.ParentClass.Implements {
			if interfaceDef.Name == class.Name {
				return true
			}
		}
	}

	return false
}

// resolveTypeAlias resolves a type name to its base type if it's an alias
// Returns the input typeName unchanged if it's not an alias
func resolveTypeAlias(typeName string) string {
	// Try to find type alias in all packages
	// First check current/default package
	for _, packageAliases := range typeAliasRegistry {
		if alias, exists := packageAliases[typeName]; exists {
			// Recursively resolve in case base type is also an alias
			return resolveTypeAlias(alias.BaseType)
		}
	}
	return typeName
}

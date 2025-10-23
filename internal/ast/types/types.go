// Package types defines built-in types for Polyloft
// This file provides the main registration and access functions for all built-in types
package types

import (
	"github.com/ArubikU/polyloft/internal/ast"
)

// BuiltinTypes contains all built-in type definitions
type BuiltinTypes struct {
	types map[string]*TypeDefinition
}

// TypeDefinition represents a built-in type with its methods
type TypeDefinition struct {
	Name            string
	Type            *ast.Type // Unified type representation
	InstanceMethods map[string]Func
	StaticMethods   map[string]Func
}

// NewBuiltinTypes creates and initializes all built-in types
func NewBuiltinTypes() *BuiltinTypes {
	bt := &BuiltinTypes{
		types: make(map[string]*TypeDefinition),
	}

	// Register all built-in types
	bt.registerStringType()
	bt.registerIntType()
	bt.registerFloatType()
	bt.registerBoolType()
	bt.registerArrayType()
	bt.registerObjectType()

	return bt
}

// registerStringType registers the string built-in type
// Instance methods moved to engine/builtin_string.go
func (bt *BuiltinTypes) registerStringType() {
	bt.types["string"] = &TypeDefinition{
		Name:            "string",
		Type:            &ast.Type{Name: "string", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_string.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}
}

// registerIntType registers the int built-in type
// Instance methods moved to engine/builtin_number.go
func (bt *BuiltinTypes) registerIntType() {
	bt.types["int"] = &TypeDefinition{
		Name:            "int",
		Type:            &ast.Type{Name: "int", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_number.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}
}

// registerFloatType registers the float built-in type
// Instance methods moved to engine/builtin_number.go
func (bt *BuiltinTypes) registerFloatType() {
	bt.types["float"] = &TypeDefinition{
		Name:            "float",
		Type:            &ast.Type{Name: "float", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_number.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}
}

// registerBoolType registers the bool built-in type
// Instance methods moved to engine/builtin_bool.go
func (bt *BuiltinTypes) registerBoolType() {
	bt.types["bool"] = &TypeDefinition{
		Name:            "bool",
		Type:            &ast.Type{Name: "bool", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_bool.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}
}

// registerArrayType registers the array built-in type
// Instance methods moved to engine/builtin_array.go
func (bt *BuiltinTypes) registerArrayType() {
	bt.types["array"] = &TypeDefinition{
		Name:            "array",
		Type:            &ast.Type{Name: "array", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_array.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}
}

// registerObjectType registers the object built-in type
// Instance methods moved to engine/builtin_map.go
func (bt *BuiltinTypes) registerObjectType() {
	bt.types["object"] = &TypeDefinition{
		Name:            "object",
		Type:            &ast.Type{Name: "map", IsBuiltin: true},
		InstanceMethods: make(map[string]Func), // Instance methods now in engine/builtin_map.go
		StaticMethods:   make(map[string]Func), // Static methods can be added here if needed
	}

	// Also register as "map" for compatibility
	bt.types["map"] = bt.types["object"]
}

// GetType returns the type definition for a given type name
func (bt *BuiltinTypes) GetType(typeName string) (*TypeDefinition, bool) {
	typedef, exists := bt.types[typeName]
	return typedef, exists
}

// GetInstanceMethod returns an instance method for a given type and method name
func (bt *BuiltinTypes) GetInstanceMethod(typeName, methodName string) (Func, bool) {
	typedef, exists := bt.types[typeName]
	if !exists {
		return nil, false
	}

	method, exists := typedef.InstanceMethods[methodName]
	return method, exists
}

// GetStaticMethod returns a static method for a given type and method name
func (bt *BuiltinTypes) GetStaticMethod(typeName, methodName string) (Func, bool) {
	typedef, exists := bt.types[typeName]
	if !exists {
		return nil, false
	}

	method, exists := typedef.StaticMethods[methodName]
	return method, exists
}

// GetTypeForValue returns the type name for a given value
func (bt *BuiltinTypes) GetTypeForValue(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case nil:
		return "nil"
	default:
		return "unknown"
	}
}

// HasInstanceMethod checks if a type has a specific instance method
func (bt *BuiltinTypes) HasInstanceMethod(typeName, methodName string) bool {
	typedef, exists := bt.types[typeName]
	if !exists {
		return false
	}

	_, exists = typedef.InstanceMethods[methodName]
	return exists
}

// HasStaticMethod checks if a type has a specific static method
func (bt *BuiltinTypes) HasStaticMethod(typeName, methodName string) bool {
	typedef, exists := bt.types[typeName]
	if !exists {
		return false
	}

	_, exists = typedef.StaticMethods[methodName]
	return exists
}

// ListTypes returns all registered type names
func (bt *BuiltinTypes) ListTypes() []string {
	types := make([]string, 0, len(bt.types))
	for typeName := range bt.types {
		types = append(types, typeName)
	}
	return types
}

// ListInstanceMethods returns all instance method names for a type
func (bt *BuiltinTypes) ListInstanceMethods(typeName string) []string {
	typedef, exists := bt.types[typeName]
	if !exists {
		return nil
	}

	methods := make([]string, 0, len(typedef.InstanceMethods))
	for methodName := range typedef.InstanceMethods {
		methods = append(methods, methodName)
	}
	return methods
}

// ListStaticMethods returns all static method names for a type
func (bt *BuiltinTypes) ListStaticMethods(typeName string) []string {
	typedef, exists := bt.types[typeName]
	if !exists {
		return nil
	}

	methods := make([]string, 0, len(typedef.StaticMethods))
	for methodName := range typedef.StaticMethods {
		methods = append(methods, methodName)
	}
	return methods
}

// Global instance for easy access
var DefaultBuiltinTypes = NewBuiltinTypes()

// Convenience functions that use the default instance

// GetInstanceMethod returns an instance method using the default builtin types
func GetInstanceMethod(typeName, methodName string) (Func, bool) {
	return DefaultBuiltinTypes.GetInstanceMethod(typeName, methodName)
}

// GetStaticMethod returns a static method using the default builtin types
func GetStaticMethod(typeName, methodName string) (Func, bool) {
	return DefaultBuiltinTypes.GetStaticMethod(typeName, methodName)
}

// GetTypeForValue returns the type name for a value using the default builtin types
func GetTypeForValue(value any) string {
	return DefaultBuiltinTypes.GetTypeForValue(value)
}

// HasInstanceMethod checks if a type has an instance method using the default builtin types
func HasInstanceMethod(typeName, methodName string) bool {
	return DefaultBuiltinTypes.HasInstanceMethod(typeName, methodName)
}

// HasStaticMethod checks if a type has a static method using the default builtin types
func HasStaticMethod(typeName, methodName string) bool {
	return DefaultBuiltinTypes.HasStaticMethod(typeName, methodName)
}

// GetTypeByName returns the ast.Type for a given type name, checking all aliases
func GetTypeByName(typeName string) *ast.Type {
	for _, t := range ast.GetBuiltinTypes() {
		if t.MatchesType(typeName) {
			return t
		}
	}
	return nil
}

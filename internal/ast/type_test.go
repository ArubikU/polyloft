package ast

import "testing"

// TestTypeMatching tests the MatchesType method
// Note: Predefined type constants like INT, BOOL, STRING, etc. have been removed.
// Types are now registered dynamically at runtime via ClassBuilder.
func TestTypeMatching(t *testing.T) {
	// Create sample types for testing
	boolType := &Type{Name: "bool", Aliases: []string{"boolean", "Bool", "Boolean"}, IsBuiltin: true}
	intType := &Type{Name: "int", Aliases: []string{"integer", "Int"}, IsBuiltin: true}
	stringType := &Type{Name: "string", Aliases: []string{"String"}, IsBuiltin: true}
	arrayType := &Type{Name: "array", Aliases: []string{"Array"}, IsBuiltin: true}
	mapType := &Type{Name: "map", Aliases: []string{"object", "Object"}, IsBuiltin: true}

	tests := []struct {
		name     string
		typ      *Type
		typeName string
		want     bool
	}{
		{"BOOL matches bool", boolType, "bool", true},
		{"BOOL matches Bool", boolType, "Bool", true},
		{"BOOL matches boolean", boolType, "boolean", true},
		{"BOOL matches Boolean", boolType, "Boolean", true},
		{"BOOL doesn't match int", boolType, "int", false},

		{"INT matches int", intType, "int", true},
		{"INT matches Int", intType, "Int", true},
		{"INT matches integer", intType, "integer", true},
		{"INT doesn't match string", intType, "string", false},

		{"STRING matches string", stringType, "string", true},
		{"STRING matches String", stringType, "String", true},
		{"STRING doesn't match bool", stringType, "bool", false},

		{"ARRAY matches array", arrayType, "array", true},
		{"ARRAY matches Array", arrayType, "Array", true},

		{"MAP matches map", mapType, "map", true},
		{"MAP matches object", mapType, "object", true},
		{"MAP matches Object", mapType, "Object", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typ.MatchesType(tt.typeName)
			if got != tt.want {
				t.Errorf("Type.MatchesType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestResolveTypeName tests that ResolveTypeName works with the remaining predefined types
func TestResolveTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		// Only ANY and NIL are predefined now
		{"Resolve nil", "nil", NIL},
		{"Resolve null", "null", NIL},
		{"Resolve any", "any", ANY},
		{"Resolve Any", "Any", ANY},
		
		// Other types should return nil (not predefined, registered at runtime)
		{"Resolve bool", "bool", nil},
		{"Resolve int", "int", nil},
		{"Resolve string", "string", nil},
		{"Resolve unknown type", "CustomClass", nil},
		{"Resolve empty string", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveTypeName(tt.typeName)
			if got != tt.want {
				t.Errorf("ResolveTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsBuiltinTypeName tests that only ANY and NIL are recognized as builtin at compile time
func TestIsBuiltinTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     bool
	}{
		// Only ANY and NIL are recognized as builtin at compile time
		{"nil is builtin", "nil", true},
		{"null is builtin", "null", true},
		{"any is builtin", "any", true},
		{"Any is builtin", "Any", true},
		
		// Other types are registered at runtime, so not recognized here
		{"bool is not predefined", "bool", false},
		{"int is not predefined", "int", false},
		{"string is not predefined", "string", false},
		{"CustomClass is not builtin", "CustomClass", false},
		{"Random is not builtin", "Random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBuiltinTypeName(tt.typeName)
			if got != tt.want {
				t.Errorf("IsBuiltinTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTypeFields tests the ANY and NIL types which are still predefined
func TestTypeFields(t *testing.T) {
	// Test that ANY and NIL have correct properties
	if !ANY.IsBuiltin {
		t.Error("ANY should be marked as built-in")
	}
	if ANY.IsClass {
		t.Error("ANY should not be marked as class")
	}

	if !NIL.IsBuiltin {
		t.Error("NIL should be marked as built-in")
	}

	// Test creating a custom type
	customType := &Type{
		Name:      "MyClass",
		IsBuiltin: false,
		IsClass:   true,
	}

	if customType.IsBuiltin {
		t.Error("Custom type should not be marked as built-in")
	}
	if !customType.IsClass {
		t.Error("Custom type should be marked as class")
	}
}

func TestGenericTypes(t *testing.T) {
tests := []struct {
name     string
typeName string
expected string
hasParams bool
paramCount int
}{
{"Simple int", "int", "int", false, 0},
{"Array with Int", "Array<Int>", "Array<Int>", true, 1},
{"Map with two params", "Map<String, Int>", "Map<String, Int>", true, 2},
{"Nested generic", "Array<Map<String, Int>>", "Array<Map<String, Int>>", true, 1},
{"List with custom type", "List<Person>", "List<Person>", true, 1},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
typ := TypeFromString(tt.typeName)
if typ == nil {
t.Fatalf("TypeFromString returned nil for %s", tt.typeName)
}

result := GetTypeNameString(typ)
if result != tt.expected {
t.Errorf("GetTypeNameString() = %v, want %v", result, tt.expected)
}

if (len(typ.TypeParams) > 0) != tt.hasParams {
t.Errorf("Type has params = %v, want %v", len(typ.TypeParams) > 0, tt.hasParams)
}

if len(typ.TypeParams) != tt.paramCount {
t.Errorf("Type param count = %v, want %v", len(typ.TypeParams), tt.paramCount)
}
})
}
}

func TestNestedGenericTypes(t *testing.T) {
// Test deeply nested generic types
typ := TypeFromString("Map<String, Array<Map<Int, String>>>")
result := GetTypeNameString(typ)
expected := "Map<String, Array<Map<Int, String>>>"

if result != expected {
t.Errorf("GetTypeNameString() = %v, want %v", result, expected)
}

// Verify structure
if len(typ.TypeParams) != 2 {
t.Fatalf("Expected 2 type params, got %d", len(typ.TypeParams))
}

// Check second parameter (Array<Map<Int, String>>)
arrayParam := typ.TypeParams[1]
if arrayParam.Name != "Array" {
t.Errorf("Second param base type = %v, want 'Array'", arrayParam.Name)
}

if len(arrayParam.TypeParams) != 1 {
t.Fatalf("Array param should have 1 type param, got %d", len(arrayParam.TypeParams))
}

// Check nested Map<Int, String>
nestedMap := arrayParam.TypeParams[0]
if nestedMap.Name != "Map" {
t.Errorf("Nested map base type = %v, want 'Map'", nestedMap.Name)
}

if len(nestedMap.TypeParams) != 2 {
t.Errorf("Nested map should have 2 type params, got %d", len(nestedMap.TypeParams))
}
}

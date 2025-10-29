package common

import (
	"unicode"
	
	"github.com/ArubikU/polyloft/internal/ast"
)

// TypeNames contains standard type name constants
const (
	// Tipos base
	TypeInt      = "int"
	TypeInteger  = "integer"
	TypeFloat    = "float"
	TypeString   = "string"
	TypeBool     = "bool"
	TypeBoolean  = "boolean"
	TypeArray    = "array"
	TypeMap      = "map"
	TypeObject   = "object"
	TypeFunction = "function"
	TypeNil      = "nil"
	TypeNull     = "null"
)

// PrivacyLevel represents visibility/access levels
type PrivacyLevel int

const (
	PrivacyPublic PrivacyLevel = iota
	PrivacyPrivate
	PrivacyProtected
	PrivacyPackage // default package visibility
)

var PrivacyLevelNames = map[PrivacyLevel]string{
	PrivacyPublic:    "public",
	PrivacyPrivate:   "private",
	PrivacyProtected: "protected",
	PrivacyPackage:   "package",
}

// Modifier represents different modifiers that can be applied to declarations
type Modifier int

const (
	ModifierNone Modifier = iota
	ModifierStatic
	ModifierAbstract
	ModifierDefault
	ModifierFinal
)

var ModifierNames = map[Modifier]string{
	ModifierNone:     "none",
	ModifierStatic:   "static",
	ModifierAbstract: "abstract",
	ModifierDefault:  "default",
	ModifierFinal:    "final",
}

// VarKind represents variable declaration kinds
type VarKind string

const (
	VarKindLet   VarKind = "let"
	VarKindVar   VarKind = "var"
	VarKindConst VarKind = "const"
)

// ModifierApplicability defines where each modifier can be used
var ModifierApplicability = map[Modifier][]string{
	ModifierStatic:   {"class", "method", "field", "enum", "record"},
	ModifierAbstract: {"class", "method"},
	ModifierDefault:  {"method"},
	ModifierFinal:    {"class", "method", "field", "variable", "parameter"},
}

// PrivacyApplicability defines where each privacy level can be used
var PrivacyApplicability = map[PrivacyLevel][]string{
	PrivacyPublic:    {"class", "method", "field", "enum", "record"},
	PrivacyPrivate:   {"method", "field", "class"},
	PrivacyProtected: {"method", "field", "class"},
	PrivacyPackage:   {"class", "method", "field", "enum", "record"},
}

// CanUseModifier checks if a modifier can be used in a specific context
func CanUseModifier(modifier Modifier, context string) bool {
	applicableContexts, exists := ModifierApplicability[modifier]
	if !exists {
		return false
	}

	for _, ctx := range applicableContexts {
		if ctx == context {
			return true
		}
	}
	return false
}

// CanUsePrivacyLevel checks if a privacy level can be used in a specific context
func CanUsePrivacyLevel(privacy PrivacyLevel, context string) bool {
	applicableContexts, exists := PrivacyApplicability[privacy]
	if !exists {
		return false
	}

	for _, ctx := range applicableContexts {
		if ctx == context {
			return true
		}
	}
	return false
}

// IsBuiltinType checks if a type name is a built-in type
func IsBuiltinType(typeName string) bool {
	return ast.IsBuiltinTypeName(typeName)
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// NormalizeTypeName normalizes type names to their canonical form
func NormalizeTypeName(typeName string) string {
	switch typeName {
	case TypeInteger, TypeInt:
		return capitalizeFirst(TypeInt)
	case TypeFloat:
		return capitalizeFirst(TypeFloat)
	case TypeString:
		return capitalizeFirst(TypeString)
	case TypeBool, TypeBoolean:
		return capitalizeFirst(TypeBool)
	case TypeArray:
		return capitalizeFirst(TypeArray)
	case TypeMap, TypeObject:
		return capitalizeFirst(TypeMap)
	case TypeNull:
		return capitalizeFirst(TypeNil)
	default:
		return typeName
	}
}

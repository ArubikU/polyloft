package ast

import "strings"

// Package ast defines Polyloft's abstract syntax tree.
// Keep node interfaces small and composable. Favor explicit fields over magic.

// Node is the common interface implemented by all AST nodes.
type Node interface {
	node()
}
type MapEntry struct {
	Key   any
	Value any
}

// Type represents a type in the Polyloft type system
// This is used for both built-in types and user-defined types (classes, interfaces, enums, records)
// It also supports generic/parameterized types like Array<Int>, Map<String, Int>
// And union types like string | int | null
type Type struct {
	Name        string   // Canonical name of the type (e.g., "bool", "int", "String", "Array")
	Aliases     []string // Alternative names for the type (e.g., ["boolean"] for bool)
	TypeParams  []*Type  // Type parameters for generic types (e.g., [Int] for Array<Int>, [String, Int] for Map<String, Int>)
	UnionTypes  []*Type  // Union type members (e.g., [string, int] for string | int)
	GoParallel  bool     // Whether this type supports parallel operations (optional, defaults to false)
	IsBuiltin   bool     // Whether this is a built-in type
	IsClass     bool     // Whether this is a class type
	IsInterface bool     // Whether this is an interface type
	IsEnum      bool     // Whether this is an enum type
	IsRecord    bool     // Whether this is a record type
	IsUnion     bool     // Whether this is a union type
}

// Predefined built-in types
var (
	ANY = &Type{Name: "Any", Aliases: []string{"any"}, GoParallel: false, IsBuiltin: true}
	NIL = &Type{Name: "nil", Aliases: []string{"null", "Nil", "Null"}, GoParallel: false, IsBuiltin: true}
)

// GetBuiltinTypes returns all built-in types
func GetBuiltinTypes() []*Type {
	return []*Type{ANY, NIL}
}

// MatchesType checks if a type name matches this type (checking name and aliases)
func (t *Type) MatchesType(name string) bool {
	if t.Name == name {
		return true
	}
	for _, alias := range t.Aliases {
		if alias == name {
			return true
		}
	}
	return false
}

// ResolveTypeName returns the ast.Type for a given type name string
// It checks built-in types first, then returns nil if not found
func ResolveTypeName(typeName string) *Type {
	if typeName == "" {
		return nil
	}

	// Check all built-in types
	for _, t := range GetBuiltinTypes() {
		if t.MatchesType(typeName) {
			return t
		}
	}

	// Not a built-in type - could be a user-defined type
	// Return nil to indicate this should be resolved at runtime
	return nil
}

// IsBuiltinTypeName checks if a type name refers to a built-in type
func IsBuiltinTypeName(typeName string) bool {
	return ResolveTypeName(typeName) != nil
}

// TypeFromString creates an ast.Type from a string type name
// This is a compatibility helper for code migration
// It supports generic type syntax like "Array<Int>", "Map<String, Int>"
// And union types like "string | int | null"
func TypeFromString(typeName string) *Type {
	if typeName == "" {
		return nil
	}

	// Check if it's a union type (contains |)
	if strings.Contains(typeName, "|") {
		return parseUnionType(typeName)
	}

	// Check if it's a generic type (contains < and >)
	if strings.Contains(typeName, "<") && strings.Contains(typeName, ">") {
		return parseGenericType(typeName)
	}

	// First try to resolve as built-in type
	if builtinType := ResolveTypeName(typeName); builtinType != nil {
		return builtinType
	}

	// Otherwise create a user-defined type placeholder
	// The actual type will be resolved at runtime
	return &Type{
		Name:      typeName,
		IsBuiltin: false,
	}
}

// parseGenericType parses a generic type string like "Array<Int>" or "Map<String, Int>"
func parseGenericType(typeName string) *Type {
	// Find the base type name and type parameters
	openBracket := strings.Index(typeName, "<")
	if openBracket == -1 {
		return TypeFromString(typeName)
	}

	baseName := strings.TrimSpace(typeName[:openBracket])
	closeBracket := strings.LastIndex(typeName, ">")
	if closeBracket == -1 || closeBracket <= openBracket {
		// Invalid syntax, return as simple type
		return TypeFromString(typeName)
	}

	paramsStr := typeName[openBracket+1 : closeBracket]

	// Parse type parameters (handle nested generics)
	typeParams := parseTypeParams(paramsStr)

	// Get the base type
	baseType := ResolveTypeName(baseName)
	if baseType == nil {
		// Create a new type for user-defined generic types
		return &Type{
			Name:       baseName,
			TypeParams: typeParams,
			IsBuiltin:  false,
		}
	}

	// Create a new type instance with type parameters
	return &Type{
		Name:        baseType.Name,
		Aliases:     baseType.Aliases,
		TypeParams:  typeParams,
		GoParallel:  baseType.GoParallel,
		IsBuiltin:   baseType.IsBuiltin,
		IsClass:     baseType.IsClass,
		IsInterface: baseType.IsInterface,
		IsEnum:      baseType.IsEnum,
		IsRecord:    baseType.IsRecord,
	}
}

// parseUnionType parses a union type string like "string | int" or "string | int | null"
func parseUnionType(typeName string) *Type {
	// Split by | but respect generic type brackets
	var unionTypes []*Type
	var currentType strings.Builder
	depth := 0

	for _, ch := range typeName {
		switch ch {
		case '<':
			depth++
			currentType.WriteRune(ch)
		case '>':
			depth--
			currentType.WriteRune(ch)
		case '|':
			if depth == 0 {
				// End of current union member
				typeStr := strings.TrimSpace(currentType.String())
				if typeStr != "" {
					unionTypes = append(unionTypes, TypeFromString(typeStr))
				}
				currentType.Reset()
			} else {
				currentType.WriteRune(ch)
			}
		default:
			currentType.WriteRune(ch)
		}
	}

	// Add the last type
	typeStr := strings.TrimSpace(currentType.String())
	if typeStr != "" {
		unionTypes = append(unionTypes, TypeFromString(typeStr))
	}

	if len(unionTypes) == 0 {
		return nil
	}

	if len(unionTypes) == 1 {
		// Only one type, not a union
		return unionTypes[0]
	}

	// Create union type name by joining all type names
	var typeNames []string
	for _, t := range unionTypes {
		typeNames = append(typeNames, t.Name)
	}
	unionName := strings.Join(typeNames, " | ")

	return &Type{
		Name:       unionName,
		UnionTypes: unionTypes,
		IsUnion:    true,
		IsBuiltin:  false, // Union types are composite types
	}
}

// parseTypeParams parses comma-separated type parameters, handling nested generics
func parseTypeParams(paramsStr string) []*Type {
	if paramsStr == "" {
		return nil
	}

	var params []*Type
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
				// End of current parameter
				paramStr := strings.TrimSpace(currentParam.String())
				if paramStr != "" {
					params = append(params, TypeFromString(paramStr))
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
		params = append(params, TypeFromString(paramStr))
	}

	return params
}

// GenericType creates a generic type with type parameters
// Example: GenericType(arrayType, intType) creates Array<Int>
// The first parameter is the base type, the rest are type parameters
func GenericType(baseType *Type, typeParams ...*Type) *Type {
	if baseType == nil {
		return nil
	}

	// Create a new type instance with the given type parameters
	return &Type{
		Name:        baseType.Name,
		Aliases:     baseType.Aliases,
		TypeParams:  typeParams,
		GoParallel:  baseType.GoParallel,
		IsBuiltin:   baseType.IsBuiltin,
		IsClass:     baseType.IsClass,
		IsInterface: baseType.IsInterface,
		IsEnum:      baseType.IsEnum,
		IsRecord:    baseType.IsRecord,
		IsUnion:     baseType.IsUnion,
	}
}

// GetTypeName returns the string name of a type
// This is a compatibility helper for code migration
func GetTypeNameString(t *Type) string {
	if t == nil {
		return ""
	}

	// If no type parameters, return simple name
	if len(t.TypeParams) == 0 {
		return t.Name
	}

	// Format with type parameters
	var paramNames []string
	for _, param := range t.TypeParams {
		paramNames = append(paramNames, GetTypeNameString(param))
	}

	return t.Name + "<" + strings.Join(paramNames, ", ") + ">"
}

// Position describes a location in a source file.
type Position struct {
	Offset int // 0-based byte offset
	Line   int // 1-based line number
	Col    int // 1-based column number (in runes)
}

// Program is the root node for a source file.
type Program struct {
	Stmts []Stmt
}

func (*Program) node() {}

// Stmt represents a statement.
type Stmt interface {
	Node
	stmt()
}

// Expr represents an expression.
type Expr interface {
	Node
	expr()
}

// Identifier
type Ident struct {
	Name string
}

func (*Ident) node() {}
func (*Ident) expr() {}

// Literals
type NumberLit struct{ Value any } // Can be int or float64
type StringLit struct{ Value string }
type InterpolatedStringLit struct {
	Parts []Expr // alternating string literals and expressions
}
type BoolLit struct{ Value bool }
type NilLit struct{}

func (*NumberLit) node()             {}
func (*NumberLit) expr()             {}
func (*StringLit) node()             {}
func (*StringLit) expr()             {}
func (*InterpolatedStringLit) node() {}
func (*InterpolatedStringLit) expr() {}
func (*BoolLit) node()               {}
func (*BoolLit) expr()               {}
func (*NilLit) node()                {}
func (*NilLit) expr()                {}

// Composite literals
type ArrayLit struct{ Elems []Expr }
type MapPair struct {
	Key   string
	Value Expr
}
type MapLit struct{ Pairs []MapPair }

func (*ArrayLit) node() {}
func (*ArrayLit) expr() {}
func (*MapLit) node()   {}
func (*MapLit) expr()   {}

// Unary and binary
type UnaryExpr struct {
	Op int
	X  Expr
}
type BinaryExpr struct {
	Op       int
	Lhs, Rhs Expr
}

func (*UnaryExpr) node()  {}
func (*UnaryExpr) expr()  {}
func (*BinaryExpr) node() {}
func (*BinaryExpr) expr() {}

// Call
type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (*CallExpr) node() {}
func (*CallExpr) expr() {}

// Generic Call: List<Int>(), Map<String, Int>()
type GenericCallExpr struct {
	Name       string      // Base type name (e.g., "List", "Set", "Map", "Lambda")
	TypeParams []TypeParam // Type parameters (e.g., [TypeParam{Name: "Int"}] for List<Int>)
	Args       []Expr      // Constructor arguments
}

func (*GenericCallExpr) node() {}
func (*GenericCallExpr) expr() {}

// TypeParam represents a type parameter in a generic call, which can be a concrete type or a wildcard
type TypeParam struct {
	IsWildcard   bool     // true if this is a wildcard (?)
	IsVariadic   bool     // true if this is a variadic type parameter (T...)
	Name         string   // Type name for concrete types (e.g., "Int", "String")
	WildcardKind string   // "unbounded", "extends", or "super" for wildcards
	Bounds       []string // Multiple bounds for wildcards (e.g., ["Number", "Comparable"] in "? extends Number & Comparable")
	Variance     string   // "in" (contravariance), "out" (covariance), or "" (invariant)
}

// Super call expression: super(args)
type SuperExpr struct {
	Args []Expr
}

func (*SuperExpr) node() {}
func (*SuperExpr) expr() {}

// Indexing: arr[idx] or map[key]
type IndexExpr struct {
	X     Expr
	Index Expr
}

func (*IndexExpr) node() {}
func (*IndexExpr) expr() {}

// Field access: obj.field
type FieldExpr struct {
	X    Expr
	Name string
}

func (*FieldExpr) node() {}
func (*FieldExpr) expr() {}

// Statements
type LetStmt struct {
	Name      string   // Single variable name (for backward compatibility)
	Names     []string // Multiple variable names for destructuring (e.g., let a, b = [1,2])
	Value     Expr
	Type      *Type    // Type annotation using unified type system
	Modifiers []string // optional modifiers: public/private/protected/static
	Kind      string   // "let", "var", "const", "final"
	Inferred  bool     // true if declared with ':=' (type inference)
}
type AssignStmt struct {
	Target Expr     // left side of assignment (could be identifier or field access)
	Value  Expr     // right side of assignment
	Pos    Position // position of the assignment operator
}
type ReturnStmt struct{ Value Expr }
type ExprStmt struct{ X Expr }
type DefStmt struct {
	Name        string
	Params      []Parameter // updated to support typed and variadic parameters
	Body        []Stmt
	ReturnType  *Type       // Return type using unified type system
	AccessLevel string      // "public", "private", "protected"
	Modifiers   []string    // all modifiers including access level
	TypeParams  []TypeParam // generic type parameters (e.g., [T, K, V])
}
type IfClause struct {
	Cond Expr
	Body []Stmt
}
type IfStmt struct {
	Clauses []IfClause
	Else    []Stmt
}
type ForInStmt struct {
	Name     string   // deprecated: use Names for single or multiple vars
	Names    []string // iteration variable names (supports destructuring)
	Iterable Expr
	Where    Expr // optional where clause for filtering
	Body     []Stmt
}
type LoopStmt struct{ Body []Stmt }
type BreakStmt struct{}
type ContinueStmt struct{}

// Import statement: import path.with.dots { Name, Name2 }
type ImportStmt struct {
	Path  []string // e.g., ["math","vector"]
	Names []string // specific symbols to import; if empty, import as namespace (future)
}

// Try-catch statement: try { ... } catch e: Type { ... } finally { ... }
type CatchClause struct {
	VarName    string // variable name to bind the exception
	Modifier   string // modifier like "final", "const"
	ExceptType string // type of exception to catch (optional)
	Body       []Stmt
}

type TryStmt struct {
	Body    []Stmt
	Catches []CatchClause // can have multiple catch clauses
	Finally []Stmt        // optional finally block
}

// Throw statement: throw expr
type ThrowStmt struct {
	Value Expr
}

// Defer statement: defer expr (usually a function call)
type DeferStmt struct {
	Call Expr
	Pos  Position
}

func (*LetStmt) node()      {}
func (*LetStmt) stmt()      {}
func (*AssignStmt) node()   {}
func (*AssignStmt) stmt()   {}
func (*ReturnStmt) node()   {}
func (*ReturnStmt) stmt()   {}
func (*ExprStmt) node()     {}
func (*ExprStmt) stmt()     {}
func (*DefStmt) node()      {}
func (*DefStmt) stmt()      {}
func (*IfStmt) node()       {}
func (*IfStmt) stmt()       {}
func (*ForInStmt) node()    {}
func (*ForInStmt) stmt()    {}
func (*LoopStmt) node()     {}
func (*LoopStmt) stmt()     {}
func (*BreakStmt) node()    {}
func (*BreakStmt) stmt()    {}
func (*ContinueStmt) node() {}
func (*ContinueStmt) stmt() {}
func (*ImportStmt) node()   {}
func (*ImportStmt) stmt()   {}
func (*TryStmt) node()      {}
func (*TryStmt) stmt()      {}
func (*ThrowStmt) node()    {}
func (*ThrowStmt) stmt()    {}
func (*DeferStmt) node()    {}
func (*DeferStmt) stmt()    {}

// Interface declaration with method signatures
type InterfaceDecl struct {
	Name        string
	Methods     []MethodSignature
	TypeParams  []TypeParam // generic type parameters (e.g., [T, K, V])
	IsSealed    bool        // whether the interface is sealed
	Permits     []string    // names permitted to implement
	Fields      []FieldDecl // static fields
	AccessLevel string      // "public", "private", "protected"
}

func (*InterfaceDecl) node() {}
func (*InterfaceDecl) stmt() {}

// Method signature for interfaces
type MethodSignature struct {
	Name        string
	Params      []Parameter
	ReturnType  *Type  // Return type using unified type system
	HasDefault  bool   // whether this method has a default implementation
	DefaultBody []Stmt // default implementation body
}

// Parameter with optional type annotation
type Parameter struct {
	Name       string
	Type       *Type // Type annotation using unified type system
	IsVariadic bool  // true if this parameter is variadic (args...)
}

// Class declaration with full OOP support
type ClassDecl struct {
	Name             string
	Parent           string           // parent class name for inheritance
	ParentTypeParams []TypeParam      // parent's generic type parameters (e.g., Container<T> has [T])
	Implements       []string         // interface names
	IsAbstract       bool             // whether the class is abstract
	AccessLevel      string           // "public", "private", "protected"
	IsSealed         bool             // whether the class is sealed
	Permits          []string         // names permitted to inherit
	Fields           []FieldDecl      // class fields
	Methods          []MethodDecl     // class methods
	Constructor      *ConstructorDecl // class constructor
	TypeParams       []TypeParam      // generic type parameters (e.g., [T, K, V])
}

func (*ClassDecl) node() {}
func (*ClassDecl) stmt() {}

// Enum declaration similar to Java enums
type EnumDecl struct {
	Name        string
	AccessLevel string // "public", "private", "protected"
	IsSealed    bool
	Permits     []string
	Values      []EnumValue
	Fields      []FieldDecl      // enum can have fields
	Methods     []MethodDecl     // enum can have methods
	Constructor *ConstructorDecl // enum constructor
}

func (*EnumDecl) node() {}
func (*EnumDecl) stmt() {}

// Enum value with optional constructor arguments
type EnumValue struct {
	Name string
	Args []Expr // constructor arguments
}

// Record declaration similar to Java records
type RecordDecl struct {
	Name        string
	AccessLevel string            // "public", "private", "protected"
	Components  []RecordComponent // record components
	Methods     []MethodDecl      // additional methods
}

func (*RecordDecl) node() {}
func (*RecordDecl) stmt() {}

// Record component (immutable field with accessor)
type RecordComponent struct {
	Name string
	Type *Type // Type annotation using unified type system
}

// Field declaration within a class
type FieldDecl struct {
	Name      string
	Type      *Type    // Type annotation using unified type system
	Modifiers []string // public, private, protected, static, final
	InitValue Expr     // optional initial value
}

func (*FieldDecl) node() {}
func (*FieldDecl) stmt() {}

// Annotation captures metadata attached to declarations such as methods.
type Annotation struct {
	Raw        string // original casing as written in source
	Normalized string // canonical lowercase form for comparisons
	Known      bool   // true if the annotation is registered in the runtime
}

// Method declaration within a class
type MethodDecl struct {
	Name        string
	Params      []Parameter
	ReturnType  *Type // Return type using unified type system
	Body        []Stmt
	Modifiers   []string // public, private, protected, static, abstract, final
	IsAbstract  bool
	IsOverride  bool         // whether this method is marked with @override
	Annotations []Annotation // annotations like @override, @deprecated, etc.
}

func (*MethodDecl) node() {}
func (*MethodDecl) stmt() {}

// Constructor declaration
type ConstructorDecl struct {
	Params []Parameter
	Body   []Stmt
}

func (*ConstructorDecl) node() {}
func (*ConstructorDecl) stmt() {}

// Instance of expression: obj instanceof Type with optional variable assignment
type InstanceOfExpr struct {
	Expr     Expr   // the object to check
	TypeName string // the type to check against
	Variable string // optional variable to assign (e.g., "foo" in "10 instanceof int final foo")
	Modifier string // optional modifier like "final"
}

func (*InstanceOfExpr) node() {}
func (*InstanceOfExpr) expr() {}

// Type expression for Sys.type(obj)
type TypeExpr struct {
	Expr Expr // the object to get type of
}

func (*TypeExpr) node() {}
func (*TypeExpr) expr() {}

// Operator codes for BinaryExpr and UnaryExpr.Op.
// Keep these AST-local to avoid import cycles.
const (
	OpPlus = iota + 1
	OpMinus
	OpMul
	OpDiv
	OpMod
	OpEq
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpAnd
	OpOr
	OpNot // unary
	OpNeg // unary minus
)

// Lambda expression: (params) => expr or (params) => do ... end
type LambdaExpr struct {
	Params     []Parameter // updated to support typed and variadic parameters
	Body       Expr        // expression body (single expression)
	BlockBody  []Stmt      // statement block body (multi-line lambda with do...end)
	ReturnType *Type       // Return type using unified type system
	IsBlock    bool        // true if this is a multi-line lambda
}

func (*LambdaExpr) node() {}
func (*LambdaExpr) expr() {}

// Thread expressions for concurrency
type ThreadSpawnExpr struct {
	Body []Stmt // thread body statements
}

func (*ThreadSpawnExpr) node() {}
func (*ThreadSpawnExpr) expr() {}

type ThreadJoinExpr struct {
	Thread Expr // thread expression to join
}

func (*ThreadJoinExpr) node() {}
func (*ThreadJoinExpr) expr() {}

// Channel creation: channel[Type]()
type ChannelExpr struct {
	ElemType string // Type of elements in the channel
}

func (*ChannelExpr) node() {}
func (*ChannelExpr) expr() {}

// Select statement for channel operations
type SelectStmt struct {
	Cases []SelectCase
	Pos   Position
}

// SelectCase represents a case in a select statement
type SelectCase struct {
	IsRecv  bool   // true for recv case, false for closed case
	RecvVar string // variable name for received value (optional)
	Channel Expr   // channel expression
	Body    []Stmt // statements to execute if this case is selected
}

func (*SelectStmt) node() {}
func (*SelectStmt) stmt() {}

// Switch statement for value/type/enum matching
type SwitchStmt struct {
	Expr    Expr         // expression to switch on (can be nil for type switches)
	Cases   []SwitchCase // switch cases
	Default []Stmt       // default case body (optional)
	Pos     Position
}

// SwitchCase represents a case in a switch statement
type SwitchCase struct {
	// For value matching: case value:
	Values []Expr // values to match (can have multiple values per case)

	// For type matching: case (x: Type):
	TypeName string // type name for type matching
	VarName  string // variable name for type matching (optional)

	// For enum matching: case Enum.VALUE:
	// Values field is used with enum member access expressions

	Body []Stmt // statements to execute if this case matches
}

func (*SwitchStmt) node() {}
func (*SwitchStmt) stmt() {}

// Ternary expression: condition ? trueBranch : falseBranch
type TernaryExpr struct {
	Condition   Expr
	TrueBranch  Expr
	FalseBranch Expr
}

func (*TernaryExpr) node() {}
func (*TernaryExpr) expr() {}

// RangeExpr represents range expressions like 1...10 or arr[1...3]
type RangeExpr struct {
	Start     Expr
	End       Expr
	Inclusive bool // true for ..., false for ..
}

func (*RangeExpr) node() {}
func (*RangeExpr) expr() {}

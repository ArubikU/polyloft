package common

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
)

// PrebuildedDefinition represents a forward reference to a class or interface
// Used for permits declarations where the target may not be defined yet
type PrebuildedDefinition struct {
	Name        string
	PackageName string
}

// ClassInstance represents an instance of a class
type ClassInstance struct {
	ClassName   string
	Fields      map[string]any
	Methods     map[string]Func
	ParentClass *ClassDefinition
}

// GenericType represents a generic type parameter (like T, E, K, V)
type GenericType struct {
	Name       string
	Bounds     []string // upper bounds (extends)
	IsVariadic bool
	Variance   string // "in" (contravariance), "out" (covariance), or "" (invariant)
}

// ClassDefinition represents a class definition
type ClassDefinition struct {
	Name         string
	Aliases      []string
	Type         *ast.Type // Unified type representation
	Parent       *ClassDefinition
	Implements   []*InterfaceDefinition // Changed from []string to proper references
	IsAbstract   bool
	AccessLevel  string // "public", "private", "protected"
	IsSealed     bool
	Permits      []*ClassDefinition     // Resolved permits (lazy resolution)
	PermitNames  []PrebuildedDefinition // Unresolved permits (stored at declaration time)
	FileName     string                 // file where class is defined
	PackageName  string                 // package/directory where class is defined
	Fields       map[string]FieldInfo
	Methods      map[string][]MethodInfo // Support overloading: multiple methods with same name
	Constructors []ConstructorInfo       // Support overloading: multiple constructors
	StaticFields map[string]any
	// Generic type support
	TypeParams []GenericType // Generic type parameters (e.g., [T, E extends Comparable])
	IsGeneric  bool          // Whether this class is generic
}

// FieldInfo contains field metadata
type FieldInfo struct {
	Name      string
	Type      *ast.Type // Type using unified type system
	Modifiers []string
	InitValue any
	IsStatic  bool
	IsPrivate bool
}

// EnumDefinition represents an enum declaration
type EnumDefinition struct {
	Name         string
	Type         *ast.Type // Unified type representation
	AccessLevel  string
	FileName     string
	PackageName  string
	Methods      map[string][]MethodInfo // Support overloading: multiple methods with same name
	Values       map[string]*EnumValueInstance
	Fields       map[string]FieldInfo
	Constructors []ConstructorInfo // Support overloading: multiple constructors
	IsSealed     bool
	Permits      []*ClassDefinition     // Resolved permits (lazy resolution)
	PermitNames  []PrebuildedDefinition // Unresolved permits (stored at declaration time)
}

// EnumValueInstance represents a single enum constant instance
type EnumValueInstance struct {
	Definition *EnumDefinition
	Name       string
	Ordinal    int
	Fields     map[string]any
	Methods    map[string]Func
}

// RecordDefinition represents a record declaration
type RecordDefinition struct {
	Name        string
	Type        *ast.Type // Unified type representation
	AccessLevel string
	FileName    string
	PackageName string
	Components  []ast.RecordComponent
	Methods     map[string][]MethodInfo // Support overloading: multiple methods with same name
}

// RecordInstance represents an instance of a record
type RecordInstance struct {
	Definition *RecordDefinition
	Values     map[string]any
	Methods    map[string]Func
}

// MethodInfo contains method metadata
type MethodInfo struct {
	Name        string
	Params      []ast.Parameter
	ReturnType  *ast.Type // Return type using unified type system
	Body        []ast.Stmt
	Modifiers   []string
	IsAbstract  bool
	IsStatic    bool
	IsPrivate   bool
	BuiltinImpl Func // Optional builtin implementation
}

// ParameterInfo contains parameter metadata - DEPRECATED: Use ast.Parameter instead
type ParameterInfo = ast.Parameter

// ConstructorInfo contains constructor metadata
type ConstructorInfo struct {
	Params      []ast.Parameter
	Body        []ast.Stmt
	BuiltinImpl Func // Optional builtin implementation
}

// Interface definition
type InterfaceDefinition struct {
	Name         string
	Type         *ast.Type                    // Unified type representation
	Methods      map[string][]MethodSignature // Support overloading: multiple methods with same name
	IsSealed     bool
	Permits      []*ClassDefinition     // Resolved permits (lazy resolution)
	PermitNames  []PrebuildedDefinition // Unresolved permits (stored at declaration time)
	StaticFields map[string]any
	AccessLevel  string
	FileName     string // file where interface is defined
	PackageName  string // package/directory where interface is defined
}

// MethodSignature for interface methods
type MethodSignature struct {
	Name        string
	Params      []ast.Parameter
	ReturnType  *ast.Type // Return type using unified type system
	HasDefault  bool
	DefaultBody []ast.Stmt
}

// Channel represents a communication channel for concurrent operations
type Channel struct {
	Ch     chan any // Exported for reflection in select
	closed bool
}

// NewChannel creates a new channel with optional buffer size
func NewChannel(bufferSize int) *Channel {
	return &Channel{
		Ch:     make(chan any, bufferSize),
		closed: false,
	}
}

// Send sends a value to the channel
func (c *Channel) Send(value any) error {
	if c.closed {
		return fmt.Errorf("send on closed channel")
	}
	c.Ch <- value
	return nil
}

// Recv receives a value from the channel
func (c *Channel) Recv() (any, bool) {
	val, ok := <-c.Ch
	return val, ok
}

// Close closes the channel
func (c *Channel) Close() {
	if !c.closed {
		c.closed = true
		close(c.Ch)
	}
}

// IsClosed returns whether the channel is closed
func (c *Channel) IsClosed() bool {
	return c.closed
}

// Env is a simple lexical environment for variables and functions.
type Env struct {
	Parent           *Env
	Vars             map[string]any
	Consts           map[string]bool
	Finals           map[string]bool
	Defers           []func() error      // stack of deferred calls (LIFO)
	FileName         string              // current file being executed
	PackageName      string              // current package/directory
	CurrentLine      int                 // current line number being executed (1-based)
	CurrentColumn    int                 // current column number being executed (1-based)
	CodeContext      []string            // context lines: [line-2, line-1, current line]
	SourceLines      []string            // all source lines (for hint generation)
	PositionStack    []PositionInfo      // stack of positions for better stack traces
	ImportedClasses  map[string]string   // className -> packageName, tracks imported classes
	ImportedPackages map[string]struct{} // packageName -> struct{}, tracks imported packages
}

//__ArrayClass__
//__BoolClass__
//__MapClass__
//__IntClass__
//__FloatClass__
//__StringClass__

//create a getter that consume env and a BuiltinType.ARRAY so we can collect them in any registry

type Builtin struct {
	Name        string
	IsClass     bool
	IsInterface bool
	IsEnum      bool
	IsRecord    bool
	IsPrimitive bool
	IsFunction  bool
}

var (
	BuiltinTypeBool              = Builtin{Name: "__BoolClass__", IsPrimitive: true}
	BuiltinTypeInt               = Builtin{Name: "__IntClass__", IsPrimitive: true}
	BuiltinTypeString            = Builtin{Name: "__StringClass__", IsPrimitive: true}
	BuiltinTypeMap               = Builtin{Name: "__MapClass__", IsPrimitive: false}
	BuiltinTypeFloat             = Builtin{Name: "__FloatClass__", IsPrimitive: true}
	BuiltinTypeNumber            = Builtin{Name: "__NumberClass__", IsPrimitive: true}
	BuiltinTypeArray             = Builtin{Name: "__ArrayClass__", IsPrimitive: false}
	BuiltinTypeGeneric           = Builtin{Name: "__GenericClass__", IsPrimitive: false}
	BuiltinTypeRange             = Builtin{Name: "__RangeClass__", IsPrimitive: false}
	BuiltinTypeList              = Builtin{Name: "__ListClass__", IsPrimitive: false}
	BuiltinTypeSet               = Builtin{Name: "__SetClass__", IsPrimitive: false}
	BuiltinTypeDeque             = Builtin{Name: "__DequeClass__", IsPrimitive: false}
	BuiltinTypePair              = Builtin{Name: "__PairClass__", IsPrimitive: false}
	BuiltinTypeTuple             = Builtin{Name: "__TupleClass__", IsPrimitive: false}
	BuiltinInterfaceIterable     = Builtin{Name: "__IterableInterface__", IsInterface: true}
	BuiltinInterfaceUnstructured = Builtin{Name: "__UnstructuredInterface__", IsInterface: true}
)

func (bt *Builtin) GetClassDefinition(env *Env) *ClassDefinition {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*ClassDefinition)
}
func (bt *Builtin) GetTypeDefinition(env *Env) *ast.Type {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*ClassDefinition).Type
}
func (bt *Builtin) GetInterfaceDefinition(env *Env) *InterfaceDefinition {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*InterfaceDefinition)
}
func (bt *Builtin) GetEnumDefinition(env *Env) *EnumDefinition {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*EnumDefinition)
}
func (bt *Builtin) GetRecordDefinition(env *Env) *RecordDefinition {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*RecordDefinition)
}
func (bt *Builtin) GetFunctionDefinition(env *Env) *FunctionDefinition {
	val, ok := env.Get(bt.Name)
	if !ok {
		return nil
	}
	return val.(*FunctionDefinition)
}

// PositionInfo holds position information for stack traces
type PositionInfo struct {
	File   string
	Line   int
	Column int
}

// NewEnv creates a new environment
func NewEnv() *Env {
	return &Env{
		Vars:             map[string]any{},
		Consts:           map[string]bool{},
		Finals:           map[string]bool{},
		ImportedClasses:  map[string]string{},
		ImportedPackages: map[string]struct{}{},
	}
}

// NewEnvWithContext creates a new environment with file and package context
func NewEnvWithContext(fileName, packageName string) *Env {
	return &Env{
		Vars:             map[string]any{},
		Consts:           map[string]bool{},
		Finals:           map[string]bool{},
		FileName:         fileName,
		PackageName:      packageName,
		ImportedClasses:  map[string]string{},
		ImportedPackages: map[string]struct{}{},
	}
}

// Child creates a child environment
func (e *Env) Child() *Env {
	// Copy imported classes and packages to child
	importedClasses := make(map[string]string)
	for k, v := range e.ImportedClasses {
		importedClasses[k] = v
	}
	importedPackages := make(map[string]struct{})
	for k, v := range e.ImportedPackages {
		importedPackages[k] = v
	}

	return &Env{
		Parent:           e,
		Vars:             map[string]any{},
		Consts:           map[string]bool{},
		Defers:           []func() error{}, // empty defer stack for child
		FileName:         e.FileName,       // inherit file context
		PackageName:      e.PackageName,    // inherit package context
		CurrentLine:      e.CurrentLine,    // inherit current line
		CurrentColumn:    e.CurrentColumn,  // inherit current column
		CodeContext:      e.CodeContext,    // inherit code context
		SourceLines:      e.SourceLines,    // inherit source lines
		PositionStack:    e.PositionStack,  // inherit position stack
		ImportedClasses:  importedClasses,  // inherit imported classes
		ImportedPackages: importedPackages, // inherit imported packages
	}
}

// GetFileName returns the current file name
func (e *Env) GetFileName() string { return e.FileName }

// GetPackageName returns the current package name
func (e *Env) GetPackageName() string { return e.PackageName }

// GetCurrentLine returns the current line number being executed
func (e *Env) GetCurrentLine() int { return e.CurrentLine }

// GetCodeContext returns the code context (previous 2 lines + current line)
func (e *Env) GetCodeContext() []string { return e.CodeContext }

// GetSourceLines returns all source lines
func (e *Env) GetSourceLines() []string { return e.SourceLines }

// SetSourceLines sets the source lines for this environment
func (e *Env) SetSourceLines(lines []string) { e.SourceLines = lines }

// SetCurrentLine updates the current line number and code context
func (e *Env) SetCurrentLine(line int, sourceLines []string) {
	e.CurrentLine = line

	// Build code context: [line-2, line-1, current line]
	e.CodeContext = make([]string, 0, 3)

	// Add line-2 if it exists
	if line >= 3 && line-3 < len(sourceLines) {
		e.CodeContext = append(e.CodeContext, sourceLines[line-3])
	}

	// Add line-1 if it exists
	if line >= 2 && line-2 < len(sourceLines) {
		e.CodeContext = append(e.CodeContext, sourceLines[line-2])
	}

	// Add current line if it exists
	if line >= 1 && line-1 < len(sourceLines) {
		e.CodeContext = append(e.CodeContext, sourceLines[line-1])
	}
}

// UpdateCurrentLine updates just the line number without updating context
// Useful when context is not available or not needed
func (e *Env) UpdateCurrentLine(line int) {
	e.CurrentLine = line
}

// PushPosition adds a position to the stack for better error tracking
func (e *Env) PushPosition(file string, line, column int) {
	if e.PositionStack == nil {
		e.PositionStack = make([]PositionInfo, 0, 10)
	}
	e.PositionStack = append(e.PositionStack, PositionInfo{
		File:   file,
		Line:   line,
		Column: column,
	})
}

// PopPosition removes the last position from the stack
func (e *Env) PopPosition() {
	if len(e.PositionStack) > 0 {
		e.PositionStack = e.PositionStack[:len(e.PositionStack)-1]
	}
}

// GetCurrentPosition returns the current position info
func (e *Env) GetCurrentPosition() PositionInfo {
	if len(e.PositionStack) > 0 {
		return e.PositionStack[len(e.PositionStack)-1]
	}
	return PositionInfo{
		File:   e.FileName,
		Line:   e.CurrentLine,
		Column: e.CurrentColumn,
	}
}

// Set sets a variable in the current environment
func (e *Env) Set(k string, v any) { e.Vars[k] = v }

// Get retrieves a variable from the environment chain
func (e *Env) Get(k string) (any, bool) {
	for cur := e; cur != nil; cur = cur.Parent {
		if v, ok := cur.Vars[k]; ok {
			return v, true
		}
	}
	return nil, false
}

// Define defines a new variable, optionally as a constant
func (e *Env) Define(k string, v any, kind string) {
	e.Vars[k] = v
	if kind == "const" {
		e.Consts[k] = true
	}
	if kind == "final" {
		e.Finals[k] = true
	}
}

// Func represents a Polyloft function that can be called from the engine
// This is the canonical function type used throughout the system
type Func func(env *Env, args []any) (any, error)

// ExtractFunc extracts a Func from various function-like wrappers
func ExtractFunc(val any) (Func, bool) {
	switch v := val.(type) {
	case Func:
		return v, true
	case *FunctionDefinition:
		return v.Func, true
	case *LambdaDefinition:
		return v.Func, true
	case *ClassConstructor:
		return v.Func, true
	default:
		return nil, false
	}
}

// FunctionDefinition represents a top-level function with metadata
type FunctionDefinition struct {
	Name        string
	Func        Func
	Params      []ast.Parameter // Function parameters with types
	ReturnType  *ast.Type       // Return type (can be inferred)
	AccessLevel string          // "public", "private", "protected"
	Modifiers   []string
	FileName    string // file where function is defined
	PackageName string // package/directory where function is defined
}

// LambdaDefinition represents a lambda expression with type information
type LambdaDefinition struct {
	Func       Func            // The actual lambda function
	Params     []ast.Parameter // Lambda parameters with types
	ReturnType *ast.Type       // Return type (can be inferred)
}

// ClassConstructor wraps a class constructor function with metadata
// This allows Sys.type to identify and format class names properly
type ClassConstructor struct {
	Definition *ClassDefinition
	Func       Func
}

// EnumConstructor wraps an enum object with metadata
// This allows Sys.type to identify and format enum names properly
type EnumConstructor struct {
	Definition *EnumDefinition
	EnumObject map[string]any
}

// TypeError represents a type error in Polyloft operations
type TypeError struct {
	Message string
}

func (e TypeError) Error() string {
	return e.Message
}

// RuntimeError represents a runtime error with more context
type RuntimeError struct {
	Message string
	Context string
}

func (e RuntimeError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s (in %s)", e.Message, e.Context)
	}
	return e.Message
}

// ArityError represents an arity mismatch error
type ArityError struct {
	Expected int
	Got      int
}

func (e ArityError) Error() string {
	return fmt.Sprintf("arity mismatch: want %d, got %d", e.Expected, e.Got)
}

// UndefinedError represents an undefined identifier error
type UndefinedError struct {
	Name string
}

func (e UndefinedError) Error() string {
	return fmt.Sprintf("undefined identifier: %s", e.Name)
}

// IndexError represents an index out of range error
type IndexError struct {
	Index int
	Size  int
}

func (e IndexError) Error() string {
	return fmt.Sprintf("index out of range: %d (size: %d)", e.Index, e.Size)
}

// Sentinel types for control flow
type BreakSentinel struct{}
type ContinueSentinel struct{}

// InferReturnType analyzes function body statements to infer the return type
func InferReturnType(body []ast.Stmt, env *Env) *ast.Type {
	returnTypes := collectReturnTypes(body, env)

	// If no return statements found, return nil (void/Any)
	if len(returnTypes) == 0 {
		return nil
	}

	// If more than 5 different return types, return Any
	if len(returnTypes) > 5 {
		return &ast.Type{Name: "any"}
	}

	// If all returns are the same type, return that type
	if len(returnTypes) == 1 {
		return returnTypes[0]
	}

	// Check if all returns are Int or Float (should be Number)
	allNumbers := true
	for _, rt := range returnTypes {
		if rt.Name != "int" && rt.Name != "float" {
			allNumbers = false
			break
		}
	}
	if allNumbers {
		return BuiltinTypeNumber.GetTypeDefinition(env)
	}

	// Multiple different types - create union type or return Any
	return ast.ANY
}

// collectReturnTypes collects all return types from statements
func collectReturnTypes(stmts []ast.Stmt, env *Env) []*ast.Type {
	var types []*ast.Type
	seen := make(map[string]bool)

	for _, stmt := range stmts {
		collectReturnTypesFromStmt(stmt, &types, seen, env)
	}

	return types
}

// collectReturnTypesFromStmt recursively collects return types from a statement
func collectReturnTypesFromStmt(stmt ast.Stmt, types *[]*ast.Type, seen map[string]bool, env *Env) {
	switch s := stmt.(type) {
	case *ast.ReturnStmt:
		// Infer type from return expression
		if s.Value != nil {
			inferredType := InferExprType(s.Value, env)
			if inferredType != nil {
				typeName := inferredType.Name
				if !seen[typeName] {
					seen[typeName] = true
					*types = append(*types, inferredType)
				}
			}
		}
	case *ast.IfStmt:
		// Check all branches
		for _, clause := range s.Clauses {
			collectReturnTypesFromStmts(clause.Body, types, seen, env)
		}
		if s.Else != nil {
			collectReturnTypesFromStmts(s.Else, types, seen, env)
		}
	case *ast.ForInStmt:
		collectReturnTypesFromStmts(s.Body, types, seen, env)
	case *ast.LoopStmt:
		collectReturnTypesFromStmts(s.Body, types, seen, env)
	case *ast.TryStmt:
		collectReturnTypesFromStmts(s.Body, types, seen, env)
		for _, catchClause := range s.Catches {
			collectReturnTypesFromStmts(catchClause.Body, types, seen, env)
		}
		if s.Finally != nil {
			collectReturnTypesFromStmts(s.Finally, types, seen, env)
		}
	}
}

// collectReturnTypesFromStmts collects return types from a slice of statements
func collectReturnTypesFromStmts(stmts []ast.Stmt, types *[]*ast.Type, seen map[string]bool, env *Env) {
	for _, stmt := range stmts {
		collectReturnTypesFromStmt(stmt, types, seen, env)
	}
}

// InferExprType infers the type of an expression
func InferExprType(expr ast.Expr, env *Env) *ast.Type {
	switch e := expr.(type) {
	case *ast.NumberLit:
		// Check if it's an int or float
		switch e.Value.(type) {
		case int, int32, int64:
			return BuiltinTypeInt.GetTypeDefinition(env)
		case float32, float64:
			return BuiltinTypeFloat.GetTypeDefinition(env)
		default:
			return BuiltinTypeNumber.GetTypeDefinition(env)
		}
	case *ast.StringLit, *ast.InterpolatedStringLit:
		return BuiltinTypeString.GetTypeDefinition(env)
	case *ast.BoolLit:
		return BuiltinTypeBool.GetTypeDefinition(env)
	case *ast.ArrayLit:
		return BuiltinTypeArray.GetTypeDefinition(env)
	case *ast.MapLit:
		return BuiltinTypeMap.GetTypeDefinition(env)
	case *ast.NilLit:
		return ast.NIL
	default:
		return ast.ANY
	}
}

// formatFunctionType formats a function type as Function<Param1,Param2,...,ReturnType>
func formatFunctionType(params []ast.Parameter, returnType *ast.Type) string {
	var parts []string

	// Add parameter types
	for _, param := range params {
		var paramType string
		if param.Type != nil {
			paramType = getTypeString(param.Type)
		} else {
			paramType = "Any"
		}

		// Add ... notation for variadic parameters
		if param.IsVariadic {
			paramType += "..."
		}

		parts = append(parts, paramType)
	}

	// Add return type
	if returnType != nil {
		parts = append(parts, getTypeString(returnType))
	} else {
		parts = append(parts, "Any")
	}

	if len(parts) == 0 {
		return "Function"
	}

	return fmt.Sprintf("Function<%s>", strings.Join(parts, ","))
}

// getTypeString returns a string representation of an ast.Type
func getTypeString(t *ast.Type) string {
	if t == nil {
		return "Any"
	}

	// Handle union types
	if t.IsUnion && len(t.UnionTypes) > 0 {
		// Check if it's Int | Float which should be simplified to Number
		if len(t.UnionTypes) == 2 {
			hasInt := false
			hasFloat := false
			for _, ut := range t.UnionTypes {
				if ut.Name == "int" {
					hasInt = true
				} else if ut.Name == "float" {
					hasFloat = true
				}
			}
			if hasInt && hasFloat {
				return "Number"
			}
		}

		// Otherwise, format as union
		var typeNames []string
		for _, ut := range t.UnionTypes {
			typeNames = append(typeNames, strings.Title(ut.Name))
		}
		return strings.Join(typeNames, "|")
	}

	// Handle generic types with type parameters
	if len(t.TypeParams) > 0 {
		var paramStrs []string
		for _, tp := range t.TypeParams {
			paramStrs = append(paramStrs, getTypeString(tp))
		}
		return fmt.Sprintf("%s<%s>", strings.Title(t.Name), strings.Join(paramStrs, ","))
	}

	// Return capitalized type name
	return strings.Title(t.Name)
}

// GetTypeName returns the type name of a value
func GetTypeName(val any) string {
	switch v := val.(type) {
	case *ClassConstructor:
		// Format as "Class {Name}@{Package}"
		pkg := v.Definition.PackageName
		if pkg == "" {
			pkg = "default"
		}
		return fmt.Sprintf("Class %s@%s", v.Definition.Name, pkg)
	case *EnumConstructor:
		// Format as "Enum {Name}@{Package}"
		pkg := v.Definition.PackageName
		if pkg == "" {
			pkg = "default"
		}
		return fmt.Sprintf("Enum %s@%s", v.Definition.Name, pkg)
	case *ClassInstance:
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
		case "Array", "List", "Set", "Deque":
			// For collections, infer element type if no type args stored

			if typeArgs, ok := v.Fields["__type_args__"].([]string); ok && len(typeArgs) > 0 {
				return fmt.Sprintf("%s<%s>", v.ClassName, strings.Join(typeArgs, ", "))
			}
			// Infer from elements - handle both array and map storage
			if items, ok := v.Fields["_items"].([]any); ok {
				return inferCollectionType(v.ClassName, items)
			}
			if itemsPtr, ok := v.Fields["_items"].(*[]any); ok {
				return inferCollectionType(v.ClassName, *itemsPtr)
			}
			if itemsMap, ok := v.Fields["_items"].(map[uint64]any); ok {
				// For Sets, extract values from the map
				values := make([]any, 0, len(itemsMap))
				for _, val := range itemsMap {
					values = append(values, val)
				}
				return inferCollectionType(v.ClassName, values)
			}
			return v.ClassName
		case "Map":
			// For maps, infer key/value types if no type args stored
			if typeArgs, ok := v.Fields["__type_args__"].([]string); ok && len(typeArgs) == 2 {
				return fmt.Sprintf("Map<%s, %s>", strings.Title(typeArgs[0]), strings.Title(typeArgs[1]))
			}

			return inferMapType(MapToObject(v))
		default:
			// Check if this is a generic class instance
			if typeArgs, ok := v.Fields["__type_args__"].([]string); ok && len(typeArgs) > 0 {
				return fmt.Sprintf("%s<%s>", v.ClassName, strings.Join(typeArgs, ", "))
			}
			return v.ClassName
		}

	case *EnumValueInstance:
		if v.Definition != nil {
			return v.Definition.Name
		}
		return "enum"
	case *RecordInstance:
		if v.Definition != nil {
			return v.Definition.Name
		}
		return "record"
	case int, int32, int64, float32, float64, string, bool, []any, map[string]any:
		// Native Go types should be wrapped in Generic builtin
		return "Generic"
	case nil:
		return ast.NIL.Name
	case Func:
		return "function"
	case *FunctionDefinition:
		return formatFunctionType(v.Params, v.ReturnType)
	case *LambdaDefinition:
		return formatFunctionType(v.Params, v.ReturnType)
	default:
		fmt.Println("Unknown type for GetTypeName:", v)
		return v.(fmt.Stringer).String()
	}
}

// mapEntry stores a key-value pair in the internal map storage (from engine package)

func MapToObject(mapInstance *ClassInstance) map[any]any {
	dataField, ok := mapInstance.Fields["_data"]
	if !ok {
		return nil
	}

	// Handle the internal hash-based storage format
	if internalMap, ok := dataField.(map[uint64]*ast.MapEntry); ok {
		result := make(map[any]any, len(internalMap))
		for _, entry := range internalMap {
			if entry != nil {
				result[entry.Key] = entry.Value
			}
		}
		return result
	}

	// Fallback for direct map[any]any format
	if directMap, ok := dataField.(map[any]any); ok {
		return directMap
	}
	return nil
}

func inferCollectionType(className string, items []any) string {
	if len(items) == 0 {
		return className
	}
	firstType := ""
	allSame := true
	hasInt := false
	hasFloat := false
	for i, item := range items {
		itemType := GetTypeName(item)
		if i == 0 {
			firstType = itemType
		} else if itemType != firstType {
			allSame = false
		}
		if itemType == "int" {
			hasInt = true
		} else if itemType == "float" {
			hasFloat = true
		}
	}

	// Si todos los elementos son del mismo tipo
	if allSame && firstType != "" {
		// Capitalizar primera letra para mostrar
		displayType := strings.Title(firstType)
		return fmt.Sprintf("%s<%s>", className, displayType)
	}

	// Si hay mezcla de int y float, es Array<Number>
	if hasInt && hasFloat {
		return fmt.Sprintf("%s<Number>", className)
	}

	// De lo contrario, es un array mixto
	return fmt.Sprintf("%s<Any>", className)

}

func inferMapType(m map[any]any) string {
	if len(m) == 0 {
		return "Map"
	}

	// Determinar tipo de clave (asumimos que todas son string, pero se valida)
	var keyType string = "String"
	hasMixedKeys := false
	for k := range m {
		kType := GetTypeName(k)
		if keyType == "String" {
			keyType = strings.Title(kType)
		} else if keyType != strings.Title(kType) {
			hasMixedKeys = true
			break
		}
	}
	if hasMixedKeys {
		keyType = "Any"
	}

	// Inferir tipo de valores
	valueTypes := make(map[string]bool)
	hasInt := false
	hasFloat := false

	for _, v := range m {
		vType := GetTypeName(v)

		// Normalizar nombres base
		switch vType {
		case "int":
			hasInt = true
		case "float":
			hasFloat = true
		}

		valueTypes[strings.Title(vType)] = true
	}

	// Si hay mezcla de int y float → simplificar a Number
	if hasInt && hasFloat {
		if len(valueTypes) == 2 || len(valueTypes) == 1 {
			return fmt.Sprintf("Map<%s, Number>", keyType)
		}
		delete(valueTypes, "Int")
		delete(valueTypes, "Float")
		valueTypes["Number"] = true
	}

	// Si todos los valores son del mismo tipo
	if len(valueTypes) == 1 {
		for t := range valueTypes {
			return fmt.Sprintf("Map<%s, %s>", keyType, t)
		}
	}

	// Si hay varios tipos → unión (ej: Map<String, Float | String>)
	union := ""
	i := 0
	for t := range valueTypes {
		if i > 0 {
			union += " | "
		}
		union += t
		i++
	}

	return fmt.Sprintf("Map<%s, %s>", keyType, union)
}

// GetType returns the ast.Type for a value
func GetType(val any) *ast.Type {
	switch v := val.(type) {
	case *ClassConstructor:
		if v.Definition != nil && v.Definition.Type != nil {
			return v.Definition.Type
		}
		return nil
	case *EnumConstructor:
		if v.Definition != nil && v.Definition.Type != nil {
			return v.Definition.Type
		}
		return nil
	case *ClassInstance:
		if v.ParentClass != nil && v.ParentClass.Type != nil {
			return v.ParentClass.Type
		}
		return nil
	case *EnumValueInstance:
		if v.Definition != nil && v.Definition.Type != nil {
			return v.Definition.Type
		}
		return nil
	case *RecordInstance:
		if v.Definition != nil && v.Definition.Type != nil {
			return v.Definition.Type
		}
		return nil
	case int, int32, int64, float32, float64, string, bool, []any, map[string]any:
		// Native Go types should be wrapped in Generic builtin
		return &ast.Type{Name: "Generic", IsBuiltin: true}
	case nil:
		return ast.NIL
	case Func:
		return &ast.Type{Name: "function", IsBuiltin: true}
	default:
		return nil
	}
}

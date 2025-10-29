package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// ========================================
// ClassBuilder - Para crear clases builtin
// ========================================

// ClassBuilder facilita la creación de clases builtin con soporte completo para genéricos
type ClassBuilder struct {
	name         string
	parent       *ClassDefinition
	implements   []*common.InterfaceDefinition
	isAbstract   bool
	accessLevel  string
	fileName     string
	packageName  string
	fields       map[string]FieldInfo
	methods      map[string][]MethodInfo
	staticFields map[string]any
	constructors []ConstructorInfo
	typeParams   []common.GenericType
	isGeneric    bool
	aliases      []string
	Type         *ast.Type
}
type InterfaceBuilder struct {
	name         string
	methods      map[string][]MethodSignature
	isSealed     bool
	permits      []string
	staticFields map[string]any
	accessLevel  string
	fileName     string
	packageName  string
	typeParams   []common.GenericType
	isGeneric    bool
	Type         *ast.Type
}
type EnumBuilder struct {
	name        string
	values      []string
	accessLevel string
	fileName    string
	packageName string
	Type        *ast.Type
}
type RecordBuilder struct {
	name        string
	components  []RecordComponent
	methods     map[string]MethodInfo
	accessLevel string
	fileName    string
	packageName string
	Type        *ast.Type
}
type GenericTypeBuilder struct {
	name        string
	castFunc    common.Func
	validators  []func(any) bool
	Type        *ast.Type
	accessLevel string
	fileName    string
	packageName string
}

// NewClassBuilder crea un nuevo builder para clases
func NewClassBuilder(name string) *ClassBuilder {
	return &ClassBuilder{
		name:         name,
		accessLevel:  "public",
		fileName:     "builtin",
		packageName:  "builtin",
		fields:       make(map[string]FieldInfo),
		methods:      make(map[string][]MethodInfo),
		staticFields: make(map[string]any),
		constructors: []ConstructorInfo{},
		aliases:      []string{},
		Type:         &ast.Type{Name: name, IsBuiltin: true, IsClass: true},
	}
}

// SetParent establece la clase padre
func (cb *ClassBuilder) SetParent(parent *ClassDefinition) *ClassBuilder {
	cb.parent = parent
	return cb
}
func (cb *ClassBuilder) AddAlias(alias string) *ClassBuilder {
	cb.aliases = append(cb.aliases, alias)
	// Also update the Type.Aliases to keep them in sync
	cb.Type.Aliases = append(cb.Type.Aliases, alias)
	return cb
}

// GetType returns the ast.Type reference for this class builder
// This should be called after all AddAlias calls to ensure aliases are synchronized
func (cb *ClassBuilder) GetType() *ast.Type {
	return cb.Type
}

// SetAbstract marca la clase como abstracta
func (cb *ClassBuilder) SetAbstract(isAbstract bool) *ClassBuilder {
	cb.isAbstract = isAbstract
	return cb
}

// AddInterface agrega una interfaz implementada
func (cb *ClassBuilder) AddInterface(interfaceDef *common.InterfaceDefinition) *ClassBuilder {
	cb.implements = append(cb.implements, interfaceDef)
	return cb
}

// AddTypeParameters agrega múltiples parámetros de tipo genérico
func (cb *ClassBuilder) AddTypeParameters(typeParams []common.GenericType) *ClassBuilder {
	if len(typeParams) > 0 {
		cb.typeParams = append(cb.typeParams, typeParams...)
		cb.isGeneric = true
	}
	return cb
}

// AddField agrega un campo a la clase
func (cb *ClassBuilder) AddField(name string, fieldType *ast.Type, modifiers []string) *ClassBuilder {
	cb.fields[name] = FieldInfo{
		Name:      name,
		Type:      fieldType,
		Modifiers: modifiers,
		IsStatic:  contains(modifiers, "static"),
		IsPrivate: contains(modifiers, "private"),
	}
	return cb
}

// AddMethod agrega un método a la clase (supports overloading)
func (cb *ClassBuilder) AddMethod(name string, returnType *ast.Type, params []ast.Parameter, body []ast.Stmt, modifiers []string) *ClassBuilder {
	methodInfo := MethodInfo{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
		Modifiers:  modifiers,
		IsAbstract: contains(modifiers, "abstract"),
		IsStatic:   contains(modifiers, "static"),
		IsPrivate:  contains(modifiers, "private"),
	}
	cb.methods[name] = append(cb.methods[name], methodInfo)
	return cb
}

// AddBuiltinMethod agrega un método builtin (implementado en Go, supports overloading)
func (cb *ClassBuilder) AddBuiltinMethod(name string, returnType *ast.Type, params []ast.Parameter, fn common.Func, modifiers []string) *ClassBuilder {
	methodInfo := MethodInfo{
		Name:        name,
		Params:      params,
		ReturnType:  returnType,
		Body:        nil,
		Modifiers:   modifiers,
		IsAbstract:  false,
		IsStatic:    contains(modifiers, "static"),
		IsPrivate:   contains(modifiers, "private"),
		BuiltinImpl: fn,
	}
	cb.methods[name] = append(cb.methods[name], methodInfo)
	return cb
}

// AddSimpleMethod agrega un método sin parámetros que retorna un valor
func (cb *ClassBuilder) AddSimpleMethod(name string, returnType *ast.Type, fn common.Func) *ClassBuilder {
	return cb.AddBuiltinMethod(name, returnType, []ast.Parameter{}, fn, []string{})
}

// AddMethodWithParams agrega un método con parámetros especificados por nombres y tipos
func (cb *ClassBuilder) AddMethodWithParams(name string, returnType *ast.Type, paramNames []string, paramTypes []*ast.Type, fn common.Func) *ClassBuilder {
	params := make([]ast.Parameter, len(paramNames))
	for i, paramName := range paramNames {
		var paramType *ast.Type
		if i < len(paramTypes) {
			paramType = paramTypes[i]
		}
		params[i] = ast.Parameter{Name: paramName, Type: paramType}
	}
	return cb.AddBuiltinMethod(name, returnType, params, fn, []string{})
}

// AddStaticField agrega un campo estático a la clase
func (cb *ClassBuilder) AddStaticField(name string, value any) *ClassBuilder {
	cb.staticFields[name] = value
	return cb
}

// AddStaticMethod agrega un método estático a la clase (supports overloading)
func (cb *ClassBuilder) AddStaticMethod(name string, returnType *ast.Type, params []ast.Parameter, fn common.Func) *ClassBuilder {
	methodInfo := MethodInfo{
		Name:        name,
		Params:      params,
		ReturnType:  returnType,
		Body:        nil,
		Modifiers:   []string{"static"},
		IsStatic:    true,
		IsPrivate:   false,
		BuiltinImpl: fn,
	}
	cb.methods[name] = append(cb.methods[name], methodInfo)
	return cb
}

// AddConstructor agrega un constructor a la clase (supports overloading)
func (cb *ClassBuilder) AddConstructor(params []ast.Parameter, body []ast.Stmt) *ClassBuilder {
	cb.constructors = append(cb.constructors, ConstructorInfo{
		Params: params,
		Body:   body,
	})
	return cb
}

// SetConstructor establece el constructor de la clase (deprecated, use AddConstructor for overloading)
func (cb *ClassBuilder) SetConstructor(params []ast.Parameter, body []ast.Stmt) *ClassBuilder {
	return cb.AddConstructor(params, body)
}

// AddBuiltinConstructor agrega un constructor builtin (implementado en Go, supports overloading)
func (cb *ClassBuilder) AddBuiltinConstructor(params []ast.Parameter, fn common.Func) *ClassBuilder {
	cb.constructors = append(cb.constructors, ConstructorInfo{
		Params:      params,
		Body:        nil,
		BuiltinImpl: fn,
	})
	return cb
}

// SetBuiltinConstructor establece un constructor builtin (deprecated, use AddBuiltinConstructor for overloading)
func (cb *ClassBuilder) SetBuiltinConstructor(params []ast.Parameter, fn common.Func) *ClassBuilder {
	cb.constructors = append(cb.constructors, ConstructorInfo{
		Params:      params,
		Body:        nil,
		BuiltinImpl: fn,
	})
	return cb
}

// Build construye y registra la clase
func (cb *ClassBuilder) Build(env *Env) (*ClassDefinition, error) {
	// Use the interface definitions directly
	implements := cb.implements

	classDef := &ClassDefinition{
		Name:         cb.name,
		Parent:       cb.parent,
		Implements:   implements,
		IsAbstract:   cb.isAbstract,
		AccessLevel:  cb.accessLevel,
		FileName:     cb.fileName,
		PackageName:  cb.packageName,
		Fields:       cb.fields,
		Methods:      cb.methods,
		Constructors: cb.constructors,
		StaticFields: cb.staticFields,
		TypeParams:   cb.typeParams,
		IsGeneric:    cb.isGeneric,
		Aliases:      cb.aliases,
		Type:         cb.Type,
	}

	// Register as builtin class (always available)
	builtinClasses[cb.name] = classDef

	constructorFunc := common.Func(func(callEnv *common.Env, args []any) (any, error) {
		return createClassInstance(classDef, (*Env)(callEnv), args)
	})

	classConstructor := &common.ClassConstructor{
		Definition: classDef,
		Func:       constructorFunc,
	}

	env.Set("__"+cb.name+"Class"+"__", classDef)
	env.Set(cb.name, classConstructor)
	
	// Register aliases
	for _, alias := range cb.aliases {
		env.Set(alias, classConstructor)
	}
	
	return classDef, nil
}

// BuildStatic construye y registra una clase con solo miembros estáticos (no instanciable)
// Para módulos como Math, Sys, IO que solo tienen métodos y campos estáticos
func (cb *ClassBuilder) BuildStatic(env *Env) (*ClassDefinition, error) {
	// Convert implements strings to interface definitions (empty for builtins)
	var implements []*common.InterfaceDefinition

	classDef := &ClassDefinition{
		Name:         cb.name,
		Parent:       cb.parent,
		Implements:   implements,
		IsAbstract:   true, // Marcar como abstracta para evitar instanciación
		AccessLevel:  cb.accessLevel,
		FileName:     cb.fileName,
		PackageName:  cb.packageName,
		Fields:       cb.fields,
		Methods:      cb.methods,
		Constructors: nil, // No constructors para clases estáticas
		StaticFields: cb.staticFields,
		TypeParams:   cb.typeParams,
		IsGeneric:    cb.isGeneric,
	}

	builtinClasses[cb.name] = classDef

	env.Set(cb.name, classDef)
	return classDef, nil
}

// BuildAndGet construye la clase y retorna tanto la definición como el constructor
func (cb *ClassBuilder) BuildAndGet(env *Env) (*ClassDefinition, common.Func, error) {
	classDef, err := cb.Build(env)
	if err != nil {
		return nil, nil, err
	}

	constructorVal, _ := env.Get(cb.name)
	switch ctor := constructorVal.(type) {
	case common.Func:
		return classDef, ctor, nil
	case *common.ClassConstructor:
		return classDef, ctor.Func, nil
	default:
		return classDef, nil, ThrowTypeError(nil, "Func or ClassConstructor", constructorVal)
	}
}

// ========================================
// InterfaceBuilder - Para crear interfaces
// ========================================

// InterfaceBuilder facilita la creación de interfaces builtin

// NewInterfaceBuilder crea un nuevo builder para interfaces
func NewInterfaceBuilder(name string) *InterfaceBuilder {
	return &InterfaceBuilder{
		name:         name,
		methods:      make(map[string][]MethodSignature),
		staticFields: make(map[string]any),
		accessLevel:  "public",
		fileName:     "builtin",
		packageName:  "builtin",
		Type:         &ast.Type{Name: name, IsBuiltin: true, IsInterface: true},
	}
}
func (ib *InterfaceBuilder) AddTypeParameters(typeParams []common.GenericType) *InterfaceBuilder {
	ib.typeParams = append(ib.typeParams, typeParams...)
	ib.isGeneric = true
	return ib
}

// GetType returns the ast.Type reference for this interface builder
func (ib *InterfaceBuilder) GetType() *ast.Type {
	return ib.Type
}

// AddMethod agrega un método a la interfaz (supports overloading)
func (ib *InterfaceBuilder) AddMethod(name string, returnType any, params []ast.Parameter) *InterfaceBuilder {
	switch returnType.(type) {
	case string:
		sig := MethodSignature{
			Name:       name,
			Params:     params,
			ReturnType: ast.TypeFromString(returnType.(string)),
			HasDefault: false,
		}
		ib.methods[name] = append(ib.methods[name], sig)
		return ib
	case common.ClassDefinition:
		sig := MethodSignature{
			Name:       name,
			Params:     params,
			ReturnType: returnType.(common.ClassDefinition).Type,
			HasDefault: false,
		}
		ib.methods[name] = append(ib.methods[name], sig)
		return ib
	case ast.Type:
		t := returnType.(ast.Type)
		sig := MethodSignature{
			Name:       name,
			Params:     params,
			ReturnType: &t,
			HasDefault: false,
		}
		ib.methods[name] = append(ib.methods[name], sig)
		return ib
	}
	return ib
}

// AddDefaultMethod agrega un método con implementación por defecto (supports overloading)
func (ib *InterfaceBuilder) AddDefaultMethod(name string, returnType *ast.Type, params []ast.Parameter, body []ast.Stmt) *InterfaceBuilder {
	sig := MethodSignature{
		Name:        name,
		Params:      params,
		ReturnType:  returnType,
		HasDefault:  true,
		DefaultBody: body,
	}
	ib.methods[name] = append(ib.methods[name], sig)
	return ib
}

// SetSealed marca la interfaz como sealed
func (ib *InterfaceBuilder) SetSealed(isSealed bool, permits []string) *InterfaceBuilder {
	ib.isSealed = isSealed
	ib.permits = permits
	return ib
}

// AddStaticField agrega un campo estático a la interfaz
func (ib *InterfaceBuilder) AddStaticField(name string, value any) *InterfaceBuilder {
	ib.staticFields[name] = value
	return ib
}

// Build construye y registra la interfaz
func (ib *InterfaceBuilder) Build(env *Env) (*InterfaceDefinition, error) {
	// Convert permits strings to class definitions (empty for builtins)
	var permits []*ClassDefinition

	interfaceDef := &InterfaceDefinition{
		Name:         ib.name,
		Methods:      ib.methods,
		IsSealed:     ib.isSealed,
		Permits:      permits,
		StaticFields: ib.staticFields,
		AccessLevel:  ib.accessLevel,
		FileName:     ib.fileName,
		PackageName:  ib.packageName,
	}

	env.Set("__"+ib.name+"Interface"+"__", interfaceDef)
	interfaceRegistry[ib.name] = interfaceDef
	return interfaceDef, nil
}

// ========================================
// EnumBuilder - Para crear enums
// ========================================

// EnumBuilder facilita la creación de enums builtin

// NewEnumBuilder crea un nuevo builder para enums
func NewEnumBuilder(name string) *EnumBuilder {
	return &EnumBuilder{
		name:        name,
		values:      []string{},
		accessLevel: "public",
		fileName:    "builtin",
		packageName: "builtin",
		Type:        &ast.Type{Name: name, IsBuiltin: true, IsEnum: true},
	}
}

// AddValue agrega un valor al enum
func (eb *EnumBuilder) AddValue(value string) *EnumBuilder {
	eb.values = append(eb.values, value)
	return eb
}

// AddValues agrega múltiples valores al enum
func (eb *EnumBuilder) AddValues(values ...string) *EnumBuilder {
	eb.values = append(eb.values, values...)
	return eb
}

// Build construye y registra el enum
func (eb *EnumBuilder) Build(env *Env) (map[string]any, error) {
	enumMap := make(map[string]any)

	// Crear objetos para cada valor del enum
	for _, value := range eb.values {
		enumMap[value] = &EnumValue{
			EnumName: eb.name,
			Value:    value,
		}
	}

	env.Set("__"+eb.name+"Enum"+"__", enumMap)
	env.Set(eb.name, enumMap)

	return enumMap, nil
}

// EnumValue representa un valor de enum
type EnumValue struct {
	EnumName string
	Value    string
}

// ========================================
// RecordBuilder - Para crear records
// ========================================

// RecordBuilder facilita la creación de records builtin

// RecordComponent representa un componente de un record
type RecordComponent struct {
	Name string
	Type string
}

// NewRecordBuilder crea un nuevo builder para records
func NewRecordBuilder(name string) *RecordBuilder {
	return &RecordBuilder{
		name:        name,
		components:  []RecordComponent{},
		methods:     make(map[string]MethodInfo),
		accessLevel: "public",
		fileName:    "builtin",
		packageName: "builtin",
		Type:        &ast.Type{Name: name, IsBuiltin: true, IsRecord: true},
	}
}

// AddComponent agrega un componente al record
func (rb *RecordBuilder) AddComponent(name, componentType string) *RecordBuilder {
	rb.components = append(rb.components, RecordComponent{
		Name: name,
		Type: componentType,
	})
	return rb
}

// AddComponents agrega múltiples componentes al record
func (rb *RecordBuilder) AddComponents(components map[string]string) *RecordBuilder {
	for name, componentType := range components {
		rb.components = append(rb.components, RecordComponent{
			Name: name,
			Type: componentType,
		})
	}
	return rb
}

// AddMethod agrega un método al record
func (rb *RecordBuilder) AddMethod(name string, returnType *ast.Type, params []ast.Parameter, fn common.Func) *RecordBuilder {
	rb.methods[name] = MethodInfo{
		Name:        name,
		Params:      params,
		ReturnType:  returnType,
		Body:        nil,
		BuiltinImpl: fn,
	}
	return rb
}

// Build construye y registra el record
func (rb *RecordBuilder) Build(env *Env) (*RecordDefinition, error) {
	// Convertir componentes al formato esperado
	components := make(map[string]string)
	for _, comp := range rb.components {
		components[comp.Name] = comp.Type
	}

	recordDef := &RecordDefinition{
		Name:        rb.name,
		Components:  components,
		Methods:     rb.methods,
		AccessLevel: rb.accessLevel,
		FileName:    rb.fileName,
		PackageName: rb.packageName,
	}

	// Crear constructor para el record
	constructor := common.Func(func(callEnv *common.Env, args []any) (any, error) {
		return createRecordInstance(recordDef, (*Env)(callEnv), args)
	})

	env.Set("__"+rb.name+"Record"+"__", recordDef)
	env.Set(rb.name, constructor)
	return recordDef, nil
}

// RecordDefinition representa la definición de un record
type RecordDefinition struct {
	Name        string
	Components  map[string]string
	Methods     map[string]MethodInfo
	AccessLevel string
	FileName    string
	PackageName string
}

// createRecordInstance crea una instancia de record
func createRecordInstance(recordDef *RecordDefinition, env *Env, args []any) (any, error) {
	if len(args) != len(recordDef.Components) {
		return nil, ThrowArityError(env, len(recordDef.Components), len(args))
	}

	// Crear instancia del record
	instance := make(map[string]any)
	i := 0
	for name := range recordDef.Components {
		instance[name] = args[i]
		i++
	}

	// Agregar métodos
	for methodName, methodInfo := range recordDef.Methods {
		if methodInfo.BuiltinImpl != nil {
			instance[methodName] = methodInfo.BuiltinImpl
		}
	}

	return instance, nil
}

// ========================================
// FunctionBuilder - Para crear funciones standalone
// ========================================

// FunctionBuilder facilita la creación de funciones builtin
type FunctionBuilder struct {
	name        string
	params      []ast.Parameter
	returnType  *ast.Type
	impl        common.Func
	accessLevel string
	modifiers   []string
	fileName    string
	packageName string
}

// NewFunctionBuilder crea un nuevo builder para funciones
func NewFunctionBuilder(name string) *FunctionBuilder {
	return &FunctionBuilder{
		name:        name,
		accessLevel: "public",
		modifiers:   []string{"public"},
		fileName:    "builtin",
		packageName: "builtin",
	}
}

// SetParams establece los parámetros de la función
func (fb *FunctionBuilder) SetParams(params []ast.Parameter) *FunctionBuilder {
	fb.params = params
	return fb
}

// SetParamsFromNames establece los parámetros usando nombres y tipos
func (fb *FunctionBuilder) SetParamsFromNames(paramNames []string, paramTypes []*ast.Type) *FunctionBuilder {
	params := make([]ast.Parameter, len(paramNames))
	for i, name := range paramNames {
		var paramType *ast.Type
		if i < len(paramTypes) {
			paramType = paramTypes[i]
		}
		params[i] = ast.Parameter{Name: name, Type: paramType}
	}
	fb.params = params
	return fb
}

// SetReturnType establece el tipo de retorno
func (fb *FunctionBuilder) SetReturnType(returnType *ast.Type) *FunctionBuilder {
	fb.returnType = returnType
	return fb
}

// SetImplementation establece la implementación de la función
func (fb *FunctionBuilder) SetImplementation(impl common.Func) *FunctionBuilder {
	fb.impl = impl
	return fb
}

// SetAccessLevel establece el nivel de acceso
func (fb *FunctionBuilder) SetAccessLevel(accessLevel string) *FunctionBuilder {
	fb.accessLevel = accessLevel
	fb.modifiers = []string{accessLevel}
	return fb
}

// Build construye y registra la función en el entorno
func (fb *FunctionBuilder) Build(env *Env) (common.Func, error) {
	if fb.impl == nil {
		return nil, ThrowRuntimeError(env, "function builder requires an implementation")
	}

	// Wrap the implementation to validate parameters
	wrappedFunc := common.Func(func(callEnv *common.Env, args []any) (any, error) {
		// Validate parameter count
		requiredParams := 0
		for _, param := range fb.params {
			if !param.IsVariadic {
				requiredParams++
			}
		}

		if len(args) < requiredParams {
			return nil, ThrowArityError((*Env)(callEnv), requiredParams, len(args))
		}

		// Call the actual implementation
		return fb.impl(callEnv, args)
	})

	// Register in environment
	env.Set(fb.name, wrappedFunc)

	return wrappedFunc, nil
}

package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

type ClassDefinition = common.ClassDefinition
type InterfaceDefinition = common.InterfaceDefinition
type MethodInfo = common.MethodInfo
type FieldInfo = common.FieldInfo
type ConstructorInfo = common.ConstructorInfo
type ClassInstance = common.ClassInstance
type MethodSignature = common.MethodSignature

// isClassPermitted checks if a class is permitted by checking PermitNames
func isClassPermitted(permitNames []common.PrebuildedDefinition, className, classPackage, targetPackage string) bool {
	// If no permits list, only classes in the same package are permitted
	if len(permitNames) == 0 {
		return classPackage == targetPackage
	}

	// Check if class name matches any permit
	for _, permit := range permitNames {
		if permit.Name == className {
			return true
		}
	}
	return false
}

// Global registry for classes and interfaces, organized by package
var (
	classRegistry     = make(map[string]map[string]*ClassDefinition) // packageName -> className -> ClassDefinition
	interfaceRegistry = make(map[string]*InterfaceDefinition)
	builtinClasses    = make(map[string]*ClassDefinition) // builtin classes (always available)
	typeAliasRegistry = make(map[string]map[string]*TypeAlias) // packageName -> aliasName -> TypeAlias
)

// TypeAlias represents a type alias definition
type TypeAlias struct {
	Name      string
	BaseType  string
	IsFinal   bool // if true, this is a nominal type (distinct from base type)
	PackageName string
}

// lookupClass looks up a class by name, checking builtins first, then the given package
func lookupClass(className string, packageName string) (*ClassDefinition, bool) {
	// Check builtins first
	if classDef, exists := builtinClasses[className]; exists {
		return classDef, true
	}

	// Check package-scoped registry
	if packageClasses, exists := classRegistry[packageName]; exists {
		if classDef, exists := packageClasses[className]; exists {
			return classDef, true
		}
	}

	return nil, false
}

// ResetGlobalRegistries clears all global registries.
// This is primarily used for testing to ensure tests don't affect each other.
func ResetGlobalRegistries() {
	classRegistry = make(map[string]map[string]*ClassDefinition)
	interfaceRegistry = make(map[string]*InterfaceDefinition)
	builtinClasses = make(map[string]*ClassDefinition)
	typeAliasRegistry = make(map[string]map[string]*TypeAlias)
	// Also reset enum registry
	enumRegistry = make(map[string]*common.EnumDefinition)
	// Reset exception classes
	exceptionClasses = map[string]common.Func{}
}

// isClassAccessible checks if a class is accessible from the current file/package context
func isClassAccessible(classDef *ClassDefinition, currentFileName, currentPackageName string) bool {
	switch classDef.AccessLevel {
	case "public":
		return true // public classes are accessible everywhere
	case "private":
		return classDef.FileName == currentFileName // private classes only accessible in same file
	case "protected":
		// protected classes accessible in same package or same file
		return classDef.PackageName == currentPackageName || classDef.FileName == currentFileName
	default:
		return true // default to public if access level is not specified
	}
}

// getAccessibleClass retrieves a class if it's accessible from the current context
func getAccessibleClass(className string, env *Env) (*ClassDefinition, error) {
	// First, check if it's a builtin class (always available)
	if classDef, exists := builtinClasses[className]; exists {
		return classDef, nil
	}

	// Check if the class has been imported into this environment
	packageName, imported := env.ImportedClasses[className]
	if !imported {
		// Class not imported, check if it's in the current package
		packageName = env.GetPackageName()
	}

	// Look up in the package-scoped registry
	packageClasses, packageExists := classRegistry[packageName]
	if !packageExists {
		return nil, ThrowNameError(env, className)
	}

	classDef, exists := packageClasses[className]
	if !exists {
		return nil, ThrowNameError(env, className)
	}

	// Check accessibility
	if !isClassAccessible(classDef, env.GetFileName(), env.GetPackageName()) {
		switch classDef.AccessLevel {
		case "private":
			return nil, ThrowAccessError(env, className, "class", "private")
		case "protected":
			return nil, ThrowAccessError(env, className, "class", "protected")
		default:
			return nil, ThrowAccessError(env, className, "class", "unknown")
		}
	}

	return classDef, nil
}

// evalClassDecl handles class declaration evaluation
func evalClassDecl(env *Env, s *ast.ClassDecl) (any, error) {
	packageName := env.GetPackageName()

	// Initialize package map if it doesn't exist
	if classRegistry[packageName] == nil {
		classRegistry[packageName] = make(map[string]*ClassDefinition)
	}

	// Check if class already exists in this package
	if _, exists := classRegistry[packageName][s.Name]; exists {
		return nil, ThrowRuntimeError(env, fmt.Sprintf("class '%s' is already defined in package '%s'", s.Name, packageName))
	}

	// Parse parent class if specified
	var parentClass *ClassDefinition
	if s.Parent != "" {
		var err error
		parentClass, err = getAccessibleClass(s.Parent, env)
		if err != nil {
			return nil, err
		}
	}

	// Validate implemented interfaces and resolve to definitions
	var implementedInterfaces []*common.InterfaceDefinition
	for _, interfaceName := range s.Implements {
		interfaceDef, exists := interfaceRegistry[interfaceName]
		if !exists {
			return nil, ThrowNameError(env, interfaceName)
		}

		// If this interface is sealed, ensure this class is permitted to implement it
		if interfaceDef.IsSealed {
			allowed := isClassPermitted(interfaceDef.PermitNames, s.Name, env.GetPackageName(), interfaceDef.PackageName)
			if !allowed {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("interface %s is sealed and does not permit %s to implement it", interfaceName, s.Name))
			}
		}

		implementedInterfaces = append(implementedInterfaces, interfaceDef)
	}

	// If this class extends a sealed parent, ensure it's permitted
	if parentClass != nil && parentClass.IsSealed {
		allowed := isClassPermitted(parentClass.PermitNames, s.Name, env.GetPackageName(), parentClass.PackageName)
		if !allowed {
			return nil, ThrowRuntimeError(env, fmt.Sprintf("class %s is sealed and does not permit %s to inherit", parentClass.Name, s.Name))
		}
	}

	// Convert AST TypeParams to common.GenericType
	var typeParams []common.GenericType
	for _, tp := range s.TypeParams {
		// Create a GenericBound from the TypeParam
		var bounds []common.GenericBound
		
		// Create the primary bound from the type parameter name
		bound := common.GenericBound{
			Name:       ast.Type{Name: tp.Name},
			Variance:   tp.Variance,
			IsVariadic: tp.IsVariadic,
		}
		
		// If there are bounds (extends constraints), resolve them now
		// tp.Bounds contains the names of the types that this type parameter extends
		if len(tp.Bounds) > 0 {
			// For now, we'll try to resolve the first bound (the extends constraint)
			extendsName := tp.Bounds[0]
			if extVal, ok := env.Get(extendsName); ok {
				if classConst, ok := extVal.(*common.ClassConstructor); ok {
					bound.Extends = classConst.Definition
				} else if classDef, ok := extVal.(*ClassDefinition); ok {
					bound.Extends = classDef
				}
			}
		}
		
		bounds = append(bounds, bound)
		
		typeParams = append(typeParams, common.GenericType{
			Bounds: bounds,
		})
	}

	// Store permits as PrebuildedDefinition for lazy resolution
	var permitNames []common.PrebuildedDefinition
	for _, permitName := range s.Permits {
		permitNames = append(permitNames, common.PrebuildedDefinition{
			Name:        permitName,
			PackageName: env.GetPackageName(), // Assume same package unless qualified
		})
	}

	// Create class definition
	classDef := &ClassDefinition{
		Name:         s.Name,
		Type:         &ast.Type{Name: s.Name, IsBuiltin: false, IsClass: true},
		Parent:       parentClass,
		Implements:   implementedInterfaces,
		IsAbstract:   s.IsAbstract,
		AccessLevel:  s.AccessLevel,
		IsSealed:     s.IsSealed,
		Permits:      nil, // Resolved lazily when needed
		PermitNames:  permitNames,
		FileName:     env.GetFileName(),
		PackageName:  env.GetPackageName(),
		Fields:       make(map[string]FieldInfo),
		Methods:      make(map[string][]MethodInfo),
		StaticFields: make(map[string]any),
		TypeParams:   typeParams,
		IsGeneric:    len(typeParams) > 0,
		Constructors: []ConstructorInfo{},
	}

	// Process fields
	for _, field := range s.Fields {
		fieldInfo := FieldInfo{
			Name:      field.Name,
			Type:      field.Type,
			Modifiers: field.Modifiers,
			IsStatic:  contains(field.Modifiers, "static"),
			IsPrivate: contains(field.Modifiers, "private"),
		}

		// Evaluate initial value if provided
		if field.InitValue != nil {
			val, err := evalExpr(env, field.InitValue)
			if err != nil {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("error evaluating field '%s' initial value in class '%s': %v", field.Name, s.Name, err))
			}
			fieldInfo.InitValue = val
		}

		classDef.Fields[field.Name] = fieldInfo
	}

	// Process methods
	for _, method := range s.Methods {
		methodInfo := MethodInfo{
			Name:       method.Name,
			ReturnType: method.ReturnType,
			Body:       method.Body,
			Modifiers:  method.Modifiers,
			IsAbstract: method.IsAbstract,
			IsStatic:   contains(method.Modifiers, "static"),
			IsPrivate:  contains(method.Modifiers, "private"),
		}

		// Convert parameters
		for _, param := range method.Params {
			methodInfo.Params = append(methodInfo.Params, param)
		}

		classDef.Methods[method.Name] = append(classDef.Methods[method.Name], methodInfo)
	}

	// Process constructor (supports overloading)
	if s.Constructor != nil {
		constructorInfo := ConstructorInfo{
			Body: s.Constructor.Body,
		}

		// Convert parameters
		for _, param := range s.Constructor.Params {
			constructorInfo.Params = append(constructorInfo.Params, param)
		}

		classDef.Constructors = append(classDef.Constructors, constructorInfo)
	}

	// Register the class in the package
	classRegistry[packageName][s.Name] = classDef

	// Create constructor function
	// For user-defined generic classes, we handle type parameters
	constructorFunc := Func(func(callEnv *Env, args []any) (any, error) {
		// Extract type arguments if present (for generic classes)
		var typeArgs []string
		var constructorArgs []any

		if classDef.IsGeneric && len(args) > 0 {
			// Check if first arguments are type parameters
			numTypeParams := len(classDef.TypeParams)
			if len(args) >= numTypeParams {
				// Try to extract type arguments
				typeArgs = make([]string, 0, numTypeParams)
				for i := 0; i < numTypeParams && i < len(args); i++ {
					if typeInfo, ok := args[i].(map[string]any); ok {
						// Handle wildcard or variance-annotated types
						if isWildcard, ok := typeInfo["isWildcard"].(bool); ok {
							if isWildcard {
								kind, _ := typeInfo["kind"].(string)
								bound, _ := typeInfo["bound"].(string)
								variance, _ := typeInfo["variance"].(string)
								wildcardStr := formatWildcard(kind, bound)
								if variance != "" {
									wildcardStr = variance + " " + wildcardStr
								}
								typeArgs = append(typeArgs, wildcardStr)
							} else {
								// Variance-annotated type
								name, _ := typeInfo["name"].(string)
								variance, _ := typeInfo["variance"].(string)
								if variance != "" && name != "" {
									typeArgs = append(typeArgs, variance+" "+name)
								} else {
									typeArgs = append(typeArgs, name)
								}
							}
						} else {
							// Not a type info map, treat as regular arg
							break
						}
					} else if typeStr, ok := args[i].(string); ok && isTypeName(typeStr) {
						// Regular type name
						typeArgs = append(typeArgs, typeStr)
					} else {
						// Not a type argument, stop extraction
						break
					}
				}

				// If we extracted the right number of type params, use them
				if len(typeArgs) == numTypeParams {
					constructorArgs = args[numTypeParams:]
				} else {
					// Couldn't extract type params, treat all as constructor args
					typeArgs = nil
					constructorArgs = args
				}
			} else {
				constructorArgs = args
			}
		} else {
			constructorArgs = args
		}

		// Create the instance
		instance, err := createClassInstance(classDef, callEnv, constructorArgs)
		if err != nil {
			return nil, err
		}

		// Store type arguments if available
		if classInst, ok := instance.(*ClassInstance); ok && len(typeArgs) > 0 {
			classInst.Fields["__type_args__"] = typeArgs

			// Also create a map for easy lookup
			typeMap := make(map[string]string)
			for i, typeParam := range classDef.TypeParams {
				if i < len(typeArgs) && len(typeParam.Bounds) > 0 {
					// Use the first bound's name as the type parameter name
					typeMap[typeParam.Bounds[0].Name.Name] = typeArgs[i]
				}
			}
			classInst.Fields["__generic_types__"] = typeMap

			// Store variance information for runtime checking
			varianceMap := make(map[string]string)
			for _, typeParam := range classDef.TypeParams {
				if len(typeParam.Bounds) > 0 && typeParam.Bounds[0].Variance != "" {
					varianceMap[typeParam.Bounds[0].Name.Name] = typeParam.Bounds[0].Variance
				}
			}
			if len(varianceMap) > 0 {
				classInst.Fields["__variance__"] = varianceMap
			}
		}

		return instance, nil
	})

	// Wrap constructor with ClassConstructor for proper type identification
	constructor := &common.ClassConstructor{
		Definition: classDef,
		Func:       constructorFunc,
	}

	// Set class in environment
	env.Set(s.Name, constructor)

	return constructor, nil
}

// createClassInstance creates a new instance of a class
func createClassInstance(classDef *ClassDefinition, env *Env, args []any) (any, error) {
	// Check if trying to instantiate abstract class
	if classDef.IsAbstract {
		return nil, ThrowTypeError(env, "concrete class", fmt.Sprintf("abstract class '%s'", classDef.Name))
	}

	// Create instance
	instance := &ClassInstance{
		ClassName:   classDef.Name,
		Fields:      make(map[string]any),
		Methods:     make(map[string]Func),
		ParentClass: classDef,
	}

	// Initialize fields from class hierarchy
	if err := initializeFields(instance, classDef); err != nil {
		return nil, err
	}

	// Bind methods
	if err := bindMethods(instance, classDef, env); err != nil {
		return nil, err
	}

	// Call constructor if exists (with overload resolution)
	if len(classDef.Constructors) > 0 {
		// Select appropriate constructor based on argument count
		constructor := common.SelectConstructorOverload(classDef.Constructors, len(args))
		if constructor == nil {
			// No matching constructor found - compute available arities
			available := make([]int, len(classDef.Constructors))
			for i, c := range classDef.Constructors {
				available[i] = len(c.Params)
			}
			return nil, ThrowRuntimeError(env, fmt.Sprintf("no constructor found for %s with %d arguments (available: %v)", classDef.Name, len(args), available))
		}

		// Create constructor environment
		constructorEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
		constructorEnv.Set("this", instance)

		// Add super function if there's a parent class
		if classDef.Parent != nil {
			// Create super object with method access
			superObj := createSuperObject(instance, classDef.Parent, constructorEnv)

			// For backward compatibility, make super callable directly for constructor
			constructorEnv.Set("super", Func(func(callEnv *Env, superArgs []any) (any, error) {
				return callParentConstructor(instance, classDef.Parent, callEnv, superArgs)
			}))

			// Also set super as an object for method access (super.methodName)
			constructorEnv.Set("Super", superObj)
		}

		// Bind constructor parameters and validate arity/types centrally
		if err := bindParametersWithVariadic(constructorEnv, constructor.Params, args); err != nil {
			return nil, err
		}

		// Execute constructor body
		if constructor.BuiltinImpl != nil {
			// Use builtin implementation
			_, err := constructor.BuiltinImpl(constructorEnv, args)
			if err != nil {
				return nil, err
			}
		} else {
			// Execute Polyloft constructor body
			for _, stmt := range constructor.Body {
				_, ret, err := evalStmt(constructorEnv, stmt)
				if err != nil {
					return nil, err
				}
				if ret {
					break
				}
			}
		}
	}

	return instance, nil
}

// callParentConstructor calls the constructor of the parent class
func callParentConstructor(instance *ClassInstance, parentClass *ClassDefinition, env *Env, args []any) (any, error) {
	if len(parentClass.Constructors) == 0 {
		// If parent has no constructor, just return nil
		return nil, nil
	}

	// Select appropriate constructor based on argument count
	constructor := common.SelectConstructorOverload(parentClass.Constructors, len(args))
	if constructor == nil {
		// No matching constructor found
		available := make([]int, len(parentClass.Constructors))
		for i, c := range parentClass.Constructors {
			available[i] = len(c.Params)
		}
		return nil, ThrowRuntimeError(env, fmt.Sprintf("no constructor found for parent %s with %d arguments (available: %v)", parentClass.Name, len(args), available))
	}

	// Create parent constructor environment
	parentEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
	parentEnv.Set("this", instance)

	// Add super function recursively if parent has a parent
	if parentClass.Parent != nil {
		parentEnv.Set("super", Func(func(callEnv *Env, superArgs []any) (any, error) {
			return callParentConstructor(instance, parentClass.Parent, callEnv, superArgs)
		}))
	}

	// Bind parent constructor parameters and validate centrally
	if err := bindParametersWithVariadic(parentEnv, constructor.Params, args); err != nil {
		return nil, err
	}

	// Execute parent constructor body or builtin implementation
	if constructor.BuiltinImpl != nil {
		if _, err := constructor.BuiltinImpl(parentEnv, args); err != nil {
			return nil, err
		}
	} else {
		for _, stmt := range constructor.Body {
			_, ret, err := evalStmt(parentEnv, stmt)
			if err != nil {
				return nil, err
			}
			if ret {
				break
			}
		}
	}

	return nil, nil
}

// callParentMethod calls a specific method from the parent class
func callParentMethod(instance *ClassInstance, parentClass *ClassDefinition, methodInfo MethodInfo, env *Env, args []any) (any, error) {
	// Create parent method environment
	parentMethodEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
	parentMethodEnv.Set("this", instance)

	// Add super recursively if parent has a parent
	if parentClass.Parent != nil {
		parentMethodEnv.Set("super", createSuperObject(instance, parentClass.Parent, parentMethodEnv))
	}

	// Bind method parameters (validates arity/types)
	if err := bindParametersWithVariadic(parentMethodEnv, methodInfo.Params, args); err != nil {
		return nil, err
	}

	// Execute parent method body
	var lastValue any
	for _, stmt := range methodInfo.Body {
		val, ret, err := evalStmt(parentMethodEnv, stmt)
		if err != nil {
			return nil, err
		}
		if ret {
			return val, nil
		}
		lastValue = val
	}

	return lastValue, nil
}

// createSuperObject creates a super object that provides access to parent methods
func createSuperObject(instance *ClassInstance, parentClass *ClassDefinition, env *Env) map[string]any {
	superObj := map[string]any{}

	// Add constructor call capability
	superObj["__call__"] = Func(func(callEnv *Env, superArgs []any) (any, error) {
		return callParentConstructor(instance, parentClass, callEnv, superArgs)
	})

	// Add all parent methods (including inherited from grandparents)
	addParentMethods(superObj, instance, parentClass, env)

	return superObj
}

// addParentMethods recursively adds methods from parent classes
func addParentMethods(superObj map[string]any, instance *ClassInstance, classDef *ClassDefinition, env *Env) {
	if classDef == nil {
		return
	}

	// Add methods from parent classes first (so they can be overridden)
	if classDef.Parent != nil {
		addParentMethods(superObj, instance, classDef.Parent, env)
	}

	// Add methods from this class level
	for methodName, methodOverloads := range classDef.Methods {
		// For each method name, create a wrapper that selects the right overload
		overloads := methodOverloads // Copy for closure
		superObj[methodName] = Func(func(callEnv *Env, args []any) (any, error) {
			// Select appropriate method based on argument count
			method := common.SelectMethodOverload(overloads, len(args))
			if method == nil {
				return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no overload found for super.%s with %d arguments", methodName, len(args)))
			}
			if !method.IsStatic {
				return callParentMethod(instance, classDef, *method, callEnv, args)
			}
			return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("cannot call static method %s via super", methodName))
		})
	}

	// Add methods from implemented interfaces with default implementations
	for _, interfaceDef := range classDef.Implements {
		for methodName, signatures := range interfaceDef.Methods {
			sigs := signatures // Copy for closure
			superObj[methodName] = Func(func(callEnv *Env, args []any) (any, error) {
				// Find matching signature based on argument count
				var matchingSig *MethodSignature
				for i := range sigs {
					if len(sigs[i].Params) == len(args) {
						matchingSig = &sigs[i]
						break
					}
				}
				if matchingSig == nil || !matchingSig.HasDefault || matchingSig.DefaultBody == nil {
					return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no default implementation for interface method %s with %d arguments", methodName, len(args)))
				}
				return callDefaultInterfaceMethod(instance, *matchingSig, callEnv, args)
			})
		}
	}
}

// initializeFields initializes instance fields from the class hierarchy
func initializeFields(instance *ClassInstance, classDef *ClassDefinition) error {
	// Initialize parent fields first
	if classDef.Parent != nil {
		if err := initializeFields(instance, classDef.Parent); err != nil {
			return err
		}
	}

	// Initialize this class's fields
	for name, fieldInfo := range classDef.Fields {
		if !fieldInfo.IsStatic {
			instance.Fields[name] = fieldInfo.InitValue
		}
	}

	return nil
}

// bindMethods binds instance methods from the class hierarchy
func bindMethods(instance *ClassInstance, classDef *ClassDefinition, env *Env) error {
	if instance == nil {
		return fmt.Errorf("bindMethods: instance is nil")
	}
	if classDef == nil {
		return fmt.Errorf("bindMethods: classDef is nil")
	}

	// Bind parent methods first
	if classDef.Parent != nil {
		if err := bindMethods(instance, classDef.Parent, env); err != nil {
			return err
		}
	}

	// Bind this class's methods (with overload resolution)
	for name, methodOverloads := range classDef.Methods {
		overloads := methodOverloads // Copy for closure
		method := Func(func(callEnv *Env, args []any) (any, error) {
			// Select appropriate method based on argument count
			selectedMethod := common.SelectMethodOverload(overloads, len(args))
			if selectedMethod == nil {
				return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no overload found for %s.%s with %d arguments", instance.ClassName, name, len(args)))
			}
			if selectedMethod.IsStatic || selectedMethod.IsAbstract {
				return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("cannot call static or abstract method %s via instance", name))
			}
			return CallInstanceMethod(instance, *selectedMethod, callEnv, args)
		})
		instance.Methods[name] = method
	}

	// Bind default interface methods for implemented interfaces
	for _, interfaceDef := range classDef.Implements {
		if interfaceDef == nil {
			continue
		}
		for methodName, signatures := range interfaceDef.Methods {
			// Only bind default methods that aren't already implemented
			if _, methodExists := instance.Methods[methodName]; !methodExists {
				sigs := signatures // Copy for closure
				method := Func(func(callEnv *Env, args []any) (any, error) {
					// Find matching signature based on argument count
					var matchingSig *MethodSignature
					for i := range sigs {
						if sigs[i].HasDefault && sigs[i].DefaultBody != nil && len(sigs[i].Params) == len(args) {
							matchingSig = &sigs[i]
							break
						}
					}
					if matchingSig == nil {
						return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no default implementation for interface method %s with %d arguments", methodName, len(args)))
					}
					return callDefaultInterfaceMethod(instance, *matchingSig, callEnv, args)
				})
				instance.Methods[methodName] = method
			}
		}
	}

	// Add default toString method if not already present
	if _, exists := instance.Methods["toString"]; !exists {
		instance.Methods["toString"] = Func(func(callEnv *Env, args []any) (any, error) {
			return fmt.Sprintf("%s@%p", instance.ClassName, instance), nil
		})
	}

	return nil
}

// validateReturnType validates that a return value matches the expected type
func validateReturnType(instance *ClassInstance, expectedTypeName string, value any, env *Env) error {

	// Check if it's a generic type parameter (single uppercase letter or T-prefixed)
	if isGenericTypeParameter(expectedTypeName) {
		// Generic type parameters are not validated at runtime
		return nil
	}

	// Strip variance annotations before validating concrete types
	expectedTypeName = stripVarianceAnnotation(expectedTypeName)

	// It's a concrete type, validate it
	return validateConcreteType(expectedTypeName, value, env)
}

// stripVarianceAnnotation removes "out " or "in " prefix from a type name
func stripVarianceAnnotation(typeName string) string {
	if strings.HasPrefix(typeName, "out ") {
		return strings.TrimPrefix(typeName, "out ")
	}
	if strings.HasPrefix(typeName, "in ") {
		return strings.TrimPrefix(typeName, "in ")
	}
	return typeName
}

// validateConcreteType validates a value against a concrete type name
func validateConcreteType(typeName string, value any, env *Env) error {
	if value == nil {
		return nil
	}

	if !IsInstanceOf(value, typeName) {
		actualType := GetTypeName(value)
		return ThrowTypeError(env, typeName, actualType)
	}
	return nil
}

// CallInstanceMethod calls a method on an instance
func CallInstanceMethod(instance *ClassInstance, methodInfo MethodInfo, env *Env, args []any) (any, error) {
	// Create method environment
	methodEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
	methodEnv.Set("this", instance)

	// Add super access if the instance has a parent class
	if instance.ParentClass != nil && instance.ParentClass.Parent != nil {
		superObj := createSuperObject(instance, instance.ParentClass.Parent, methodEnv)
		methodEnv.Set("super", superObj)
	}

	// Bind method parameters and validate centrally
	if err := bindParametersWithVariadic(methodEnv, methodInfo.Params, args); err != nil {
		return nil, err
	}

	// Execute method body
	var result any
	var err error

	if methodInfo.BuiltinImpl != nil {
		// Use builtin implementation
		result, err = methodInfo.BuiltinImpl(methodEnv, args)
		if err != nil {
			return nil, err
		}
	} else {
		// Execute Polyloft method body
		var lastValue any
		for _, stmt := range methodInfo.Body {
			val, ret, errStmt := evalStmt(methodEnv, stmt)
			if errStmt != nil {
				return nil, errStmt
			}
			if ret {
				result = val
				break
			}
			lastValue = val
		}
		if result == nil {
			result = lastValue
		}
	}

	returnTypeName := ast.GetTypeNameString(methodInfo.ReturnType)
	if returnTypeName != "" && returnTypeName != "Any" && returnTypeName != "Void" {
		if err := validateReturnType(instance, returnTypeName, result, methodEnv); err != nil {
			return nil, err
		}
	}

	// For Void return type, don't return the last value (return nil instead)
	if returnTypeName == "Void" {
		return nil, nil
	}

	return result, nil
}

// callDefaultInterfaceMethod calls a default method from an interface
func callDefaultInterfaceMethod(instance *ClassInstance, signature MethodSignature, env *Env, args []any) (any, error) {
	// Create method environment
	methodEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}
	methodEnv.Set("this", instance)

	// Bind method parameters (handles variadic/validation)
	if err := bindParametersWithVariadic(methodEnv, signature.Params, args); err != nil {
		return nil, err
	}

	// Execute default method body
	var lastValue any
	for _, stmt := range signature.DefaultBody {
		val, ret, err := evalStmt(methodEnv, stmt)
		if err != nil {
			return nil, err
		}
		if ret {
			return val, nil
		}
		lastValue = val
	}

	return lastValue, nil
}

// evalInterfaceDecl handles interface declaration evaluation
func evalInterfaceDecl(env *Env, s *ast.InterfaceDecl) (any, error) {
	// Check if interface already exists
	if _, exists := interfaceRegistry[s.Name]; exists {
		return nil, ThrowRuntimeError(env, fmt.Sprintf("interface '%s' is already defined", s.Name))
	}

	// Store permits as PrebuildedDefinition for lazy resolution
	var permitNames []common.PrebuildedDefinition
	for _, permitName := range s.Permits {
		permitNames = append(permitNames, common.PrebuildedDefinition{
			Name:        permitName,
			PackageName: env.GetPackageName(), // Assume same package unless qualified
		})
	}

	interfaceDef := &InterfaceDefinition{
		Name:         s.Name,
		Type:         &ast.Type{Name: s.Name, IsBuiltin: false, IsInterface: true},
		Methods:      make(map[string][]MethodSignature),
		IsSealed:     s.IsSealed,
		Permits:      nil, // Resolved lazily when needed
		PermitNames:  permitNames,
		StaticFields: make(map[string]any),
		AccessLevel:  s.AccessLevel,
		FileName:     env.GetFileName(),
		PackageName:  env.GetPackageName(),
	}

	// Process method signatures (with overload support)
	for _, method := range s.Methods {
		signature := MethodSignature{
			Name:        method.Name,
			ReturnType:  method.ReturnType,
			HasDefault:  method.HasDefault,
			DefaultBody: method.DefaultBody,
		}

		// Convert parameters
		for _, param := range method.Params {
			signature.Params = append(signature.Params, param)
		}

		interfaceDef.Methods[method.Name] = append(interfaceDef.Methods[method.Name], signature)
	}

	// Process static fields
	for _, field := range s.Fields {
		// Initialize field value
		var fieldValue any
		if field.InitValue != nil {
			val, err := evalExpr(env, field.InitValue)
			if err != nil {
				return nil, err
			}
			fieldValue = val
		}
		interfaceDef.StaticFields[field.Name] = fieldValue
	}

	// Register the interface
	interfaceRegistry[s.Name] = interfaceDef

	return nil, nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetClassInstance safely casts a value to a ClassInstance
func GetClassInstance(val any) (*ClassInstance, bool) {
	instance, ok := val.(*ClassInstance)
	return instance, ok
}

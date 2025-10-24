package engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"reflect"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// Package engine hosts evaluation for Polyloft.
// Options control execution behavior (flags, limits, debug hooks, etc.).

// rangeMarker is an internal type to mark range expressions in index context
type rangeMarker struct {
	start int
	end   int
}

// bindParametersWithVariadic binds function parameters to arguments, handling variadic parameters
// isGenericTypeParameter checks if a type name is a generic type parameter or Any
// Generic type parameters can be:
// - Single uppercase letters: T, K, V, E, etc.
// - T followed by uppercase letter: TKey, TValue, TElement, etc.
// - Any (pass-through type with no validation)
func isGenericTypeParameter(typeName string) bool {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return false
	}
	// Any is treated as a pass-through type (no validation)
	if strings.EqualFold(typeName, "any") {
		return true
	}
	// Check if it's a single uppercase letter (T, K, V, E, etc.)
	if len(typeName) == 1 && typeName[0] >= 'A' && typeName[0] <= 'Z' {
		return true
	}
	// Check if it starts with T followed by uppercase (TKey, TValue, TElement, etc.)
	if len(typeName) > 1 && typeName[0] == 'T' && typeName[1] >= 'A' && typeName[1] <= 'Z' {
		return true
	}
	return false
}

func bindParametersWithVariadic(env *common.Env, params []ast.Parameter, args []any) error {
	// Check minimum required parameters (non-variadic)
	requiredParams := 0
	variadicParam := -1
	var variadicType string

	for i, param := range params {
		if param.IsVariadic {
			variadicParam = i
			variadicType = ast.GetTypeNameString(param.Type)
			break
		}
		requiredParams++
	}

	// Check if we have enough arguments for required parameters
	if len(args) < requiredParams {
		return ThrowArityError((*Env)(env), requiredParams, len(args))
	}

	// Get generic type mappings from 'this' if available (for method calls on generic class instances)
	var genericTypes map[string]string
	var varianceMap map[string]string // Maps type parameter name to variance ("in", "out", or "")
	if thisVal, ok := env.Get("this"); ok {
		if classInst, ok := thisVal.(*common.ClassInstance); ok {
			if typeMap, ok := classInst.Fields["__generic_types__"].(map[string]string); ok {
				genericTypes = typeMap
			}
			// Get variance information from __variance__ field if available
			if varMap, ok := classInst.Fields["__variance__"].(map[string]string); ok {
				varianceMap = varMap
			}
		}
	}

	// Bind regular parameters with type validation
	for i := 0; i < requiredParams; i++ {
		paramTypeName := ast.GetTypeNameString(params[i].Type)

		// Check variance constraints - covariant (out) type parameters cannot appear in parameter positions
		if paramTypeName != "" && isGenericTypeParameter(paramTypeName) && varianceMap != nil {
			if variance, found := varianceMap[paramTypeName]; found && variance == "out" {
				return ThrowRuntimeError((*Env)(env), fmt.Sprintf(
					"covariant type parameter '%s' (declared with 'out') cannot be used in parameter position - it can only appear in return types",
					paramTypeName))
			}
		}

		// Resolve generic type parameter to concrete type if possible
		resolvedType := paramTypeName
		if paramTypeName != "" && isGenericTypeParameter(paramTypeName) && genericTypes != nil {
			if concreteType, found := genericTypes[paramTypeName]; found {
				resolvedType = concreteType
			}
		}

		// Validate type if we have a concrete type (not a generic parameter)
		// Skip validation for wildcards as they have special semantics
		if resolvedType != "" && !isGenericTypeParameter(resolvedType) && !isWildcardType(resolvedType) {
			if err := ValidateArgumentType(args[i], resolvedType); err != nil {
				return err
			}
		}
		env.Set(params[i].Name, args[i])
	}

	// Handle variadic parameter if present
	if variadicParam >= 0 {
		variadicArgs := args[requiredParams:]

		// Check variance constraints for variadic parameters
		if variadicType != "" && isGenericTypeParameter(variadicType) && varianceMap != nil {
			if variance, found := varianceMap[variadicType]; found && variance == "out" {
				return ThrowRuntimeError((*Env)(env), fmt.Sprintf(
					"covariant type parameter '%s' (declared with 'out') cannot be used in parameter position - it can only appear in return types",
					variadicType))
			}
		}

		// Resolve generic type parameter for variadic type
		resolvedVariadicType := variadicType
		if variadicType != "" && isGenericTypeParameter(variadicType) && genericTypes != nil {
			if concreteType, found := genericTypes[variadicType]; found {
				resolvedVariadicType = concreteType
			}
		}

		// Validate variadic arguments if we have a concrete type (not wildcard)
		if resolvedVariadicType != "" && !isGenericTypeParameter(resolvedVariadicType) && !isWildcardType(resolvedVariadicType) {
			// Validate and convert variadic arguments to array of the specified type
			validatedArgs, err := ValidateVariadicArguments(variadicArgs, resolvedVariadicType)
			if err != nil {
				return err
			}
			env.Set(params[variadicParam].Name, validatedArgs)
		} else {
			// For unresolved generic variadic types or wildcards, just bind the args as-is
			env.Set(params[variadicParam].Name, variadicArgs)
		}
	} else if len(args) > requiredParams {
		return ThrowArityError((*Env)(env), requiredParams, len(args))
	}

	return nil
}

// isWildcardType checks if a type name represents a wildcard type or variance-annotated type
func isWildcardType(typeName string) bool {
	typeName = strings.TrimSpace(typeName)
	// Check for wildcard types (?, ? extends T, ? super T)
	if strings.HasPrefix(typeName, "?") {
		return true
	}
	// Check for variance annotations (in T, out T)
	if strings.HasPrefix(typeName, "in ") || strings.HasPrefix(typeName, "out ") {
		return true
	}
	return false
}

func Eval(prog *ast.Program, opts Options) (any, error) {
	return EvalWithContext(prog, opts, "", "")
}

func EvalWithContext(prog *ast.Program, opts Options, fileName, packageName string) (any, error) {
	return EvalWithContextAndSource(prog, opts, fileName, packageName, "")
}

func EvalWithContextAndSource(prog *ast.Program, opts Options, fileName, packageName, source string) (any, error) {
	var env *common.Env
	if fileName != "" {
		env = common.NewEnvWithContext(fileName, packageName)
	} else {
		env = common.NewEnv()
	}

	// Store source lines in env for context tracking and hint generation
	if source != "" {
		sourceLines := strings.Split(source, "\n")
		env.SetSourceLines(sourceLines)
	}

	installBuiltins(env, opts)
	InstallSysModule(env, opts)         // Install enhanced Sys module
	InstallMathModule(env)              // Install Math module
	InstallExceptionBuiltins(env)       // Install exception system
	InitFunctionInterfaces((*Env)(env)) // Install Function and BiFunction interfaces
	// module loader bound names
	env.Set("import", common.Func(func(e *common.Env, args []any) (any, error) {
		return nil, ThrowRuntimeError((*Env)(e), "use 'import' statement, not function")
	}))

	// Set file environment variables (like Python's __name__ but with $ prefix)
	// These variables are available in the main execution context
	if fileName != "" {
		env.Set("$name", filepath.Base(fileName))                                                            // e.g., "main.pf"
		env.Set("$file", fileName)                                                                           // e.g., "src/main.pf"
		env.Set("$package", packageName)                                                                     // e.g., "src"
		env.Set("$stem", strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(filepath.Base(fileName)))) // e.g., "main"
	}

	var last any
	for _, st := range prog.Stmts {
		v, ret, err := evalStmtWithSource(env, st, env.GetSourceLines())
		if err != nil {
			return nil, err
		}
		if ret {
			return v, nil
		}
		last = v
	}
	return last, nil
}

// evalStmtWithSource evaluates a statement with source context for better error messages
func evalStmtWithSource(env *common.Env, st ast.Stmt, sourceLines []string) (val any, returned bool, err error) {
	return evalStmt(env, st)
}

func evalStmt(env *common.Env, st ast.Stmt) (val any, returned bool, err error) {
	switch s := st.(type) {
	case *ast.ImportStmt:
		err := handleImport(env, s)
		return nil, false, err
	case *ast.IfStmt:
		for _, c := range s.Clauses {
			v, err := evalExpr(env, c.Cond)
			if err != nil {
				return nil, false, err
			}
			if utils.AsBool(v) {
				for _, sub := range c.Body {
					v, ret, err := evalStmt(env, sub)
					if err != nil {
						return nil, false, err
					}
					// Check for break/continue sentinels
					switch v.(type) {
					case breakSentinel:
						return v, false, nil
					case continueSentinel:
						return v, false, nil
					}
					if ret {
						return v, true, nil
					}
				}
				return nil, false, nil
			}
		}
		for _, sub := range s.Else {
			v, ret, err := evalStmt(env, sub)
			if err != nil {
				return nil, false, err
			}
			// Check for break/continue sentinels
			switch v.(type) {
			case breakSentinel:
				return v, false, nil
			case continueSentinel:
				return v, false, nil
			}
			if ret {
				return v, true, nil
			}
		}
		return nil, false, nil
	case *ast.ForInStmt:
		it, err := evalExpr(env, s.Iterable)
		if err != nil {
			return nil, false, err
		}

		// Determine if we're using destructuring (multiple iteration variables)
		useDestructuring := len(s.Names) > 1

		// Handle plain Go maps (map[string]any)
		if mapVal, ok := it.(map[string]any); ok {
			// Iterate over map entries
			for key, value := range mapVal {
				if useDestructuring {
					// Set key and value to the two variables
					if len(s.Names) >= 2 {
						env.Set(s.Names[0], key)
						env.Set(s.Names[1], value)
					}
					// Set any additional variables to nil
					for i := 2; i < len(s.Names); i++ {
						env.Set(s.Names[i], nil)
					}
				} else {
					// Single variable: set it to the key
					varName := s.Name
					if len(s.Names) > 0 {
						varName = s.Names[0]
					}
					env.Set(varName, key)
				}

				// Evaluate where clause if present
				if s.Where != nil {
					whereResult, err := evalExpr(env, s.Where)
					if err != nil {
						return nil, false, err
					}
					// Skip this iteration if where clause is false
					if !utils.AsBool(whereResult) {
						continue
					}
				}

				// Execute loop body
				brk, cont, ret, val, err := runBlock(env, s.Body)
				if err != nil {
					return nil, false, err
				}
				if ret {
					return val, true, nil
				}
				if brk {
					break
				}
				if cont {
					continue
				}
			}
			return nil, false, nil
		}

		// Check if the object implements Iterable interface
		if instance, ok := it.(*ClassInstance); ok {
			// Check if class implements Iterable
			iterableInterfaceDef := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
			if iterableInterfaceDef == nil {
				return nil, false, fmt.Errorf("Iterable interface not found")
			}

			// Check if instance implements Iterable interface
			implementsIterable := false
			if instance.ParentClass != nil {
				for _, interfaceDef := range instance.ParentClass.Implements {
					if interfaceDef == iterableInterfaceDef {
						implementsIterable = true
						break
					}
				}
			}

			if !implementsIterable {
				return nil, false, fmt.Errorf("object of type %s does not implement Iterable interface", instance.ClassName)
			}
			// New methods of iterable are __length and __get(index)
			__lengthMethod, ok := instance.Methods["__length"]
			if !ok {
				return nil, false, fmt.Errorf("Iterable object missing __length() method")
			}
			__lengthFunc, ok := common.ExtractFunc(__lengthMethod)
			if !ok {
				return nil, false, fmt.Errorf("__length is not a function")
			}
			__lengthResult, err := __lengthFunc(env, []any{})
			if err != nil {
				return nil, false, err
			}
			length, ok := utils.AsInt(__lengthResult)
			if !ok {
				return nil, false, fmt.Errorf("__length() must return an integer")
			}
			__getMethod, ok := instance.Methods["__get"]
			if !ok {
				return nil, false, fmt.Errorf("Iterable object missing __get() method")
			}
			__getFunc, ok := common.ExtractFunc(__getMethod)
			if !ok {
				return nil, false, fmt.Errorf("__get is not a function")
			}
			// Iterate from 0 to length-1
			for idx := 0; idx < length; idx++ {
				// Call __get(index) method
				el, err := __getFunc(env, []any{idx})
				if err != nil {
					return nil, false, err
				}

				// Handle destructuring using Unstructured interface if needed
				if useDestructuring {
					// Check if element implements Unstructured interface
					if elInstance, ok := el.(*ClassInstance); ok {
						unstructuredInterfaceDef := common.BuiltinInterfaceUnstructured.GetInterfaceDefinition(env)
						implementsUnstructured := false
						if unstructuredInterfaceDef != nil && elInstance.ParentClass != nil {
							for _, interfaceDef := range elInstance.ParentClass.Implements {
								if interfaceDef == unstructuredInterfaceDef {
									implementsUnstructured = true
									break
								}
							}
						}

						if implementsUnstructured {
							// Element implements Unstructured, use pieces() and getPiece()
							piecesMethod, ok := elInstance.Methods["__pieces"]
							if !ok {
								return nil, false, fmt.Errorf("Unstructured object missing pieces() method")
							}
							piecesFunc, ok := common.ExtractFunc(piecesMethod)
							if !ok {
								return nil, false, fmt.Errorf("pieces is not a function")
							}
							piecesResult, err := piecesFunc(env, []any{})
							if err != nil {
								return nil, false, err
							}
							numPieces, ok := utils.AsInt(piecesResult)
							if !ok {
								return nil, false, fmt.Errorf("pieces() must return an integer")
							}

							// Check if number of variables matches number of pieces
							if len(s.Names) != numPieces {
								return nil, false, fmt.Errorf("destructuring mismatch: expected %d variables, got %d pieces", len(s.Names), numPieces)
							}

							// Get each piece using getPiece(index)
							getPieceMethod, ok := elInstance.Methods["__get_piece"]
							if !ok {
								return nil, false, fmt.Errorf("Unstructured object missing getPiece() method")
							}
							getPieceFunc, ok := common.ExtractFunc(getPieceMethod)
							if !ok {
								return nil, false, fmt.Errorf("getPiece is not a function")
							}
							for i, name := range s.Names {
								piece, err := getPieceFunc(env, []any{i})
								if err != nil {
									return nil, false, err
								}
								env.Set(name, piece)
							}
						} else {
							// Element doesn't implement Unstructured, set first var to element, rest to nil
							for i, name := range s.Names {
								if i == 0 {
									env.Set(name, el)
								} else {
									env.Set(name, nil)
								}
							}
						}
					} else {
						// Element is not a ClassInstance - check if it's an array for destructuring
						if elArray, ok := el.([]any); ok {
							// Destructure the array into variables
							for i, name := range s.Names {
								if i < len(elArray) {
									env.Set(name, elArray[i])
								} else {
									env.Set(name, nil)
								}
							}
						} else {
							// Not an array, set first var to element, rest to nil
							for i, name := range s.Names {
								if i == 0 {
									env.Set(name, el)
								} else {
									env.Set(name, nil)
								}
							}
						}
					}
				} else {
					// Single variable: set it to the element
					varName := s.Name
					if len(s.Names) > 0 {
						varName = s.Names[0]
					}
					env.Set(varName, el)
				}

				// Evaluate where clause if present
				if s.Where != nil {
					whereResult, err := evalExpr(env, s.Where)
					if err != nil {
						return nil, false, err
					}
					// Skip this iteration if where clause is false
					if !utils.AsBool(whereResult) {
						continue
					}
				}

				// Execute loop body
				brk, cont, ret, val, err := runBlock(env, s.Body)
				if err != nil {
					return nil, false, err
				}
				if ret {
					return val, true, nil
				}
				if brk {
					break
				}
				if cont {
					continue
				}
			}
		} else {
			// Not a ClassInstance - can't be iterable
			return nil, false, ThrowTypeError(env, "iterable", it)
		}
		return nil, false, nil
	case *ast.LoopStmt:
		for {
			brk, cont, ret, val, err := runBlock(env, s.Body)
			if err != nil {
				return nil, false, err
			}
			if ret {
				return val, true, nil
			}
			if brk {
				break
			}
			if cont {
				continue
			}
		}
		return nil, false, nil
	case *ast.BreakStmt:
		return breakSentinel{}, false, nil
	case *ast.ContinueStmt:
		return continueSentinel{}, false, nil
	case *ast.LetStmt:
		v, err := evalExpr(env, s.Value)
		if err != nil {
			return nil, false, err
		}
		
		// Handle destructuring if multiple names are present
		if len(s.Names) > 1 {
			// Try to destructure the value
			var values []any
			
			// Check if value is a plain Go array first
			if arr, ok := v.([]any); ok {
				values = arr
			} else if instance, ok := v.(*ClassInstance); ok {
				// For Array ClassInstance, extract the underlying array
				if instance.ClassName == "Array" {
					if arrData, ok := instance.Fields["_items"].([]any); ok {
						values = arrData
					} else {
						return nil, false, fmt.Errorf("Array instance missing _items field")
					}
				} else {
					// Check if it implements Unstructured interface
					unstructuredInterfaceDef := common.BuiltinInterfaceUnstructured.GetInterfaceDefinition(env)
					implementsUnstructured := false
					if unstructuredInterfaceDef != nil && instance.ParentClass != nil {
						for _, interfaceDef := range instance.ParentClass.Implements {
							if interfaceDef == unstructuredInterfaceDef {
								implementsUnstructured = true
								break
							}
						}
					}
					
					if implementsUnstructured {
						// Use __pieces() and __getPiece(i) methods
						piecesMethod, ok := instance.Methods["__pieces"]
						if !ok {
							return nil, false, fmt.Errorf("Unstructured object missing __pieces() method")
						}
						piecesFunc, ok := common.ExtractFunc(piecesMethod)
						if !ok {
							return nil, false, fmt.Errorf("__pieces is not a function")
						}
						piecesResult, err := piecesFunc(env, []any{})
						if err != nil {
							return nil, false, err
						}
						numPieces, ok := utils.AsInt(piecesResult)
						if !ok {
							return nil, false, fmt.Errorf("__pieces() must return an integer")
						}
						
						getPieceMethod, ok := instance.Methods["__getPiece"]
						if !ok {
							return nil, false, fmt.Errorf("Unstructured object missing __getPiece() method")
						}
						getPieceFunc, ok := common.ExtractFunc(getPieceMethod)
						if !ok {
							return nil, false, fmt.Errorf("__getPiece is not a function")
						}
						
						values = make([]any, numPieces)
						for i := 0; i < numPieces; i++ {
							piece, err := getPieceFunc(env, []any{i})
							if err != nil {
								return nil, false, err
							}
							values[i] = piece
						}
					} else {
						// Can't destructure this object
						return nil, false, fmt.Errorf("cannot destructure value of type %s", instance.ClassName)
					}
				}
			} else {
				// Can't destructure non-array, non-Unstructured value
				return nil, false, fmt.Errorf("cannot destructure value of type %T", v)
			}
			
			// Assign values to variables
			for i, name := range s.Names {
				var val any
				if i < len(values) {
					val = values[i]
				} else {
					val = nil // Not enough values, assign nil
				}
				env.Define(name, val, s.Kind)
			}
			return v, false, nil
		}
		
		// Single variable assignment (backward compatible)
		env.Define(s.Name, v, s.Kind)
		return v, false, nil
	case *ast.AssignStmt:
		// Handle assignment statements like x = value or this.field = value
		value, err := evalExpr(env, s.Value)
		if err != nil {
			return nil, false, err
		}

		// Handle different types of assignment targets

		switch target := s.Target.(type) {
		case *ast.Ident:
			for cur := env; cur != nil; cur = cur.Parent {
				if _, ok := cur.Vars[target.Name]; ok {
					if cur.Finals[target.Name] {
						return nil, false, ThrowRuntimeError(env, fmt.Sprintf("cannot assign to final variable '%s'", target.Name))
					}
					if cur.Consts[target.Name] {
						return nil, false, ThrowRuntimeError(env, fmt.Sprintf("cannot assign to constant '%s'", target.Name))
					}
					cur.Vars[target.Name] = value
					return value, false, nil
				}
			}
			// Variable doesn't exist, create it
			env.Vars[target.Name] = value
		case *ast.FieldExpr:
			// Field assignment: obj.field = value
			// First check if this is a static field assignment (ClassName.field or InterfaceName.field)
			if ident, ok := target.X.(*ast.Ident); ok {
				// Check if it's a class static field assignment
				if classDef, exists := lookupClass(ident.Name, env.GetPackageName()); exists {
					if _, fieldExists := classDef.StaticFields[target.Name]; fieldExists {
						classDef.StaticFields[target.Name] = value
						return value, false, nil
					}
				}
				// Check if it's an interface static field assignment
				if interfaceDef, exists := interfaceRegistry[ident.Name]; exists {
					if _, fieldExists := interfaceDef.StaticFields[target.Name]; fieldExists {
						interfaceDef.StaticFields[target.Name] = value
						return value, false, nil
					}
				}
			}

			// Regular instance field assignment
			obj, err := evalExpr(env, target.X)
			if err != nil {
				return nil, false, err
			}
			cur := env
			if instance, ok := obj.(*ClassInstance); ok {
				// Special handling for Map instances - set data in _data map
				if instance.ClassName == "Map" {
					if hashData, ok := instance.Fields["_data"].(map[uint64][]*mapEntry); ok {
						hash := hashValue(env, target.Name)
						// Check if hash already exists
						if entries, exists := hashData[hash]; exists {
							// Look for existing key
							found := false
							for i, entry := range entries {
								if equals(entry.Key, target.Name) {
									hashData[hash][i] = &mapEntry{Key: target.Name, Value: value}
									found = true
									break
								}
							}
							if !found {
								hashData[hash] = append(entries, &mapEntry{Key: target.Name, Value: value})
							}
						} else {
							hashData[hash] = []*mapEntry{{Key: target.Name, Value: value}}
						}
						return value, false, nil
					}
				}

				//cur.Parent.Vars get any of Final and check if is a ClassInstance and if its the same class
				if instance.Fields[target.Name] == nil {
					instance.Fields[target.Name] = value
				} else {
					for k := range cur.Parent.Finals {
						c, cok := cur.Parent.Vars[k].(*ClassInstance)
						if cok && c == instance {
							if isFinal, exists := cur.Parent.Finals[k]; exists && isFinal {
								return nil, false, ThrowRuntimeError(cur.Parent, fmt.Sprintf("cannot modify to final object '%s'", k))
							}
						}
					}
					instance.Fields[target.Name] = value
				}

			} else if enumVal, ok := obj.(*common.EnumValueInstance); ok {
				if enumVal.Fields[target.Name] == nil {
					enumVal.Fields[target.Name] = value
				} else {
					return nil, false, ThrowRuntimeError(env, fmt.Sprintf("cannot modify field '%s' of the enum '%s'", target.Name, enumVal.Definition.Name))
				}

			} else if recordVal, ok := obj.(*common.RecordInstance); ok {

				if recordVal.Values[target.Name] == nil {
					recordVal.Values[target.Name] = value
				} else {
					for k := range cur.Parent.Finals {
						c, cok := cur.Parent.Vars[k].(*common.RecordInstance)
						if cok && c == recordVal {
							if isFinal, exists := cur.Parent.Finals[k]; exists && isFinal {
								return nil, false, ThrowRuntimeError(cur.Parent, fmt.Sprintf("cannot modify to final record '%s'", k))
							}
						}
					}
					recordVal.Values[target.Name] = value
				}
			} else if m, ok := obj.(map[string]any); ok {
				if cur != nil {
					for k := range cur.Finals {
						c, cok := cur.Vars[k].(map[string]any)
						if cok && reflect.ValueOf(c).Pointer() == reflect.ValueOf(m).Pointer() {
							if isFinal, exists := cur.Finals[k]; exists && isFinal {
								return nil, false, ThrowRuntimeError(cur, fmt.Sprintf("cannot modify to final map '%s'", k))
							}
						}
					}
				}
				m[target.Name] = value
			} else {
				return nil, false, ThrowTypeError(env, "object with assignable fields", obj)
			}
		case *ast.IndexExpr:
			base, err := evalExpr(env, target.X)
			if err != nil {
				return nil, false, err
			}

			index, err := evalExpr(env, target.Index)
			if err != nil {
				return nil, false, err
			}

			instance, ok := base.(*ClassInstance)
			if !ok {
				return nil, false, ThrowTypeError(env, "indexable type", base)
			}
			indexableDef := common.BuiltinIndexableInterface.GetInterfaceDefinition(env)
			if indexableDef == nil {
				return nil, false, fmt.Errorf("Indexable interface not found")
			}
			if instance.ParentClass.ImplementsInterface(indexableDef) {
				setOverloads, valid := instance.ParentClass.Methods["__set"]
				if !valid {
					return nil, false, fmt.Errorf("Indexable object missing __set() method")
				}
				method := common.SelectMethodOverload(setOverloads, 2)
				if method == nil {
					return nil, false, ThrowRuntimeError(env, fmt.Sprintf("no overload found for %s.__set with %d arguments", instance.ClassName, 2))
				}
				_, err = CallInstanceMethod(instance, *method, env, []any{index, value})
				if err != nil {
					return nil, false, err
				}
				return value, false, nil
			}

		default:
			return nil, false, ThrowRuntimeError(env, fmt.Sprintf("invalid assignment target: %T", target))
		}

		return value, false, nil
	case *ast.ReturnStmt:
		v, err := evalExpr(env, s.Value)
		return v, true, err
	case *ast.ExprStmt:
		v, err := evalExpr(env, s.X)
		return v, false, err
	case *ast.DefStmt:
		// Check if this is a generic function
		isGeneric := len(s.TypeParams) > 0

		// Capture current env for closure
		fn := common.Func(func(callEnv *common.Env, args []any) (any, error) {
			local := &common.Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}, Defers: []func() error{}}

			// For generic functions, we need to handle type parameters
			// In a simple implementation, we just make them available as types in the local scope
			if isGeneric {
				// Store type parameter names in the local environment
				// This allows the function body to reference them
				for _, tp := range s.TypeParams {
					// Make type parameters available as type identifiers
					local.Set("__type_"+tp.Name, tp.Name)
				}
			}

			// Use the new parameter binding with variadic support
			err := bindParametersWithVariadic(local, s.Params, args)
			if err != nil {
				return nil, err
			}

			// Ensure defers are executed even if function returns early or errors
			defer func() {
				// Execute deferred calls in LIFO order
				for i := len(local.Defers) - 1; i >= 0; i-- {
					// Ignore errors in deferred calls (or could collect them)
					_ = local.Defers[i]()
				}
			}()

			for _, st := range s.Body {
				v, ret, err := evalStmt(local, st)
				if err != nil {
					return nil, err
				}
				if ret {
					return v, nil
				}
			}
			return nil, nil
		})

		// Infer return type if not explicitly specified
		returnType := s.ReturnType
		if returnType == nil {
			returnType = common.InferReturnType(s.Body, env)
		}

		// Wrap function in FunctionDefinition with metadata
		funcDef := &common.FunctionDefinition{
			Name:        s.Name,
			Func:        fn,
			Params:      s.Params,
			ReturnType:  returnType,
			AccessLevel: s.AccessLevel,
			Modifiers:   s.Modifiers,
			FileName:    env.GetFileName(),
			PackageName: env.GetPackageName(),
		}

		env.Set(s.Name, funcDef)
		return funcDef, false, nil
	case *ast.TypeAliasStmt:
		// Register type alias
		packageName := env.GetPackageName()
		if typeAliasRegistry[packageName] == nil {
			typeAliasRegistry[packageName] = make(map[string]*TypeAlias)
		}
		
		// Check if alias already exists
		if _, exists := typeAliasRegistry[packageName][s.Name]; exists {
			return nil, false, ThrowRuntimeError(env, fmt.Sprintf("type alias '%s' already defined in package '%s'", s.Name, packageName))
		}
		
		// Create and register the alias
		alias := &TypeAlias{
			Name:        s.Name,
			BaseType:    s.BaseType,
			IsFinal:     s.IsFinal,
			PackageName: packageName,
		}
		typeAliasRegistry[packageName][s.Name] = alias
		
		return nil, false, nil
	case *ast.InterfaceDecl:
		_, err := evalInterfaceDecl(env, s)
		return nil, false, err
	case *ast.ClassDecl:
		val, err := evalClassDecl(env, s)
		return val, false, err
	case *ast.EnumDecl:
		val, err := evalEnumDecl(env, s)
		return val, false, err
	case *ast.RecordDecl:
		val, err := evalRecordDecl(env, s)
		return val, false, err
	case *ast.TryStmt:
		return evalTryStmt(env, s)
	case *ast.ThrowStmt:
		return evalThrowStmt(env, s)
	case *ast.DeferStmt:
		return evalDeferStmt(env, s)
	case *ast.SelectStmt:
		return evalSelectStmt(env, s)
	case *ast.SwitchStmt:
		return evalSwitchStmt(env, s)
	default:
		return nil, false, ThrowNotImplementedError(env, fmt.Sprintf("statement type %T", st))
	}
}

// installBuiltins populates env with standard namespaces and functions.
func installBuiltins(env *common.Env, opts Options) {
	out := opts.Stdout
	if out == nil {
		out = io.Discard
	}
	env.Set("print", common.Func(func(_ *common.Env, args []any) (any, error) {
		for i, a := range args {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, utils.ToString(a))
		}
		return nil, nil
	}))
	env.Set("println", common.Func(func(_ *common.Env, args []any) (any, error) {
		for i, a := range args {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, utils.ToString(a))
		}
		fmt.Fprintln(out)
		return nil, nil
	}))

	env.Set("int", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		i, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "int-convertible value", args[0])
		}
		return CreateIntInstance(e, i)
	}))
	env.Set("float", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		f, ok := utils.AsFloat(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "float-convertible value", args[0])
		}
		return CreateFloatInstance(e, f)
	}))
	env.Set("str", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		return CreateStringInstance(e, utils.ToString(args[0]))
	}))
	env.Set("bool", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		b := utils.AsBool(args[0])
		return CreateBoolInstance(e, b)
	}))

	// len() - get length of string, array (ClassInstance), or map (ClassInstance)
	env.Set("len", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		switch v := args[0].(type) {
		case string:
			return float64(len(v)), nil
		case *ClassInstance:
			// For Array and Map ClassInstances, use their length() method
			if v.ClassName == "Array" {
				if items, ok := v.Fields["_items"].([]any); ok {
					return float64(len(items)), nil
				}
			} else if v.ClassName == "Map" {
				if data, ok := v.Fields["_data"].(map[uint64][]*mapEntry); ok {
					size := 0
					for _, entries := range data {
						size += len(entries)
					}
					return float64(size), nil
				}
			}
			return nil, ThrowTypeError(e, "string, array, or map", args[0])
		default:
			return nil, ThrowTypeError(e, "string, array, or map", args[0])
		}
	}))

	// range() - generate a range of numbers (global convenience function)
	env.Set("range", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) == 0 || len(args) > 3 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}

		var start, stop, step int
		if len(args) == 1 {
			// range(stop) -> 0..stop-1
			stop, _ = utils.AsInt(args[0])
			start = 0
			step = 1
		} else if len(args) == 2 {
			// range(start, stop) -> start..stop-1
			start, _ = utils.AsInt(args[0])
			stop, _ = utils.AsInt(args[1])
			step = 1
		} else {
			// range(start, stop, step)
			start, _ = utils.AsInt(args[0])
			stop, _ = utils.AsInt(args[1])
			step, _ = utils.AsInt(args[2])
			if step == 0 {
				return nil, ThrowValueError(e, "range() step argument must not be zero")
			}
		}

		result := []any{}
		if step > 0 {
			for i := start; i < stop; i += step {
				result = append(result, float64(i))
			}
		} else {
			for i := start; i > stop; i += step {
				result = append(result, float64(i))
			}
		}
		// Return an Array instance, not a native slice
		return CreateArrayInstance((*Env)(e), result)
	}))

	// str() - convert to string (global convenience function)
	env.Set("str", common.Func(func(e *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(e), 1, len(args))
		}
		return utils.ToString(args[0]), nil
	}))

	// Install Net module
	InstallNetModule(env, opts)
	InstallHttpModule(env, opts)

	// Install Iterable interface (base for all collections)
	if err := InstallIterableInterface((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Iterable interface: %v\n", err)
	}
	if err := InstallCollectionInterface((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Collection interface: %v\n", err)
	}
	if err := InstallSliceableInterface((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Sliceable interface: %v\n", err)
	}
	if err := InstallIndexableInterface((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Indexable interface: %v\n", err)
	}
	// Install Unstructured interface (for destructuring support)
	if err := InstallUnstructuredInterface((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Unstructured interface: %v\n", err)
	}

	// Install MapEntry builtin class
	if err := InstallPairBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Pair builtin: %v\n", err)
	}

	// Install Tuple builtin (immutable tuple with Unstructured support)
	if err := InstallTupleClass((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Tuple builtin: %v\n", err)
	}

	// Install Generic builtin FIRST (wraps native Go types)
	if err := InstallGenericBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Generic builtin: %v\n", err)
	}

	// Install Int and Float builtins as classes (so other types can reference them)
	if err := InstallNumberBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Number builtins: %v\n", err)
	}

	// Install Bool builtin as a class
	if err := InstallBoolBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Bool builtin: %v\n", err)
	}

	// Install String builtin as a class (can now reference Int for parameters)
	if err := InstallStringBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install String builtin: %v\n", err)
	}

	// Install unified Array builtin
	if err := InstallArrayBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Array builtin: %v\n", err)
	}

	// Install unified Map builtin (replaces Object and old Map)
	if err := InstallMapBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Map builtin: %v\n", err)
	}

	// Install Pair builtin (for key-value pairs)
	if err := InstallPairBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Pair builtin: %v\n", err)
	}

	// Install Range builtin as a class (iterable but not unstructured)
	if err := InstallRangeBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Range builtin: %v\n", err)
	}

	// Install Channel builtin as a class
	if err := InstallChannelBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Channel builtin: %v\n", err)
	}

	// Install List<T> builtin as a class
	if err := InstallListBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install List builtin: %v\n", err)
	}

	// Install Set<T> builtin as a class
	if err := InstallSetBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Set builtin: %v\n", err)
	}

	// Install Deque<T> builtin as a class
	if err := InstallDequeBuiltin((*Env)(env)); err != nil {
		fmt.Printf("Warning: Failed to install Deque builtin: %v\n", err)
	}

	// Install legacy generic types system (Lambda - legacy support)
	InstallGenerics(env)

	// Install Async/Await with Promise and CompletableFuture
	InstallAsyncAwait(env)

	// Initialize base Annotation class
	InitializeAnnotationBase(env)

	// Initialize annotation interfaces
	InitializeAnnotationInterfaces(env)
}

// Builtin modules registry
var builtinModules = map[string]func() map[string]any{}

// moduleCache avoids reloading same module
var moduleCache = map[string]map[string]any{}

// handleImport loads a module from libs/ or src/ and binds exported names.
func handleImport(env *common.Env, im *ast.ImportStmt) error {
	// resolve path to .pf file; try libs first then src
	rel := filepath.Join(im.Path...)
	// Builtin registry first
	if ctor, ok := builtinModules[strings.Join(im.Path, ".")]; ok {
		symbols := ctor()
		return bindImports(env, im, symbols)
	}

	homeDir, _ := os.UserHomeDir()
	candidates := []string{}

	// If we have a current file context, try relative imports from current directory first
	if env.FileName != "" {
		currentDir := filepath.Dir(env.FileName)
		// Relative import: import from same directory or subdirectory
		candidates = append(candidates,
			filepath.Join(currentDir, rel+".pf"),                     // same directory: helper.pf
			filepath.Join(currentDir, rel, "index.pf"),               // subdirectory with index
			filepath.Join(currentDir, rel, filepath.Base(rel)+".pf"), // subdirectory/subdirectory.pf
		)
	}

	// Standard library paths
	candidates = append(candidates,
		// Try user-specific library directory
		filepath.Join("libs", rel)+".pf",                                      // libs/math/vector.pf (single file)
		filepath.Join("libs", rel, "index.pf"),                                // libs/math/vector/index.pf (public API aggregator)
		filepath.Join("libs", rel, rel[strings.LastIndex(rel, "/")+1:]+".pf"), // libs/math/vector/vector.pf
		filepath.Join("src", rel+".pf"),
		filepath.Join("src", rel, "index.pf"),
	)
	if homeDir != "" {
		globalLib := filepath.Join(homeDir, ".polyloft", "libs")
		globalSrc := filepath.Join(homeDir, ".polyloft", "src")

		candidates = append(candidates,
			filepath.Join(globalLib, rel)+".pf",
			filepath.Join(globalLib, rel, "index.pf"),
			filepath.Join(globalLib, rel, filepath.Base(rel)+".pf"),
			filepath.Join(globalSrc, rel+".pf"),
			filepath.Join(globalSrc, rel, "index.pf"),
		)
	}
	var modKey string
	for _, cand := range candidates {
		if fi, err := os.Stat(cand); err == nil {
			if fi.IsDir() {
				continue
			}
			modKey = cand
			break
		}
	}
	if modKey == "" {
		// Try directory with multiple .pf files
		dir := filepath.Join("libs", rel)
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			modKey = dir + string(os.PathSeparator)
		} else {
			dir = filepath.Join("src", rel)
			if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
				modKey = dir + string(os.PathSeparator)
			}
		}
	}
	if modKey == "" {
		return ThrowRuntimeError(env, fmt.Sprintf("module not found: %s", rel))
	}

	if cached, ok := moduleCache[modKey]; ok {
		return bindImports(env, im, cached)
	}

	// Load and eval
	symbols := map[string]any{}
	if strings.HasSuffix(modKey, string(os.PathSeparator)) {
		// directory: load all .pf files and merge exports
		entries, _ := os.ReadDir(strings.TrimSuffix(modKey, string(os.PathSeparator)))
		for _, e := range entries {
			if e.IsDir() || filepath.Ext(e.Name()) != ".pf" {
				continue
			}
			fp := filepath.Join(strings.TrimSuffix(modKey, string(os.PathSeparator)), e.Name())
			m, err := loadModuleFile(fp, env)
			if err != nil {
				return err
			}
			for k, v := range m {
				symbols[k] = v
			}
		}
	} else {
		m, err := loadModuleFile(modKey, env)
		if err != nil {
			return err
		}
		for k, v := range m {
			symbols[k] = v
		}
	}
	moduleCache[modKey] = symbols
	return bindImports(env, im, symbols)
}

func bindImports(env *common.Env, im *ast.ImportStmt, symbols map[string]any) error {
	// Determine the source package for this import
	sourcePackage := strings.Join(im.Path, "/")

	// Mark package as imported
	if env.ImportedPackages == nil {
		env.ImportedPackages = make(map[string]struct{})
	}
	env.ImportedPackages[sourcePackage] = struct{}{}

	if len(im.Names) == 0 {
		// Namespace import: create nested structure for full path (e.g., test.math)
		ns := map[string]any{}
		for k, v := range symbols {
			// If symbol is a sealed enum, ensure importer is permitted
			if def, ok := enumRegistry[k]; ok && def.IsSealed {
				if !isEnumAccessPermitted(env, def) {
					return ThrowRuntimeError(env, fmt.Sprintf("cannot import sealed enum %s into %s", k, env.GetPackageName()))
				}
			}

			// Track imported classes
			if classCtor, ok := v.(*common.ClassConstructor); ok {
				if env.ImportedClasses == nil {
					env.ImportedClasses = make(map[string]string)
				}
				// Don't re-import if already imported
				if _, alreadyImported := env.ImportedClasses[k]; !alreadyImported {
					env.ImportedClasses[k] = classCtor.Definition.PackageName
				}
			}

			ns[k] = v
		}

		// Create nested namespace structure: test.math.Point becomes test -> math -> Point
		// This allows accessing as test.math.Point
		if len(im.Path) > 1 {
			// Build from innermost to outermost
			current := ns
			for i := len(im.Path) - 1; i > 0; i-- {
				wrapper := map[string]any{
					im.Path[i]: current,
				}
				current = wrapper
			}

			// Set the top-level namespace
			topLevel := im.Path[0]

			// Merge with existing namespace if it exists
			if existing, exists := env.Get(topLevel); exists {
				if existingMap, ok := existing.(map[string]any); ok {
					// Merge the new namespace into existing
					for k, v := range current {
						existingMap[k] = v
					}
					return nil
				}
			}

			env.Set(topLevel, current)
		} else {
			// Single-segment path, just set directly
			name := im.Path[0]
			env.Set(name, ns)
		}
		return nil
	}

	for _, n := range im.Names {
		v, ok := symbols[n]
		if !ok {
			return ThrowNameError(env, n)
		}

		// If symbol is a sealed enum, ensure importer is permitted
		if def, ok := enumRegistry[n]; ok && def.IsSealed {
			if !isEnumAccessPermitted(env, def) {
				return ThrowRuntimeError(env, fmt.Sprintf("cannot import sealed enum %s into %s", n, env.GetPackageName()))
			}
		}

		// Track imported classes
		if classCtor, ok := v.(*common.ClassConstructor); ok {
			if env.ImportedClasses == nil {
				env.ImportedClasses = make(map[string]string)
			}
			// Don't re-import if already imported
			if _, alreadyImported := env.ImportedClasses[n]; !alreadyImported {
				env.ImportedClasses[n] = classCtor.Definition.PackageName
			}
		}

		env.Set(n, v)
	}
	return nil
}

// isEnumAccessPermitted returns true if the current env/package is allowed to access a sealed enum
func isEnumAccessPermitted(env *common.Env, def *common.EnumDefinition) bool {
	// if not sealed it's always permitted
	if def == nil || !def.IsSealed {
		return true
	}
	// allow same package as definition
	if env.GetPackageName() != "" && env.GetPackageName() == def.PackageName {
		return true
	}
	pkg := env.GetPackageName()
	fileName := env.GetFileName()
	fileBase := filepath.Base(fileName)
	fileStem := strings.TrimSuffix(fileBase, filepath.Ext(fileBase))
	// allow if package name, file base, or file stem matches permitted entries
	for _, p := range def.Permits {
		if p.Name == pkg || p.Name == fileBase || p.Name == fileStem || p.Name == def.PackageName || p.Name == def.Name {
			return true
		}
	}
	return false
}

// loadModuleFile parses and evaluates a .pf file, returning its exported symbols.
// It inherits builtins from the parent environment to avoid re-creating them.
func loadModuleFile(path string, parentEnv *common.Env) (map[string]any, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lx := &lexer.Lexer{}
	items := lx.Scan(b)
	p := parser.NewWithFile(items, path)
	prog, err := p.Parse()
	if err != nil {
		return nil, err
	}
	// evaluate program in a child env that inherits builtins from the parent
	// This avoids re-creating builtin modules for each import
	packageName := filepath.Dir(path)

	// Create an env that inherits from parent to reuse builtins
	// but has its own file context
	var env *Env
	if parentEnv != nil {
		// Create a child that inherits builtins but with new file context
		env = &Env{
			Parent:      parentEnv,
			Vars:        map[string]any{},
			Consts:      map[string]bool{},
			Finals:      map[string]bool{},
			FileName:    path,
			PackageName: packageName,
		}
	} else {
		// Fallback: create fresh env with builtins if no parent
		env = NewEnvWithContext(path, packageName)
		opts := Options{}
		installBuiltins(env, opts)
		InstallSysModule(env, opts)
		InstallMathModule(env)
		InstallExceptionBuiltins(env)
	}

	// Set file environment variables (like Python's __name__ but with $ prefix)
	// These variables cannot be exported from the module
	env.Set("$name", filepath.Base(path))                                                        // e.g., "Vec2.pf"
	env.Set("$file", path)                                                                       // e.g., "libs/math/vector/Vec2.pf"
	env.Set("$package", packageName)                                                             // e.g., "libs/math/vector"
	env.Set("$stem", strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))) // e.g., "Vec2"

	_, err = evalProgramWithEnv(env, prog)
	if err != nil {
		return nil, err
	}
	// All vars in env are considered exported if declared with 'export' (not yet) or top-level names at end.
	// For now, expose all top-level function definitions and classes produced by names set in env.
	// Important: File environment variables (starting with $) are NOT exported
	// Respect access modifiers: private functions/classes are NOT exported
	out := map[string]any{}
	for k, v := range env.Vars {
		// Skip file environment variables (they start with $)
		if strings.HasPrefix(k, "$") {
			continue
		}

		// Check if this is a FunctionDefinition and respect access level
		if funcDef, ok := v.(*common.FunctionDefinition); ok {
			// Only export public functions (private/protected are not exported)
			if funcDef.AccessLevel == "private" {
				continue
			}
			// For protected, only export if it's the same package (not really applicable for exports)
			// In module exports, we only export public items
			if funcDef.AccessLevel == "protected" {
				continue
			}
			// Export the actual function, not the wrapper
			out[k] = funcDef.Func
			continue
		}

		// Check if this is a ClassConstructor and respect access level
		if classCtor, ok := v.(*common.ClassConstructor); ok {
			// Only export public classes
			if classCtor.Definition.AccessLevel == "private" {
				continue
			}
			if classCtor.Definition.AccessLevel == "protected" {
				continue
			}
			// Export the constructor as-is
			out[k] = v
			continue
		}

		// For other values (enums, records, variables), export them
		// TODO: Add access control for variables too
		out[k] = v
	}
	return out, nil
}

// evalProgramWithEnv runs statements into provided env using same evaluator.
func evalProgramWithEnv(env *common.Env, prog *ast.Program) (any, error) {
	var last any
	for _, st := range prog.Stmts {
		v, ret, err := evalStmt(env, st)
		if err != nil {
			return nil, err
		}
		if ret {
			return v, nil
		}
		last = v
	}
	return last, nil
}

// runBlock executes a list of statements with handling for break/continue.
func runBlock(env *common.Env, body []ast.Stmt) (brk, cont, ret bool, val any, err error) {
	for _, st := range body {
		v, r, err := evalStmt(env, st)
		if err != nil {
			return false, false, false, nil, err
		}
		switch v.(type) {
		case breakSentinel:
			return true, false, false, nil, nil
		case continueSentinel:
			return false, true, false, nil, nil
		}
		if r {
			return false, false, true, v, nil
		}
	}
	return false, false, false, nil, nil
}

func evalExpr(env *common.Env, e ast.Expr) (any, error) {
	switch x := e.(type) {
	case *ast.Ident:
		v, ok := env.Get(x.Name)
		if !ok {
			return nil, ThrowNameError(env, x.Name)
		}
		return v, nil
	case *ast.NumberLit:
		switch v := x.Value.(type) {
		case int:
			return CreateIntInstance(env, v)
		case float64:
			return CreateFloatInstance(env, v)
		case float32:
			return CreateFloatInstance(env, float64(v))
		}
		return x.Value, nil
	case *ast.StringLit:
		// Check if string contains interpolation
		if strings.Contains(x.Value, "#{") {
			return processStringInterpolation(env, x.Value)
		}
		return CreateStringInstance(env, x.Value)
	case *ast.BoolLit:
		return CreateBoolInstance(env, x.Value)
	case *ast.NilLit:
		return nil, nil
	case *ast.ArrayLit:
		arr := make([]any, 0, len(x.Elems))
		for _, e := range x.Elems {
			v, err := evalExpr(env, e)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		// Create an Array instance - builtin class must be available
		arrayInstance, err := CreateArrayInstance(env, arr)
		if err != nil {
			return nil, err
		}
		return arrayInstance, nil
	case *ast.MapLit:
		m := map[string]any{}
		for _, p := range x.Pairs {
			v, err := evalExpr(env, p.Value)
			if err != nil {
				return nil, err
			}
			m[p.Key] = v
		}
		// Create a Map instance - builtin class must be available
		mapInstance, err := CreateMapInstance(env, m)
		if err != nil {
			return nil, err
		}
		return mapInstance, nil
	case *ast.IndexExpr:
		base, err := evalExpr(env, x.X)
		if err != nil {
			return nil, err
		}

		// Check if this is a range index expression
		if rangeIdx, ok := x.Index.(*ast.RangeExpr); ok {
			// Handle range slicing: arr[start...end]
			startVal, err := evalExpr(env, rangeIdx.Start)
			if err != nil {
				return nil, err
			}
			endVal, err := evalExpr(env, rangeIdx.End)
			if err != nil {
				return nil, err
			}

			start, ok1 := utils.AsInt(startVal)
			end, ok2 := utils.AsInt(endVal)
			if !ok1 || !ok2 {
				return nil, ThrowTypeError(env, "integer", "range indices")
			}

			instance, ok := base.(*ClassInstance)
			if !ok {
				return nil, ThrowTypeError(env, "sliceable type", base)
			}

			sliceableInterface := common.BuiltinSliceableInterface.GetInterfaceDefinition(env)
			if instance.ParentClass.ImplementsInterface(sliceableInterface) {
				methodOverloads, exists := instance.ParentClass.Methods["__slice"]
				if !exists {
					return nil, ThrowAttributeError(env, "__slice", fmt.Sprintf("class '%s'", instance.ClassName))
				}
				// Select the correct overload based on argument count (start, end)
				method := common.SelectMethodOverload(methodOverloads, 2)
				if method == nil {
					return nil, ThrowRuntimeError(env, fmt.Sprintf("no overload found for %s.__slice with %d arguments", instance.ClassName, 2))
				}
				result, err := CallInstanceMethod(instance, *method, env, []any{start, end})
				if err != nil {
					return nil, err
				}
				return result, nil
			}
		}

		idx, err := evalExpr(env, x.Index)
		if err != nil {
			return nil, err
		}

		instance, ok := base.(*ClassInstance)
		if !ok {
			return nil, ThrowTypeError(env, "indexable type", base)
		}
		// Check if instance implements Indexable interface
		indexableInterface := common.BuiltinIndexableInterface.GetInterfaceDefinition(env)
		if instance.ParentClass.ImplementsInterface(indexableInterface) {
			containsOverloads, exists := instance.ParentClass.Methods["__contains"]
			if !exists {
				return nil, ThrowAttributeError(env, "__contains", fmt.Sprintf("class '%s'", instance.ClassName))
			}
			// Select the correct overload based on argument count (index)
			method := common.SelectMethodOverload(containsOverloads, 1)
			if method == nil {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("no overload found for %s.__contains with %d arguments", instance.ClassName, 1))
			}
			result, err := CallInstanceMethod(instance, *method, env, []any{idx})
			if err != nil {
				return nil, err
			}
			if !utils.AsBool(result) {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("Key not found: %v", idx))
			}
			// now call __get to retrieve the value
			getOverloads, exists := instance.ParentClass.Methods["__get"]
			if !exists {
				return nil, ThrowAttributeError(env, "__get", fmt.Sprintf("class '%s'", instance.ClassName))
			}
			// Select the correct overload based on argument count (index)
			method = common.SelectMethodOverload(getOverloads, 1)
			if method == nil {
				return nil, ThrowRuntimeError(env, fmt.Sprintf("no overload found for %s.__get with %d arguments", instance.ClassName, 1))
			}
			result, err = CallInstanceMethod(instance, *method, env, []any{idx})
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, ThrowTypeError(env, "indexable type", base)
	case *ast.FieldExpr:
		// Check if this is a static method call on a built-in type
		if ident, ok := x.X.(*ast.Ident); ok {
			// Check if it's a class static method or field access
			if classDef, exists := lookupClass(ident.Name, env.GetPackageName()); exists {
				// Check for static fields first
				if value, fieldExists := classDef.StaticFields[x.Name]; fieldExists {
					return value, nil
				}
				// Check for static methods (with overload support)
				if methodOverloads, methodExists := classDef.Methods[x.Name]; methodExists {
					// Return a function wrapper that selects the right overload
					return common.Func(func(callEnv *common.Env, args []any) (any, error) {
						// Select appropriate method based on argument count
						method := common.SelectMethodOverload(methodOverloads, len(args))
						if method == nil {
							return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no overload found for %s.%s with %d arguments", classDef.Name, x.Name, len(args)))
						}

						// Check if the method is static
						if !method.IsStatic {
							return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("method %s.%s is not static", classDef.Name, x.Name))
						}

						// Create a new environment for the static method
						methodEnv := callEnv.Child()

						// Bind parameters (including variadic) - validates and binds args
						if method.Params != nil {
							err := bindParametersWithVariadic(methodEnv, method.Params, args)
							if err != nil {
								return nil, err
							}
						}

						// Execute builtin implementation if available
						if method.BuiltinImpl != nil {
							// For builtin methods, parameters are already bound in methodEnv
							// The builtin can access them by name
							return method.BuiltinImpl(methodEnv, args)
						}

						// Execute method body for non-builtin methods
						var result any
						for _, stmt := range method.Body {
							var err error
							val, returned, err := evalStmt(methodEnv, stmt)
							if err != nil {
								return nil, err
							}
							if returned {
								result = val
								break
							}
						}
						return result, nil
					}), nil
				}
				return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("class '%s' (static access required)", ident.Name))
			}
			// Check if it's an interface static field access
			if interfaceDef, exists := interfaceRegistry[ident.Name]; exists {
				if value, fieldExists := interfaceDef.StaticFields[x.Name]; fieldExists {
					return value, nil
				}
				return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("interface '%s'", ident.Name))
			}
		}

		base, err := evalExpr(env, x.X)
		if err != nil {
			return nil, err
		}
		switch b := base.(type) {
		case *common.EnumConstructor:
			// Access fields from the wrapped enum object
			return b.EnumObject[x.Name], nil
		case *ClassDefinition:
			// Access static fields and methods on ClassDefinition
			// Check for static fields first
			if value, fieldExists := b.StaticFields[x.Name]; fieldExists {
				return value, nil
			}
			// Check for static methods (with overload support)
			if methodOverloads, methodExists := b.Methods[x.Name]; methodExists {
				// Return a function wrapper that selects the right overload
				return common.Func(func(callEnv *common.Env, args []any) (any, error) {
					// Select appropriate method based on argument count
					method := common.SelectMethodOverload(methodOverloads, len(args))
					if method == nil {
						return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("no overload found for %s.%s with %d arguments", b.Name, x.Name, len(args)))
					}

					if !method.IsStatic {
						return nil, ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("method %s.%s is not static", b.Name, x.Name))
					}

					// Create a new environment for the static method
					methodEnv := callEnv.Child()

					// Bind parameters (including variadic) - validates and binds args
					if method.Params != nil {
						err := bindParametersWithVariadic(methodEnv, method.Params, args)
						if err != nil {
							return nil, err
						}
					}

					// Execute builtin implementation if available
					if method.BuiltinImpl != nil {
						// For builtin methods, parameters are already bound in methodEnv
						// The builtin can access them by name
						return method.BuiltinImpl(methodEnv, args)
					}

					// Execute method body for non-builtin methods
					var result any
					for _, stmt := range method.Body {
						var err error
						val, returned, err := evalStmt(methodEnv, stmt)
						if err != nil {
							return nil, err
						}
						if returned {
							result = val
							break
						}
					}

					return result, nil
				}), nil
			}
			return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("class '%s'", b.Name))
		case *ClassInstance:
			// Special handling for Map instances to support field access syntax
			if b.ClassName == "Map" {
				if hashData, ok := b.Fields["_data"].(map[uint64][]*mapEntry); ok {
					// Look for the key by hashing the field name and checking entries
					hash := hashValue(env, x.Name)
					if entries, exists := hashData[hash]; exists {
						for _, entry := range entries {
							if equals(entry.Key, x.Name) {
								return entry.Value, nil
							}
						}
					}
				}
			}

			// Check fields first
			if value, exists := b.Fields[x.Name]; exists {
				return value, nil
			}
			// Check methods
			if method, exists := b.Methods[x.Name]; exists {
				return method, nil
			}
			return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("'%s' instance", b.ClassName))
		case *common.EnumValueInstance:
			if value, exists := b.Fields[x.Name]; exists {
				return value, nil
			}
			if method, exists := b.Methods[x.Name]; exists {
				return method, nil
			}
			if b.Definition != nil {
				return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("enum value '%s.%s'", b.Definition.Name, b.Name))
			}
			return nil, ThrowAttributeError(env, x.Name, "enum value")
		case *common.RecordInstance:
			if value, exists := b.Values[x.Name]; exists {
				return value, nil
			}
			if method, exists := b.Methods[x.Name]; exists {
				return method, nil
			}
			if b.Definition != nil {
				return nil, ThrowAttributeError(env, x.Name, fmt.Sprintf("record '%s'", b.Definition.Name))
			}
			return nil, ThrowAttributeError(env, x.Name, "record")
		case float64:
			// Wrap primitive float in Float class instance
			floatInstance, err := CreateFloatInstance(env, b)
			if err != nil {
				return nil, ThrowAttributeError(env, x.Name, "float")
			}
			if method, exists := floatInstance.Methods[x.Name]; exists {
				return method, nil
			}
			return nil, ThrowAttributeError(env, x.Name, "Float")
		case int:
			// Wrap primitive int in Int class instance
			intInstance, err := CreateIntInstance(env, b)
			if err != nil {
				return nil, ThrowAttributeError(env, x.Name, "int")
			}
			if method, exists := intInstance.Methods[x.Name]; exists {
				return method, nil
			}
			return nil, ThrowAttributeError(env, x.Name, "Int")
		case string:
			// Wrap primitive string in String class instance
			stringInstance, err := CreateStringInstance(env, b)
			if err != nil {
				return nil, ThrowAttributeError(env, x.Name, "string")
			}
			if method, exists := stringInstance.Methods[x.Name]; exists {
				return method, nil
			}
			return nil, ThrowAttributeError(env, x.Name, "String")
		case bool:
			// Wrap primitive bool in Bool class instance
			boolInstance, err := CreateBoolInstance(env, b)
			if err != nil {
				return nil, ThrowAttributeError(env, x.Name, "bool")
			}
			if method, exists := boolInstance.Methods[x.Name]; exists {
				return method, nil
			}
			return nil, ThrowAttributeError(env, x.Name, "Bool")
		case map[string]any:
			// Support namespace imports: allow accessing map fields with dot notation
			if value, exists := b[x.Name]; exists {
				return value, nil
			}
			return nil, ThrowAttributeError(env, x.Name, "namespace")
		default:
			return nil, ThrowTypeError(env, "object with field access", base)
		}
	case *ast.UnaryExpr:
		v, err := evalExpr(env, x.X)
		if err != nil {
			return nil, err
		}
		switch x.Op {
		case ast.OpNot:
			return !utils.AsBool(v), nil
		case ast.OpNeg:
			f, ok := utils.AsFloat(v)
			if !ok {
				return nil, typeError("number", v)
			}
			return -f, nil
		default:
			return nil, ThrowNotImplementedError(env, fmt.Sprintf("unary operator %d", x.Op))
		}
	case *ast.BinaryExpr:
		a, err := evalExpr(env, x.Lhs)
		if err != nil {
			return nil, err
		}
		b, err := evalExpr(env, x.Rhs)
		if err != nil {
			return nil, err
		}
		switch x.Op {
		case ast.OpPlus:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "+", "add", a, b); handled {
				return result, err
			}
			// strings (check both primitive and String class instances)
			aStr := extractPrimitiveValue(a)
			if sa, ok := aStr.(string); ok {
				return sa + utils.ToString(b), nil
			}
			// numbers
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateFloatInstance(env, fa+fb)
		case ast.OpMinus:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "-", "subtract", a, b); handled {
				return result, err
			}
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateFloatInstance(env, fa-fb)
		case ast.OpMul:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "*", "multiply", a, b); handled {
				return result, err
			}

			//a and b must be *ClassInstance
			aClass, ok := a.(*ClassInstance)
			bClass, ok2 := b.(*ClassInstance)
			if !ok || !ok2 {
				return nil, typeError("ClassInstance", a, b)
			}
			if aClass.ParentClass == common.BuiltinTypeString.GetClassDefinition(env) && bClass.ParentClass == common.BuiltinTypeInt.GetClassDefinition(env) {
				count, ok := utils.AsInt(b)
				if !ok {
					return nil, typeError("int", b)
				}
				originalStr := utils.ToString(a)
				return CreateStringInstance(env, strings.Repeat(originalStr, count))
			}
			if bClass.ParentClass == common.BuiltinTypeString.GetClassDefinition(env) && aClass.ParentClass == common.BuiltinTypeInt.GetClassDefinition(env) {
				count, ok := utils.AsInt(a)
				if !ok {
					return nil, typeError("int", a)
				}
				originalStr := utils.ToString(b)
				return CreateStringInstance(env, strings.Repeat(originalStr, count))
			}
			//check if any of the 2 are float
			if bClass.ParentClass == common.BuiltinTypeFloat.GetClassDefinition(env) || aClass.ParentClass == common.BuiltinTypeFloat.GetClassDefinition(env) {
				fa, oka := utils.AsFloat(a)
				fb, okb := utils.AsFloat(b)
				if !oka || !okb {
					return nil, typeError("number", a, b)
				}
				return CreateFloatInstance(env, fa*fb)
			}
			ia, oka := utils.AsInt(a)
			ib, okb := utils.AsInt(b)
			if !oka || !okb {
				return nil, typeError("int", a, b)
			}
			return CreateIntInstance(env, ia*ib)

		case ast.OpDiv:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "/", "divide", a, b); handled {
				return result, err
			}
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateFloatInstance(env, fa/fb)
		case ast.OpMod:
			ia, oka := utils.AsInt(a)
			ib, okb := utils.AsInt(b)
			if !oka || !okb {
				return nil, typeError("int", a, b)
			}
			return CreateIntInstance(env, ia%ib)
		case ast.OpEq:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "==", "equals", a, b); handled {
				return result, err
			}

			return CreateBoolInstance(env, equal(a, b))
		case ast.OpNeq:
			// Check for operator overloading first
			if result, handled, err := tryOperatorOverload(env, "!=", "notequals", a, b); handled {
				return result, err
			}
			return CreateBoolInstance(env, !equal(a, b))
		case ast.OpLt:
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateBoolInstance(env, fa < fb)
		case ast.OpLte:
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateBoolInstance(env, fa <= fb)
		case ast.OpGt:
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateBoolInstance(env, fa > fb)
		case ast.OpGte:
			fa, oka := utils.AsFloat(a)
			fb, okb := utils.AsFloat(b)
			if !oka || !okb {
				return nil, typeError("number", a, b)
			}
			return CreateBoolInstance(env, fa >= fb)
		case ast.OpAnd:
			if !utils.AsBool(a) {
				return CreateBoolInstance(env, false)
			}
			return CreateBoolInstance(env, utils.AsBool(b))
		case ast.OpOr:
			if utils.AsBool(a) {
				return CreateBoolInstance(env, true)
			}
			return CreateBoolInstance(env, utils.AsBool(b))
		default:
			return nil, ThrowNotImplementedError(env, fmt.Sprintf("binary operator %d", x.Op))
		}
	case *ast.CallExpr:
		cal, err := evalExpr(env, x.Callee)
		if err != nil {
			return nil, err
		}

		// Handle ClassConstructor wrapper
		var fn Func
		if classConstructor, ok := cal.(*common.ClassConstructor); ok {
			fn = classConstructor.Func
		} else if funcDef, ok := cal.(*common.FunctionDefinition); ok {
			// Unwrap FunctionDefinition to get the actual function
			fn = funcDef.Func
		} else if lambdaDef, ok := cal.(*common.LambdaDefinition); ok {
			// Unwrap LambdaDefinition to get the actual function
			fn = lambdaDef.Func
		} else if funcVal, ok := cal.(Func); ok {
			fn = funcVal
		} else {
			// Provide more detailed error information
			valueInfo := "nil"
			if cal != nil {
				valueInfo = utils.ToString(cal)
				// Truncate long values
				if len(valueInfo) > 50 {
					valueInfo = valueInfo[:47] + "..."
				}
			}

			return nil, ThrowNotCallableError(env, fmt.Sprintf("%T", cal), valueInfo)
		}

		args := make([]any, 0, len(x.Args))
		for _, a := range x.Args {
			v, err := evalExpr(env, a)
			if err != nil {
				return nil, err
			}
			args = append(args, v)
		}
		return fn(env, args)
	case *ast.GenericCallExpr:
		// Handle generic type constructor: List<Int>(), Map<String, Int>(), etc.
		return evalGenericCallExpr(env, x)
	case *ast.InstanceOfExpr:
		return evalInstanceOfExpr(env, x)
	case *ast.TypeExpr:
		return evalTypeExpr(env, x)
	case *ast.ThreadSpawnExpr:
		return evalThreadSpawnExpr(env, x)
	case *ast.ThreadJoinExpr:
		return evalThreadJoinExpr(env, x)
	case *ast.ChannelExpr:
		return evalChannelExpr(env, x)
	case *ast.RangeExpr:
		// Evaluate range expression to create a Range instance
		startVal, err := evalExpr(env, x.Start)
		if err != nil {
			return nil, err
		}
		endVal, err := evalExpr(env, x.End)
		if err != nil {
			return nil, err
		}

		start, ok1 := utils.AsInt(startVal)
		end, ok2 := utils.AsInt(endVal)
		if !ok1 || !ok2 {
			return nil, ThrowTypeError(env, "integer", "range bounds")
		}

		if start > end {
			return nil, ThrowValueError(env, "range start must be <= end")
		}

		// Create a Range instance (memory-efficient iterable)
		return CreateRangeInstance(env, start, end, 1)
	case *ast.TernaryExpr:
		condition, err := evalExpr(env, x.Condition)
		if err != nil {
			return nil, err
		}
		if utils.AsBool(condition) {
			return evalExpr(env, x.TrueBranch)
		} else {
			return evalExpr(env, x.FalseBranch)
		}
	case *ast.LambdaExpr:
		// Create a closure that captures the current environment
		fn := common.Func(func(callEnv *common.Env, args []any) (any, error) {
			// Create new environment for lambda execution
			lambdaEnv := &common.Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}, Defers: []func() error{}}

			// Bind parameters with type validation and variadic support
			err := bindParametersWithVariadic(lambdaEnv, x.Params, args)
			if err != nil {
				return nil, err
			}

			// Ensure defers are executed after lambda completes
			defer func() {
				// Execute deferred calls in LIFO order
				for i := len(lambdaEnv.Defers) - 1; i >= 0; i-- {
					_ = lambdaEnv.Defers[i]()
				}
			}()

			// Execute lambda body
			if x.IsBlock {
				// Multi-line lambda with statement block
				for _, stmt := range x.BlockBody {
					v, ret, err := evalStmt(lambdaEnv, stmt)
					if err != nil {
						return nil, err
					}
					if ret {
						return v, nil
					}
				}
				return nil, nil
			} else {
				// Single expression lambda
				return evalExpr(lambdaEnv, x.Body)
			}
		})

		// Infer return type if not explicitly specified
		returnType := x.ReturnType
		if returnType == nil {
			if x.IsBlock {
				returnType = common.InferReturnType(x.BlockBody, env)
			} else {
				// For single expression lambdas, infer from the expression
				returnType = common.InferExprType(x.Body, env)
			}
		}

		// Wrap lambda in LambdaDefinition with type information
		lambdaDef := &common.LambdaDefinition{
			Func:       fn,
			Params:     x.Params,
			ReturnType: returnType,
		}

		return lambdaDef, nil
	default:
		return nil, ThrowNotImplementedError(env, fmt.Sprintf("expression type %T", e))
	}
}

// tryOperatorOverload checks if an object has an operator overload method
func tryOperatorOverload(env *Env, op string, methodName string, left, right any) (any, bool, error) {
	// Check if left operand is a class instance with the operator method
	if instance, ok := left.(*ClassInstance); ok {
		// First try the operator symbol (e.g., "+", "-", "==")
		if method, exists := instance.Methods[op]; exists {
			result, err := method(env, []any{right})
			return result, true, err
		}
		// Fallback to method name (e.g., "add", "subtract")
		if method, exists := instance.Methods[methodName]; exists {
			result, err := method(env, []any{right})
			return result, true, err
		}
	}
	return nil, false, nil
}

func ThrowAttributeErrorWithHint(env *Env, attrName string, typeName string, availableMethods []string) error {
	message := fmt.Sprintf("'%s' object has no attribute '%s'", typeName, attrName)

	// Create hint with available methods
	hint := &ExceptionHint{
		Message:     fmt.Sprintf("Available methods: %v", availableMethods),
		Suggestions: availableMethods,
		HintType:    "method",
	}

	// Create exception with hint
	exc := &HyException{
		Message: message,
		Type:    "AttributeError",
		Hint:    hint,
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}
func equal(a, b any) bool {
	// Extract primitive values from class instances
	aVal := extractPrimitiveValue(a)
	bVal := extractPrimitiveValue(b)

	switch aa := aVal.(type) {
	case nil:
		return bVal == nil
	case string:
		bb, ok := bVal.(string)
		return ok && aa == bb
	case bool:
		bb, ok := bVal.(bool)
		return ok && aa == bb
	case float32:
		bb, ok := utils.AsFloat(bVal)
		return ok && float64(aa) == bb
	case float64:
		bb, ok := utils.AsFloat(bVal)
		return ok && aa == bb
	case int64:
		bb, ok := utils.AsFloat(bVal)
		return ok && float64(aa) == bb
	case int:
		bb, ok := utils.AsFloat(bVal)
		return ok && float64(aa) == bb
	default:
		return aVal == bVal
	}
}

// extractPrimitiveValue extracts the underlying primitive value from a class instance
func extractPrimitiveValue(v any) any {
	if instance, ok := v.(*ClassInstance); ok {
		switch instance.ClassName {
		case "String":
			if val, ok := instance.Fields["_value"].(string); ok {
				return val
			}
		case "Int":
			if val, ok := instance.Fields["_value"].(int); ok {
				return val
			}
		case "Float":
			if val, ok := instance.Fields["_value"].(float64); ok {
				return val
			}
		case "Bool":
			if val, ok := instance.Fields["_value"].(bool); ok {
				return val
			}
		}
	}
	return v
}

func typeError(exp string, got ...any) error {
	return ThrowTypeError(nil, exp, got...)
}

// Use common sentinels
type breakSentinel = common.BreakSentinel
type continueSentinel = common.ContinueSentinel

// Options control execution behavior (flags, limits, debug hooks, etc.).
type Options struct {
	Stdout io.Writer // where println/print write to
}

// Use common definitions for Env and Func
var NewEnv = common.NewEnv
var NewEnvWithContext = common.NewEnvWithContext

// Func represents a callable value, using the common definition
type Func = common.Func

// Env represents an environment, using the common definition
type Env = common.Env

// Thread represents a running thread

// Thread represents a running thread
type Thread struct {
	result chan any
	err    chan error
	done   bool
}

// evalThreadSpawnExpr evaluates thread spawn expressions
func evalThreadSpawnExpr(env *Env, expr *ast.ThreadSpawnExpr) (any, error) {
	thread := &Thread{
		result: make(chan any, 1),
		err:    make(chan error, 1),
		done:   false,
	}

	// Start goroutine to execute thread body
	go func() {
		defer func() {
			if r := recover(); r != nil {
				thread.err <- ThrowRuntimeError(env, fmt.Sprintf("thread panic: %v", r))
			}
		}()

		// Create a new environment for the thread
		threadEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}

		var lastResult any
		for _, stmt := range expr.Body {
			result, returned, err := evalStmt(threadEnv, stmt)
			if err != nil {
				thread.err <- err
				return
			}
			if returned {
				thread.result <- result
				thread.done = true
				return
			}
			lastResult = result
		}

		// If no explicit return, return last result
		thread.result <- lastResult
		thread.done = true
	}()

	return thread, nil
}

// evalThreadJoinExpr evaluates thread join expressions
func evalThreadJoinExpr(env *Env, expr *ast.ThreadJoinExpr) (any, error) {
	threadVal, err := evalExpr(env, expr.Thread)
	if err != nil {
		return nil, err
	}

	thread, ok := threadVal.(*Thread)
	if !ok {
		return nil, ThrowTypeError(env, "thread", threadVal)
	}

	// Wait for thread to complete
	select {
	case result := <-thread.result:
		return result, nil
	case err := <-thread.err:
		return nil, err
	}
}

// processStringInterpolation processes string interpolation in the format "text #{expr} more text"
func processStringInterpolation(env *Env, str string) (string, error) {
	result := ""
	i := 0

	for i < len(str) {
		// Find the next interpolation
		start := findNext(str, i, "#{")
		if start == -1 {
			// No more interpolations, append rest of string
			result += str[i:]
			break
		}

		// Append the text before interpolation
		result += str[i:start]

		// Find the end of the interpolation
		end := findMatchingBrace(str, start+2)
		if end == -1 {
			return "", ThrowRuntimeError(env, "unclosed interpolation expression in string")
		}

		// Extract and evaluate the expression
		exprStr := str[start+2 : end]
		value, err := evaluateInterpolationExpr(env, exprStr)
		if err != nil {
			return "", err
		}

		// Convert value to string and append
		result += utils.ToString(value)

		// Move past the closing brace
		i = end + 1
	}

	return result, nil
}

// findNext finds the next occurrence of substr starting from index start
func findNext(str string, start int, substr string) int {
	for i := start; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// findMatchingBrace finds the matching closing brace for an interpolation expression
func findMatchingBrace(str string, start int) int {
	depth := 1
	for i := start; i < len(str); i++ {
		if str[i] == '{' {
			depth++
		} else if str[i] == '}' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// evalGenericCallExpr evaluates generic type constructor calls like List<Int>(), Map<String, Int>(), List<?>()
func evalGenericCallExpr(env *common.Env, expr *ast.GenericCallExpr) (any, error) {
	// Get the base constructor function
	constructor, ok := env.Get(expr.Name)
	if !ok {
		return nil, ThrowNameError(env, expr.Name)
	}

	// Prepare arguments: type parameters come first, then constructor args
	var allArgs []any
	var gtypes []GenericType
	for _, tp := range expr.TypeParams {
		if tp.IsWildcard {
			// Create a GenericBound for the wildcard type parameter
			// For wildcards, the Bounds array contains the type names (e.g., ["Number"])
			// and WildcardKind is "extends", "super", or "unbounded"
			var boundTypeName string
			if len(tp.Bounds) > 0 {
				boundTypeName = tp.Bounds[0]
			}
			
			bound := common.GenericBound{
				Name:       ast.Type{Name: boundTypeName},
				Variance:   tp.WildcardKind,  // "extends", "super", or "unbounded"
				IsVariadic: tp.IsVariadic,
			}
			
			gtypes = append(gtypes, GenericType{
				Bounds: []common.GenericBound{bound},
			})
		} else {
			// Regular type parameter without variance
			allArgs = append(allArgs, tp.Name)
			
			// Also create a GenericType for regular type parameters
			bound := common.GenericBound{
				Name:       ast.Type{Name: tp.Name},
				Variance:   tp.Variance,  // This will be "" for non-wildcard types
				IsVariadic: tp.IsVariadic,
			}
			gtypes = append(gtypes, GenericType{
				Bounds: []common.GenericBound{bound},
			})
		}
	}

	// Evaluate and add constructor arguments
	for _, arg := range expr.Args {
		val, err := evalExpr(env, arg)
		if err != nil {
			return nil, err
		}
		allArgs = append(allArgs, val)
	}

	// Call the constructor function
	if fn, ok := common.ExtractFunc(constructor); ok {

		instance, err := fn(env, allArgs)
		switch instance.(type) {
		case *ClassInstance:
			val := instance.(*ClassInstance)
			if len(gtypes) > 0 {
				val.GenericTypes = gtypes
			}
			return instance, err
		}
		return instance, err
	}

	return nil, ThrowNotCallableError(env, fmt.Sprintf("%T", constructor), utils.ToString(constructor))
}

// evaluateInterpolationExpr evaluates a simple interpolation expression
// For now, we'll support basic field access like "this.x"
func evaluateInterpolationExpr(env *Env, exprStr string) (any, error) {
	exprStr = strings.TrimSpace(exprStr)

	// Handle "this.field" pattern
	if strings.HasPrefix(exprStr, "this.") {
		fieldName := exprStr[5:] // Remove "this."

		// Get "this" from environment
		thisVal, ok := env.Get("this")
		if !ok {
			return nil, ThrowNameError(env, "this")
		}

		// Get the field value
		switch obj := thisVal.(type) {
		case *ClassInstance:
			if value, exists := obj.Fields[fieldName]; exists {
				return value, nil
			}
			return nil, ThrowAttributeError(env, fieldName, fmt.Sprintf("class '%s'", obj.ClassName))
		default:
			return nil, ThrowTypeError(env, "object with fields", thisVal)
		}
	}

	// Handle simple variable access
	if val, ok := env.Get(exprStr); ok {
		return val, nil
	}

	return nil, ThrowRuntimeError(env, fmt.Sprintf("unsupported interpolation expression: %s", exprStr))
}

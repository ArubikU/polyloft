package engine

import (
	"fmt"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// HyException represents a Polyloft exception
type HyException struct {
	Message    string
	Type       string
	StackTrace []string
	Cause      *HyException
	Instance   any // The actual exception instance from Polyloft
	File       string
	Line       int
	Column     int
	Hint       *ExceptionHint
}

func (e *HyException) Error() string {
	return fmt.Sprintf(e.Message)
}

// WithLocation adds location information to the exception
func (e *HyException) WithLocation(file string, line, column int) *HyException {
	e.File = file
	e.Line = line
	e.Column = column
	return e
}

// WithHint adds a hint to the exception
func (e *HyException) WithHint(hint *ExceptionHint) *HyException {
	e.Hint = hint
	return e
}

func NewHyException(exceptionType, message string) *HyException {
	// Si existe un constructor registrado en exceptionClasses, usarlo
	if constructor, ok := exceptionClasses[exceptionType]; ok {
		instance, err := constructor(nil, []any{message})
		if err == nil {
			return &HyException{
				Message:    message,
				Type:       exceptionType,
				StackTrace: []string{},
				Instance:   instance,
			}
		}
	}

	// fallback si no se encuentra el constructor
	return &HyException{
		Message:    message,
		Type:       exceptionType,
		StackTrace: []string{},
	}
}

// Exception class registry
var exceptionClasses = map[string]common.Func{}

// createExceptionClassesProgrammatically creates exception classes using the ClassBuilder
func createExceptionClassesProgrammatically(env *Env) error {
	// Create Throwable base class
	throwableClass, throwableConstructor, err := NewClassBuilder("Throwable").
		AddField("message", &ast.Type{Name: "string", IsBuiltin: true}, []string{}).
		AddField("type", &ast.Type{Name: "string", IsBuiltin: true}, []string{}).
		AddField("stackTrace", &ast.Type{Name: "array", IsBuiltin: true}, []string{}).
		SetBuiltinConstructor(
			[]ast.Parameter{{Name: "message", Type: ast.TypeFromString("string")}},
			func(callEnv *common.Env, args []any) (any, error) {
				// Get the instance from "this"
				thisVal, _ := callEnv.Get("this")
				if instance, ok := thisVal.(*common.ClassInstance); ok {
					message := "Error"
					if len(args) > 0 {
						if msg, ok := args[0].(string); ok {
							message = msg
						}
					}
					instance.Fields["message"] = message
					instance.Fields["type"] = "Throwable"
					instance.Fields["stackTrace"] = []any{}
				}
				return nil, nil
			},
		).
		AddBuiltinMethod("toString", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{},
			func(callEnv *common.Env, args []any) (any, error) {
				thisVal, _ := callEnv.Get("this")
				if instance, ok := thisVal.(*common.ClassInstance); ok {
					msgType, _ := instance.Fields["type"].(string)
					message, _ := instance.Fields["message"].(string)
					if msgType == "" {
						msgType = "Throwable"
					}
					if message == "" {
						message = "Error"
					}
					return fmt.Sprintf("%s: %s", msgType, message), nil
				}
				return "Throwable: Error", nil
			}, []string{}).
		AddBuiltinMethod("getMessage", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{},
			func(callEnv *common.Env, args []any) (any, error) {
				thisVal, _ := callEnv.Get("this")
				if instance, ok := thisVal.(*common.ClassInstance); ok {
					return instance.Fields["message"], nil
				}
				return "Error", nil
			}, []string{}).
		AddBuiltinMethod("getType", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{},
			func(callEnv *common.Env, args []any) (any, error) {
				thisVal, _ := callEnv.Get("this")
				if instance, ok := thisVal.(*common.ClassInstance); ok {
					return instance.Fields["type"], nil
				}
				return "Throwable", nil
			}, []string{}).
		BuildAndGet(env)

	if err != nil {
		return err
	}
	exceptionClasses["Throwable"] = throwableConstructor

	// Create RuntimeError class
	_, runtimeErrorConstructor, err := NewClassBuilder("RuntimeError").
		SetParent(throwableClass).
		SetBuiltinConstructor(
			[]ast.Parameter{{Name: "message", Type: ast.TypeFromString("string")}},
			func(callEnv *common.Env, args []any) (any, error) {
				// Call parent constructor through super()
				if classDef, exists := builtinClasses["Throwable"]; exists {
					thisVal, _ := callEnv.Get("this")
					if instance, ok := thisVal.(*common.ClassInstance); ok {
						_, err := callParentConstructor(instance, classDef, callEnv, args)
						if err != nil {
							return nil, err
						}
						// Override type to RuntimeError
						instance.Fields["type"] = "RuntimeError"
					}
				}
				return nil, nil
			},
		).
		BuildAndGet(env)

	if err != nil {
		return err
	}
	exceptionClasses["RuntimeError"] = runtimeErrorConstructor

	// Create TypeError class
	_, typeErrorConstructor, err := NewClassBuilder("TypeError").
		SetParent(builtinClasses["RuntimeError"]).
		SetBuiltinConstructor(
			[]ast.Parameter{{Name: "message", Type: ast.TypeFromString("string")}},
			func(callEnv *common.Env, args []any) (any, error) {
				// Call parent constructor through super()
				if classDef, exists := builtinClasses["RuntimeError"]; exists {
					thisVal, _ := callEnv.Get("this")
					if instance, ok := thisVal.(*common.ClassInstance); ok {
						_, err := callParentConstructor(instance, classDef, callEnv, args)
						if err != nil {
							return nil, err
						}
						// Override type to TypeError
						instance.Fields["type"] = "TypeError"
					}
				}
				return nil, nil
			},
		).
		BuildAndGet(env)

	if err != nil {
		return err
	}
	exceptionClasses["TypeError"] = typeErrorConstructor

	// Create ArityError class
	_, arityErrorConstructor, err := NewClassBuilder("ArityError").
		SetParent(builtinClasses["RuntimeError"]).
		AddField("expected", &ast.Type{Name: "int", IsBuiltin: true}, []string{}).
		AddField("got", &ast.Type{Name: "int", IsBuiltin: true}, []string{}).
		SetBuiltinConstructor(
			[]ast.Parameter{
				{Name: "expected", Type: ast.TypeFromString("int")},
				{Name: "got", Type: ast.TypeFromString("int")},
			},
			func(callEnv *common.Env, args []any) (any, error) {
				expected := 0
				got := 0
				if len(args) > 0 {
					if e, ok := args[0].(float64); ok {
						expected = int(e)
					}
				}
				if len(args) > 1 {
					if g, ok := args[1].(float64); ok {
						got = int(g)
					}
				}

				message := fmt.Sprintf("arity mismatch: expected %d, got %d", expected, got)

				// Call parent constructor through super()
				if classDef, exists := builtinClasses["RuntimeError"]; exists {
					thisVal, _ := callEnv.Get("this")
					if instance, ok := thisVal.(*common.ClassInstance); ok {
						_, err := callParentConstructor(instance, classDef, callEnv, []any{message})
						if err != nil {
							return nil, err
						}
						// Set fields
						instance.Fields["type"] = "ArityError"
						instance.Fields["expected"] = expected
						instance.Fields["got"] = got
					}
				}
				return nil, nil
			},
		).
		BuildAndGet(env)

	if err != nil {
		return err
	}
	exceptionClasses["ArityError"] = arityErrorConstructor

	return nil
}

// loadExceptionClasses loads the built-in exception classes
func loadExceptionClasses(env *Env) error {
	return createExceptionClassesProgrammatically(env)
}

// ThrowRuntimeError throws a RuntimeError exception
// Position information (file, line, column) is automatically retrieved from env
func ThrowRuntimeError(env *Env, message string) error {
	exc := &HyException{
		Message: message,
		Type:    "RuntimeError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowTypeError throws a TypeError exception
// Position information (file, line, column) is automatically retrieved from env
func ThrowTypeError(env *Env, expected string, got ...any) error {
	message := fmt.Sprintf("expected %s", expected)
	if len(got) > 0 {
		message += fmt.Sprintf(", got %T", got[0])
		if len(got) > 1 {
			for _, g := range got[1:] {
				message += fmt.Sprintf(", %T", g)
			}
		}
	}

	exc := &HyException{
		Message: message,
		Type:    "TypeError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["TypeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowArityError throws an ArityError exception
func ThrowArityError(env *Env, expected, got int) error {
	message := fmt.Sprintf("Arity mismatch: expected %d, got %d", expected, got)

	exc := &HyException{
		Message: message,
		Type:    "ArityError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["ArityError"]; exists {
		instance, err := constructor(env, []any{float64(expected), float64(got)})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowAttributeError throws an AttributeError exception for missing attributes/fields
func ThrowAttributeError(env *Env, attrName string, typeName string) error {
	message := fmt.Sprintf("'%s' object has no attribute '%s'", typeName, attrName)

	// Generate hint
	hintProvider := NewHintProvider(env)
	hint := hintProvider.GetHintForAttribute(attrName, typeName)

	// Create exception with location from env
	exc := &HyException{
		Message: message,
		Type:    "AttributeError",
		Hint:    hint,
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowNameError throws a NameError exception for undefined names/identifiers
func ThrowNameError(env *Env, name string) error {
	message := fmt.Sprintf("name '%s' is not defined", name)

	// Generate hint with file context from env
	hintProvider := NewHintProvider(env)
	var hint *ExceptionHint
	if env != nil {
		// Use the code context from env for better hints
		hint = hintProvider.GetHintForUndefinedNameWithEnv(name)
	} else {
		hint = hintProvider.GetHintForUndefinedName(name)
	}

	// Create exception with location from env
	exc := &HyException{
		Message: message,
		Type:    "NameError",
		Hint:    hint,
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowValueError throws a ValueError exception for invalid values
func ThrowValueError(env *Env, message string) error {
	// Check if this is an enum value error and generate hint
	var hint *ExceptionHint
	if strings.Contains(message, "not found in enum") {
		// Extract enum name and value from message
		// Format: "enum value 'X' not found in enum 'Y'"
		parts := strings.Split(message, "'")
		if len(parts) >= 4 {
			valueName := parts[1]
			enumName := parts[3]
			hintProvider := NewHintProvider(env)
			hint = hintProvider.GetHintForEnumValue(enumName, valueName)
		}
	}

	// Create exception with location from env
	exc := &HyException{
		Message: message,
		Type:    "ValueError",
		Hint:    hint,
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowNotCallableError throws an error when trying to call a non-callable object
func ThrowNotCallableError(env *Env, objectType string, value string) error {
	message := fmt.Sprintf("'%s' object is not callable", objectType)
	if value != "" {
		message = fmt.Sprintf("'%s' object with value '%s' is not callable", objectType, value)
	}

	exc := &HyException{
		Message: message,
		Type:    "TypeError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["TypeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowAccessError throws an error when trying to access private/protected members
func ThrowAccessError(env *Env, memberName string, containerName string, accessLevel string) error {
	message := fmt.Sprintf("cannot access %s member '%s' of '%s'", accessLevel, memberName, containerName)

	exc := &HyException{
		Message: message,
		Type:    "AccessError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowNotImplementedError throws an error for unimplemented features
func ThrowNotImplementedError(env *Env, feature string) error {
	message := fmt.Sprintf("%s is not implemented", feature)

	exc := &HyException{
		Message: message,
		Type:    "NotImplementedError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowIndexError throws an IndexError exception for out-of-bounds access
func ThrowIndexError(env *Env, index, size int, collectionType string) error {
	message := fmt.Sprintf("index out of bounds: %d (size: %d)", index, size)

	// Generate hint
	hintProvider := NewHintProvider(env)
	var hint *ExceptionHint

	if index < 0 {
		hint = &ExceptionHint{
			Message:     "Index cannot be negative.",
			Suggestions: []string{fmt.Sprintf("Valid indices are 0 to %d", size-1)},
			HintType:    "general",
		}
	} else if index >= size {
		if size == 0 {
			hint = &ExceptionHint{
				Message:     "Cannot access elements in an empty collection.",
				Suggestions: []string{"Add elements first using .add() method"},
				HintType:    "general",
			}
		} else {
			hint = &ExceptionHint{
				Message:     "Index is beyond the collection size.",
				Suggestions: []string{fmt.Sprintf("Valid indices are 0 to %d", size-1)},
				HintType:    "general",
			}
		}
	}

	_ = hintProvider // suppress unused warning for now

	exc := &HyException{
		Message: message,
		Type:    "IndexError",
		Hint:    hint,
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowConversionError throws a ConversionError exception for type conversion failures
func ThrowConversionError(env *Env, value any, targetType string) error {
	message := fmt.Sprintf("cannot convert %T to %s", value, targetType)
	if strVal, ok := value.(string); ok {
		message = fmt.Sprintf("cannot convert string '%s' to %s", strVal, targetType)
	}

	exc := &HyException{
		Message: message,
		Type:    "ConversionError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["TypeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowStateError throws a StateError exception for invalid state operations
func ThrowStateError(env *Env, message string) error {
	exc := &HyException{
		Message: message,
		Type:    "StateError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ThrowInitializationError throws an InitializationError exception
func ThrowInitializationError(env *Env, what string) error {
	message := fmt.Sprintf("%s not initialized", what)

	exc := &HyException{
		Message: message,
		Type:    "InitializationError",
	}
	if env != nil {
		exc.File = env.GetFileName()
		exc.Line = env.GetCurrentLine()
		exc.Column = env.CurrentColumn
	}

	if constructor, exists := exceptionClasses["RuntimeError"]; exists {
		instance, err := constructor(env, []any{message})
		if err == nil {
			exc.Instance = instance
		}
	}

	return exc
}

// ValidateArgumentType validates that an argument matches the expected type
func ValidateArgumentType(value any, expectedType string) error {
	if expectedType == "" {
		return nil // No type constraint
	}

	// Check if it's a union type (contains |)
	if strings.Contains(expectedType, "|") {
		// Parse union type and check if value matches any member
		unionTypeStrs := strings.Split(expectedType, "|")
		for _, typeStr := range unionTypeStrs {
			typeStr = strings.TrimSpace(typeStr)
			if IsInstanceOf(value, typeStr) {
				return nil // Value matches at least one union member
			}
		}
		// Value doesn't match any union member
		actualType := GetTypeName(value)
		return ThrowTypeError(nil, expectedType, actualType)
	}

	if IsInstanceOf(value, expectedType) {
		return nil
	}
	fmt.Println("Validation failed:", GetTypeName(value), "is not", expectedType)

	actualType := GetTypeName(value)
	return ThrowTypeError(nil, expectedType, actualType)
}

// ValidateVariadicArguments validates variadic arguments and converts them to an array
func ValidateVariadicArguments(args []any, expectedType string) ([]any, error) {
	result := make([]any, len(args))

	for i, arg := range args {
		if err := ValidateArgumentType(arg, expectedType); err != nil {
			return nil, err
		}
		result[i] = arg
	}

	return result, nil
}

// ValidateFunctionArguments validates all function arguments including variadic ones
func ValidateFunctionArguments(args []any, paramTypes []string, hasVariadic bool, variadicType string) ([]any, error) {
	regularParamCount := len(paramTypes)
	if hasVariadic {
		regularParamCount-- // Last param is variadic
	}

	// Check minimum argument count
	if len(args) < regularParamCount {
		return nil, common.ArityError{Expected: regularParamCount, Got: len(args)}
	}

	// If no variadic and exact match required
	if !hasVariadic && len(args) != regularParamCount {
		return nil, common.ArityError{Expected: regularParamCount, Got: len(args)}
	}

	result := make([]any, regularParamCount)

	// Validate regular parameters
	for i := 0; i < regularParamCount; i++ {
		if err := ValidateArgumentType(args[i], paramTypes[i]); err != nil {
			return nil, err
		}
		result[i] = args[i]
	}

	// Handle variadic arguments
	if hasVariadic {
		variadicArgs := args[regularParamCount:]
		validatedVariadic, err := ValidateVariadicArguments(variadicArgs, variadicType)
		if err != nil {
			return nil, err
		}

		// Add variadic arguments as an array at the end
		result = append(result, validatedVariadic)
	}

	return result, nil
}

// evalTryStmt handles try-catch-finally statements
func evalTryStmt(env *Env, stmt *ast.TryStmt) (val any, returned bool, err error) {
	var lastValue any
	var caughtException *HyException

	// Execute try block
	for _, st := range stmt.Body {
		v, ret, err := evalStmt(env, st)
		if err != nil {
			// Check if it's a HyException
			if hyErr, ok := err.(*HyException); ok {
				caughtException = hyErr
				break
			} else {
				// Convert regular error to HyException
				caughtException = NewHyException("RuntimeError", err.Error())
				break
			}
		}
		if ret {
			// If we have a return, we need to handle it after finally
			lastValue = v
			returned = true
			break
		}
		lastValue = v
	}

	// Handle catch blocks if there was an exception
	if caughtException != nil {
		handled := false
		for _, catch := range stmt.Catches {
			// Check if exception type matches (if specified)
			if catch.ExceptType == "" || catch.ExceptType == caughtException.Type {
				// Create new scope for catch block
				catchEnv := &Env{Parent: env, Vars: map[string]any{}, Consts: map[string]bool{}}

				// Bind exception variable if specified
				if catch.VarName != "" {
					catchEnv.Define(catch.VarName, caughtException.Instance, catch.Modifier)
				}

				// Execute catch block
				for _, st := range catch.Body {
					v, ret, err := evalStmt(catchEnv, st)
					if err != nil {
						return nil, false, err
					}
					if ret {
						lastValue = v
						returned = true
						break
					}
					lastValue = v
				}
				handled = true
				break
			}
		}

		// If exception wasn't handled, re-throw it
		if !handled {
			return nil, false, caughtException
		}
	}

	// Execute finally block
	if len(stmt.Finally) > 0 {
		for _, st := range stmt.Finally {
			_, ret, err := evalStmt(env, st)
			if err != nil {
				return nil, false, err
			}
			if ret {
				// Finally block return overrides everything
				return nil, true, nil
			}
		}
	}

	return lastValue, returned, nil
}

// evalThrowStmt handles throw statements
func evalThrowStmt(env *Env, stmt *ast.ThrowStmt) (val any, returned bool, err error) {
	// Evaluate the expression to throw
	value, err := evalExpr(env, stmt.Value)
	if err != nil {
		return nil, false, err
	}
	et := &HyException{
		Message:    fmt.Sprintf("unknown error: %v", value),
		Type:       "UnknownError",
		StackTrace: []string{},
		Instance:   value,
	}
	switch val := value.(type) {
	case *ClassInstance:
		if msg, ok := val.Fields["message"].(string); ok {
			typ, ok := val.Fields["type"].(string)
			if !ok {
				break
			}

			var stackTrace []string
			if st, ok := val.Fields["stackTrace"].([]string); ok {
				stackTrace = st
			}
			et = &HyException{
				Message:    msg,
				Type:       typ,
				StackTrace: stackTrace,
				Instance:   val,
			}
		}
	}

	return nil, false, et
}

// evalDeferStmt handles defer statements
func evalDeferStmt(env *Env, stmt *ast.DeferStmt) (val any, returned bool, err error) {
	// Defer statements schedule a function call to be executed when the current
	// function/scope exits. We save a closure that will evaluate the call later.
	deferredCall := func() error {
		_, evalErr := evalExpr(env, stmt.Call)
		return evalErr
	}

	// Add to defer stack (will be executed in LIFO order)
	env.Defers = append(env.Defers, deferredCall)
	return nil, false, nil
}

// InstallExceptionBuiltins installs exception-related built-ins
func InstallExceptionBuiltins(env *Env) {
	// Load exception classes
	loadExceptionClasses(env)
}

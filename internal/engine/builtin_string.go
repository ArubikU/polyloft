package engine

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallStringBuiltin installs String builtin type as a class
func InstallStringBuiltin(env *Env) error {
	// Step 1: Create basic class structure first with fields
	stringClass := NewClassBuilder("String")

	// Add field with native string type
	nativeStringType := &ast.Type{Name: "string", IsBuiltin: true}
	stringClass.AddField("_value", nativeStringType, []string{"private"})

	// Step 2: Now get type references for method signatures
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	stringType := stringClass.GetType()
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)

	stringClass.AddField("_currentIndex", intType, []string{"private"})

	// Step 3: Get interface implementations
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)
	stringClass.AddInterface(iterableInterface)

	//hasNext() -> Bool
	stringClass.AddBuiltinMethod("hasNext", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		return instance.Fields["_currentIndex"].(int) < utf8.RuneCountInString(instance.Fields["_value"].(string)), nil
	}, []string{})
	//next() -> String
	stringClass.AddBuiltinMethod("next", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		currentIndex := instance.Fields["_currentIndex"].(int)
		runes := []rune(str)
		if currentIndex >= len(runes) {
			return CreateStringInstance((*Env)(callEnv), "")
		}
		currentIndex++
		return CreateStringInstance((*Env)(callEnv), string(runes[currentIndex-1]))
	}, []string{})
	// length() -> Int
	stringClass.AddBuiltinMethod("length", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return utf8.RuneCountInString(str), nil
	}, []string{})

	// isEmpty() -> Bool
	stringClass.AddBuiltinMethod("isEmpty", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return len(str) == 0, nil
	}, []string{})

	// charAt(index: Int) -> String
	stringClass.AddBuiltinMethod("charAt", stringType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)

		index, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		runes := []rune(str)
		if index < 0 || index >= len(runes) {
			return CreateStringInstance((*Env)(callEnv), "")
		}

		return CreateStringInstance((*Env)(callEnv), string(runes[index]))
	}, []string{})

	// indexOf(substr: String) -> Int
	stringClass.AddBuiltinMethod("indexOf", intType, []ast.Parameter{
		{Name: "substr", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)

		substr := StringValue(args[0])
		index := strings.Index(str, substr)
		return index, nil
	}, []string{})

	// substring(start: Int, end: Int) -> String
	stringClass.AddBuiltinMethod("substring", stringType, []ast.Parameter{
		{Name: "start", Type: intType},
		{Name: "end", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)

		start, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		runes := []rune(str)
		length := len(runes)

		if start < 0 {
			start = 0
		}
		if start > length {
			start = length
		}

		end := length
		if len(args) >= 2 {
			if endArg, ok := utils.AsInt(args[1]); ok {
				end = endArg
			}
		}

		if end < 0 {
			end = 0
		}
		if end > length {
			end = length
		}
		if end < start {
			end = start
		}

		return CreateStringInstance((*Env)(callEnv), string(runes[start:end]))
	}, []string{})

	// toUpperCase() -> String
	stringClass.AddBuiltinMethod("toUpperCase", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return CreateStringInstance((*Env)(callEnv), strings.ToUpper(str))
	}, []string{})

	// toLowerCase() -> String
	stringClass.AddBuiltinMethod("toLowerCase", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return CreateStringInstance((*Env)(callEnv), strings.ToLower(str))
	}, []string{})

	// trim() -> String
	stringClass.AddBuiltinMethod("trim", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return CreateStringInstance((*Env)(callEnv), strings.TrimSpace(str))
	}, []string{})

	// startsWith(prefix: String) -> Bool
	stringClass.AddBuiltinMethod("startsWith", boolType, []ast.Parameter{
		{Name: "prefix", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		prefix := StringValue(args[0])
		return strings.HasPrefix(str, prefix), nil
	}, []string{})

	// endsWith(suffix: String) -> Bool
	stringClass.AddBuiltinMethod("endsWith", boolType, []ast.Parameter{
		{Name: "suffix", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		suffix := StringValue(args[0])
		return strings.HasSuffix(str, suffix), nil
	}, []string{})

	// contains(substr: String) -> Bool
	stringClass.AddBuiltinMethod("contains", boolType, []ast.Parameter{
		{Name: "substr", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		substr := StringValue(args[0])
		return strings.Contains(str, substr), nil
	}, []string{})

	// replace(old: String, new: String) -> String
	stringClass.AddBuiltinMethod("replace", stringType, []ast.Parameter{
		{Name: "old", Type: stringType},
		{Name: "new", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		old := StringValue(args[0])
		new := StringValue(args[1])
		return CreateStringInstance((*Env)(callEnv), strings.ReplaceAll(str, old, new))
	}, []string{})

	// split(sep: String) -> Array
	stringClass.AddBuiltinMethod("split", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{
		{Name: "sep", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		sep := StringValue(args[0])
		parts := strings.Split(str, sep)
		result := make([]any, len(parts))
		for i, p := range parts {
			strInstance, err := CreateStringInstance((*Env)(callEnv), p)
			if err != nil {
				return nil, err
			}
			result[i] = strInstance
		}
		return CreateArrayInstance((*Env)(callEnv), result)
	}, []string{})

	// utils.ToString() -> String
	stringClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return CreateStringInstance((*Env)(callEnv), str)
	}, []string{})

	// repeat(count: Int) -> String
	stringClass.AddBuiltinMethod("repeat", intType, []ast.Parameter{
		{Name: "count", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		count, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		return CreateStringInstance((*Env)(callEnv), strings.Repeat(str, count))
	}, []string{})

	// padStart(length: Int, pad: String) -> String
	stringClass.AddBuiltinMethod("padStart", stringType, []ast.Parameter{
		{Name: "length", Type: intType},
		{Name: "pad", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		length, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		pad := StringValue(args[1])

		if len(str) >= length {
			return CreateStringInstance((*Env)(callEnv), str)
		}
		padding := strings.Repeat(pad, (length-len(str)+len(pad)-1)/len(pad))
		return CreateStringInstance((*Env)(callEnv), padding[:length-len(str)]+str)
	}, []string{})

	// padEnd(length: Int, pad: String) -> String
	stringClass.AddBuiltinMethod("padEnd", stringType, []ast.Parameter{
		{Name: "length", Type: intType},
		{Name: "pad", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		length, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		pad := StringValue(args[1])

		if len(str) >= length {
			return CreateStringInstance((*Env)(callEnv), str)
		}
		padding := strings.Repeat(pad, (length-len(str)+len(pad)-1)/len(pad))
		return CreateStringInstance((*Env)(callEnv), str+padding[:length-len(str)])
	}, []string{})

	// serialize() -> String
	stringClass.AddBuiltinMethod("serialize", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		str := instance.Fields["_value"].(string)
		return CreateStringInstance((*Env)(callEnv), strconv.Quote(str))
	}, []string{})

	// Build the class
	_, err := stringClass.Build(env)
	return err
}

// CreateStringInstance creates a String instance from a Go string
// This is used when evaluating string literals
func CreateStringInstance(env *Env, value string) (*ClassInstance, error) {
	stringClassVal, ok := env.Get("__StringClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "String class")
	}

	stringClass := stringClassVal.(*ClassDefinition)

	// Create instance
	instance, err := createClassInstance(stringClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_value"] = value

	return classInstance, nil
}

// StringValue extracts the Go string value from a String instance or converts a value to string
func StringValue(v any) string {
	switch val := v.(type) {
	case *ClassInstance:
		if val.ClassName == "String" {
			if strVal, ok := val.Fields["_value"].(string); ok {
				return strVal
			}
		}
		return utils.ToString(v)
	case string:
		return val
	default:
		return utils.ToString(v)
	}
}

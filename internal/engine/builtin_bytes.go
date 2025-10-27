package engine

import (
	"encoding/hex"
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallBytesBuiltin installs the Bytes builtin class for handling byte arrays
func InstallBytesBuiltin(env *Env) error {
	// Get type references
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	floatType := common.BuiltinTypeFloat.GetTypeDefinition(env)

	// Create Bytes class
	bytesBuilder := NewClassBuilder("Bytes").
		AddField("_data", arrayType, []string{"private"})

	bytesType := bytesBuilder.GetType()

	// Constructor: Bytes() - empty bytes
	bytesBuilder.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_data"] = []byte{}
		return nil, nil
	})

	// Constructor: Bytes(size) - create byte array of specific size
	bytesBuilder.AddBuiltinConstructor([]ast.Parameter{
		{Name: "size", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		size, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		instance.Fields["_data"] = make([]byte, size)
		return nil, nil
	})

	// Constructor: Bytes(data: Array) - from array of ints
	bytesBuilder.AddBuiltinConstructor([]ast.Parameter{
		{Name: "data", Type: arrayType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		byteData, _ := AsBytes((*common.Env)(callEnv), args[0])
		instance.Fields["_data"] = byteData

		return nil, nil
	})

	// size() -> Int
	bytesBuilder.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		return len(data), nil
	}, []string{})

	// get(index: Int) -> Int
	bytesBuilder.AddBuiltinMethod("get", intType, []ast.Parameter{
		{Name: "index", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		if idx < 0 || idx >= len(data) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(data), "Bytes")
		}
		return int(data[idx]), nil
	}, []string{})

	// set(index: Int, value: Int) -> Void
	bytesBuilder.AddBuiltinMethod("set", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "index", Type: intType},
		{Name: "value", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		idx, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		val, ok := utils.AsInt(args[1])
		if !ok || val < 0 || val > 255 {
			return nil, ThrowValueError((*Env)(callEnv), "byte value must be 0-255")
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		if idx < 0 || idx >= len(data) {
			return nil, ThrowIndexError((*Env)(callEnv), idx, len(data), "Bytes")
		}
		data[idx] = byte(val)
		return nil, nil
	}, []string{})

	// toString() -> String
	bytesBuilder.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		binStr := ""
		for _, b := range data {
			binStr += fmt.Sprintf("%08b", b)
		}
		return CreateStringInstance(env, "0b"+binStr)
	}, []string{})

	// asHex() -> String
	bytesBuilder.AddBuiltinMethod("asHex", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		return hex.EncodeToString(data), nil
	}, []string{})

	// asString() -> String - converts bytes to UTF-8 string
	bytesBuilder.AddBuiltinMethod("asString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		return CreateStringInstance(env, string(data))
	}, []string{})

	// asInt() -> Int - converts first byte to int (error if empty)
	bytesBuilder.AddBuiltinMethod("asInt", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		if len(data) == 0 {
			return nil, ThrowValueError((*Env)(callEnv), "cannot convert empty Bytes to Int")
		}
		return int(data[0]), nil
	}, []string{})

	// asFloat() -> Float - converts first byte to float (error if empty)
	bytesBuilder.AddBuiltinMethod("asFloat", floatType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		if len(data) == 0 {
			return nil, ThrowValueError((*Env)(callEnv), "cannot convert empty Bytes to Float")
		}
		return float64(data[0]), nil
	}, []string{})

	// asBool() -> Bool - returns true if any byte is non-zero
	bytesBuilder.AddBuiltinMethod("asBool", boolType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		for _, b := range data {
			if b != 0 {
				return true, nil
			}
		}
		return false, nil
	}, []string{})

	// asHex() -> String - converts bytes to hexadecimal string with 0x prefix
	bytesBuilder.AddBuiltinMethod("asHex", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		return CreateStringInstance(env, "0x"+hex.EncodeToString(data))
	}, []string{})

	// asBinary() -> String - converts bytes to binary string with 0b prefix
	bytesBuilder.AddBuiltinMethod("asBinary", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		binStr := ""
		for _, b := range data {
			binStr += fmt.Sprintf("%08b", b)
		}
		return CreateStringInstance(env, "0b"+binStr)
	}, []string{})

	// asArray() -> Array - same as toArray() but for naming consistency
	bytesBuilder.AddBuiltinMethod("asArray", arrayType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		result := make([]any, len(data))
		for i, b := range data {
			result[i] = int(b)
		}
		return CreateArrayInstance(callEnv, result)
	}, []string{})

	// slice(start: Int, end: Int) -> Bytes
	bytesBuilder.AddBuiltinMethod("slice", bytesType, []ast.Parameter{
		{Name: "start", Type: intType},
		{Name: "end", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		start, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		end, ok := utils.AsInt(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[1])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)

		if start < 0 || start > len(data) || end < start || end > len(data) {
			return nil, ThrowIndexError((*Env)(callEnv), start, len(data), "Bytes")
		}

		// Create new Bytes instance
		bytesInst, err := CreateBytesInstance((*common.Env)(callEnv), data[start:end])
		if err != nil {
			return nil, err
		}
		return bytesInst, nil
	}, []string{})

	// equals(other: Bytes) -> Bool
	bytesBuilder.AddBuiltinMethod("equals", boolType, []ast.Parameter{
		{Name: "other", Type: bytesType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)

		if otherInst, ok := args[0].(*ClassInstance); ok {
			otherData := otherInst.Fields["_data"].([]byte)
			if len(data) != len(otherData) {
				return false, nil
			}
			for i := range data {
				if data[i] != otherData[i] {
					return false, nil
				}
			}
			return true, nil
		}
		return false, nil
	}, []string{})

	// --- STATIC METHODS ---

	// fromString(str: String)
	bytesBuilder.AddStaticMethod("fromString", bytesType, []ast.Parameter{
		{Name: "str", Type: stringType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		return CreateBytesInstance(env, args[0])
	}))

	// fromArray(arr: Array)
	bytesBuilder.AddStaticMethod("fromArray", bytesType, []ast.Parameter{
		{Name: "arr", Type: arrayType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		return CreateBytesInstance(env, args[0])
	}))

	// fromHex(hex: String)
	bytesBuilder.AddStaticMethod("fromHex", bytesType, []ast.Parameter{
		{Name: "hex", Type: stringType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		data, err := hex.DecodeString(utils.ToString(args[0]))
		if err != nil {
			return nil, err
		}
		return CreateBytesInstance(env, data)
	}))

	// fromFloat(num: Float)
	bytesBuilder.AddStaticMethod("fromFloat", bytesType, []ast.Parameter{
		{Name: "num", Type: floatType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		return CreateBytesInstance(env, args[0])
	}))

	// fromBool(flag: Bool)
	bytesBuilder.AddStaticMethod("fromBool", bytesType, []ast.Parameter{
		{Name: "flag", Type: boolType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		return CreateBytesInstance(env, args[0])
	}))

	_, err := bytesBuilder.Build(env)
	return err
}

func AsBytes(env *common.Env, value any) ([]byte, bool) {
	// Use the unified type converter
	result, ok := ConvertTo(env, "Bytes", value)
	if !ok {
		return []byte{}, false
	}
	if byteSlice, ok := result.([]byte); ok {
		return byteSlice, true
	}
	return []byte{}, false
}

func CreateBytesInstance(env *common.Env, data any) (*ClassInstance, error) {
	bytesClassDef := common.BuiltinTypeBytes.GetClassDefinition(env)
	if bytesClassDef == nil {
		return nil, ThrowInitializationError(env, "Bytes class not found")
	}

	instance, err := createClassInstance(bytesClassDef, env, []any{})
	if err != nil {
		return nil, err
	}

	ins := instance.(*ClassInstance)

	// Usamos AsBytes para convertir el dato
	byteData, ok := AsBytes(env, data)
	if !ok {
		byteData = []byte{}
	}

	ins.Fields["_data"] = byteData
	return ins, nil
}

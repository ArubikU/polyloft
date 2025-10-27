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
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_data"] = []byte{}
		return nil, nil
	})

	// Constructor: Bytes(size) - create byte array of specific size
	bytesBuilder.AddBuiltinConstructor([]ast.Parameter{
		{Name: "size", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)

		byteData, _ := AsBytes((*common.Env)(callEnv), args[0])
		instance.Fields["_data"] = byteData

		return nil, nil
	})

	// size() -> Int
	bytesBuilder.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		binStr := ""
		for _, b := range data {
			binStr += fmt.Sprintf("%08b", b)
		}
		return CreateStringInstance(env, "0b"+binStr)
	}, []string{})

	// toHex() -> String
	bytesBuilder.AddBuiltinMethod("toHex", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].([]byte)
		return hex.EncodeToString(data), nil
	}, []string{})

	// toArray() -> Array
	bytesBuilder.AddBuiltinMethod("toArray", arrayType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
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
		thisVal, _ := callEnv.Get("this")
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
	if value == nil {
		return []byte{}, false
	}

	switch v := value.(type) {
	case string:

		//check if its a hex with 0x prefix
		if len(v) >= 2 && v[0:2] == "0x" {
			decoded, err := hex.DecodeString(v[2:])
			if err == nil {
				return decoded, true
			}
		}
		//check if its a binary with 0b prefix
		if len(v) >= 2 && v[0:2] == "0b" {
			binStr := v[2:]
			byteLen := (len(binStr) + 7) / 8
			result := make([]byte, byteLen)
			for i, char := range binStr {
				if char == '1' {
					byteIndex := i / 8
					bitIndex := 7 - (i % 8)
					result[byteIndex] |= 1 << bitIndex
				}
			}
			return result, true
		}

		return []byte(v), true
	case []byte:
		return v, true
	case float32:
		return []byte{byte(v)}, true
	case float64:
		return []byte{byte(v)}, true
	case int:
		return []byte{byte(v)}, true
	case int8:
		return []byte{byte(v)}, true
	case int16:
		return []byte{byte(v)}, true
	case int32:
		return []byte{byte(v)}, true
	case int64:
		return []byte{byte(v)}, true
	case uint:
		return []byte{byte(v)}, true
	case uint8:
		return []byte{byte(v)}, true
	case uint16:
		return []byte{byte(v)}, true
	case uint32:
		return []byte{byte(v)}, true
	case uint64:
		return []byte{byte(v)}, true
	case bool:
		if v {
			return []byte{1}, true
		}
	case *ClassInstance:
		stringDef := common.BuiltinTypeString.GetClassDefinition(env)
		intDef := common.BuiltinTypeInt.GetClassDefinition(env)
		arrayDef := common.BuiltinTypeArray.GetClassDefinition(env)
		boolDef := common.BuiltinTypeBool.GetClassDefinition(env)
		floatDef := common.BuiltinTypeFloat.GetClassDefinition(env)
		listDef := common.BuiltinTypeList.GetClassDefinition(env)
		mapDef := common.BuiltinTypeMap.GetClassDefinition(env)
		tupleDef := common.BuiltinTypeTuple.GetClassDefinition(env)
		setDef := common.BuiltinTypeSet.GetClassDefinition(env)
		dequeDef := common.BuiltinTypeDeque.GetClassDefinition(env)
		rangeDef := common.BuiltinTypeRange.GetClassDefinition(env)
		bytesDef := common.BuiltinTypeBytes.GetClassDefinition(env)

		if v.ParentClass.IsSubclassOf(stringDef) {
			content := utils.ToString(v)
			//check if its a hex with 0x prefix
			if len(content) >= 2 && content[0:2] == "0x" {
				decoded, err := hex.DecodeString(content[2:])
				if err == nil {
					return decoded, true
				}
			}
			//check if its a binary with 0b prefix like 0b101010
			if len(content) >= 2 && content[0:2] == "0b" {
				binStr := content[2:]
				byteLen := (len(binStr) + 7) / 8
				result := make([]byte, byteLen)
				for i, char := range binStr {
					if char == '1' {
						byteIndex := i / 8
						bitIndex := 7 - (i % 8)
						result[byteIndex] |= 1 << bitIndex
					}
				}
				return result, true
			}

			return []byte(utils.ToString(v)), true
		} else if v.ParentClass.IsSubclassOf(arrayDef) {
			items := v.Fields["items"].([]any)
			buf := []byte{}
			for _, item := range items {
				b, ok := AsBytes(env, item)
				if ok {
					buf = append(buf, b...)
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(listDef) {
			// List: _items is a pointer to []any
			itemsPtr := v.Fields["_items"].(*[]any)
			buf := []byte{}
			for _, item := range *itemsPtr {
				b, ok := AsBytes(env, item)
				if ok {
					buf = append(buf, b...)
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(dequeDef) {
			// Deque: _items is a pointer to []any
			itemsPtr := v.Fields["_items"].(*[]any)
			buf := []byte{}
			for _, item := range *itemsPtr {
				b, ok := AsBytes(env, item)
				if ok {
					buf = append(buf, b...)
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(tupleDef) {
			// Tuple: _elements is []any
			elements := v.Fields["_elements"].([]any)
			buf := []byte{}
			for _, item := range elements {
				b, ok := AsBytes(env, item)
				if ok {
					buf = append(buf, b...)
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(setDef) {
			// Set: _keys is a pointer to []any (ordered keys)
			keysPtr := v.Fields["_keys"].(*[]any)
			buf := []byte{}
			for _, item := range *keysPtr {
				b, ok := AsBytes(env, item)
				if ok {
					buf = append(buf, b...)
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(mapDef) {
			// Map: _keys is a pointer to []any
			keysPtr := v.Fields["_keys"].(*[]any)
			buf := []byte{}
			for _, key := range *keysPtr {
				// Convert key to bytes
				keyBytes, ok := AsBytes(env, key)
				if ok {
					buf = append(buf, keyBytes...)
				}
				// Convert value to bytes
				entriesPtr := v.Fields["_entries"].(*[]ast.MapEntry)
				for _, entry := range *entriesPtr {
					if utils.ToString(entry.Key) == utils.ToString(key) {
						valBytes, ok := AsBytes(env, entry.Value)
						if ok {
							buf = append(buf, valBytes...)
						}
						break
					}
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(rangeDef) {
			// Range: convert range to bytes by iterating through values
			start, _ := v.Fields["_start"].(int)
			end, _ := v.Fields["_end"].(int)
			step, _ := v.Fields["_step"].(int)

			buf := []byte{}
			if step > 0 {
				for i := start; i < end; i += step {
					buf = append(buf, byte(i))
				}
			} else if step < 0 {
				for i := start; i > end; i += step {
					buf = append(buf, byte(i))
				}
			}
			return buf, true
		} else if v.ParentClass.IsSubclassOf(bytesDef) {
			// Bytes: _data is []byte
			data := v.Fields["_data"].([]byte)
			return data, true
		} else if v.ParentClass.IsSubclassOf(intDef) {
			return []byte{byte(v.Fields["_value"].(int))}, true
		} else if v.ParentClass.IsSubclassOf(boolDef) {
			if v.Fields["_value"].(bool) {
				return []byte{1}, true
			}
		} else if v.ParentClass.IsSubclassOf(floatDef) {
			return []byte{byte(v.Fields["_value"].(float64))}, true
		}
	}

	return []byte{}, false
}

func CreateBytesInstance(env *common.Env, data any) (*ClassInstance, error) {
	bytesClass, ok := env.Get("__BytesClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Bytes class not found")
	}

	arrayClassDef := bytesClass.(*common.ClassDefinition)
	instance, err := createClassInstance(arrayClassDef, env, []any{})
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

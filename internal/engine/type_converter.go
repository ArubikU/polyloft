package engine

import (
	"encoding/hex"
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// TypeConverter is a function that converts a value to a specific type
// Returns the converted value and a boolean indicating success
type TypeConverter func(env *common.Env, value any) (any, bool)

// TypeRegistry holds all registered type converters
type TypeRegistry struct {
	converters map[string]TypeConverter
}

var globalTypeRegistry = &TypeRegistry{
	converters: make(map[string]TypeConverter),
}

// RegisterTypeConverter registers a converter for a specific type
func RegisterTypeConverter(typeName string, converter TypeConverter) {
	globalTypeRegistry.converters[typeName] = converter
}

// ConvertTo attempts to convert a value to the specified type using registered converters
func ConvertTo(env *common.Env, typeName string, value any) (any, bool) {
	if converter, ok := globalTypeRegistry.converters[typeName]; ok {
		return converter(env, value)
	}
	return nil, false
}

// GetTypeConverter returns the converter for a specific type
func GetTypeConverter(typeName string) (TypeConverter, bool) {
	converter, ok := globalTypeRegistry.converters[typeName]
	return converter, ok
}

// InitializeBuiltinTypeConverters registers all builtin type converters
func InitializeBuiltinTypeConverters() {
	// Register bytes converter
	RegisterTypeConverter("Bytes", func(env *common.Env, value any) (any, bool) {
		if value == nil {
			return []byte{}, false
		}

		switch v := value.(type) {
		case string:
			// Check if it's a hex with 0x prefix
			if len(v) >= 2 && v[0:2] == "0x" {
				decoded, err := hex.DecodeString(v[2:])
				if err == nil {
					return decoded, true
				}
			}
			// Check if it's a binary with 0b prefix
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
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			intVal, _ := utils.AsInt(v)
			return []byte{byte(intVal)}, true
		case bool:
			if v {
				return []byte{1}, true
			}
			return []byte{0}, true
		case *ClassInstance:
			return convertClassInstanceToBytes(env, v)
		}

		return []byte{}, false
	})

	// Register array converter
	RegisterTypeConverter("Array", func(env *common.Env, value any) (any, bool) {
		if value == nil {
			return []any{}, true
		}

		switch v := value.(type) {
		case []any:
			return v, true
		case *ClassInstance:
			arrayDef := common.BuiltinTypeArray.GetClassDefinition(env)
			if v.ParentClass.IsSubclassOf(arrayDef) {
				if items, ok := v.Fields["_items"].([]any); ok {
					return items, true
				}
			}
		}

		return nil, false
	})

	// Register string converter
	RegisterTypeConverter("String", func(env *common.Env, value any) (any, bool) {
		if value == nil {
			return "", true
		}

		switch v := value.(type) {
		case string:
			return v, true
		case *ClassInstance:
			stringDef := common.BuiltinTypeString.GetClassDefinition(env)
			if v.ParentClass.IsSubclassOf(stringDef) {
				if val, ok := v.Fields["_value"].(string); ok {
					return val, true
				}
			}
		}

		return utils.ToString(value), true
	})

	// Register int converter
	RegisterTypeConverter("Int", func(env *common.Env, value any) (any, bool) {
		return utils.AsInt(value)
	})

	// Register float converter
	RegisterTypeConverter("Float", func(env *common.Env, value any) (any, bool) {
		return utils.AsFloat(value)
	})

	// Register bool converter
	RegisterTypeConverter("Bool", func(env *common.Env, value any) (any, bool) {
		return utils.AsBool(value), true
	})

	// Register map converter
	RegisterTypeConverter("Map", func(env *common.Env, value any) (any, bool) {
		if value == nil {
			return make(map[string]any), true
		}

		switch v := value.(type) {
		case map[string]any:
			return v, true
		case map[any]any:
			result := make(map[string]any)
			for key, val := range v {
				result[utils.ToString(key)] = val
			}
			return result, true
		case *ClassInstance:
			mapDef := common.BuiltinTypeMap.GetClassDefinition(env)
			if v.ParentClass.IsSubclassOf(mapDef) {
				// Extract the map data
				if data, ok := v.Fields["_data"]; ok {
					return data, true
				}
			}
		}

		return nil, false
	})
}

// convertClassInstanceToBytes handles conversion of ClassInstance to bytes
func convertClassInstanceToBytes(env *common.Env, v *ClassInstance) ([]byte, bool) {
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
		// Check if it's a hex with 0x prefix
		if len(content) >= 2 && content[0:2] == "0x" {
			decoded, err := hex.DecodeString(content[2:])
			if err == nil {
				return decoded, true
			}
		}
		// Check if it's a binary with 0b prefix
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
		items := v.Fields["_items"].([]any)
		buf := []byte{}
		for _, item := range items {
			if b, ok := ConvertTo(env, "Bytes", item); ok {
				if byteSlice, ok := b.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(listDef) {
		itemsPtr := v.Fields["_items"].(*[]any)
		buf := []byte{}
		for _, item := range *itemsPtr {
			if b, ok := ConvertTo(env, "Bytes", item); ok {
				if byteSlice, ok := b.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(dequeDef) {
		itemsPtr := v.Fields["_items"].(*[]any)
		buf := []byte{}
		for _, item := range *itemsPtr {
			if b, ok := ConvertTo(env, "Bytes", item); ok {
				if byteSlice, ok := b.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(tupleDef) {
		elements := v.Fields["_elements"].([]any)
		buf := []byte{}
		for _, item := range elements {
			if b, ok := ConvertTo(env, "Bytes", item); ok {
				if byteSlice, ok := b.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(setDef) {
		keysPtr := v.Fields["_keys"].(*[]any)
		buf := []byte{}
		for _, item := range *keysPtr {
			if b, ok := ConvertTo(env, "Bytes", item); ok {
				if byteSlice, ok := b.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(mapDef) {
		keysPtr := v.Fields["_keys"].(*[]any)
		buf := []byte{}
		for _, key := range *keysPtr {
			if keyBytes, ok := ConvertTo(env, "Bytes", key); ok {
				if byteSlice, ok := keyBytes.([]byte); ok {
					buf = append(buf, byteSlice...)
				}
			}
			entriesPtr := v.Fields["_entries"].(*[]ast.MapEntry)
			for _, entry := range *entriesPtr {
				if utils.ToString(entry.Key) == utils.ToString(key) {
					if valBytes, ok := ConvertTo(env, "Bytes", entry.Value); ok {
						if byteSlice, ok := valBytes.([]byte); ok {
							buf = append(buf, byteSlice...)
						}
					}
					break
				}
			}
		}
		return buf, true
	} else if v.ParentClass.IsSubclassOf(rangeDef) {
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
		data := v.Fields["_data"].([]byte)
		return data, true
	} else if v.ParentClass.IsSubclassOf(intDef) {
		return []byte{byte(v.Fields["_value"].(int))}, true
	} else if v.ParentClass.IsSubclassOf(boolDef) {
		if v.Fields["_value"].(bool) {
			return []byte{1}, true
		}
		return []byte{0}, true
	} else if v.ParentClass.IsSubclassOf(floatDef) {
		return []byte{byte(v.Fields["_value"].(float64))}, true
	}

	return []byte{}, false
}

// TypeConversionError represents an error during type conversion
type TypeConversionError struct {
	FromType string
	ToType   string
	Value    any
}

func (e *TypeConversionError) Error() string {
	return fmt.Sprintf("cannot convert %v (type %s) to %s", e.Value, e.FromType, e.ToType)
}

// NewTypeConversionError creates a new type conversion error
func NewTypeConversionError(fromType, toType string, value any) *TypeConversionError {
	return &TypeConversionError{
		FromType: fromType,
		ToType:   toType,
		Value:    value,
	}
}

// InstanceCreator is a function that creates a ClassInstance from a native value
type InstanceCreator func(env *common.Env, value any) (*ClassInstance, error)

// InstanceCreatorRegistry holds all registered instance creators
type InstanceCreatorRegistry struct {
	creators map[string]InstanceCreator
}

var globalInstanceCreatorRegistry = &InstanceCreatorRegistry{
	creators: make(map[string]InstanceCreator),
}

// RegisterInstanceCreator registers a creator for a specific type
func RegisterInstanceCreator(typeName string, creator InstanceCreator) {
	globalInstanceCreatorRegistry.creators[typeName] = creator
}

// CreateInstanceFor attempts to create a ClassInstance for the specified type
func CreateInstanceFor(env *common.Env, typeName string, value any) (*ClassInstance, error) {
	if creator, ok := globalInstanceCreatorRegistry.creators[typeName]; ok {
		return creator(env, value)
	}
	return nil, fmt.Errorf("no instance creator registered for type: %s", typeName)
}

// GetInstanceCreator returns the creator for a specific type
func GetInstanceCreator(typeName string) (InstanceCreator, bool) {
	creator, ok := globalInstanceCreatorRegistry.creators[typeName]
	return creator, ok
}

// InitializeBuiltinInstanceCreators registers all builtin instance creators
func InitializeBuiltinInstanceCreators() {
	// Register String instance creator
	RegisterInstanceCreator("String", func(env *common.Env, value any) (*ClassInstance, error) {
		strVal, ok := ConvertTo(env, "String", value)
		if !ok {
			return nil, fmt.Errorf("cannot convert value to String")
		}
		return CreateStringInstance(env, strVal.(string))
	})

	// Register Int instance creator
	RegisterInstanceCreator("Int", func(env *common.Env, value any) (*ClassInstance, error) {
		intVal, ok := ConvertTo(env, "Int", value)
		if !ok {
			return nil, fmt.Errorf("cannot convert value to Int")
		}
		return CreateIntInstance(env, intVal.(int))
	})

	// Register Float instance creator
	RegisterInstanceCreator("Float", func(env *common.Env, value any) (*ClassInstance, error) {
		floatVal, ok := ConvertTo(env, "Float", value)
		if !ok {
			return nil, fmt.Errorf("cannot convert value to Float")
		}
		return CreateFloatInstance(env, floatVal.(float64))
	})

	// Register Bool instance creator
	RegisterInstanceCreator("Bool", func(env *common.Env, value any) (*ClassInstance, error) {
		boolVal, ok := ConvertTo(env, "Bool", value)
		if !ok {
			return nil, fmt.Errorf("cannot convert value to Bool")
		}
		return CreateBoolInstance(env, boolVal.(bool))
	})

	// Register Bytes instance creator
	RegisterInstanceCreator("Bytes", func(env *common.Env, value any) (*ClassInstance, error) {
		return CreateBytesInstance(env, value)
	})

	// Register Array instance creator
	RegisterInstanceCreator("Array", func(env *common.Env, value any) (*ClassInstance, error) {
		arrVal, ok := ConvertTo(env, "Array", value)
		if !ok {
			return nil, fmt.Errorf("cannot convert value to Array")
		}
		if slice, ok := arrVal.([]any); ok {
			return CreateArrayInstance(env, slice)
		}
		return nil, fmt.Errorf("array conversion did not return []any")
	})

	// Register Map instance creator
	RegisterInstanceCreator("Map", func(env *common.Env, value any) (*ClassInstance, error) {
		return CreateMapInstance(env, value)
	})
}


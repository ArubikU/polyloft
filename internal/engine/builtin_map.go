package engine

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

type mapEntry = ast.MapEntry

// hashValue computes a hash for a value
func hashValue(v any) uint64 {
	h := fnv.New64a()
	switch val := v.(type) {
	case string:
		h.Write([]byte(val))
	case int:
		h.Write([]byte(fmt.Sprintf("%d", val)))
	case float64:
		h.Write([]byte(fmt.Sprintf("%f", val)))
	case bool:
		if val {
			h.Write([]byte("true"))
		} else {
			h.Write([]byte("false"))
		}
	default:
		h.Write([]byte(fmt.Sprintf("%v", v)))
	}
	return h.Sum64()
}

// equals checks if two values are equal
func equals(a, b any) bool {
	// Simple equality check - can be enhanced
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func InstallSerializableInterface(env *Env) error {
	serializableInterface := NewInterfaceBuilder("Serializable").
		AddMethod("serialize", "String | Map", []ast.Parameter{}).
		AddMethod("deserialize", "T", []ast.Parameter{
			{Name: "data", Type: nil},
		})
	_, err := serializableInterface.Build(env)
	return err
}

func InstallMapBuiltin(env *Env) error {
	// Get type references from already-installed builtin types
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	iterableInterface := common.BuiltinInterfaceIterable.GetInterfaceDefinition(env)

	mapClass := NewClassBuilder("Map").
		AddTypeParameter("K", []string{}, false).
		AddTypeParameter("V", []string{}, false).
		AddInterface(iterableInterface).
		AddField("_data", &ast.Type{Name: "map", IsBuiltin: true}, []string{"private"})

	// Instance methods

	// get(key: any) -> Any
	mapClass.AddBuiltinMethod("get", ast.ANY, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		if entry, exists := data[hash]; exists {
			return entry.Value, nil
		}
		return nil, nil
	}, []string{})

	// set(key: any, value: Any) -> Void
	mapClass.AddBuiltinMethod("set", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
		{Name: "value", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		data[hash] = &mapEntry{Key: args[0], Value: args[1]}
		return nil, nil
	}, []string{})

	// put(key: any, value: Any) -> Void (alias for set)
	mapClass.AddBuiltinMethod("put", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
		{Name: "value", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		data[hash] = &mapEntry{Key: args[0], Value: args[1]}
		return nil, nil
	}, []string{})

	// has(key: any) -> Bool
	mapClass.AddBuiltinMethod("has", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		_, exists := data[hash]
		return exists, nil
	}, []string{})

	// hasKey(key: any) -> Bool (alias for has)
	mapClass.AddBuiltinMethod("hasKey", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		_, exists := data[hash]
		return exists, nil
	}, []string{})

	// remove(key: any) -> Void
	mapClass.AddBuiltinMethod("remove", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		delete(data, hash)
		return nil, nil
	}, []string{})

	// delete(key: any) -> Void (alias for remove)
	mapClass.AddBuiltinMethod("delete", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		hash := hashValue(args[0])
		delete(data, hash)
		return nil, nil
	}, []string{})

	// clear() -> Void
	mapClass.AddBuiltinMethod("clear", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_data"] = make(map[uint64]*mapEntry)
		return nil, nil
	}, []string{})

	// keys() -> Array
	mapClass.AddBuiltinMethod("keys", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		keys := make([]any, 0, len(data))
		for _, entry := range data {
			keys = append(keys, entry.Key)
		}
		return keys, nil
	}, []string{})

	// values() -> Array
	mapClass.AddBuiltinMethod("values", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		values := make([]any, 0, len(data))
		for _, entry := range data {
			values = append(values, entry.Value)
		}
		return values, nil
	}, []string{})

	// entries() -> Array
	mapClass.AddBuiltinMethod("entries", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		entries := make([]any, 0, len(data))
		for _, entry := range data {
			entryArr := []any{entry.Key, entry.Value}
			entries = append(entries, entryArr)
		}
		return entries, nil
	}, []string{})

	// size() -> Int
	mapClass.AddBuiltinMethod("size", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		return len(data), nil
	}, []string{})

	// length() -> Int (alias for size)
	mapClass.AddBuiltinMethod("length", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		return len(data), nil
	}, []string{})

	// isEmpty() -> Bool
	mapClass.AddBuiltinMethod("isEmpty", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		return len(data) == 0, nil
	}, []string{})

	// utils.ToString() -> String
	mapClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		if len(data) == 0 {
			return "{}", nil
		}

		result := "{"
		first := true
		for _, entry := range data {
			if !first {
				result += ", "
			}
			first = false
			result += utils.ToString(entry.Key) + ": " + utils.ToString(entry.Value)
		}
		result += "}"

		return result, nil
	}, []string{})

	// serialize() -> String (convert to JSON string)
	mapClass.AddBuiltinMethod("serialize", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)

		// Convert to Go map for JSON encoding
		objMap, err := MapToObject(instance)
		if err != nil {
			return nil, err
		}

		// Encode as JSON
		jsonBytes, err := json.Marshal(objMap)
		if err != nil {
			return nil, err
		}

		return string(jsonBytes), nil
	}, []string{})

	// deserialize(data: String) -> Map (parse JSON string)
	mapClass.AddStaticMethod("deserialize", &ast.Type{Name: "Map", IsBuiltin: true}, []ast.Parameter{
		{Name: "data", Type: stringType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		jsonStr := utils.ToString(args[0])

		var data map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			return nil, err
		}

		return CreateMapInstance(env, data)
	}))

	mapClass.AddStaticMethod("fromString", &ast.Type{Name: "Map", IsBuiltin: true}, []ast.Parameter{
		{Name: "str", Type: stringType},
	}, common.Func(func(env *common.Env, args []any) (any, error) {
		str := utils.ToString(args[0])
		data := make(map[string]any)

		// Simple parsing logic (for demonstration purposes)
		entries := strings.Split(str, ",")
		for _, entry := range entries {
			kv := strings.SplitN(entry, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				data[key] = value
			}
		}

		return CreateMapInstance(env, data)
	}))

	// getEntries() -> List<MapEntry<K,V>>
	mapClass.AddBuiltinMethod("getEntries", &ast.Type{Name: "List", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64]*mapEntry)

		// Get MapEntry class definition
		mapEntryClassDef, ok := builtinClasses["MapEntry"]
		if !ok {
			return nil, ThrowRuntimeError((*Env)(callEnv), "MapEntry class not found")
		}

		// Create a slice of MapEntry instances
		entries := make([]any, 0, len(data))
		for _, entry := range data {
			// Create MapEntry instance
			mapEntryInstance := &ClassInstance{
				ClassName: "MapEntry",
				Fields: map[string]any{
					"key":   entry.Key,
					"value": entry.Value,
				},
				Methods:     make(map[string]common.Func),
				ParentClass: mapEntryClassDef,
			}
			entries = append(entries, mapEntryInstance)
		}

		// Create List instance containing the entries
		listClass, exists := callEnv.Get("List")
		if !exists {
			return nil, ThrowRuntimeError((*Env)(callEnv), "List class not found")
		}

		if listCtor, ok := listClass.(*common.ClassConstructor); ok {
			// Create Array instance for the entries
			arrayInstance, err := CreateArrayInstance((*Env)(callEnv), entries)
			if err != nil {
				return nil, err
			}

			// Call List constructor with the array
			listInstance, err := listCtor.Func(callEnv, []any{arrayInstance})
			if err != nil {
				return nil, err
			}
			return listInstance, nil
		}

		return nil, ThrowRuntimeError((*Env)(callEnv), "List is not a constructor")
	}, []string{})

	// Build the class
	_, err := mapClass.Build(env)
	return err
}

// CreateMapInstance creates a Map instance from a map[string]any
// This is used when evaluating {} literals
func CreateMapInstance(env *Env, data map[string]any) (*ClassInstance, error) {
	mapClassVal, ok := env.Get("__MapClass__")
	if !ok {
		return nil, ThrowInitializationError(env, "Map class")
	}

	mapClass := mapClassVal.(*ClassDefinition)

	// Create instance
	instance, err := createClassInstance(mapClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)

	// Convert to hash-based storage
	hashData := make(map[uint64]*mapEntry)
	for k, v := range data {
		hash := hashValue(k)
		hashData[hash] = &mapEntry{Key: k, Value: v}
	}
	classInstance.Fields["_data"] = hashData

	return classInstance, nil
}

// MapToObject converts a Map instance back to a Go map[string]any for JSON serialization
func MapToObject(mapInstance *ClassInstance) (map[string]any, error) {
	dataField, ok := mapInstance.Fields["_data"]
	if !ok {
		return nil, ThrowAttributeError(nil, "_data", "Map")
	}

	hashData, ok := dataField.(map[uint64]*mapEntry)
	if !ok {
		return nil, ThrowTypeError(nil, "map", dataField)
	}

	// Convert hash-based storage to regular map for JSON
	result := make(map[string]any)
	for _, entry := range hashData {
		// Convert key to string for JSON
		keyStr := utils.ToString(entry.Key)

		// Recursively convert nested ClassInstance objects
		if nestedInstance, ok := entry.Value.(*ClassInstance); ok {
			if nestedInstance.ClassName == "Map" {
				nestedObj, err := MapToObject(nestedInstance)
				if err != nil {
					return nil, err
				}
				result[keyStr] = nestedObj
			} else if nestedInstance.ClassName == "Array" {
				nestedSlice, err := ArrayToSlice(nestedInstance)
				if err != nil {
					return nil, err
				}
				result[keyStr] = nestedSlice
			} else {
				result[keyStr] = entry.Value
			}
		} else {
			result[keyStr] = entry.Value
		}
	}

	return result, nil
}

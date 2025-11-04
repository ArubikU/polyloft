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
func hashValue(env *Env, v any) uint64 {
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
	case *ClassInstance:
		// Extract primitive value from ClassInstance
		if val.ClassName == "String" {
			if strVal, ok := val.Fields["_value"].(string); ok {
				h.Write([]byte(strVal))
				return h.Sum64()
			}
		} else if val.ClassName == "Int" {
			if intVal, ok := val.Fields["_value"].(int); ok {
				h.Write([]byte(fmt.Sprintf("%d", intVal)))
				return h.Sum64()
			}
		} else if val.ClassName == "Float" {
			if floatVal, ok := val.Fields["_value"].(float64); ok {
				h.Write([]byte(fmt.Sprintf("%f", floatVal)))
				return h.Sum64()
			}
		} else if val.ClassName == "Bool" {
			if boolVal, ok := val.Fields["_value"].(bool); ok {
				if boolVal {
					h.Write([]byte("true"))
				} else {
					h.Write([]byte("false"))
				}
				return h.Sum64()
			}
		}
		methods := val.ParentClass.GetMethods("hash")
		method := common.SelectMethodOverload(methods, 0)
		if method == nil {
			h.Write([]byte(fmt.Sprintf("%p", v)))
			return h.Sum64()
		}
		hashResult, err := CallInstanceMethod(val, *method, env, []any{})
		if err == nil {
			intval, _ := utils.AsInt(hashResult)
			return uint64(intval)
		}
		return hashValue(env, hashResult)
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
	indexableInterface := common.BuiltinIndexableInterface.GetInterfaceDefinition(env)

	mapClass := NewClassBuilder("Map").
		AddTypeParameters(common.KBound.AsGenericType().AsArray()).
		AddTypeParameters(common.VBound.AsGenericType().AsArray()).
		AddInterface(iterableInterface).
		AddInterface(indexableInterface).
		AddField("_data", &ast.Type{Name: "map", IsBuiltin: true}, []string{"private"})

	// Instance methods

	// get(key: K) -> V
	mapClass.AddBuiltinMethod("get", &common.VBound.Name, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for _, entry := range entries {
				if equals(entry.Key, args[0]) {
					return entry.Value, nil
				}
			}
		}
		return nil, nil
	}, []string{})

	// set(key: K, value: V) -> Void
	mapClass.AddBuiltinMethod("set", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
		{Name: "value", Type: &common.VBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for i, entry := range entries {
				if equals(entry.Key, args[0]) {
					entries[i].Value = args[1]
					return nil, nil
				}
			}
			// Key not found in this bucket, add it
			data[hash] = append(entries, &mapEntry{Key: args[0], Value: args[1]})
		} else {
			data[hash] = []*mapEntry{{Key: args[0], Value: args[1]}}
		}
		return nil, nil
	}, []string{})

	mapClass.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_data"] = make(map[uint64][]*mapEntry)
		instance.Fields["_entries"] = make([]*mapEntry, 0)
		return nil, nil
	})

	// put(key: any, value: Any) -> Void (alias for set)
	mapClass.AddBuiltinMethod("put", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
		{Name: "value", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for i, entry := range entries {
				if equals(entry.Key, args[0]) {
					entries[i].Value = args[1]
					return nil, nil
				}
			}
			// Key not found in this bucket, add it
			data[hash] = append(entries, &mapEntry{Key: args[0], Value: args[1]})
		} else {
			data[hash] = []*mapEntry{{Key: args[0], Value: args[1]}}
		}
		return nil, nil
	}, []string{})

	// has(key: K) -> Bool
	mapClass.AddBuiltinMethod("has", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for _, entry := range entries {
				if equals(entry.Key, args[0]) {
					return true, nil
				}
			}
		}
		return false, nil
	}, []string{})

	// hasKey(key: K) -> Bool (alias for has)
	mapClass.AddBuiltinMethod("hasKey", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for _, entry := range entries {
				if equals(entry.Key, args[0]) {
					return true, nil
				}
			}
		}
		return false, nil
	}, []string{})

	// __get(key: K) -> V (Indexable interface)
	// Also supports numeric index for iteration (returns Pair<K, V>)
	mapClass.AddBuiltinMethod("__get", &common.VBound.Name, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		// Check if the argument is a numeric index (for iteration)
		if idx, ok := utils.AsInt(args[0]); ok {
			// Use stable _entries slice for iteration
			entries, hasEntries := instance.Fields["_entries"].([]*mapEntry)
			if hasEntries && idx >= 0 && idx < len(entries) {
				entry := entries[idx]
				// Create a Pair instance
				pairClass, exists := lookupClass("Pair", "")
				if !exists {
					// Fallback to array if Pair not available
					return []any{entry.Key, entry.Value}, nil
				}

				// Construct Pair
				pairInstance, err := constructPairInstance(pairClass, entry.Key, entry.Value, (*Env)(callEnv))
				if err != nil {
					// Fallback to array on error
					return []any{entry.Key, entry.Value}, nil
				}
				return pairInstance, nil
			}
			return nil, nil // Index out of bounds
		}

		// Otherwise, treat as key lookup - return value directly
		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for _, entry := range entries {
				if equals(entry.Key, args[0]) {
					return entry.Value, nil
				}
			}
		}
		return nil, nil
	}, []string{})

	// __set(key: K, value: V) -> Void (Indexable interface)
	mapClass.AddBuiltinMethod("__set", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
		{Name: "value", Type: &common.VBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)
		entries, hasEntries := instance.Fields["_entries"].([]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if bucketEntries, exists := data[hash]; exists {
			// Check if key already exists
			for i, entry := range bucketEntries {
				if equals(entry.Key, args[0]) {
					// Update existing entry value
					bucketEntries[i].Value = args[1]
					// Also update in _entries if it exists
					if hasEntries {
						for j, e := range entries {
							if equals(e.Key, args[0]) {
								entries[j].Value = args[1]
								break
							}
						}
					}
					return nil, nil
				}
			}
			// Key not found in this bucket, add it
			newEntry := &mapEntry{Key: args[0], Value: args[1]}
			data[hash] = append(bucketEntries, newEntry)
			if hasEntries {
				instance.Fields["_entries"] = append(entries, newEntry)
			}
		} else {
			// New bucket
			newEntry := &mapEntry{Key: args[0], Value: args[1]}
			data[hash] = []*mapEntry{newEntry}
			if hasEntries {
				instance.Fields["_entries"] = append(entries, newEntry)
			}
		}
		return nil, nil
	}, []string{})

	// __contains(key: K) -> Bool (Indexable interface)
	mapClass.AddBuiltinMethod("__contains", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: &common.KBound.Name},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for _, entry := range entries {
				if equals(entry.Key, args[0]) {
					return true, nil
				}
			}
		}
		return false, nil
	}, []string{})

	// __length() -> Int (Iterable interface)
	// Returns the number of key-value pairs in the map
	mapClass.AddBuiltinMethod("__length", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		// Use _entries if available for accurate count
		if entries, ok := instance.Fields["_entries"].([]*mapEntry); ok {
			return len(entries), nil
		}

		// Fallback to counting from hash map
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)
		size := 0
		for _, entries := range data {
			size += len(entries)
		}
		return size, nil
	}, []string{})

	// remove(key: any) -> Void
	mapClass.AddBuiltinMethod("remove", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for i, entry := range entries {
				if equals(entry.Key, args[0]) {
					// Remove the entry from the slice
					data[hash] = append(entries[:i], entries[i+1:]...)
					// If the slice is empty, remove the hash entry
					if len(data[hash]) == 0 {
						delete(data, hash)
					}
					break
				}
			}
		}
		return nil, nil
	}, []string{})

	// delete(key: any) -> Void (alias for remove)
	mapClass.AddBuiltinMethod("delete", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "key", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		hash := hashValue(callEnv, args[0])
		if entries, exists := data[hash]; exists {
			for i, entry := range entries {
				if equals(entry.Key, args[0]) {
					// Remove the entry from the slice
					data[hash] = append(entries[:i], entries[i+1:]...)
					// If the slice is empty, remove the hash entry
					if len(data[hash]) == 0 {
						delete(data, hash)
					}
					break
				}
			}
		}
		return nil, nil
	}, []string{})

	// clear() -> Void
	mapClass.AddBuiltinMethod("clear", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_data"] = make(map[uint64][]*mapEntry)
		return nil, nil
	}, []string{})

	// keys() -> Array
	mapClass.AddBuiltinMethod("keys", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		keys := make([]any, 0, len(data))
		for _, entries := range data {
			for _, entry := range entries {
				keys = append(keys, entry.Key)
			}
		}
		return keys, nil
	}, []string{})

	// values() -> Array
	mapClass.AddBuiltinMethod("values", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		values := make([]any, 0, len(data))
		for _, entries := range data {
			for _, entry := range entries {
				values = append(values, entry.Value)
			}
		}
		return values, nil
	}, []string{})

	// entries() -> Array
	mapClass.AddBuiltinMethod("entries", &ast.Type{Name: "array", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		entries := make([]any, 0, len(data))
		for _, entrySlice := range data {
			for _, entry := range entrySlice {
				entryArr := []any{entry.Key, entry.Value}
				entries = append(entries, entryArr)
			}
		}
		return entries, nil
	}, []string{})

	// size() -> Int
	mapClass.AddBuiltinMethod("size", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		size := 0
		for _, entries := range data {
			size += len(entries)
		}
		return size, nil
	}, []string{})

	// length() -> Int (alias for size)
	mapClass.AddBuiltinMethod("length", &ast.Type{Name: "int", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		size := 0
		for _, entries := range data {
			size += len(entries)
		}
		return size, nil
	}, []string{})

	// isEmpty() -> Bool
	mapClass.AddBuiltinMethod("isEmpty", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		for _, entries := range data {
			if len(entries) > 0 {
				return false, nil
			}
		}
		return true, nil
	}, []string{})

	// utils.ToString() -> String
	mapClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		size := 0
		for _, entries := range data {
			size += len(entries)
		}
		if size == 0 {
			return "{}", nil
		}

		result := "{"
		first := true
		for _, entries := range data {
			for _, entry := range entries {
				if !first {
					result += ", "
				}
				first = false
				result += utils.ToString(entry.Key) + ": " + utils.ToString(entry.Value)
			}
		}
		result += "}"

		return result, nil
	}, []string{})

	// serialize() -> String (convert to JSON string)
	mapClass.AddBuiltinMethod("serialize", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		// Convert to Go map for JSON encoding
		objMap, err := MapToObject((*Env)(callEnv), instance)
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
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		data := instance.Fields["_data"].(map[uint64][]*mapEntry)

		// Get MapEntry class definition
		mapEntryClassDef, ok := builtinClasses["MapEntry"]
		if !ok {
			return nil, ThrowRuntimeError((*Env)(callEnv), "MapEntry class not found")
		}

		size := 0
		for _, entries := range data {
			size += len(entries)
		}

		// Create a slice of MapEntry instances
		entries := make([]any, 0, size)
		for _, entrySlice := range data {
			for _, entry := range entrySlice {
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
		}

		// Create List instance containing the entries
		listCtor := common.BuiltinTypeList.GetConstructor(callEnv)
		if listCtor == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "List class not found")
		}

		if listCtor != nil {
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

func CreateMapInstance(env *Env, data any) (*ClassInstance, error) {
	mapClass := common.BuiltinTypeMap.GetClassDefinition(env)
	if mapClass == nil {
		return nil, ThrowInitializationError(env, "Map class")
	}

	// Create instance
	instance, err := createClassInstance(mapClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)

	// Convert to hash-based storage for fast lookups
	hashData := make(map[uint64][]*mapEntry)
	// Also maintain insertion order for stable iteration
	entries := make([]*mapEntry, 0)

	// Convert data to map entries with automatic type conversion
	convertedData, ok := MapToData(env, data)
	if !ok {
		return nil, ThrowTypeError(env, "map-compatible type", data)
	}

	for k, v := range convertedData {
		// Convert key and value to appropriate types
		convertedKey := ConvertMapKey(env, k)
		convertedValue := ConvertMapValue(env, v)

		entry := &mapEntry{Key: convertedKey, Value: convertedValue}
		hash := hashValue(env, convertedKey)
		hashData[hash] = append(hashData[hash], entry)
		entries = append(entries, entry)
	}

	classInstance.Fields["_data"] = hashData
	classInstance.Fields["_entries"] = entries

	return classInstance, nil
}

// MapToData converts various types to a map[string]any representation
func MapToData(env *Env, value any) (map[string]any, bool) {
	// Use the unified type converter
	result, ok := ConvertTo(env, "Map", value)
	if !ok {
		return nil, false
	}
	
	// The converter might return the internal hash structure or a map[string]any
	switch v := result.(type) {
	case map[string]any:
		return v, true
	case map[any]any:
		// Convert map[any]any to map[string]any
		converted := make(map[string]any)
		for key, val := range v {
			converted[utils.ToString(key)] = val
		}
		return converted, true
	default:
		return nil, false
	}
}

// ConvertMapKey converts a key to the appropriate type for use in a Map
func ConvertMapKey(env *Env, key any) any {
	if key == nil {
		// Create a String instance for "null"
		if strInst, err := CreateInstanceFor(env, "String", "null"); err == nil {
			return strInst
		}
		return "null"
	}

	// If already a ClassInstance, return as is
	if inst, ok := key.(*ClassInstance); ok {
		return inst
	}

	// Use ConvertToClassInstance for all other cases
	return ConvertToClassInstance(env, key)
}

// ConvertMapValue converts a value to the appropriate type for storage in a Map
func ConvertMapValue(env *Env, value any) any {
	if value == nil {
		return nil
	}

	// Use ConvertToClassInstance which handles all types uniformly
	return ConvertToClassInstance(env, value)
}

// MapToObject converts a Map instance back to a Go map[string]any for JSON serialization
func MapToObject(env *Env, mapInstance *ClassInstance) (map[string]any, error) {
	dataField, ok := mapInstance.Fields["_data"]
	if !ok {
		return nil, ThrowAttributeError(env, "_data", "Map")
	}

	hashData, ok := dataField.(map[uint64][]*mapEntry)
	if !ok {
		return nil, ThrowTypeError(env, "map", dataField)
	}

	// Convert hash-based storage to regular map for JSON
	result := make(map[string]any)
	for _, entries := range hashData {
		for _, entry := range entries {
			// Convert key to string for JSON
			var keyStr string
			if keyInst, ok := entry.Key.(*ClassInstance); ok {
				// Extract primitive value from key wrapper classes
				stringDef := common.BuiltinTypeString.GetClassDefinition((*common.Env)(env))
				intDef := common.BuiltinTypeInt.GetClassDefinition((*common.Env)(env))
				floatDef := common.BuiltinTypeFloat.GetClassDefinition((*common.Env)(env))
				boolDef := common.BuiltinTypeBool.GetClassDefinition((*common.Env)(env))

				if keyInst.ParentClass.IsSubclassOf(stringDef) {
					if val, ok := keyInst.Fields["_value"].(string); ok {
						keyStr = val
					} else {
						keyStr = utils.ToString(keyInst)
					}
				} else if keyInst.ParentClass.IsSubclassOf(intDef) {
					if val, ok := keyInst.Fields["_value"].(int); ok {
						keyStr = utils.ToString(val)
					} else {
						keyStr = utils.ToString(keyInst)
					}
				} else if keyInst.ParentClass.IsSubclassOf(floatDef) {
					if val, ok := keyInst.Fields["_value"].(float64); ok {
						keyStr = utils.ToString(val)
					} else {
						keyStr = utils.ToString(keyInst)
					}
				} else if keyInst.ParentClass.IsSubclassOf(boolDef) {
					if val, ok := keyInst.Fields["_value"].(bool); ok {
						keyStr = utils.ToString(val)
					} else {
						keyStr = utils.ToString(keyInst)
					}
				} else {
					keyStr = utils.ToString(keyInst)
				}
			} else {
				keyStr = utils.ToString(entry.Key)
			}

			// Recursively convert nested ClassInstance objects and extract primitives
			if nestedInstance, ok := entry.Value.(*ClassInstance); ok {
				mapDef := common.BuiltinTypeMap.GetClassDefinition((*common.Env)(env))
				arrayDef := common.BuiltinTypeArray.GetClassDefinition((*common.Env)(env))
				listDef := common.BuiltinTypeList.GetClassDefinition((*common.Env)(env))
				stringDef := common.BuiltinTypeString.GetClassDefinition((*common.Env)(env))
				intDef := common.BuiltinTypeInt.GetClassDefinition((*common.Env)(env))
				floatDef := common.BuiltinTypeFloat.GetClassDefinition((*common.Env)(env))
				boolDef := common.BuiltinTypeBool.GetClassDefinition((*common.Env)(env))

				if nestedInstance.ParentClass.IsSubclassOf(mapDef) {
					nestedObj, err := MapToObject(env, nestedInstance)
					if err != nil {
						return nil, err
					}
					result[keyStr] = nestedObj
				} else if nestedInstance.ParentClass.IsSubclassOf(arrayDef) {
					nestedSlice, err := ArrayToSlice(env, nestedInstance)
					if err != nil {
						return nil, err
					}
					result[keyStr] = nestedSlice
				} else if nestedInstance.ParentClass.IsSubclassOf(listDef) {
					// List uses _items which is a pointer to []any
					if itemsPtr, ok := nestedInstance.Fields["_items"].(*[]any); ok {
						result[keyStr] = *itemsPtr
					} else {
						result[keyStr] = nestedInstance
					}
				} else if nestedInstance.ParentClass.IsSubclassOf(stringDef) {
					// Extract primitive string value
					if val, ok := nestedInstance.Fields["_value"].(string); ok {
						result[keyStr] = val
					} else {
						result[keyStr] = utils.ToString(nestedInstance)
					}
				} else if nestedInstance.ParentClass.IsSubclassOf(intDef) {
					// Extract primitive int value
					if val, ok := nestedInstance.Fields["_value"].(int); ok {
						result[keyStr] = val
					} else {
						result[keyStr] = nestedInstance
					}
				} else if nestedInstance.ParentClass.IsSubclassOf(floatDef) {
					// Extract primitive float value
					if val, ok := nestedInstance.Fields["_value"].(float64); ok {
						result[keyStr] = val
					} else {
						result[keyStr] = nestedInstance
					}
				} else if nestedInstance.ParentClass.IsSubclassOf(boolDef) {
					// Extract primitive bool value
					if val, ok := nestedInstance.Fields["_value"].(bool); ok {
						result[keyStr] = val
					} else {
						result[keyStr] = nestedInstance
					}
				} else {
					result[keyStr] = entry.Value
				}
			} else {
				result[keyStr] = entry.Value
			}
		}
	}

	return result, nil
}

func MapToClassMap(mapInstance *ClassInstance) (map[*ClassInstance]*ClassInstance, error) {
	dataField, ok := mapInstance.Fields["_data"]
	if !ok {
		return nil, ThrowAttributeError(nil, "_data", "Map")
	}

	hashData, ok := dataField.(map[uint64][]*mapEntry)
	if !ok {
		return nil, ThrowTypeError(nil, "map", dataField)
	}

	classMap := make(map[*ClassInstance]*ClassInstance)
	for _, entries := range hashData {
		for _, entry := range entries {
			if entry != nil && entry.Key != nil && entry.Value != nil {
				keyInstance, ok := entry.Key.(*ClassInstance)
				if !ok {
					return nil, ThrowTypeError(nil, "*ClassInstance", entry.Key)
				}
				valueInstance, ok := entry.Value.(*ClassInstance)
				if !ok {
					return nil, ThrowTypeError(nil, "*ClassInstance", entry.Value)
				}
				classMap[keyInstance] = valueInstance
			}
		}
	}

	return classMap, nil
}

// constructPairInstance creates a Pair instance with the given key and value
func constructPairInstance(pairClass *ClassDefinition, key, value any, env *Env) (*ClassInstance, error) {
	// Create instance using the proper initialization function
	instance, err := createClassInstance(pairClass, env, []any{})
	if err != nil {
		return nil, err
	}

	classInstance := instance.(*ClassInstance)

	// Initialize fields (Pair uses "key" and "value")
	classInstance.Fields["key"] = key
	classInstance.Fields["value"] = value

	return classInstance, nil
}

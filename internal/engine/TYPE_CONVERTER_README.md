# Unified Type Converter System

This document describes the unified type converter system implemented in `type_converter.go`.

## Overview

Previously, type conversion logic was scattered across multiple files:
- `AsBytes` in `builtin_bytes.go` (~210 lines)
- `MapToData` in `builtin_map.go` (~30 lines)
- `ConvertToClassInstance` in `builtin_array.go` (~90 lines)
- `ConvertMapKey` in `builtin_map.go` (~70 lines)
- `ConvertMapValue` in `builtin_map.go` (~90 lines)

**Total duplicated code: ~490 lines**

The unified system centralizes this logic to make it easier to:
1. Add new type conversions
2. Register custom type converters
3. Maintain consistency across the codebase

**After unification: ~350 lines total (including registry infrastructure)**
**Net reduction: ~140 lines + significantly better organization**

## Components

### Type Converters

Type converters transform values from one type to another (e.g., any value to bytes, any value to int).

**Registration:**
```go
RegisterTypeConverter("TypeName", func(env *common.Env, value any) (any, bool) {
    // Conversion logic here
    return convertedValue, success
})
```

**Usage:**
```go
result, ok := ConvertTo(env, "TypeName", value)
if !ok {
    // Handle conversion failure
}
```

### Instance Creators

Instance creators create ClassInstance objects from native Go types.

**Registration:**
```go
RegisterInstanceCreator("TypeName", func(env *common.Env, value any) (*ClassInstance, error) {
    // Creation logic here
    return instance, nil
})
```

**Usage:**
```go
instance, err := CreateInstanceFor(env, "TypeName", value)
if err != nil {
    // Handle creation failure
}
```

## Built-in Converters

The following type converters are registered by default:

- **Bytes**: Converts any value to `[]byte`
  - Supports hex strings (0x prefix)
  - Supports binary strings (0b prefix)
  - Converts primitive types, arrays, maps, and ClassInstances
- **Array**: Converts any value to `[]any`
- **String**: Converts any value to `string`
- **Int**: Converts any value to `int`
- **Float**: Converts any value to `float64`
- **Bool**: Converts any value to `bool`
- **Map**: Converts any value to `map[string]any`

## Built-in Instance Creators

The following instance creators are registered by default:

- **String**: Creates a String ClassInstance
- **Int**: Creates an Int ClassInstance
- **Float**: Creates a Float ClassInstance
- **Bool**: Creates a Bool ClassInstance
- **Bytes**: Creates a Bytes ClassInstance
- **Array**: Creates an Array ClassInstance
- **Map**: Creates a Map ClassInstance

## Adding New Types

To add support for a new type:

1. Register a type converter:
```go
RegisterTypeConverter("MyType", func(env *common.Env, value any) (any, bool) {
    // Convert value to MyType's native representation
    return nativeValue, true
})
```

2. Register an instance creator:
```go
RegisterInstanceCreator("MyType", func(env *common.Env, value any) (*ClassInstance, error) {
    // Create MyType ClassInstance
    return CreateMyTypeInstance(env, value)
})
```

3. The converters should be registered in `InitializeBuiltinTypeConverters()` or `InitializeBuiltinInstanceCreators()`.

## Migration Guide

### Before (scattered approach):

```go
// In builtin_bytes.go
func AsBytes(env *common.Env, value any) ([]byte, bool) {
    // 200+ lines of conversion logic with nested type switches
    switch v := value.(type) {
    case string:
        // Handle string...
    case int, int8, int16, int32, int64:
        // Handle ints...
    case *ClassInstance:
        // 150+ lines of ClassInstance handling...
    }
}

// In builtin_map.go
func MapToData(env *Env, value any) (map[string]any, bool) {
    // Similar custom conversion logic
}

// In builtin_array.go
func ConvertToClassInstance(env *Env, value any) any {
    // Another 90 lines of custom conversion
}

// In builtin_map.go
func ConvertMapKey(env *Env, key any) any {
    // 70 lines of key-specific conversion
}

func ConvertMapValue(env *Env, value any) any {
    // 90 lines of value-specific conversion
}
```

### After (unified approach):

```go
// Type converter is registered once in type_converter.go
RegisterTypeConverter("Bytes", func(env *common.Env, value any) (any, bool) {
    // Conversion logic (handles all cases)
})

// Used everywhere consistently
result, ok := ConvertTo(env, "Bytes", value)

// Functions simplified to use the registry
func AsBytes(env *common.Env, value any) ([]byte, bool) {
    result, ok := ConvertTo(env, "Bytes", value)
    if !ok {
        return []byte{}, false
    }
    return result.([]byte), true
}

func ConvertToClassInstance(env *Env, value any) any {
    // Simplified to ~50 lines using CreateInstanceFor
}

func ConvertMapKey(env *Env, key any) any {
    // Simplified to ~10 lines using ConvertToClassInstance
}

func ConvertMapValue(env *Env, value any) any {
    // Simplified to ~5 lines using ConvertToClassInstance
}
```

## Benefits

1. **Centralized**: All type conversion logic is in one place
2. **Extensible**: Easy to add new types without modifying existing code
3. **Consistent**: Same API for all type conversions
4. **Maintainable**: Changes to conversion logic only need to be made once
5. **Testable**: Each converter can be tested independently
6. **Less Duplication**: Reduced ~140 lines of duplicate code
7. **Better Organization**: Clear separation between conversion logic and usage

## Examples

### Converting to Bytes
```go
// Convert a string to bytes
bytes, ok := ConvertTo(env, "Bytes", "hello")

// Convert an array to bytes
bytes, ok := ConvertTo(env, "Bytes", arrayInstance)

// Convert hex string to bytes
bytes, ok := ConvertTo(env, "Bytes", "0xFF00")

// Convert binary string to bytes
bytes, ok := ConvertTo(env, "Bytes", "0b11110000")
```

### Creating Class Instances
```go
// Create an Int instance from a native int
intInstance, err := CreateInstanceFor(env, "Int", 42)

// Create a String instance from a native string
strInstance, err := CreateInstanceFor(env, "String", "hello")

// Create an Array instance from a native slice
arrayInstance, err := CreateInstanceFor(env, "Array", []any{1, 2, 3})
```

### Custom Type Conversion
```go
// Register a custom converter for a hypothetical "Color" type
RegisterTypeConverter("Color", func(env *common.Env, value any) (any, bool) {
    switch v := value.(type) {
    case string:
        return parseColorFromString(v)
    case int:
        return colorFromInt(v)
    default:
        return nil, false
    }
})

// Use it
color, ok := ConvertTo(env, "Color", "#FF0000")
```

## Implementation Details

### Type Aliases
The code uses `*Env` and `*common.Env` interchangeably because `Env` is defined as:
```go
type Env = common.Env
```
This allows the registry functions to accept either type, maintaining backward compatibility.

### Initialization Order
The converters are initialized after all builtin types are installed in `engine.go`:
```go
// Install all builtin types first...
InstallStringBuiltin(env)
InstallIntBuiltin(env)
// ... etc

// Then initialize converters
InitializeBuiltinTypeConverters()
InitializeBuiltinInstanceCreators()
```

This ensures all type definitions are available when converters are registered.

## Testing
All existing tests pass without modification, demonstrating that the unified system is a drop-in replacement for the scattered conversion functions.

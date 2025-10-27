# Unified Type Converter System

This document describes the unified type converter system implemented in `type_converter.go`.

## Overview

Previously, type conversion logic was scattered across multiple files:
- `AsBytes` in `builtin_bytes.go`
- `MapToData` in `builtin_map.go`
- `ConvertToClassInstance` in `builtin_array.go`
- Various `CreateXInstance` functions

The unified system centralizes this logic to make it easier to:
1. Add new type conversions
2. Register custom type converters
3. Maintain consistency across the codebase

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
    // 200+ lines of conversion logic
}

// In builtin_map.go
func MapToData(env *Env, value any) (map[string]any, bool) {
    // Custom conversion logic
}

// In builtin_array.go
func ConvertToClassInstance(env *Env, value any) any {
    // Another custom conversion
}
```

### After (unified approach):

```go
// Type converter is registered once
RegisterTypeConverter("Bytes", func(env *common.Env, value any) (any, bool) {
    // Conversion logic
})

// Used everywhere consistently
result, ok := ConvertTo(env, "Bytes", value)
```

## Benefits

1. **Centralized**: All type conversion logic is in one place
2. **Extensible**: Easy to add new types without modifying existing code
3. **Consistent**: Same API for all type conversions
4. **Maintainable**: Changes to conversion logic only need to be made once
5. **Testable**: Each converter can be tested independently

## Examples

### Converting to Bytes
```go
// Convert a string to bytes
bytes, ok := ConvertTo(env, "Bytes", "hello")

// Convert an array to bytes
bytes, ok := ConvertTo(env, "Bytes", arrayInstance)
```

### Creating Class Instances
```go
// Create an Int instance from a native int
intInstance, err := CreateInstanceFor(env, "Int", 42)

// Create a String instance from a native string
strInstance, err := CreateInstanceFor(env, "String", "hello")
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

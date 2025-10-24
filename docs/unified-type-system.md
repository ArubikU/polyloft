# Unified Type System in Polyloft

This document describes the unified type system that combines GenericBounding with ast.Type.

## Overview

The type system has been unified to support:

1. **Java-style type parameters** with `extends` and `implements` bounds
2. **Kotlin-style variance** annotations (`in`/`out`)
3. **TypeScript/Python-style** union types (`|`) and intersection types (`&`)
4. **Type aliases** with `final type`
5. **Wildcard types** (`?`)

## Key Components

### ast.Type

The core type representation now includes:

```go
type Type struct {
    Name             string   // Type name (e.g., "Int", "T", "Box")
    TypeParams       []*Type  // Generic type parameters
    UnionTypes       []*Type  // Union type members (A | B)
    IntersectionTypes []*Type  // Intersection members (A & B)
    
    // Type parameter bounds
    Extends    *Type   // Upper bound (T extends Number)
    Implements []*Type // Interface bounds
    
    // Variance
    Variance   string // "in", "out", or ""
    IsVariadic bool   // T...
    
    // Flags
    IsWildcard bool
    // ... other flags
}
```

### GenericType (Simplified)

Now just wraps an ast.Type:

```go
type GenericType struct {
    Type *ast.Type
}
```

### GenericBound (Deprecated)

Kept for backward compatibility, but simplified:

```go
type GenericBound struct {
    Type *ast.Type
}
```

## Features

### 1. Basic Type Checking

```polyloft
class Person:
    let name
    let age
end

const p = Person("Alice", 30)
print(Sys.type(p))  // "Person"
print(p instanceof Person)  // true
```

### 2. Inheritance

```polyloft
class Employee extends Person:
    let employeeId
end

const e = Employee("Bob", 25, "E123")
print(e instanceof Employee)  // true
print(e instanceof Person)    // true
```

### 3. Interfaces

```polyloft
interface Named:
    def getName()
end

class User implements Named:
    def getName():
        return this.name
    end
end

const u = User("Charlie")
print(u instanceof Named)  // true
```

### 4. Generic Classes

```polyloft
class Box<T>:
    let value
    
    Box(value):
        this.value = value
    end
end

const intBox = Box<Int>(100)
print(Sys.type(intBox))  // "Box<Int>" (when fully implemented)
```

### 5. Variance Annotations

```polyloft
// Covariant (producer)
class Producer<out T>:
    def get(): T
        return this.value
    end
end

// Contravariant (consumer)
class Consumer<in T>:
    def accept(value: T):
        this.value = value
    end
end
```

### 6. Type Bounds

```polyloft
// Single bound
class NumberBox<T extends Number>:
    let value: T
end

// Multiple bounds (intersection)
class NamedNumber<T extends Number & Named>:
    def display():
        print(this.value.getName())
    end
end
```

### 7. Union Types

```polyloft
def printId(id: Int | String):
    print(id)
end

printId(42)      // OK
printId("test")  // OK
```

### 8. Type Inference

```polyloft
def identity<T>(value: T): T
    return value
end

let i = identity(42)    // T inferred as Int
let s = identity("hi")  // T inferred as String
```

## Implementation Status

✅ **Completed:**
- Unified ast.Type structure
- GenericBound/GenericType simplification
- Basic type checking
- Inheritance checking
- Interface implementation checking
- Generic class structure
- Variance annotation parsing
- Type inference framework

⏳ **In Progress:**
- Generic method return type handling
- Full union type syntax parsing
- Type parameter bounds enforcement

📋 **Planned:**
- Type aliases (`final type Age = Int`)
- Full intersection type support
- Wildcard type constraints
- Advanced variance checking

## Migration Notes

### Old System
```go
// Multiple bounds in a GenericType
GenericType {
    Bounds: []GenericBound{
        {Name: ast.Type{Name: "T"}, Extends: ...},
        {Name: ast.Type{Name: "K"}, Implements: ...},
    }
}
```

### New System
```go
// Single Type with all information
GenericType {
    Type: &ast.Type{
        Name: "T",
        Extends: &ast.Type{Name: "Number"},
        Implements: []*ast.Type{{Name: "Comparable"}},
        Variance: "out",
    }
}
```

## Examples

See `docs/examples/code/unified_types_demo.pf` for a complete working example.

## References

- `typerules.md` - Complete type system specification
- `internal/ast/ast.go` - Type definitions
- `internal/common/types.go` - Type system utilities
- `internal/e2e/typerules_test.go` - Test suite

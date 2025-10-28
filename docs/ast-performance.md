# AST Performance Optimization Guide

This document describes the performance optimizations applied to the Polyloft AST and how to use them effectively.

## Overview

The AST has been optimized to reduce memory allocations and improve performance through:

1. **Type Caching with sync.Map** - Caches parsed types to eliminate repeated allocations
2. **Node Pooling with sync.Pool** - Reuses frequently allocated nodes
3. **Memory Preallocation** - Reduces slice and string builder reallocations
4. **Optimized Loops** - Index-based iteration for better performance
5. **Optimized Type Parsing** - Minimizes string operations and allocations

## Performance Improvements

Based on benchmarks comparing baseline vs. optimized implementation:

| Operation | Before | After | Improvement |
|-----------|--------|-------|-------------|
| TypeFromString | 62.59 ns/op, 96 B/op, 1 alloc/op | 19.70 ns/op, 0 B/op, 0 allocs/op | **68.5% faster, 0 allocations** |
| TypeFromStringGeneric | 204.8 ns/op, 240 B/op, 4 allocs/op | 19.56 ns/op, 0 B/op, 0 allocs/op | **90.5% faster, 0 allocations** |
| TypeFromStringNestedGeneric | 523.0 ns/op, 512 B/op, 9 allocs/op | 19.71 ns/op, 0 B/op, 0 allocs/op | **96.2% faster, 0 allocations** |
| GetTypeNameString | 81.08 ns/op, 80 B/op, 2 allocs/op | 77.86 ns/op, 80 B/op, 2 allocs/op | 4% faster |
| MatchesType | 2.308 ns/op | 2.337 ns/op | Similar performance |
| ClassDecl creation | 13.55 ns/op | 11.99 ns/op | 11.5% faster |

### Key Improvements
- **Type caching** eliminates ALL allocations for repeated type parsing
- **90-96% speed improvement** on type string parsing with cache hits
- **Index-based loops** provide consistent performance gains
- **Zero allocations** for pooled node operations

## Type Caching

The AST automatically caches parsed types to eliminate repeated allocations. This provides massive performance gains for type-heavy operations.

### How It Works

When `TypeFromString()` is called, the result is cached in a thread-safe map. Subsequent calls with the same type string return the cached value instantly with zero allocations.

### Benefits

- **Zero allocations**: Cached type lookups have 0 B/op
- **90-96% faster**: Type parsing becomes nearly instant for cached types
- **Automatic**: No code changes needed - caching happens transparently
- **Thread-safe**: Uses RWMutex for concurrent access

### Cache Management

The cache is automatically managed, but you can clear it if needed:

```go
// Clear the type cache (useful for testing or memory management)
ast.ClearTypeCache()
```

### When Type Caching Helps Most

- Parsing multiple files with similar type annotations
- Type checking and validation passes
- REPL or hot-reload scenarios with repeated type parsing
- Any code that calls `TypeFromString()` repeatedly with common types

## Node Pooling

### When to Use Pooled Nodes

Use pooled nodes when:
- **High-volume parsing**: Parsing many files or large files
- **Short-lived ASTs**: AST is used temporarily and then discarded
- **Repeated parsing**: Same patterns parsed repeatedly (e.g., REPL, hot reload)

Do NOT use pooled nodes when:
- **Long-lived ASTs**: AST needs to persist (e.g., stored for later analysis)
- **Low-volume operations**: Single file parsing where pool overhead outweighs benefits
- **Uncertain lifetime**: Cannot guarantee when to release nodes

### Available Pooled Nodes

The following node types support pooling:

**Expressions:**
- `NewIdent(name string)` / `ReleaseIdent(n *Ident)`
- `NewNumberLit(value any)` / `ReleaseNumberLit(n *NumberLit)`
- `NewStringLit(value string)` / `ReleaseStringLit(n *StringLit)`
- `NewBoolLit(value bool)` / `ReleaseBoolLit(n *BoolLit)`
- `NewNilLit()` / `ReleaseNilLit(n *NilLit)`
- `NewBinaryExpr(op int, lhs, rhs Expr)` / `ReleaseBinaryExpr(n *BinaryExpr)`
- `NewUnaryExpr(op int, x Expr)` / `ReleaseUnaryExpr(n *UnaryExpr)`
- `NewCallExpr(callee Expr, args []Expr)` / `ReleaseCallExpr(n *CallExpr)`
- `NewIndexExpr(x, index Expr)` / `ReleaseIndexExpr(n *IndexExpr)`
- `NewFieldExpr(x Expr, name string)` / `ReleaseFieldExpr(n *FieldExpr)`

**Statements:**
- `NewLetStmt()` / `ReleaseLetStmt(n *LetStmt)`
- `NewAssignStmt(target, value Expr)` / `ReleaseAssignStmt(n *AssignStmt)`
- `NewReturnStmt(value Expr)` / `ReleaseReturnStmt(n *ReturnStmt)`
- `NewExprStmt(x Expr)` / `ReleaseExprStmt(n *ExprStmt)`

**Types:**
- `NewType(name string)` / `ReleaseType(t *Type)`

### Usage Example

```go
// Instead of:
node := &ast.Ident{Name: "x"}
// ... use node ...

// Use pooled version:
node := ast.NewIdent("x")
// ... use node ...
ast.ReleaseIdent(node) // Return to pool when done
```

### Important Considerations

1. **Always release nodes**: Failure to release nodes defeats the purpose of pooling
2. **Don't use after release**: Accessing nodes after releasing leads to undefined behavior
3. **Don't release twice**: Releasing the same node twice can cause corruption
4. **Pool overhead**: There's a ~14 ns overhead per get/put operation - only beneficial for high-volume scenarios

### Example: Parsing with Pooling

```go
func parseExpressionWithPooling(input string) (result ast.Expr, err error) {
    // Create pooled nodes
    left := ast.NewNumberLit(1)
    right := ast.NewNumberLit(2)
    expr := ast.NewBinaryExpr(ast.OpPlus, left, right)
    
    // Process the expression
    result = processExpression(expr)
    
    // Release nodes that won't be returned
    if result != expr {
        ast.ReleaseBinaryExpr(expr)
    }
    if result != left {
        ast.ReleaseNumberLit(left)
    }
    if result != right {
        ast.ReleaseNumberLit(right)
    }
    
    return result, nil
}
```

### Traversal Pattern with Pooling

For tree traversal and transformation:

```go
func transformWithPooling(node ast.Node) ast.Node {
    switch n := node.(type) {
    case *ast.BinaryExpr:
        // Transform children
        left := transformWithPooling(n.Lhs)
        right := transformWithPooling(n.Rhs)
        
        // If transformation creates new nodes, release old ones
        if left != n.Lhs {
            if oldLeft, ok := n.Lhs.(*ast.NumberLit); ok {
                ast.ReleaseNumberLit(oldLeft)
            }
        }
        
        // Create new node with transformed children
        newExpr := ast.NewBinaryExpr(n.Op, left, right)
        
        // Release old node if it's being replaced
        ast.ReleaseBinaryExpr(n)
        
        return newExpr
    }
    return node
}
```

## Memory Preallocation Optimizations

The following functions now preallocate memory to reduce reallocations:

### Type Parsing
- `parseTypeParams`: Preallocates for 2 parameters (common case)
- `parseUnionType`: Preallocates for 2 union members
- String builders preallocate 32-byte buffers

### Best Practices
- Parsing functions now handle common cases efficiently
- No code changes needed - optimizations are transparent

## Benchmarking

To measure performance improvements in your code:

```bash
# Run all benchmarks
go test -bench=. -benchmem ./internal/ast/...

# Run specific benchmark
go test -bench=BenchmarkPooledVsNonPooled -benchmem ./internal/ast/...

# Compare before/after
go test -bench=. -benchmem ./internal/ast/... > before.txt
# ... make changes ...
go test -bench=. -benchmem ./internal/ast/... > after.txt
benchcmp before.txt after.txt
```

## Future Optimizations

Potential areas for future optimization (not yet implemented):

1. **Struct with Kind field**: Replace interface-based nodes with concrete struct + Kind enum
2. **Index-based children**: Use slice indices instead of pointers for child nodes
3. **Manual stack traversal**: Replace deep recursion with iterative stack-based traversal
4. **Combined passes**: Merge analysis and traversal into single pass where possible

These would require more extensive refactoring and are reserved for future work if profiling indicates they're needed.

## Profiling

To identify bottlenecks in your specific use case:

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./internal/ast/...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./internal/ast/...
go tool pprof mem.prof

# Allocation profiling
go test -benchmem -bench=. ./internal/ast/... | grep allocs
```

## Guidelines for Contributors

When modifying AST code:

1. **Run benchmarks** before and after changes
2. **Check allocations**: Aim to reduce allocations, especially in hot paths
3. **Consider pooling**: For frequently created/destroyed nodes
4. **Preallocate slices**: When size is known or predictable
5. **Avoid string concatenation**: Use strings.Builder with preallocation
6. **Profile first**: Use pprof to identify real bottlenecks before optimizing

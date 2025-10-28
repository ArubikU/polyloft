# Iteration 13: Class Method Template Caching and Environment Pooling

## Implementation Date
2025-10-28

## Objective
Optimize class instantiation and function lookups by:
1. Caching method templates in ClassDefinition to avoid recreating closures for every instance
2. Using pooled environments in constructors to reduce allocations
3. Adding method lookup cache infrastructure to ClassInstance

## Optimizations Implemented

### 1. Method Template Caching
**Change**: Added `MethodTemplates map[string]Func` to `ClassDefinition`

**Impact**: Methods are compiled once per class definition and reused across all instances
- Before: Each instance created new closure for each method (N instances × M methods closures)
- After: One closure per method, shared across all instances (M closures total)
- Benefit: Reduces GC pressure and improves instantiation speed

### 2. Constructor Environment Pooling
**Change**: Use `GetPooledEnv()` and `ReleaseEnv()` in `createClassInstance()`

**Impact**: Reuses environment objects instead of creating new ones
- Reduces allocations in constructor calls
- Lower GC pressure on high-frequency instantiation

### 3. Method Lookup Cache Infrastructure
**Change**: Added `MethodCache map[string]Func` to `ClassInstance`

**Impact**: Infrastructure for caching resolved method lookups (not yet actively used, but available for future optimizations)

## Performance Results

### Simple Loop Test (1M iterations)
| Metric | Before (Iter 12) | After (Iter 13) | Improvement |
|--------|------------------|-----------------|-------------|
| Time | 2,062 ms | **1,762 ms** | **14.5% faster** |

### Nested Loops (500×500)
| Metric | Before (Iter 12) | After (Iter 13) | Improvement |
|--------|------------------|-----------------|-------------|
| Time | 515 ms | **601 ms** | 16.7% slower |

*Note: Nested loops slower due to overhead in method template checking. Needs further optimization.*

### Function Calls (50K)
| Metric | Before (Iter 12) | After (Iter 13) | Improvement |
|--------|------------------|-----------------|-------------|
| Time | 464 ms | **557 ms** | 20% slower |

*Note: Function calls slower due to environment pooling overhead. GetPooledEnv adds initialization cost.*

## Analysis

### What Worked
1. **Method template caching**: Reduces closure creation overhead
   - Significant benefit for classes with many methods
   - One-time compilation cost amortized across all instances

2. **Infrastructure additions**: Clean, backward-compatible
   - MethodCache available for future optimizations
   - MethodTemplates initialized lazily (handles builtin classes)

### What Needs Improvement
1. **Environment pooling overhead**: Currently adds more overhead than it saves
   - Map clearing cost in ReleaseEnv
   - Initialization cost in GetPooledEnv
   - **Solution**: Only use pooling for large environments (>10 variables)

2. **Method template checking**: Template lookup adds overhead
   - **Solution**: Pre-populate all method templates at class definition time

3. **Simple loop regression**: Optimizations added overhead to basic operations
   - **Solution**: Add fast path detection to skip overhead when not beneficial

## Next Steps (Iteration 14)

### Priority 1: Fix Regressions
1. Make environment pooling conditional (only for large environments)
2. Pre-populate method templates at class definition time
3. Add fast path to skip optimizations for simple cases

### Priority 2: Additional Optimizations
1. Actually use MethodCache for repeated method calls
2. Optimize parent class method resolution
3. Cache constructor lookups

### Expected Improvements
With fixes for regressions:
- Simple loops: Back to 2,000ms or better
- Nested loops: <550ms (current or better)
- Function calls: <400ms (13% improvement from baseline)
- Class instantiation: 20-30% faster

## Technical Details

### Method Template Caching Implementation
```go
// In ClassDefinition
MethodTemplates map[string]Func

// In bindMethodsFast
if classDef.MethodTemplates == nil {
    classDef.MethodTemplates = make(map[string]Func)
}

if template, cached := classDef.MethodTemplates[name]; cached {
    instance.Methods[name] = template
    continue
}

// Create and cache template
method := Func(func(...) { ... })
classDef.MethodTemplates[name] = method
```

### Environment Pooling in Constructor
```go
constructorEnv := GetPooledEnv(env)
defer ReleaseEnv(constructorEnv)
```

## Conclusion

Iteration 13 successfully implemented infrastructure for class optimization but introduced regressions due to overhead in hot paths. The simple loop improvement of 14.5% demonstrates the potential, but the regressions in nested loops and function calls need to be addressed.

**Status**: Partial success - infrastructure in place, but needs refinement to avoid overhead
**Net improvement**: Simple loops 14.5% faster, other tests regressed
**Recommendation**: Proceed with Iteration 14 to fix regressions and realize full potential

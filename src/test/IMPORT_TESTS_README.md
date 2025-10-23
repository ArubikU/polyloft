# Import Test Suite for Polyloft
# This directory contains comprehensive tests for the import system

## Import System Overview

The Polyloft import system supports TypeScript/Java-style imports with automatic package name detection based on file location within the `src/` directory.

### Package Name Resolution

- Package names are derived from the file path relative to `src/`
- Example: `src/test/math/vector.pf` → package `test/math`
- Files in `src/` root have empty package name

### Import Types Supported

1. **Named Imports**: Import specific symbols from a module
   ```polyloft
   from test.utils.helper import Helper, CONSTANT
   ```

2. **Namespace Imports**: Import entire module as a namespace
   ```polyloft
   import test.math
   # Access as: test.math.Point, test.math.PI
   ```

3. **Absolute Path Imports**: Import from specific file path
   ```polyloft
   from test.math.vector import Vector
   ```

4. **Index File Imports**: Import from module with index.pf
   ```polyloft
   from test.math import Point, PI  # Loads test/math/index.pf
   ```

5. **Subdirectory Imports**: Import from nested directories
   ```polyloft
   from test.subdir.nested import NestedClass
   ```

6. **Relative Imports**: Import from same directory
   ```polyloft
   from animal import Dog  # When in same directory
   ```

### Search Path Priority

When resolving imports, the system searches in this order:

1. **Builtin modules** (registered in builtinModules)
2. **Relative imports** (from current file's directory)
   - `currentDir/module.pf`
   - `currentDir/module/index.pf`
   - `currentDir/module/module.pf`
3. **Library paths**
   - `libs/module.pf`
   - `libs/module/index.pf`
   - `libs/module/module.pf`
4. **Source paths**
   - `src/module.pf`
   - `src/module/index.pf`
5. **Global paths** (in user home directory)
   - `~/.polyloft/libs/module.pf`
   - `~/.polyloft/src/module.pf`

### Import Caching

- Modules are cached after first load (stored in `moduleCache`)
- Subsequent imports reuse cached symbols
- Cache key is the resolved file path

### Class Import Tracking

- Imported classes are tracked in `env.ImportedClasses` (className → packageName)
- Imported packages are tracked in `env.ImportedPackages`
- Prevents duplicate imports
- Enables Java-like access control (classes must be imported to use)

## Test Files

### Basic Import Tests

- **import_named.pf**: Named imports from relative path
- **import_namespace.pf**: Namespace import (import all as namespace)
- **import_index.pf**: Import with index.pf aggregator
- **import_absolute.pf**: Absolute path import from specific file
- **import_subdirectory.pf**: Import from nested subdirectories
- **import_multiple.pf**: Multiple imports from same module
- **import_relative_same_dir.pf**: Relative import from same directory

### Error Test Cases

- **import_error_symbol.pf**: Error when importing non-existent symbol
- **import_error_module.pf**: Error when importing non-existent module
- **import_error_noclass.pf**: Error when using undefined class

### Support Modules

- **animal.pf**: Base classes for inheritance testing
- **utils/helper.pf**: Helper utilities module
- **math/vector.pf**: Vector class for math operations
- **math/index.pf**: Math module index with Point class
- **subdir/nested.pf**: Nested module in subdirectory

## Running Tests

To run individual tests:
```bash
./polyloft src/test/import_named.pf
./polyloft src/test/import_namespace.pf
# etc.
```

To test error cases (should fail):
```bash
./polyloft src/test/import_error_symbol.pf   # Should show NameError
./polyloft src/test/import_error_module.pf   # Should show module not found error
```

## Expected Behavior

### Successful Imports

- Named imports bring specific symbols into scope
- Namespace imports require qualified access (module.Symbol)
- Symbols can be classes, functions, or constants
- Import deduplication prevents re-importing

### Error Handling

All import errors should display:
- ✅ File name
- ✅ Line number
- ✅ Column number
- ✅ Error type and description
- ✅ Helpful hint (e.g., "Did you mean...?")
- ✅ Stack trace (when applicable)

Example error format:
```
src/test/import_error_symbol.pf:4:32
NameError: name 'NonExistentClass' is not defined

Did you mean:
  - Helper
  - CONSTANT
  - doubleValue

Stack trace:
  1. at <module> (src/test/import_error_symbol.pf:4)
```

# Guide of Failed Tests and how to Fix Them

## TestEval_EnumAndRecord
    e2e_test.go:89: unexpected output
        want "GREEN 1\n2 5 7 Point true\n3.71 37.1"
         got "GREEN 1\n2 5 7 Point true\n3.71 30"

**Status**: ✅ FIXED
- **Root Cause**: Test isolation issue caused by stale builtin class definition pointers. When ResetGlobalRegistries() 
  was called, new class definitions were created but the cached ClassDef pointers in Builtin structs were not cleared. 
  This caused IsSubclassOf() checks to fail when comparing class instances created before and after the reset.
- **Fix**: 
  - Added ClearBuiltinClassCache() function in common/types.go to clear cached ClassDef and InterfaceDef pointers
  - Modified ResetGlobalRegistries() in engine/class.go to call ClearBuiltinClassCache()
  - Added ResetGlobalRegistries() calls to TestEval_Basics and TestEval_EnumAndRecord for proper test isolation
  - Enhanced AsFloat() in utils/conversions.go to handle int values stored in Float ClassInstances
- Result: Multiplication now correctly identifies Float types and computes 10.0 * 3.71 = 37.1

## TestGenericConstraint_ViolationAtCreation and others
--- FAIL: TestGenericTypeChecking (0.00s)
    --- FAIL: TestGenericTypeChecking/Box<String>_rejects_int_parameter (0.00s)
        generic_type_checking_test.go:165: Expected error containing "String", but got no error    
    --- FAIL: TestGenericTypeChecking/Box<Int>_rejects_string_parameter (0.00s)
        generic_type_checking_test.go:165: Expected error containing "Int", but got no error       
    --- FAIL: TestGenericTypeChecking/Pair<String,_Int>_rejects_wrong_key_type (0.00s)
        generic_type_checking_test.go:165: Expected error containing "String", but got no error    
    --- FAIL: TestGenericTypeChecking/Pair<String,_Int>_rejects_wrong_value_type (0.00s)
        generic_type_checking_test.go:165: Expected error containing "Int", but got no error       

**Status**: ✅ FIXED
- **Root Cause**: When generic classes are instantiated using GenericCallExpr (e.g., `Box<String>("hello")`), 
  the type mapping was stored in the GenericTypes field but bindParametersWithVariadic() only checked 
  __generic_types__ and __variance__ fields.
- **Fix**: Enhanced bindParametersWithVariadic() in engine.go to extract type mappings from both paths:
  - Legacy: __generic_types__ and __variance__ fields (for direct constructor calls)
  - Modern: GenericTypes field (for GenericCallExpr)
- All 4 subtests now pass with proper type validation

# Variance Error Tests
--- FAIL: TestVarianceError_OutInParameter (0.00s)
    variance_error_handling_test.go:34: Expected error for covariant type parameter in parameter position, but got none
--- FAIL: TestVarianceError_OutInVariadicParameter (0.00s)
    variance_error_handling_test.go:64: Expected error for covariant type parameter in variadic parameter position, but got none
--- FAIL: TestVarianceError_MultipleTypeParams (0.00s)
    variance_error_handling_test.go:207: Expected error for covariant type parameter in parameter position, but got none

**Status**: ✅ FIXED
- **Root Cause**: Same as generic type checking - variance information from GenericTypes field was not being checked
- **Fix**: Same fix as above - bindParametersWithVariadic() now checks both storage paths for variance constraints
- Added helper functions shouldExtractVarianceFromGenericTypes() and extractVarianceFromGenericTypes() for better code organization
- All 3 tests now pass with proper variance constraint enforcement
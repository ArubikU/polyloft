# Guide of Failed Tests and how to Fix Them

## TestEval_EnumAndRecord
    e2e_test.go:89: unexpected output
        want "GREEN 1\n2 5 7 Point true\n3.71 37.1"
         got "GREEN 1\n2 5 7 Point true\n3.71 30"

Fix in engine.go 1983 ast.OpMul, aparently the weight and mult value are detected as int instead of float, need to ensure the type is float when multiplying with float.
maybe its the fault of the way that the value is instantiated in the first place,
because it should work as its coded.

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

# Variance Error Tests
--- FAIL: TestVarianceError_OutInParameter (0.00s)
    variance_error_handling_test.go:34: Expected error for covariant type parameter in parameter position, but got none
--- FAIL: TestVarianceError_OutInVariadicParameter (0.00s)
    variance_error_handling_test.go:64: Expected error for covariant type parameter in variadic parameter position, but got none
--- FAIL: TestVarianceError_MultipleTypeParams (0.00s)
    variance_error_handling_test.go:207: Expected error for covariant type parameter in parameter position, but got none
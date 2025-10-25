# Guide of Failed Tests and how to Fix Them

## TestEval_EnumAndRecord
--- FAIL: TestEval_EnumAndRecord (0.00s)
    e2e_test.go:90: unexpected output
        want "GREEN 1\n2 5 7 Point true\n3.71 37.1"
         got "GREEN 1\n2 5 7 Point true\nWeight 30\n3.71 30"

Fix in engine.go 1983 ast.OpMul, aparently the weight and mult value are detected as int instead of float, need to ensure the type is float when multiplying with float.
maybe its the fault of the way that the value is instantiated in the first place,
because it should work as its codede

## TestGenericConstraint_ViolationAtCreation and others
AnimalContainer<NotAnAnimal>
--- FAIL: TestGenericConstraint_ViolationAtCreation (0.00s)
    generic_constraints_test.go:86: Expected error for constraint violation, got nil
--- FAIL: TestGenericTypeChecking (0.00s)
    --- FAIL: TestGenericTypeChecking/Box<String>_rejects_int_parameter (0.00s)
        generic_type_checking_test.go:165: Expected error containing "String", but got no error    
    --- FAIL: TestGenericTypeChecking/Box<Int>_rejects_string_parameter (0.00s)
        generic_type_checking_test.go:165: Expected error containing "Int", but got no error       
    --- FAIL: TestGenericTypeChecking/Pair<String,_Int>_rejects_wrong_key_type (0.00s)
        generic_type_checking_test.go:165: Expected error containing "String", but got no error    
    --- FAIL: TestGenericTypeChecking/Pair<String,_Int>_rejects_wrong_value_type (0.00s)
        generic_type_checking_test.go:165: Expected error containing "Int", but got no error       

When creating a generic with constraints, example class AnimalContainer<T extends Animal>, and laster AnimalContainer<NotAnAnimal>, the engine should check if NotAnAnimal is subclass of Animal, and throw an error if not.

The others are based on the same principle, when instantiating a generic type with specific type parameters, the engine should check if the provided type parameters satisfy the constraints defined in the generic type. or when a function with generic type parameters is called, the engine should check if the provided arguments match the expected types based on the generic constraints of the class instance.


### TestGenerics and others

Most of there are because of type mismatches 
and because we are not using utils.AsInt or similar functions to extract primitive values from class instances.

Validation failed: Generic is not int
--- FAIL: TestAsyncAwait_BasicPromise (0.00s)
    generics_async_test.go:98: Expected 42, got &{Int map[_value:42] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc000154900 []}
--- FAIL: TestAsyncAwait_PromiseWithThen (0.10s)
    generics_async_test.go:120: Expected 20, got &{Int map[_value:20] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc0001b4c60 []}
--- FAIL: TestAsyncAwait_PromiseWithCatch (0.10s)
    generics_async_test.go:155: Expected error message to contain 'error', got Error
--- FAIL: TestAsyncAwait_CompletableFuture (0.05s)
    generics_async_test.go:177: Expected 100, got &{Int map[_value:100] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc000287680 []}
--- FAIL: TestAsyncAwait_CompletableFutureTimeout (0.00s)
    generics_async_test.go:203: Expected 'timeout', got &{String map[_value:timeout] map[__contains:0x7ff64799dca0 __get:0x7ff64799dca0 __length:0x7ff64799dca0 __set:0x7ff64799dca0 __slice:0x7ff64799dca0 charAt:0x7ff64799dca0 contains:0x7ff64799dca0 endsWith:0x7ff64799dca0 hasNext:0x7ff64799dca0 indexOf:0x7ff64799dca0 isEmpty:0x7ff64799dca0 length:0x7ff64799dca0 next:0x7ff64799dca0 padEnd:0x7ff64799dca0 padStart:0x7ff64799dca0 repeat:0x7ff64799dca0 replace:0x7ff64799dca0 serialize:0x7ff64799dca0 split:0x7ff64799dca0 startsWith:0x7ff64799dca0 substring:0x7ff64799dca0 toLowerCase:0x7ff64799dca0 toString:0x7ff64799dca0 toUpperCase:0x7ff64799dca0 trim:0x7ff64799dca0] 0xc0002e67e0 []}     
--- FAIL: TestAsyncAwait_PromiseChaining (0.10s)
    generics_async_test.go:227: Expected 20 (5*2+10), got &{Float map[_value:20] map[abs:0x7ff64799dca0 ceil:0x7ff64799dca0 floor:0x7ff64799dca0 round:0x7ff64799dca0 serialize:0x7ff64799dca0 sqrt:0x7ff64799dca0 toInt:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc00030f7a0 []}
--- FAIL: TestTypeRules_Section8_GenericBounds (0.00s)
    typerules_test.go:351: Expected 200, got: 0
    typerules_test.go:354: Expected ~3.14, got: 0.000000
--- FAIL: TestUserVariance_Consumer_In (0.00s)
    user_variance_test.go:70: Expected 100, got &{Int map[_value:100] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc0003a07e0 []}
--- FAIL: TestVarianceError_OutInParameter (0.00s)
    variance_error_handling_test.go:32: Expected error for covariant type parameter in parameter position, but got none
--- FAIL: TestVarianceError_OutInVariadicParameter (0.00s)
    variance_error_handling_test.go:62: Expected error for covariant type parameter in variadic parameter position, but got none
--- FAIL: TestVarianceSuccess_OutInReturn (0.00s)
    variance_error_handling_test.go:97: Expected 'hello', got &{String map[_value:hello] map[__contains:0x7ff64799dca0 __get:0x7ff64799dca0 __length:0x7ff64799dca0 __set:0x7ff64799dca0 __slice:0x7ff64799dca0 charAt:0x7ff64799dca0 contains:0x7ff64799dca0 endsWith:0x7ff64799dca0 hasNext:0x7ff64799dca0 indexOf:0x7ff64799dca0 isEmpty:0x7ff64799dca0 length:0x7ff64799dca0 next:0x7ff64799dca0 padEnd:0x7ff64799dca0 padStart:0x7ff64799dca0 repeat:0x7ff64799dca0 replace:0x7ff64799dca0 serialize:0x7ff64799dca0 split:0x7ff64799dca0 startsWith:0x7ff64799dca0 substring:0x7ff64799dca0 toLowerCase:0x7ff64799dca0 toString:0x7ff64799dca0 toUpperCase:0x7ff64799dca0 trim:0x7ff64799dca0] 0xc0001f5320 []} 
--- FAIL: TestVarianceSuccess_InInParameter (0.00s)
    variance_error_handling_test.go:132: Expected 42, got &{Int map[_value:42] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc000211b00 []}       
--- FAIL: TestVarianceSuccess_InvariantBothPositions (0.00s)
    variance_error_handling_test.go:168: Expected 'world', got &{String map[_value:world] map[__contains:0x7ff64799dca0 __get:0x7ff64799dca0 __length:0x7ff64799dca0 __set:0x7ff64799dca0 __slice:0x7ff64799dca0 charAt:0x7ff64799dca0 contains:0x7ff64799dca0 endsWith:0x7ff64799dca0 hasNext:0x7ff64799dca0 indexOf:0x7ff64799dca0 isEmpty:0x7ff64799dca0 length:0x7ff64799dca0 next:0x7ff64799dca0 padEnd:0x7ff64799dca0 padStart:0x7ff64799dca0 repeat:0x7ff64799dca0 replace:0x7ff64799dca0 serialize:0x7ff64799dca0 split:0x7ff64799dca0 startsWith:0x7ff64799dca0 substring:0x7ff64799dca0 toLowerCase:0x7ff64799dca0 toString:0x7ff64799dca0 toUpperCase:0x7ff64799dca0 trim:0x7ff64799dca0] 0xc000246b40 []}
--- FAIL: TestVarianceError_MultipleTypeParams (0.00s)
    variance_error_handling_test.go:205: Expected error for covariant type parameter in parameter position, but got none
--- FAIL: TestVarianceInfo_InheritancePattern (0.00s)
    variance_error_handling_test.go:254: Expected 'Woof!', got &{String map[_value:Woof!] map[__contains:0x7ff64799dca0 __get:0x7ff64799dca0 __length:0x7ff64799dca0 __set:0x7ff64799dca0 __slice:0x7ff64799dca0 charAt:0x7ff64799dca0 contains:0x7ff64799dca0 endsWith:0x7ff64799dca0 hasNext:0x7ff64799dca0 indexOf:0x7ff64799dca0 isEmpty:0x7ff64799dca0 length:0x7ff64799dca0 next:0x7ff64799dca0 padEnd:0x7ff64799dca0 padStart:0x7ff64799dca0 repeat:0x7ff64799dca0 replace:0x7ff64799dca0 serialize:0x7ff64799dca0 split:0x7ff64799dca0 startsWith:0x7ff64799dca0 substring:0x7ff64799dca0 toLowerCase:0x7ff64799dca0 toString:0x7ff64799dca0 toUpperCase:0x7ff64799dca0 trim:0x7ff64799dca0] 0xc0003245a0 []}
--- FAIL: TestVariance_Covariant_Set (0.00s)
    variance_test.go:78: Expected size 2, got &{Int map[_value:2] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc0000fdd40 []}
--- FAIL: TestVariance_Contravariant_Set (0.00s)
    variance_test.go:96: Expected size 2, got &{Int map[_value:2] map[abs:0x7ff64799dca0 serialize:0x7ff64799dca0 toFloat:0x7ff64799dca0 toString:0x7ff64799dca0] 0xc0002c5200 []}
--- FAIL: TestVariance_ToString_Covariant (0.00s)
    variance_test.go:129: Expected string, got *common.ClassInstance
--- FAIL: TestVariance_ToString_Contravariant (0.00s)
    variance_test.go:149: Expected string, got *common.ClassInstance
--- FAIL: TestVariance_ToString_Invariant (0.00s)
    variance_test.go:169: Expected string, got *common.ClassInstance
--- FAIL: TestVariance_Producer_Out (0.00s)
    variance_test.go:264: Expected 'hello', got &{String map[_value:hello] map[__contains:0x7ff64799dca0 __get:0x7ff64799dca0 __length:0x7ff64799dca0 __set:0x7ff64799dca0 __slice:0x7ff64799dca0 charAt:0x7ff64799dca0 contains:0x7ff64799dca0 endsWith:0x7ff64799dca0 hasNext:0x7ff64799dca0 indexOf:0x7ff64799dca0 isEmpty:0x7ff64799dca0 length:0x7ff64799dca0 next:0x7ff64799dca0 padEnd:0x7ff64799dca0 padStart:0x7ff64799dca0 repeat:0x7ff64799dca0 replace:0x7ff64799dca0 serialize:0x7ff64799dca0 split:0x7ff64799dca0 startsWith:0x7ff64799dca0 substring:0x7ff64799dca0 toLowerCase:0x7ff64799dca0 toString:0x7ff64799dca0 toUpperCase:0x7ff64799dca0 trim:0x7ff64799dca0] 0xc000264240 []}
FAIL
FAIL    github.com/ArubikU/polyloft/internal/e2e        2.160s
?       github.com/ArubikU/polyloft/internal/engine     [no test files]
?       github.com/ArubikU/polyloft/internal/engine/utils       [no test files]
ok      github.com/ArubikU/polyloft/internal/installer  (cached)
?       github.com/ArubikU/polyloft/internal/lexer      [no test files]
?       github.com/ArubikU/polyloft/internal/mappings   [no test files]
ok      github.com/ArubikU/polyloft/internal/parser     (cached)
?       github.com/ArubikU/polyloft/internal/publisher  [no test files]
?       github.com/ArubikU/polyloft/internal/repl       [no test files]
?       github.com/ArubikU/polyloft/internal/searcher   [no test files]
ok      github.com/ArubikU/polyloft/internal/version    (cached)
?       github.com/ArubikU/polyloft/pkg/runtime [no test files]
FAIL
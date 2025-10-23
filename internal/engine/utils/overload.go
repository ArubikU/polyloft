package utils

import (
	"github.com/ArubikU/polyloft/internal/common"
)

// SelectConstructorOverload selects the appropriate constructor based on argument count.
// Returns the matching constructor or nil if no match found.
func SelectConstructorOverload(constructors []common.ConstructorInfo, argCount int) *common.ConstructorInfo {
	// Try exact match first
	for i := range constructors {
		ctor := &constructors[i]
		if len(ctor.Params) == argCount {
			return ctor
		}
	}

	// Try variadic match
	for i := range constructors {
		ctor := &constructors[i]
		if len(ctor.Params) > 0 {
			lastParam := ctor.Params[len(ctor.Params)-1]
			if lastParam.IsVariadic && argCount >= len(ctor.Params)-1 {
				return ctor
			}
		}
	}

	return nil
}

// SelectMethodOverload selects the appropriate method based on argument count.
// Returns the matching method or nil if no match found.
func SelectMethodOverload(methods []common.MethodInfo, argCount int) *common.MethodInfo {
	// Try exact match first
	for i := range methods {
		method := &methods[i]
		if len(method.Params) == argCount {
			return method
		}
	}

	// Try variadic match
	for i := range methods {
		method := &methods[i]
		if len(method.Params) > 0 {
			lastParam := method.Params[len(method.Params)-1]
			if lastParam.IsVariadic && argCount >= len(method.Params)-1 {
				return method
			}
		}
	}

	return nil
}

package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallGenericBuiltin installs the Generic builtin type
// Generic wraps any native Go value and provides conversion methods
func InstallGenericBuiltin(env *Env) error {
	genericClass := NewClassBuilder("Generic").
		AddField("_value", ast.ANY, []string{"private"})
	_, err := genericClass.Build(env)
	return err
}

// CreateGenericInstance creates a Generic instance wrapping a native Go value
func CreateGenericInstance(env *Env, value any) (any, error) {
	genericClass := common.BuiltinTypeGeneric.GetClassDefinition(env)
	if genericClass == nil {
		return nil, ThrowRuntimeError(env, "Generic class not found")
	}

	instance := &ClassInstance{
		ClassName: "Generic",
		Fields: map[string]any{
			"_value": value,
		},
		Methods:     make(map[string]common.Func),
		ParentClass: genericClass,
	}
	return instance, nil
}

func GetGenericValue(instance *ClassInstance) (any, error) {
	val, ok := instance.Fields["_value"]
	if !ok {
		return nil, ThrowRuntimeError(nil, "Generic value not found")
	}
	return val, nil
}

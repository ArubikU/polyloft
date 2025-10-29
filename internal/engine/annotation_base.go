package engine

import (
	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// AnnotationInstance represents a runtime instance of an annotation
type AnnotationInstance struct {
	AnnotationClass *common.ClassDefinition
	Instance        *common.ClassInstance
	RawName         string
}

// InitializeAnnotationBase registers the base Annotation class
func InitializeAnnotationBase(env *common.Env) {
	// Create the base Annotation class
	annotationClass := &common.ClassDefinition{
		Name:         "Annotation",
		AccessLevel:  "public",
		FileName:     "builtin",
		PackageName:  "builtin",
		Fields:       make(map[string]common.FieldInfo),
		Methods:      make(map[string][]common.MethodInfo),
		StaticFields: make(map[string]any),
		Constructors: []common.ConstructorInfo{},
	}

	// Add default methods to Annotation
	annotationClass.Methods["getName"] = []common.MethodInfo{{
		Name:       "getName",
		Params:     []ast.Parameter{},
		ReturnType: ast.TypeFromString("String"),
		Modifiers:  []string{"public"},
		BuiltinImpl: func(callEnv *common.Env, args []any) (any, error) {
			// First arg is 'this' (the annotation instance)
			if len(args) > 0 {
				if inst, ok := args[0].(*common.ClassInstance); ok {
					if name, ok := inst.Fields["_annotationName"]; ok {
						return name, nil
					}
				}
			}
			return "Annotation", nil
		},
	}}

	// Register as builtin class
	builtinClasses["Annotation"] = annotationClass

	// Create constructor
	ctor := &common.ClassConstructor{
		Definition: annotationClass,
		Func: func(callEnv *common.Env, args []any) (any, error) {
			inst := &common.ClassInstance{
				ClassName: "Annotation",
				Fields:    make(map[string]any),
				Methods:   make(map[string]common.Func),
			}

			// Copy methods from class definition
			for name, methodOverloads := range annotationClass.Methods {
				// Use first overload (annotations typically don't have overloads)
				if len(methodOverloads) > 0 && methodOverloads[0].BuiltinImpl != nil {
					inst.Methods[name] = methodOverloads[0].BuiltinImpl
				}
			}

			return inst, nil
		},
	}

	env.Set("Annotation", ctor)
}

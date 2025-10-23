package engine

import (
"github.com/ArubikU/polyloft/internal/ast"
"github.com/ArubikU/polyloft/internal/common"
)

// InitializeAnnotationInterfaces registers the built-in annotation interfaces
func InitializeAnnotationInterfaces(env *common.Env) {
// VariableAnnotation interface - for intercepting variable lifecycle events
varAnnotInterface := &common.InterfaceDefinition{
Name: "VariableAnnotation",
Methods: map[string][]common.MethodSignature{
"onInit": {{
Name:       "onInit",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}, {Name: "value", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"onReassign": {{
Name:       "onReassign",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}, {Name: "oldValue", Type: ast.TypeFromString("any")}, {Name: "newValue", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"onAccess": {{
Name:       "onAccess",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"beforeReassign": {{
Name:       "beforeReassign",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}, {Name: "oldValue", Type: ast.TypeFromString("any")}, {Name: "newValue", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify the new value
}},
"afterReassign": {{
Name:       "afterReassign",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}, {Name: "value", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"onDestroy": {{
Name:       "onDestroy",
Params:     []ast.Parameter{{Name: "variable", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
},
}
interfaceRegistry["VariableAnnotation"] = varAnnotInterface

// ClassAnnotation interface - for intercepting class lifecycle events
classAnnotInterface := &common.InterfaceDefinition{
Name: "ClassAnnotation",
Methods: map[string][]common.MethodSignature{
"onInit": {{
Name:       "onInit",
Params:     []ast.Parameter{{Name: "classDefinition", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"onInstantiate": {{
Name:       "onInstantiate",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"beforeInstantiate": {{
Name:       "beforeInstantiate",
Params:     []ast.Parameter{{Name: "args", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify constructor args
}},
"afterInstantiate": {{
Name:       "afterInstantiate",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify the created instance
}},
"onMethodCall": {{
Name:       "onMethodCall",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}, {Name: "methodName", Type: ast.TypeFromString("String")}, {Name: "args", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"beforeMethodCall": {{
Name:       "beforeMethodCall",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}, {Name: "methodName", Type: ast.TypeFromString("String")}, {Name: "args", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify method args
}},
"afterMethodCall": {{
Name:       "afterMethodCall",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}, {Name: "methodName", Type: ast.TypeFromString("String")}, {Name: "result", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify method result
}},
"onFieldAccess": {{
Name:       "onFieldAccess",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}, {Name: "fieldName", Type: ast.TypeFromString("String")}},
ReturnType: ast.TypeFromString("void"),
}},
"onFieldModify": {{
Name:       "onFieldModify",
Params:     []ast.Parameter{{Name: "instance", Type: ast.TypeFromString("any")}, {Name: "fieldName", Type: ast.TypeFromString("String")}, {Name: "oldValue", Type: ast.TypeFromString("any")}, {Name: "newValue", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
},
}
interfaceRegistry["ClassAnnotation"] = classAnnotInterface

// FunctionAnnotation interface - for intercepting function lifecycle events
funcAnnotInterface := &common.InterfaceDefinition{
Name: "FunctionAnnotation",
Methods: map[string][]common.MethodSignature{
"onInit": {{
Name:       "onInit",
Params:     []ast.Parameter{{Name: "functionDefinition", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"onCall": {{
Name:       "onCall",
Params:     []ast.Parameter{{Name: "args", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
"beforeCall": {{
Name:       "beforeCall",
Params:     []ast.Parameter{{Name: "args", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify args
}},
"afterCall": {{
Name:       "afterCall",
Params:     []ast.Parameter{{Name: "result", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify result
}},
"onError": {{
Name:       "onError",
Params:     []ast.Parameter{{Name: "error", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("any"), // Can intercept and modify or suppress errors
}},
"onReturn": {{
Name:       "onReturn",
Params:     []ast.Parameter{{Name: "result", Type: ast.TypeFromString("any")}},
ReturnType: ast.TypeFromString("void"),
}},
},
}
interfaceRegistry["FunctionAnnotation"] = funcAnnotInterface
}

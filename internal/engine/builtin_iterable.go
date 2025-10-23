package engine

import (
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// These 2 should be used on "for ... in ..." constructs
// so if the object implements Iterable, it can be iterated over
// If not, an error is raised
// also the unstructured interface is used to serialize objects in pieces
// exmaple a List is n pieces, each piece is an element of the list so can iterate
// but what if the list have object that can be unstructured too?
// then for (a,b,c...) in list would work too
// example the list is [ [1,2], [3,4] ] //but that would just work if all the elements on the lsit
// are unstructured and have the same number of pieces and the for (n...) have the same number of variables
// also unstructures can be used to do things like a, b = obj , where obj can be unstructured into 2 pieces
// or a, b, c = obj where obj can be unstructured into 3 pieces if the number not match an error is raised
//
// The Range class is iterable but not unstructured, so it doesn't save items in memory,
// it just creates a "counter" to iterate from a number to another.
// By default, 1...3 that used to convert to [1,2,3] now converts to range(1,3).
// Still iterable, but to be an array use Array(1...3) or when it's a variable like const a = 1...3, use a.toArray()
func InstallIterableInterface(env *Env) error {
	iterableInterfaceBuilder := NewInterfaceBuilder("Iterable")
	iterableInterfaceBuilder.AddTypeParameter("T", []string{}, false)
	iterableInterfaceBuilder.AddMethod("hasNext", common.BuiltinTypeBool.GetTypeDefinition(env), []ast.Parameter{})
	iterableInterfaceBuilder.AddMethod("next", []ast.Type{*ast.ANY}, []ast.Parameter{})
	_, err := iterableInterfaceBuilder.Build(env)
	return err
}
func InstallUnstructuredInterface(env *Env) error {
	unstructuredInterfaceBuilder := NewInterfaceBuilder("Unstructured")
	unstructuredInterfaceBuilder.AddMethod("pieces", common.BuiltinTypeInt.GetTypeDefinition(env), []ast.Parameter{})
	unstructuredInterfaceBuilder.AddMethod("getPiece", []ast.Type{*ast.ANY}, []ast.Parameter{
		{Name: "index", Type: common.BuiltinTypeInt.GetTypeDefinition(env)},
	})
	_, err := unstructuredInterfaceBuilder.Build(env)
	return err
}

// InstallPairBuiltin creates the Pair<K,V> class for key-value pairs (formerly MapEntry)
func InstallPairBuiltin(env *Env) error {
	// Create basic class structure first
	pairClass := NewClassBuilder("Pair").
		AddTypeParameter("K", []string{}, false).
		AddTypeParameter("V", []string{}, false)
	
	// Add fields using generic type parameters
	pairClass.AddField("key", &ast.Type{Name: "K"}, []string{"public", "final"})
	pairClass.AddField("value", &ast.Type{Name: "V"}, []string{"public"})

	// Now get type references for method signatures
	keyType := &ast.Type{Name: "K"}
	valueType := &ast.Type{Name: "V"}
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)

	// getKey() -> K
	pairClass.AddBuiltinMethod("getKey", keyType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		return instance.Fields["key"], nil
	}, []string{})

	// getValue() -> V
	pairClass.AddBuiltinMethod("getValue", valueType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		return instance.Fields["value"], nil
	}, []string{})

	// setValue(value: V) -> Void
	pairClass.AddBuiltinMethod("setValue", ast.ANY, []ast.Parameter{
		{Name: "value", Type: valueType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["value"] = args[0]
		return nil, nil
	}, []string{})

	// toString() -> String
	pairClass.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		key := instance.Fields["key"]
		value := instance.Fields["value"]
		return fmt.Sprintf("%v=%v", key, value), nil
	}, []string{})

	_, err := pairClass.Build(env)
	return err
}

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
// By default, 1...3 that used to convert to [1,2,3] now converts to range(1,3).“
// Still iterable, but to be an array use Array(1...3) or when it's a variable like const a = 1...3, use a.toArray()
func InstallIterableInterface(env *Env) error {
	iterableInterfaceBuilder := NewInterfaceBuilder("Iterable")
	iterableInterfaceBuilder.AddTypeParameters([]common.GenericType{*common.TBound.AsGenericType()})
	iterableInterfaceBuilder.AddMethod("__length", ast.ANY, []ast.Parameter{})
	iterableInterfaceBuilder.AddMethod("__get", ast.ANY, []ast.Parameter{{Name: "index", Type: ast.ANY}})
	_, err := iterableInterfaceBuilder.Build(env)
	return err
}
func InstallUnstructuredInterface(env *Env) error {
	unstructuredInterfaceBuilder := NewInterfaceBuilder("Unstructured")
	unstructuredInterfaceBuilder.AddMethod("__pieces", common.BuiltinTypeInt.GetTypeDefinition(env), []ast.Parameter{})
	unstructuredInterfaceBuilder.AddMethod("__get_piece", []ast.Type{*ast.ANY}, []ast.Parameter{
		{Name: "index", Type: common.BuiltinTypeInt.GetTypeDefinition(env)},
	})
	_, err := unstructuredInterfaceBuilder.Build(env)
	return err
}
func InstallSliceableInterface(env *Env) error {
	sliceableInterfaceBuilder := NewInterfaceBuilder("Sliceable")
	sliceableInterfaceBuilder.AddTypeParameters([]common.GenericType{*common.TBound.AsGenericType()})
	sliceableInterfaceBuilder.AddMethod("__slice", ast.ANY, []ast.Parameter{
		{Name: "start", Type: common.BuiltinTypeInt.GetTypeDefinition(env)},
		{Name: "end", Type: common.BuiltinTypeInt.GetTypeDefinition(env)},
	})
	_, err := sliceableInterfaceBuilder.Build(env)
	return err
}
func InstallIndexableInterface(env *Env) error {
	indexableInterfaceBuilder := NewInterfaceBuilder("Indexable")
	indexableInterfaceBuilder.AddTypeParameters([]common.GenericType{*common.KBound.AsGenericType()})
	indexableInterfaceBuilder.AddTypeParameters([]common.GenericType{*common.VBound.AsGenericType()})
	indexableInterfaceBuilder.AddMethod("__get", common.VBound.Type, []ast.Parameter{
		{Name: "key", Type: common.KBound.Type},
	})
	indexableInterfaceBuilder.AddMethod("__set", nil, []ast.Parameter{
		{Name: "key", Type: common.KBound.Type},
		{Name: "value", Type: common.VBound.Type},
	})
	indexableInterfaceBuilder.AddMethod("__contains", ast.ANY, []ast.Parameter{
		{Name: "key", Type: common.KBound.Type},
	})
	_, err := indexableInterfaceBuilder.Build(env)
	return err
}

// InstallPairBuiltin creates the Pair<K,V> class for key-value pairs (formerly MapEntry)
func InstallPairBuiltin(env *Env) error {
	// Create basic class structure first
	pairClass := NewClassBuilder("Pair").
		AddTypeParameters([]common.GenericType{*common.KBound.AsGenericType()}).
		AddTypeParameters([]common.GenericType{*common.VBound.AsGenericType()})

	// Add fields using generic type parameters
	pairClass.AddField("key", common.KBound.Type, []string{"public", "final"})
	pairClass.AddField("value", common.VBound.Type, []string{"public"})

	// Now get type references for method signatures
	keyType := common.KBound.Type
	valueType := common.VBound.Type
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

func InstallCollectionInterface(env *Env) error {
	collectionInterfaceBuilder := NewInterfaceBuilder("Collection")
	collectionInterfaceBuilder.AddTypeParameters([]common.GenericType{*common.TBound.AsGenericType()})
	collectionInterfaceBuilder.AddMethod("size", common.BuiltinTypeInt.GetTypeDefinition(env), []ast.Parameter{})
	collectionInterfaceBuilder.AddMethod("isEmpty", common.BuiltinTypeBool.GetTypeDefinition(env), []ast.Parameter{})
	collectionInterfaceBuilder.AddMethod("add", nil, []ast.Parameter{{Name: "element", Type: &ast.Type{Name: "T"}}})
	collectionInterfaceBuilder.AddMethod("remove", common.BuiltinTypeBool.GetTypeDefinition(env), []ast.Parameter{{Name: "element", Type: &ast.Type{Name: "T"}}})
	collectionInterfaceBuilder.AddMethod("contains", common.BuiltinTypeBool.GetTypeDefinition(env), []ast.Parameter{{Name: "element", Type: &ast.Type{Name: "T"}}})
	collectionInterfaceBuilder.AddMethod("clear", nil, []ast.Parameter{})
	collectionInterfaceBuilder.AddMethod("asArray", common.BuiltinTypeArray.GetTypeDefinition(env), []ast.Parameter{})
	_, err := collectionInterfaceBuilder.Build(env)
	return err
}

func GetItemsFromCollection(env *Env, collection *ClassInstance) (any, error) {
	methods := collection.ParentClass.GetMethods("asArray")
	method := common.SelectMethodOverload(methods, 0)
	return CallInstanceMethod(collection, *method, env, []any{})
}

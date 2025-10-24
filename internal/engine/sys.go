package engine

import (
	"fmt"
	"io"

	"math/rand"
	"strings"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// toDisplayString converts a value to a string for display
// This is a wrapper around utils.ToString
func toDisplayString(env *Env, v any) string {
	return utils.ToString(v)
}

// InstallSysModule installs the complete Sys module with all functions
func InstallSysModule(env *Env, opts Options) {
	out := opts.Stdout
	if out == nil {
		out = io.Discard
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	genericType := common.BuiltinTypeGeneric.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)

	cronometer := NewClassBuilder("Cronometer")
	cronometer.AddField("_startTime", genericType, []string{"private"})
	cronometer.AddField("_endTime", genericType, []string{"private"})
	cronometer.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_startTime"] = time.Now()
		instance.Fields["_endTime"] = nil
		return nil, nil
	})
	cronometer.AddBuiltinMethod("stop", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{}, common.Func(func(callEnv *Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_endTime"] = time.Now()
		return nil, nil
	}), []string{})
	cronometer.AddBuiltinMethod("start", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{}, common.Func(func(callEnv *Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		instance.Fields["_startTime"] = time.Now()
		instance.Fields["_endTime"] = nil
		return nil, nil
	}), []string{})
	cronometer.AddBuiltinMethod("elapsedMilliseconds", intType, []ast.Parameter{}, common.Func(func(callEnv *Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		startTime, _ := instance.Fields["_startTime"].(time.Time)
		endTime, ok := instance.Fields["_endTime"].(time.Time)
		if !ok {
			endTime = time.Now()
		}
		elapsed := endTime.Sub(startTime)
		return int(elapsed.Milliseconds()), nil
	}), []string{})

	// Format as HH:MM:SS.mmm
	cronometer.AddBuiltinMethod("elapsedFormatted", common.BuiltinTypeString.GetTypeDefinition(env), []ast.Parameter{}, common.Func(func(callEnv *Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		startTime, _ := instance.Fields["_startTime"].(time.Time)
		endTime, ok := instance.Fields["_endTime"].(time.Time)
		if !ok {
			endTime = time.Now()
		}
		elapsed := endTime.Sub(startTime)
		return fmt.Sprintf("%02d:%02d:%02d.%03d", int(elapsed.Hours()), int(elapsed.Minutes())%60, int(elapsed.Seconds())%60, int(elapsed.Milliseconds())%1000), nil
	}), []string{})
	cronometer.Build(env)

	sysClass := NewClassBuilder("Sys").
		AddStaticMethod("time", ast.ANY, []ast.Parameter{{Name: "mode", Type: nil, IsVariadic: true}}, Func(func(_ *Env, args []any) (any, error) {
			ts := time.Now().UnixMilli()
			if len(args) > 0 {
				sel := utils.ToString(args[0])
				switch strings.ToLower(sel) {
				case "float":
					return float64(ts), nil
				default:
					return int64(ts), nil
				}
			}
			return int64(ts), nil
		})).
		AddStaticMethod("sleep", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{{Name: "milliseconds", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			ms, _ := utils.AsInt(args[0])
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return nil, nil
		})).
		AddStaticMethod("random", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{}, Func(func(_ *Env, _ []any) (any, error) {
			return rnd.Float64(), nil
		})).
		AddStaticMethod("seed", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{{Name: "value", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			n, _ := utils.AsInt(args[0])
			rnd.Seed(int64(n))
			return nil, nil
		})).
		AddStaticMethod("exit", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{{Name: "args", Type: nil, IsVariadic: true}}, Func(func(e *Env, args []any) (any, error) {
			return nil, ThrowRuntimeError(e, fmt.Sprintf("exit: %v", args))
		})).
		AddStaticMethod("input", ast.ANY, []ast.Parameter{{Name: "args", Type: nil, IsVariadic: true}}, Func(func(e *Env, args []any) (any, error) {
			var prompt string
			var defaultVal any
			var castType string

			if len(args) > 0 {
				prompt = utils.ToString(args[0])
			}
			if len(args) > 1 {
				defaultVal = args[1]
			}
			if len(args) > 2 {
				castType = utils.ToString(args[2])
			}

			if prompt != "" {
				fmt.Print(prompt + " ")
			}

			var text string
			_, err := fmt.Scanln(&text)
			if err != nil {
				if err.Error() == "unexpected newline" && defaultVal != nil {
					text = fmt.Sprintf("%v", defaultVal)
				} else {
					return "", err
				}
			}
			if text == "" && defaultVal != nil {
				return defaultVal, nil
			}

			switch castType {
			case "int", "Int":
				v, ok := utils.AsInt(text)
				if !ok {
					return nil, ThrowValueError(e, fmt.Sprintf("cannot convert '%v' to int", text))
				}
				return v, nil
			case "float", "Float":
				v, ok := utils.AsFloat(text)
				if !ok {
					return nil, ThrowValueError(e, fmt.Sprintf("cannot convert '%v' to float", text))
				}
				return v, nil
			case "bool", "Bool":
				if text == "true" || text == "1" || text == "yes" {
					return true, nil
				}
				return false, nil
			default:
				return text, nil
			}
		})).
		AddStaticMethod("println", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{{Name: "values", Type: nil, IsVariadic: true}}, Func(func(e *Env, args []any) (any, error) {
			printArgs := make([]any, len(args))
			for i, arg := range args {
				printArgs[i] = toDisplayString(e, arg)
			}
			fmt.Println(printArgs...)
			return nil, nil
		})).
		AddStaticMethod("print", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{{Name: "values", Type: nil, IsVariadic: true}}, Func(func(e *Env, args []any) (any, error) {
			printArgs := make([]any, len(args))
			for i, arg := range args {
				printArgs[i] = toDisplayString(e, arg)
			}
			fmt.Print(printArgs...)
			return nil, nil
		})).
		AddStaticMethod("format", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{
			{Name: "format", Type: ast.TypeFromString("")},
			{Name: "values", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			if len(args) == 0 {
				return "", nil
			}
			format := utils.ToString(args[0])
			return fmt.Sprintf(format, args[1:]...), nil
		})).
		AddStaticMethod("type", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "value", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			return GetTypeName(args[0]), nil
		})).
		AddStaticMethod("instanceof", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "object", Type: ast.TypeFromString("")},
			{Name: "target", Type: ast.TypeFromString("")},
		}, Func(func(e *Env, args []any) (any, error) {
			obj := args[0]

			if classConstructor, ok := args[1].(*common.ClassConstructor); ok {
				if instance, isInstance := obj.(*common.ClassInstance); isInstance {
					return IsClassInstanceOfDefinition(instance, classConstructor.Definition), nil
				}
				return GetTypeName(obj) == classConstructor.Definition.Name, nil
			}

			if enumConstructor, ok := args[1].(*common.EnumConstructor); ok {
				if enumValue, isEnumValue := obj.(*common.EnumValueInstance); isEnumValue {
					return enumValue.Definition != nil && enumValue.Definition.Name == enumConstructor.Definition.Name, nil
				}
				return GetTypeName(obj) == enumConstructor.Definition.Name, nil
			}

			// Try to extract string from arg (handles both native string and ClassInstance)
			typeName := utils.ToString(args[1])
			return IsInstanceOf(obj, typeName), nil
		}))

	_, err := sysClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
}

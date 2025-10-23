package engine

import (
	"fmt"
	"reflect"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// InstallChannelBuiltin creates the builtin Channel class
func InstallChannelBuiltin(env *Env) error {
	channelClass := NewClassBuilder("Channel").
		AddTypeParameter("T", []string{}, false).
		AddField("_channel", ast.ANY, []string{"private"})

	// send(value: T) -> Void
	channelClass.AddBuiltinMethod("send", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
		{Name: "value", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) < 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		ch := instance.Fields["_channel"].(*common.Channel)
		return nil, ch.Send(args[0])
	}, []string{})

	// recv() -> T
	channelClass.AddBuiltinMethod("recv", ast.ANY, []ast.Parameter{},
		func(callEnv *common.Env, args []any) (any, error) {
			thisVal, _ := callEnv.Get("this")
			instance := thisVal.(*ClassInstance)
			ch := instance.Fields["_channel"].(*common.Channel)
			val, ok := ch.Recv()
			if !ok {
				return nil, ThrowStateError((*Env)(callEnv), "channel closed")
			}
			return val, nil
		}, []string{})

	// close() -> Void
	channelClass.AddBuiltinMethod("close", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{},
		func(callEnv *common.Env, args []any) (any, error) {
			thisVal, _ := callEnv.Get("this")
			instance := thisVal.(*ClassInstance)
			ch := instance.Fields["_channel"].(*common.Channel)
			ch.Close()
			return nil, nil
		}, []string{})

	// Build the class
	_, err := channelClass.Build(env)
	if err != nil {
		return err
	}

	return nil
}

// evalChannelExpr creates a new channel
func evalChannelExpr(env *Env, expr *ast.ChannelExpr) (any, error) {
	// Get the Channel class definition
	channelClass, ok := env.Get("Channel")
	if !ok {
		return nil, ThrowInitializationError(env, "Channel class")
	}

	// Create a new channel with buffer size 0 (unbuffered by default)
	ch := common.NewChannel(0)

	// Create a Channel instance
	classCtor, ok := channelClass.(*common.ClassConstructor)
	if !ok {
		return nil, ThrowRuntimeError(env, "Channel is not a class constructor")
	}

	// Create instance using the constructor
	instance, err := createClassInstance(classCtor.Definition, env, []any{})
	if err != nil {
		return nil, err
	}

	// Set the internal channel
	channelInstance := instance.(*ClassInstance)
	channelInstance.Fields["_channel"] = ch

	return channelInstance, nil
}

// evalSelectStmt evaluates a select statement
func evalSelectStmt(env *Env, stmt *ast.SelectStmt) (val any, returned bool, err error) {
	if len(stmt.Cases) == 0 {
		return nil, false, nil
	}

	// Build list of channels and their operations
	cases := make([]reflect.SelectCase, 0, len(stmt.Cases))
	caseInfo := make([]selectCaseInfo, 0, len(stmt.Cases))

	// Process all cases and check for closed channel cases
	closedCaseIdx := -1
	var closedCaseBody []ast.Stmt
	
	for i, c := range stmt.Cases {
		var ch *common.Channel
		
		if c.IsRecv {
			// For receive case, the Channel expression is typically ch.recv()
			// We need to extract the channel object without actually calling recv()
			// Check if it's a method call expression
			if callExpr, ok := c.Channel.(*ast.CallExpr); ok {
				// It's a call expression, evaluate the callee to get the channel object
				if fieldExpr, ok := callExpr.Callee.(*ast.FieldExpr); ok {
					// It's ch.recv(), evaluate ch to get the channel object
					channelVal, err := evalExpr(env, fieldExpr.X)
					if err != nil {
						return nil, false, err
					}
					
					// Extract the actual channel from ClassInstance
					if channelInstance, ok := channelVal.(*ClassInstance); ok {
						if channelInstance.ClassName == "Channel" {
							if chVal, exists := channelInstance.Fields["_channel"]; exists {
								if channel, ok := chVal.(*common.Channel); ok {
									ch = channel
								}
							}
						}
					}
				}
			}
			
			if ch == nil {
				return nil, false, ThrowRuntimeError(env, "receive case must use ch.recv() pattern")
			}
		} else {
			// For closed case, evaluate the channel expression normally
			channelVal, err := evalExpr(env, c.Channel)
			if err != nil {
				return nil, false, err
			}

			// Extract the actual channel from ClassInstance
			if channelInstance, ok := channelVal.(*ClassInstance); ok {
				if channelInstance.ClassName == "Channel" {
					if chVal, exists := channelInstance.Fields["_channel"]; exists {
						if channel, ok := chVal.(*common.Channel); ok {
							ch = channel
						}
					}
				}
			}

			if ch == nil {
				return nil, false, ThrowRuntimeError(env, fmt.Sprintf("closed case must use a channel, got %T", channelVal))
			}
		}

		if c.IsRecv {
			// Receive case - try to recv and check if channel is closed
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch.Ch),
			})
			caseInfo = append(caseInfo, selectCaseInfo{
				isRecv:  true,
				recvVar: c.RecvVar,
				body:    c.Body,
			})
		} else {
			// Closed case - remember its index and body for later
			closedCaseIdx = i
			closedCaseBody = c.Body
		}
	}

	// If there are no receive cases, nothing to select on
	if len(cases) == 0 {
		return nil, false, nil
	}

	// Perform select operation
	chosen, recv, recvOK := reflect.Select(cases)
	
	info := caseInfo[chosen]
	
	// If this is a receive case, check if channel was closed
	if info.isRecv {
		if !recvOK {
			// Channel was closed, execute closed case if present
			if closedCaseIdx >= 0 {
				// Use runBlock to properly handle break/continue/return
				brk, cont, ret, val, err := runBlock(env, closedCaseBody)
				if err != nil {
					return nil, false, err
				}
				// Propagate break/continue up (select is in a loop context)
				if brk {
					// Return break sentinel so outer loop can handle it
					return breakSentinel{}, false, nil
				}
				if cont {
					return continueSentinel{}, false, nil
				}
				if ret {
					return val, true, nil
				}
			}
			return nil, false, nil
		}
		
		// Bind the received value if variable name provided
		if info.recvVar != "" {
			env.Set(info.recvVar, recv.Interface())
		}
	}

	// Execute the chosen case body using runBlock for proper control flow
	brk, cont, ret, val, err := runBlock(env, info.body)
	if err != nil {
		return nil, false, err
	}
	// Propagate break/continue up
	if brk {
		return breakSentinel{}, false, nil
	}
	if cont {
		return continueSentinel{}, false, nil
	}
	if ret {
		return val, true, nil
	}

	return nil, false, nil
}

type selectCaseInfo struct {
	isRecv  bool
	recvVar string
	body    []ast.Stmt
}

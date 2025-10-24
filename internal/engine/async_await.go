package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
)

// Promise represents a JavaScript-style promise
type Promise struct {
	mu              sync.Mutex
	state           string // "pending", "fulfilled", "rejected"
	value           any
	err             error
	done            chan struct{}
	thenHandlers    []func(any) (any, error)
	catchHandlers   []func(error) (any, error)
	finallyHandlers []func()
}

// CompletableFuture represents a Java-style completable future
type CompletableFuture struct {
	mu        sync.Mutex
	done      chan struct{}
	value     any
	err       error
	completed bool
	cancelled bool
}

// InstallAsyncAwait sets up async/await support in the environment using ClassBuilder
func InstallAsyncAwait(env *common.Env) {
	envTyped := (*Env)(env)

	// Get common type references
	promiseTypeRef := &ast.Type{Name: "Promise", IsBuiltin: true}
	futureTypeRef := &ast.Type{Name: "CompletableFuture", IsBuiltin: true}

	// Pre-register the class definitions so they're available in constructors
	// Build Promise class first, then store definition
	promiseClass := NewClassBuilder("Promise").
		AddTypeParameters(common.TBound.AsGenericType().AsArray()).
		AddField("_promise", promiseTypeRef, []string{"private"})

	promiseDef, _ := buildPromiseClass(promiseClass, envTyped)
	envTyped.Set("__PromiseClass__", promiseDef)

	// Build CompletableFuture class, then store definition
	futureClass := NewClassBuilder("CompletableFuture").
		AddTypeParameters(common.TBound.AsGenericType().AsArray()).
		AddField("_future", futureTypeRef, []string{"private"})

	futureDef, _ := buildCompletableFutureClass(futureClass, envTyped)
	envTyped.Set("__CompletableFutureClass__", futureDef)

	// Install async helper function
	installAsyncFunction(env)
}

// buildPromiseClass builds the Promise class with all its methods
func buildPromiseClass(promiseClass *ClassBuilder, env *Env) (*ClassDefinition, error) {

	// Default constructor: Promise() - for internal use
	promiseClass.AddBuiltinConstructor(
		[]ast.Parameter{},
		func(callEnv *common.Env, args []any) (any, error) {
			// This constructor is not meant to be called directly
			// It's used when creating Promise instances internally
			return nil, ThrowRuntimeError((*Env)(callEnv), "Promise() cannot be called directly. Use Promise(executor) or async()")
		},
	)

	// Main constructor: Promise(executor: Function)
	promiseClass.AddBuiltinConstructor(
		[]ast.Parameter{{Name: "executor", Type: nil}},
		func(callEnv *common.Env, args []any) (any, error) {
			if len(args) != 1 {
				return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
			}

			executor, ok := common.ExtractFunc(args[0])
			if !ok {
				return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
			}

			promise := &Promise{
				state:           "pending",
				thenHandlers:    []func(any) (any, error){},
				catchHandlers:   []func(error) (any, error){},
				finallyHandlers: []func(){},
				done:            make(chan struct{}),
			}

			// Create resolve and reject functions
			resolve := common.Func(func(env *common.Env, args []any) (any, error) {
				if len(args) > 0 {
					promise.resolve(args[0])
				} else {
					promise.resolve(nil)
				}
				return nil, nil
			})

			reject := common.Func(func(env *common.Env, args []any) (any, error) {
				if len(args) > 0 {
					if err, ok := args[0].(error); ok {
						promise.reject(err)
					} else {
						promise.reject(ThrowRuntimeError((*Env)(env), fmt.Sprintf("%v", args[0])))
					}
				} else {
					promise.reject(ThrowRuntimeError((*Env)(env), "promise rejected"))
				}
				return nil, nil
			})

			// Execute the executor function asynchronously
			go func() {
				defer func() {
					if r := recover(); r != nil {
						promise.reject(ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("panic in promise executor: %v", r)))
					}
				}()

				_, err := executor(callEnv, []any{resolve, reject})
				if err != nil {
					promise.reject(err)
				}
			}()

			// Get 'this' instance from constructor environment
			thisVal, ok := callEnv.Get("this")
			if !ok {
				return nil, ThrowRuntimeError((*Env)(callEnv), "constructor called without 'this'")
			}
			instance := thisVal.(*ClassInstance)

			// Set the _promise field on the existing instance
			instance.Fields["_promise"] = promise

			return nil, nil
		},
	)

	// then(onFulfilled: Function) -> Promise<U>
	promiseClass.AddBuiltinMethod("then", &ast.Type{Name: "Promise", IsBuiltin: true}, []ast.Parameter{
		{Name: "onFulfilled", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) < 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		handler, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		// Create a new promise for chaining
		newPromise := &Promise{
			state:           "pending",
			thenHandlers:    []func(any) (any, error){},
			catchHandlers:   []func(error) (any, error){},
			finallyHandlers: []func(){},
			done:            make(chan struct{}),
		}

		promise.mu.Lock()
		if promise.state == "fulfilled" {
			// Original promise already fulfilled
			promise.mu.Unlock()
			go func() {
				result, err := handler(callEnv, []any{promise.value})
				if err != nil {
					newPromise.reject(err)
				} else {
					newPromise.resolve(result)
				}
			}()
		} else if promise.state == "pending" {
			// Original promise still pending
			promise.thenHandlers = append(promise.thenHandlers, func(val any) (any, error) {
				result, err := handler(callEnv, []any{val})
				if err != nil {
					newPromise.reject(err)
					return nil, err
				}
				newPromise.resolve(result)
				return result, nil
			})
			promise.mu.Unlock()
		} else {
			// Original promise rejected
			promise.mu.Unlock()
			newPromise.reject(promise.err)
		}

		// Return new Promise instance
		promiseClassVal, ok := (*Env)(callEnv).Get("__PromiseClass__")
		if !ok {
			return nil, ThrowInitializationError((*Env)(callEnv), "Promise class")
		}
		promiseClassDef := promiseClassVal.(*ClassDefinition)

		newInstance, err := createClassInstanceDirect(promiseClassDef, (*Env)(callEnv))
		if err != nil {
			return nil, err
		}

		newClassInstance := newInstance.(*ClassInstance)
		newClassInstance.Fields["_promise"] = newPromise

		return newClassInstance, nil
	}, []string{})

	// catch(onRejected: Function) -> Promise<T>
	promiseClass.AddBuiltinMethod("catch", &ast.Type{Name: "Promise", IsBuiltin: true}, []ast.Parameter{
		{Name: "onRejected", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) < 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		handler, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		promise.mu.Lock()
		if promise.state == "rejected" {
			promise.mu.Unlock()
			handler(callEnv, []any{promise.err.Error()})
		} else if promise.state == "pending" {
			promise.catchHandlers = append(promise.catchHandlers, func(e error) (any, error) {
				return handler(callEnv, []any{e.Error()})
			})
			promise.mu.Unlock()
		} else {
			promise.mu.Unlock()
		}

		return instance, nil
	}, []string{})

	// finally(onFinally: Function) -> Promise<T>
	promiseClass.AddBuiltinMethod("finally", &ast.Type{Name: "Promise", IsBuiltin: true}, []ast.Parameter{
		{Name: "onFinally", Type: nil},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) < 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		handler, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		promise.mu.Lock()
		if promise.state != "pending" {
			promise.mu.Unlock()
			handler(callEnv, []any{})
		} else {
			promise.finallyHandlers = append(promise.finallyHandlers, func() {
				handler(callEnv, []any{})
			})
			promise.mu.Unlock()
		}

		return instance, nil
	}, []string{})

	// await() -> T
	promiseClass.AddBuiltinMethod("await", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 0 {
			return nil, ThrowArityError((*Env)(callEnv), 0, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		<-promise.done

		promise.mu.Lock()
		defer promise.mu.Unlock()

		if promise.state == "rejected" {
			return nil, promise.err
		}

		return promise.value, nil
	}, []string{})

	// getState() -> String
	promiseClass.AddBuiltinMethod("getState", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		promise.mu.Lock()
		defer promise.mu.Unlock()
		return promise.state, nil
	}, []string{})

	// toString() -> String
	promiseClass.AddBuiltinMethod("toString", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		promise := instance.Fields["_promise"].(*Promise)

		promise.mu.Lock()
		defer promise.mu.Unlock()
		return fmt.Sprintf("Promise{state=%s}", promise.state), nil
	}, []string{})

	_, err := promiseClass.Build(env)
	if err != nil {
		return nil, err
	}

	// Get the class definition from the registered constructor
	promiseClassVal, _ := env.Get("Promise")
	if promiseConstructor, ok := promiseClassVal.(*common.ClassConstructor); ok {
		return promiseConstructor.Definition, nil
	}

	return nil, fmt.Errorf("failed to get Promise class definition")
}

// buildCompletableFutureClass builds the CompletableFuture class with all its methods
func buildCompletableFutureClass(futureClass *ClassBuilder, env *Env) (*ClassDefinition, error) {

	// Constructor: CompletableFuture()
	futureClass.AddBuiltinConstructor(
		[]ast.Parameter{},
		func(callEnv *common.Env, args []any) (any, error) {
			if len(args) != 0 {
				return nil, ThrowArityError((*Env)(callEnv), 0, len(args))
			}

			future := &CompletableFuture{
				done: make(chan struct{}),
			}

			// Get 'this' instance from constructor environment
			thisVal, ok := callEnv.Get("this")
			if !ok {
				return nil, ThrowRuntimeError((*Env)(callEnv), "constructor called without 'this'")
			}
			instance := thisVal.(*ClassInstance)

			// Set the _future field on the existing instance
			instance.Fields["_future"] = future

			return nil, nil
		},
	)

	// complete(value: T) -> Bool
	futureClass.AddBuiltinMethod("complete", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "value", Type: ast.TypeFromString("Any")},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()

		if future.completed || future.cancelled {
			return false, nil
		}

		future.value = args[0]
		future.completed = true
		close(future.done)

		return true, nil
	}, []string{})

	// completeExceptionally(error: Any) -> Bool
	futureClass.AddBuiltinMethod("completeExceptionally", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
		{Name: "error", Type: ast.TypeFromString("Any")},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()

		if future.completed || future.cancelled {
			return false, nil
		}

		if err, ok := args[0].(error); ok {
			future.err = err
		} else {
			future.err = ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("%v", args[0]))
		}
		future.completed = true
		close(future.done)

		return true, nil
	}, []string{})

	// get() -> T
	futureClass.AddBuiltinMethod("get", ast.ANY, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 0 {
			return nil, ThrowArityError((*Env)(callEnv), 0, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		<-future.done

		future.mu.Lock()
		defer future.mu.Unlock()

		if future.err != nil {
			return nil, future.err
		}

		return future.value, nil
	}, []string{})

	// getTimeout(timeout: Int) -> T
	futureClass.AddBuiltinMethod("getTimeout", ast.ANY, []ast.Parameter{
		{Name: "timeout", Type: ast.TypeFromString("Int")},
	}, func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		timeout, ok := args[0].(int)
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		select {
		case <-future.done:
			future.mu.Lock()
			defer future.mu.Unlock()

			if future.err != nil {
				return nil, future.err
			}

			return future.value, nil
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			return nil, ThrowRuntimeError((*Env)(callEnv), "timeout waiting for CompletableFuture")
		}
	}, []string{})

	// isDone() -> Bool
	futureClass.AddBuiltinMethod("isDone", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()
		return future.completed, nil
	}, []string{})

	// isCancelled() -> Bool
	futureClass.AddBuiltinMethod("isCancelled", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()
		return future.cancelled, nil
	}, []string{})

	// cancel() -> Bool
	futureClass.AddBuiltinMethod("cancel", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()

		if future.completed {
			return false, nil
		}

		if !future.cancelled {
			future.cancelled = true
			close(future.done)
		}

		return true, nil
	}, []string{})

	// toString() -> String
	futureClass.AddBuiltinMethod("toString", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.Get("this")
		instance := thisVal.(*ClassInstance)
		future := instance.Fields["_future"].(*CompletableFuture)

		future.mu.Lock()
		defer future.mu.Unlock()

		status := "pending"
		if future.completed {
			status = "completed"
		} else if future.cancelled {
			status = "cancelled"
		}

		return fmt.Sprintf("CompletableFuture{status=%s}", status), nil
	}, []string{})

	_, err := futureClass.Build(env)
	if err != nil {
		return nil, err
	}

	// Get the class definition from the registered constructor
	futureClassVal, _ := env.Get("CompletableFuture")
	if futureConstructor, ok := futureClassVal.(*common.ClassConstructor); ok {
		return futureConstructor.Definition, nil
	}

	return nil, fmt.Errorf("failed to get CompletableFuture class definition")
}

// installAsyncFunction installs the async helper function
func installAsyncFunction(env *common.Env) {
	env.Set("async", common.Func(func(callEnv *common.Env, args []any) (any, error) {
		if len(args) != 1 {
			return nil, ThrowArityError((*Env)(callEnv), 1, len(args))
		}

		fn, ok := common.ExtractFunc(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "function", args[0])
		}

		promise := &Promise{
			state:           "pending",
			thenHandlers:    []func(any) (any, error){},
			catchHandlers:   []func(error) (any, error){},
			finallyHandlers: []func(){},
			done:            make(chan struct{}),
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					promise.reject(ThrowRuntimeError((*Env)(callEnv), fmt.Sprintf("panic in async function: %v", r)))
				}
			}()

			result, err := fn(callEnv, []any{})
			if err != nil {
				promise.reject(err)
			} else {
				promise.resolve(result)
			}
		}()

		// Get the Promise class definition
		promiseClassVal, ok := (*Env)(callEnv).Get("__PromiseClass__")
		if !ok {
			return nil, ThrowInitializationError((*Env)(callEnv), "Promise class")
		}
		promiseClassDef := promiseClassVal.(*ClassDefinition)

		// Create instance directly without calling constructor
		instance, err := createClassInstanceDirect(promiseClassDef, (*Env)(callEnv))
		if err != nil {
			return nil, err
		}

		// Set the _promise field
		classInstance := instance.(*ClassInstance)
		classInstance.Fields["_promise"] = promise

		return classInstance, nil
	}))
}

// createClassInstanceDirect creates a class instance without calling the constructor
// This is used for Promise (async function and then chaining) where we need to set up the internal state
func createClassInstanceDirect(classDef *ClassDefinition, env *Env) (any, error) {
	// Create instance
	instance := &ClassInstance{
		ClassName:   classDef.Name,
		Fields:      make(map[string]any),
		Methods:     make(map[string]common.Func),
		ParentClass: classDef,
	}

	// Initialize fields from class hierarchy
	if err := initializeFields(instance, classDef); err != nil {
		return nil, err
	}

	// Bind methods
	if err := bindMethods(instance, classDef, env); err != nil {
		return nil, err
	}

	return instance, nil
}

func (p *Promise) resolve(value any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != "pending" {
		return
	}

	p.state = "fulfilled"
	p.value = value

	// Execute then handlers
	for _, handler := range p.thenHandlers {
		handler(value)
	}

	// Execute finally handlers
	for _, handler := range p.finallyHandlers {
		handler()
	}

	close(p.done)
}

func (p *Promise) reject(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != "pending" {
		return
	}

	p.state = "rejected"
	p.err = err

	// Execute catch handlers
	for _, handler := range p.catchHandlers {
		handler(err)
	}

	// Execute finally handlers
	for _, handler := range p.finallyHandlers {
		handler()
	}

	close(p.done)
}

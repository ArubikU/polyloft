package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
	"github.com/gorilla/websocket"
)

// InstallHttpModule installs the HTTP module using the builder pattern
func InstallHttpModule(env *Env, opts Options) {
	// Get type references from already-installed builtin types
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	mapType := common.BuiltinTypeMap.GetTypeDefinition(env)
	voidType := &ast.Type{Name: "void", IsBuiltin: true}

	// Step 1: Create HttpRequest builder and get its type BEFORE building
	httpRequestBuilder := NewClassBuilder("HttpRequest").
		AddField("method", stringType, []string{"public"}).
		AddField("path", stringType, []string{"public"}).
		AddField("url", stringType, []string{"public"}).
		AddField("headers", mapType, []string{"public"}).
		AddField("query", mapType, []string{"public"}).
		AddField("params", mapType, []string{"public"}).  // 3.7: Route parameters
		AddField("body", ast.ANY, []string{"public"})

	// Step 2: Create HttpResponse builder and get its type BEFORE building
	httpResponseBuilder := NewClassBuilder("HttpResponse").
		AddField("_writer", ast.ANY, []string{"private"}).
		AddField("_statusCode", intType, []string{"private"}).
		AddField("_headers", mapType, []string{"private"}).
		AddField("_sent", boolType, []string{"private"})

	// Get HttpResponse type reference BEFORE adding methods
	httpResponseType := httpResponseBuilder.GetType()

	// Now add methods using the type reference
	httpResponseBuilder.
		AddBuiltinMethod("status", httpResponseType, []ast.Parameter{
			{Name: "code", Type: intType},
		}, common.Func(httpResponseStatus), []string{}).
		AddBuiltinMethod("header", httpResponseType, []ast.Parameter{
			{Name: "name", Type: stringType},
			{Name: "value", Type: stringType},
		}, common.Func(httpResponseHeader), []string{}).
		AddBuiltinMethod("json", voidType, []ast.Parameter{
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpResponseJson), []string{}).
		AddBuiltinMethod("send", voidType, []ast.Parameter{
			{Name: "text", Type: stringType},
		}, common.Func(httpResponseSend), []string{}).
		AddBuiltinMethod("html", voidType, []ast.Parameter{
			{Name: "html", Type: stringType},
		}, common.Func(httpResponseHtml), []string{}).
		// Response shortcuts (3.8)
		AddBuiltinMethod("ok", voidType, []ast.Parameter{
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpResponseOk), []string{}).
		AddBuiltinMethod("created", voidType, []ast.Parameter{
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpResponseCreated), []string{}).
		AddBuiltinMethod("noContent", voidType, []ast.Parameter{}, common.Func(httpResponseNoContent), []string{}).
		AddBuiltinMethod("notFound", voidType, []ast.Parameter{
			{Name: "message", Type: stringType},
		}, common.Func(httpResponseNotFound), []string{}).
		AddBuiltinMethod("notFound", voidType, []ast.Parameter{}, common.Func(httpResponseNotFound), []string{}).
		AddBuiltinMethod("error", voidType, []ast.Parameter{
			{Name: "code", Type: intType},
			{Name: "message", Type: stringType},
		}, common.Func(httpResponseError), []string{}).
		// Template rendering - 3.12
		AddBuiltinMethod("render", voidType, []ast.Parameter{
			{Name: "template", Type: stringType},
			{Name: "data", Type: mapType},
		}, common.Func(httpResponseRender), []string{})

	// Step 3: Create HttpServer builder and get its type BEFORE building
	httpServerBuilder := NewClassBuilder("HttpServer").
		AddField("router", ast.ANY, []string{"private"}).
		AddField("_config", mapType, []string{"private"}).
		AddField("_errorHandler", ast.ANY, []string{"private"}).
		AddField("_globalMiddlewares", ast.ANY, []string{"private"}).
		AddField("_logLevel", stringType, []string{"private"}).
		SetBuiltinConstructor([]ast.Parameter{}, common.Func(newHttpServer)).
		AddBuiltinMethod("get", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerGet), []string{}).
		AddBuiltinMethod("get", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "middlewares", Type: ast.ANY},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerGet), []string{}).
		AddBuiltinMethod("post", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPost), []string{}).
		AddBuiltinMethod("post", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "middlewares", Type: ast.ANY},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPost), []string{}).
		AddBuiltinMethod("put", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPut), []string{}).
		AddBuiltinMethod("put", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "middlewares", Type: ast.ANY},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPut), []string{}).
		AddBuiltinMethod("delete", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerDelete), []string{}).
		AddBuiltinMethod("delete", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "middlewares", Type: ast.ANY},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerDelete), []string{}).
		AddBuiltinMethod("use", voidType, []ast.Parameter{
			{Name: "middleware", Type: ast.ANY},
		}, common.Func(httpServerUse), []string{}).
		AddBuiltinMethod("onError", voidType, []ast.Parameter{
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerOnError), []string{}).
		AddBuiltinMethod("config", voidType, []ast.Parameter{
			{Name: "options", Type: mapType},
		}, common.Func(httpServerConfig), []string{}).
		AddBuiltinMethod("log", voidType, []ast.Parameter{
			{Name: "message", Type: stringType},
			{Name: "level", Type: stringType},
		}, common.Func(httpServerLog), []string{}).
		AddBuiltinMethod("log", voidType, []ast.Parameter{
			{Name: "message", Type: stringType},
		}, common.Func(httpServerLog), []string{}).
		// WebSocket support - 3.2
		AddBuiltinMethod("ws", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerWs), []string{}).
		AddBuiltinMethod("listen", mapType, []ast.Parameter{
			{Name: "port", Type: stringType},
		}, common.Func(httpServerListen), []string{})

	httpServerType := httpServerBuilder.GetType()

	// Get Promise type for async methods
	promiseType := common.BuiltinTypePromise.GetTypeDefinition(env)

	// Step 5: Create Http class with static methods using proper type references
	httpStaticClassBuilder := NewClassBuilder("Http").
		// Existing synchronous methods
		AddStaticMethod("get", mapType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "timeout", Type: intType},
		}, common.Func(httpGet)).
		AddStaticMethod("post", mapType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
		}, common.Func(httpPost)).
		AddStaticMethod("put", mapType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
		}, common.Func(httpPut)).
		AddStaticMethod("delete", mapType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "timeout", Type: intType},
		}, common.Func(httpDelete)).
		// Simplified request - 3.3: Allow body without options wrapper
		AddStaticMethod("request", mapType, []ast.Parameter{
			{Name: "method", Type: stringType},
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpRequest)).
		AddStaticMethod("request", mapType, []ast.Parameter{
			{Name: "method", Type: stringType},
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
		}, common.Func(httpRequest)).
		AddStaticMethod("request", mapType, []ast.Parameter{
			{Name: "method", Type: stringType},
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
			{Name: "headers", Type: mapType},
		}, common.Func(httpRequest)).
		// Async methods returning Promises - 3.10
		AddStaticMethod("getAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
		}, common.Func(httpGetAsync)).
		AddStaticMethod("getAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "timeout", Type: intType},
		}, common.Func(httpGetAsync)).
		AddStaticMethod("postAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpPostAsync)).
		AddStaticMethod("postAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
		}, common.Func(httpPostAsync)).
		AddStaticMethod("putAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpPutAsync)).
		AddStaticMethod("putAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
		}, common.Func(httpPutAsync)).
		AddStaticMethod("deleteAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
		}, common.Func(httpDeleteAsync)).
		AddStaticMethod("deleteAsync", promiseType, []ast.Parameter{
			{Name: "url", Type: stringType},
			{Name: "timeout", Type: intType},
		}, common.Func(httpDeleteAsync)).
		AddStaticMethod("requestAsync", promiseType, []ast.Parameter{
			{Name: "method", Type: stringType},
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
		}, common.Func(httpRequestAsync)).
		AddStaticMethod("createServer", httpServerType, []ast.Parameter{
			{Name: "debug", Type: boolType},
		}, common.Func(createHttpServer)).
		AddStaticMethod("createServer", httpServerType, []ast.Parameter{}, common.Func(createHttpServer))

	// Step 4: NOW build all classes after getting their type references
	_, _ = httpRequestBuilder.Build(env)
	_, _ = httpResponseBuilder.Build(env)
	_, _ = httpServerBuilder.Build(env)
	_, _ = httpStaticClassBuilder.BuildStatic(env)
}

// httpGet performs an HTTP GET request
func httpGet(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])
	timeout := 30 * time.Second
	if len(args) > 1 {
		if t, ok := utils.AsInt(args[1]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return createHttpResponse((*Env)(e), resp, body), nil
}

// httpPost performs an HTTP POST request
func httpPost(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])

	bodyBytes, err := prepareRequestBody(args[1])
	if err != nil {
		return nil, err
	}

	timeout := 30 * time.Second
	if len(args) > 2 {
		if t, ok := utils.AsInt(args[2]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return createHttpResponse((*Env)(e), resp, body), nil
}

// httpPut performs an HTTP PUT request
func httpPut(e *common.Env, args []any) (any, error) {
	if len(args) < 2 {
		return nil, ThrowArityError((*Env)(e), 2, len(args))
	}
	url := utils.ToString(args[0])

	bodyBytes, err := prepareRequestBody(args[1])
	if err != nil {
		return nil, err
	}

	timeout := 30 * time.Second
	if len(args) > 2 {
		if t, ok := utils.AsInt(args[2]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return createHttpResponse((*Env)(e), resp, body), nil
}

// httpDelete performs an HTTP DELETE request
func httpDelete(e *common.Env, args []any) (any, error) {
	if len(args) < 1 {
		return nil, ThrowArityError((*Env)(e), 1, len(args))
	}
	url := utils.ToString(args[0])

	timeout := 30 * time.Second
	if len(args) > 1 {
		if t, ok := utils.AsInt(args[1]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return createHttpResponse((*Env)(e), resp, body), nil
}

// httpRequest performs a custom HTTP request
func httpRequest(e *common.Env, args []any) (any, error) {
	if len(args) < 2 {
		return nil, ThrowArityError((*Env)(e), 2, len(args))
	}
	method := utils.ToString(args[0])
	url := utils.ToString(args[1])

	var bodyBytes []byte
	if len(args) > 2 && args[2] != nil {
		var err error
		bodyBytes, err = prepareRequestBody(args[2])
		if err != nil {
			return nil, err
		}
	}

	timeout := 30 * time.Second
	if len(args) > 3 {
		if t, ok := utils.AsInt(args[3]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	client := &http.Client{Timeout: timeout}
	var req *http.Request
	var err error
	if len(bodyBytes) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if provided
	if len(args) > 4 {
		if headers, ok := args[4].(map[string]any); ok {
			for key, value := range headers {
				req.Header.Set(key, utils.ToString(value))
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return createHttpResponse((*Env)(e), resp, body), nil
}

// createHttpServer creates a new HTTP server instance
func createHttpServer(e *common.Env, args []any) (any, error) {
	// Create a new HttpServer instance
	ctor := common.BuiltinTypeHttpServer.GetConstructor(e)
	if ctor == nil {
		return nil, ThrowInitializationError((*Env)(e), "HttpServer class")
	}

	return ctor.Func(e, []any{})
}

// newHttpServer creates a new HttpServer instance
func newHttpServer(e *common.Env, args []any) (any, error) {
	// Get the instance from the environment (created by createClassInstance)
	thisVal, exists := e.This()
	if !exists {
		return nil, ThrowRuntimeError((*Env)(e), "no instance context found")
	}

	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "ClassInstance", thisVal)
	}

	// Initialize the router field
	router := &httpRouter{
		routes:             make(map[string]map[string]*routeHandler),
		dynamicRoutes:      make(map[string][]*routeHandler),
		mu:                 &sync.RWMutex{},
		globalMiddlewares:  []common.Func{},
		errorHandler:       nil,
		config:             make(map[string]any),
		logLevel:           "info",
		wsHandlers:         make(map[string]common.Func),
	}

	instance.Fields["router"] = router
	instance.Fields["_config"] = make(map[string]any)
	instance.Fields["_errorHandler"] = nil
	instance.Fields["_globalMiddlewares"] = []common.Func{}
	instance.Fields["_logLevel"] = "info"

	return nil, nil // Constructors shouldn't return the instance
}

// httpServerGet registers a GET route handler
func httpServerGet(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}
	router := instance.Fields["router"].(*httpRouter)
	
	path := utils.ToString(args[0])
	var middlewares []common.Func
	var handler common.Func
	
	// Support both forms: (path, handler) and (path, middlewares, handler)
	if len(args) == 2 {
		h, ok := common.ExtractFunc(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[1])
		}
		handler = h
	} else if len(args) == 3 {
		// Extract middlewares
		middlewares = extractMiddlewares(args[1])
		h, ok := common.ExtractFunc(args[2])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[2])
		}
		handler = h
	}

	router.addRoute("GET", path, handler, middlewares)
	return nil, nil
}

// httpServerPost registers a POST route handler
func httpServerPost(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	path := utils.ToString(args[0])
	var middlewares []common.Func
	var handler common.Func
	
	// Support both forms: (path, handler) and (path, middlewares, handler)
	if len(args) == 2 {
		h, ok := common.ExtractFunc(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[1])
		}
		handler = h
	} else if len(args) == 3 {
		// Extract middlewares
		middlewares = extractMiddlewares(args[1])
		h, ok := common.ExtractFunc(args[2])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[2])
		}
		handler = h
	}

	router.addRoute("POST", path, handler, middlewares)
	return nil, nil
}

// httpServerPut registers a PUT route handler
func httpServerPut(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	path := utils.ToString(args[0])
	var middlewares []common.Func
	var handler common.Func
	
	// Support both forms: (path, handler) and (path, middlewares, handler)
	if len(args) == 2 {
		h, ok := common.ExtractFunc(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[1])
		}
		handler = h
	} else if len(args) == 3 {
		// Extract middlewares
		middlewares = extractMiddlewares(args[1])
		h, ok := common.ExtractFunc(args[2])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[2])
		}
		handler = h
	}

	router.addRoute("PUT", path, handler, middlewares)
	return nil, nil
}

// httpServerDelete registers a DELETE route handler
func httpServerDelete(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	path := utils.ToString(args[0])
	var middlewares []common.Func
	var handler common.Func
	
	// Support both forms: (path, handler) and (path, middlewares, handler)
	if len(args) == 2 {
		h, ok := common.ExtractFunc(args[1])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[1])
		}
		handler = h
	} else if len(args) == 3 {
		// Extract middlewares
		middlewares = extractMiddlewares(args[1])
		h, ok := common.ExtractFunc(args[2])
		if !ok {
			return nil, ThrowTypeError((*Env)(e), "function", args[2])
		}
		handler = h
	}

	router.addRoute("DELETE", path, handler, middlewares)
	return nil, nil
}

// httpServerListen starts the HTTP server
func httpServerListen(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	port := utils.ToString(args[0])
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	// Create HTTP handler with timeout support - 3.13
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a WebSocket request
		if wsHandler, isWs := router.isWebSocketRequest(r.URL.Path); isWs {
			router.handleWebSocket((*common.Env)(e), w, r, wsHandler)
			return
		}

		// Check for timeout in config
		timeoutMs := 0
		if timeout, ok := router.config["timeout"]; ok {
			if t, ok := utils.AsInt(timeout); ok {
				timeoutMs = t
			}
		}

		if timeoutMs > 0 {
			// Use context with timeout
			ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeoutMs)*time.Millisecond)
			defer cancel()
			
			// Replace request context
			r = r.WithContext(ctx)
		}
		
		router.handleRequest(e, w, r)
	})

	// Start server in background
	go func() {
		fmt.Printf("HTTP Server listening on %s\n", port)
		if err := http.ListenAndServe(port, httpHandler); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return map[string]any{
		"address": port,
		"message": "Server started successfully",
	}, nil
}

// httpResponseStatus sets the HTTP status code
func httpResponseStatus(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	code, ok := utils.AsInt(args[0])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "number", args[0])
	}

	instance.Fields["_statusCode"] = code
	return instance, nil
}

// httpResponseHeader sets an HTTP header
func httpResponseHeader(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	name := utils.ToString(args[0])
	value := utils.ToString(args[1])

	headers := instance.Fields["_headers"].(map[string]string)
	headers[name] = value

	return instance, nil
}

// httpResponseJson sends a JSON response
func httpResponseJson(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	resp := instance.Fields["_writer"].(*httpResponse)
	resp.sendJSON(args[0])

	return nil, nil
}

// httpResponseSend sends a text response
func httpResponseSend(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	resp := instance.Fields["_writer"].(*httpResponse)
	resp.send(utils.ToString(args[0]))

	return nil, nil
}

// httpResponseHtml sends an HTML response
func httpResponseHtml(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	resp := instance.Fields["_writer"].(*httpResponse)
	resp.sendHTML(utils.ToString(args[0]))

	return nil, nil
}

// Helper functions

// prepareRequestBody converts request data to JSON bytes
func prepareRequestBody(data any) ([]byte, error) {
	// Handle Map instances
	if mapInstance, ok := data.(*ClassInstance); ok && mapInstance.ClassName == "Map" {
		// Note: We can't pass env here since we don't have it, so we'll use a simple extraction
		// This is a limitation but works for most cases
		objMap := make(map[string]any)
		if hashData, ok := mapInstance.Fields["_data"].(map[uint64][]*mapEntry); ok {
			for _, entries := range hashData {
				for _, entry := range entries {
					keyStr := utils.ToString(entry.Key)
					objMap[keyStr] = entry.Value
				}
			}
		}
		return json.Marshal(objMap)
	}
	// Handle plain Go maps
	if dataMap, ok := data.(map[string]any); ok {
		return json.Marshal(dataMap)
	}
	return []byte(utils.ToString(data)), nil
}

// createHttpResponse creates a standardized HTTP response object
func createHttpResponse(env *Env, resp *http.Response, body []byte) any {
	// Parse body based on Content-Type
	var bodyData any = string(body)
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") && len(body) > 0 {
		var jsonData any
		if err := json.Unmarshal(body, &jsonData); err == nil {
			// Convert JSON to Polyloft Map if it's an object
			if jsonMap, ok := jsonData.(map[string]any); ok {
				if mapInstance, err := CreateMapInstance(env, jsonMap); err == nil {
					bodyData = mapInstance
				} else {
					bodyData = jsonData // Fallback to raw Go map
				}
			} else {
				bodyData = jsonData // Arrays, strings, numbers, etc.
			}
		}
	}

	// Create response Map as Polyloft Map instance
	responseMap := map[string]any{
		"status":     float64(resp.StatusCode),
		"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
		"statusText": resp.Status,
		"body":       bodyData,
		"headers":    convertHeaders(resp.Header),
	}
	
	// Convert the response map to a Polyloft Map instance
	if mapInstance, err := CreateMapInstance(env, responseMap); err == nil {
		return mapInstance
	}
	
	// Fallback to plain Go map if conversion fails
	return responseMap
}

// convertHeaders converts http.Header to map[string]any
func convertHeaders(headers http.Header) map[string]any {
	result := make(map[string]any)
	for key, values := range headers {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			valsAny := make([]any, len(values))
			for i, v := range values {
				valsAny[i] = v
			}
			result[key] = valsAny
		}
	}
	return result
}

// routeHandler holds a handler and its middlewares
type routeHandler struct {
	handler     common.Func
	middlewares []common.Func
	pattern     *routePattern // For dynamic routes
}

// routePattern represents a parsed route pattern with parameters
type routePattern struct {
	segments []routeSegment
	isStatic bool
	original string
}

// routeSegment represents a part of the route
type routeSegment struct {
	isParam   bool
	isWildcard bool
	name      string
	value     string
	validator *regexp.Regexp // For validation like :id([0-9]+)
}

// httpRouter manages HTTP routes
type httpRouter struct {
	routes            map[string]map[string]*routeHandler // method -> path -> handler (static routes)
	dynamicRoutes     map[string][]*routeHandler          // method -> []handler (dynamic routes)
	mu                *sync.RWMutex
	globalMiddlewares []common.Func
	errorHandler      common.Func
	config            map[string]any
	logLevel          string
	wsHandlers        map[string]common.Func // WebSocket handlers
}

func (r *httpRouter) addRoute(method, path string, handler common.Func, middlewares []common.Func) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Parse the route pattern to check for dynamic segments
	pattern := parseRoutePattern(path)
	
	rHandler := &routeHandler{
		handler:     handler,
		middlewares: middlewares,
		pattern:     pattern,
	}

	if pattern.isStatic {
		// Static route - use map for O(1) lookup
		if r.routes[method] == nil {
			r.routes[method] = make(map[string]*routeHandler)
		}
		r.routes[method][path] = rHandler
	} else {
		// Dynamic route - add to slice for matching
		if r.dynamicRoutes == nil {
			r.dynamicRoutes = make(map[string][]*routeHandler)
		}
		r.dynamicRoutes[method] = append(r.dynamicRoutes[method], rHandler)
	}
}

// parseRoutePattern parses a route pattern into segments
// Supports: /users/:id, /users/:id([0-9]+), /files/*filepath
func parseRoutePattern(path string) *routePattern {
	pattern := &routePattern{
		original: path,
		isStatic: true,
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		segment := routeSegment{}

		// Check for wildcard: *name
		if strings.HasPrefix(part, "*") {
			segment.isWildcard = true
			segment.name = part[1:]
			pattern.isStatic = false
		} else if strings.HasPrefix(part, ":") {
			// Parameter with optional validation: :id or :id([0-9]+)
			segment.isParam = true
			pattern.isStatic = false

			// Check for validation pattern
			if idx := strings.Index(part, "("); idx > 0 {
				segment.name = part[1:idx]
				validatorStr := part[idx+1 : len(part)-1] // Remove ( and )
				if compiled, err := regexp.Compile("^" + validatorStr + "$"); err == nil {
					segment.validator = compiled
				}
			} else {
				segment.name = part[1:]
			}
		} else {
			// Static segment
			segment.value = part
		}

		pattern.segments = append(pattern.segments, segment)
	}

	return pattern
}

// matchRoute attempts to match a request path against a route pattern
// Returns matched params if successful, nil otherwise
func (rh *routeHandler) matchRoute(reqPath string) map[string]string {
	if rh.pattern.isStatic {
		if reqPath == rh.pattern.original {
			return make(map[string]string)
		}
		return nil
	}

	params := make(map[string]string)
	reqParts := strings.Split(strings.Trim(reqPath, "/"), "/")
	patternSegs := rh.pattern.segments

	// Handle wildcards
	if len(patternSegs) > 0 && patternSegs[len(patternSegs)-1].isWildcard {
		// Wildcard must match all remaining parts
		if len(reqParts) < len(patternSegs) {
			return nil
		}
		
		// Match all segments before wildcard
		for i := 0; i < len(patternSegs)-1; i++ {
			if !matchSegment(patternSegs[i], reqParts[i], params) {
				return nil
			}
		}
		
		// Wildcard captures remaining path
		wildcardSeg := patternSegs[len(patternSegs)-1]
		params[wildcardSeg.name] = strings.Join(reqParts[len(patternSegs)-1:], "/")
		return params
	}

	// Regular matching - must have same number of segments
	if len(reqParts) != len(patternSegs) {
		return nil
	}

	for i, seg := range patternSegs {
		if !matchSegment(seg, reqParts[i], params) {
			return nil
		}
	}

	return params
}

// matchSegment matches a single segment
func matchSegment(seg routeSegment, value string, params map[string]string) bool {
	if seg.isParam {
		// Validate if validator exists
		if seg.validator != nil && !seg.validator.MatchString(value) {
			return false
		}
		params[seg.name] = value
		return true
	}
	
	// Static segment must match exactly
	return seg.value == value
}

func (r *httpRouter) handleRequest(env *common.Env, w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	staticRoutes := r.routes[req.Method]
	dynamicRoutes := r.dynamicRoutes[req.Method]
	errorHandler := r.errorHandler
	globalMiddlewares := r.globalMiddlewares
	r.mu.RUnlock()

	// Try to find matching route (static first, then dynamic)
	var routeHandler *routeHandler
	var routeParams map[string]string

	// 1. Try static route first (O(1) lookup)
	if staticRoutes != nil {
		if handler, found := staticRoutes[req.URL.Path]; found {
			routeHandler = handler
			routeParams = make(map[string]string)
		}
	}

	// 2. Try dynamic routes if no static match
	if routeHandler == nil && dynamicRoutes != nil {
		for _, handler := range dynamicRoutes {
			if params := handler.matchRoute(req.URL.Path); params != nil {
				routeHandler = handler
				routeParams = params
				break
			}
		}
	}

	// 3. No route found
	if routeHandler == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not Found"}`))
		return
	}

	// Parse request body
	var bodyData any
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		if len(bodyBytes) > 0 {
			// Try to parse as JSON
			var jsonData map[string]any
			if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
				bodyData = jsonData
			} else {
				bodyData = string(bodyBytes)
			}
		}
	}

	// Parse query parameters
	queryParams := make(map[string]any)
	for key, values := range req.URL.Query() {
		if len(values) == 1 {
			queryParams[key] = values[0]
		} else {
			valsAny := make([]any, len(values))
			for i, v := range values {
				valsAny[i] = v
			}
			queryParams[key] = valsAny
		}
	}

	// Create HttpRequest instance using the class constructor
	requestCtor := common.BuiltinTypeHttpRequest.GetConstructor(env)
	var requestInstance *ClassInstance
	if requestCtor != nil {
		inst, err := requestCtor.Func(env, []any{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "Failed to create HttpRequest: %v"}`, err)))
			return
		}
		if ri, ok := inst.(*ClassInstance); ok {
			requestInstance = ri
			// Set the fields
			requestInstance.Fields["method"], _ = CreateStringInstance(env, req.Method)
			requestInstance.Fields["path"], _ = CreateStringInstance(env, req.URL.Path)
			requestInstance.Fields["url"], _ = CreateStringInstance(env, req.URL.String())
			requestInstance.Fields["headers"], _ = CreateMapInstance(env, convertHeaders(req.Header))
			requestInstance.Fields["query"], _ = CreateMapInstance(env, queryParams)
			// Convert routeParams (map[string]string) to map[string]any for CreateMapInstance
			routeParamsAny := make(map[string]any)
			for k, v := range routeParams {
				routeParamsAny[k] = v
			}
			requestInstance.Fields["params"], _ = CreateMapInstance(env, routeParamsAny)
			requestInstance.Fields["body"], _ = CreateGenericInstance(env, bodyData)
		}
	}

	if requestInstance == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create HttpRequest instance"}`))
		return
	}

	// Create underlying httpResponse
	responseObj := &httpResponse{
		writer:     w,
		statusCode: 200,
		headers:    make(map[string]string),
		env:        (*Env)(env),
	}

	// Create HttpResponse instance using the class constructor
	responseCtor := common.BuiltinTypeHttpResponse.GetConstructor(env)
	var responseInstance *ClassInstance
	if responseCtor != nil {
		inst, err := responseCtor.Func(env, []any{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "Failed to create HttpResponse: %v"}`, err)))
			return
		}
		if ri, ok := inst.(*ClassInstance); ok {
			responseInstance = ri
			// Set the fields
			responseInstance.Fields["_writer"] = responseObj
			responseInstance.Fields["_statusCode"] = 200
			responseInstance.Fields["_headers"] = make(map[string]string)
			responseInstance.Fields["_sent"] = false
		}
	}

	if responseInstance == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create HttpResponse instance"}`))
		return
	}

	// Execute middleware chain and handler
	middlewareChain := append([]common.Func{}, globalMiddlewares...)
	middlewareChain = append(middlewareChain, routeHandler.middlewares...)
	
	// Create next function for middleware chain
	var executeChain func(int) error
	executeChain = func(index int) error {
		if index < len(middlewareChain) {
			// Create next function to be passed to middleware
			nextFunc := common.Func(func(e *common.Env, args []any) (any, error) {
				return nil, executeChain(index + 1)
			})
			
			// Call middleware with req, res, next
			_, err := middlewareChain[index](env, []any{requestInstance, responseInstance, nextFunc})
			return err
		}
		
		// All middlewares passed, call the actual handler
		_, err := routeHandler.handler(env, []any{requestInstance, responseInstance})
		return err
	}
	
	// Execute the chain with error handling
	err := executeChain(0)
	if err != nil {
		// Use custom error handler if available
		if errorHandler != nil {
			errorHandler(env, []any{err, requestInstance, responseInstance})
		} else {
			// Default error response
			if !responseObj.sent {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
			}
		}
	}
}

// extractMiddlewares extracts middleware functions from an array or single function
// Only checks for _items field, validates that extracted values are valid middleware functions with 3 parameters
func extractMiddlewares(arg any) []common.Func {
	var middlewares []common.Func
	
	// Try to extract from any object with _items field (List or Array-like)
	if instance, ok := arg.(*ClassInstance); ok {
		// Check if it has _items field
		if items, hasItems := instance.Fields["_items"]; hasItems {
			if itemSlice, ok := items.([]any); ok {
				for _, item := range itemSlice {
					// Extract function and validate it's a proper middleware
					if fn, ok := common.ExtractFunc(item); ok {
						// Validate middleware has exactly 3 parameters (req, res, next)
						if isValidMiddleware(item) {
							middlewares = append(middlewares, fn)
						}
					}
				}
			}
		}
	} else if slice, ok := arg.([]any); ok {
		// Handle plain Go slice
		for _, item := range slice {
			if fn, ok := common.ExtractFunc(item); ok {
				// Validate middleware has exactly 3 parameters (req, res, next)
				if isValidMiddleware(item) {
					middlewares = append(middlewares, fn)
				}
			}
		}
	} else {
		// Try single middleware function
		if fn, ok := common.ExtractFunc(arg); ok {
			// Validate middleware has exactly 3 parameters (req, res, next)
			if isValidMiddleware(arg) {
				middlewares = append(middlewares, fn)
			}
		}
	}
	
	return middlewares
}

// isValidMiddleware checks if a function has exactly 3 parameters (req, res, next)
func isValidMiddleware(fn any) bool {
	// Check FunctionDefinition
	if funcDef, ok := fn.(*common.FunctionDefinition); ok {
		return len(funcDef.Params) == 3
	}
	
	// Check LambdaDefinition
	if lambdaDef, ok := fn.(*common.LambdaDefinition); ok {
		return len(lambdaDef.Params) == 3
	}
	
	// If we can't determine parameter count, reject it to be safe
	return false
}

// httpServerUse registers a global middleware
func httpServerUse(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	middleware, ok := common.ExtractFunc(args[0])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[0])
	}

	router.mu.Lock()
	router.globalMiddlewares = append(router.globalMiddlewares, middleware)
	router.mu.Unlock()

	return nil, nil
}

// httpServerOnError registers a global error handler
func httpServerOnError(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	handler, ok := common.ExtractFunc(args[0])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[0])
	}

	router.mu.Lock()
	router.errorHandler = handler
	router.mu.Unlock()

	return nil, nil
}

// httpServerConfig configures the server
func httpServerConfig(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	
	// Extract config map using MapToObject
	var configMap map[string]any
	if mapInstance, ok := args[0].(*ClassInstance); ok && mapInstance.ClassName == "Map" {
		objMap, err := MapToObject((*Env)(e), mapInstance)
		if err != nil {
			return nil, err
		}
		configMap = objMap
	} else if m, ok := args[0].(map[string]any); ok {
		configMap = m
	}

	router.mu.Lock()
	for key, value := range configMap {
		router.config[key] = value
	}
	router.mu.Unlock()

	return nil, nil
}

// httpServerLog logs a message with an optional level
func httpServerLog(e *common.Env, args []any) (any, error) {
	message := utils.ToString(args[0])
	level := "info"
	if len(args) > 1 {
		level = utils.ToString(args[1])
	}

	// Simple logging implementation
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] %s\n", timestamp, strings.ToUpper(level), message)

	return nil, nil
}

// Response shortcut methods (3.8)

// httpResponseOk sends a 200 OK response with data
func httpResponseOk(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	instance.Fields["_statusCode"] = 200
	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = 200
	resp.sendJSON(args[0])

	return nil, nil
}

// httpResponseCreated sends a 201 Created response with data
func httpResponseCreated(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	instance.Fields["_statusCode"] = 201
	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = 201
	resp.sendJSON(args[0])

	return nil, nil
}

// httpResponseNoContent sends a 204 No Content response
func httpResponseNoContent(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = 204
	if !resp.sent {
		resp.sent = true
		for k, v := range resp.headers {
			resp.writer.Header().Set(k, v)
		}
		resp.writer.WriteHeader(204)
	}

	return nil, nil
}

// httpResponseNotFound sends a 404 Not Found response
func httpResponseNotFound(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	message := "Not Found"
	if len(args) > 0 {
		message = utils.ToString(args[0])
	}

	instance.Fields["_statusCode"] = 404
	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = 404
	resp.sendJSON(map[string]any{"error": message})

	return nil, nil
}

// httpResponseError sends a custom error response
func httpResponseError(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	code, ok := utils.AsInt(args[0])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "number", args[0])
	}
	message := utils.ToString(args[1])

	instance.Fields["_statusCode"] = code
	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = code
	resp.sendJSON(map[string]any{"error": message})

	return nil, nil
}

// httpResponse wraps http.ResponseWriter with convenience methods
type httpResponse struct {
	writer     http.ResponseWriter
	statusCode int
	headers    map[string]string
	sent       bool
	env        *Env // Store env for MapToObject
}

func (r *httpResponse) sendJSON(data any) {
	if r.sent {
		return
	}
	r.sent = true

	// Set headers
	r.writer.Header().Set("Content-Type", "application/json")
	for k, v := range r.headers {
		r.writer.Header().Set(k, v)
	}
	r.writer.WriteHeader(r.statusCode)

	// Handle Map instances - convert to Go map for JSON encoding
	if mapInstance, ok := data.(*ClassInstance); ok && mapInstance.ClassName == "Map" {
		if r.env != nil {
			objMap, err := MapToObject(r.env, mapInstance)
			if err != nil {
				json.NewEncoder(r.writer).Encode(map[string]string{"error": "failed to convert Map: " + err.Error()})
				return
			}
			json.NewEncoder(r.writer).Encode(objMap)
			return
		}
	}

	// Handle regular data
	json.NewEncoder(r.writer).Encode(data)
}

func (r *httpResponse) send(text string) {
	if r.sent {
		return
	}
	r.sent = true

	// Set headers
	for k, v := range r.headers {
		r.writer.Header().Set(k, v)
	}
	r.writer.WriteHeader(r.statusCode)
	r.writer.Write([]byte(text))
}

func (r *httpResponse) sendHTML(html string) {
	if r.sent {
		return
	}
	r.sent = true

	// Set headers
	r.writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	for k, v := range r.headers {
		r.writer.Header().Set(k, v)
	}
	r.writer.WriteHeader(r.statusCode)
	r.writer.Write([]byte(html))
}

// Async HTTP Methods - 3.10: Native Async/Await Integration

// httpGetAsync performs an async HTTP GET request returning a Promise
func httpGetAsync(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])
	timeout := 30 * time.Second
	if len(args) > 1 {
		if t, ok := utils.AsInt(args[1]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	return createHttpPromise((*Env)(e), func() (any, error) {
		client := &http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return createHttpResponse((*Env)(e), resp, body), nil
	})
}

// httpPostAsync performs an async HTTP POST request returning a Promise
func httpPostAsync(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])
	
	bodyBytes, err := prepareRequestBody(args[1])
	if err != nil {
		return nil, err
	}

	timeout := 30 * time.Second
	if len(args) > 2 {
		if t, ok := utils.AsInt(args[2]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	return createHttpPromise((*Env)(e), func() (any, error) {
		client := &http.Client{Timeout: timeout}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return createHttpResponse((*Env)(e), resp, body), nil
	})
}

// httpPutAsync performs an async HTTP PUT request returning a Promise
func httpPutAsync(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])
	
	bodyBytes, err := prepareRequestBody(args[1])
	if err != nil {
		return nil, err
	}

	timeout := 30 * time.Second
	if len(args) > 2 {
		if t, ok := utils.AsInt(args[2]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	return createHttpPromise((*Env)(e), func() (any, error) {
		client := &http.Client{Timeout: timeout}
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return createHttpResponse((*Env)(e), resp, body), nil
	})
}

// httpDeleteAsync performs an async HTTP DELETE request returning a Promise
func httpDeleteAsync(e *common.Env, args []any) (any, error) {
	url := utils.ToString(args[0])

	timeout := 30 * time.Second
	if len(args) > 1 {
		if t, ok := utils.AsInt(args[1]); ok {
			timeout = time.Duration(t) * time.Second
		}
	}

	return createHttpPromise((*Env)(e), func() (any, error) {
		client := &http.Client{Timeout: timeout}
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return createHttpResponse((*Env)(e), resp, body), nil
	})
}

// httpRequestAsync performs an async custom HTTP request returning a Promise
func httpRequestAsync(e *common.Env, args []any) (any, error) {
	if len(args) < 2 {
		return nil, ThrowArityError((*Env)(e), 2, len(args))
	}
	method := utils.ToString(args[0])
	url := utils.ToString(args[1])

	var bodyBytes []byte
	if len(args) > 2 && args[2] != nil {
		var err error
		bodyBytes, err = prepareRequestBody(args[2])
		if err != nil {
			return nil, err
		}
	}

	return createHttpPromise((*Env)(e), func() (any, error) {
		timeout := 30 * time.Second
		client := &http.Client{Timeout: timeout}
		
		var req *http.Request
		var err error
		if len(bodyBytes) > 0 {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return createHttpResponse((*Env)(e), resp, body), nil
	})
}

// createHttpPromise creates a Promise that executes an HTTP request asynchronously
func createHttpPromise(env *Env, requestFunc func() (any, error)) (any, error) {
	// Get the Promise class definition
	promiseClassDef := common.BuiltinTypePromise.GetClassDefinition((*common.Env)(env))
	if promiseClassDef == nil {
		return nil, ThrowInitializationError(env, "Promise class")
	}

	// Create Promise instance directly without calling constructor
	instance, err := createClassInstanceDirect(promiseClassDef, env)
	if err != nil {
		return nil, err
	}

	// Create the underlying Promise structure
	promise := &Promise{
		state:           "pending",
		thenHandlers:    []func(any) (any, error){},
		catchHandlers:   []func(error) (any, error){},
		finallyHandlers: []func(){},
		done:            make(chan struct{}),
	}

	// Execute the request asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				promise.reject(ThrowRuntimeError(env, fmt.Sprintf("panic in HTTP request: %v", r)))
			}
		}()

		result, err := requestFunc()
		if err != nil {
			promise.reject(err)
		} else {
			promise.resolve(result)
		}
	}()

	// Set the _promise field
	classInstance := instance.(*ClassInstance)
	classInstance.Fields["_promise"] = promise

	return classInstance, nil
}

// httpResponseRender renders an HTML template with data - 3.12
func httpResponseRender(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpResponse", thisVal)
	}

	templatePath := utils.ToString(args[0])
	
	// Extract data from Map instance
	var dataMap map[string]any
	if mapInstance, ok := args[1].(*ClassInstance); ok && mapInstance.ClassName == "Map" {
		objMap, err := MapToObject((*Env)(e), mapInstance)
		if err != nil {
			return nil, err
		}
		dataMap = objMap
	} else if m, ok := args[1].(map[string]any); ok {
		dataMap = m
	}

	// Simple template rendering - replace {{key}} with values
	// In a real implementation, you'd read from a file and use a proper template engine
	templateContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
</head>
<body>
    <h1>Rendered Template: %s</h1>
    <pre>%v</pre>
</body>
</html>`, templatePath, templatePath, dataMap)

	// Replace template variables if any
	for key, value := range dataMap {
		placeholder := fmt.Sprintf("{{%s}}", key)
		templateContent = strings.ReplaceAll(templateContent, placeholder, utils.ToString(value))
	}

	resp := instance.Fields["_writer"].(*httpResponse)
	resp.statusCode = 200
	resp.sendHTML(templateContent)

	return nil, nil
}

// WebSocket Support - 3.2

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	conn     *websocket.Conn
	env      *Env
	handlers map[string][]common.Func // event -> handlers
	mu       sync.RWMutex
	done     chan struct{}
}

// httpServerWs registers a WebSocket route
func httpServerWs(e *common.Env, args []any) (any, error) {
	thisVal, _ := e.This()
	instance, ok := thisVal.(*ClassInstance)
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "HttpServer", thisVal)
	}

	router := instance.Fields["router"].(*httpRouter)
	path := utils.ToString(args[0])
	handler, ok := common.ExtractFunc(args[1])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[1])
	}

	router.mu.Lock()
	router.wsHandlers[path] = handler
	router.mu.Unlock()

	return nil, nil
}

// handleWebSocket handles WebSocket upgrade and connection
func (r *httpRouter) handleWebSocket(env *common.Env, w http.ResponseWriter, req *http.Request, handler common.Func) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	// Create WebSocket connection wrapper
	wsConn := &WebSocketConnection{
		conn:     conn,
		env:      env,
		handlers: make(map[string][]common.Func),
		done:     make(chan struct{}),
	}

	// Create WebSocket object for Polyloft
	wsInstance := createWebSocketInstance(env, wsConn)

	// Call the handler with the WebSocket instance
	_, err = handler((*common.Env)(env), []any{wsInstance})
	if err != nil {
		fmt.Printf("WebSocket handler error: %v\n", err)
		return
	}

	// Start message reading loop
	go wsConn.readMessages()

	// Wait until connection is closed
	<-wsConn.done
}

// readMessages reads messages from the WebSocket connection
func (ws *WebSocketConnection) readMessages() {
	defer close(ws.done)

	for {
		messageType, message, err := ws.conn.ReadMessage()
		if err != nil {
			// Connection closed or error
			ws.triggerEvent("close", string(message))
			break
		}

		// Handle different message types
		switch messageType {
		case websocket.TextMessage:
			ws.triggerEvent("message", string(message))
		case websocket.BinaryMessage:
			ws.triggerEvent("binary", string(message))
		}
	}
}

// triggerEvent triggers all handlers for an event
func (ws *WebSocketConnection) triggerEvent(event string, data string) {
	ws.mu.RLock()
	handlers := ws.handlers[event]
	ws.mu.RUnlock()

	for _, handler := range handlers {
		// Call handler with data
		handler((*common.Env)(ws.env), []any{data})
	}
}

// createWebSocketInstance creates a Polyloft WebSocket instance
func createWebSocketInstance(env *Env, wsConn *WebSocketConnection) *ClassInstance {
	// Create a simple object with WebSocket methods
	instance := &ClassInstance{
		ClassName: "WebSocket",
		Fields:    make(map[string]any),
		Methods:   make(map[string]common.Func),
	}

	// Store the connection
	instance.Fields["_conn"] = wsConn

	// Add send method
	instance.Methods["send"] = common.Func(func(e *common.Env, args []any) (any, error) {
		message := utils.ToString(args[0])
		return nil, wsConn.conn.WriteMessage(websocket.TextMessage, []byte(message))
	})

	// Add broadcast method (sends to all connections - simplified version)
	instance.Methods["broadcast"] = common.Func(func(e *common.Env, args []any) (any, error) {
		message := utils.ToString(args[0])
		// For now, just send to this connection
		// In a full implementation, this would send to all connected clients
		return nil, wsConn.conn.WriteMessage(websocket.TextMessage, []byte(message))
	})

	// Add on method for event handlers
	instance.Methods["on"] = common.Func(func(e *common.Env, args []any) (any, error) {
		event := utils.ToString(args[0])
		handler, ok := common.ExtractFunc(args[1])
		if !ok {
			return nil, ThrowTypeError(env, "function", args[1])
		}

		wsConn.mu.Lock()
		wsConn.handlers[event] = append(wsConn.handlers[event], handler)
		wsConn.mu.Unlock()

		return nil, nil
	})

	// Add close method
	instance.Methods["close"] = common.Func(func(e *common.Env, args []any) (any, error) {
		return nil, wsConn.conn.Close()
	})

	return instance
}

// Update handleRequest to check for WebSocket routes
func (r *httpRouter) isWebSocketRequest(path string) (common.Func, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	handler, ok := r.wsHandlers[path]
	return handler, ok
}

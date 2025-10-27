package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
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
		}, common.Func(httpResponseHtml), []string{})

	// Step 3: Create HttpServer builder and get its type BEFORE building
	httpServerBuilder := NewClassBuilder("HttpServer").
		AddField("router", ast.ANY, []string{"private"}).
		SetBuiltinConstructor([]ast.Parameter{}, common.Func(newHttpServer)).
		AddBuiltinMethod("get", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerGet), []string{}).
		AddBuiltinMethod("post", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPost), []string{}).
		AddBuiltinMethod("put", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerPut), []string{}).
		AddBuiltinMethod("delete", voidType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "handler", Type: ast.ANY},
		}, common.Func(httpServerDelete), []string{}).
		AddBuiltinMethod("listen", mapType, []ast.Parameter{
			{Name: "port", Type: stringType},
		}, common.Func(httpServerListen), []string{})

	httpServerType := httpServerBuilder.GetType()

	// Step 5: Create Http class with static methods using proper type references
	httpStaticClassBuilder := NewClassBuilder("Http").
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
		AddStaticMethod("request", mapType, []ast.Parameter{
			{Name: "method", Type: stringType},
			{Name: "url", Type: stringType},
			{Name: "data", Type: ast.ANY},
			{Name: "timeout", Type: intType},
			{Name: "headers", Type: mapType},
		}, common.Func(httpRequest)).
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

	return createHttpResponse(resp, body), nil
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

	return createHttpResponse(resp, body), nil
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

	return createHttpResponse(resp, body), nil
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

	return createHttpResponse(resp, body), nil
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

	return createHttpResponse(resp, body), nil
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
		routes: make(map[string]map[string]common.Func),
		mu:     &sync.RWMutex{},
	}

	instance.Fields["router"] = router

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
	handler, ok := common.ExtractFunc(args[1])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[1])
	}

	router.addRoute("GET", path, handler)
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
	handler, ok := common.ExtractFunc(args[1])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[1])
	}

	router.addRoute("POST", path, handler)
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
	handler, ok := common.ExtractFunc(args[1])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[1])
	}

	router.addRoute("PUT", path, handler)
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
	handler, ok := common.ExtractFunc(args[1])
	if !ok {
		return nil, ThrowTypeError((*Env)(e), "function", args[1])
	}

	router.addRoute("DELETE", path, handler)
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

	// Create HTTP handler
	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		objMap, err := MapToObject(mapInstance)
		if err != nil {
			return nil, err
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
func createHttpResponse(resp *http.Response, body []byte) map[string]any {
	return map[string]any{
		"status":     float64(resp.StatusCode),
		"statusText": resp.Status,
		"body":       string(body),
		"headers":    convertHeaders(resp.Header),
	}
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

// httpRouter manages HTTP routes
type httpRouter struct {
	routes map[string]map[string]common.Func // method -> path -> handler
	mu     *sync.RWMutex
}

func (r *httpRouter) addRoute(method, path string, handler common.Func) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]common.Func)
	}
	r.routes[method][path] = handler
}

func (r *httpRouter) handleRequest(env *common.Env, w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	methodRoutes := r.routes[req.Method]
	r.mu.RUnlock()

	// Find matching route
	handler, found := methodRoutes[req.URL.Path]
	if !found {
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

	// Call handler with class instances
	_, err := handler(env, []any{requestInstance, responseInstance})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, err)))
	}
}

// httpResponse wraps http.ResponseWriter with convenience methods
type httpResponse struct {
	writer     http.ResponseWriter
	statusCode int
	headers    map[string]string
	sent       bool
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
		objMap, err := MapToObject(mapInstance)
		if err != nil {
			json.NewEncoder(r.writer).Encode(map[string]string{"error": "failed to convert Map instance: " + err.Error()})
			return
		}
		json.NewEncoder(r.writer).Encode(objMap)
		return
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

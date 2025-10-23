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
	// Create Http class with static methods for HTTP client operations
	NewClassBuilder("Http").
		AddStaticMethod("get", ast.ANY, []ast.Parameter{
			{Name: "url", Type: ast.TypeFromString("String")},
			{Name: "timeout", Type: ast.TypeFromString("Int")},
		}, common.Func(httpGet)).
		AddStaticMethod("post", ast.ANY, []ast.Parameter{
			{Name: "url", Type: ast.TypeFromString("String")},
			{Name: "data", Type: ast.TypeFromString("any")},
			{Name: "timeout", Type: ast.TypeFromString("Int")},
		}, common.Func(httpPost)).
		AddStaticMethod("put", ast.ANY, []ast.Parameter{
			{Name: "url", Type: ast.TypeFromString("String")},
			{Name: "data", Type: ast.TypeFromString("any")},
			{Name: "timeout", Type: ast.TypeFromString("Int")},
		}, common.Func(httpPut)).
		AddStaticMethod("delete", ast.ANY, []ast.Parameter{
			{Name: "url", Type: ast.TypeFromString("String")},
			{Name: "timeout", Type: ast.TypeFromString("Int")},
		}, common.Func(httpDelete)).
		AddStaticMethod("request", ast.ANY, []ast.Parameter{
			{Name: "method", Type: ast.TypeFromString("String")},
			{Name: "url", Type: ast.TypeFromString("String")},
			{Name: "data", Type: ast.TypeFromString("any")},
			{Name: "timeout", Type: ast.TypeFromString("Int")},
			{Name: "headers", Type: ast.TypeFromString("Map")},
		}, common.Func(httpRequest)).
		AddStaticMethod("createServer", &ast.Type{Name: "HttpServer", IsBuiltin: true}, []ast.Parameter{
			{Name: "debug", Type: ast.TypeFromString("Bool")},
		}, common.Func(createHttpServer)).
		AddStaticMethod("createServer", &ast.Type{Name: "HttpServer", IsBuiltin: true}, []ast.Parameter{}, common.Func(createHttpServer)).
		BuildStatic(env)

	// Create HttpServer class (instanciable)
	NewClassBuilder("HttpServer").
		AddField("router", &ast.Type{Name: "HttpRouter", IsBuiltin: true}, []string{"private"}).
		SetBuiltinConstructor([]ast.Parameter{}, common.Func(newHttpServer)).
		AddBuiltinMethod("get", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("String")},
			{Name: "handler", Type: ast.TypeFromString("Function")},
		}, common.Func(httpServerGet), []string{}).
		AddBuiltinMethod("post", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("String")},
			{Name: "handler", Type: ast.TypeFromString("Function")},
		}, common.Func(httpServerPost), []string{}).
		AddBuiltinMethod("put", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("String")},
			{Name: "handler", Type: ast.TypeFromString("Function")},
		}, common.Func(httpServerPut), []string{}).
		AddBuiltinMethod("delete", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("String")},
			{Name: "handler", Type: ast.TypeFromString("Function")},
		}, common.Func(httpServerDelete), []string{}).
		AddBuiltinMethod("listen", ast.ANY, []ast.Parameter{
			{Name: "port", Type: ast.TypeFromString("String")},
		}, common.Func(httpServerListen), []string{}).
		Build(env)

	// Create HttpRequest class (represents incoming requests)
	NewClassBuilder("HttpRequest").
		AddField("method", &ast.Type{Name: "string", IsBuiltin: true}, []string{"public"}).
		AddField("path", &ast.Type{Name: "string", IsBuiltin: true}, []string{"public"}).
		AddField("url", &ast.Type{Name: "string", IsBuiltin: true}, []string{"public"}).
		AddField("headers", &ast.Type{Name: "Map", IsBuiltin: true}, []string{"public"}).
		AddField("query", &ast.Type{Name: "Map", IsBuiltin: true}, []string{"public"}).
		AddField("body", ast.ANY, []string{"public"}).
		Build(env)

	// Create HttpResponse class (represents outgoing responses)
	NewClassBuilder("HttpResponse").
		AddField("_writer", ast.ANY, []string{"private"}).
		AddField("_statusCode", &ast.Type{Name: "int", IsBuiltin: true}, []string{"private"}).
		AddField("_headers", &ast.Type{Name: "Map", IsBuiltin: true}, []string{"private"}).
		AddField("_sent", &ast.Type{Name: "bool", IsBuiltin: true}, []string{"private"}).
		AddBuiltinMethod("status", &ast.Type{Name: "HttpResponse", IsBuiltin: true}, []ast.Parameter{
			{Name: "code", Type: ast.TypeFromString("Int")},
		}, common.Func(httpResponseStatus), []string{}).
		AddBuiltinMethod("header", &ast.Type{Name: "HttpResponse", IsBuiltin: true}, []ast.Parameter{
			{Name: "name", Type: ast.TypeFromString("String")},
			{Name: "value", Type: ast.TypeFromString("String")},
		}, common.Func(httpResponseHeader), []string{}).
		AddBuiltinMethod("json", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "data", Type: ast.TypeFromString("any")},
		}, common.Func(httpResponseJson), []string{}).
		AddBuiltinMethod("send", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "text", Type: ast.TypeFromString("String")},
		}, common.Func(httpResponseSend), []string{}).
		AddBuiltinMethod("html", &ast.Type{Name: "void", IsBuiltin: true}, []ast.Parameter{
			{Name: "html", Type: ast.TypeFromString("String")},
		}, common.Func(httpResponseHtml), []string{}).
		Build(env)
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
	env := (*Env)(e)
	serverConstructor, _ := env.Get("HttpServer")

	if ctor, ok := serverConstructor.(*common.ClassConstructor); ok {
		return ctor.Func(e, []any{})
	}

	return nil, ThrowInitializationError((*Env)(e), "HttpServer class")
}

// newHttpServer creates a new HttpServer instance
func newHttpServer(e *common.Env, args []any) (any, error) {
	// Get the instance from the environment (created by createClassInstance)
	thisVal, exists := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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
	thisVal, _ := e.Get("this")
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

	// Get class definitions to create proper instances
	envTyped := (*Env)(env)

	// Create HttpRequest instance using the class constructor
	requestClassVal, _ := envTyped.Get("HttpRequest")
	var requestInstance *ClassInstance
	if requestCtor, ok := requestClassVal.(*common.ClassConstructor); ok {
		inst, err := requestCtor.Func(env, []any{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "Failed to create HttpRequest: %v"}`, err)))
			return
		}
		if ri, ok := inst.(*ClassInstance); ok {
			requestInstance = ri
			// Set the fields
			requestInstance.Fields["method"] = req.Method
			requestInstance.Fields["path"] = req.URL.Path
			requestInstance.Fields["url"] = req.URL.String()
			requestInstance.Fields["headers"] = convertHeaders(req.Header)
			requestInstance.Fields["query"] = queryParams
			requestInstance.Fields["body"] = bodyData
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
	responseClassVal, _ := envTyped.Get("HttpResponse")
	var responseInstance *ClassInstance
	if responseCtor, ok := responseClassVal.(*common.ClassConstructor); ok {
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

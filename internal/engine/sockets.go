package engine

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallSocketsModule installs Socket classes for network communication
func InstallSocketsModule(env *Env, opts Options) error {
	// Get type references
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	voidType := &ast.Type{Name: "void", IsBuiltin: true}
	bytesType := common.BuiltinTypeBytes.GetTypeDefinition(env)

	// ========================================
	// Socket class - TCP socket
	// ========================================
	socketBuilder := NewClassBuilder("Socket").
		AddField("_conn", ast.ANY, []string{"private"}).
		AddField("_reader", ast.ANY, []string{"private"}).
		AddField("connected", boolType, []string{"public"}).
		AddField("remoteAddr", stringType, []string{"public"}).
		AddField("localAddr", stringType, []string{"public"})

	socketType := socketBuilder.GetType()

	// Constructor: Socket() - not connected
	socketBuilder.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_conn"] = nil
		instance.Fields["_reader"] = nil
		instance.Fields["connected"] = false
		instance.Fields["remoteAddr"] = ""
		instance.Fields["localAddr"] = ""
		return nil, nil
	})

	// connect(host: String, port: Int, timeout: Int) -> Bool
	socketBuilder.AddBuiltinMethod("connect", boolType, []ast.Parameter{
		{Name: "host", Type: stringType},
		{Name: "port", Type: intType},
		{Name: "timeout", Type: intType, IsVariadic: false},
	}, func(callEnv *common.Env, args []any) (any, error) {
		host := utils.ToString(args[0])
		port, ok := utils.AsInt(args[1])
		if !ok {
			return false, ThrowTypeError((*Env)(callEnv), "int", args[1])
		}

		timeout := 10 * time.Second
		if len(args) > 2 {
			if t, ok := utils.AsInt(args[2]); ok {
				timeout = time.Duration(t) * time.Second
			}
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return false, nil
		}

		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_conn"] = conn
		instance.Fields["_reader"] = bufio.NewReader(conn)
		instance.Fields["connected"] = true
		instance.Fields["remoteAddr"] = conn.RemoteAddr().String()
		instance.Fields["localAddr"] = conn.LocalAddr().String()
		return true, nil
	}, []string{})

	// send(data: String) -> Int
	socketBuilder.AddBuiltinMethod("send", intType, []ast.Parameter{
		{Name: "data", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return 0, ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		data := utils.ToString(args[0])
		n, err := conn.Write([]byte(data))
		if err != nil {
			return 0, err
		}
		return n, nil
	}, []string{})

	// sendBytes(data: Bytes) -> Int
	socketBuilder.AddBuiltinMethod("sendBytes", intType, []ast.Parameter{
		{Name: "data", Type: bytesType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return 0, ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		if bytesInst, ok := args[0].(*ClassInstance); ok {
			data := bytesInst.Fields["_data"].([]byte)
			n, err := conn.Write(data)
			if err != nil {
				return 0, err
			}
			return n, nil
		}
		return 0, ThrowTypeError((*Env)(callEnv), "Bytes", args[0])
	}, []string{})

	// recv(size: Int, timeout: Int) -> String
	socketBuilder.AddBuiltinMethod("recv", stringType, []ast.Parameter{
		{Name: "size", Type: intType, IsVariadic: false},
		{Name: "timeout", Type: intType, IsVariadic: false},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return "", ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		reader := instance.Fields["_reader"].(*bufio.Reader)

		size := 1024
		if len(args) > 0 {
			if s, ok := utils.AsInt(args[0]); ok {
				size = s
			}
		}

		timeout := 5 * time.Second
		if len(args) > 1 {
			if t, ok := utils.AsInt(args[1]); ok {
				timeout = time.Duration(t) * time.Second
			}
		}

		buf := make([]byte, size)
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := reader.Read(buf)
		if err != nil {
			return "", nil
		}
		return string(buf[:n]), nil
	}, []string{})

	// recvBytes(size: Int, timeout: Int) -> Bytes
	socketBuilder.AddBuiltinMethod("recvBytes", bytesType, []ast.Parameter{
		{Name: "size", Type: intType, IsVariadic: false},
		{Name: "timeout", Type: intType, IsVariadic: false},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		reader := instance.Fields["_reader"].(*bufio.Reader)

		size := 1024
		if len(args) > 0 {
			if s, ok := utils.AsInt(args[0]); ok {
				size = s
			}
		}

		timeout := 5 * time.Second
		if len(args) > 1 {
			if t, ok := utils.AsInt(args[1]); ok {
				timeout = time.Duration(t) * time.Second
			}
		}

		buf := make([]byte, size)
		conn.SetReadDeadline(time.Now().Add(timeout))
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}

		// Create Bytes instance
		bytesClassDef := common.BuiltinTypeBytes.GetClassDefinition(callEnv)
		if bytesClassDef == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Bytes class not found")
		}

		bytesInst := &ClassInstance{
			ClassName:   "Bytes",
			Fields:      make(map[string]any),
			Methods:     make(map[string]common.Func),
			ParentClass: bytesClassDef,
		}
		bytesInst.Fields["_data"] = buf[:n]
		return bytesInst, nil
	}, []string{})

	// close() -> Void
	socketBuilder.AddBuiltinMethod("close", voidType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if ok && conn != nil {
			conn.Close()
		}
		instance.Fields["_conn"] = nil
		instance.Fields["_reader"] = nil
		instance.Fields["connected"] = false
		return nil, nil
	}, []string{})

	// setReadTimeout(timeout: Int) -> Void
	socketBuilder.AddBuiltinMethod("setReadTimeout", voidType, []ast.Parameter{
		{Name: "timeout", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		timeout, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		return nil, nil
	}, []string{})

	// setWriteTimeout(timeout: Int) -> Void
	socketBuilder.AddBuiltinMethod("setWriteTimeout", voidType, []ast.Parameter{
		{Name: "timeout", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		conn, ok := instance.Fields["_conn"].(net.Conn)
		if !ok || conn == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "socket not connected")
		}

		timeout, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}

		conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		return nil, nil
	}, []string{})

	_, err := socketBuilder.Build(env)
	if err != nil {
		return err
	}

	// ========================================
	// ServerSocket class - TCP server socket
	// ========================================
	serverSocketBuilder := NewClassBuilder("ServerSocket").
		AddField("_listener", ast.ANY, []string{"private"}).
		AddField("listening", boolType, []string{"public"}).
		AddField("address", stringType, []string{"public"})

	// Constructor: ServerSocket() - not listening
	serverSocketBuilder.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_listener"] = nil
		instance.Fields["listening"] = false
		instance.Fields["address"] = ""
		return nil, nil
	})

	// bind(host: String, port: Int) -> Bool
	serverSocketBuilder.AddBuiltinMethod("bind", boolType, []ast.Parameter{
		{Name: "host", Type: stringType},
		{Name: "port", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		host := utils.ToString(args[0])
		port, ok := utils.AsInt(args[1])
		if !ok {
			return false, ThrowTypeError((*Env)(callEnv), "int", args[1])
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return false, nil
		}

		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_listener"] = listener
		instance.Fields["listening"] = true
		instance.Fields["address"] = listener.Addr().String()
		return true, nil
	}, []string{})

	// accept() -> Socket
	serverSocketBuilder.AddBuiltinMethod("accept", socketType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		listener, ok := instance.Fields["_listener"].(net.Listener)
		if !ok || listener == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "server socket not listening")
		}

		conn, err := listener.Accept()
		if err != nil {
			return nil, err
		}

		// Create Socket instance
		socketClassDef := common.BuiltinTypeSocket.GetClassDefinition(callEnv)
		if socketClassDef == nil {
			return nil, ThrowRuntimeError((*Env)(callEnv), "Socket class not found")
		}

		socketInst := &ClassInstance{
			ClassName:   "Socket",
			Fields:      make(map[string]any),
			Methods:     make(map[string]common.Func),
			ParentClass: socketClassDef,
		}
		socketInst.Fields["_conn"] = conn
		socketInst.Fields["_reader"] = bufio.NewReader(conn)
		socketInst.Fields["connected"] = true
		socketInst.Fields["remoteAddr"] = conn.RemoteAddr().String()
		socketInst.Fields["localAddr"] = conn.LocalAddr().String()
		return socketInst, nil
	}, []string{})

	// close() -> Void
	serverSocketBuilder.AddBuiltinMethod("close", voidType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)

		listener, ok := instance.Fields["_listener"].(net.Listener)
		if ok && listener != nil {
			listener.Close()
		}
		instance.Fields["_listener"] = nil
		instance.Fields["listening"] = false
		return nil, nil
	}, []string{})

	_, err = serverSocketBuilder.Build(env)
	return err
}

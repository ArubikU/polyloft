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

// InstallNetModule installs the complete Net module with networking functions
func InstallNetModule(env *Env, opts Options) {
	// Get type references from already-installed builtin types
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	mapType := common.BuiltinTypeMap.GetTypeDefinition(env)
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)

	netClass := NewClassBuilder("Net").
		AddStaticMethod("listen", mapType, []ast.Parameter{
			{Name: "addr", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			addr, _ := env.Get("addr")
			addrStr := utils.ToString(addr)
			ln, err := net.Listen("tcp", addrStr)
			if err != nil {
				return nil, err
			}
			server := map[string]any{}
			server["addr"] = ln.Addr().String()
			server["close"] = Func(func(_ *Env, _ []any) (any, error) { return nil, ln.Close() })
			server["accept"] = Func(func(_ *Env, _ []any) (any, error) {
				c, err := ln.Accept()
				if err != nil {
					return nil, err
				}
				br := bufio.NewReader(c)
				conn := map[string]any{}
				conn["remote"] = c.RemoteAddr().String()
				conn["send"] = Func(func(_ *Env, args []any) (any, error) {
					if len(args) < 1 {
						return float64(0), nil
					}
					s := utils.ToString(args[0])
					n, err := c.Write([]byte(s))
					if err != nil {
						return nil, err
					}
					return float64(n), nil
				})
				conn["recv"] = Func(func(_ *Env, args []any) (any, error) {
					n := 1024
					if len(args) > 0 {
						if v, ok := utils.AsInt(args[0]); ok {
							n = v
						}
					}
					buf := make([]byte, n)
					c.SetReadDeadline(time.Now().Add(5 * time.Second))
					r, err := br.Read(buf)
					if err != nil {
						return "", nil
					}
					return string(buf[:r]), nil
				})
				conn["close"] = Func(func(_ *Env, _ []any) (any, error) { return nil, c.Close() })
				return conn, nil
			})
			return server, nil
		})).
		AddStaticMethod("connect", mapType, []ast.Parameter{
			{Name: "addr", Type: stringType},
			{Name: "timeout", Type: intType, IsVariadic: false},
		}, Func(func(env *Env, _ []any) (any, error) {
			addr, _ := env.Get("addr")
			addrStr := utils.ToString(addr)

			timeout := 10 * time.Second
			if timeoutVal, ok := env.Get("timeout"); ok {
				if t, ok := utils.AsInt(timeoutVal); ok {
					timeout = time.Duration(t) * time.Second
				}
			}

			conn, err := net.DialTimeout("tcp", addrStr, timeout)
			if err != nil {
				return nil, err
			}

			br := bufio.NewReader(conn)
			client := map[string]any{}
			client["remote"] = conn.RemoteAddr().String()
			client["local"] = conn.LocalAddr().String()

			client["send"] = Func(func(_ *Env, args []any) (any, error) {
				if len(args) < 1 {
					return float64(0), nil
				}
				s := utils.ToString(args[0])
				n, err := conn.Write([]byte(s))
				if err != nil {
					return nil, err
				}
				return float64(n), nil
			})

			client["recv"] = Func(func(_ *Env, args []any) (any, error) {
				n := 1024
				if len(args) > 0 {
					if v, ok := utils.AsInt(args[0]); ok {
						n = v
					}
				}
				buf := make([]byte, n)

				timeout := 5 * time.Second
				if len(args) > 1 {
					if t, ok := utils.AsInt(args[1]); ok {
						timeout = time.Duration(t) * time.Second
					}
				}

				conn.SetReadDeadline(time.Now().Add(timeout))
				r, err := br.Read(buf)
				if err != nil {
					return "", nil
				}
				return string(buf[:r]), nil
			})

			client["close"] = Func(func(_ *Env, _ []any) (any, error) { return nil, conn.Close() })

			return client, nil
		})).
		AddStaticMethod("listenUdp", mapType, []ast.Parameter{
			{Name: "addr", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			addr, _ := env.Get("addr")
			addrStr := utils.ToString(addr)

			udpAddr, err := net.ResolveUDPAddr("udp", addrStr)
			if err != nil {
				return nil, err
			}

			conn, err := net.ListenUDP("udp", udpAddr)
			if err != nil {
				return nil, err
			}

			server := map[string]any{}
			server["addr"] = conn.LocalAddr().String()
			server["close"] = Func(func(_ *Env, _ []any) (any, error) { return nil, conn.Close() })

			server["recv"] = Func(func(_ *Env, args []any) (any, error) {
				n := 1024
				if len(args) > 0 {
					if v, ok := utils.AsInt(args[0]); ok {
						n = v
					}
				}
				buf := make([]byte, n)

				timeout := 5 * time.Second
				if len(args) > 1 {
					if t, ok := utils.AsInt(args[1]); ok {
						timeout = time.Duration(t) * time.Second
					}
				}

				conn.SetReadDeadline(time.Now().Add(timeout))
				r, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					return nil, err
				}

				result := map[string]any{}
				result["data"] = string(buf[:r])
				result["addr"] = addr.String()
				return result, nil
			})

			server["send"] = Func(func(e2 *Env, args []any) (any, error) {
				if len(args) < 2 {
					return nil, ThrowArityError(e2, 2, len(args))
				}
				data := utils.ToString(args[0])
				addrArg := utils.ToString(args[1])

				udpAddr, err := net.ResolveUDPAddr("udp", addrArg)
				if err != nil {
					return nil, err
				}

				n, err := conn.WriteToUDP([]byte(data), udpAddr)
				if err != nil {
					return nil, err
				}
				return float64(n), nil
			})

			return server, nil
		})).
		AddStaticMethod("dialUdp", mapType, []ast.Parameter{
			{Name: "addr", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			addr, _ := env.Get("addr")
			addrStr := utils.ToString(addr)

			udpAddr, err := net.ResolveUDPAddr("udp", addrStr)
			if err != nil {
				return nil, err
			}

			conn, err := net.DialUDP("udp", nil, udpAddr)
			if err != nil {
				return nil, err
			}

			client := map[string]any{}
			client["remote"] = conn.RemoteAddr().String()
			client["local"] = conn.LocalAddr().String()
			client["close"] = Func(func(_ *Env, _ []any) (any, error) { return nil, conn.Close() })

			client["send"] = Func(func(_ *Env, args []any) (any, error) {
				if len(args) < 1 {
					return float64(0), nil
				}
				data := utils.ToString(args[0])
				n, err := conn.Write([]byte(data))
				if err != nil {
					return nil, err
				}
				return float64(n), nil
			})

			client["recv"] = Func(func(_ *Env, args []any) (any, error) {
				n := 1024
				if len(args) > 0 {
					if v, ok := utils.AsInt(args[0]); ok {
						n = v
					}
				}
				buf := make([]byte, n)

				timeout := 5 * time.Second
				if len(args) > 1 {
					if t, ok := utils.AsInt(args[1]); ok {
						timeout = time.Duration(t) * time.Second
					}
				}

				conn.SetReadDeadline(time.Now().Add(timeout))
				r, err := conn.Read(buf)
				if err != nil {
					return "", nil
				}
				return string(buf[:r]), nil
			})

			return client, nil
		})).
		AddStaticMethod("resolveIp", arrayType, []ast.Parameter{
			{Name: "hostname", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			hostname, _ := env.Get("hostname")
			hostnameStr := utils.ToString(hostname)

			ips, err := net.LookupIP(hostnameStr)
			if err != nil {
				return nil, err
			}

			result := make([]any, len(ips))
			for i, ip := range ips {
				result[i] = ip.String()
			}
			return result, nil
		})).
		AddStaticMethod("resolveHost", arrayType, []ast.Parameter{
			{Name: "ip", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			ip, _ := env.Get("ip")
			ipStr := utils.ToString(ip)

			names, err := net.LookupAddr(ipStr)
			if err != nil {
				return nil, err
			}

			result := make([]any, len(names))
			for i, name := range names {
				result[i] = name
			}
			return result, nil
		})).
		AddStaticMethod("getLocalIPs", arrayType, []ast.Parameter{}, Func(func(_ *Env, _ []any) (any, error) {
			addrs, err := net.InterfaceAddrs()
			if err != nil {
				return nil, err
			}

			var ips []any
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
			return ips, nil
		})).
		AddStaticMethod("isPortOpen", boolType, []ast.Parameter{
			{Name: "host", Type: stringType},
			{Name: "port", Type: intType},
			{Name: "timeout", Type: intType, IsVariadic: false},
		}, Func(func(env *Env, _ []any) (any, error) {
			host, _ := env.Get("host")
			portVal, _ := env.Get("port")
			hostStr := utils.ToString(host)
			port, ok := utils.AsInt(portVal)
			if !ok {
				return nil, fmt.Errorf("port must be a number")
			}

			timeout := 3 * time.Second
			if timeoutVal, ok := env.Get("timeout"); ok {
				if t, ok := utils.AsInt(timeoutVal); ok {
					timeout = time.Duration(t) * time.Second
				}
			}

			addr := fmt.Sprintf("%s:%d", hostStr, port)
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				return false, nil
			}
			conn.Close()
			return true, nil
		}))

	_, err := netClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
}

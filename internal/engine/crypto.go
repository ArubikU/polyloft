package engine

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

// InstallCryptoModule installs the complete Crypto module with cryptographic functions
func InstallCryptoModule(env *Env, opts Options) {
	// Get type references from already-installed builtin types
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)

	cryptoClass := NewClassBuilder("Crypto").
		AddStaticMethod("md5", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			hash := md5.Sum([]byte(dataStr))
			return hex.EncodeToString(hash[:]), nil
		})).
		AddStaticMethod("sha1", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			hash := sha1.Sum([]byte(dataStr))
			return hex.EncodeToString(hash[:]), nil
		})).
		AddStaticMethod("sha256", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			hash := sha256.Sum256([]byte(dataStr))
			return hex.EncodeToString(hash[:]), nil
		})).
		AddStaticMethod("sha512", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			hash := sha512.Sum512([]byte(dataStr))
			return hex.EncodeToString(hash[:]), nil
		})).
		AddStaticMethod("base64Encode", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			encoded := base64.StdEncoding.EncodeToString([]byte(dataStr))
			return encoded, nil
		})).
		AddStaticMethod("base64Decode", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			decoded, err := base64.StdEncoding.DecodeString(dataStr)
			if err != nil {
				return nil, err
			}
			return string(decoded), nil
		})).
		AddStaticMethod("hexEncode", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			encoded := hex.EncodeToString([]byte(dataStr))
			return encoded, nil
		})).
		AddStaticMethod("hexDecode", stringType, []ast.Parameter{
			{Name: "data", Type: stringType},
		}, Func(func(env *Env, _ []any) (any, error) {
			data, _ := env.Get("data")
			dataStr := utils.ToString(data)
			decoded, err := hex.DecodeString(dataStr)
			if err != nil {
				return nil, err
			}
			return string(decoded), nil
		}))

	_, err := cryptoClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
}

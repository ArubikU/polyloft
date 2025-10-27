package engine

import (
	"bufio"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func InstallIOModule(env *Env, opts Options) error {
	// Get type references from already-installed builtin types
	stringType := common.BuiltinTypeString.GetTypeDefinition(env)
	intType := common.BuiltinTypeInt.GetTypeDefinition(env)
	boolType := common.BuiltinTypeBool.GetTypeDefinition(env)
	mapType := common.BuiltinTypeMap.GetTypeDefinition(env)
	arrayType := common.BuiltinTypeArray.GetTypeDefinition(env)
	floatType := common.BuiltinTypeFloat.GetTypeDefinition(env)
	voidType := ast.ANY
	bytesType := common.BuiltinTypeBytes.GetTypeDefinition(env)

	// ========================================
	// Buffer class - In-memory buffer for reading/writing
	// ========================================
	bufferBuilder := NewClassBuilder("Buffer").
		AddField("_buffer", ast.ANY, []string{"private"})

	// Constructor: Buffer() - empty buffer
	bufferBuilder.AddBuiltinConstructor([]ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		instance.Fields["_buffer"] = &bytes.Buffer{}
		return nil, nil
	})

	// Constructor: Buffer(data: String) - buffer with initial string data
	bufferBuilder.AddBuiltinConstructor([]ast.Parameter{
		{Name: "data", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		str := utils.ToString(args[0])
		instance.Fields["_buffer"] = bytes.NewBufferString(str)
		return nil, nil
	})

	// write(data: String) -> Int - returns bytes written
	bufferBuilder.AddBuiltinMethod("write", intType, []ast.Parameter{
		{Name: "data", Type: stringType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)
		data := utils.ToString(args[0])
		n, err := buf.WriteString(data)
		if err != nil {
			return nil, err
		}
		return n, nil
	}, []string{})

	// writeBytes(data: Bytes) -> Int
	bufferBuilder.AddBuiltinMethod("writeBytes", intType, []ast.Parameter{
		{Name: "data", Type: bytesType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)

		if bytesInst, ok := args[0].(*ClassInstance); ok {
			data := bytesInst.Fields["_data"].([]byte)
			n, err := buf.Write(data)
			if err != nil {
				return nil, err
			}
			return n, nil
		}
		return 0, ThrowTypeError((*Env)(callEnv), "Bytes", args[0])
	}, []string{})

	// read(size: Int) -> String
	bufferBuilder.AddBuiltinMethod("read", stringType, []ast.Parameter{
		{Name: "size", Type: intType},
	}, func(callEnv *common.Env, args []any) (any, error) {
		size, ok := utils.AsInt(args[0])
		if !ok {
			return nil, ThrowTypeError((*Env)(callEnv), "int", args[0])
		}
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)

		data := make([]byte, size)
		n, err := buf.Read(data)
		if err != nil && err != io.EOF {
			return "", err
		}
		return string(data[:n]), nil
	}, []string{})

	// readAll() -> String
	bufferBuilder.AddBuiltinMethod("readAll", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)
		return buf.String(), nil
	}, []string{})

	// size() -> Int
	bufferBuilder.AddBuiltinMethod("size", intType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)
		return buf.Len(), nil
	}, []string{})

	// clear() -> Void
	bufferBuilder.AddBuiltinMethod("clear", voidType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)
		buf.Reset()
		return nil, nil
	}, []string{})

	// toString() -> String
	bufferBuilder.AddBuiltinMethod("toString", stringType, []ast.Parameter{}, func(callEnv *common.Env, args []any) (any, error) {
		thisVal, _ := callEnv.This()
		instance := thisVal.(*ClassInstance)
		buf := instance.Fields["_buffer"].(*bytes.Buffer)
		return buf.String(), nil
	}, []string{})

	_, err := bufferBuilder.Build(env)
	if err != nil {
		panic(err)
	}

	// ========================================
	// Writer interface - For types that can write data
	// ========================================
	writerBuilder := NewInterfaceBuilder("Writer").
		AddMethod("write", "Int", []ast.Parameter{
			{Name: "data", Type: stringType},
		}).
		AddMethod("writeBytes", "Int", []ast.Parameter{
			{Name: "data", Type: bytesType},
		}).
		AddMethod("close", "Void", []ast.Parameter{})

	_, err = writerBuilder.Build(env)
	if err != nil {
		panic(err)
	}

	// ========================================
	// Reader interface - For types that can read data
	// ========================================
	readerBuilder := NewInterfaceBuilder("Reader").
		AddMethod("read", "String", []ast.Parameter{
			{Name: "size", Type: intType},
		}).
		AddMethod("readAll", "String", []ast.Parameter{}).
		AddMethod("close", "Void", []ast.Parameter{})

	_, err = readerBuilder.Build(env)
	if err != nil {
		panic(err)
	}

	// ========================================
	// IO class - File system operations
	// ========================================
	ioClass := NewClassBuilder("IO").
		// File operations
		AddStaticMethod("readFile", stringType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return string(data), nil
		})).
		AddStaticMethod("writeFile", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "content", Type: stringType},
			{Name: "mode", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			content := utils.ToString(args[1])

			perm := os.FileMode(0644)
			if len(args) > 2 {
				if p, ok := utils.AsInt(args[2]); ok {
					perm = os.FileMode(p)
				}
			}

			if err := os.WriteFile(path, []byte(content), perm); err != nil {
				return nil, err
			}
			return true, nil
		})).
		AddStaticMethod("appendFile", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "content", Type: stringType},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			content := utils.ToString(args[1])

			file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			if _, err = file.WriteString(content); err != nil {
				return nil, err
			}
			return true, nil
		})).
		AddStaticMethod("copyFile", floatType, []ast.Parameter{
			{Name: "source", Type: stringType},
			{Name: "destination", Type: stringType},
		}, Func(func(_ *Env, args []any) (any, error) {
			src := utils.ToString(args[0])
			dst := utils.ToString(args[1])

			sourceFile, err := os.Open(src)
			if err != nil {
				return nil, err
			}
			defer sourceFile.Close()

			destFile, err := os.Create(dst)
			if err != nil {
				return nil, err
			}
			defer destFile.Close()

			bytesWritten, err := io.Copy(destFile, sourceFile)
			if err != nil {
				return nil, err
			}

			return float64(bytesWritten), nil
		})).
		AddStaticMethod("deleteFile", boolType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			if err := os.Remove(path); err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("moveFile", boolType, []ast.Parameter{
			{Name: "source", Type: stringType},
			{Name: "destination", Type: stringType},
		}, Func(func(_ *Env, args []any) (any, error) {
			src := utils.ToString(args[0])
			dst := utils.ToString(args[1])

			if err := os.Rename(src, dst); err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("fileExists", boolType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			_, err := os.Stat(path)
			return err == nil, nil
		})).
		AddStaticMethod("fileSize", floatType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			info, err := os.Stat(path)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"name":    info.Name(),
				"size":    float64(info.Size()),
				"mode":    float64(info.Mode()),
				"modTime": info.ModTime().Unix(),
				"isDir":   info.IsDir(),
			}, nil
		})).

		// openFile(path: String, mode: String, encoding?: String) -> Map
		// Modes: "r" (read), "w" (write), "a" (append), "r+" (read/write)
		// Encodings: "utf-8", "utf-16", "utf-16le", "utf-16be", "latin1", "windows-1252", "shift_jis", etc.
		AddStaticMethod("openFile", mapType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "mode", Type: stringType},
			{Name: "encoding", Type: nil, IsVariadic: true},
		}, Func(func(e *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			mode := utils.ToString(args[1])

			encodingName := "utf-8"
			if len(args) > 2 {
				encodingName = utils.ToString(args[2])
			}

			// Parse mode
			var flags int
			var isReadMode, isWriteMode bool

			switch mode {
			case "r":
				flags = os.O_RDONLY
				isReadMode = true
			case "w":
				flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				isWriteMode = true
			case "a":
				flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
				isWriteMode = true
			case "r+", "rw":
				flags = os.O_RDWR
				isReadMode = true
				isWriteMode = true
			case "w+":
				flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
				isReadMode = true
				isWriteMode = true
			case "a+":
				flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
				isReadMode = true
				isWriteMode = true
			default:
				return nil, io.ErrUnexpectedEOF
			}

			file, err := os.OpenFile(path, flags, 0644)
			if err != nil {
				return nil, err
			}

			// Get encoder/decoder based on encoding
			enc := getEncoding(encodingName)

			var reader *bufio.Reader
			var writer *bufio.Writer

			if isReadMode && enc != nil {
				decoder := enc.NewDecoder()
				reader = bufio.NewReader(decoder.Reader(file))
			} else if isReadMode {
				reader = bufio.NewReader(file)
			}

			if isWriteMode && enc != nil {
				encoder := enc.NewEncoder()
				writer = bufio.NewWriter(encoder.Writer(file))
			} else if isWriteMode {
				writer = bufio.NewWriter(file)
			}

			// Create file handle object
			fileHandle := map[string]any{
				"_file":     file,
				"_reader":   reader,
				"_writer":   writer,
				"_path":     path,
				"_mode":     mode,
				"_encoding": encodingName,
				"_closed":   false,
			}

			return fileHandle, nil
		})).

		// readFileWithEncoding(path: String, encoding?: String) -> String
		AddStaticMethod("readFileWithEncoding", stringType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "encoding", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			encodingName := "utf-8"
			if len(args) > 1 {
				encodingName = utils.ToString(args[1])
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			// If UTF-8 or no encoding specified, return as-is
			if encodingName == "utf-8" || encodingName == "" {
				return string(data), nil
			}

			// Decode from specified encoding
			enc := getEncoding(encodingName)
			if enc == nil {
				return string(data), nil // Fallback to default
			}

			decoder := enc.NewDecoder()
			decoded, err := decoder.Bytes(data)
			if err != nil {
				return nil, err
			}

			return string(decoded), nil
		})).

		// writeFileWithEncoding(path: String, content: String, encoding?: String) -> Bool
		AddStaticMethod("writeFileWithEncoding", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "content", Type: stringType},
			{Name: "encoding", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			content := utils.ToString(args[1])

			encodingName := "utf-8"
			if len(args) > 2 {
				encodingName = utils.ToString(args[2])
			}

			var data []byte

			// If UTF-8 or no encoding specified, write as-is
			if encodingName == "utf-8" || encodingName == "" {
				data = []byte(content)
			} else {
				// Encode to specified encoding
				enc := getEncoding(encodingName)
				if enc == nil {
					data = []byte(content) // Fallback to default
				} else {
					encoder := enc.NewEncoder()
					encoded, err := encoder.Bytes([]byte(content))
					if err != nil {
						return false, err
					}
					data = encoded
				}
			}

			if err := os.WriteFile(path, data, 0644); err != nil {
				return false, err
			}
			return true, nil
		})).

		// closeFile(fileHandle: Map) -> Bool
		AddStaticMethod("closeFile", boolType, []ast.Parameter{
			{Name: "fileHandle", Type: mapType},
		}, Func(func(e *Env, args []any) (any, error) {
			handle, ok := args[0].(map[string]any)
			if !ok {
				return false, ThrowTypeError(e, "map", args[0])
			}

			if closed, ok := handle["_closed"].(bool); ok && closed {
				return true, nil // Already closed
			}

			// Flush writer if exists
			if writer, ok := handle["_writer"].(*bufio.Writer); ok && writer != nil {
				writer.Flush()
			}

			// Close file
			if file, ok := handle["_file"].(*os.File); ok && file != nil {
				if err := file.Close(); err != nil {
					return false, err
				}
			}

			handle["_closed"] = true
			return true, nil
		})).

		// readFromFile(fileHandle: Map, size?: Int) -> String
		AddStaticMethod("readFromFile", stringType, []ast.Parameter{
			{Name: "fileHandle", Type: mapType},
			{Name: "size", Type: nil, IsVariadic: true},
		}, Func(func(e *Env, args []any) (any, error) {
			handle, ok := args[0].(map[string]any)
			if !ok {
				return "", ThrowTypeError(e, "map", args[0])
			}

			reader, ok := handle["_reader"].(*bufio.Reader)
			if !ok || reader == nil {
				return "", ThrowTypeError(e, "file opened in read mode", handle["_reader"])
			}

			// Read specific size or all
			if len(args) > 1 {
				size, ok := utils.AsInt(args[1])
				if !ok {
					return "", ThrowTypeError(e, "int", args[1])
				}
				buf := make([]byte, size)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					return "", err
				}
				return string(buf[:n]), nil
			}

			// Read all remaining
			data, err := io.ReadAll(reader)
			if err != nil {
				return "", err
			}
			return string(data), nil
		})).

		// writeToFile(fileHandle: Map, content: String) -> Int
		AddStaticMethod("writeToFile", intType, []ast.Parameter{
			{Name: "fileHandle", Type: mapType},
			{Name: "content", Type: stringType},
		}, Func(func(e *Env, args []any) (any, error) {
			handle, ok := args[0].(map[string]any)
			if !ok {
				return 0, ThrowTypeError(e, "map", args[0])
			}

			writer, ok := handle["_writer"].(*bufio.Writer)
			if !ok || writer == nil {
				return 0, ThrowTypeError(e, "file opened in write mode", handle["_writer"])
			}

			content := utils.ToString(args[1])
			n, err := writer.WriteString(content)
			if err != nil {
				return 0, err
			}

			return n, nil
		})).

		// flushFile(fileHandle: Map) -> Bool
		AddStaticMethod("flushFile", boolType, []ast.Parameter{
			{Name: "fileHandle", Type: mapType},
		}, Func(func(e *Env, args []any) (any, error) {
			handle, ok := args[0].(map[string]any)
			if !ok {
				return false, ThrowTypeError(e, "map", args[0])
			}

			writer, ok := handle["_writer"].(*bufio.Writer)
			if !ok || writer == nil {
				return true, nil // No writer, nothing to flush
			}

			if err := writer.Flush(); err != nil {
				return false, err
			}

			return true, nil
		})).

		// Directory operations
		AddStaticMethod("createDir", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "recursive", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			recursive := false
			if len(args) > 1 {
				if r, ok := args[1].(bool); ok {
					recursive = r
				}
			}

			var err error
			if recursive {
				err = os.MkdirAll(path, 0755)
			} else {
				err = os.Mkdir(path, 0755)
			}

			if err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("removeDir", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "recursive", Type: nil, IsVariadic: true},
		}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			recursive := false
			if len(args) > 1 {
				if r, ok := args[1].(bool); ok {
					recursive = r
				}
			}

			var err error
			if recursive {
				err = os.RemoveAll(path)
			} else {
				err = os.Remove(path)
			}

			if err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("listDir", arrayType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}

			result := make([]any, len(entries))
			for i, entry := range entries {
				info, _ := entry.Info()
				item := map[string]any{
					"name":  entry.Name(),
					"isDir": entry.IsDir(),
				}
				if info != nil {
					item["size"] = float64(info.Size())
					item["modTime"] = info.ModTime().Unix()
				}
				result[i] = item
			}
			return result, nil
		})).
		AddStaticMethod("walkDir", arrayType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			root := utils.ToString(args[0])

			var files []any
			err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				info, _ := d.Info()
				item := map[string]any{
					"path":  path,
					"name":  d.Name(),
					"isDir": d.IsDir(),
				}
				if info != nil {
					item["size"] = float64(info.Size())
					item["modTime"] = info.ModTime().Unix()
				}
				files = append(files, item)
				return nil
			})

			if err != nil {
				return nil, err
			}
			return files, nil
		})).
		AddStaticMethod("workingDir", stringType, []ast.Parameter{}, Func(func(_ *Env, _ []any) (any, error) {
			wd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			return wd, nil
		})).
		AddStaticMethod("changeDir", boolType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			if err := os.Chdir(path); err != nil {
				return false, err
			}
			return true, nil
		})).

		// Path operations
		AddStaticMethod("pathJoin", stringType, []ast.Parameter{{Name: "parts", Type: nil, IsVariadic: true}}, Func(func(_ *Env, args []any) (any, error) {
			if len(args) == 0 {
				return "", nil
			}
			parts := make([]string, len(args))
			for i, arg := range args {
				parts[i] = utils.ToString(arg)
			}
			return filepath.Join(parts...), nil
		})).
		AddStaticMethod("pathSplit", mapType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			dir, file := filepath.Split(path)
			return map[string]any{
				"dir":  dir,
				"file": file,
			}, nil
		})).
		AddStaticMethod("pathExt", stringType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Ext(path), nil
		})).
		AddStaticMethod("pathBase", stringType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Base(path), nil
		})).
		AddStaticMethod("pathDir", stringType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Dir(path), nil
		})).
		AddStaticMethod("pathAbs", stringType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			abs, err := filepath.Abs(path)
			if err != nil {
				return nil, err
			}
			return abs, nil
		})).

		// Stream operations
		AddStaticMethod("readLines", arrayType, []ast.Parameter{{Name: "path", Type: stringType}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			var lines []any
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				return nil, err
			}

			return lines, nil
		})).
		AddStaticMethod("writeLines", boolType, []ast.Parameter{
			{Name: "path", Type: stringType},
			{Name: "lines", Type: stringType},
		}, Func(func(e *Env, args []any) (any, error) {
			path := utils.ToString(args[0])

			lines, ok := args[1].([]any)
			if !ok {
				return nil, ThrowTypeError(e, "array", args[1])
			}

			file, err := os.Create(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			for _, line := range lines {
				if _, err := writer.WriteString(utils.ToString(line) + "\n"); err != nil {
					return nil, err
				}
			}

			if err := writer.Flush(); err != nil {
				return nil, err
			}

			return true, nil
		})).

		// Text processing
		AddStaticMethod("grep", arrayType, []ast.Parameter{
			{Name: "pattern", Type: stringType},
			{Name: "text", Type: stringType},
		}, Func(func(_ *Env, args []any) (any, error) {
			pattern := utils.ToString(args[0])
			text := utils.ToString(args[1])

			lines := strings.Split(text, "\n")
			var matches []any

			for i, line := range lines {
				if strings.Contains(line, pattern) {
					match := map[string]any{
						"line":  float64(i + 1),
						"text":  line,
						"match": pattern,
					}
					matches = append(matches, match)
				}
			}

			return matches, nil
		}))

	_, err = ioClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
	return nil
}

// getEncoding returns the appropriate encoding based on the encoding name
func getEncoding(name string) encoding.Encoding {
	name = strings.ToLower(strings.ReplaceAll(name, "-", ""))
	name = strings.ReplaceAll(name, "_", "")

	switch name {
	// UTF encodings
	case "utf8":
		return nil // UTF-8 is default, no transformation needed
	case "utf16", "utf16be":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "utf16le":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	// Latin/Western European
	case "latin1", "iso88591":
		return charmap.ISO8859_1
	case "latin2", "iso88592":
		return charmap.ISO8859_2
	case "latin3", "iso88593":
		return charmap.ISO8859_3
	case "latin4", "iso88594":
		return charmap.ISO8859_4
	case "latin5", "iso88595":
		return charmap.ISO8859_5
	case "latin6", "iso88596":
		return charmap.ISO8859_6
	case "latin7", "iso88597":
		return charmap.ISO8859_7
	case "latin8", "iso88598":
		return charmap.ISO8859_8
	case "latin9", "iso88599":
		return charmap.ISO8859_9
	case "latin10", "iso885910":
		return charmap.ISO8859_10

	// Windows code pages
	case "windows1250", "cp1250":
		return charmap.Windows1250
	case "windows1251", "cp1251":
		return charmap.Windows1251
	case "windows1252", "cp1252":
		return charmap.Windows1252
	case "windows1253", "cp1253":
		return charmap.Windows1253
	case "windows1254", "cp1254":
		return charmap.Windows1254
	case "windows1255", "cp1255":
		return charmap.Windows1255
	case "windows1256", "cp1256":
		return charmap.Windows1256
	case "windows1257", "cp1257":
		return charmap.Windows1257
	case "windows1258", "cp1258":
		return charmap.Windows1258

	// Japanese
	case "shiftjis", "sjis":
		return japanese.ShiftJIS
	case "eucjp":
		return japanese.EUCJP
	case "iso2022jp":
		return japanese.ISO2022JP

	// Korean
	case "euckr":
		return korean.EUCKR

	// Chinese
	case "gbk":
		return simplifiedchinese.GBK
	case "gb18030":
		return simplifiedchinese.GB18030
	case "big5":
		return traditionalchinese.Big5

	// IBM code pages
	case "ibm437", "cp437":
		return charmap.CodePage437
	case "ibm850", "cp850":
		return charmap.CodePage850
	case "ibm852", "cp852":
		return charmap.CodePage852
	case "ibm855", "cp855":
		return charmap.CodePage855
	case "ibm858", "cp858":
		return charmap.CodePage858
	case "ibm860", "cp860":
		return charmap.CodePage860
	case "ibm862", "cp862":
		return charmap.CodePage862
	case "ibm863", "cp863":
		return charmap.CodePage863
	case "ibm865", "cp865":
		return charmap.CodePage865
	case "ibm866", "cp866":
		return charmap.CodePage866

	// Mac encodings
	case "macintosh", "macroman":
		return charmap.Macintosh
	case "maccyrillic":
		return charmap.MacintoshCyrillic

	// KOI8
	case "koi8r":
		return charmap.KOI8R
	case "koi8u":
		return charmap.KOI8U

	default:
		return nil // Unknown encoding, fallback to UTF-8
	}
}

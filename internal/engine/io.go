package engine

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/engine/utils"
)

func InstallIOModule(env *Env, opts Options) {
	ioClass := NewClassBuilder("IO").
		// File operations
		AddStaticMethod("readFile", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return string(data), nil
		})).
		AddStaticMethod("writeFile", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("")},
			{Name: "content", Type: ast.TypeFromString("")},
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
		AddStaticMethod("appendFile", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("")},
			{Name: "content", Type: ast.TypeFromString("")},
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
		AddStaticMethod("copyFile", &ast.Type{Name: "float", IsBuiltin: true}, []ast.Parameter{
			{Name: "source", Type: ast.TypeFromString("")},
			{Name: "destination", Type: ast.TypeFromString("")},
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
		AddStaticMethod("deleteFile", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			if err := os.Remove(path); err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("moveFile", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "source", Type: ast.TypeFromString("")},
			{Name: "destination", Type: ast.TypeFromString("")},
		}, Func(func(_ *Env, args []any) (any, error) {
			src := utils.ToString(args[0])
			dst := utils.ToString(args[1])

			if err := os.Rename(src, dst); err != nil {
				return false, err
			}
			return true, nil
		})).
		AddStaticMethod("fileExists", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			_, err := os.Stat(path)
			return err == nil, nil
		})).
		AddStaticMethod("fileInfo", &ast.Type{Name: "Map", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
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

		// Directory operations
		AddStaticMethod("createDir", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("")},
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
		AddStaticMethod("removeDir", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("")},
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
		AddStaticMethod("listDir", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
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
		AddStaticMethod("walkDir", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
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
		AddStaticMethod("workingDir", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{}, Func(func(_ *Env, _ []any) (any, error) {
			wd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			return wd, nil
		})).
		AddStaticMethod("changeDir", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			if err := os.Chdir(path); err != nil {
				return false, err
			}
			return true, nil
		})).

		// Path operations
		AddStaticMethod("pathJoin", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "parts", Type: nil, IsVariadic: true}}, Func(func(_ *Env, args []any) (any, error) {
			if len(args) == 0 {
				return "", nil
			}
			parts := make([]string, len(args))
			for i, arg := range args {
				parts[i] = utils.ToString(arg)
			}
			return filepath.Join(parts...), nil
		})).
		AddStaticMethod("pathSplit", &ast.Type{Name: "Map", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			dir, file := filepath.Split(path)
			return map[string]any{
				"dir":  dir,
				"file": file,
			}, nil
		})).
		AddStaticMethod("pathExt", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Ext(path), nil
		})).
		AddStaticMethod("pathBase", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Base(path), nil
		})).
		AddStaticMethod("pathDir", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			return filepath.Dir(path), nil
		})).
		AddStaticMethod("pathAbs", &ast.Type{Name: "string", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
			path := utils.ToString(args[0])
			abs, err := filepath.Abs(path)
			if err != nil {
				return nil, err
			}
			return abs, nil
		})).

		// Stream operations
		AddStaticMethod("readLines", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{{Name: "path", Type: ast.TypeFromString("")}}, Func(func(_ *Env, args []any) (any, error) {
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
		AddStaticMethod("writeLines", &ast.Type{Name: "bool", IsBuiltin: true}, []ast.Parameter{
			{Name: "path", Type: ast.TypeFromString("")},
			{Name: "lines", Type: ast.TypeFromString("")},
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
		AddStaticMethod("grep", &ast.Type{Name: "Array", IsBuiltin: true}, []ast.Parameter{
			{Name: "pattern", Type: ast.TypeFromString("")},
			{Name: "text", Type: ast.TypeFromString("")},
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

	_, err := ioClass.BuildStatic(env)
	if err != nil {
		panic(err)
	}
}

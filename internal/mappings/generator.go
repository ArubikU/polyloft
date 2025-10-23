package mappings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Symbol represents a symbol (class, function, variable, etc.) in a Polyloft file
type Symbol struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "class", "function", "variable", "constant"
	ReturnType  string            `json:"returnType,omitempty"`
	Parameters  []Parameter       `json:"parameters,omitempty"`
	Fields      []Field           `json:"fields,omitempty"`
	Methods     []Symbol          `json:"methods,omitempty"`
	Description string            `json:"description,omitempty"`
	File        string            `json:"file"`
	Line        int               `json:"line"`
	Parent      string            `json:"parent,omitempty"`
	Implements  []string          `json:"implements,omitempty"`
	Modifiers   []string          `json:"modifiers,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Parameter represents a function parameter
type Parameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
	Variadic bool   `json:"variadic,omitempty"`
}

// Field represents a class field
type Field struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Modifiers  []string `json:"modifiers,omitempty"`
	Visibility string   `json:"visibility,omitempty"`
}

// PackageMapping represents all symbols in a package/module
type PackageMapping struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Description string            `json:"description,omitempty"`
	Version     string            `json:"version"`
	Symbols     []Symbol          `json:"symbols"`
	Imports     []string          `json:"imports"`
	Exports     []string          `json:"exports"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Mappings represents the complete mapping file
type Mappings struct {
	Version  string                    `json:"version"`
	Packages map[string]PackageMapping `json:"packages"`
}

// Generator generates mappings.json from Polyloft source files
type Generator struct {
	rootPath string
	libsPath string
}

// NewGenerator creates a new mappings generator
func NewGenerator(rootPath string) *Generator {
	return &Generator{
		rootPath: rootPath,
		libsPath: filepath.Join(rootPath, "libs"),
	}
}

// Generate creates a mappings.json file by scanning all .pf files
func (g *Generator) Generate(outputPath string) error {
	mappings := Mappings{
		Version:  "1.0.0",
		Packages: make(map[string]PackageMapping),
	}

	// Scan libs directory for packages
	if _, err := os.Stat(g.libsPath); err == nil {
		err = filepath.Walk(g.libsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".pf") {
				relPath, _ := filepath.Rel(g.libsPath, filepath.Dir(path))
				packageName := strings.ReplaceAll(relPath, string(filepath.Separator), ".")

				// Parse the file and extract symbols
				symbols, imports, exports, err := g.parseFile(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", path, err)
					return nil
				}

				// Get or create package mapping
				pkgMapping, exists := mappings.Packages[packageName]
				if !exists {
					pkgMapping = PackageMapping{
						Name:     packageName,
						Path:     relPath,
						Version:  "1.0.0",
						Symbols:  []Symbol{},
						Imports:  []string{},
						Exports:  []string{},
						Metadata: make(map[string]string),
					}
				}

				// Add symbols
				pkgMapping.Symbols = append(pkgMapping.Symbols, symbols...)

				// Merge imports and exports
				pkgMapping.Imports = append(pkgMapping.Imports, imports...)
				pkgMapping.Exports = append(pkgMapping.Exports, exports...)

				// Update package
				mappings.Packages[packageName] = pkgMapping
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk libs directory: %w", err)
		}
	}

	// Write mappings to file
	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mappings: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write mappings file: %w", err)
	}

	fmt.Printf("Generated mappings.json with %d packages\n", len(mappings.Packages))
	return nil
}

// parseFile parses a Polyloft file and extracts symbols
func (g *Generator) parseFile(filePath string) ([]Symbol, []string, []string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, nil, err
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	var symbols []Symbol
	var imports []string
	var exports []string

	// Regular expressions for parsing
	classRegex := regexp.MustCompile(`^\s*(public\s+|private\s+)?class\s+([A-Z][a-zA-Z0-9_]*)(?:\s*<\s*([A-Z][a-zA-Z0-9_]*))?(?:\s+implements\s+([^\n{]+))?`)
	functionRegex := regexp.MustCompile(`^\s*(public\s+|private\s+|protected\s+|static\s+)*def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)(?:\s*->\s*([a-zA-Z][a-zA-Z0-9_]*))?`)
	varRegex := regexp.MustCompile(`^\s*(public\s+|private\s+|protected\s+|static\s+)?var\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*([a-zA-Z][a-zA-Z0-9_]*)`)
	importRegex := regexp.MustCompile(`^\s*import\s+([a-zA-Z._/]+)`)

	currentClass := ""
	classSymbol := (*Symbol)(nil)

	for lineNum, line := range lines {
		// Parse imports
		if match := importRegex.FindStringSubmatch(line); match != nil {
			imports = append(imports, match[1])
		}

		// Parse class declarations
		if match := classRegex.FindStringSubmatch(line); match != nil {
			modifiers := []string{}
			if match[1] != "" {
				modifiers = append(modifiers, strings.TrimSpace(match[1]))
			}

			className := match[2]
			parent := match[3]
			implementsStr := match[4]

			var implements []string
			if implementsStr != "" {
				for _, iface := range strings.Split(implementsStr, ",") {
					implements = append(implements, strings.TrimSpace(iface))
				}
			}

			symbol := Symbol{
				Name:       className,
				Type:       "class",
				File:       filePath,
				Line:       lineNum + 1,
				Parent:     parent,
				Implements: implements,
				Modifiers:  modifiers,
				Methods:    []Symbol{},
				Fields:     []Field{},
			}

			currentClass = className
			classSymbol = &symbol
			exports = append(exports, className)
		}

		// Parse function/method declarations
		if match := functionRegex.FindStringSubmatch(line); match != nil {
			modifiersStr := strings.TrimSpace(match[1])
			funcName := match[2]
			paramsStr := match[3]
			returnType := match[4]

			if returnType == "" {
				returnType = "Void"
			}

			modifiers := []string{}
			if modifiersStr != "" {
				for _, mod := range strings.Fields(modifiersStr) {
					modifiers = append(modifiers, mod)
				}
			}

			// Parse parameters
			params := []Parameter{}
			if paramsStr != "" {
				for _, param := range strings.Split(paramsStr, ",") {
					param = strings.TrimSpace(param)
					if param == "" {
						continue
					}

					// Parse parameter (name: Type)
					parts := strings.Split(param, ":")
					if len(parts) == 2 {
						paramName := strings.TrimSpace(parts[0])
						paramType := strings.TrimSpace(parts[1])

						params = append(params, Parameter{
							Name: paramName,
							Type: paramType,
						})
					}
				}
			}

			symbol := Symbol{
				Name:       funcName,
				Type:       "function",
				ReturnType: returnType,
				Parameters: params,
				File:       filePath,
				Line:       lineNum + 1,
				Modifiers:  modifiers,
			}

			// If inside a class, add as method
			if currentClass != "" && classSymbol != nil {
				classSymbol.Methods = append(classSymbol.Methods, symbol)
			} else {
				symbols = append(symbols, symbol)
				exports = append(exports, funcName)
			}
		}

		// Parse variable declarations
		if match := varRegex.FindStringSubmatch(line); match != nil {
			modifiersStr := strings.TrimSpace(match[1])
			varName := match[2]
			varType := match[3]

			modifiers := []string{}
			if modifiersStr != "" {
				for _, mod := range strings.Fields(modifiersStr) {
					modifiers = append(modifiers, mod)
				}
			}

			// If inside a class, add as field
			if currentClass != "" && classSymbol != nil {
				visibility := "private"
				if strings.Contains(modifiersStr, "public") {
					visibility = "public"
				} else if strings.Contains(modifiersStr, "protected") {
					visibility = "protected"
				}

				field := Field{
					Name:       varName,
					Type:       varType,
					Modifiers:  modifiers,
					Visibility: visibility,
				}
				classSymbol.Fields = append(classSymbol.Fields, field)
			}
		}

		// End of class
		if strings.TrimSpace(line) == "end" && currentClass != "" {
			if classSymbol != nil {
				symbols = append(symbols, *classSymbol)
			}
			currentClass = ""
			classSymbol = nil
		}
	}

	// If we still have an open class (no 'end' found), add it anyway
	if classSymbol != nil {
		symbols = append(symbols, *classSymbol)
	}

	return symbols, imports, exports, nil
}

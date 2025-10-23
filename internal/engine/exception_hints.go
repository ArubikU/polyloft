package engine

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// ExceptionHint provides contextual suggestions for errors
type ExceptionHint struct {
	Message     string
	Suggestions []string
	HintType    string // "typo", "language_conversion", "general"
}

// HintProvider generates hints based on error context
type HintProvider struct {
	env *Env
}

// NewHintProvider creates a new hint provider
func NewHintProvider(env *Env) *HintProvider {
	return &HintProvider{env: env}
}

// getSourceLine reads a specific line from a source file
func getSourceLine(filepath string, lineNum int) string {
	if filepath == "" || lineNum <= 0 {
		return ""
	}

	// Try to read the file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		return lines[lineNum-1]
	}

	return ""
}

// LanguagePattern represents a pattern from another language
type LanguagePattern struct {
	Pattern     string
	Language    string
	HyLangEquiv string
	Description string
}

// getLanguagePatterns returns common patterns from other languages
func getLanguagePatterns() []LanguagePattern {
	return []LanguagePattern{
		// JavaScript/TypeScript
		{"console.log", "JavaScript", "println", "In HyLang, use: println(...)"},
		{"console.error", "JavaScript", "Sys.println", "In HyLang, use: Sys.println(...)"},
		{"typeof", "JavaScript", "Sys.type", "In HyLang, use: Sys.type(value)"},
		{"function", "JavaScript", "def", "In HyLang, functions are defined as: def name(...) { ... } or def name(...) = expression"},
		{"const ", "JavaScript", "let", "In HyLang, use: let name = value (immutable) or var name = value (mutable)"},
		{"=>", "JavaScript", "=>", "Lambda syntax is similar: (params) => expression"},

		// Java
		{"System.out.println", "Java", "println", "In HyLang, use: println(...)"},
		{"System.out.print", "Java", "print", "In HyLang, use: print(...)"},
		{"public static void main", "Java", "def main", "In HyLang, use: def main(...): ... end"},
		{"new ", "Java", "", "In HyLang, call constructors directly: ClassName(...)"},

		// Python
		{"def ", "Python", "def", "In HyLang, use: def name(...): ... end or def name(...) = expression"},
		{"print(", "Python", "println", "Both work similarly, but HyLang also has: Sys.println(...)"},
		{"__init__", "Python", "ClassName", "In HyLang, constructors use the class name: ClassName(...): ... end"},
		{"self", "Python", "this", "In HyLang, use 'this' to refer to the current instance"},

		// C#
		{"Console.WriteLine", "C#", "println", "In HyLang, use: println(...)"},
		{"Console.Write", "C#", "print", "In HyLang, use: print(...)"},

		// Common function names
		{"log", "common", "println", "In HyLang, use: println(...) or Sys.println(...)"},
		{"printf", "C/C++", "println", "In HyLang, use: println(...) for formatted output"},
		{"echo", "PHP", "println", "In HyLang, use: println(...)"},
	}
}

// detectFullPattern analyzes the full line of code to find complete pattern matches
func detectFullPattern(sourceLine string, name string) *ExceptionHint {
	if sourceLine == "" {
		return nil
	}

	line := strings.TrimSpace(sourceLine)
	lineLower := strings.ToLower(line)

	// Define complete pattern mappings with more precise detection
	patterns := getLanguagePatterns()

	for _, p := range patterns {
		// Get fuzzy distance
		distance := fuzzy.RankMatch(p.Pattern, lineLower)
		if distance != -1 && distance <= 2 { // <=2 permite pequeñas diferencias
			return &ExceptionHint{
				Message:     p.Description,
				Suggestions: []string{p.Description},
				HintType:    "language_conversion",
			}
		}
	}

	return nil
}

// searchFileForPattern searches the entire file for patterns containing the identifier
func searchFileForPattern(filepath string, name string) *ExceptionHint {
	if filepath == "" || name == "" {
		return nil
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	nameLower := strings.ToLower(name)

	for _, line := range lines {
		lineLower := strings.ToLower(line)
		if strings.Contains(lineLower, nameLower+".") {
			if hint := detectFullPattern(line, name); hint != nil {
				return hint
			}
		}
	}

	return nil
}

// getContextualHint provides context-aware hints for common mistakes
// This is now a fallback when we don't have source context
func getContextualHint(name string) *ExceptionHint {
	nameLower := strings.ToLower(name)

	// Only provide hints for very specific standalone identifiers
	// that clearly indicate another language
	standaloneHints := map[string]*ExceptionHint{
		"typeof": {
			Message:     "In HyLang, use:",
			Suggestions: []string{"Sys.type(value)"},
			HintType:    "language_conversion",
		},
		"printf": {
			Message:     "In HyLang, use:",
			Suggestions: []string{"println(...)"},
			HintType:    "language_conversion",
		},
		"echo": {
			Message:     "In HyLang, use:",
			Suggestions: []string{"println(...)"},
			HintType:    "language_conversion",
		},
	}

	if hint, ok := standaloneHints[nameLower]; ok {
		return hint
	}

	return nil
}

// getKeywordHint provides hints for typos in keywords, modifiers, and annotations
func getKeywordHint(name string) *ExceptionHint {
	nameLower := strings.ToLower(name)

	// Common HyLang keywords, modifiers, and annotations
	keywords := []string{
		// Access modifiers
		"public", "private", "protected",
		// Other modifiers
		"static", "final", "const", "abstract", "sealed",
		// Keywords
		"class", "interface", "enum", "record",
		"def", "let", "var",
		"if", "elif", "else",
		"for", "in", "loop", "do", "end",
		"break", "continue", "return",
		"try", "catch", "finally", "throw",
		"import", "from",
		"this", "super",
		"instanceof", "new",
		"true", "false", "nil",
	}

	// Find close matches (Levenshtein distance <= 2)
	var suggestions []string
	for _, keyword := range keywords {
		if levenshteinDistance(nameLower, keyword) <= 2 {
			suggestions = append(suggestions, keyword)
		}
	}

	if len(suggestions) == 0 {
		return nil
	}

	// Limit to top 3
	suggestions = uniqueAndLimit(suggestions, 3)

	return &ExceptionHint{
		Message:     "Did you mean:",
		Suggestions: suggestions,
		HintType:    "keyword_typo",
	}
}

// getVariedHintMessage returns a varied hint message based on context
func getVariedHintMessage(hintType string, isExactMatch bool) string {
	if isExactMatch {
		return "In HyLang this would be:"
	}

	switch hintType {
	case "language_conversion":
		return "In HyLang, try:"
	case "typo":
		messages := []string{
			"Maybe you meant:",
			"Perhaps you were referring to:",
			"Possible match:",
		}
		// Simple rotation based on length for variety
		return messages[len(hintType)%len(messages)]
	case "enum":
		return "Perhaps you were referring to:"
	case "attribute":
		return "Maybe you meant:"
	default:
		return "Possible match:"
	}
}

// GetHintForUndefinedName provides suggestions for undefined identifiers
func (hp *HintProvider) GetHintForUndefinedName(name string) *ExceptionHint {
	return hp.GetHintForUndefinedNameWithContext(name, "", 0)
}

// GetHintForUndefinedNameWithEnv provides suggestions using the env's code context
func (hp *HintProvider) GetHintForUndefinedNameWithEnv(name string) *ExceptionHint {
	if hp.env == nil {
		return hp.GetHintForUndefinedName(name)
	}

	// Try to detect patterns from the code context stored in env
	codeContext := hp.env.GetCodeContext()
	if len(codeContext) > 0 {
		// Check the current line (last in context)
		currentLine := codeContext[len(codeContext)-1]
		if currentLine != "" {
			if fullPatternHint := detectFullPattern(currentLine, name); fullPatternHint != nil {
				return fullPatternHint
			}
		}
	}

	// Try to search through all source lines if available
	sourceLines := hp.env.GetSourceLines()
	if len(sourceLines) > 0 {
		nameLower := strings.ToLower(name)
		// Look for patterns that contain this identifier
		patterns := getLanguagePatterns()
		for _, line := range sourceLines {
			lineLower := strings.ToLower(line)
			// Check if this line contains the identifier being used
			if strings.Contains(lineLower, nameLower) {
				// Check against all patterns - must match the FULL pattern, not just prefix
				for _, p := range patterns {
					patternLower := strings.ToLower(p.Pattern)
					// Only match if the FULL pattern is present in the line
					// This ensures "console" alone won't match, but "console.log" will
					if strings.Contains(lineLower, patternLower) {
						return &ExceptionHint{
							Message:     p.Description,
							Suggestions: []string{p.HyLangEquiv},
							HintType:    "language_conversion",
						}
					}
				}
			}
		}
	}

	// Fall back to file-based search if we have a filename
	fileName := hp.env.GetFileName()
	currentLine := hp.env.GetCurrentLine()
	if fileName != "" && currentLine > 0 {
		return hp.GetHintForUndefinedNameWithContext(name, fileName, currentLine)
	}

	// Fall back to basic hint without context
	return hp.GetHintForUndefinedNameWithContext(name, "", 0)
}

// GetHintForUndefinedNameWithContext provides suggestions with source file context
func (hp *HintProvider) GetHintForUndefinedNameWithContext(name string, filepath string, line int) *ExceptionHint {
	suggestions := []string{}
	hintType := "typo"

	// First priority: Try to read the actual source line and detect full patterns
	if filepath != "" && line > 0 {
		sourceLine := getSourceLine(filepath, line)
		if sourceLine != "" {
			if fullPatternHint := detectFullPattern(sourceLine, name); fullPatternHint != nil {
				return fullPatternHint
			}
		}
	}

	// If we don't have line number but have file path, search the entire file
	if filepath != "" && line == 0 {
		if filePatternHint := searchFileForPattern(filepath, name); filePatternHint != nil {
			return filePatternHint
		}
	}

	// Second priority: Check for contextual hints (common language patterns)
	// but only for very specific standalone identifiers
	if contextHint := getContextualHint(name); contextHint != nil {
		return contextHint
	}

	// Third priority: Check for keyword typos (modifiers, annotations, etc.)
	if keywordHint := getKeywordHint(name); keywordHint != nil {
		return keywordHint
	}

	// Check environment variables for similar names
	if hp.env != nil {
		envSuggestions := hp.findSimilarNamesInEnv(name)
		suggestions = append(suggestions, envSuggestions...)
	}

	// Check common built-in functions
	builtins := []string{
		"println", "print", "Sys.println", "Sys.print",
		"Array", "Map", "String", "Int", "Float", "Bool",
		"Math", "Sys", "IO", "Net", "Crypto",
	}

	for _, builtin := range builtins {
		// Only suggest if:
		// 1. Very similar (Levenshtein distance <= 2)
		// 2. The name is a substantial part of the builtin (at least 3 chars and is substring)
		if levenshteinDistance(name, builtin) <= 2 {
			suggestions = append(suggestions, builtin)
		} else if len(name) >= 3 && strings.Contains(strings.ToLower(builtin), strings.ToLower(name)) {
			suggestions = append(suggestions, builtin)
		}
	}

	if len(suggestions) == 0 {
		return nil
	}

	// Remove duplicates and limit to top 3
	suggestions = uniqueAndLimit(suggestions, 3)

	return &ExceptionHint{
		Message:     getVariedHintMessage(hintType, false),
		Suggestions: suggestions,
		HintType:    hintType,
	}
}

// GetHintForAttribute provides suggestions for missing attributes
func (hp *HintProvider) GetHintForAttribute(attrName string, typeName string) *ExceptionHint {
	suggestions := []string{}

	// Extract class name from typeName (it may be formatted like "'ClassName' instance")
	className := typeName
	if strings.Contains(typeName, "'") {
		parts := strings.Split(typeName, "'")
		if len(parts) >= 2 {
			className = parts[1]
		}
	}

	// Try to find the class/enum definition and suggest similar names
	// Note: For suggestions, we check all packages since we're trying to be helpful
	if classDef, ok := builtinClasses[className]; ok {
		// Check fields
		for fieldName := range classDef.Fields {
			if levenshteinDistance(attrName, fieldName) <= 2 {
				suggestions = append(suggestions, fieldName)
			}
		}
		// Check methods
		for methodName := range classDef.Methods {
			if levenshteinDistance(attrName, methodName) <= 2 {
				suggestions = append(suggestions, methodName)
			}
		}
	} else {
		// Check all packages for class definitions
		for _, packageClasses := range classRegistry {
			if classDef, ok := packageClasses[className]; ok {
				// Check fields
				for fieldName := range classDef.Fields {
					if levenshteinDistance(attrName, fieldName) <= 2 {
						suggestions = append(suggestions, fieldName)
					}
				}
				// Check methods
				for methodName := range classDef.Methods {
					if levenshteinDistance(attrName, methodName) <= 2 {
						suggestions = append(suggestions, methodName)
					}
				}
				break
			}
		}
	}

	// Check enum values
	for enumName, enumDef := range enumRegistry {
		if strings.Contains(className, enumName) {
			for valueName := range enumDef.Values {
				if levenshteinDistance(attrName, valueName) <= 2 {
					suggestions = append(suggestions, fmt.Sprintf("%s.%s", enumName, valueName))
				}
			}
		}
	}

	if len(suggestions) == 0 {
		return nil
	}

	suggestions = uniqueAndLimit(suggestions, 3)

	return &ExceptionHint{
		Message:     getVariedHintMessage("attribute", false),
		Suggestions: suggestions,
		HintType:    "typo",
	}
}

// GetHintForEnumValue provides suggestions for enum value typos
func (hp *HintProvider) GetHintForEnumValue(enumName string, valueName string) *ExceptionHint {
	enumDef, ok := enumRegistry[enumName]
	if !ok {
		return nil
	}

	suggestions := []string{}
	for existingValue := range enumDef.Values {
		if levenshteinDistance(valueName, existingValue) <= 2 {
			suggestions = append(suggestions, fmt.Sprintf("%s.%s", enumName, existingValue))
		}
	}

	if len(suggestions) == 0 {
		// Show all available values
		for existingValue := range enumDef.Values {
			suggestions = append(suggestions, fmt.Sprintf("%s.%s", enumName, existingValue))
		}
		suggestions = uniqueAndLimit(suggestions, 5)

		return &ExceptionHint{
			Message:     "Available values:",
			Suggestions: suggestions,
			HintType:    "enum",
		}
	}

	suggestions = uniqueAndLimit(suggestions, 3)

	return &ExceptionHint{
		Message:     getVariedHintMessage("enum", false),
		Suggestions: suggestions,
		HintType:    "typo",
	}
}

// findSimilarNamesInEnv searches for similar variable names in the environment
func (hp *HintProvider) findSimilarNamesInEnv(name string) []string {
	suggestions := []string{}

	for cur := hp.env; cur != nil; cur = cur.Parent {
		for varName := range cur.Vars {
			if levenshteinDistance(name, varName) <= 2 {
				suggestions = append(suggestions, varName)
			}
		}
	}

	return suggestions
}

// levenshteinDistance computes the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	s1Lower := strings.ToLower(s1)
	s2Lower := strings.ToLower(s2)

	if len(s1Lower) == 0 {
		return len(s2Lower)
	}
	if len(s2Lower) == 0 {
		return len(s1Lower)
	}

	// Create a matrix
	matrix := make([][]int, len(s1Lower)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2Lower)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len(s1Lower); i++ {
		for j := 1; j <= len(s2Lower); j++ {
			cost := 1
			if s1Lower[i-1] == s2Lower[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1Lower)][len(s2Lower)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// uniqueAndLimit removes duplicates and limits the slice
func uniqueAndLimit(items []string, limit int) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
			if len(result) >= limit {
				break
			}
		}
	}

	return result
}

// GetGeneralHint provides context-specific hints based on error type
func GetGeneralHint(errorType string, context string) *ExceptionHint {
	hints := map[string]map[string][]string{
		"NameError": {
			"log": {"Use 'println' or 'sys.println' for output", "Use 'print' or 'sys.print' for output without newline"},
		},
		"TypeError": {
			"callable": {"Make sure the object is a function", "Check if you're calling the right variable"},
		},
		"ArityError": {
			"": {"Check the function signature", "Verify the number of arguments matches the parameters"},
		},
	}

	if contextHints, ok := hints[errorType]; ok {
		for key, suggestions := range contextHints {
			if key == "" || strings.Contains(strings.ToLower(context), strings.ToLower(key)) {
				return &ExceptionHint{
					Message:     "Hint:",
					Suggestions: suggestions,
				}
			}
		}
	}

	return nil
}

// SortSuggestionsByRelevance sorts suggestions by their similarity to the target
func SortSuggestionsByRelevance(suggestions []string, target string) []string {
	type scoredItem struct {
		suggestion string
		distance   int
	}

	scoredItems := make([]scoredItem, len(suggestions))
	for i, s := range suggestions {
		scoredItems[i] = scoredItem{
			suggestion: s,
			distance:   levenshteinDistance(target, s),
		}
	}

	sort.Slice(scoredItems, func(i, j int) bool {
		return scoredItems[i].distance < scoredItems[j].distance
	})

	result := make([]string, len(scoredItems))
	for i, s := range scoredItems {
		result[i] = s.suggestion
	}

	return result
}

package engine

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorWhite  = "\033[97m"
	
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldCyan   = "\033[1;36m"
)

// formatHintSuggestions formats hint suggestions with or without bullets
func formatHintSuggestions(hint *ExceptionHint, withColor bool) string {
	if hint == nil || len(hint.Suggestions) == 0 {
		return ""
	}
	
	var builder strings.Builder
	
	if withColor {
		builder.WriteString(fmt.Sprintf("\n%s%s%s ", 
			ColorBoldYellow, hint.Message, ColorReset))
	} else {
		builder.WriteString(fmt.Sprintf("\n%s ", hint.Message))
	}
	
	// If only one suggestion, show it inline without bullet
	if len(hint.Suggestions) == 1 {
		if withColor {
			builder.WriteString(fmt.Sprintf("%s%s%s\n", 
				ColorGreen, hint.Suggestions[0], ColorReset))
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", hint.Suggestions[0]))
		}
	} else {
		// Multiple suggestions, show as bulleted list
		builder.WriteString("\n")
		for _, suggestion := range hint.Suggestions {
			if withColor {
				builder.WriteString(fmt.Sprintf("  %s-%s %s%s%s\n", 
					ColorYellow, ColorReset, ColorGreen, suggestion, ColorReset))
			} else {
				builder.WriteString(fmt.Sprintf("  - %s\n", suggestion))
			}
		}
	}
	
	return builder.String()
}


// FormatError formats a HyException with colors and hints
func FormatError(err error) string {
	hyErr, ok := err.(*HyException)
	if !ok {
		// Not a HyException, format as regular error
		return fmt.Sprintf("%s%s%s", ColorBoldRed, err.Error(), ColorReset)
	}
	
	var builder strings.Builder
	
	// File:Line location (if available)
	if hyErr.File != "" {
		builder.WriteString(fmt.Sprintf("%s%s:%d%s", ColorBoldCyan, hyErr.File, hyErr.Line, ColorReset))
		if hyErr.Column > 0 {
			builder.WriteString(fmt.Sprintf(":%d", hyErr.Column))
		}
		builder.WriteString("\n")
	}
	
	// Error type and message
	errorType := hyErr.Type
	if errorType == "" {
		errorType = "Error"
	}
	
	builder.WriteString(fmt.Sprintf("%s%s%s: %s\n", 
		ColorBoldRed, errorType, ColorReset, hyErr.Message))
	
	// Hints (if available)
	if hyErr.Hint != nil && len(hyErr.Hint.Suggestions) > 0 {
		builder.WriteString(formatHintSuggestions(hyErr.Hint, true))
	}
	
	// Stack trace (if available and not empty)
	if len(hyErr.StackTrace) > 0 {
		builder.WriteString(fmt.Sprintf("\n%sStack trace:%s\n", ColorGray, ColorReset))
		for i, trace := range hyErr.StackTrace {
			builder.WriteString(fmt.Sprintf("  %s%d. %s%s\n", ColorGray, i+1, trace, ColorReset))
		}
	}
	
	return builder.String()
}

// FormatErrorWithContext formats a HyException with code context if available
func FormatErrorWithContext(err error, codeContext []string, currentLine int) string {
	hyErr, ok := err.(*HyException)
	if !ok {
		// Not a HyException, format as regular error
		return fmt.Sprintf("%s%s%s", ColorBoldRed, err.Error(), ColorReset)
	}
	
	var builder strings.Builder
	
	// File:Line location (if available)
	if hyErr.File != "" {
		builder.WriteString(fmt.Sprintf("%s%s:%d%s", ColorBoldCyan, hyErr.File, hyErr.Line, ColorReset))
		if hyErr.Column > 0 {
			builder.WriteString(fmt.Sprintf(":%d", hyErr.Column))
		}
		builder.WriteString("\n")
	}
	
	// Code context (if available)
	if len(codeContext) > 0 && currentLine > 0 {
		builder.WriteString(fmt.Sprintf("\n%sCode context:%s\n", ColorGray, ColorReset))
		
		// Show lines with line numbers
		startLine := currentLine - len(codeContext) + 1
		for i, line := range codeContext {
			lineNum := startLine + i
			isCurrentLine := lineNum == currentLine
			
			if isCurrentLine {
				// Highlight the current line
				builder.WriteString(fmt.Sprintf("  %s>%s %s%-4d%s %s%s%s\n", 
					ColorBoldRed, ColorReset, ColorGray, lineNum, ColorReset, 
					ColorWhite, line, ColorReset))
			} else {
				builder.WriteString(fmt.Sprintf("   %s%-4d%s %s\n", 
					ColorGray, lineNum, ColorReset, line))
			}
		}
		builder.WriteString("\n")
	}
	
	// Error type and message
	errorType := hyErr.Type
	if errorType == "" {
		errorType = "Error"
	}
	
	builder.WriteString(fmt.Sprintf("%s%s%s: %s\n", 
		ColorBoldRed, errorType, ColorReset, hyErr.Message))
	
	// Hints (if available)
	if hyErr.Hint != nil && len(hyErr.Hint.Suggestions) > 0 {
		builder.WriteString(formatHintSuggestions(hyErr.Hint, true))
	}
	
	// Stack trace (if available and not empty)
	if len(hyErr.StackTrace) > 0 {
		builder.WriteString(fmt.Sprintf("\n%sStack trace:%s\n", ColorGray, ColorReset))
		for i, trace := range hyErr.StackTrace {
			builder.WriteString(fmt.Sprintf("  %s%d. %s%s\n", ColorGray, i+1, trace, ColorReset))
		}
	}
	
	return builder.String()
}

// FormatErrorPlain formats an error without colors (for non-terminal output)
func FormatErrorPlain(err error) string {
	hyErr, ok := err.(*HyException)
	if !ok {
		return err.Error()
	}
	
	var builder strings.Builder
	
	// File:Line location (if available)
	if hyErr.File != "" {
		builder.WriteString(fmt.Sprintf("%s:%d", hyErr.File, hyErr.Line))
		if hyErr.Column > 0 {
			builder.WriteString(fmt.Sprintf(":%d", hyErr.Column))
		}
		builder.WriteString("\n")
	}
	
	// Error type and message
	errorType := hyErr.Type
	if errorType == "" {
		errorType = "Error"
	}
	
	builder.WriteString(fmt.Sprintf("%s: %s\n", errorType, hyErr.Message))
	
	// Hints (if available)
	if hyErr.Hint != nil && len(hyErr.Hint.Suggestions) > 0 {
		builder.WriteString(formatHintSuggestions(hyErr.Hint, false))
	}
	
	// Stack trace (if available and not empty)
	if len(hyErr.StackTrace) > 0 {
		builder.WriteString("\nStack trace:\n")
		for i, trace := range hyErr.StackTrace {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, trace))
		}
	}
	
	return builder.String()
}

// FormatErrorPlainWithContext formats an error without colors but with code context
func FormatErrorPlainWithContext(err error, codeContext []string, currentLine int) string {
	hyErr, ok := err.(*HyException)
	if !ok {
		return err.Error()
	}
	
	var builder strings.Builder
	
	// File:Line location (if available)
	if hyErr.File != "" {
		builder.WriteString(fmt.Sprintf("%s:%d", hyErr.File, hyErr.Line))
		if hyErr.Column > 0 {
			builder.WriteString(fmt.Sprintf(":%d", hyErr.Column))
		}
		builder.WriteString("\n")
	}
	
	// Code context (if available)
	if len(codeContext) > 0 && currentLine > 0 {
		builder.WriteString("\nCode context:\n")
		
		// Show lines with line numbers
		startLine := currentLine - len(codeContext) + 1
		for i, line := range codeContext {
			lineNum := startLine + i
			isCurrentLine := lineNum == currentLine
			
			if isCurrentLine {
				builder.WriteString(fmt.Sprintf("  > %-4d %s\n", lineNum, line))
			} else {
				builder.WriteString(fmt.Sprintf("    %-4d %s\n", lineNum, line))
			}
		}
		builder.WriteString("\n")
	}
	
	// Error type and message
	errorType := hyErr.Type
	if errorType == "" {
		errorType = "Error"
	}
	
	builder.WriteString(fmt.Sprintf("%s: %s\n", errorType, hyErr.Message))
	
	// Hints (if available)
	if hyErr.Hint != nil && len(hyErr.Hint.Suggestions) > 0 {
		builder.WriteString(formatHintSuggestions(hyErr.Hint, false))
	}
	
	// Stack trace (if available and not empty)
	if len(hyErr.StackTrace) > 0 {
		builder.WriteString("\nStack trace:\n")
		for i, trace := range hyErr.StackTrace {
			builder.WriteString(fmt.Sprintf("  %d. %s\n", i+1, trace))
		}
	}
	
	return builder.String()
}

// IsTerminal checks if the output supports colors (simplified check)
func IsTerminal() bool {
	// For now, we'll just return true. In a real implementation,
	// you'd check if stdout is a terminal using syscall or a library
	return true
}

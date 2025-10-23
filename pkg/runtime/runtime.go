package runtime

import (
	"fmt"
	"os"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// ExecuteSource compiles and executes Polyloft source code
func ExecuteSource(source, filename string) error {
	// Tokenize
	lx := &lexer.Lexer{}
	items := lx.Scan([]byte(source))
	
	// Parse
	p := parser.NewWithSource(items, filename, source)
	prog, err := p.Parse()
	if err != nil {
		formattedErr := engine.FormatError(err)
		fmt.Fprint(os.Stderr, formattedErr)
		return err
	}
	
	// Execute
	_, err = engine.EvalWithContextAndSource(prog, engine.Options{Stdout: os.Stdout}, filename, ".", source)
	if err != nil {
		formattedErr := engine.FormatError(err)
		fmt.Fprint(os.Stderr, formattedErr)
		return err
	}

	return nil
}

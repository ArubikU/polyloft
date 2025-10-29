package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ArubikU/polyloft/internal/engine"
	"github.com/ArubikU/polyloft/internal/lexer"
	"github.com/ArubikU/polyloft/internal/parser"
)

// Start launches a minimal line-oriented REPL with simple evaluation.
// Meta commands:
//
//	:quit  - exit the REPL
//	:help  - show brief help
func Start(in io.Reader, out io.Writer, prompt string) {
	s := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, prompt)
		if !s.Scan() {
			fmt.Fprintln(out)
			return
		}
		line := strings.TrimSpace(s.Text())
		switch line {
		case ":quit", ":q":
			fmt.Fprintln(out, "bye")
			return
		case ":help", ":h":
			fmt.Fprintln(out, "Polyloft REPL commands:")
			fmt.Fprintln(out, "  :help  Show this help")
			fmt.Fprintln(out, "  :quit  Exit the REPL")
			continue
		case "":
			continue
		}
		// Wrap line as a program (expression statement)
		lx := &lexer.Lexer{}
		items := lx.Scan([]byte(line))
		p := parser.NewWithFile(items, "<repl>")
		prog, err := p.Parse()
		if err != nil {
			fmt.Fprintln(out, "error:", err)
			continue
		}
		v, err := engine.Eval(prog, engine.Options{Stdout: out})
		if err != nil {
			fmt.Fprintln(out, "error:", err)
			continue
		}
		if v != nil {
			fmt.Fprintln(out, "=", v)
		}
	}
}

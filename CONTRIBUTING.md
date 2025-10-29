# Contributing to Polyloft (Go implementation scaffold)

Thanks for your interest in contributing! This repo contains the language docs and a Go-based scaffold for an interpreter/compiler pipeline.

Guiding principles:
- Keep packages small and cohesive (`internal/{lexer,parser,ast,engine}`).
- Write tests alongside new behavior.
- Prefer clear code and comments over cleverness.
- Leave TODOs with enough detail for the next contributor.

Suggested workflow:
1. Open an issue describing the change.
2. Create a feature branch.
3. Add tests for the new capability (happy path + 1 edge case).
4. Implement until tests pass.
5. Submit a PR referencing the issue.

Project layout (Go):
- `cmd/polyloft`: CLI entrypoint.
- `internal/ast`: AST node types.
- `internal/lexer`: Tokenizer.
- `internal/parser`: Parser to AST.
- `internal/engine`: Execution/runtime.
- `internal/repl`: REPL utilities used by the CLI.
- `internal/version`: Version metadata.

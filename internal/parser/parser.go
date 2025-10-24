package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ArubikU/polyloft/internal/ast"
	"github.com/ArubikU/polyloft/internal/common"
	"github.com/ArubikU/polyloft/internal/lexer"
)

// Package parser turns tokens into an AST.

type Parser struct {
	items      []lexer.Item
	pos        int
	file       string
	sourceCode string // Store source code for better error messages
}

func New(items []lexer.Item) *Parser                      { return &Parser{items: items} }
func NewWithFile(items []lexer.Item, file string) *Parser { return &Parser{items: items, file: file} }

// NewWithSource creates a parser with source code for better error messages
func NewWithSource(items []lexer.Item, file string, source string) *Parser {
	return &Parser{items: items, file: file, sourceCode: source}
}

// ParseError includes file and precise position for better diagnostics.
type ParseError struct {
	File       string
	Pos        ast.Position
	Msg        string
	SourceLine string
	Token      lexer.Item
}

func (e ParseError) Error() string {
	var buf strings.Builder

	if e.File != "" {
		fmt.Fprintf(&buf, "%s:%d:%d: ", e.File, e.Pos.Line, e.Pos.Col)
	} else {
		fmt.Fprintf(&buf, "%d:%d: ", e.Pos.Line, e.Pos.Col)
	}

	buf.WriteString(e.Msg)

	// Add source line preview if available
	if e.SourceLine != "" {
		buf.WriteString("\n")
		buf.WriteString(e.SourceLine)
		buf.WriteString("\n")

		// Add pointer to the exact column
		if e.Pos.Col > 0 {
			for i := 1; i < e.Pos.Col; i++ {
				// Preserve tabs in the source line for proper alignment
				if i <= len(e.SourceLine) && e.SourceLine[i-1] == '\t' {
					buf.WriteString("\t")
				} else {
					buf.WriteString(" ")
				}
			}
			buf.WriteString("^")

			// Add additional context about the token
			if e.Token.Tok != lexer.EOF && e.Token.Tok != lexer.ILLEGAL {
				buf.WriteString(" (found ")
				buf.WriteString(lexer.FormatToken(e.Token))
				buf.WriteString(")")
			}
		}
	}

	return buf.String()
}

// ModifierError represents errors related to invalid modifier usage
type ModifierError struct {
	File     string
	Pos      ast.Position
	Modifier common.Modifier
	Context  string
}

func (e ModifierError) Error() string {
	modifierName := common.ModifierNames[e.Modifier]
	msg := fmt.Sprintf("modifier '%s' cannot be used in %s context", modifierName, e.Context)
	if e.File != "" {
		return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Pos.Line, e.Pos.Col, msg)
	}
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Col, msg)
}

// PrivacyError represents errors related to invalid privacy level usage
type PrivacyError struct {
	File    string
	Pos     ast.Position
	Privacy common.PrivacyLevel
	Context string
}

func (e PrivacyError) Error() string {
	privacyName := common.PrivacyLevelNames[e.Privacy]
	msg := fmt.Sprintf("privacy level '%s' cannot be used in %s context", privacyName, e.Context)
	if e.File != "" {
		return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Pos.Line, e.Pos.Col, msg)
	}
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Col, msg)
}

func (p *Parser) errf(format string, a ...any) error {
	item := p.curr()
	pos := item.Start
	msg := fmt.Sprintf(format, a...)

	// Get the source line for context
	sourceLine := p.getSourceLine(pos.Line)

	return ParseError{
		File:       p.file,
		Pos:        pos,
		Msg:        msg,
		SourceLine: sourceLine,
		Token:      item,
	}
}

// getSourceLine retrieves the source code line at the given line number
func (p *Parser) getSourceLine(lineNum int) string {
	if p.sourceCode == "" {
		return ""
	}

	lines := strings.Split(p.sourceCode, "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		return lines[lineNum-1]
	}
	return ""
}

// validateModifier checks if a modifier can be used in the given context
func (p *Parser) validateModifier(modifier common.Modifier, context string) error {
	if !common.CanUseModifier(modifier, context) {
		pos := p.curr().Start
		return ModifierError{File: p.file, Pos: pos, Modifier: modifier, Context: context}
	}
	return nil
}

// validatePrivacyLevel checks if a privacy level can be used in the given context
func (p *Parser) validatePrivacyLevel(privacy common.PrivacyLevel, context string) error {
	if !common.CanUsePrivacyLevel(privacy, context) {
		pos := p.curr().Start
		return PrivacyError{File: p.file, Pos: pos, Privacy: privacy, Context: context}
	}
	return nil
}

// parsePrivacyLevel parses a privacy level from tokens
func (p *Parser) parsePrivacyLevel() (common.PrivacyLevel, error) {
	switch p.curr().Tok {
	case lexer.KW_PUBLIC:
		p.next()
		return common.PrivacyPublic, nil
	case lexer.KW_PRIVATE:
		p.next()
		return common.PrivacyPrivate, nil
	case lexer.KW_PROTECTED:
		p.next()
		return common.PrivacyProtected, nil
	default:
		return common.PrivacyPackage, nil // default visibility
	}
}

// parseModifier parses a modifier from tokens
func (p *Parser) parseModifier() (common.Modifier, error) {
	switch p.curr().Tok {
	case lexer.KW_STATIC:
		p.next()
		return common.ModifierStatic, nil
	case lexer.KW_ABSTRACT:
		p.next()
		return common.ModifierAbstract, nil
	default:
		return common.ModifierNone, nil
	}
}

func (p *Parser) curr() lexer.Item {
	if p.pos >= len(p.items) {
		return lexer.Item{Tok: lexer.EOF}
	}
	return p.items[p.pos]
}
func (p *Parser) previous() lexer.Item {
	if p.pos-1 < 0 || p.pos-1 >= len(p.items) {
		return lexer.Item{Tok: lexer.EOF}
	}
	return p.items[p.pos-1]
}
func (p *Parser) next() lexer.Item { i := p.curr(); p.pos++; return i }
func (p *Parser) accept(tok lexer.Token) bool {
	if p.curr().Tok == tok {
		p.pos++
		return true
	}
	return false
}

// Parse program: sequence of statements until EOF.
func (p *Parser) Parse() (*ast.Program, error) {
	prog := &ast.Program{}
	for p.curr().Tok != lexer.EOF {
		st, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if st != nil {
			prog.Stmts = append(prog.Stmts, st)
		}
		// optional semicolons between statements
		for p.accept(lexer.SEMI) {
		}
	}
	return prog, nil
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch p.curr().Tok {
	case lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED:
		// Check if this is an access modifier for a class, enum, or function
		savedPos := p.pos
		accessLevel := p.curr().Lit
		p.next()

		switch p.curr().Tok {
		case lexer.KW_CLASS, lexer.KW_ABSTRACT:
			return p.parseClassWithAccess(accessLevel)
		case lexer.KW_ENUM:
			return p.parseEnum(accessLevel)
		case lexer.KW_DEF:
			// Parse function declaration with access modifier
			// Skip the 'def' keyword and continue with normal function parsing
			p.next()
			id := p.curr()
			if id.Tok != lexer.IDENT {
				return nil, p.errf("expected function name after def")
			}
			name := id.Lit
			p.next()

			// Optional generic type parameters: def identity<T>(value: T) -> T
			var typeParams []ast.TypeParam
			if p.curr().Tok == lexer.LT {
				// Try to parse generic type parameters
				params, err := p.tryParseGenericTypeParams()
				if err == nil && params != nil {
					typeParams = params
				}
				// If parsing fails, just continue
			}

			if !p.accept(lexer.LPAREN) {
				return nil, p.errf("expected '('")
			}
			var params []ast.Parameter
			if p.curr().Tok != lexer.RPAREN {
				for {
					tok := p.curr()
					if tok.Tok != lexer.IDENT {
						return nil, p.errf("expected parameter name")
					}
					paramName := tok.Lit
					p.next()

					var paramType string
					var isVariadic bool
					if p.accept(lexer.COLON) {
						if p.curr().Tok != lexer.IDENT {
							return nil, p.errf("expected parameter type")
						}
						paramType = p.curr().Lit
						p.next()

						// Check for variadic parameter (type...)
						if p.accept(lexer.ELLIPSIS) {
							isVariadic = true
						}
					}

					params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

					// If this is a variadic parameter, it must be the last one
					if isVariadic && p.curr().Tok == lexer.COMMA {
						return nil, p.errf("variadic parameter must be the last parameter")
					}

					if p.accept(lexer.COMMA) {
						continue
					}
					break
				}
			}
			if !p.accept(lexer.RPAREN) {
				return nil, p.errf("expected ')'")
			}

			// Parse optional return type
			var returnType string
			if p.accept(lexer.RARROW) {
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected return type after '->'")
				}
				returnType = p.curr().Lit
				p.next()
			}

			var body []ast.Stmt
			expectsEnd := false
			if p.accept(lexer.ASSIGN) {
				expr, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}
				body = append(body, &ast.ReturnStmt{Value: expr})
			} else {
				if !p.accept(lexer.COLON) {
					return nil, p.errf("expected ':' or '=' before function body")
				}

				var err error
				body, err = p.parseBlock()
				if err != nil {
					return nil, err
				}
				expectsEnd = true
			}

			if expectsEnd {
				if !p.accept(lexer.KW_END) {
					return nil, p.errf("expected 'end' to close function body")
				}
			}

			return &ast.DefStmt{
				Name:        name,
				Params:      params,
				Body:        body,
				ReturnType:  ast.TypeFromString(returnType),
				AccessLevel: accessLevel,
				Modifiers:   []string{accessLevel},
				TypeParams:  typeParams,
			}, nil
		case lexer.KW_SEALED:
			if len(p.items) > p.pos+1 {
				nextTok := p.items[p.pos+1].Tok
				if nextTok == lexer.KW_CLASS {
					return p.parseClassWithAccess(accessLevel)
				}
				if nextTok == lexer.KW_ENUM {
					return p.parseEnum(accessLevel)
				}
			}
		}

		// Not a class, enum, or function; treat as variable/field declaration
		p.pos = savedPos
		return p.parseVarLike()

	case lexer.KW_SEALED:
		if len(p.items) > p.pos+1 {
			nextTok := p.items[p.pos+1].Tok
			// Check for "sealed abstract class"
			if nextTok == lexer.KW_ABSTRACT {
				if len(p.items) > p.pos+2 && p.items[p.pos+2].Tok == lexer.KW_CLASS {
					return p.parseClassWithAccess("public")
				}
			}
			if nextTok == lexer.KW_CLASS {
				return p.parseClassWithAccess("public")
			}
			if nextTok == lexer.KW_ENUM {
				return p.parseEnum("public")
			}
			if nextTok == lexer.KW_INTERFACE {
				return p.parseInterfaceWithModifiers("public", true)
			}
		}
		return nil, p.errf("expected 'class', 'enum', 'interface', or 'abstract class' after 'sealed'")
	case lexer.KW_STATIC, lexer.KW_VAR, lexer.KW_LET, lexer.KW_CONST, lexer.KW_FINAL:
		return p.parseVarLike()
	case lexer.KW_IF:
		return p.parseIf()
	case lexer.KW_FOR:
		return p.parseForIn()
	case lexer.KW_LOOP:
		return p.parseLoop()
	case lexer.KW_INTERFACE:
		return p.parseInterfaceWithModifiers("public", false)
	case lexer.KW_CLASS:
		return p.parseClassWithAccess("public") // default to public
	case lexer.KW_ENUM:
		return p.parseEnum("public") // default to public
	case lexer.KW_RECORD:
		return p.parseRecord("public") // default to public
	case lexer.KW_ABSTRACT:
		return p.parseClassWithAccess("public") // default to public
	case lexer.KW_IMPORT:
		return p.parseImport()
	case lexer.KW_TRY:
		return p.parseTry()
	case lexer.KW_THROW:
		return p.parseThrow()
	case lexer.KW_DEFER:
		return p.parseDefer()
	case lexer.KW_SELECT:
		return p.parseSelect()
	case lexer.KW_SWITCH:
		return p.parseSwitch()
	case lexer.KW_BREAK:
		p.next()
		return &ast.BreakStmt{}, nil
	case lexer.KW_CONTINUE:
		p.next()
		return &ast.ContinueStmt{}, nil
	case lexer.KW_RETURN:
		p.next()
		// Check if return has a value or is just "return" alone
		// If next token is a block terminator, return has no value
		var expr ast.Expr
		var err error
		curr := p.curr()
		if curr.Tok != lexer.KW_END && curr.Tok != lexer.KW_ELSE && curr.Tok != lexer.KW_ELIF && curr.Tok != lexer.EOF && curr.Tok != lexer.RBRACE {
			expr, err = p.parseExpr(0)
			if err != nil {
				return nil, err
			}
		}
		return &ast.ReturnStmt{Value: expr}, nil
	case lexer.KW_DEF:
		p.next()
		id := p.curr()
		if id.Tok != lexer.IDENT {
			return nil, p.errf("expected function name after def")
		}
		name := id.Lit
		p.next()

		// Optional generic type parameters: def identity<T>(value: T) -> T
		var typeParams []ast.TypeParam
		if p.curr().Tok == lexer.LT {
			// Try to parse generic type parameters
			params, err := p.tryParseGenericTypeParams()
			if err == nil && params != nil {
				typeParams = params
			}
			// If parsing fails, just continue (might be a comparison in the function)
		}

		if !p.accept(lexer.LPAREN) {
			return nil, p.errf("expected '('")
		}
		var params []ast.Parameter
		if p.curr().Tok != lexer.RPAREN {
			for {
				tok := p.curr()
				if tok.Tok != lexer.IDENT {
					return nil, p.errf("expected parameter name")
				}
				paramName := tok.Lit
				p.next()

				var paramType string
				var isVariadic bool
				if p.accept(lexer.COLON) {
					if p.curr().Tok != lexer.IDENT {
						return nil, p.errf("expected parameter type")
					}
					paramType = p.curr().Lit
					p.next()

					// Check for variadic parameter (type...)
					if p.accept(lexer.ELLIPSIS) {
						isVariadic = true
					}
				}

				params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

				// If this is a variadic parameter, it must be the last one
				if isVariadic && p.curr().Tok == lexer.COMMA {
					return nil, p.errf("variadic parameter must be the last parameter")
				}

				if p.accept(lexer.COMMA) {
					continue
				}
				break
			}
		}
		if !p.accept(lexer.RPAREN) {
			return nil, p.errf("expected ')'")
		}

		// Parse optional return type
		var returnType string
		if p.accept(lexer.RARROW) {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected return type after '->'")
			}
			returnType = p.curr().Lit
			p.next()
		}

		var body []ast.Stmt
		expectsEnd := false
		if p.accept(lexer.ASSIGN) {
			expr, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			body = append(body, &ast.ReturnStmt{Value: expr})
		} else {
			if !p.accept(lexer.COLON) {
				return nil, p.errf("expected ':' or '=' before function body")
			}

			var err error
			body, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
			expectsEnd = true
		}

		if expectsEnd {
			if !p.accept(lexer.KW_END) {
				return nil, p.errf("expected 'end' to close function body")
			}
		}

		return &ast.DefStmt{
			Name:        name,
			Params:      params,
			Body:        body,
			ReturnType:  ast.TypeFromString(returnType),
			AccessLevel: "public", // default to public
			Modifiers:   []string{"public"},
			TypeParams:  typeParams,
		}, nil
	default:
		// Try to parse as assignment statement first
		if p.curr().Tok == lexer.IDENT || p.curr().Tok == lexer.KW_THIS {
			// Look ahead to see if this is an assignment
			saved_pos := p.pos

			// Parse the left side (could be identifier or field access)
			lhs, err := p.parseExpr(0)
			if err != nil {
				p.pos = saved_pos // restore position
				return nil, err
			}

			// Check if next token is assignment
			if p.curr().Tok == lexer.ASSIGN {
				assignPos := p.curr().Start // capture position of '=' token
				p.next()                    // consume '='
				rhs, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}

				// Create assignment statement
				return &ast.AssignStmt{Target: lhs, Value: rhs, Pos: assignPos}, nil
			} else {
				// Not an assignment, treat as expression statement
				p.pos = saved_pos // restore position
				e, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}
				return &ast.ExprStmt{X: e}, nil
			}
		} else {
			// Parse as expression statement
			e, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			return &ast.ExprStmt{X: e}, nil
		}
	}
}

// parseVarLike handles declarations with modifiers and optional types:
// [public|private|protected]? [static]? (var|let|const|final) name ( ':' Type )? ( '=' expr | ':=' expr )?
func (p *Parser) parseVarLike() (ast.Stmt, error) {
	mods := []string{}
	for {
		switch p.curr().Tok {
		case lexer.KW_PUBLIC:
			mods = append(mods, "public")
			p.next()
			continue
		case lexer.KW_PRIVATE:
			mods = append(mods, "private")
			p.next()
			continue
		case lexer.KW_PROTECTED:
			mods = append(mods, "protected")
			p.next()
			continue
		case lexer.KW_STATIC:
			mods = append(mods, "static")
			p.next()
			continue
		}
		break
	}

	kindTok := p.curr()
	kind := ""
	switch kindTok.Tok {
	case lexer.KW_VAR:
		kind = "var"
	case lexer.KW_LET:
		kind = "let"
	case lexer.KW_CONST:
		kind = "const"
	case lexer.KW_FINAL:
		kind = "final"
	default:
		return nil, p.errf("expected declaration keyword after modifiers")
	}
	p.next()

	id := p.curr()
	if id.Tok != lexer.IDENT {
		return nil, p.errf("expected identifier after %s", kind)
	}
	name := id.Lit
	p.next()

	typ := ""
	inferred := false
	if p.accept(lexer.COLON) {
		if p.curr().Tok != lexer.IDENT {
			return nil, p.errf("expected type name after ':'")
		}
		typ = p.curr().Lit
		p.next()
		if p.accept(lexer.ASSIGN) {
			// typed with '='
		} else if p.accept(lexer.COLONASSIGN) {
			inferred = true
		} else if kind == "var" || kind == "let" {
			// allow typed without initializer -> default nil
			return &ast.LetStmt{Name: name, Value: &ast.NilLit{}, Type: ast.TypeFromString(typ), Modifiers: mods, Kind: kind, Inferred: inferred}, nil
		} else {
			return nil, p.errf("expected '=' after typed %s %s", kind, name)
		}
	} else if p.accept(lexer.COLONASSIGN) {
		inferred = true
	} else if !p.accept(lexer.ASSIGN) {
		// allow declaration without initializer only for var/let
		if kind == "var" || kind == "let" {
			return &ast.LetStmt{Name: name, Value: &ast.NilLit{}, Type: ast.TypeFromString(typ), Modifiers: mods, Kind: kind, Inferred: inferred}, nil
		}
		return nil, p.errf("expected '=' or ':=' after %s %s", kind, name)
	}

	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	return &ast.LetStmt{Name: name, Value: expr, Type: ast.TypeFromString(typ), Modifiers: mods, Kind: kind, Inferred: inferred}, nil
}

// parseBlock parses a sequence of statements until a terminator token.
func (p *Parser) parseBlock() ([]ast.Stmt, error) {
	var body []ast.Stmt
	for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.KW_ELSE && p.curr().Tok != lexer.KW_ELIF {
		st, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		body = append(body, st)
	}
	return body, nil
}

func (p *Parser) parseIf() (ast.Stmt, error) {
	// if <expr>: <block> (elif <expr>: <block>)* (else: <block>)? end
	p.next()
	cond, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}
	if !p.accept(lexer.COLON) {
		return nil, p.errf("expected ':' after if condition")
	}
	thenB, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	clauses := []ast.IfClause{{Cond: cond, Body: thenB}}
	for p.curr().Tok == lexer.KW_ELIF {
		p.next()
		c, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		if !p.accept(lexer.COLON) {
			return nil, p.errf("expected ':' after elif condition")
		}
		b, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, ast.IfClause{Cond: c, Body: b})
	}
	var elseB []ast.Stmt
	if p.curr().Tok == lexer.KW_ELSE {
		p.next()
		if !p.accept(lexer.COLON) {
			return nil, p.errf("expected ':' after else")
		}
		b, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		elseB = b
	}
	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close if")
	}
	return &ast.IfStmt{Clauses: clauses, Else: elseB}, nil
}

func (p *Parser) parseForIn() (ast.Stmt, error) {
	// for <ident>[,<ident>...] in <expr>: <block> end
	p.next()

	// Parse iteration variable(s)
	var names []string
	id := p.curr()
	if id.Tok != lexer.IDENT {
		return nil, p.errf("expected identifier after 'for'")
	}
	names = append(names, id.Lit)
	p.next()

	// Check for additional iteration variables (destructuring)
	for p.accept(lexer.COMMA) {
		id := p.curr()
		if id.Tok != lexer.IDENT {
			return nil, p.errf("expected identifier after ',' in for-in")
		}
		names = append(names, id.Lit)
		p.next()
	}

	if !p.accept(lexer.KW_IN) {
		return nil, p.errf("expected 'in' in for-in")
	}
	it, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	// Check for optional where clause
	var whereClause ast.Expr
	if p.accept(lexer.KW_WHERE) {
		whereClause, err = p.parseExpr(0)
		if err != nil {
			return nil, err
		}
	}

	if !p.accept(lexer.COLON) {
		return nil, p.errf("expected ':' after for-in header")
	}
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close for-in")
	}

	// For backward compatibility, set Name to the first name
	name := ""
	if len(names) > 0 {
		name = names[0]
	}

	return &ast.ForInStmt{Name: name, Names: names, Iterable: it, Where: whereClause, Body: body}, nil
}

func (p *Parser) parseLoop() (ast.Stmt, error) {
	// loop <block> end
	p.next()
	body, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close loop")
	}
	return &ast.LoopStmt{Body: body}, nil
}

// parseInterfaceWithModifiers parses: [sealed] interface Name[<TypeParams>] [(permits)] NEWLINE? (method signatures... | static fields... ) end
func (p *Parser) parseInterfaceWithModifiers(accessLevel string, isSealed bool) (ast.Stmt, error) {
	// Consume 'sealed' if present
	if p.curr().Tok == lexer.KW_SEALED {
		p.next()
	}

	// interface <Ident> ... end
	if p.curr().Tok != lexer.KW_INTERFACE {
		return nil, p.errf("expected 'interface'")
	}
	p.next() // consume 'interface'

	id := p.curr()
	if id.Tok != lexer.IDENT {
		return nil, p.errf("expected interface name")
	}
	name := id.Lit
	p.next()

	// Parse optional generic type parameters: <T>, <K, V>, etc.
	var typeParams []ast.TypeParam
	if p.curr().Tok == lexer.LT {
		p.next() // consume '<'
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected type parameter name")
			}
			paramName := p.curr().Lit
			p.next()

			// Check for bound: T extends Bound
			var bounds []string
			if p.curr().Tok == lexer.KW_EXTENDS {
				p.next() // consume 'extends'
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected bound type")
				}
				bounds = append(bounds, p.curr().Lit)
				p.next()
			}

			typeParams = append(typeParams, ast.TypeParam{
				Name:   paramName,
				Bounds: bounds,
			})

			if p.curr().Tok == lexer.COMMA {
				p.next() // consume ','
				continue
			}
			break
		}
		if p.curr().Tok != lexer.GT {
			return nil, p.errf("expected '>' after type parameters")
		}
		p.next() // consume '>'
	}

	// Optional permits list for sealed interfaces
	var permits []string
	if isSealed && p.curr().Tok == lexer.LPAREN {
		p.next() // consume '('
		for p.curr().Tok == lexer.IDENT {
			permits = append(permits, p.curr().Lit)
			p.next()
			if !p.accept(lexer.COMMA) {
				break
			}
		}
		if !p.accept(lexer.RPAREN) {
			return nil, p.errf("expected ')' after sealed permit list")
		}
	}

	var methods []ast.MethodSignature
	var fields []ast.FieldDecl

	// Parse method signatures and static fields until 'end'
	depth := 1
	for depth > 0 && p.curr().Tok != lexer.EOF {
		if p.curr().Tok == lexer.KW_END {
			depth--
			if depth == 0 {
				p.next() // consume final 'end'
				break
			}
			p.next() // consume nested 'end'
			continue
		}
		if p.curr().Tok == lexer.KW_INTERFACE || p.curr().Tok == lexer.KW_CLASS {
			depth++
		}

		// Parse method signature: def methodName(params) -> ReturnType
		if p.curr().Tok == lexer.KW_DEF {
			method, err := p.parseMethodSignature()
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		} else if p.curr().Tok == lexer.KW_STATIC {
			// Parse static field: static var/let/const name: Type = value
			field, err := p.parseFieldDecl()
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		} else if p.curr().Tok == lexer.KW_VAR || p.curr().Tok == lexer.KW_LET || p.curr().Tok == lexer.KW_CONST {
			// Parse static field without explicit static keyword: var/let/const name: Type = value
			field, err := p.parseFieldDecl()
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		} else {
			// Skip unknown tokens
			p.next()
		}
	}

	return &ast.InterfaceDecl{
		Name:        name,
		Methods:     methods,
		TypeParams:  typeParams,
		IsSealed:    isSealed,
		Permits:     permits,
		Fields:      fields,
		AccessLevel: accessLevel,
	}, nil
}

// parseInterface parses: interface Name NEWLINE? (method signatures... ) end
// Kept for backward compatibility, delegates to parseInterfaceWithModifiers
func (p *Parser) parseInterface() (ast.Stmt, error) {
	return p.parseInterfaceWithModifiers("public", false)
}

// parseMethodSignature parses method signatures for interfaces
func (p *Parser) parseMethodSignature() (ast.MethodSignature, error) {
	p.next() // consume 'def'

	// Parse method name (can be identifier or operator for operator overloading)
	var name string
	switch p.curr().Tok {
	case lexer.IDENT:
		name = p.curr().Lit
	case lexer.PLUS:
		name = "+"
	case lexer.MINUS:
		name = "-"
	case lexer.STAR:
		name = "*"
	case lexer.SLASH:
		name = "/"
	case lexer.PERCENT:
		name = "%"
	case lexer.EQ:
		name = "=="
	case lexer.NEQ:
		name = "!="
	case lexer.LT:
		name = "<"
	case lexer.LTE:
		name = "<="
	case lexer.GT:
		name = ">"
	case lexer.GTE:
		name = ">="
	default:
		return ast.MethodSignature{}, p.errf("expected method name")
	}
	p.next()

	if !p.accept(lexer.LPAREN) {
		return ast.MethodSignature{}, p.errf("expected '(' after method name")
	}

	var params []ast.Parameter
	if p.curr().Tok != lexer.RPAREN {
		for {
			if p.curr().Tok != lexer.IDENT {
				return ast.MethodSignature{}, p.errf("expected parameter name")
			}
			paramName := p.curr().Lit
			p.next()

			var paramType string
			var isVariadic bool
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return ast.MethodSignature{}, p.errf("expected parameter type")
				}
				paramType = p.curr().Lit
				p.next()

				// Check for variadic parameter (type...)
				if p.accept(lexer.ELLIPSIS) {
					isVariadic = true
				}
			}

			params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

			// If this is a variadic parameter, it must be the last one
			if isVariadic && p.curr().Tok == lexer.COMMA {
				return ast.MethodSignature{}, p.errf("variadic parameter must be the last parameter")
			}

			if !p.accept(lexer.COMMA) {
				break
			}
		}
	}

	if !p.accept(lexer.RPAREN) {
		return ast.MethodSignature{}, p.errf("expected ')' after parameters")
	}

	var returnType string
	if p.accept(lexer.RARROW) {
		if p.curr().Tok != lexer.IDENT {
			return ast.MethodSignature{}, p.errf("expected return type after '->'")
		}
		returnType = p.curr().Lit
		p.next()
	}

	// Check for default method implementation
	var hasDefault bool
	var defaultBody []ast.Stmt
	if p.accept(lexer.COLON) {
		hasDefault = true
		// Parse default method body until 'end' - use parseBlock() for cleaner handling
		var err error
		defaultBody, err = p.parseBlock()
		if err != nil {
			return ast.MethodSignature{}, err
		}
		if !p.accept(lexer.KW_END) {
			return ast.MethodSignature{}, p.errf("expected 'end' to close default method body")
		}
	}

	return ast.MethodSignature{
		Name:        name,
		Params:      params,
		ReturnType:  ast.TypeFromString(returnType),
		HasDefault:  hasDefault,
		DefaultBody: defaultBody,
	}, nil
}

// parseFieldDecl parses field declarations: [modifiers] var/let/const/final name: Type [= value]
func (p *Parser) parseFieldDecl() (ast.FieldDecl, error) {
	var modifiers []string

	// Parse modifiers
	for {
		switch p.curr().Tok {
		case lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED, lexer.KW_STATIC:
			modifiers = append(modifiers, p.curr().Lit)
			p.next()
		default:
			goto done_modifiers
		}
	}
done_modifiers:

	// Parse var/let/const/final
	if p.curr().Tok != lexer.KW_VAR && p.curr().Tok != lexer.KW_LET &&
		p.curr().Tok != lexer.KW_CONST && p.curr().Tok != lexer.KW_FINAL {
		return ast.FieldDecl{}, p.errf("expected var, let, final or const for field declaration")
	}
	kind := p.curr().Lit
	if kind == "const" || kind == "final" {
		modifiers = append(modifiers, kind)
	}
	p.next()

	// Parse field name
	if p.curr().Tok != lexer.IDENT {
		return ast.FieldDecl{}, p.errf("expected field name")
	}
	name := p.curr().Lit
	p.next()

	// Parse type annotation
	var fieldType string
	if p.accept(lexer.COLON) {
		if p.curr().Tok != lexer.IDENT {
			return ast.FieldDecl{}, p.errf("expected field type")
		}
		fieldType = p.curr().Lit
		p.next()
	}

	// Parse default value
	var value ast.Expr
	if p.accept(lexer.ASSIGN) {
		var err error
		value, err = p.parseExpr(0)
		if err != nil {
			return ast.FieldDecl{}, err
		}
	}

	return ast.FieldDecl{
		Name:      name,
		Type:      ast.TypeFromString(fieldType),
		Modifiers: modifiers,
		InitValue: value,
	}, nil
}

// parseMethodDecl parses method declarations: [annotations] [modifiers] def name(params): ReturnType body end
func (p *Parser) parseMethodDecl() (ast.MethodDecl, error) {
	var (
		annotations     []ast.Annotation
		modifiers       []string
		annotationFlags common.AnnotationFlags
	)

	// Parse annotations and capture metadata once so we can expand behaviour later.
	for p.curr().Tok == lexer.AT {
		p.next() // consume @
		if p.curr().Tok != lexer.IDENT {
			return ast.MethodDecl{}, p.errf("expected annotation name after @")
		}
		rawName := p.curr().Lit
		p.next()

		info, known := common.LookupAnnotation(rawName)
		if known {
			annotationFlags = annotationFlags.Merge(info.Flags)
		}
		annotations = append(annotations, ast.Annotation{
			Raw:        rawName,
			Normalized: info.Name,
			Known:      known,
		})
	}

	// Parse modifiers (handled before getting here, but could be enhanced)
	if p.curr().Tok == lexer.KW_PUBLIC || p.curr().Tok == lexer.KW_PRIVATE ||
		p.curr().Tok == lexer.KW_PROTECTED || p.curr().Tok == lexer.KW_STATIC ||
		p.curr().Tok == lexer.KW_ABSTRACT {
		for {
			switch p.curr().Tok {
			case lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED,
				lexer.KW_STATIC, lexer.KW_ABSTRACT:
				modifiers = append(modifiers, p.curr().Lit)
				p.next()
			default:
				goto done_method_modifiers
			}
		}
	}
done_method_modifiers:

	// Determine whether the method is abstract before we parse its body.
	isAbstract := false
	for _, mod := range modifiers {
		if mod == "abstract" {
			isAbstract = true
			break
		}
	}

	// Parse def
	if p.curr().Tok != lexer.KW_DEF {
		return ast.MethodDecl{}, p.errf("expected 'def' for method declaration")
	}
	p.next()

	// Parse method name (can be identifier or operator for operator overloading)
	var name string
	switch p.curr().Tok {
	case lexer.IDENT:
		name = p.curr().Lit
	case lexer.PLUS:
		name = "+"
	case lexer.MINUS:
		name = "-"
	case lexer.STAR:
		name = "*"
	case lexer.SLASH:
		name = "/"
	case lexer.PERCENT:
		name = "%"
	case lexer.EQ:
		name = "=="
	case lexer.NEQ:
		name = "!="
	case lexer.LT:
		name = "<"
	case lexer.LTE:
		name = "<="
	case lexer.GT:
		name = ">"
	case lexer.GTE:
		name = ">="
	default:
		return ast.MethodDecl{}, p.errf("expected method name")
	}
	p.next()

	// Parse parameters
	if !p.accept(lexer.LPAREN) {
		return ast.MethodDecl{}, p.errf("expected '(' after method name")
	}

	var params []ast.Parameter
	if p.curr().Tok != lexer.RPAREN {
		for {
			if p.curr().Tok != lexer.IDENT {
				return ast.MethodDecl{}, p.errf("expected parameter name")
			}
			paramName := p.curr().Lit
			p.next()

			var paramType string
			var isVariadic bool
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return ast.MethodDecl{}, p.errf("expected parameter type")
				}
				paramType = p.curr().Lit
				p.next()

				// Check for variadic parameter (type...)
				if p.accept(lexer.ELLIPSIS) {
					isVariadic = true
				}
			}

			params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

			// If this is a variadic parameter, it must be the last one
			if isVariadic && p.curr().Tok == lexer.COMMA {
				return ast.MethodDecl{}, p.errf("variadic parameter must be the last parameter")
			}

			if !p.accept(lexer.COMMA) {
				break
			}
		}
	}

	if !p.accept(lexer.RPAREN) {
		return ast.MethodDecl{}, p.errf("expected ')' after parameters")
	}

	// Parse return type (optional)
	var returnType string
	if p.accept(lexer.RARROW) {
		if p.curr().Tok != lexer.IDENT {
			return ast.MethodDecl{}, p.errf("expected return type after '->'")
		}
		returnType = p.curr().Lit
		p.next()
	} else {
		// If no return type specified, default to Void for methods without return
		returnType = "Void"
	}

	var body []ast.Stmt

	switch {
	case isAbstract:
		if p.accept(lexer.ASSIGN) {
			return ast.MethodDecl{}, p.errf("abstract methods cannot have an expression body")
		}
		if p.accept(lexer.COLON) {
			if !p.accept(lexer.KW_END) {
				return ast.MethodDecl{}, p.errf("expected 'end' after abstract method declaration")
			}
		}
	case p.accept(lexer.ASSIGN):
		expr, err := p.parseExpr(0)
		if err != nil {
			return ast.MethodDecl{}, err
		}
		body = append(body, &ast.ReturnStmt{Value: expr})
	default:
		if !p.accept(lexer.COLON) {
			return ast.MethodDecl{}, p.errf("expected ':' or '=' before method body")
		}
		// Parse method body until 'end' - use parseBlock() for cleaner handling
		var err error
		body, err = p.parseBlock()
		if err != nil {
			return ast.MethodDecl{}, err
		}
		if !p.accept(lexer.KW_END) {
			return ast.MethodDecl{}, p.errf("expected 'end' to close method body")
		}
	}

	return ast.MethodDecl{
		Name:        name,
		Params:      params,
		ReturnType:  ast.TypeFromString(returnType),
		Body:        body,
		Modifiers:   modifiers,
		IsAbstract:  isAbstract,
		IsOverride:  annotationFlags.IsOverride,
		Annotations: annotations,
	}, nil
}

// parseConstructorDecl parses constructor declarations: ClassName(params): body end
func (p *Parser) parseConstructorDecl() (*ast.ConstructorDecl, error) {
	// Constructor name (should match class name)
	p.next() // consume constructor name

	// Parse parameters
	if !p.accept(lexer.LPAREN) {
		return nil, p.errf("expected '(' after constructor name")
	}

	var params []ast.Parameter
	if p.curr().Tok != lexer.RPAREN {
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected parameter name")
			}
			paramName := p.curr().Lit
			p.next()

			var paramType string
			var isVariadic bool
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected parameter type")
				}
				paramType = p.curr().Lit
				p.next()

				// Check for variadic parameter (type...)
				if p.accept(lexer.ELLIPSIS) {
					isVariadic = true
				}
			}

			params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

			// If this is a variadic parameter, it must be the last one
			if isVariadic && p.curr().Tok == lexer.COMMA {
				return nil, p.errf("variadic parameter must be the last parameter")
			}

			if !p.accept(lexer.COMMA) {
				break
			}
		}
	}

	if !p.accept(lexer.RPAREN) {
		return nil, p.errf("expected ')' after constructor parameters")
	}

	// Parse colon for constructor body
	if !p.accept(lexer.COLON) {
		return nil, p.errf("expected ':' before constructor body")
	}

	// Parse constructor body until 'end' - use parseBlock() for cleaner handling
	var body []ast.Stmt
	var err error
	body, err = p.parseBlock()
	if err != nil {
		return nil, err
	}
	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close constructor body")
	}

	return &ast.ConstructorDecl{
		Params: params,
		Body:   body,
	}, nil
}

// parseClassWithAccess parses a class with the given access level
func (p *Parser) parseClassWithAccess(accessLevel string) (ast.Stmt, error) {
	return p.parseClassInternal(accessLevel)
}

// parseClassInternal is the internal implementation for parsing classes
func (p *Parser) parseClassInternal(accessLevel string) (ast.Stmt, error) {
	// Check for modifiers: abstract, sealed, or both (may have been consumed earlier)
	isAbstract := false
	isSealed := false
	permits := []string{}

	// Handle sealed abstract class or abstract sealed class
	if p.curr().Tok == lexer.KW_SEALED || p.curr().Tok == lexer.KW_ABSTRACT {
		if p.curr().Tok == lexer.KW_SEALED {
			isSealed = true
			p.next()
			// Check if next is abstract
			if p.curr().Tok == lexer.KW_ABSTRACT {
				isAbstract = true
				p.next()
			}
			if p.curr().Tok != lexer.KW_CLASS {
				return nil, p.errf("expected 'class' after 'sealed' or 'sealed abstract'")
			}
		} else if p.curr().Tok == lexer.KW_ABSTRACT {
			isAbstract = true
			p.next()
			// Check if next is sealed
			if p.curr().Tok == lexer.KW_SEALED {
				isSealed = true
				p.next()
			}
			if p.curr().Tok != lexer.KW_CLASS {
				return nil, p.errf("expected 'class' after 'abstract' or 'abstract sealed'")
			}
		}
	}

	// class <Ident> <TypeParams>? (< <Parent>)? (implements <Interface1>, <Interface2>...)? ... end
	p.next() // consume 'class'
	id := p.curr()
	if id.Tok != lexer.IDENT {
		return nil, p.errf("expected class name")
	}
	name := id.Lit
	p.next()

	// Optional generic type parameters: class Builder<T>, class Pair<K, V>
	var typeParams []ast.TypeParam
	if p.curr().Tok == lexer.LT {
		// Try to parse generic type parameters
		params, err := p.tryParseGenericTypeParams()
		if err == nil && params != nil {
			typeParams = params
		} else {
		}
	}

	// Optional permits list for sealed classes
	if isSealed && p.accept(lexer.LPAREN) {
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected identifier in sealed permit list")
			}
			permits = append(permits, p.curr().Lit)
			p.next()
			if p.accept(lexer.COMMA) {
				continue
			}
			break
		}
		if !p.accept(lexer.RPAREN) {
			return nil, p.errf("expected ')' after sealed permit list")
		}
	}

	var parent string
	var parentTypeParams []ast.TypeParam
	var implements []string

	// Check for inheritance: < Parent (only if we didn't parse type params)
	// If we have type params, inheritance uses 'extends' keyword instead
	if len(typeParams) == 0 && p.curr().Tok == lexer.LT {
		p.next() // consume '<'
		if p.curr().Tok != lexer.IDENT {
			return nil, p.errf("expected parent class name after '<'")
		}
		parent = p.curr().Lit
		p.next()

		// Check if parent has generic type parameters: < Parent<T>
		if p.curr().Tok == lexer.LT {
			params, err := p.tryParseGenericTypeParams()
			if err == nil && params != nil {
				parentTypeParams = params
			}
		}
	} else if p.curr().Tok == lexer.KW_EXTENDS {
		// New style: class MyClass<T> extends BaseClass or BaseClass<T>
		p.next() // consume 'extends'
		if p.curr().Tok != lexer.IDENT {
			return nil, p.errf("expected parent class name after 'extends'")
		}
		parent = p.curr().Lit
		p.next()

		// Check if parent has generic type parameters: extends Parent<T>
		if p.curr().Tok == lexer.LT {
			params, err := p.tryParseGenericTypeParams()
			if err == nil && params != nil {
				parentTypeParams = params
			}
		}
	}

	// Check for implements: implements Interface1, Interface2
	if p.curr().Tok == lexer.KW_IMPLEMENTS {
		p.next() // consume 'implements'
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected interface name after 'implements'")
			}
			implements = append(implements, p.curr().Lit)
			p.next()
			if !p.accept(lexer.COMMA) {
				break
			}
		}
	}

	// Parse class body until 'end'
	var fields []ast.FieldDecl
	var methods []ast.MethodDecl
	var constructor *ast.ConstructorDecl

	depth := 1
	for depth > 0 && p.curr().Tok != lexer.EOF {
		if p.curr().Tok == lexer.KW_END {
			depth--
			if depth == 0 {
				p.next() // consume final 'end'
				break
			}
			p.next() // consume nested 'end'
			continue
		}
		if p.curr().Tok == lexer.KW_CLASS || p.curr().Tok == lexer.KW_INTERFACE {
			depth++
		}

		// Parse class members
		switch p.curr().Tok {
		case lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED:
			// Check if this is a method or field by looking ahead
			saved_pos := p.pos
			// Skip modifiers
			for p.curr().Tok == lexer.KW_PUBLIC || p.curr().Tok == lexer.KW_PRIVATE ||
				p.curr().Tok == lexer.KW_PROTECTED || p.curr().Tok == lexer.KW_STATIC {
				p.next()
			}

			if p.curr().Tok == lexer.KW_DEF {
				// Method declaration
				p.pos = saved_pos // restore position
				method, err := p.parseMethodDecl()
				if err != nil {
					return nil, err
				}
				methods = append(methods, method)
			} else {
				// Field declaration
				p.pos = saved_pos // restore position
				field, err := p.parseFieldDecl()
				if err != nil {
					return nil, err
				}
				fields = append(fields, field)
			}
		case lexer.KW_STATIC:
			// Check if this is a static method or static field
			saved_pos := p.pos
			p.next() // consume 'static'

			if p.curr().Tok == lexer.KW_DEF {
				// Static method
				p.pos = saved_pos // restore position
				method, err := p.parseMethodDecl()
				if err != nil {
					return nil, err
				}
				methods = append(methods, method)
			} else {
				// Static field
				p.pos = saved_pos // restore position
				field, err := p.parseFieldDecl()
				if err != nil {
					return nil, err
				}
				fields = append(fields, field)
			}
		case lexer.KW_VAR, lexer.KW_LET, lexer.KW_CONST, lexer.KW_FINAL:
			// Field declaration
			field, err := p.parseFieldDecl()
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		case lexer.KW_ABSTRACT:
			// Abstract method declaration (abstract def methodName)
			method, err := p.parseMethodDecl()
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		case lexer.KW_DEF:
			// Method declaration
			method, err := p.parseMethodDecl()
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		case lexer.IDENT:
			// Check if this is a constructor (same name as class)
			if p.curr().Lit == name && len(p.items) > p.pos+1 && p.items[p.pos+1].Tok == lexer.LPAREN {
				// Constructor
				constr, err := p.parseConstructorDecl()
				if err != nil {
					return nil, err
				}
				constructor = constr
			} else {
				// Skip unknown identifier
				p.next()
			}
		default:
			// Skip unknown tokens
			p.next()
		}

		// optional semicolons
		for p.accept(lexer.SEMI) {
		}
	}

	return &ast.ClassDecl{
		Name:             name,
		Parent:           parent,
		ParentTypeParams: parentTypeParams,
		Implements:       implements,
		IsAbstract:       isAbstract,
		AccessLevel:      accessLevel,
		IsSealed:         isSealed,
		Permits:          permits,
		Fields:           fields,
		Methods:          methods,
		Constructor:      constructor,
		TypeParams:       typeParams,
	}, nil
}

// parseImport: import a.b.c { X, Y }
func (p *Parser) parseImport() (ast.Stmt, error) {
	// import <dotted.ident> ( '{' ident (',' ident)* '}' )?
	p.next() // consume 'import'
	// dotted path
	parts := []string{}
	if p.curr().Tok != lexer.IDENT {
		return nil, p.errf("expected module path after import")
	}
	parts = append(parts, p.curr().Lit)
	p.next()
	for p.accept(lexer.DOT) {
		if p.curr().Tok != lexer.IDENT {
			return nil, p.errf("expected identifier after '.' in import path")
		}
		parts = append(parts, p.curr().Lit)
		p.next()
	}
	names := []string{}
	if p.accept(lexer.LBRACE) {
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected identifier in import list")
			}
			names = append(names, p.curr().Lit)
			p.next()
			if p.accept(lexer.COMMA) {
				continue
			}
			break
		}
		if !p.accept(lexer.RBRACE) {
			return nil, p.errf("expected '}' to close import list")
		}
	}
	return &ast.ImportStmt{Path: parts, Names: names}, nil
}

// Pratt parser precedence levels
const (
	precTernary = iota
	precRange   // for ... range operator
	precOr
	precAnd
	precEq
	precCmp
	precAdd
	precMul
	precUnary
	precCall
)

func (p *Parser) precedence(tok lexer.Token) int {
	switch tok {
	case lexer.QUESTION:
		return precTernary
	case lexer.ELLIPSIS:
		return precRange
	case lexer.OR:
		return precOr
	case lexer.AND:
		return precAnd
	case lexer.EQ, lexer.NEQ:
		return precEq
	case lexer.LT, lexer.LTE, lexer.GT, lexer.GTE:
		return precCmp
	case lexer.KW_INSTANCEOF:
		return precCmp
	case lexer.PLUS, lexer.MINUS:
		return precAdd
	case lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return precMul
	default:
		return -1
	}
}

func (p *Parser) parseExpr(minPrec int) (ast.Expr, error) {
	// Parse prefix
	var left ast.Expr
	tok := p.curr()
	switch tok.Tok {
	case lexer.IDENT:
		name := tok.Lit
		p.next()
		if p.curr().Tok == lexer.LT {
			// Try to parse as generic type
			typeParams, err := p.tryParseGenericTypeParams()
			if err == nil && typeParams != nil {
				// Successfully parsed generic type params
				// Now expect () for constructor call
				if p.curr().Tok == lexer.LPAREN {
					p.next() // consume '('
					var args []ast.Expr
					if p.curr().Tok != lexer.RPAREN {
						for {
							arg, err := p.parseExpr(0)
							if err != nil {
								return nil, err
							}
							args = append(args, arg)
							if !p.accept(lexer.COMMA) {
								break
							}
						}
					}
					if !p.accept(lexer.RPAREN) {
						return nil, p.errf("expected ')' after generic type constructor arguments")
					}
					left = &ast.GenericCallExpr{
						Name:       name,
						TypeParams: typeParams,
						Args:       args,
					}
				} else {
					// Generic type reference without constructor call
					// For now, treat as regular identifier
					left = &ast.Ident{Name: name}
				}
			} else {
				// Not a generic type, just an identifier
				left = &ast.Ident{Name: name}
			}
		} else {
			left = &ast.Ident{Name: name}
		}
	case lexer.NUMBER:
		// Legacy support - parse as float
		f, err := strconv.ParseFloat(tok.Lit, 64)
		if err != nil {
			return nil, err
		}
		left = &ast.NumberLit{Value: f}
		p.next()
	case lexer.INT:
		// Parse as integer
		i, err := strconv.ParseInt(tok.Lit, 10, 64)
		if err != nil {
			return nil, err
		}
		left = &ast.NumberLit{Value: int(i)}
		p.next()
	case lexer.FLOAT:
		// Parse float literals with explicit suffix as float64 to keep runtime consistent
		litStr := strings.TrimSuffix(tok.Lit, "f")
		f, err := strconv.ParseFloat(litStr, 64)
		if err != nil {
			return nil, err
		}
		left = &ast.NumberLit{Value: f}
		p.next()
	case lexer.STRING:
		// strip quotes and simple escapes
		s, err := unquote(tok.Lit)
		if err != nil {
			return nil, err
		}
		left = &ast.StringLit{Value: s}
		p.next()
	case lexer.KW_TRUE:
		left = &ast.BoolLit{Value: true}
		p.next()
	case lexer.KW_FALSE:
		left = &ast.BoolLit{Value: false}
		p.next()
	case lexer.KW_NIL:
		left = &ast.NilLit{}
		p.next()
	case lexer.KW_THIS:
		left = &ast.Ident{Name: "this"}
		p.next()
	case lexer.KW_SUPER:
		// Parse super as a special identifier that can be called
		left = &ast.Ident{Name: "super"}
		p.next()
	case lexer.KW_THREAD:
		// thread spawn do ... end or thread join expr
		p.next() // consume 'thread'
		if p.curr().Tok == lexer.KW_SPAWN {
			p.next() // consume 'spawn'
			if p.curr().Tok != lexer.KW_DO {
				return nil, p.errf("expected 'do' after 'thread spawn'")
			}
			p.next() // consume 'do'

			// Parse thread body until 'end'
			var body []ast.Stmt
			for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					body = append(body, stmt)
				}
			}

			if !p.accept(lexer.KW_END) {
				return nil, p.errf("expected 'end' to close thread spawn block")
			}

			left = &ast.ThreadSpawnExpr{Body: body}
		} else if p.curr().Tok == lexer.KW_JOIN {
			p.next() // consume 'join'
			threadExpr, err := p.parseExpr(precUnary)
			if err != nil {
				return nil, err
			}
			left = &ast.ThreadJoinExpr{Thread: threadExpr}
		} else {
			return nil, p.errf("expected 'spawn' or 'join' after 'thread'")
		}
	case lexer.KW_CHANNEL:
		// channel[Type]()
		p.next() // consume 'channel'
		if p.curr().Tok != lexer.LBRACK {
			return nil, p.errf("expected '[' after 'channel'")
		}
		p.next() // consume '['

		if p.curr().Tok != lexer.IDENT {
			return nil, p.errf("expected type name in channel[Type]")
		}
		elemType := p.curr().Lit
		p.next() // consume type name

		if p.curr().Tok != lexer.RBRACK {
			return nil, p.errf("expected ']' after channel type")
		}
		p.next() // consume ']'

		// Expect () for channel creation
		if p.curr().Tok != lexer.LPAREN {
			return nil, p.errf("expected '()' after channel[Type]")
		}
		p.next() // consume '('
		if p.curr().Tok != lexer.RPAREN {
			return nil, p.errf("expected ')' in channel[Type]()")
		}
		p.next() // consume ')'

		left = &ast.ChannelExpr{ElemType: elemType}
	case lexer.MINUS, lexer.NOT:
		p.next()
		x, err := p.parseExpr(precUnary)
		if err != nil {
			return nil, err
		}
		op := ast.OpNeg
		if tok.Tok == lexer.NOT {
			op = ast.OpNot
		}
		left = &ast.UnaryExpr{Op: op, X: x}
	case lexer.LBRACK:
		// array literal: [a, b, c]
		p.next()
		var elems []ast.Expr
		if p.curr().Tok != lexer.RBRACK {
			for {
				e, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}
				elems = append(elems, e)
				if p.accept(lexer.COMMA) {
					continue
				}
				break
			}
		}
		if !p.accept(lexer.RBRACK) {
			return nil, p.errf("expected ']' in array literal")
		}
		left = &ast.ArrayLit{Elems: elems}
	case lexer.LBRACE:
		// map literal: { key: expr, ... } with string keys
		p.next()
		var pairs []ast.MapPair
		if p.curr().Tok != lexer.RBRACE {
			for {
				k := p.curr()
				if k.Tok != lexer.IDENT && k.Tok != lexer.STRING {
					return nil, p.errf("expected key in map literal")
				}
				key := k.Lit
				if k.Tok == lexer.STRING {
					s, err := unquote(k.Lit)
					if err != nil {
						return nil, err
					}
					key = s
				}
				p.next()
				if !p.accept(lexer.COLON) {
					return nil, p.errf("expected ':' after key in map literal")
				}
				v, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}
				pairs = append(pairs, ast.MapPair{Key: key, Value: v})
				if p.accept(lexer.COMMA) {
					continue
				}
				break
			}
		}
		if !p.accept(lexer.RBRACE) {
			return nil, p.errf("expected '}' in map literal")
		}
		left = &ast.MapLit{Pairs: pairs}
	case lexer.LPAREN:
		expr, err := p.parseParenthesized()
		if err != nil {
			return nil, err
		}
		left = expr
	default:
		// Provide more helpful error message based on context
		tokenName := lexer.TokenName(tok.Tok)
		if tok.Lit != "" {
			return nil, p.errf("unexpected %s, expected expression (literal, identifier, or operator)", tokenName)
		}
		return nil, p.errf("unexpected %s, expected expression", tokenName)
	}

	// Parse infix/postfix
	for {
		tok = p.curr()
		// Normalize word operators if present
		_ = p.maybeWordOp(left)
		tok = p.curr()
		// call
		if tok.Tok == lexer.LPAREN {
			p.next()
			var args []ast.Expr
			if p.curr().Tok != lexer.RPAREN {
				for {
					e, err := p.parseExpr(0)
					if err != nil {
						return nil, err
					}
					args = append(args, e)
					if p.accept(lexer.COMMA) {
						continue
					}
					break
				}
			}
			if !p.accept(lexer.RPAREN) {
				return nil, p.errf("expected ')'")
			}
			left = &ast.CallExpr{Callee: left, Args: args}
			continue
		}

		// index
		if tok.Tok == lexer.LBRACK {
			p.next()
			idx, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}

			// Check for range expression: arr[start...end]
			if p.curr().Tok == lexer.ELLIPSIS {
				p.next() // consume ...
				endExpr, err := p.parseExpr(0)
				if err != nil {
					return nil, err
				}
				idx = &ast.RangeExpr{Start: idx, End: endExpr, Inclusive: true}
			}

			if !p.accept(lexer.RBRACK) {
				return nil, p.errf("expected ']' in index expression")
			}
			left = &ast.IndexExpr{X: left, Index: idx}
			continue
		}

		// field access or method call
		if tok.Tok == lexer.DOT {
			p.next()
			id := p.curr()
			var fieldName string

			// Allow certain keywords as field/method names, but primarily expect IDENT
			switch id.Tok {
			case lexer.IDENT:
				fieldName = id.Lit
			case lexer.KW_INSTANCEOF:
				fieldName = "instanceof" // special case for instanceof operator
			case lexer.KW_CATCH:
				fieldName = "catch" // allow .catch() method calls
			case lexer.KW_FINALLY:
				fieldName = "finally" // allow .finally() method calls
			default:
				return nil, p.errf("expected field or method name after '.', got token: %v", id.Tok)
			}

			p.next()
			left = &ast.FieldExpr{X: left, Name: fieldName}
			continue
		}

		prec := p.precedence(tok.Tok)
		if prec < minPrec {
			break
		}

		// Special handling for instanceof as operator (not method call)
		// Only process if we're not immediately after a field access
		if tok.Tok == lexer.KW_INSTANCEOF {
			// Check if the left side is a field access to avoid conflicts with method calls
			if _, isFieldExpr := left.(*ast.FieldExpr); isFieldExpr {
				// This is likely Obj.instanceof(...) method call, don't treat as operator
				break
			}

			p.next() // consume 'instanceof'
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected type name after 'instanceof'")
			}
			typeName := p.curr().Lit
			p.next()

			var variable, modifier string
			// Check for optional modifier and variable assignment
			// Examples: "10 instanceof int final foo" or "obj instanceof String var myVar"
			if p.curr().Tok == lexer.KW_VAR ||
				p.curr().Tok == lexer.KW_LET || p.curr().Tok == lexer.KW_CONST ||
				p.curr().Tok == lexer.KW_FINAL {
				modifier = p.curr().Lit
				p.next()
				if p.curr().Tok == lexer.IDENT {
					variable = p.curr().Lit
					p.next()
				}
			} else if p.curr().Tok == lexer.IDENT {
				// Check if this is a variable name (not another operator)
				next := p.items[p.pos+1 : p.pos+2]
				if len(next) == 0 || (next[0].Tok != lexer.EQ && next[0].Tok != lexer.PLUS &&
					next[0].Tok != lexer.MINUS && next[0].Tok != lexer.STAR && next[0].Tok != lexer.SLASH) {
					variable = p.curr().Lit
					p.next()
				}
			}

			left = &ast.InstanceOfExpr{
				Expr:     left,
				TypeName: typeName,
				Variable: variable,
				Modifier: modifier,
			}
			continue
		}

		// Special handling for ternary operator
		if tok.Tok == lexer.QUESTION {
			prec := p.precedence(tok.Tok)
			if prec < minPrec {
				break
			}

			p.next() // consume '?'
			trueBranch, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}

			if !p.accept(lexer.COLON) {
				return nil, p.errf("expected ':' in ternary expression")
			}

			falseBranch, err := p.parseExpr(prec + 1)
			if err != nil {
				return nil, err
			}

			left = &ast.TernaryExpr{
				Condition:   left,
				TrueBranch:  trueBranch,
				FalseBranch: falseBranch,
			}
			continue
		}

		// Special handling for range operator
		if tok.Tok == lexer.ELLIPSIS {
			prec := p.precedence(tok.Tok)
			if prec < minPrec {
				break
			}

			p.next() // consume '...'
			right, err := p.parseExpr(prec + 1)
			if err != nil {
				return nil, err
			}

			left = &ast.RangeExpr{
				Start:     left,
				End:       right,
				Inclusive: true,
			}
			continue
		}

		op := tok
		p.next()
		right, err := p.parseExpr(prec + 1)
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{Op: p.toOp(op.Tok), Lhs: left, Rhs: right}
	}
	return left, nil
}

// parseParenthesized handles parentheses in expressions:
// - Grouped expressions: (expr)
// - Lambda expressions: (a, b) => expr
// - Multi-line expressions with proper comma handling
func (p *Parser) parseParenthesized() (ast.Expr, error) {
	if !p.accept(lexer.LPAREN) {
		return nil, p.errf("expected '('")
	}

	// Handle empty parentheses case
	if p.curr().Tok == lexer.RPAREN {
		p.next() // consume ')'
		// Check if this is an empty lambda: () => expr or () -> Type => expr
		if p.curr().Tok == lexer.RARROW {
			// Parse return type
			p.next() // consume '->'
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected return type after '->'")
			}
			returnType := p.curr().Lit
			p.next()

			if !p.accept(lexer.ARROW) {
				return nil, p.errf("expected '=>' after return type")
			}

			// Check for multi-line lambda with do...end
			if p.curr().Tok == lexer.KW_DO {
				p.next() // consume 'do'

				// Parse block body
				var blockBody []ast.Stmt
				for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
					stmt, err := p.parseStmt()
					if err != nil {
						return nil, err
					}
					if stmt != nil {
						blockBody = append(blockBody, stmt)
					}
				}

				if !p.accept(lexer.KW_END) {
					return nil, p.errf("expected 'end' to close multi-line lambda")
				}

				return &ast.LambdaExpr{Params: []ast.Parameter{}, BlockBody: blockBody, ReturnType: ast.TypeFromString(returnType), IsBlock: true}, nil
			}

			// Single expression lambda
			body, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			return &ast.LambdaExpr{Params: []ast.Parameter{}, Body: body, ReturnType: ast.TypeFromString(returnType)}, nil
		} else if p.curr().Tok == lexer.ARROW {
			p.next() // consume '=>'

			// Check for multi-line lambda with do...end
			if p.curr().Tok == lexer.KW_DO {
				p.next() // consume 'do'

				// Parse block body
				var blockBody []ast.Stmt
				for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
					stmt, err := p.parseStmt()
					if err != nil {
						return nil, err
					}
					if stmt != nil {
						blockBody = append(blockBody, stmt)
					}
				}

				if !p.accept(lexer.KW_END) {
					return nil, p.errf("expected 'end' to close multi-line lambda")
				}

				return &ast.LambdaExpr{Params: []ast.Parameter{}, BlockBody: blockBody, IsBlock: true}, nil
			}

			// Single expression lambda
			body, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}
			return &ast.LambdaExpr{Params: []ast.Parameter{}, Body: body}, nil
		}
		// Empty parentheses group - not valid in most contexts
		return nil, p.errf("empty parentheses not allowed")
	}

	// Try to determine if this is a lambda or grouped expression
	// Look ahead to see if we have identifier(s) followed by '=>'
	savedPos := p.pos
	var params []ast.Parameter
	isLambda := false

	// Try to parse as lambda parameters
	if p.curr().Tok == lexer.IDENT {
		for {
			if p.curr().Tok != lexer.IDENT {
				break
			}
			paramName := p.curr().Lit
			p.next()

			var paramType string
			var isVariadic bool
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					break // Not a valid lambda parameter
				}
				paramType = p.curr().Lit
				p.next()

				// Check for variadic parameter (type...)
				if p.accept(lexer.ELLIPSIS) {
					isVariadic = true
				}
			}

			params = append(params, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

			// If this is a variadic parameter, it must be the last one
			if isVariadic && p.curr().Tok == lexer.COMMA {
				break // Invalid lambda
			}

			if p.curr().Tok == lexer.COMMA {
				p.next() // consume ','
				continue
			}
			break
		}

		// Check if we have closing paren followed by arrow (possibly with return type)
		if p.curr().Tok == lexer.RPAREN {
			p.next() // consume ')'

			// Check for optional return type
			if p.curr().Tok == lexer.RARROW {
				p.next() // consume '->'
				if p.curr().Tok == lexer.IDENT {
					p.next() // skip return type for now, we'll parse it later
				}
			}

			if p.curr().Tok == lexer.ARROW {
				isLambda = true
			}
		}
	}

	if isLambda {
		// Parse lambda: we already consumed the ')' and found '=>'
		// First, restore position to re-parse return type properly
		p.pos = savedPos

		// Re-parse parameters with proper type handling
		var finalParams []ast.Parameter
		if p.curr().Tok == lexer.IDENT {
			for {
				if p.curr().Tok != lexer.IDENT {
					break
				}
				paramName := p.curr().Lit
				p.next()

				var paramType string
				var isVariadic bool
				if p.accept(lexer.COLON) {
					if p.curr().Tok != lexer.IDENT {
						return nil, p.errf("expected parameter type")
					}
					paramType = p.curr().Lit
					p.next()

					// Check for variadic parameter (type...)
					if p.accept(lexer.ELLIPSIS) {
						isVariadic = true
					}
				}

				finalParams = append(finalParams, ast.Parameter{Name: paramName, Type: ast.TypeFromString(paramType), IsVariadic: isVariadic})

				// If this is a variadic parameter, it must be the last one
				if isVariadic && p.curr().Tok == lexer.COMMA {
					return nil, p.errf("variadic parameter must be the last parameter")
				}

				if p.curr().Tok == lexer.COMMA {
					p.next() // consume ','
					continue
				}
				break
			}
		}

		if !p.accept(lexer.RPAREN) {
			return nil, p.errf("expected ')' in lambda parameters")
		}

		// Parse optional return type
		var returnType string
		if p.accept(lexer.RARROW) {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected return type after '->'")
			}
			returnType = p.curr().Lit
			p.next()
		}

		if !p.accept(lexer.ARROW) {
			return nil, p.errf("expected '=>' in lambda expression")
		}

		// Check for multi-line lambda with do...end
		if p.curr().Tok == lexer.KW_DO {
			p.next() // consume 'do'

			// Parse block body
			var blockBody []ast.Stmt
			for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					blockBody = append(blockBody, stmt)
				}
			}

			if !p.accept(lexer.KW_END) {
				return nil, p.errf("expected 'end' to close multi-line lambda")
			}

			return &ast.LambdaExpr{Params: finalParams, BlockBody: blockBody, ReturnType: ast.TypeFromString(returnType), IsBlock: true}, nil
		}

		// Single expression lambda
		body, err := p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		return &ast.LambdaExpr{Params: finalParams, Body: body, ReturnType: ast.TypeFromString(returnType), IsBlock: false}, nil
	}

	// Not a lambda, restore position and parse as grouped expression
	p.pos = savedPos

	// Parse the expression inside parentheses
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	// Expect closing parenthesis
	if !p.accept(lexer.RPAREN) {
		return nil, p.errf("expected ')' to close grouped expression")
	}

	return expr, nil
}

func unquote(q string) (string, error) {
	if len(q) >= 2 {
		// Handle float-quoted strings
		if q[0] == '"' && q[len(q)-1] == '"' {
			// handle simple escapes
			b := make([]rune, 0, len(q)-2)
			escaped := false
			for _, r := range q[1 : len(q)-1] {
				if escaped {
					switch r {
					case 'n':
						b = append(b, '\n')
					case 't':
						b = append(b, '\t')
					case '"':
						b = append(b, '"')
					case '\\':
						b = append(b, '\\')
					default:
						b = append(b, r)
					}
					escaped = false
					continue
				}
				if r == '\\' {
					escaped = true
					continue
				}
				b = append(b, r)
			}
			return string(b), nil
		}
		// Handle single-quoted strings
		if q[0] == '\'' && q[len(q)-1] == '\'' {
			// handle simple escapes
			b := make([]rune, 0, len(q)-2)
			escaped := false
			for _, r := range q[1 : len(q)-1] {
				if escaped {
					switch r {
					case 'n':
						b = append(b, '\n')
					case 't':
						b = append(b, '\t')
					case '\'':
						b = append(b, '\'')
					case '\\':
						b = append(b, '\\')
					default:
						b = append(b, r)
					}
					escaped = false
					continue
				}
				if r == '\\' {
					escaped = true
					continue
				}
				b = append(b, r)
			}
			return string(b), nil
		}
	}
	// not quoted
	return q, nil
}

// precedence maps lexer tokens to parse precedence.
func (p *Parser) toOp(tok lexer.Token) int {
	switch tok {
	case lexer.PLUS:
		return ast.OpPlus
	case lexer.MINUS:
		return ast.OpMinus
	case lexer.STAR:
		return ast.OpMul
	case lexer.SLASH:
		return ast.OpDiv
	case lexer.PERCENT:
		return ast.OpMod
	case lexer.EQ:
		return ast.OpEq
	case lexer.NEQ:
		return ast.OpNeq
	case lexer.LT:
		return ast.OpLt
	case lexer.LTE:
		return ast.OpLte
	case lexer.GT:
		return ast.OpGt
	case lexer.GTE:
		return ast.OpGte
	case lexer.AND:
		return ast.OpAnd
	case lexer.OR:
		return ast.OpOr
	default:
		return 0
	}
}

// Support word operators by rewriting identifiers 'and'/'or'/'not' into the corresponding tokens during expression parsing.
func (p *Parser) maybeWordOp(left ast.Expr) bool {
	if p.curr().Tok != lexer.IDENT {
		return false
	}
	switch strings.ToLower(p.curr().Lit) {
	case "and":
	case "&&":
		p.items[p.pos].Tok = lexer.AND
		return true
	case "or":
	case "||":
		p.items[p.pos].Tok = lexer.OR
		return true
	case "not":
	case "!":
		p.items[p.pos].Tok = lexer.NOT
		return true
	}
	return false
}

// parseEnum parses an enum declaration
func (p *Parser) parseEnum(accessLevel string) (ast.Stmt, error) {
	isSealed := false

	if p.curr().Tok == lexer.KW_SEALED {
		isSealed = true
		p.next()
		if p.curr().Tok != lexer.KW_ENUM {
			return nil, p.errf("expected 'enum' after 'sealed'")
		}
	}

	if p.curr().Tok != lexer.KW_ENUM {
		return nil, p.errf("expected 'enum'")
	}

	// consume 'enum' (might have been prefixed by modifiers)
	p.next()

	// optional sealed modifier handled in parseStmt() via lookahead; here just parse name
	if p.curr().Tok != lexer.IDENT {
		return nil, p.errf("expected enum name")
	}
	name := p.curr().Lit
	p.next()

	// Optional permits: enum Name(AllowedA, AllowedB)
	permits := []string{}
	if p.curr().Tok == lexer.LPAREN {
		p.next()
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected identifier in enum permit list")
			}
			permits = append(permits, p.curr().Lit)
			p.next()
			if p.accept(lexer.COMMA) {
				continue
			}
			break
		}
		if !p.accept(lexer.RPAREN) {
			return nil, p.errf("expected ')' after enum permit list")
		}
	}

	var values []ast.EnumValue
	var fields []ast.FieldDecl
	var methods []ast.MethodDecl
	var constructor *ast.ConstructorDecl
	parsingMembers := false

	for p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
		if !parsingMembers {
			if p.curr().Tok == lexer.COLON {
				p.next()
				continue
			}

			switch p.curr().Tok {
			case lexer.IDENT:
				if p.curr().Lit == name && len(p.items) > p.pos+1 && p.items[p.pos+1].Tok == lexer.LPAREN {
					parsingMembers = true

					continue
				}

				// Parse enum value
				valueName := p.curr().Lit
				p.next()

				var args []ast.Expr
				if p.curr().Tok == lexer.LPAREN {
					p.next() // consume '('
					if p.curr().Tok != lexer.RPAREN {
						for {
							arg, err := p.parseExpr(0)
							if err != nil {
								return nil, err
							}
							args = append(args, arg)
							if p.accept(lexer.COMMA) {
								continue
							}
							break
						}
					}
					if !p.accept(lexer.RPAREN) {
						return nil, p.errf("expected ')' after enum value arguments")
					}
				}

				values = append(values, ast.EnumValue{Name: valueName, Args: args})

				if p.accept(lexer.COMMA) {
					continue
				}
			case lexer.SEMI:
				parsingMembers = true
				p.next()
				continue
			case lexer.KW_DEF, lexer.KW_VAR, lexer.KW_LET, lexer.KW_CONST, lexer.KW_FINAL,
				lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED, lexer.KW_STATIC:
				parsingMembers = true
				continue
			default:
				parsingMembers = true
				continue
			}
			continue
		}

		// Parsing members (fields, methods, constructor)
		switch p.curr().Tok {
		case lexer.KW_DEF:
			method, err := p.parseMethodDecl()
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		case lexer.KW_VAR, lexer.KW_LET, lexer.KW_CONST, lexer.KW_FINAL:
			field, err := p.parseFieldDecl()
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		case lexer.KW_PUBLIC, lexer.KW_PRIVATE, lexer.KW_PROTECTED, lexer.KW_STATIC:
			savedPos := p.pos
			for p.curr().Tok == lexer.KW_PUBLIC || p.curr().Tok == lexer.KW_PRIVATE ||
				p.curr().Tok == lexer.KW_PROTECTED || p.curr().Tok == lexer.KW_STATIC {
				p.next()
			}

			switch p.curr().Tok {
			case lexer.KW_DEF:
				p.pos = savedPos
				method, err := p.parseMethodDecl()
				if err != nil {
					return nil, err
				}
				methods = append(methods, method)
			case lexer.KW_VAR, lexer.KW_LET, lexer.KW_CONST, lexer.KW_FINAL:
				p.pos = savedPos
				field, err := p.parseFieldDecl()
				if err != nil {
					return nil, err
				}
				fields = append(fields, field)
			default:
				p.pos = savedPos
				field, err := p.parseFieldDecl()
				if err != nil {
					return nil, err
				}
				fields = append(fields, field)
			}
		case lexer.IDENT:
			if p.curr().Lit == name && len(p.items) > p.pos+1 && p.items[p.pos+1].Tok == lexer.LPAREN {
				t, err := p.parseConstructorDecl()
				if err != nil {
					return nil, err
				}
				constructor = t
			} else {
				p.next()
			}
		default:
			p.next()
		}

		for p.accept(lexer.SEMI) {
		}
	}

	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' after enum body")
	}

	return &ast.EnumDecl{
		Name:        name,
		AccessLevel: accessLevel,
		IsSealed:    isSealed,
		Permits:     permits,
		Values:      values,
		Fields:      fields,
		Methods:     methods,
		Constructor: constructor,
	}, nil
}

// parseRecord parses a record declaration
func (p *Parser) parseRecord(accessLevel string) (ast.Stmt, error) {
	p.next() // consume 'record'

	if p.curr().Tok != lexer.IDENT {
		return nil, p.errf("expected record name")
	}
	name := p.curr().Lit
	p.next()

	if !p.accept(lexer.LPAREN) {
		return nil, p.errf("expected '(' after record name")
	}

	var components []ast.RecordComponent
	if p.curr().Tok != lexer.RPAREN {
		for {
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected component name in record")
			}
			compName := p.curr().Lit
			p.next()

			var compType string
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected component type after ':'")
				}
				compType = p.curr().Lit
				p.next()
			}

			components = append(components, ast.RecordComponent{Name: compName, Type: ast.TypeFromString(compType)})

			if p.accept(lexer.COMMA) {
				continue
			}
			break
		}
	}

	if !p.accept(lexer.RPAREN) {
		return nil, p.errf("expected ')' after record components")
	}

	var methods []ast.MethodDecl

	// Parse optional methods
	for p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
		if p.curr().Tok == lexer.KW_DEF {
			method, err := p.parseMethodDecl()
			if err != nil {
				return nil, err
			}
			methods = append(methods, method)
		} else {
			// Skip unknown tokens
			p.next()
		}
	}

	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' after record body")
	}

	return &ast.RecordDecl{
		Name:        name,
		AccessLevel: accessLevel,
		Components:  components,
		Methods:     methods,
	}, nil
}

// parseTry parses a try-catch-finally statement
func (p *Parser) parseTry() (ast.Stmt, error) {
	p.next() // consume 'try'

	// Parse try block
	tryBody, err := p.parseBlockUntilKeywords([]lexer.Token{lexer.KW_CATCH, lexer.KW_FINALLY, lexer.KW_END})
	if err != nil {
		return nil, err
	}

	var catches []ast.CatchClause
	var finallyBody []ast.Stmt

	// Parse catch clauses
	for p.curr().Tok == lexer.KW_CATCH {
		p.next() // consume 'catch'

		var varName, modifier, exceptType string

		// Parse exception binding: (e: Type) or (e) or e or final e or e: Type
		if p.accept(lexer.LPAREN) {
			// Old style with parentheses: catch (e) or catch (e: Type)
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected identifier in catch clause")
			}
			varName = p.curr().Lit
			p.next()

			// Parse optional type annotation
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected type in catch clause")
				}
				exceptType = p.curr().Lit
				p.next()
			}

			if !p.accept(lexer.RPAREN) {
				return nil, p.errf("expected ')' after catch clause parameters")
			}
		} else {
			// New style without parentheses: catch e or catch final e or catch e: Type

			// Check for modifier
			if p.curr().Tok == lexer.IDENT {
				if p.curr().Lit == "final" || p.curr().Lit == "const" {
					modifier = p.curr().Lit
					p.next()
				}
			}

			// Parse variable name
			if p.curr().Tok != lexer.IDENT {
				return nil, p.errf("expected identifier in catch clause")
			}
			varName = p.curr().Lit
			p.next()

			// Parse optional type annotation
			if p.accept(lexer.COLON) {
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected type in catch clause")
				}
				exceptType = p.curr().Lit
				p.next()
			}
		}

		// Parse catch body
		catchBody, err := p.parseBlockUntilKeywords([]lexer.Token{lexer.KW_CATCH, lexer.KW_FINALLY, lexer.KW_END})
		if err != nil {
			return nil, err
		}

		catches = append(catches, ast.CatchClause{
			VarName:    varName,
			Modifier:   modifier,
			ExceptType: exceptType,
			Body:       catchBody,
		})
	}

	// Parse optional finally clause
	if p.accept(lexer.KW_FINALLY) {
		finallyBody, err = p.parseBlockUntilKeywords([]lexer.Token{lexer.KW_END})
		if err != nil {
			return nil, err
		}
	}

	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' after try statement")
	}

	return &ast.TryStmt{
		Body:    tryBody,
		Catches: catches,
		Finally: finallyBody,
	}, nil
}

// parseThrow parses a throw statement
func (p *Parser) parseThrow() (ast.Stmt, error) {
	p.next() // consume 'throw'

	// Parse the expression to throw
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	return &ast.ThrowStmt{Value: expr}, nil
}

// parseDefer parses a defer statement
func (p *Parser) parseDefer() (ast.Stmt, error) {
	pos := p.curr().Start
	p.next() // consume 'defer'

	// Parse the expression (should be a function call)
	expr, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	return &ast.DeferStmt{Call: expr, Pos: pos}, nil
}

// parseSelect parses a select statement
// select
//
//	case let x = ch.recv(): ...
//	case closed ch: ...
//
// end
func (p *Parser) parseSelect() (ast.Stmt, error) {
	pos := p.curr().Start
	p.next() // consume 'select'

	var cases []ast.SelectCase

	for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
		if p.curr().Tok != lexer.KW_CASE {
			return nil, p.errf("expected 'case' in select statement")
		}
		p.next() // consume 'case'

		// Check if it's a 'closed' case
		if p.curr().Tok == lexer.KW_CLOSED {
			p.next() // consume 'closed'

			// Parse channel expression
			channelExpr, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}

			if p.curr().Tok != lexer.COLON {
				return nil, p.errf("expected ':' after case condition")
			}
			p.next() // consume ':'

			// Parse case body until next case or end
			var body []ast.Stmt
			for p.curr().Tok != lexer.KW_CASE && p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					body = append(body, stmt)
				}
			}

			cases = append(cases, ast.SelectCase{
				IsRecv:  false,
				Channel: channelExpr,
				Body:    body,
			})
		} else {
			// Receive case: case let x = ch.recv():
			var recvVar string

			// Check for let/var variable declaration
			if p.curr().Tok == lexer.KW_LET || p.curr().Tok == lexer.KW_VAR {
				p.next() // consume 'let' or 'var'
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected variable name after let/var in select case")
				}
				recvVar = p.curr().Lit
				p.next() // consume variable name

				if p.curr().Tok != lexer.ASSIGN {
					return nil, p.errf("expected '=' in select receive case")
				}
				p.next() // consume '='
			}

			// Parse channel.recv() call
			channelExpr, err := p.parseExpr(0)
			if err != nil {
				return nil, err
			}

			if p.curr().Tok != lexer.COLON {
				return nil, p.errf("expected ':' after case condition")
			}
			p.next() // consume ':'

			// Parse case body
			var body []ast.Stmt
			for p.curr().Tok != lexer.KW_CASE && p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					body = append(body, stmt)
				}
			}

			cases = append(cases, ast.SelectCase{
				IsRecv:  true,
				RecvVar: recvVar,
				Channel: channelExpr,
				Body:    body,
			})
		}
	}

	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close select statement")
	}

	return &ast.SelectStmt{Cases: cases, Pos: pos}, nil
}

// parseSwitch parses a switch statement
// switch expr
//
//	case value1, value2: ...
//	case (x: Type): ...  // type matching
//	default: ...
//
// end
func (p *Parser) parseSwitch() (ast.Stmt, error) {
	pos := p.curr().Start
	p.next() // consume 'switch'

	// Parse the switch expression
	var expr ast.Expr
	var err error

	// Check if there's an expression (not just 'switch' followed by newline/case)
	if p.curr().Tok != lexer.KW_CASE && p.curr().Tok != lexer.KW_DEFAULT && p.curr().Tok != lexer.KW_END {
		expr, err = p.parseExpr(0)
		if err != nil {
			return nil, err
		}
		// Optional colon after switch expression (like in if statements)
		if p.curr().Tok == lexer.COLON {
			p.next() // consume ':'
		}
	}

	var cases []ast.SwitchCase
	var defaultBody []ast.Stmt

	for p.curr().Tok != lexer.EOF && p.curr().Tok != lexer.KW_END {
		if p.curr().Tok == lexer.KW_CASE {
			p.next() // consume 'case'

			var switchCase ast.SwitchCase

			// Check for type matching syntax: case (varName: TypeName):
			if p.curr().Tok == lexer.LPAREN {
				p.next() // consume '('

				// Parse variable name
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected variable name in type case")
				}
				switchCase.VarName = p.curr().Lit
				p.next()

				if p.curr().Tok != lexer.COLON {
					return nil, p.errf("expected ':' after variable name in type case")
				}
				p.next() // consume ':'

				// Parse type name
				if p.curr().Tok != lexer.IDENT {
					return nil, p.errf("expected type name in type case")
				}
				switchCase.TypeName = p.curr().Lit
				p.next()

				if p.curr().Tok != lexer.RPAREN {
					return nil, p.errf("expected ')' to close type case")
				}
				p.next() // consume ')'
			} else {
				// Parse value(s) for this case (can be multiple comma-separated values)
				for {
					value, err := p.parseExpr(0)
					if err != nil {
						return nil, err
					}
					switchCase.Values = append(switchCase.Values, value)

					// Check for comma (multiple values in one case)
					if p.curr().Tok == lexer.COMMA {
						p.next() // consume ','
						continue
					}
					break
				}
			}

			if p.curr().Tok != lexer.COLON {
				return nil, p.errf("expected ':' after case value(s)")
			}
			p.next() // consume ':'

			// Parse case body until next case, default, or end
			for p.curr().Tok != lexer.KW_CASE && p.curr().Tok != lexer.KW_DEFAULT && p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					switchCase.Body = append(switchCase.Body, stmt)
				}
			}

			cases = append(cases, switchCase)

		} else if p.curr().Tok == lexer.KW_DEFAULT {
			p.next() // consume 'default'

			if p.curr().Tok != lexer.COLON {
				return nil, p.errf("expected ':' after default")
			}
			p.next() // consume ':'

			// Parse default body until end
			for p.curr().Tok != lexer.KW_END && p.curr().Tok != lexer.EOF {
				stmt, err := p.parseStmt()
				if err != nil {
					return nil, err
				}
				if stmt != nil {
					defaultBody = append(defaultBody, stmt)
				}
			}
		} else {
			return nil, p.errf("expected 'case' or 'default' in switch statement")
		}
	}

	if !p.accept(lexer.KW_END) {
		return nil, p.errf("expected 'end' to close switch statement")
	}

	return &ast.SwitchStmt{
		Expr:    expr,
		Cases:   cases,
		Default: defaultBody,
		Pos:     pos,
	}, nil
}

// parseBlockUntilKeywords parses statements until one of the specified keywords is found
func (p *Parser) parseBlockUntilKeywords(keywords []lexer.Token) ([]ast.Stmt, error) {
	var stmts []ast.Stmt

	for p.curr().Tok != lexer.EOF {
		// Check if current token matches any of the terminating keywords
		found := false
		for _, kw := range keywords {
			if p.curr().Tok == kw {
				found = true
				break
			}
		}
		if found {
			break
		}

		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

// tryParseGenericTypeParams attempts to parse generic type parameters like <Int> or <String, Int> or <? extends Number>
// Returns nil if this is not a generic type (likely a comparison operator instead)
func (p *Parser) tryParseGenericTypeParams() ([]ast.TypeParam, error) {
	if p.curr().Tok != lexer.LT {
		return nil, nil
	}

	// Save position in case we need to backtrack
	saved := p.pos

	p.next() // consume '<'

	var typeParams []ast.TypeParam

	// Parse first type parameter
	param, err := p.parseTypeParam()
	if err != nil {
		// Error during parsing, backtrack
		p.pos = saved
		return nil, nil
	}
	if param == nil {
		// Not a valid type parameter, backtrack
		p.pos = saved
		return nil, nil
	}

	typeParams = append(typeParams, *param)

	// Parse additional type parameters
	for p.curr().Tok == lexer.COMMA {
		p.next() // consume ','
		param, err := p.parseTypeParam()
		if err != nil {
			// Invalid generic type syntax, backtrack
			p.pos = saved
			return nil, nil
		}
		if param == nil {
			// Not a valid type parameter, backtrack
			p.pos = saved
			return nil, nil
		}
		typeParams = append(typeParams, *param)
	}

	// Must have closing '>'
	if p.curr().Tok != lexer.GT {
		// Not a generic type, backtrack
		p.pos = saved
		return nil, nil
	}

	p.next() // consume '>'
	return typeParams, nil
}

// parseTypeParam parses a single type parameter, which can be:
// - A concrete type: Int, String, etc.
// - A variance-annotated type: in T, out T
// - An unbounded wildcard: ?
// - An upper-bounded wildcard: ? extends Number
// - A lower-bounded wildcard: ? super Integer
func (p *Parser) parseTypeParam() (*ast.TypeParam, error) {
	// Check for variance annotations first (in/out)
	variance := ""
	if p.curr().Tok == lexer.KW_IN {
		variance = "in"
		p.next() // consume 'in'
	} else if p.curr().Tok == lexer.KW_OUT {
		variance = "out"
		p.next() // consume 'out'
	}

	if p.curr().Tok == lexer.QUESTION {
		// Wildcard type
		p.next() // consume '?'

		// Check for bounds
		if p.curr().Tok == lexer.KW_EXTENDS {
			// Upper bound: ? extends T
			p.next() // consume 'extends'
			if p.curr().Tok != lexer.IDENT {
				return nil, nil // Invalid syntax
			}
			bound := p.curr().Lit
			p.next()
			return &ast.TypeParam{
				IsWildcard:   true,
				WildcardKind: "extends",
				Bounds:       []string{bound},
				Variance:     variance,
			}, nil
		} else if p.curr().Tok == lexer.KW_SUPER {
			// Lower bound: ? super T
			p.next() // consume 'super'
			if p.curr().Tok != lexer.IDENT {
				return nil, nil // Invalid syntax
			}
			bound := p.curr().Lit
			p.next()
			return &ast.TypeParam{
				IsWildcard:   true,
				WildcardKind: "super",
				Bounds:       []string{bound},
				Variance:     variance,
			}, nil
		} else {
			// Unbounded wildcard: ?
			return &ast.TypeParam{
				IsWildcard:   true,
				WildcardKind: "unbounded",
				Variance:     variance,
			}, nil
		}
	} else if p.curr().Tok == lexer.IDENT {
		// Concrete type (possibly with variance and/or variadic and/or bounds)
		typeName := p.curr().Lit
		p.next()

		// Check for variadic marker (...)
		isVariadic := false
		if p.curr().Tok == lexer.ELLIPSIS {
			isVariadic = true
			p.next() // consume '...'
		}

		// Check for bounds (extends or super)
		var bounds []string
		wildcardKind := ""

		if p.curr().Tok == lexer.KW_EXTENDS {
			// T extends Animal
			wildcardKind = "extends"
			p.next() // consume 'extends'
			if p.curr().Tok == lexer.IDENT {
				bounds = append(bounds, p.curr().Lit)
				p.next()
			}
		} else if p.curr().Tok == lexer.KW_SUPER {
			// T super SomeType
			wildcardKind = "super"
			p.next() // consume 'super'
			if p.curr().Tok == lexer.IDENT {
				bounds = append(bounds, p.curr().Lit)
				p.next()
			}
		}

		return &ast.TypeParam{
			IsWildcard:   false,
			IsVariadic:   isVariadic,
			Name:         typeName,
			Variance:     variance,
			WildcardKind: wildcardKind,
			Bounds:       bounds,
		}, nil
	}

	// Not a valid type parameter
	return nil, nil
}

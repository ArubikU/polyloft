package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/ArubikU/polyloft/internal/ast"
)

// Package lexer converts raw source into a stream of tokens.

// Lexer implements a minimal state machine to scan tokens.
type Lexer struct{}

func (l *Lexer) Scan(src []byte) []Item {
	var items []Item
	var off, line, col = 0, 1, 1
	add := func(tok Token, lit string, start ast.Position, end ast.Position) {
		items = append(items, Item{Tok: tok, Lit: lit, Start: start, End: end})
	}
	for off < len(src) {
		r, size := utf8.DecodeRune(src[off:])
		if r == utf8.RuneError && size == 1 {
			// invalid byte
			start := ast.Position{Offset: off, Line: line, Col: col}
			off++
			col++
			add(ILLEGAL, string(src[off-1:off]), start, ast.Position{Offset: off, Line: line, Col: col})
			continue
		}

		// Skip whitespace
		if unicode.IsSpace(r) {
			if r == '\n' {
				line++
				col = 1
			} else {
				col++
			}
			off += size
			continue
		}

		// Line comment //...
		if r == '/' && off+1 < len(src) && src[off+1] == '/' {
			for off < len(src) && src[off] != '\n' {
				off++
			}
			continue
		}
		// Block comment /* ... */
		if r == '/' && off+1 < len(src) && src[off+1] == '*' {
			off += 2
			for off+1 < len(src) && !(src[off] == '*' && src[off+1] == '/') {
				if src[off] == '\n' {
					line++
					col = 1
				} else {
					col++
				}
				off++
			}
			if off+1 < len(src) {
				off += 2
				col += 2
			}
			continue
		}

		start := ast.Position{Offset: off, Line: line, Col: col}

		// Ident or keyword (including special identifiers starting with $)
		if unicode.IsLetter(r) || r == '_' || r == '$' {
			i := off + size
			ccol := col + 1
			for i < len(src) {
				rr, sz := utf8.DecodeRune(src[i:])
				if !(unicode.IsLetter(rr) || unicode.IsDigit(rr) || rr == '_') {
					break
				}
				i += sz
				ccol++
			}
			lit := string(src[off:i])
			tok := IDENT
			if kw, ok := keywords[lit]; ok {
				tok = kw
			}
			add(tok, lit, start, ast.Position{Offset: i, Line: line, Col: ccol})
			col = ccol
			off = i
			continue
		}

		// Number (int, float, hex, binary with proper type detection)
		if unicode.IsDigit(r) {
			i := off + size
			ccol := col + 1

			// Check for hex (0x) or binary (0b) prefixes
			if r == '0' && i < len(src) {
				next := src[i]
				// Hexadecimal: 0x or 0X
				if next == 'x' || next == 'X' {
					i++
					ccol++
					// Scan hex digits (0-9, a-f, A-F) and underscores
					for i < len(src) {
						rr, sz := utf8.DecodeRune(src[i:])
						if rr == '_' {
							i += sz
							ccol++
							continue
						}
						if !((rr >= '0' && rr <= '9') || (rr >= 'a' && rr <= 'f') || (rr >= 'A' && rr <= 'F')) {
							break
						}
						i += sz
						ccol++
					}
					add(HEX, string(src[off:i]), start, ast.Position{Offset: i, Line: line, Col: ccol})
					col = ccol
					off = i
					continue
				}
				// Binary: 0b or 0B
				if next == 'b' || next == 'B' {
					i++
					ccol++
					// Scan binary digits (0-1) and underscores
					for i < len(src) {
						rr, sz := utf8.DecodeRune(src[i:])
						if rr == '_' {
							i += sz
							ccol++
							continue
						}
						if rr != '0' && rr != '1' {
							break
						}
						i += sz
						ccol++
					}
					add(BYTES, string(src[off:i]), start, ast.Position{Offset: i, Line: line, Col: ccol})
					col = ccol
					off = i
					continue
				}
			}

			// Regular decimal number
			dot := false
			// Scan digits, underscores (as separators), and optional decimal point
			for i < len(src) {
				rr, sz := utf8.DecodeRune(src[i:])
				if rr == '.' && !dot {
					// Check if this is part of ... operator
					if i+2 < len(src) && src[i+1] == '.' && src[i+2] == '.' {
						// This is ..., don't consume it as decimal point
						break
					}
					dot = true
					i += sz
					ccol++
					continue
				}
				if rr == '_' {
					// Allow underscore as digit separator (like Python)
					i += sz
					ccol++
					continue
				}
				if !unicode.IsDigit(rr) {
					break
				}
				i += sz
				ccol++
			}

			// Check for 'f' suffix for float
			var token Token
			if i < len(src) && src[i] == 'f' {
				token = FLOAT // ends with 'f' -> float
				i++
				ccol++
			} else if dot {
				token = FLOAT // has decimal point -> float
			} else {
				token = INT // no decimal, no suffix -> int
			}

			add(token, string(src[off:i]), start, ast.Position{Offset: i, Line: line, Col: ccol})
			col = ccol
			off = i
			continue
		}

		// String "..." with simple escapes \" and \n
		if r == '"' {
			i := off + size
			ccol := col + 1
			for i < len(src) {
				rr, sz := utf8.DecodeRune(src[i:])
				if rr == '\\' {
					i += sz
					ccol++
					if i < len(src) { // skip escaped char
						_, sz2 := utf8.DecodeRune(src[i:])
						i += sz2
						ccol++
					}
					continue
				}
				if rr == '"' {
					i += sz
					ccol++
					break
				}
				if rr == '\n' {
					line++
					ccol = 1
				} else {
					ccol++
				}
				i += sz
			}
			add(STRING, string(src[off:i]), start, ast.Position{Offset: i, Line: line, Col: ccol})
			col = ccol
			off = i
			continue
		}

		// String '...' with simple escapes \' and \n
		if r == '\'' {
			i := off + size
			ccol := col + 1
			for i < len(src) {
				rr, sz := utf8.DecodeRune(src[i:])
				if rr == '\\' {
					i += sz
					ccol++
					if i < len(src) { // skip escaped char
						_, sz2 := utf8.DecodeRune(src[i:])
						i += sz2
						ccol++
					}
					continue
				}
				if rr == '\'' {
					i += sz
					ccol++
					break
				}
				if rr == '\n' {
					line++
					ccol = 1
				} else {
					ccol++
				}
				i += sz
			}
			add(STRING, string(src[off:i]), start, ast.Position{Offset: i, Line: line, Col: ccol})
			col = ccol
			off = i
			continue
		}

		// Three-char operators
		if off+2 < len(src) {
			switch string(src[off : off+3]) {
			case "...":
				add(ELLIPSIS, "...", start, ast.Position{Offset: off + 3, Line: line, Col: col + 3})
				off += 3
				col += 3
				continue
			}
		}

		// Two-char operators
		if off+1 < len(src) {
			switch string(src[off : off+2]) {
			case ":=":
				add(COLONASSIGN, ":=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "+=":
				add(PLUS_ASSIGN, "+=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "-=":
				add(MINUS_ASSIGN, "-=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "*=":
				add(STAR_ASSIGN, "*=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "/=":
				add(SLASH_ASSIGN, "/=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "==":
				add(EQ, "==", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "!=":
				add(NEQ, "!=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "<=":
				add(LTE, "<=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case ">=":
				add(GTE, ">=", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "&&":
				add(AND, "&&", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "||":
				add(OR, "||", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "=>":
				add(ARROW, "=>", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			case "->":
				add(RARROW, "->", start, ast.Position{Offset: off + 2, Line: line, Col: col + 2})
				off += 2
				col += 2
				continue
			}
		}

		// Single-char tokens
		switch r {
		case '=':
			add(ASSIGN, "=", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '+':
			add(PLUS, "+", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '-':
			add(MINUS, "-", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '*':
			add(STAR, "*", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '/':
			add(SLASH, "/", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '%':
			add(PERCENT, "%", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '<':
			add(LT, "<", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '>':
			add(GT, ">", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '!':
			add(NOT, "!", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case ',':
			add(COMMA, ",", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case ':':
			add(COLON, ":", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case ';':
			add(SEMI, ";", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '(':
			add(LPAREN, "(", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case ')':
			add(RPAREN, ")", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '{':
			add(LBRACE, "{", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '}':
			add(RBRACE, "}", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '[':
			add(LBRACK, "[", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case ']':
			add(RBRACK, "]", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '.':
			add(DOT, ".", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '@':
			add(AT, "@", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '?':
			add(QUESTION, "?", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
		case '|':
			// Check if it's || (OR) or single | (PIPE)
			if off+1 < len(src) && src[off+1] == '|' {
				// This is already handled in the two-char operators section, so this should not reach here
				add(ILLEGAL, string(r), start, ast.Position{Offset: off + size, Line: line, Col: col + 1})
			} else {
				add(PIPE, "|", start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
			}
		default:
			add(ILLEGAL, string(r), start, ast.Position{Offset: off + 1, Line: line, Col: col + 1})
			size = utf8.RuneLen(r)
		}
		off += size
		col++
	}

	items = append(items, Item{Tok: EOF, Start: ast.Position{Offset: off, Line: line, Col: col}, End: ast.Position{Offset: off, Line: line, Col: col}})
	return items
}

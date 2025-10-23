package lexer

import (
	"fmt"

	"github.com/ArubikU/polyloft/internal/ast"
)

// Token identifies the kind of lexical item.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	// Literals/identifiers
	IDENT  // foo, Bar, _baz
	NUMBER // 10, 3.14, -5, 2.5e10
	INT    // 10, 42, -5
	FLOAT  // 10.0f, 3.14f
	STRING
	INTERPOLATED_STRING // "text #{expr} more text"

	// Keywords
	KW_VAR
	KW_LET
	KW_CONST
	KW_FINAL
	KW_PUBLIC
	KW_PRIVATE
	KW_PROTECTED
	KW_STATIC
	KW_DEF
	KW_INTERFACE
	KW_CLASS
	KW_IMPORT
	KW_IMPLEMENTS
	KW_ABSTRACT
	KW_SEALED
	KW_RETURN
	KW_TRUE
	KW_FALSE
	KW_NIL
	KW_THREAD
	KW_SPAWN
	KW_JOIN
	KW_IF
	KW_ELIF
	KW_ELSE
	KW_FOR
	KW_IN
	KW_BREAK
	KW_CONTINUE
	KW_LOOP
	KW_END
	KW_DO
	KW_INSTANCEOF
	KW_THIS
	KW_SUPER
	KW_ENUM
	KW_RECORD
	KW_TRY
	KW_CATCH
	KW_FINALLY
	KW_THROW
	KW_DEFER
	KW_CHANNEL
	KW_SELECT
	KW_SWITCH
	KW_CASE
	KW_DEFAULT
	KW_CLOSED
	KW_WHERE
	KW_EXTENDS
	KW_OUT

	// Operators and delimiters
	ASSIGN      // =
	PLUS        // +
	MINUS       // -
	STAR        // *
	SLASH       // /
	PERCENT     // %
	EQ          // ==
	NEQ         // !=
	LT          // <
	LTE         // <=
	GT          // >
	GTE         // >=
	AND         // &&
	OR          // ||
	NOT         // !
	ARROW       // =>
	RARROW      // ->
	COLONASSIGN // :=

	COMMA    // ,
	COLON    // :
	SEMI     // ;
	QUESTION // ? (for ternary operator)

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACK   // [
	RBRACK   // ]
	DOT      // .
	ELLIPSIS // ... (for variadic parameters)
	AT       // @ (for annotations)
	PIPE     // | (for union types)
)

var keywords = map[string]Token{
	"var":        KW_VAR,
	"let":        KW_LET,
	"const":      KW_CONST,
	"public":     KW_PUBLIC,
	"pub":        KW_PUBLIC,
	"private":    KW_PRIVATE,
	"priv":       KW_PRIVATE,
	"protected":  KW_PROTECTED,
	"prot":       KW_PROTECTED,
	"static":     KW_STATIC,
	"def":        KW_DEF,
	"interface":  KW_INTERFACE,
	"class":      KW_CLASS,
	"import":     KW_IMPORT,
	"implements": KW_IMPLEMENTS,
	"abstract":   KW_ABSTRACT,
	"sealed":     KW_SEALED,
	"thread":     KW_THREAD,
	"spawn":      KW_SPAWN,
	"join":       KW_JOIN,
	"return":     KW_RETURN,
	"true":       KW_TRUE,
	"false":      KW_FALSE,
	"nil":        KW_NIL,
	"if":         KW_IF,
	"elif":       KW_ELIF,
	"else":       KW_ELSE,
	"for":        KW_FOR,
	"in":         KW_IN,
	"break":      KW_BREAK,
	"continue":   KW_CONTINUE,
	"loop":       KW_LOOP,
	"end":        KW_END,
	"do":         KW_DO,
	"instanceof": KW_INSTANCEOF,
	"this":       KW_THIS,
	"super":      KW_SUPER,
	"enum":       KW_ENUM,
	"record":     KW_RECORD,
	"try":        KW_TRY,
	"catch":      KW_CATCH,
	"finally":    KW_FINALLY,
	"throw":      KW_THROW,
	"final":      KW_FINAL,
	"defer":      KW_DEFER,
	"channel":    KW_CHANNEL,
	"select":     KW_SELECT,
	"switch":     KW_SWITCH,
	"case":       KW_CASE,
	"default":    KW_DEFAULT,
	"closed":     KW_CLOSED,
	"where":      KW_WHERE,
	"extends":    KW_EXTENDS,
	"out":        KW_OUT,
}

// Item represents a scanned token with its literal text and position.
type Item struct {
	Tok   Token
	Lit   string
	Start ast.Position
	End   ast.Position
}

// TokenName returns a human-readable name for a token type.
func TokenName(tok Token) string {
	switch tok {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "identifier"
	case NUMBER:
		return "number"
	case INT:
		return "integer"
	case FLOAT:
		return "float"
	case STRING:
		return "string"
	case INTERPOLATED_STRING:
		return "interpolated string"
	case KW_VAR:
		return "keyword 'var'"
	case KW_LET:
		return "keyword 'let'"
	case KW_CONST:
		return "keyword 'const'"
	case KW_FINAL:
		return "keyword 'final'"
	case KW_PUBLIC:
		return "keyword 'public'"
	case KW_PRIVATE:
		return "keyword 'private'"
	case KW_PROTECTED:
		return "keyword 'protected'"
	case KW_STATIC:
		return "keyword 'static'"
	case KW_DEF:
		return "keyword 'def'"
	case KW_INTERFACE:
		return "keyword 'interface'"
	case KW_CLASS:
		return "keyword 'class'"
	case KW_IMPORT:
		return "keyword 'import'"
	case KW_IMPLEMENTS:
		return "keyword 'implements'"
	case KW_ABSTRACT:
		return "keyword 'abstract'"
	case KW_RETURN:
		return "keyword 'return'"
	case KW_TRUE:
		return "keyword 'true'"
	case KW_FALSE:
		return "keyword 'false'"
	case KW_NIL:
		return "keyword 'nil'"
	case KW_THREAD:
		return "keyword 'thread'"
	case KW_SPAWN:
		return "keyword 'spawn'"
	case KW_JOIN:
		return "keyword 'join'"
	case KW_IF:
		return "keyword 'if'"
	case KW_ELIF:
		return "keyword 'elif'"
	case KW_ELSE:
		return "keyword 'else'"
	case KW_FOR:
		return "keyword 'for'"
	case KW_IN:
		return "keyword 'in'"
	case KW_BREAK:
		return "keyword 'break'"
	case KW_CONTINUE:
		return "keyword 'continue'"
	case KW_LOOP:
		return "keyword 'loop'"
	case KW_END:
		return "keyword 'end'"
	case KW_DO:
		return "keyword 'do'"
	case KW_INSTANCEOF:
		return "keyword 'instanceof'"
	case KW_THIS:
		return "keyword 'this'"
	case KW_SUPER:
		return "keyword 'super'"
	case KW_ENUM:
		return "keyword 'enum'"
	case KW_RECORD:
		return "keyword 'record'"
	case KW_TRY:
		return "keyword 'try'"
	case KW_CATCH:
		return "keyword 'catch'"
	case KW_FINALLY:
		return "keyword 'finally'"
	case KW_THROW:
		return "keyword 'throw'"
	case KW_DEFER:
		return "keyword 'defer'"
	case KW_CHANNEL:
		return "keyword 'channel'"
	case KW_SELECT:
		return "keyword 'select'"
	case KW_SWITCH:
		return "keyword 'switch'"
	case KW_CASE:
		return "keyword 'case'"
	case KW_DEFAULT:
		return "keyword 'default'"
	case KW_CLOSED:
		return "keyword 'closed'"
	case KW_WHERE:
		return "keyword 'where'"
	case KW_EXTENDS:
		return "keyword 'extends'"
	case KW_OUT:
		return "keyword 'out'"
	case ASSIGN:
		return "'='"
	case PLUS:
		return "'+'"
	case MINUS:
		return "'-'"
	case STAR:
		return "'*'"
	case SLASH:
		return "'/'"
	case PERCENT:
		return "'%'"
	case EQ:
		return "'=='"
	case NEQ:
		return "'!='"
	case LT:
		return "'<'"
	case LTE:
		return "'<='"
	case GT:
		return "'>'"
	case GTE:
		return "'>='"
	case AND:
		return "'&&'"
	case OR:
		return "'||'"
	case NOT:
		return "'!'"
	case ARROW:
		return "'=>'"
	case RARROW:
		return "'->'"
	case COLONASSIGN:
		return "':='"
	case COMMA:
		return "','"
	case COLON:
		return "':'"
	case SEMI:
		return "';'"
	case QUESTION:
		return "'?'"
	case LPAREN:
		return "'('"
	case RPAREN:
		return "')'"
	case LBRACE:
		return "'{'"
	case RBRACE:
		return "'}'"
	case LBRACK:
		return "'['"
	case RBRACK:
		return "']'"
	case DOT:
		return "'.'"
	case ELLIPSIS:
		return "'...'"
	case AT:
		return "'@'"
	case PIPE:
		return "'|'"
	default:
		return fmt.Sprintf("unknown token (%d)", tok)
	}
}

// FormatToken returns a formatted string for a token including its literal value if available.
func FormatToken(item Item) string {
	name := TokenName(item.Tok)
	if item.Lit != "" && item.Tok != EOF && item.Tok != ILLEGAL {
		return fmt.Sprintf("%s (%q)", name, item.Lit)
	}
	return name
}

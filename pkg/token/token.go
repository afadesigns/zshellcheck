// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package token

type Type string

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
	// EndLine records the source line where the token ended.
	// Most tokens are single-line (EndLine == Line), but
	// multi-line strings (`'…\n…\n…'`, `"…\n…"`) and heredoc
	// bodies span lines; the parser's "same-line" arg-gathering
	// checks should consult EndLine of the just-consumed token
	// when comparing against peek.Line. Zero means "use Line".
	EndLine           int
	HasPrecedingSpace bool
	// HasPrecedingContinuation is set when the lexer consumed a
	// `\<NL>` pair immediately before this token. The Line / Column
	// fields still refer to the physical source position so error
	// messages stay accurate, but parsers can treat the token as if
	// it were on the previous line — matching how Zsh collapses
	// line-continuations into one logical command.
	HasPrecedingContinuation bool
	// ClosesDollarBrace is set on an RBRACE token that closes a
	// `${ … }` parameter expansion (rather than a regular brace
	// block). The parser uses it to distinguish `${X} }` — two
	// RBRACE tokens whose semantics differ — when deciding whether
	// a `}` is a block terminator.
	ClosesDollarBrace bool
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "hello world"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	INC      = "++"
	DEC      = "--"
	PLUSEQ   = "+="

	GT       = ">"
	LT       = "<"
	GTGT     = ">>"
	LTLT     = "<<"
	LTLTLT   = "<<<"
	GTAMP    = ">&"
	LTAMP    = "<&"
	EQ       = "=="
	NotEq    = "!="
	EQTILDE  = "=~"
	QUESTION = "?"

	// Arithmetic comparison operators (Zsh/Bash [[ ... ]])
	EQ_NUM = "-eq"
	NE_NUM = "-ne"
	LT_NUM = "-lt"
	LE_NUM = "-le"
	GT_NUM = "-gt"
	GE_NUM = "-ge"

	// Process Substitution / Array Assignment
	LT_LPAREN = "<("
	GT_LPAREN = ">("
	EQ_LPAREN = "=("

	// Delimiters
	COMMA        = ","
	SEMICOLON    = ";"
	COLON        = ":"
	LPAREN       = "("
	RPAREN       = ")"
	LBRACE       = "{"
	RBRACE       = "}"
	LBRACKET     = "["
	RBRACKET     = "]"
	LDBRACKET    = "[["
	RDBRACKET    = "]]"
	DoubleLparen = "(("
	DoubleRparen = "))"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	If       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	THEN     = "THEN"
	Fi       = "FI"
	FOR      = "FOR"
	WHILE    = "WHILE"
	DO       = "DO"
	DONE     = "DONE"
	IN       = "IN"
	CASE     = "CASE"
	ESAC     = "ESAC"
	ELIF     = "ELIF"
	SELECT   = "SELECT"
	COPROC   = "COPROC"
	TYPESET  = "TYPESET"
	DECLARE  = "DECLARE"

	// Zsh-specific tokens
	DOLLAR        = "$"
	DollarLbrace  = "${"
	DOLLAR_LPAREN = "$("
	VARIABLE      = "VARIABLE"
	HASH          = "#"
	AMPERSAND     = "&"
	PIPE          = "|"
	BACKTICK      = "`"
	TILDE         = "~"
	CARET         = "^"
	PERCENT       = "%"
	DOT           = "."
	SHEBANG       = "#!"

	// Zsh-specific operators
	AND   = "&&"
	OR    = "||"
	DSEMI = ";;"

	// Zsh-specific delimiters
	LARRAY = "("
	RARRAY = ")"
)

var keywords = map[string]Type{
	"function": FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       If,
	"else":     ELSE,
	"return":   RETURN,
	"then":     THEN,
	"fi":       Fi,
	"for":      FOR,
	"while":    WHILE,
	"do":       DO,
	"done":     DONE,
	"in":       IN,
	"case":     CASE,
	"esac":     ESAC,
	"elif":     ELIF,
	"select":   SELECT,
	"coproc":   COPROC,
	"typeset":  TYPESET,
	"declare":  DECLARE,
	"-eq":      EQ_NUM,
	"-ne":      NE_NUM,
	"-lt":      LT_NUM,
	"-le":      LE_NUM,
	"-gt":      GT_NUM,
	"-ge":      GE_NUM,
	"/":        SLASH,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

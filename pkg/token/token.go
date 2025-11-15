package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
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

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA         = ","
	SEMICOLON     = ";"
	COLON         = ":"
	LPAREN        = "("
	RPAREN        = ")"
	LBRACE        = "{"
	RBRACE        = "}"
	LBRACKET      = "["
	RBRACKET      = "]"
	LDBRACKET     = "[["
	RDBRACKET     = "]]"
	DOUBLE_LPAREN = "(("

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	THEN     = "THEN"
	FI       = "FI"
	FOR      = "FOR"
	DONE     = "DONE"

	// Zsh-specific tokens (initial)
	DOLLAR        = "$"
	DOLLAR_LBRACE = "${"
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

	// Zsh-specific operators (initial)
	AND = "&&"
	OR  = "||"

	// Zsh-specific delimiters (initial)
	LARRAY = "("
	RARRAY = ")"
)

var keywords = map[string]TokenType{
	"function": FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"then":     THEN,
	"fi":       FI,
	"for":      FOR,
	"done":     DONE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

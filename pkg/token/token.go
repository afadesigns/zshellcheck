package token

type Type string

type Token struct {
	Type              Type
	Literal           string
	Line              int
	Column            int
	HasPrecedingSpace bool
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

	EQ    = "=="
	NotEq = "!="

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

	// Zsh-specific tokens (initial)
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

	// Zsh-specific operators (initial)
	AND = "&&"
	OR  = "||"
	DSEMI = ";;"

	// Zsh-specific delimiters (initial)
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
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

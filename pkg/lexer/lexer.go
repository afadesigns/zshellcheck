package lexer

import (
	"github.com/afadesigns/zshellcheck/pkg/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0} // column is 0-indexed initially
	l.readChar()                                  // This initializes l.ch and l.position, setting column to 1
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == 10 { // \n
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	hasSpace := l.skipWhitespace()
	// Store hasSpace but set it on tok later because tok is re-assigned below.

	if l.ch == '#' {
		if l.peekChar() == '!' {
			start := l.position
			for l.ch != 10 && l.ch != 0 { // \n
				l.readChar()
			}
			literal := l.input[start:l.position]
			return token.Token{Type: token.SHEBANG, Literal: literal, Line: l.line, Column: l.column}
		}

		if hasSpace || l.column == 1 {
			l.skipComment()
			return l.NextToken()
		}
	}

	switch l.ch {
	case '#':
		tok = newToken(token.HASH, l.ch, l.line, l.column)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Column: l.column}
		} else if l.peekChar() == '~' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQTILDE, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
		}
	case ';':
		if l.peekChar() == ';' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DSEMI, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.SEMICOLON, l.ch, l.line, l.column)
		}
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.column)
	case '?':
		tok = newToken(token.QUESTION, l.ch, l.line, l.column)
	case '(':
		if l.peekChar() == '(' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DoubleLparen, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.LPAREN, l.ch, l.line, l.column)
		}
	case ')':
		if l.peekChar() == ')' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DoubleRparen, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.RPAREN, l.ch, l.line, l.column)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.column)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: l.line, Column: l.column}
		} else if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.PLUS, l.ch, l.line, l.column)
		}
	case '-':
		if isLetter(l.peekChar()) {
			savedLexer := *l
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			if tokType != token.IDENT {
				tok.Type = tokType
				tok.Literal = literal
				tok.Line = savedLexer.line
				tok.Column = savedLexer.column
				tok.HasPrecedingSpace = hasSpace
				return tok
			}
			*l = savedLexer
		}

		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.MINUS, l.ch, l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NotEq, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.BANG, l.ch, l.line, l.column)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.line, l.column)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.line, l.column)
	case '<':
		switch l.peekChar() {
		case '<':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTLT, Literal: literal, Line: l.line, Column: l.column}
		case '&':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTAMP, Literal: literal, Line: l.line, Column: l.column}
		default:
			tok = newToken(token.LT, l.ch, l.line, l.column)
		}
	case '>':
		switch l.peekChar() {
		case '>':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTGT, Literal: literal, Line: l.line, Column: l.column}
		case '&':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTAMP, Literal: literal, Line: l.line, Column: l.column}
		default:
			tok = newToken(token.GT, l.ch, l.line, l.column)
		}
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.column)
	case '[':
		if l.peekChar() == '[' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LDBRACKET, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
		}
	case ']':
		if l.peekChar() == ']' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RDBRACKET, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
		}
	case '$':
		switch l.peekChar() {
		case '{':
			tok.Type = token.DollarLbrace
			tok.Literal = "${"
			tok.Line = l.line
			tok.Column = l.column
			l.readChar()
		case '(':
			tok.Type = token.DOLLAR_LPAREN
			tok.Literal = "$("
			tok.Line = l.line
			tok.Column = l.column
			l.readChar()
		default:
			if isLetter(l.peekChar()) {
				col := l.column
				l.readChar() // consume '$'
				tok.Type = token.VARIABLE
				tok.Literal = "$" + l.readIdentifier()
				tok.Line = l.line
				tok.Column = col
				tok.HasPrecedingSpace = hasSpace
				return tok
			}
			tok = newToken(token.DOLLAR, l.ch, l.line, l.column)
		}

	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.AMPERSAND, l.ch, l.line, l.column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.PIPE, l.ch, l.line, l.column)
		}
	case '`':
		tok = newToken(token.BACKTICK, l.ch, l.line, l.column)
	case '~':
		tok = newToken(token.TILDE, l.ch, l.line, l.column)
	case '^':
		tok = newToken(token.CARET, l.ch, l.line, l.column)
	case '%':
		tok = newToken(token.PERCENT, l.ch, l.line, l.column)
	case '.':
		tok = newToken(token.DOT, l.ch, l.line, l.column)
	case '"':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readString('"')
	case '\'':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readString('\'')
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		switch {
		case isLetter(l.ch):
			line, col := l.line, l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = col
			tok.HasPrecedingSpace = hasSpace
			return tok
		case isDigit(l.ch):
			line, col := l.line, l.column
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			tok.Line = line
			tok.Column = col
			tok.HasPrecedingSpace = hasSpace
			return tok
		default:
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	tok.HasPrecedingSpace = hasSpace
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString(quote byte) string {
	position := l.position // include opening quote
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // skip escaped char
		}
	}
	return l.input[position : l.position+1] // include closing quote
}

func (l *Lexer) skipWhitespace() bool {
	skipped := false
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		skipped = true
		if l.ch == 10 { // \n
			l.line++
			l.column = 0
		}
		l.readChar()
	}
	return skipped
}

func (l *Lexer) skipComment() {
	for l.ch != 10 && l.ch != 0 { // \n
		l.readChar()
	}
}

func newToken(tokenType token.Type, ch byte, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' ||
		ch == '/' || ch == '.' || ch == '@' || ch == ':' || ch == '~'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

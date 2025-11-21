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
	l.readChar() // This initializes l.ch and l.position, setting column to 1
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == '\n' {
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
	tok.HasPrecedingSpace = hasSpace

	if l.ch == '#' {
		// Shebang check (must be at start of file, or maybe just check sequence #!)
		// Usually only valid at start. But strict check: col 1, line 1?
		// For robustness, if #! appears, it's shebang token.
		if l.peekChar() == '!' {
			// Consume until end of line
			start := l.position
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			literal := l.input[start:l.position]
			return token.Token{Type: token.SHEBANG, Literal: literal, Line: l.line, Column: l.column}
		}

		// Comment check: only if preceded by space or start of line
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
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.line, l.column)
	case ':':
		tok = newToken(token.COLON, l.ch, l.line, l.column)
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
		} else {
			tok = newToken(token.PLUS, l.ch, l.line, l.column)
		}
	case '-':
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
		if l.peekChar() == '<' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTLT, Literal: literal, Line: l.line, Column: l.column}
		} else if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTAMP, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTGT, Literal: literal, Line: l.line, Column: l.column}
		} else if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTAMP, Literal: literal, Line: l.line, Column: l.column}
		} else {
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
		tok.Literal = l.readString()
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
			return tok
		case isDigit(l.ch):
			line, col := l.line, l.column
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			tok.Line = line
			tok.Column = col
			return tok
		default:
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readShebang() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() bool {
	skipped := false
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		skipped = true
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
	return skipped
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
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

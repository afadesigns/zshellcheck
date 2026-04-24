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

	// dbracketDepth tracks nesting of `[[ … ]]` conditional blocks.
	// Without this we cannot distinguish a conditional close (the
	// fused `]]` token) from two consecutive single-bracket closes
	// in `arr[$m[i]]`. When depth is zero the lexer emits two
	// RBRACKETs instead of RDBRACKET.
	dbracketDepth int

	// pendingContinuation is set when skipWhitespace has just consumed
	// a `\<NL>` line-continuation pair. It is read and cleared by
	// NextToken so the next emitted token carries
	// HasPrecedingContinuation, letting the parser treat it as if it
	// were on the previous line without altering the physical Line
	// used for error messages.
	pendingContinuation bool

	// parenStack records the kind of every paren-like opener that is
	// still awaiting its close. 'D' for `((` (arithmetic), 'P' for
	// plain `(` and for `$(` command substitution. The lexer fuses
	// `))` into DoubleRparen only when the innermost open context is
	// 'D'; a plain `(` inside `(( … ))` closes with a single RPAREN
	// so that `(( x = $((1+1)) + 2 ))` emits two inner RPARENs for
	// the `$(` + `(` pair and a final DoubleRparen for the outer `((`.
	parenStack []byte
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

func (l *Lexer) NextToken() (tok token.Token) {
	// skipWhitespace sets pendingContinuation when a `\<NL>` pair
	// was absorbed; stamp the flag onto the returned token via this
	// named-return defer so every early return path inherits it.
	defer func() {
		if l.pendingContinuation {
			tok.HasPrecedingContinuation = true
			l.pendingContinuation = false
		}
	}()
	hasSpace := l.skipWhitespace()

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
		switch l.peekChar() {
		case '=':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Column: l.column}
		case '~':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQTILDE, Literal: literal, Line: l.line, Column: l.column}
		case '(':
			if hasSpace {
				ch := l.ch
				l.readChar()
				literal := string(ch) + string(l.ch)
				tok = token.Token{Type: token.EQ_LPAREN, Literal: literal, Line: l.line, Column: l.column}
				l.parenStack = append(l.parenStack, 'P')
			} else {
				tok = newToken(token.ASSIGN, l.ch, l.line, l.column)
			}
		default:
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
			l.parenStack = append(l.parenStack, 'D')
		} else {
			tok = newToken(token.LPAREN, l.ch, l.line, l.column)
			l.parenStack = append(l.parenStack, 'P')
		}
	case ')':
		tok = l.readCloseParen()
	case ',':
		tok = newToken(token.COMMA, l.ch, l.line, l.column)
	case '+':
		switch l.peekChar() {
		case '+':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: l.line, Column: l.column}
		case '=':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Column: l.column}
		default:
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
		case '(':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LT_LPAREN, Literal: literal, Line: l.line, Column: l.column}
			l.parenStack = append(l.parenStack, 'P')
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
		case '(':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT_LPAREN, Literal: literal, Line: l.line, Column: l.column}
			l.parenStack = append(l.parenStack, 'P')
		case '|':
			// Zsh `>|file` and `>!file` force-clobber a file even
			// when `NO_CLOBBER` is set. The trailing `|` / `!`
			// belongs to the redirection, not to a pipeline or
			// negation that follows. Emit the pair as a plain GT
			// so parseCommandPipeline's redirection path handles
			// it unchanged — the AST form is identical to `>file`.
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT, Literal: literal, Line: l.line, Column: l.column}
		case '!':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT, Literal: literal, Line: l.line, Column: l.column}
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
			l.dbracketDepth++
		} else {
			tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
		}
	case ']':
		// `]]` only means RDBRACKET when there is a pending
		// `[[` to close. In array-subscript contexts like
		// `arr[$m[i]]` the two brackets close two independent
		// subscripts and must lex as two RBRACKET tokens.
		if l.peekChar() == ']' && l.dbracketDepth > 0 {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RDBRACKET, Literal: literal, Line: l.line, Column: l.column}
			l.dbracketDepth--
		} else {
			tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
		}
	case '$':
		if dollarTok, ok := l.readDollarToken(hasSpace); ok {
			return dollarTok
		}
		tok = newToken(token.DOLLAR, l.ch, l.line, l.column)

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
			// An identifier that happens to match a keyword but is immediately
			// followed by `=` is a flag/argument assignment (e.g. `if=foo`
			// inside `dd if=foo of=bar`), not the keyword itself. Demote it
			// to a plain identifier so the parser treats the following `=`
			// as part of the same word rather than trying to open an
			// if-statement.
			if tok.Type != token.IDENT && l.ch == '=' {
				tok.Type = token.IDENT
			}
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
		case l.ch == '\\':
			// Backslash outside a string quotes exactly one following
			// character. Zsh glob escapes (`\(`, `\)`, `\*`, `\?`,
			// etc.) surface in oh-my-zsh themes. Emit the escaped
			// character's natural token — backslash-newline is
			// already handled by skipWhitespace. For non-alphanumeric
			// escapes we emit the raw escaped char as an IDENT-style
			// word so parseCommandWord folds it into the surrounding
			// word naturally. We only do this when the next byte is
			// one of the commonly-escaped glob / shell metacharacters
			// to avoid destabilising token-aware contexts.
			if next := l.peekChar(); next == '(' || next == ')' || next == '*' ||
				next == '?' || next == '[' || next == ']' || next == '|' ||
				next == '&' || next == ';' || next == '<' || next == '>' ||
				next == '{' || next == '}' || next == '$' || next == '\\' ||
				next == '/' || next == '.' || next == '!' || next == '~' ||
				next == '^' {
				line, col := l.line, l.column
				l.readChar() // consume '\'
				tok = token.Token{
					Type:    token.IDENT,
					Literal: "\\" + string(l.ch),
					Line:    line,
					Column:  col,
				}
			} else {
				tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
			}
		default:
			tok = newToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	tok.HasPrecedingSpace = hasSpace
	return tok
}

// readCloseParen resolves a `)` to either DoubleRparen (fused `))`)
// or a plain RPAREN by consulting parenStack. `))` fuses only when
// the innermost still-open context is 'D' (a `((` opener). Plain `(`
// / `$(` / `<(` / `>(` / `=(` openers close with a single `)` so the
// inner `))` in `(( x = $((1+1)) + 2 ))` does not swallow the outer
// `((`'s closer.
func (l *Lexer) readCloseParen() token.Token {
	top := byte(0)
	if n := len(l.parenStack); n > 0 {
		top = l.parenStack[n-1]
	}
	if l.peekChar() == ')' && top == 'D' {
		ch := l.ch
		l.readChar()
		literal := string(ch) + string(l.ch)
		tok := token.Token{Type: token.DoubleRparen, Literal: literal, Line: l.line, Column: l.column}
		l.parenStack = l.parenStack[:len(l.parenStack)-1]
		return tok
	}
	tok := newToken(token.RPAREN, l.ch, l.line, l.column)
	if len(l.parenStack) > 0 {
		l.parenStack = l.parenStack[:len(l.parenStack)-1]
	}
	return tok
}

// readDollarToken dispatches the specialised forms that follow a
// leading `$`. It returns (tok, true) when it has consumed a recognised
// form — parameter expansion opener (${ or $(), ANSI-C / gettext string
// ($'…' or $"…"), a named variable ($name), or a single-character
// special parameter ($? / $@ / $$ / $_). Otherwise it returns
// (zero, false) and the caller falls back to emitting a bare DOLLAR
// token.
func (l *Lexer) readDollarToken(hasSpace bool) (token.Token, bool) {
	var tok token.Token
	switch l.peekChar() {
	case '{':
		tok.Type = token.DollarLbrace
		tok.Literal = "${"
		tok.Line = l.line
		tok.Column = l.column
		l.readChar() // consume '$'
		l.readChar() // step past the shared tail; advances past '{'
		tok.HasPrecedingSpace = hasSpace
		return tok, true
	case '(':
		tok.Type = token.DOLLAR_LPAREN
		tok.Literal = "$("
		tok.Line = l.line
		tok.Column = l.column
		l.readChar() // consume '$'
		l.readChar() // advance past '('
		// `$(` opens a command-substitution that closes with a single
		// `)`. Record it as 'P' so a nested `))` does not get fused
		// into DoubleRparen when only the inner `(` / `$(` are being
		// closed.
		l.parenStack = append(l.parenStack, 'P')
		tok.HasPrecedingSpace = hasSpace
		return tok, true
	case '\'':
		// Zsh ANSI-C quoting: $'…' processes backslash escapes like
		// \n, \t, and crucially \' for an embedded single quote.
		// Must honour escapes so `$'\''` does not terminate early.
		col := l.column
		l.readChar() // consume '$'
		body := l.readStringFlavour('\'', true)
		tok.Type = token.STRING
		tok.Literal = "$" + body
		tok.Line = l.line
		tok.Column = col
		tok.HasPrecedingSpace = hasSpace
		l.readChar() // step past the closing quote
		return tok, true
	case '"':
		// Zsh gettext quoting: $"…" marks a string for translation.
		// The payload is otherwise a regular double-quoted string.
		col := l.column
		l.readChar() // consume '$'
		body := l.readString('"')
		tok.Type = token.STRING
		tok.Literal = "$" + body
		tok.Line = l.line
		tok.Column = col
		tok.HasPrecedingSpace = hasSpace
		l.readChar() // step past the closing quote
		return tok, true
	}
	if isLetter(l.peekChar()) {
		col := l.column
		l.readChar() // consume '$'
		tok.Type = token.VARIABLE
		tok.Literal = "$" + l.readIdentifier()
		tok.Line = l.line
		tok.Column = col
		tok.HasPrecedingSpace = hasSpace
		return tok, true
	}
	// Zsh single-character special parameters: $? (exit status),
	// $@ (all positional), $$ (PID), $_ (last arg). The other
	// single-char specials ($#, $*, $!, $-) are assembled by the
	// parser from DOLLAR + the following punctuation token.
	if c := l.peekChar(); c == '?' || c == '@' || c == '$' || c == '_' {
		col := l.column
		l.readChar() // consume '$'
		tok.Type = token.VARIABLE
		tok.Literal = "$" + string(l.ch)
		tok.Line = l.line
		tok.Column = col
		tok.HasPrecedingSpace = hasSpace
		l.readChar() // consume the special char
		return tok, true
	}
	return token.Token{}, false
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
	return l.readStringFlavour(quote, quote == '"')
}

// readStringFlavour is the shared body of string lexing. When
// honourEscapes is true, `\X` mid-string consumes both bytes so
// that `\<quote>` does not terminate the string — used by double-
// quoted strings and by the `$'…'` ANSI-C form. Plain single quotes
// in Zsh never honour escapes: `\` inside `'…'` is a literal
// backslash, and the only way to embed `'` is to close and reopen
// the quotation.
func (l *Lexer) readStringFlavour(quote byte, honourEscapes bool) string {
	position := l.position // include opening quote
	// braceDepth tracks `${ … }` parameter expansions embedded in a
	// double-quoted string. Zsh suspends outer-quote termination
	// while inside `${…}`, so nested quotes like
	// `"${var="default"}"` must not split the string at the inner
	// `"`. Single quotes and ANSI-C strings never embed expansions,
	// so braceDepth only grows when honourEscapes is true.
	braceDepth := 0
	for {
		l.readChar()
		if l.ch == 0 {
			break
		}
		if l.ch == quote && braceDepth == 0 {
			break
		}
		if honourEscapes && l.ch == '\\' {
			l.readChar() // skip escaped char
			if l.ch == 0 {
				break
			}
			continue
		}
		if honourEscapes && l.ch == '$' && l.peekChar() == '{' {
			l.readChar() // consume `$`
			braceDepth++
			continue
		}
		if braceDepth > 0 {
			switch l.ch {
			case '{':
				braceDepth++
			case '}':
				braceDepth--
			}
		}
	}
	if l.ch == 0 {
		end := l.position
		if end > len(l.input) {
			end = len(l.input)
		}
		return l.input[position:end]
	}
	end := l.position + 1
	if end > len(l.input) {
		end = len(l.input)
	}
	return l.input[position:end] // include closing quote
}

func (l *Lexer) skipWhitespace() bool {
	skipped := false
	for {
		switch l.ch {
		case ' ', '\t', '\n', '\r':
			skipped = true
			l.readChar()
			continue
		case '\\':
			// Line continuation: an unquoted backslash immediately
			// followed by a newline joins the next line to the
			// current one. Skip both characters so the lexer treats
			// `cmd \<NL>arg` as `cmd arg`. Any other use of `\`
			// falls through to the regular token handler.
			if l.peekChar() == '\n' {
				l.readChar() // consume '\'
				l.readChar() // consume '\n'
				skipped = true
				l.pendingContinuation = true
				continue
			}
		}
		return skipped
	}
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

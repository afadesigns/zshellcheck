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

	// dbracketDepth is kept for historical parity with katas that
	// look at it, but the primary source of truth is now
	// bracketStack below. Incremented alongside every `[[` push and
	// decremented alongside every `]]` close.
	dbracketDepth int

	// bracketStack records every open bracket so `]]` only fuses
	// into RDBRACKET when it actually closes a `[[`. 'D' marks a
	// `[[` opener; 'B' marks a plain `[` (array subscript or a
	// glob bracket class). Without this a POSIX class inside a
	// conditional like `[[ $x == *[[:alnum:]] ]]` collapsed the
	// class's `]]` into the conditional's closer and left the
	// outer `]]` unfused.
	bracketStack []byte

	// suppressLparenFusion is set after emitting DOLLAR_LPAREN (`$(`)
	// so the next `(` is NOT fused with its peek into DoubleLparen.
	// Real code writes `$(((expr) * 2))` which is `$((` arithmetic
	// plus a nested `(expr)` grouping; without this flag the lexer
	// consumed the inner `(` into a DoubleLparen and the outer
	// arithmetic `))` never found its match.
	suppressLparenFusion bool

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

// peekAt returns the byte n positions ahead of the current reading
// position (1 == peekChar). Used by the `[` handler to look past the
// immediate peek so `[[:alnum:]]` can disambiguate the bracket class
// opener from the `[[` keyword without rewinding state.
func (l *Lexer) peekAt(n int) byte {
	idx := l.readPosition + n - 1
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

func (l *Lexer) NextToken() (tok token.Token) {
	// skipWhitespace sets pendingContinuation when a `\<NL>` pair
	// was absorbed; stamp the flag onto the returned token via this
	// named-return defer so every early return path inherits it.
	// Also clear suppressLparenFusion unless the returned token is
	// another `$(` — the flag is one-shot: the NEXT `(` (paren
	// token) should skip fusion, any other token drops the flag so
	// later `((` pairs (e.g. a fresh arithmetic open) fuse normally.
	prevSuppress := l.suppressLparenFusion
	defer func() {
		if l.pendingContinuation {
			tok.HasPrecedingContinuation = true
			l.pendingContinuation = false
		}
		// If the just-emitted token is LPAREN (we consumed the
		// suppress), or anything that doesn't open another `$(`,
		// clear the flag. DOLLAR_LPAREN re-sets it explicitly.
		if prevSuppress && tok.Type != token.DOLLAR_LPAREN {
			l.suppressLparenFusion = false
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
		if l.peekChar() == '(' && !l.suppressLparenFusion {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DoubleLparen, Literal: literal, Line: l.line, Column: l.column}
			l.parenStack = append(l.parenStack, 'D')
		} else {
			tok = newToken(token.LPAREN, l.ch, l.line, l.column)
			l.parenStack = append(l.parenStack, 'P')
		}
		l.suppressLparenFusion = false
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

		switch l.peekChar() {
		case '-':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: l.line, Column: l.column}
		case '=':
			// Zsh arithmetic compound-assign `-=`. Fuse into
			// PLUSEQ like the other compound forms.
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Column: l.column}
		default:
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
		// Zsh compound-assign `*=`. Inside arithmetic
		// `(( x *= 2 ))`, `*=` updates-with-multiply; outside,
		// bare `*` keeps its ASTERISK role for globs / expansions.
		// Fuse the two-char form into PLUSEQ so the parser's
		// existing compound-assign infix covers modulo/multiply/
		// divide variants uniformly.
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.ASTERISK, l.ch, l.line, l.column)
		}
	case '<':
		tok = l.readAngleBracket(true)
	case '>':
		tok = l.readAngleBracket(false)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.line, l.column)
	case '[':
		// Fuse `[[` into LDBRACKET unless it opens a POSIX
		// character class like `[[:alnum:]]`. The keyword `[[` is
		// always followed by whitespace (and never by `:`); a
		// `[[:` run belongs to a glob bracket expression where the
		// outer `[` opens the class and `[:name:]` is the POSIX
		// indicator. Emit two independent LBRACKETs in that case
		// so parseCommandWord can pack them into the pattern word.
		if l.peekChar() == '[' && l.peekAt(2) != ':' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LDBRACKET, Literal: literal, Line: l.line, Column: l.column}
			l.dbracketDepth++
			l.bracketStack = append(l.bracketStack, 'D')
		} else {
			tok = newToken(token.LBRACKET, l.ch, l.line, l.column)
			l.bracketStack = append(l.bracketStack, 'B')
		}
	case ']':
		// `]]` fuses to RDBRACKET only when the innermost open
		// bracket is a `[[`. Plain `[` / glob bracket classes
		// close one at a time, so `[[:alnum:]]` inside a `[[ ]]`
		// conditional keeps the outer close intact instead of
		// collapsing the class's `]]` into the conditional's.
		top := byte(0)
		if n := len(l.bracketStack); n > 0 {
			top = l.bracketStack[n-1]
		}
		if l.peekChar() == ']' && top == 'D' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.RDBRACKET, Literal: literal, Line: l.line, Column: l.column}
			if l.dbracketDepth > 0 {
				l.dbracketDepth--
			}
			l.bracketStack = l.bracketStack[:len(l.bracketStack)-1]
		} else {
			tok = newToken(token.RBRACKET, l.ch, l.line, l.column)
			if len(l.bracketStack) > 0 {
				l.bracketStack = l.bracketStack[:len(l.bracketStack)-1]
			}
		}
	case '$':
		if dollarTok, ok := l.readDollarToken(hasSpace); ok {
			return dollarTok
		}
		tok = newToken(token.DOLLAR, l.ch, l.line, l.column)

	case '&':
		switch l.peekChar() {
		case '&':
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Column: l.column}
		case '|', '!':
			// Zsh disown-in-background shortcuts: `&|` and `&!` both
			// background the command AND disown it in one step. Fuse
			// with AMPERSAND so the parser treats them like a plain
			// `&` terminator; the trailing `|` / `!` is semantic
			// metadata that downstream katas can read from the
			// Literal if needed.
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AMPERSAND, Literal: literal, Line: l.line, Column: l.column}
		default:
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
		// Zsh arithmetic compound-assign `%=`: `(( x %= 60 ))`.
		// Fuse into PLUSEQ like the other compound-assign forms
		// so the parser reuses its existing EQUALS-level infix.
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUSEQ, Literal: literal, Line: l.line, Column: l.column}
		} else {
			tok = newToken(token.PERCENT, l.ch, l.line, l.column)
		}
	case '.':
		tok = newToken(token.DOT, l.ch, l.line, l.column)
	case '"':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readString('"')
		tok.EndLine = l.line
	case '\'':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readString('\'')
		tok.EndLine = l.line
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
				next == '^' || next == ' ' || next == '\t' || next == '#' ||
				next == '"' || next == '\'' || next == '=' || next == '%' ||
				next == ',' || next == ':' || next == '@' || next == '+' ||
				next == '-' {
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

// readAngleBracket emits the token for a leading `<` or `>` and the
// punctuation pair that follows it. `isLeft` selects the LT-family
// mappings (`<<`, `<&`, `<(`, `<=`); otherwise the GT-family
// (`>>`, `>&`, `>=`, `>(`, `>|`, `>!`). Fusing here keeps the main
// NextToken switch short enough for golangci's funlen limit without
// duplicating the depth / parenStack bookkeeping.
func (l *Lexer) readAngleBracket(isLeft bool) token.Token {
	lead := l.ch
	peek := l.peekChar()
	two := func(t token.Type) token.Token {
		l.readChar()
		return token.Token{Type: t, Literal: string(lead) + string(l.ch), Line: l.line, Column: l.column}
	}
	if isLeft {
		switch peek {
		case '<':
			tok := two(token.LTLT)
			l.consumeHeredocBody()
			return tok
		case '&':
			return two(token.LTAMP)
		case '(':
			t := two(token.LT_LPAREN)
			l.parenStack = append(l.parenStack, 'P')
			return t
		case '=':
			return two(token.LE_NUM)
		}
		return newToken(token.LT, l.ch, l.line, l.column)
	}
	switch peek {
	case '>':
		return two(token.GTGT)
	case '&':
		return two(token.GTAMP)
	case '=':
		return two(token.GE_NUM)
	case '(':
		t := two(token.GT_LPAREN)
		l.parenStack = append(l.parenStack, 'P')
		return t
	case '|', '!':
		// Zsh force-clobber redirections `>|file` and `>!file`
		// override NO_CLOBBER. The trailing `|` / `!` belongs to
		// the redirect, not to a following pipeline or negation;
		// emit GT so parseCommandPipeline's redirection path
		// handles them the same as plain `>`.
		return two(token.GT)
	}
	return newToken(token.GT, l.ch, l.line, l.column)
}

// consumeHeredocBody is called immediately after emitting LTLT (the
// `<<` token). Zsh heredocs — `cat <<EOF … EOF`, `cat <<-EOF … EOF`,
// `cat <<"EOF" … EOF`, `cat <<\EOF … EOF` — have a body that begins
// on the next line and ends at a line consisting of just the
// delimiter (optionally preceded by tabs when `<<-` is used). The
// body must be opaque to the lexer, because it routinely contains
// pipes, backticks, and brace groups that would otherwise lex as
// real tokens. Peek forward for the delimiter word, then fast-
// forward past the end of the matching closer line.
//
// The parser currently has no heredoc AST node. Dropping the body
// lets real scripts parse cleanly; detection katas that care about
// heredoc content can walk source directly.
func (l *Lexer) consumeHeredocBody() {
	// Skip trailing `-` (strip-tabs flavour) and intervening
	// whitespace to land on the delimiter's first byte.
	pos := l.readPosition
	if pos < len(l.input) && l.input[pos] == '-' {
		pos++
	}
	for pos < len(l.input) && (l.input[pos] == ' ' || l.input[pos] == '\t') {
		pos++
	}
	if pos >= len(l.input) {
		return
	}
	// Pull the delimiter word. Accept quoted and backslash-escaped
	// forms by stripping those characters but using the literal
	// body as the match target. When no plausible delimiter is
	// present, leave the lexer alone.
	delim, delimEnd := extractHeredocDelim(l.input, pos)
	if delim == "" {
		return
	}
	// Scan forward to the first byte after a line consisting of
	// just the delimiter (optionally indented by tabs for `<<-`).
	i := delimEnd
	// Walk to the first newline that starts the body.
	for i < len(l.input) && l.input[i] != '\n' {
		i++
	}
	// Loop over lines until we find one equal to the delimiter.
	for i < len(l.input) {
		// Step past the newline to the next line's first byte.
		i++
		lineStart := i
		// Allow leading tabs for `<<-` form.
		for i < len(l.input) && l.input[i] == '\t' {
			i++
		}
		lineBodyStart := i
		// Advance to end of line.
		for i < len(l.input) && l.input[i] != '\n' {
			i++
		}
		// Compare the (possibly tab-stripped) line body against
		// the delimiter.
		if l.input[lineBodyStart:i] == delim {
			// Consume through to just after this closer line's
			// newline (or EOF). Update position state.
			_ = lineStart
			l.fastForwardTo(i)
			return
		}
	}
	// Delimiter never found; leave lexer state alone so the caller
	// continues as if no heredoc was detected.
}

// extractHeredocDelim scans a heredoc delimiter word starting at
// pos, handling quoted and backslash-escaped forms. Returns the
// effective delimiter text and the input index immediately past it,
// or ("", pos) when no delimiter is present.
func extractHeredocDelim(input string, pos int) (string, int) {
	if pos >= len(input) {
		return "", pos
	}
	switch input[pos] {
	case '"', '\'':
		quote := input[pos]
		pos++
		start := pos
		for pos < len(input) && input[pos] != quote {
			pos++
		}
		delim := input[start:pos]
		if pos < len(input) {
			pos++
		}
		return delim, pos
	case '\\':
		pos++
	}
	start := pos
	for pos < len(input) && (isWordByte(input[pos])) {
		pos++
	}
	if pos == start {
		return "", pos
	}
	return input[start:pos], pos
}

func isWordByte(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') ||
		('0' <= ch && ch <= '9') || ch == '_' || ch == '-'
}

// fastForwardTo re-anchors the lexer at the given input index,
// advancing line/column counters along the way so subsequent tokens
// carry accurate source coordinates.
func (l *Lexer) fastForwardTo(target int) {
	for l.readPosition <= target && l.ch != 0 {
		l.readChar()
	}
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
		// closed. Also tell the `(` handler to emit its next token
		// as a plain LPAREN regardless of peek, so `$((` reads as
		// `$(` + `(` (arithmetic open) rather than `$(` +
		// DoubleLparen which would desync nested groupings like
		// `$(((expr) * 2))`.
		l.parenStack = append(l.parenStack, 'P')
		l.suppressLparenFusion = true
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

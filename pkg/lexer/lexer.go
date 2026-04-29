// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
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

	// lastEmittedType is the type of the last token NextToken returned.
	// Used by readOpenParen to suppress `((` fusion when the previous
	// token's value position rules out an arithmetic command head
	// (e.g. `_os=$_sys*~*((alt|alt))` is a glob, not arith).
	lastEmittedType     token.Type
	lastEmittedHadSpace bool

	// dollarBraceDepth tracks how many `${` openers are still awaiting
	// their `}`. Inside a `${ … }` expansion `#` is the length /
	// pattern operator, never a comment opener.
	dollarBraceDepth int

	// precedingNewline is set by skipWhitespace when at least one
	// physical newline was consumed before the current token. Used
	// by atArithCommandPos so `((` after `cmd<NL>` fuses into
	// DoubleLparen as a fresh statement head.
	precedingNewline bool
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
	prevSuppress := l.suppressLparenFusion
	defer func() {
		if l.pendingContinuation {
			tok.HasPrecedingContinuation = true
			l.pendingContinuation = false
		}
		if prevSuppress && tok.Type != token.DOLLAR_LPAREN {
			l.suppressLparenFusion = false
		}
		switch tok.Type {
		case token.DollarLbrace:
			l.dollarBraceDepth++
		case token.RBRACE:
			if l.dollarBraceDepth > 0 {
				l.dollarBraceDepth--
				tok.ClosesDollarBrace = true
			}
		}
		l.lastEmittedType = tok.Type
		l.lastEmittedHadSpace = tok.HasPrecedingSpace
	}()
	hasSpace := l.skipWhitespace()
	if shebang, ok := l.tryShebangOrComment(hasSpace); ok {
		return shebang
	}
	if early, ok := l.dispatchEarlyReturn(hasSpace); ok {
		return early
	}
	tok = l.dispatchSimpleByte(hasSpace)
	l.readChar()
	tok.HasPrecedingSpace = hasSpace
	return tok
}

// tryShebangOrComment returns the SHEBANG token when the cursor sits
// on `#!` and reports the result; for a regular comment it consumes
// the comment and recurses into NextToken. ok=false leaves state
// unchanged for callers to continue the dispatch.
func (l *Lexer) tryShebangOrComment(hasSpace bool) (token.Token, bool) {
	if l.ch != '#' {
		return token.Token{}, false
	}
	if l.peekChar() == '!' {
		start := l.position
		for l.ch != 10 && l.ch != 0 {
			l.readChar()
		}
		return token.Token{Type: token.SHEBANG, Literal: l.input[start:l.position], Line: l.line, Column: l.column}, true
	}
	// Inside `((…))` / `$((…))` arithmetic, `#` is a special parameter
	// (count of positional args), not a comment opener. Detect by an
	// open `D` ('((') marker on the paren stack. zimfw uses
	// `(( ! # ))` and `(( # > 0 ))` heavily.
	if l.inArithmetic() {
		return token.Token{}, false
	}
	// Inside `${…}` parameter expansion, `#` is the length / pattern
	// operator, never a comment opener. Without this guard
	// `${${X## ##}%%y##}` lost its trailing `##}` to skipComment.
	if l.dollarBraceDepth > 0 {
		return token.Token{}, false
	}
	if hasSpace || l.column == 1 {
		l.skipComment()
		return l.NextToken(), true
	}
	return token.Token{}, false
}

// atArithCommandPos reports whether the upcoming `((` should fuse
// into DoubleLparen (treated as arithmetic command head) rather than
// being two separate `(` opens. The previous emitted token determines
// statement-head context: only after a separator-like token can `((`
// open arithmetic. Mid-word `*((` (glob alternation) and assignment
// RHS like `=((` are left as plain LPAREN+LPAREN so the parser can
// route them to glob / grouped-expression paths.
func (l *Lexer) atArithCommandPos() bool {
	switch l.lastEmittedType {
	// Empty (start-of-input) and statement-separator contexts open
	// a new statement: `((` here is arithmetic.
	case "", token.SEMICOLON, token.DSEMI, token.AMPERSAND,
		token.AND, token.OR, token.PIPE, token.BANG,
		token.LBRACE, token.LPAREN, token.DoubleLparen,
		token.LDBRACKET, token.If, token.THEN, token.ELSE, token.ELIF,
		token.DO, token.WHILE, token.FOR, token.CASE,
		token.LET, token.RETURN:
		return true
	}
	// A newline-separated `((` always opens a fresh statement, even
	// after an IDENT / value token (e.g. `cmd<NL>(( expr ))`). The
	// physical newline is the statement separator that the case list
	// above covers for `;` / `&&` / `||`.
	if l.precedingNewline {
		return true
	}
	// A space-separated `((` after a value-like token is also a
	// statement head (e.g. `cmd ; (( … ))` after a newline/space —
	// the preceding token is whatever ended the previous statement).
	return l.lastEmittedHadSpace && l.lastTokenIsValueTerminator()
}

func (l *Lexer) lastTokenIsValueTerminator() bool {
	// Tokens that conventionally terminate a value/expression and
	// allow a fresh statement to follow without a separator.
	switch l.lastEmittedType {
	case token.RBRACE, token.RPAREN, token.RBRACKET, token.RDBRACKET,
		token.DoubleRparen, token.Fi, token.DONE, token.ESAC:
		return true
	}
	return false
}

// inArithmetic reports whether the lexer is currently inside a `((`
// arithmetic group. The parser's `((`-marker is 'D' on parenStack.
func (l *Lexer) inArithmetic() bool {
	for i := len(l.parenStack) - 1; i >= 0; i-- {
		if l.parenStack[i] == 'D' {
			return true
		}
	}
	return false
}

// dispatchEarlyReturn covers the byte categories whose handler fully
// owns the return — they advance the cursor themselves and the
// outer NextToken must skip the trailing readChar. The simple-byte
// switch in dispatchSimpleByte runs first; this path only fires when
// the character is part of an identifier / number / dollar-form / or
// the dash-keyword shortcut. The narrower letter check
// (isAsciiLetterOrUnderscore) keeps `~`, `:`, `.`, `/`, `@` —
// which `isLetter` accepts mid-word — out of the identifier path so
// those bytes hit their own switch arms first.
func (l *Lexer) dispatchEarlyReturn(hasSpace bool) (token.Token, bool) {
	if l.ch == '-' && isLetter(l.peekChar()) {
		if tok, ok := l.tryDashKeyword(hasSpace); ok {
			return tok, true
		}
	}
	switch {
	case l.ch == '$':
		if tok, ok := l.readDollarToken(hasSpace); ok {
			return tok, true
		}
		return token.Token{}, false
	case l.isIdentLeadByte(l.ch):
		return l.readIdentifierToken(hasSpace), true
	case isDigit(l.ch):
		return l.readNumberToken(hasSpace), true
	case l.ch >= 0x80:
		return l.readUnicodeIdent(hasSpace), true
	}
	return token.Token{}, false
}

// isAsciiLetterOrUnderscore is the narrower letter test used by
// dispatchEarlyReturn to decide whether to consume an identifier
// run. It deliberately excludes the punctuation aliases that the
// wider isLetter accepts mid-word but that have their own explicit
// switch arm in dispatchSimpleByte / dispatchBracketsAndOps
// (`~`, `:`, `.`). It includes `/` and `@` because the original
// switch had no explicit arm for them — they only ever entered the
// identifier path via the wide isLetter check.
func (l *Lexer) isIdentLeadByte(ch byte) bool {
	if ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_' {
		return true
	}
	if ch == '@' {
		return true
	}
	// `/` triggers an identifier read only outside arithmetic. Inside
	// `((…))` / `$((…))` `/` is the division operator and must lex as
	// SLASH; the parenStack 'D' marker tracks the arithmetic frame.
	return ch == '/' && !l.inArithmetic()
}

// dispatchSimpleByte handles the per-byte switch where the handler
// only sets tok and the outer NextToken advances the cursor.
func (l *Lexer) dispatchSimpleByte(hasSpace bool) token.Token {
	switch l.ch {
	case '#':
		return newToken(token.HASH, l.ch, l.line, l.column)
	case '=':
		return l.readEqualsLead(hasSpace)
	case ';':
		return l.readSemicolonLead()
	case ':':
		return newToken(token.COLON, l.ch, l.line, l.column)
	case '?':
		return newToken(token.QUESTION, l.ch, l.line, l.column)
	case '(':
		return l.readOpenParen()
	case ')':
		return l.readCloseParen()
	case ',':
		return newToken(token.COMMA, l.ch, l.line, l.column)
	case '+':
		return l.readPlusLead()
	case '-':
		return l.readMinusLead(hasSpace)
	case '!':
		return l.readBangLead()
	case '*':
		return l.readArithCompoundOr(token.ASTERISK)
	case '<':
		return l.readAngleBracket(true)
	case '>':
		return l.readAngleBracket(false)
	}
	return l.dispatchBracketsAndOps(hasSpace)
}

// simpleBracketTokens maps single-byte tokens whose lex result has no
// state side-effect — they just become the matching token type.
var simpleBracketTokens = map[byte]token.Type{
	'{': token.LBRACE, '}': token.RBRACE,
	'`': token.BACKTICK, '~': token.TILDE,
	'^': token.CARET, '/': token.SLASH,
	'.': token.DOT, '$': token.DOLLAR,
}

func (l *Lexer) dispatchBracketsAndOps(hasSpace bool) token.Token {
	if t, ok := simpleBracketTokens[l.ch]; ok {
		return newToken(t, l.ch, l.line, l.column)
	}
	switch l.ch {
	case '[':
		return l.readOpenBracket()
	case ']':
		return l.readCloseBracket()
	case '&':
		return l.readAmpersandLead()
	case '|':
		return l.readPipeLead()
	case '%':
		return l.readArithCompoundOr(token.PERCENT)
	case '"', '\'':
		return l.readQuotedString(l.ch)
	case 0:
		return token.Token{Type: token.EOF, Line: l.line, Column: l.column}
	}
	return l.dispatchTerminalDefault(hasSpace)
}

func (l *Lexer) dispatchTerminalDefault(_ bool) token.Token {
	if l.ch == '\\' {
		return l.readBackslashEscape()
	}
	return newToken(token.ILLEGAL, l.ch, l.line, l.column)
}

func (l *Lexer) readEqualsLead(hasSpace bool) token.Token {
	switch l.peekChar() {
	case '=':
		return l.readFusedToken(token.EQ)
	case '~':
		return l.readFusedToken(token.EQTILDE)
	case '(':
		if !hasSpace {
			return newToken(token.ASSIGN, l.ch, l.line, l.column)
		}
		tok := l.readFusedToken(token.EQ_LPAREN)
		l.parenStack = append(l.parenStack, 'P')
		return tok
	}
	return newToken(token.ASSIGN, l.ch, l.line, l.column)
}

func (l *Lexer) readSemicolonLead() token.Token {
	if l.peekChar() == ';' {
		return l.readFusedToken(token.DSEMI)
	}
	// Zsh case-clause fall-through markers `;|` and `;&` terminate
	// the clause body just like `;;`. Fuse them into DSEMI so the
	// parser's case loop honours the boundary; the literal preserves
	// the source spelling for katas that walk it.
	if l.peekChar() == '|' || l.peekChar() == '&' {
		return l.readFusedToken(token.DSEMI)
	}
	return newToken(token.SEMICOLON, l.ch, l.line, l.column)
}

func (l *Lexer) readOpenParen() token.Token {
	defer func() { l.suppressLparenFusion = false }()
	if l.peekChar() == '(' && !l.suppressLparenFusion && !l.inArithmetic() && l.atArithCommandPos() {
		// `(((` at command position is ambiguous: `( ((` (subshell
		// open + arith) when a space separates the inner `((` from
		// the operand (`if ((( cond ))`), or `(( (` (arith with
		// grouping) when no space appears (`(((x * y) + z))`).
		// Disambiguate by peeking past the third `(`.
		if l.peekAt(2) == '(' && l.peekAt(3) == ' ' {
			tok := newToken(token.LPAREN, l.ch, l.line, l.column)
			l.parenStack = append(l.parenStack, 'P')
			return tok
		}
		tok := l.readFusedToken(token.DoubleLparen)
		l.parenStack = append(l.parenStack, 'D')
		return tok
	}
	tok := newToken(token.LPAREN, l.ch, l.line, l.column)
	l.parenStack = append(l.parenStack, 'P')
	return tok
}

func (l *Lexer) readPlusLead() token.Token {
	switch l.peekChar() {
	case '+':
		return l.readFusedToken(token.INC)
	case '=':
		return l.readFusedToken(token.PLUSEQ)
	}
	return newToken(token.PLUS, l.ch, l.line, l.column)
}

func (l *Lexer) readMinusLead(_ bool) token.Token {
	switch l.peekChar() {
	case '-':
		return l.readFusedToken(token.DEC)
	case '=':
		return l.readFusedToken(token.PLUSEQ)
	}
	return newToken(token.MINUS, l.ch, l.line, l.column)
}

func (l *Lexer) tryDashKeyword(hasSpace bool) (token.Token, bool) {
	saved := *l
	literal := l.readIdentifier()
	if t := token.LookupIdent(literal); t != token.IDENT {
		return token.Token{
			Type: t, Literal: literal,
			Line: saved.line, Column: saved.column,
			HasPrecedingSpace: hasSpace,
		}, true
	}
	*l = saved
	return token.Token{}, false
}

func (l *Lexer) readBangLead() token.Token {
	if l.peekChar() == '=' {
		return l.readFusedToken(token.NotEq)
	}
	return newToken(token.BANG, l.ch, l.line, l.column)
}

func (l *Lexer) readArithCompoundOr(plain token.Type) token.Token {
	if l.peekChar() == '=' {
		return l.readFusedToken(token.PLUSEQ)
	}
	return newToken(plain, l.ch, l.line, l.column)
}

func (l *Lexer) readOpenBracket() token.Token {
	if l.peekChar() == '[' && l.peekAt(2) != ':' {
		tok := l.readFusedToken(token.LDBRACKET)
		l.dbracketDepth++
		l.bracketStack = append(l.bracketStack, 'D')
		return tok
	}
	tok := newToken(token.LBRACKET, l.ch, l.line, l.column)
	l.bracketStack = append(l.bracketStack, 'B')
	return tok
}

func (l *Lexer) readCloseBracket() token.Token {
	top := byte(0)
	if n := len(l.bracketStack); n > 0 {
		top = l.bracketStack[n-1]
	}
	if l.peekChar() == ']' && top == 'D' {
		tok := l.readFusedToken(token.RDBRACKET)
		if l.dbracketDepth > 0 {
			l.dbracketDepth--
		}
		l.bracketStack = l.bracketStack[:len(l.bracketStack)-1]
		return tok
	}
	tok := newToken(token.RBRACKET, l.ch, l.line, l.column)
	if len(l.bracketStack) > 0 {
		l.bracketStack = l.bracketStack[:len(l.bracketStack)-1]
	}
	return tok
}

func (l *Lexer) readAmpersandLead() token.Token {
	switch l.peekChar() {
	case '&':
		return l.readFusedToken(token.AND)
	case '|', '!':
		return l.readFusedToken(token.AMPERSAND)
	}
	return newToken(token.AMPERSAND, l.ch, l.line, l.column)
}

func (l *Lexer) readPipeLead() token.Token {
	if l.peekChar() == '|' {
		return l.readFusedToken(token.OR)
	}
	// Zsh `|&` is the stderr-pipe shorthand (`2>&1 |`). Fuse to a
	// PIPE token so the parser routes through the pipeline path; the
	// literal preserves the source spelling for katas that need it.
	if l.peekChar() == '&' {
		tok := l.readFusedToken(token.PIPE)
		return tok
	}
	return newToken(token.PIPE, l.ch, l.line, l.column)
}

func (l *Lexer) readQuotedString(quote byte) token.Token {
	line, col := l.line, l.column
	literal := l.readString(quote)
	return token.Token{
		Type: token.STRING, Literal: literal,
		Line: line, Column: col, EndLine: l.line,
	}
}

// readFusedToken consumes the two-byte run starting at the current
// cursor and returns a fused token of the requested type.
func (l *Lexer) readFusedToken(t token.Type) token.Token {
	ch := l.ch
	l.readChar()
	return token.Token{
		Type:    t,
		Literal: string(ch) + string(l.ch),
		Line:    l.line,
		Column:  l.column,
	}
}

func (l *Lexer) readIdentifierToken(hasSpace bool) token.Token {
	line, col := l.line, l.column
	tok := token.Token{Line: line, Column: col, HasPrecedingSpace: hasSpace}
	tok.Literal = l.readIdentifier()
	tok.Type = token.LookupIdent(tok.Literal)
	if tok.Type != token.IDENT && l.ch == '=' {
		tok.Type = token.IDENT
	}
	return tok
}

func (l *Lexer) readNumberToken(hasSpace bool) token.Token {
	line, col := l.line, l.column
	return token.Token{
		Type: token.INT, Literal: l.readNumber(),
		Line: line, Column: col, HasPrecedingSpace: hasSpace,
	}
}

func (l *Lexer) readUnicodeIdent(hasSpace bool) token.Token {
	line, col := l.line, l.column
	start := l.position
	for l.ch >= 0x80 || isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return token.Token{
		Type: token.IDENT, Literal: l.input[start:l.position],
		Line: line, Column: col, HasPrecedingSpace: hasSpace,
	}
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
			// Inside `[[ … ]]`, `<(` is not process substitution
			// — `<->` numeric-range glob followed by `(...)` glob
			// alternation. Emit a plain LT so the parser sees the
			// `(` as the alternation opener.
			if l.dbracketDepth > 0 {
				return newToken(token.LT, l.ch, l.line, l.column)
			}
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
		// Inside `[[ … ]]`, `>(` is not process substitution —
		// `<->(...)` numeric-range glob followed by alternation.
		// Emit a plain GT so the parser sees `(` as the opener.
		if l.dbracketDepth > 0 {
			return newToken(token.GT, l.ch, l.line, l.column)
		}
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
	pos := skipHeredocPrefix(l.input, l.readPosition)
	if pos >= len(l.input) {
		return
	}
	delim, delimEnd := extractHeredocDelim(l.input, pos)
	if delim == "" {
		return
	}
	if closer := findHeredocCloser(l.input, delim, delimEnd); closer >= 0 {
		l.fastForwardTo(closer)
	}
}

// skipHeredocPrefix advances past the optional `-` strip-tabs marker
// and any whitespace between `<<` and the delimiter word.
func skipHeredocPrefix(input string, pos int) int {
	if pos < len(input) && input[pos] == '-' {
		pos++
	}
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t') {
		pos++
	}
	return pos
}

// findHeredocCloser locates the line whose body matches delim
// (optionally indented by tabs for `<<-`) and returns the byte index
// of the newline that ends that closer line, or -1 when no closer
// is found.
func findHeredocCloser(input, delim string, delimEnd int) int {
	i := delimEnd
	for i < len(input) && input[i] != '\n' {
		i++
	}
	for i < len(input) {
		i++
		for i < len(input) && input[i] == '\t' {
			i++
		}
		bodyStart := i
		for i < len(input) && input[i] != '\n' {
			i++
		}
		if input[bodyStart:i] == delim {
			return i
		}
	}
	return -1
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

// readBackslashEscape handles a `\X` sequence outside string contexts.
// Backslash quotes exactly one following character; for any of the
// common glob / shell / regex metacharacters or alphanumerics we
// emit the pair as an IDENT word so parseCommandWord folds it into
// the surrounding word. Anything else falls back to ILLEGAL to keep
// token-aware contexts stable.
// backslashEscapable lists the bytes that may follow `\` outside a
// string context. Anything else falls through to the ILLEGAL token.
var backslashEscapable = map[byte]struct{}{
	'(': {}, ')': {}, '*': {}, '?': {}, '[': {}, ']': {},
	'|': {}, '&': {}, ';': {}, '<': {}, '>': {}, '{': {}, '}': {},
	'$': {}, '\\': {}, '/': {}, '.': {}, '!': {}, '~': {},
	'^': {}, ' ': {}, '\t': {}, '#': {}, '"': {}, '\'': {},
	'=': {}, '%': {}, ',': {}, ':': {}, '@': {}, '+': {}, '-': {},
	'`': {},
}

func (l *Lexer) readBackslashEscape() token.Token {
	next := l.peekChar()
	if !isBackslashEscapable(next) {
		return newToken(token.ILLEGAL, l.ch, l.line, l.column)
	}
	line, col := l.line, l.column
	l.readChar()
	return token.Token{
		Type:    token.IDENT,
		Literal: "\\" + string(l.ch),
		Line:    line,
		Column:  col,
	}
}

func isBackslashEscapable(ch byte) bool {
	if isAlphaNumByte(ch) {
		return true
	}
	_, hit := backslashEscapable[ch]
	return hit
}

func isAlphaNumByte(ch byte) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		('0' <= ch && ch <= '9')
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
		line := l.line
		l.readChar() // consume '$'
		tok.Type = token.VARIABLE
		tok.Literal = "$" + l.readIdentifier()
		// Stamp the START line, not the line after readIdentifier
		// — readIdentifier advances PAST the last identifier byte,
		// which on `$name\n` lands l.ch on the newline and bumps
		// l.line. Without capturing `line` first, every variable
		// followed by a newline got tagged as the next line, which
		// confused the parser's same-line argument check (e.g.
		// `for x in $files\ndo` exited the items loop too early
		// because $files looked like it was on the do-line).
		tok.Line = line
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
	inArith := l.inArithmetic()
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '-' {
		// Inside arithmetic `/` is the division operator, never part
		// of an identifier. Outside arithmetic `/` extends an
		// identifier so `/usr/bin/foo` survives as a single word.
		if inArith && l.ch == '/' {
			break
		}
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	if l.tryReadBaseLiteral() {
		return l.input[position:l.position]
	}
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '#' && isDigit(l.peekChar()) {
		// `BASE#NUM` — Zsh custom-base literal, e.g. `16#ff`.
		l.readChar() // #
		for isDigit(l.ch) || isHexDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

// tryReadBaseLiteral consumes a Zsh `0x…` / `0b…` / `0o…` integer
// literal when the prefix is followed by at least one digit. Returns
// true iff bytes were consumed. `0x${var}` (Zsh string concat with a
// parameter expansion) must NOT consume the prefix — the parser
// recovers via INT(0) + IDENT(x) + DollarLbrace concatenation.
func (l *Lexer) tryReadBaseLiteral() bool {
	if l.ch != '0' {
		return false
	}
	if !isBasePrefix(l.peekChar()) {
		return false
	}
	third := l.peekAt(2)
	if !isDigit(third) && !isHexDigit(third) {
		return false
	}
	l.readChar() // 0
	l.readChar() // base prefix letter
	for isDigit(l.ch) || isHexDigit(l.ch) {
		l.readChar()
	}
	return true
}

func isBasePrefix(ch byte) bool {
	switch ch {
	case 'x', 'X', 'b', 'B', 'o', 'O':
		return true
	}
	return false
}

func isHexDigit(ch byte) bool {
	return (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
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
	position := l.position
	braceDepth := 0
	for {
		l.readChar()
		if l.ch == 0 {
			break
		}
		if l.ch == quote && braceDepth == 0 {
			break
		}
		if l.absorbStringEscape(honourEscapes) {
			if l.ch == 0 {
				break
			}
			continue
		}
		if l.absorbDollarBraceOpen(honourEscapes, &braceDepth) {
			continue
		}
		if honourEscapes && l.absorbEmbeddedDollarParen() {
			continue
		}
		l.trackBraceDepth(&braceDepth)
	}
	return l.sliceClosedString(position)
}

// absorbEmbeddedDollarParen handles `$(…)` command substitution
// inside a double-quoted string body. Walks past the matching `)`
// while honouring nested `'…'` single-quoted runs so a `"` inside a
// single-quoted regex (`'[^"]+' `) doesn't close the outer
// double-quoted string.
func (l *Lexer) absorbEmbeddedDollarParen() bool {
	if l.ch != '$' || l.peekChar() != '(' {
		return false
	}
	l.readChar() // onto `(`
	depth := 1
	for depth > 0 {
		l.readChar()
		switch l.ch {
		case 0:
			return true
		case '\\':
			l.readChar()
		case '\'':
			for l.ch != 0 {
				l.readChar()
				if l.ch == '\'' {
					break
				}
			}
		case '(':
			depth++
		case ')':
			depth--
		}
	}
	return true
}

func (l *Lexer) absorbStringEscape(honourEscapes bool) bool {
	if !honourEscapes || l.ch != '\\' {
		return false
	}
	l.readChar()
	return true
}

func (l *Lexer) absorbDollarBraceOpen(honourEscapes bool, braceDepth *int) bool {
	if !honourEscapes || l.ch != '$' || l.peekChar() != '{' {
		return false
	}
	l.readChar()
	*braceDepth++
	return true
}

func (l *Lexer) trackBraceDepth(braceDepth *int) {
	if *braceDepth == 0 {
		return
	}
	switch l.ch {
	case '{':
		*braceDepth++
	case '}':
		*braceDepth--
	}
}

func (l *Lexer) sliceClosedString(position int) string {
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
	return l.input[position:end]
}

func (l *Lexer) skipWhitespace() bool {
	skipped := false
	l.precedingNewline = false
	for {
		switch l.ch {
		case ' ', '\t', '\r':
			skipped = true
			l.readChar()
			continue
		case '\n':
			skipped = true
			l.precedingNewline = true
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

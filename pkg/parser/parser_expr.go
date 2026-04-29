// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

// expressionTerminators lists the tokens that end the current
// expression silently — usually because they belong to the enclosing
// statement frame (`]]`, `&&`, keyword openers, …).
var expressionTerminators = map[token.Type]struct{}{
	token.RDBRACKET: {}, token.AND: {}, token.OR: {},
	token.THEN: {}, token.ELSE: {}, token.ELIF: {}, token.Fi: {},
	token.DO: {}, token.DONE: {}, token.ESAC: {},
	token.SEMICOLON: {}, token.DSEMI: {},
	token.FOR: {}, token.WHILE: {}, token.If: {}, token.CASE: {},
	token.SELECT: {}, token.LET: {}, token.RETURN: {},
	token.TYPESET: {}, token.DECLARE: {},
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if _, hit := expressionTerminators[p.curToken.Type]; hit {
		return nil
	}
	// Inside `[[ … ]]` a `((` is glob alternation grouping (not
	// arithmetic). Decompose the fused DoubleLparen into a LPAREN so
	// parseGroupedExpression's glob-alt path handles it.
	if p.inDoubleBracket && p.curTokenIs(token.DoubleLparen) {
		p.curToken.Type = token.LPAREN
		p.curToken.Literal = "("
	}
	// Inside `${…[KEY]}` subscripts and `[[ … ]]` tests, Zsh keywords
	// are literal pattern words, not statement-block openers. Return
	// an Identifier so the surrounding parse keeps moving.
	if (p.inDoubleBracket || p.inArithmetic) && isDoubleBracketLiteralKeyword(p.curToken.Type) {
		tok := p.curToken
		return &ast.Identifier{Token: tok, Value: tok.Literal}
	}
	// Inside `[[ … ]]`, a leading `[` opens a glob bracket-class
	// fragment (`[abc]`, `[[:alnum:]]`, `[^[:blank:]]`), not the `[`
	// test-builtin or an array subscript. The default LBRACKET prefix
	// (parseSingleCommand) gobbles every glued token until a command
	// delimiter and walks past the closing `]]`. Consume the bracket
	// body as a literal pattern fragment instead.
	if p.inDoubleBracket && p.curTokenIs(token.LBRACKET) {
		return p.parseDoubleBracketGlobBracket()
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		if p.inDoubleBracket {
			tok := p.curToken
			return &ast.Identifier{Token: tok, Value: tok.Literal}
		}
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		if p.expressionInfixShouldBreak() {
			break
		}
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// isDoubleBracketLiteralKeyword reports whether a Zsh keyword token
// should be treated as a literal pattern word inside `[[ … ]]` or a
// `${var[KEY]}` subscript. Most reserved words (FUNCTION, IF, FOR, …)
// only mean "statement head" in command position; in pattern context
// they are simply pattern strings.
func isDoubleBracketLiteralKeyword(t token.Type) bool {
	return t == token.FUNCTION
}

// expressionInfixShouldBreak reports whether the infix chain in
// parseExpression should stop before the next infix call. It captures
// the various Zsh syntactic guards that prevent the infix loop from
// crossing a statement / bracket / glob boundary.
func (p *Parser) expressionInfixShouldBreak() bool {
	if p.peekTokenIs(token.RDBRACKET) || p.curTokenIs(token.RDBRACKET) {
		return true
	}
	if p.inDoubleBracket && p.peekTokenIs(token.LPAREN) {
		return true
	}
	if p.peekTokenIs(token.LPAREN) && p.peekToken.HasPrecedingSpace {
		return true
	}
	if !p.inArithmetic && p.peekShouldBreakInfix() {
		return true
	}
	return false
}

// peekShouldBreakInfix groups the not-in-arithmetic infix breaks so
// expressionInfixShouldBreak stays under the gocyclo threshold. The
// SLASH/LBRACKET arms guard glob-context shapes; AMPERSAND/CARET/COMMA
// arms guard shell-control bytes that are only infix inside `((…))`.
func (p *Parser) peekShouldBreakInfix() bool {
	if p.peekTokenIs(token.LBRACKET) && p.peekToken.HasPrecedingSpace {
		return true
	}
	if p.peekTokenIs(token.SLASH) && !p.peekToken.HasPrecedingSpace {
		return true
	}
	if p.peekTokenIs(token.AMPERSAND) || p.peekTokenIs(token.CARET) {
		return true
	}
	if p.peekTokenIs(token.COMMA) {
		return true
	}
	return false
}

func (p *Parser) parseIdentifier() ast.Expression {
	tok := p.curToken
	value := tok.Literal
	// Inside arithmetic, Zsh concatenates an IDENT with a glued
	// VARIABLE / INT to form a dynamic variable name:
	// `(( X_$y == 2 ))` reads as the value of `X_<expanded y>`.
	// Absorb the glued tokens so the closing `))` lines up.
	if p.inArithmetic {
		for !p.peekToken.HasPrecedingSpace &&
			(p.peekTokenIs(token.VARIABLE) || p.peekTokenIs(token.INT)) {
			p.nextToken()
			value += p.curToken.Literal
		}
	}
	return &ast.Identifier{Token: tok, Value: value}
}

// parseEqualsForm handles Zsh's `=cmd` notation where a leading `=`
// with no preceding space substitutes the absolute path of the named
// command (roughly equivalent to `$(which cmd)`). The lexer splits
// this into ASSIGN + IDENT; we fuse them into a single identifier so
// the downstream pipeline / redirection code treats it as a command
// head. Used at statement head position; infix ASSIGN still handles
// plain assignments.
func (p *Parser) parseEqualsForm() ast.Expression {
	eqTok := p.curToken
	if p.peekToken.HasPrecedingSpace || p.peekTokenIs(token.EOF) {
		return nil
	}
	p.nextToken()
	name := "=" + p.curToken.Literal
	cmd := &ast.SimpleCommand{
		Token: eqTok,
		Name:  &ast.Identifier{Token: eqTok, Value: name},
	}
	for !p.isCommandDelimiter(p.peekToken) && p.peekToken.Line == eqTok.Line {
		p.nextToken()
		cmd.Arguments = append(cmd.Arguments, p.parseCommandWord())
	}
	return cmd
}

// parseKeywordAsCommand wraps a Zsh keyword (currently RETURN) as a
// SimpleCommand so it can appear as the right-hand side of a logical
// expression chain like `cmd || return 0`. Any arguments on the same
// line are collected via parseCommandWord so `return 1` or
// `break 2` (when BREAK/CONTINUE get their own tokens) round-trip
// through the expression layer.
func (p *Parser) parseKeywordAsCommand() ast.Expression {
	tok := p.curToken
	cmd := &ast.SimpleCommand{
		Token: tok,
		Name:  &ast.Identifier{Token: tok, Value: tok.Literal},
	}
	for !p.isCommandDelimiter(p.peekToken) && p.peekToken.Line == tok.Line {
		p.nextToken()
		cmd.Arguments = append(cmd.Arguments, p.parseCommandWord())
	}
	return cmd
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := parseZshIntLiteral(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	// Inside arithmetic (`(( … ))`, `$(( … ))`), Zsh accepts
	// floating-point literals like `1.0` and the trailing-dot
	// variant `1000.`. The lexer emits these as INT + DOT (+ INT
	// for the fractional part). Absorb the DOT (and any following
	// digit run) so the closing `))` aligns. The AST keeps the
	// integer part as Value; katas that need the source form walk
	// Token.Literal which still names the int run.
	if p.inArithmetic && p.peekTokenIs(token.DOT) && !p.peekToken.HasPrecedingSpace {
		p.nextToken() // consume DOT
		if p.peekTokenIs(token.INT) && !p.peekToken.HasPrecedingSpace {
			p.nextToken() // consume fractional INT
		}
	}
	// Zsh number-base concat: `0x${var}`, `16#$base`, `0b${b}` —
	// arithmetic treats these as a single literal whose value comes
	// from the surrounding string concat. Absorb glued IDENT / `${…}`
	// / VARIABLE / HASH+INT tail tokens so the closing `))` lines up.
	if p.inArithmetic {
		p.absorbArithmeticNumberTail()
	}
	return lit
}

// parseHashSpecial returns `#` as a special-parameter identifier when
// curToken is HASH. Used inside `((…))` arithmetic where HASH stands
// for the count of positional arguments. Outside arithmetic HASH
// opens a comment which the lexer skips before the parser sees it,
// so this prefix only fires in arithmetic context.
func (p *Parser) parseHashSpecial() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: "#"}
}

// parseZshIntLiteral converts a Zsh integer literal to int64. Handles
// the standard 0x / 0b / 0o / decimal forms via strconv plus the
// custom-base `BASE#NUM` form (e.g. `16#ff`).
func parseZshIntLiteral(s string) (int64, error) {
	if hash := strings.IndexByte(s, '#'); hash > 0 {
		base, err := strconv.Atoi(s[:hash])
		if err != nil {
			return 0, err
		}
		return strconv.ParseInt(s[hash+1:], base, 64)
	}
	return strconv.ParseInt(s, 0, 64)
}

// absorbArithmeticNumberTail walks the no-preceding-space tail after
// an INT, swallowing the tokens that complete a Zsh numeric literal
// with concat or custom-base form. Stops at the first whitespace gap
// or non-eligible token.
func (p *Parser) absorbArithmeticNumberTail() {
	for !p.peekToken.HasPrecedingSpace {
		switch p.peekToken.Type {
		case token.IDENT, token.VARIABLE, token.INT:
			p.nextToken()
		case token.DollarLbrace:
			p.nextToken()
			p.skipDollarBraceBody()
		case token.HASH:
			p.nextToken()
		default:
			return
		}
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{Token: p.curToken, Left: left, Operator: p.curToken.Literal}
}

// parseDoubleBracketGlobBracket consumes a balanced `[ … ]` glob
// bracket-class while inside `[[ … ]]`. POSIX classes (`[:alnum:]`)
// nest a leading `[` that does NOT increment depth; their closing
// `]` does NOT decrement the outer depth. Returns an Identifier
// carrying the literal text and leaves curToken on the matching
// outer `]`.
func (p *Parser) parseDoubleBracketGlobBracket() ast.Expression {
	startTok := p.curToken
	literal := p.curToken.Literal
	depth := 1
	for depth > 0 && !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RDBRACKET) {
		p.nextToken()
		literal += p.curToken.Literal
		switch {
		case p.curTokenIs(token.LBRACKET) && p.peekTokenIs(token.COLON):
			// POSIX class `[:name:]`: drain `:`, IDENT, `]` without
			// touching outer depth.
			for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.RDBRACKET) {
				p.nextToken()
				literal += p.curToken.Literal
				if p.curTokenIs(token.RBRACKET) {
					break
				}
			}
		case p.curTokenIs(token.LBRACKET):
			depth++
		case p.curTokenIs(token.RBRACKET):
			depth--
		}
	}
	return &ast.Identifier{Token: startTok, Value: literal}
}

func (p *Parser) parseDoubleBracketExpression() ast.Expression {
	bracketToken := p.curToken
	p.nextToken()
	prevInDB := p.inDoubleBracket
	p.inDoubleBracket = true
	defer func() { p.inDoubleBracket = prevInDB }()
	expressions := []ast.Expression{}
	for !p.curTokenIs(token.RDBRACKET) && !p.curTokenIs(token.EOF) {
		exp := p.parseExpression(LOWEST)
		if exp != nil {
			expressions = append(expressions, exp)
		}
		if !p.curTokenIs(token.RDBRACKET) && !p.curTokenIs(token.EOF) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.RDBRACKET) {
		p.peekError(token.RDBRACKET)
		return nil
	}
	return &ast.DoubleBracketExpression{Token: bracketToken, Elements: expressions}
}

// parseTernaryExpression handles Zsh arithmetic ternary
// `cond ? then : else`. Builds an InfixExpression `?` whose Right
// is itself an InfixExpression `:` between the then- and else-
// branches — keeps the AST shape simple while letting katas walk
// either branch via Left / Right traversal.
func (p *Parser) parseTernaryExpression(left ast.Expression) ast.Expression {
	tok := p.curToken
	p.nextToken() // past `?`
	thenExpr := p.parseExpression(LOWEST)
	// expectPeek consumes `:` if present; otherwise return the
	// partial ternary so the parser can recover.
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		colonTok := p.curToken
		p.nextToken()
		elseExpr := p.parseExpression(LOWEST)
		thenExpr = &ast.InfixExpression{
			Token:    colonTok,
			Operator: ":",
			Left:     thenExpr,
			Right:    elseExpr,
		}
	}
	return &ast.InfixExpression{Token: tok, Operator: "?", Left: left, Right: thenExpr}
}

// isEmptyRhsTerminator reports whether the token signals an empty
// right-hand side for an assignment-like infix. Used by ASSIGN /
// PLUSEQ to bail before advancing past a keyword that belongs to the
// next statement.
func isEmptyRhsTerminator(t token.Type) bool {
	switch t {
	case token.SEMICOLON, token.EOF, token.PIPE, token.AMPERSAND,
		token.AND, token.OR, token.RPAREN, token.RBRACE,
		token.RDBRACKET, token.RBRACKET,
		token.FOR, token.WHILE, token.If, token.CASE,
		token.SELECT, token.LET, token.RETURN,
		token.TYPESET, token.DECLARE,
		token.THEN, token.ELSE, token.ELIF, token.Fi,
		token.DO, token.DONE, token.ESAC, token.DSEMI:
		return true
	}
	return false
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	// Empty RHS for assignment-shaped infixes (`X=<NL>for …`,
	// `Y+=<NL>case …`): if peek already signals a terminator /
	// next-statement keyword, leave Right=nil and DON'T advance —
	// otherwise parseStatement's outer nextToken would skip past
	// the keyword.
	isAssign := p.curTokenIs(token.ASSIGN) || p.curTokenIs(token.PLUSEQ)
	if isAssign && isEmptyRhsTerminator(p.peekToken.Type) {
		return expression
	}
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	tok := p.curToken
	p.nextToken()

	if p.inArithmetic {
		exp := p.parseExpression(LOWEST)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		return &ast.GroupedExpression{Token: tok, Expression: exp}
	}

	// When the group opens with a statement keyword (FOR, WHILE,
	// IF, CASE, LBRACE, DoubleLparen, LDBRACKET, BANG, TYPESET,
	// DECLARE, LOCAL, LET, SELECT, COPROC), the `( … )` is a
	// subshell body, not an array literal. Dispatch through
	// parseStatement so loops and conditionals inside subshells
	// like `time (for x in a; do …; done)` parse correctly.
	//
	// Exception: inside `[[ … ]]` the group is a glob alternation,
	// so keywords are literal pattern words (`[[ $x = (select|cont) ]]`).
	// Fall through to the glob-alt branch when inDoubleBracket.
	if !p.inDoubleBracket {
		switch p.curToken.Type {
		case token.FOR, token.WHILE, token.If, token.CASE,
			token.DoubleLparen, token.LDBRACKET, token.BANG, token.TYPESET,
			token.DECLARE, token.LET, token.SELECT, token.COPROC:
			statements := []ast.Statement{}
			for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
				stmt := p.parseStatement()
				if stmt != nil {
					statements = append(statements, stmt)
				}
				p.nextToken()
			}
			return &ast.GroupedExpression{Token: tok, Expression: &ast.BlockStatement{Statements: statements}}
		}
	}

	// Array Literal / glob alternation mode. Inside `[[ ]]` a
	// parenthesised group `(a|b|c)` is a glob alternation where `|`
	// is the pattern separator, not a pipe. Skip bare PIPE tokens
	// between elements so patterns like `(wip|WIP)` and the richer
	// p10k forms `(|*[^[:alnum:]])(wip|WIP)(|[^[:alnum:]]*)` parse
	// as a sequence of alternatives rather than erroring on the
	// first `|`. Outside `[[ ]]` a `|` inside `( )` would be unusual
	// — array literals never contain pipe-separated elements — so
	// swallowing the token there is safe.
	elements := []ast.Expression{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.PIPE) {
			p.nextToken()
			continue
		}
		elem := p.parseCommandWord()
		elements = append(elements, elem)
		p.nextToken()
	}

	// Advance past the array literal's `)` and signal the enclosing
	// block to treat it as already-consumed. Without this, a `( arr=(
	// "x" ); list=( "y" ) )` subshell broke at the array's `)` —
	// parseBlockStatement saw curToken=RPAREN and ended the subshell
	// body prematurely. Mirrors parseDollarParenExpression's flag.
	if p.curTokenIs(token.RPAREN) {
		p.consumedParenTerminator = true
	}
	return &ast.ArrayLiteral{Token: tok, Elements: elements}
}

func (p *Parser) parseArrayAccess() ast.Expression {
	exp := &ast.ArrayAccess{Token: p.curToken}
	p.consumeArrayAccessFlags()
	hasLengthOp := p.consumeLengthOp()
	p.consumePreflags()

	subjectWasNested := false
	if p.subjectIsEmpty() {
		exp.Left = nil
	} else {
		// A nested `${INNER}` subject leaves curToken on the inner's
		// RBRACE. Track that so the early-return below doesn't fire
		// and skip the outer's modifier tail (`${${INNER}MOD}`).
		subjectWasNested = p.peekTokenIs(token.DollarLbrace)
		p.nextToken()
		exp.Left, exp.Index = p.parseArrayAccessSubject()
	}
	if hasLengthOp && exp.Left != nil {
		exp.Left = &ast.PrefixExpression{
			Token:    token.Token{Type: token.HASH, Literal: "#"},
			Operator: "#",
			Right:    exp.Left,
		}
	}

	// Modifier tail: `${var#glob}`, `${var##glob}`, `${var%glob}`,
	// `${var%%glob}`, `${var/pat/repl}`, `${var:-default}`,
	// `${var:+alt}`, `${var:?err}`, `${var:=default}`,
	// `${var:offset:length}` and the composed forms all introduce
	// operator tokens that parseExpression does not yet model (see
	// issue #129 for the richer design). Until that lands we walk
	// through the remaining tokens, tracking matching brace/paren
	// depth, so the closing `}` is found correctly. The AST still
	// exposes the subject — katas that only care about the variable
	// name keep working — but the modifier body is opaque.
	// When curToken is already at the closing `}` (e.g. ASSIGN's
	// RHS in `${X=}` was empty and parseExpression bailed on
	// RBRACE without consuming it), short-circuit IF this is the
	// outermost expansion — detected by peek being EOF or another
	// terminator. Inside a nested form like `${#${=name}}`, the
	// inner ArrayAccess already left curToken on its close `}`
	// while the outer's close is the next `}` (peek RBRACE), and
	// we still need expectPeek(RBRACE) to advance.
	if p.curTokenIs(token.RBRACE) && !p.peekTokenIs(token.RBRACE) && !subjectWasNested {
		return exp
	}
	if !p.peekTokenIs(token.RBRACE) {
		p.consumeArrayAccessModifierTail()
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return exp
}

// consumeArrayAccessFlags drains the optional `(flags)` group between
// `${` and the subject, used by Zsh parameter-expansion flag tuples
// like `${(j:,:)arr}`.
func (p *Parser) consumeArrayAccessFlags() {
	if !p.peekTokenIs(token.LPAREN) {
		return
	}
	p.nextToken()
	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
	}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}
}

// consumeLengthOp consumes the `#` length operator if present.
func (p *Parser) consumeLengthOp() bool {
	if p.peekTokenIs(token.HASH) {
		p.nextToken()
		return true
	}
	return false
}

// consumePreflags drains Zsh single-character pre-flags `=`, `~`, `^`
// that precede the subject inside `${ … }`.
func (p *Parser) consumePreflags() {
	for p.peekTokenIs(token.ASSIGN) || p.peekTokenIs(token.TILDE) ||
		p.peekTokenIs(token.CARET) || p.peekTokenIs(token.EQ) {
		p.nextToken()
	}
}

// subjectIsEmpty reports whether the upcoming token starts a modifier
// tail directly (no parameter name).
func (p *Parser) subjectIsEmpty() bool {
	return p.peekTokenIs(token.COLON) || p.peekTokenIs(token.HASH) ||
		p.peekTokenIs(token.PERCENT) || p.peekTokenIs(token.SLASH)
}

// parseArrayAccessSubject parses the parameter-name subject and
// returns (left, index) so callers can populate ArrayAccess.Left
// and ArrayAccess.Index.
func (p *Parser) parseArrayAccessSubject() (ast.Expression, ast.Expression) {
	switch {
	case p.curTokenIs(token.IDENT) && strings.Contains(p.curToken.Literal, "/"):
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
	case p.curTokenIs(token.IDENT):
		return p.parseArrayAccessIdent()
	case p.subjectIsSpecialName():
		return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
	case p.curTokenIs(token.DollarLbrace):
		// Nested `${INNER}` subject. Call parseArrayAccess directly
		// rather than going through parseExpression so the infix loop
		// in parseExpression doesn't eat the outer's modifier tail
		// (`%%pat`, `//pat/repl`) as misinterpreted operators.
		return p.parseArrayAccess(), nil
	}
	expr := p.parseExpression(LOWEST)
	if idx, ok := expr.(*ast.IndexExpression); ok {
		return idx.Left, idx.Index
	}
	return expr, nil
}

func (p *Parser) parseArrayAccessIdent() (ast.Expression, ast.Expression) {
	id := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.peekTokenIs(token.LBRACKET) || p.peekToken.HasPrecedingSpace {
		return id, nil
	}
	p.nextToken()
	idx, ok := p.parseIndexExpression(id).(*ast.IndexExpression)
	if !ok {
		return id, nil
	}
	return idx.Left, idx.Index
}

func (p *Parser) subjectIsSpecialName() bool {
	return p.curTokenIs(token.VARIABLE) || p.curTokenIs(token.INT) ||
		p.curTokenIs(token.ASTERISK) || p.curTokenIs(token.QUESTION) ||
		p.curTokenIs(token.MINUS) || p.curTokenIs(token.BANG)
}

func (p *Parser) consumeArrayAccessModifierTail() {
	depth := 0
	for !p.peekTokenIs(token.EOF) {
		switch {
		case p.peekTokenIs(token.DollarLbrace) || p.peekTokenIs(token.LBRACE):
			depth++
			p.nextToken()
		case p.peekTokenIs(token.RBRACE):
			if depth == 0 {
				return
			}
			depth--
			p.nextToken()
		default:
			p.nextToken()
		}
	}
}

func (p *Parser) parseInvalidArrayAccessPrefix() ast.Expression {
	dollarToken := p.curToken
	if p.peekIsDollarTerminator() {
		return &ast.Identifier{Token: dollarToken, Value: "$"}
	}
	if p.peekTokenIs(token.LBRACKET) {
		return p.parseDollarArithExpansion(dollarToken)
	}
	if p.peekIsDollarSpecialOp() {
		return p.parseDollarSpecialOp(dollarToken)
	}
	if p.peekTokenIs(token.PLUS) {
		return p.parseDollarPlusName(dollarToken)
	}
	return p.parseDollarIdent(dollarToken)
}

func (p *Parser) peekIsDollarTerminator() bool {
	switch {
	case p.peekTokenIs(token.SEMICOLON), p.peekTokenIs(token.EOF):
		return true
	case p.peekTokenIs(token.PIPE), p.peekTokenIs(token.AMPERSAND):
		return true
	case p.peekTokenIs(token.AND), p.peekTokenIs(token.OR):
		return true
	case p.peekTokenIs(token.RPAREN), p.peekTokenIs(token.RBRACE):
		return true
	case p.peekTokenIs(token.RDBRACKET), p.peekTokenIs(token.RBRACKET):
		return true
	}
	return false
}

func (p *Parser) peekIsDollarSpecialOp() bool {
	return p.peekTokenIs(token.HASH) || p.peekTokenIs(token.INT) ||
		p.peekTokenIs(token.ASTERISK) || p.peekTokenIs(token.BANG) ||
		p.peekTokenIs(token.MINUS) || p.peekTokenIs(token.CARET) ||
		p.peekTokenIs(token.EQ) || p.peekTokenIs(token.TILDE)
}

func (p *Parser) parseDollarArithExpansion(dollarToken token.Token) ast.Expression {
	p.nextToken() // onto [
	bdepth := 1
	for bdepth > 0 && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		switch {
		case p.curTokenIs(token.LBRACKET):
			bdepth++
		case p.curTokenIs(token.RBRACKET):
			bdepth--
		}
	}
	return &ast.Identifier{Token: dollarToken, Value: "$[…]"}
}

func (p *Parser) parseDollarSpecialOp(dollarToken token.Token) ast.Expression {
	p.nextToken()
	opToken := p.curToken
	if opToken.Type == token.HASH && p.peekIsHashLengthOperand() {
		p.nextToken()
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return &ast.PrefixExpression{
			Token:    dollarToken,
			Operator: "$",
			Right: &ast.PrefixExpression{
				Token:    opToken,
				Operator: "#",
				Right:    name,
			},
		}
	}
	// Zsh `$<flag>name` forms: `$^name` (array-broadcast),
	// `$=name` (split), `$~name` (glob), `$+name` already handled.
	// Absorb the IDENT name when it directly follows the flag.
	if isDollarFlagOp(opToken.Type) && p.peekTokenIs(token.IDENT) && !p.peekToken.HasPrecedingSpace {
		flagToken := opToken
		p.nextToken()
		name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return &ast.PrefixExpression{
			Token:    dollarToken,
			Operator: "$",
			Right: &ast.PrefixExpression{
				Token:    flagToken,
				Operator: flagToken.Literal,
				Right:    name,
			},
		}
	}
	ident := &ast.Identifier{Token: opToken, Value: opToken.Literal}
	return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: ident}
}

func isDollarFlagOp(t token.Type) bool {
	switch t {
	case token.CARET, token.EQ, token.TILDE:
		return true
	}
	return false
}

// peekIsHashLengthOperand reports whether the upcoming token is a valid
// operand for `$#` (length of). Zsh accepts a name (`$#name`), a
// positional digit (`$#1`), or `$#*` / `$#?` / `$##` for the count of
// the special parameter.
func (p *Parser) peekIsHashLengthOperand() bool {
	if p.peekToken.HasPrecedingSpace {
		return false
	}
	switch p.peekToken.Type {
	case token.IDENT, token.INT, token.ASTERISK, token.QUESTION, token.HASH:
		return true
	}
	return false
}

func (p *Parser) parseDollarPlusName(dollarToken token.Token) ast.Expression {
	p.nextToken()
	plusToken := p.curToken
	// Zsh accepts a name (`$+commands`), a positional digit
	// (`$+3`), or `$+*` / `$+@` / `$+?` for the existence test of the
	// special parameter. The lexer emits these as IDENT/INT plus a
	// punctuation token; we accept either.
	if !p.peekTokenIs(token.IDENT) && !p.peekTokenIs(token.INT) {
		return nil
	}
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	plus := &ast.PrefixExpression{Token: plusToken, Operator: "+", Right: ident}
	if !p.peekTokenIs(token.LBRACKET) {
		return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: plus}
	}
	p.nextToken()
	exp := &ast.InvalidArrayAccess{Token: dollarToken, Left: plus}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	return p.finalizeInvalidArrayAccess(exp)
}

func (p *Parser) parseDollarIdent(dollarToken token.Token) ast.Expression {
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.peekTokenIs(token.LBRACKET) {
		return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: ident}
	}
	p.nextToken()
	exp := &ast.InvalidArrayAccess{Token: dollarToken, Left: ident}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// finalizeInvalidArrayAccess closes the subscript body, draining
// trailing tokens to the matching `]` so the caller does not leak
// half-parsed input back into the dispatch loop.
func (p *Parser) finalizeInvalidArrayAccess(exp *ast.InvalidArrayAccess) ast.Expression {
	if !p.peekTokenIs(token.RBRACKET) {
		bdepth := 0
		for !p.peekTokenIs(token.EOF) {
			p.nextToken()
			switch {
			case p.curTokenIs(token.LBRACKET):
				bdepth++
			case p.curTokenIs(token.RBRACKET):
				if bdepth == 0 {
					return exp
				}
				bdepth--
			}
		}
		return exp
	}
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// peekIsFunctionDefinitionContinuation reports whether the token
// after `function` shapes a Zsh function definition: a name token, a
// `${…}`-spliced name, a leading `-` for dashed names, an opening
// `(` for `function name()` form, or `{` for `function { body }`.
// Used to guard parseFunctionLiteral so a stray `function` keyword
// in expression position (assignment RHS, case label) degrades to a
// literal identifier instead of erroring on the missing brace body.
func (p *Parser) peekIsFunctionDefinitionContinuation() bool {
	switch p.peekToken.Type {
	case token.IDENT, token.STRING, token.VARIABLE,
		token.DollarLbrace, token.MINUS, token.LPAREN, token.LBRACE:
		return true
	}
	return false
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	// `function` only opens a Zsh function definition when followed
	// by a name token, a `${…}`-spliced name, a leading `-` (dashed
	// name), or directly by `(` or `{`. Anywhere else — e.g. as the
	// RHS of an assignment (`REPLY=function`) or as a case-label
	// pattern — it is a literal identifier. Without this guard the
	// expectPeek(LBRACE) below errored on the next statement's
	// keyword (`expected next token to be {, got ELIF instead`).
	if !p.peekIsFunctionDefinitionContinuation() {
		tok := p.curToken
		return &ast.Identifier{Token: tok, Value: tok.Literal}
	}
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if p.peekTokenIs(token.DollarLbrace) {
		nameTok := p.peekToken
		p.nextToken()
		p.skipDollarBraceBody()
		lit.Name = &ast.Identifier{Token: nameTok, Value: nameTok.Literal}
	}
	// Zsh allows function names that start with `-` (e.g.
	// `function -coreutils-alias-setup { … }`). The lexer emits the
	// leading `-` as a MINUS token followed by an IDENT (with no
	// space between them), so glue the pair into the name.
	if p.peekTokenIs(token.MINUS) {
		dashTok := p.peekToken
		p.nextToken()
		if p.peekTokenIs(token.IDENT) && !p.peekToken.HasPrecedingSpace {
			p.nextToken()
			lit.Name = &ast.Identifier{Token: dashTok, Value: "-" + p.curToken.Literal}
			p.consumeCompositeFunctionName()
		}
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.consumeCompositeFunctionName()
	}
	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
	}
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		lit.Params = p.parseFunctionParameters()
	} else {
		lit.Params = []*ast.Identifier{}
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken()
	lit.Body = p.parseBlockStatement(token.RBRACE)
	return lit
}

// skipDollarBraceBody walks past the matching `}` of a `${ … }`
// expansion that we don't care to parse structurally.
func (p *Parser) skipDollarBraceBody() {
	depth := 1
	for depth > 0 && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		switch {
		case p.curTokenIs(token.DollarLbrace) || p.curTokenIs(token.LBRACE):
			depth++
		case p.curTokenIs(token.RBRACE):
			depth--
		}
	}
}

// consumeCompositeFunctionName absorbs adjacent (no preceding space)
// IDENT / STRING / VARIABLE / `${…}` tokens onto the function name so
// `function n${1:-}suffix() { … }` survives as a single name.
func (p *Parser) consumeCompositeFunctionName() {
	for !p.peekToken.HasPrecedingSpace {
		switch {
		case p.peekTokenIs(token.IDENT),
			p.peekTokenIs(token.STRING),
			p.peekTokenIs(token.VARIABLE):
			p.nextToken()
		case p.peekTokenIs(token.DollarLbrace):
			p.nextToken()
			p.skipDollarBraceBody()
		default:
			return
		}
	}
}

func (p *Parser) parseCommandSubstitution() ast.Expression {
	exp := &ast.CommandSubstitution{Token: p.curToken}
	p.nextToken()

	p.inBackticks++
	// Backtick body can be a multi-statement list `\`a; b; c\`` —
	// parseCommandList alone only handles one pipeline + logical
	// chain, so subsequent `;` separators fired
	// "expected `, got ;". Drain extras opaquely after the first
	// command, mirroring parseDollarParenExpression's fix.
	exp.Command = p.parseCommandList()
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		p.nextToken()
		_ = p.parseStatement()
	}
	p.inBackticks--

	if !p.expectPeek(token.BACKTICK) {
		return nil
	}
	return exp
}

func (p *Parser) parseDollarParenExpression() ast.Expression {
	exp := &ast.DollarParenExpression{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		p.nextToken() // consume `(`

		prevInArithmetic := p.inArithmetic
		p.inArithmetic = true
		p.consumeArithmeticRadixPrefix()
		cmd := p.parseExpression(LOWEST)
		p.inArithmetic = prevInArithmetic

		if p.peekTokenIs(token.DoubleRparen) {
			p.nextToken() // consume ))
			exp.Command = cmd
			return exp
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		exp.Command = cmd
		return exp
	}

	p.nextToken()
	// First parse as a command list so `$(cmd arg1 arg2 | other)`
	// returns a SimpleCommand / pipeline — detection katas like
	// ZC1050 walk `n.Command` expecting that shape. If the body
	// continues past the first pipeline with a `;` separator,
	// drain the rest opaquely so `$(cmd1; cmd2)` reaches its
	// closing `)` cleanly. The AST keeps the first command; katas
	// that care about the full body can walk source.
	exp.Command = p.parseCommandList()
	// Drain any further statements inside the `$( … )` body.
	// Zsh separates statements with `;` or a newline, so advance
	// past either and re-enter parseStatement. Stops at RPAREN or
	// EOF; callers handle unexpected EOF as a parse error.
	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		switch {
		case p.peekTokenIs(token.SEMICOLON):
			p.nextToken() // onto ;
		case p.peekToken.Line > p.curToken.Line:
			// implicit newline separator: fall through to nextToken
		default:
			// Unknown continuation — bail so the RPAREN expectPeek
			// below reports a meaningful error.
			goto drainDone
		}
		p.nextToken() // onto next stmt head
		_ = p.parseStatement()
	}
drainDone:
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	// Signal that an inner expression consumed its own closing `)`.
	// parseBlockStatement uses the flag to skip RPAREN as its own
	// terminator and advance past it instead, so an enclosing
	// `( … )` subshell body doesn't end at the inner `)` of
	// `HOST=$(cmd)`.
	p.consumedParenTerminator = true
	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken() // curToken is first IDENT
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()                   // Consume COMMA
		if p.peekTokenIs(token.IDENT) { // Expect IDENT after COMMA
			p.nextToken() // Consume IDENT
			ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			identifiers = append(identifiers, ident)
		} else {
			// Expected an identifier, got something else. Report error and break.
			p.peekError(token.IDENT) // Report expected IDENT
			break
		}
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	// Parse each argument at LOGICAL so the comma-separator (precedence
	// LOWEST+1) does not get absorbed as a binary operator inside the
	// argument expression itself.
	args = append(args, p.parseExpression(LOGICAL))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOGICAL))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseDoubleParenExpression() ast.Expression {
	p.nextToken()
	p.consumeArithmeticRadixPrefix()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	return exp
}

// consumeArithmeticRadixPrefix drains an optional Zsh arithmetic radix
// prefix `[#]`, `[#N]`, or `[##N]` that prints the result in a non-
// decimal base (`(([#16] 0xff))`). Only valid at the start of an
// arithmetic expression. Caller must have already advanced curToken
// onto the first body token.
func (p *Parser) consumeArithmeticRadixPrefix() {
	if !p.curTokenIs(token.LBRACKET) || !p.peekTokenIs(token.HASH) {
		return
	}
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
	if p.curTokenIs(token.RBRACKET) {
		p.nextToken()
	}
}

func (p *Parser) parseRedirection(left ast.Expression) ast.Expression {
	expr := &ast.Redirection{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expr.Right = p.parseExpression(LOWEST)

	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	if p.curTokenIs(token.LPAREN) {
		return p.parseFlaggedSubscript(exp)
	}
	p.parseArithmeticSubscript(exp)
	if p.curTokenIs(token.RBRACKET) && p.peekTokenIs(token.RBRACE) {
		return exp
	}
	if !p.peekTokenIs(token.RBRACKET) {
		p.drainSubscriptBody()
		return exp
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseFlaggedSubscript handles `arr[(R)pattern]` style subscripts
// where a parenthesised flag tuple precedes a glob-pattern subject.
func (p *Parser) parseFlaggedSubscript(exp *ast.IndexExpression) ast.Expression {
	depth := 1
	for depth > 0 && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		switch {
		case p.curTokenIs(token.LPAREN):
			depth++
		case p.curTokenIs(token.RPAREN):
			depth--
		}
	}
	p.nextToken()
	exp.Index = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	bdepth := 0
	for !p.curTokenIs(token.EOF) {
		switch {
		case p.curTokenIs(token.RBRACKET):
			if bdepth == 0 {
				return exp
			}
			bdepth--
		case p.curTokenIs(token.LBRACKET):
			bdepth++
		}
		p.nextToken()
	}
	return exp
}

// parseArithmeticSubscript parses the arithmetic body of a regular
// `arr[expr]` or slice `arr[a,b]` subscript, leaving curToken on the
// last consumed token.
func (p *Parser) parseArithmeticSubscript(exp *ast.IndexExpression) {
	prev := p.inArithmetic
	p.inArithmetic = true
	exp.Index = p.parseExpression(LOWEST)
	p.inArithmetic = prev
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		_ = p.parseExpression(LOWEST)
	}
}

// drainSubscriptBody walks forward to the matching `]`, tolerating
// nested `[ ]` pairs that the arithmetic parser could not consume.
func (p *Parser) drainSubscriptBody() {
	bdepth := 0
	for !p.peekTokenIs(token.EOF) {
		p.nextToken()
		switch {
		case p.curTokenIs(token.LBRACKET):
			bdepth++
		case p.curTokenIs(token.RBRACKET):
			if bdepth == 0 {
				return
			}
			bdepth--
		}
	}
}

func (p *Parser) parseProcessSubstitution() ast.Expression {
	exp := &ast.ProcessSubstitution{Token: p.curToken}
	p.nextToken()

	// Process substitution body can be a multi-statement command
	// list: `<( cmd1; cmd2; cmd3 )` or a multi-line block. Parse
	// statements until we reach the matching RPAREN. parseCommandList
	// alone only handles a single pipeline + logical chain, so
	// subsequent `;` separators were crashing as "expected ), got ;".
	statements := []ast.Statement{}
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		// Brace-form / case statements advance past their own
		// terminator and set consumedBraceTerminator. Skipping the
		// trailing nextToken keeps the next statement's head live.
		if p.consumedBraceTerminator {
			p.consumedBraceTerminator = false
			continue
		}
		p.nextToken()
	}
	if len(statements) == 1 {
		if es, ok := statements[0].(*ast.ExpressionStatement); ok {
			exp.Command = es.Expression
		}
	}
	if !p.curTokenIs(token.RPAREN) {
		p.peekError(token.RPAREN)
		return nil
	}
	return exp
}

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// The infix chain may recurse into parseInfixExpression's
	// right-hand side just before a statement terminator or
	// block-structure keyword (`]]`, `&&`, `||`, `then`, `else`,
	// `elif`, `fi`, `do`, `done`, `esac`). Bail silently in all
	// those cases so the caller's partial infix result stays
	// well-formed and the outer statement parser resumes cleanly.
	// Typical trigger: a bare `VAR=` at end of line followed by
	// the next statement's keyword, or a glob pattern inside a
	// conditional like `"foo"* && …`.
	switch p.curToken.Type {
	case token.RDBRACKET, token.AND, token.OR,
		token.THEN, token.ELSE, token.ELIF, token.Fi,
		token.DO, token.DONE, token.ESAC,
		token.SEMICOLON, token.DSEMI,
		// Statement keywords on the next line indicate the previous
		// expression's RHS was empty (`name=<NL>for x in …; do`).
		// Return nil so the InfixExpression's Right stays unset and
		// the dispatcher picks up FOR/IF/etc. as a fresh statement.
		token.FOR, token.WHILE, token.If, token.CASE,
		token.SELECT, token.LET, token.RETURN,
		token.TYPESET, token.DECLARE:
		return nil
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// Inside a `[[ … ]]` conditional, tokens without a prefix
		// (`--`, `++`, `,`, bare punctuation words) routinely
		// appear as literal test arguments: `[[ $1 == -- ]]`,
		// `[[ $x != ++ ]]`. Treat them as identifiers rather than
		// errroring so the bracket expression closes cleanly.
		if p.inDoubleBracket {
			tok := p.curToken
			return &ast.Identifier{Token: tok, Value: tok.Literal}
		}
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		if !p.inArithmetic && p.peekTokenIs(token.LBRACKET) && p.peekToken.HasPrecedingSpace {
			break
		}
		// Stop an infix chain at `]]`. Inside a `[[ … ]]`
		// conditional the expression parser would otherwise
		// consume `*]]` as an infix multiplication with a
		// non-existent right-hand side, producing
		// "no prefix parse function for ]]" on glob patterns
		// like `.*` or `*.zsh`. RDBRACKET has no precedence
		// entry, but ASTERISK's PRODUCT outranks LOWEST and
		// lures the loop in before the peek check fires.
		if p.peekTokenIs(token.RDBRACKET) {
			break
		}
		// When the infix chain has already landed curToken on
		// `]]`, the conditional is finished. Stop the outer loop
		// before it reaches across the bracket and picks up the
		// next statement's `&&` / `||` as a continuation of the
		// bracket expression. Without this, patterns that end in
		// an ASTERISK (`[[ $x = /* ]]`) left curToken on RDBRACKET
		// with peek on OR, and LOWEST's precedence table let OR
		// win, swallowing the following command into the bracket
		// body.
		if p.curTokenIs(token.RDBRACKET) {
			break
		}
		// Outside arithmetic, a `/` glued to the previous token is
		// a path separator, not a division operator. Treat
		// `$(cmd)/` or `x/y` at statement level as end-of-
		// expression so SLASH's PRODUCT precedence doesn't sweep
		// the next line's keyword (`if`, `for`, etc.) in as the
		// division RHS. Spaced forms `5 / 5` still enter the infix
		// path so arithmetic tests keep working.
		if !p.inArithmetic && p.peekTokenIs(token.SLASH) && !p.peekToken.HasPrecedingSpace {
			break
		}
		// Inside a `[[ … ]]` conditional, adjacent `(…)` groups are
		// glob alternations being concatenated — not function calls
		// on the left-hand expression. Stop the infix loop from
		// picking up the LPAREN as a CALL so parseGroupedExpression
		// handles the pattern group on its own.
		if p.inDoubleBracket && p.peekTokenIs(token.LPAREN) {
			break
		}
		// `cmd (subshell)` with space between `cmd` and `(` is a
		// command followed by a subshell argument, not a function
		// call (`cmd(args)` with no space). The caller's
		// parseSingleCommand handles argument gathering; bail out
		// here so LPAREN doesn't get treated as a CALL infix.
		if p.peekTokenIs(token.LPAREN) && p.peekToken.HasPrecedingSpace {
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
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
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
	return lit
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

	return &ast.ArrayLiteral{Token: tok, Elements: elements}
}

func (p *Parser) parseArrayAccess() ast.Expression {
	exp := &ast.ArrayAccess{Token: p.curToken}

	// Handle Zsh flags: ${(flags)...}
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // consume (
		// consume until )
		for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
			p.nextToken()
		}
		if p.peekTokenIs(token.RPAREN) {
			p.nextToken() // consume )
		}
	}

	// Handle length operator #
	hasLengthOp := false
	if p.peekTokenIs(token.HASH) {
		p.nextToken() // consume #
		hasLengthOp = true
	}

	// Zsh single-character pre-flags inside `${X name}` that modify the
	// expansion rather than naming a parameter: `=` (split), `~` (glob
	// interpret), `^` (rc-style expansion). They precede the subject and
	// are consumed without producing an AST node — detection katas that
	// care about them walk the source directly. Without this guard the
	// generic prefix-expression path rejects `=` and `^`, breaking real
	// scripts like `strategies=(${=VAR})` from zsh-autosuggestions.
	for p.peekTokenIs(token.ASSIGN) || p.peekTokenIs(token.TILDE) ||
		p.peekTokenIs(token.CARET) || p.peekTokenIs(token.EQ) {
		// `${==X}` is the double-`=` form (strip an outer `=`
		// flag); the lexer fuses `==` into a single EQ token, so
		// allow that as a pre-flag here too.
		p.nextToken()
	}

	// The subject is optional when the only body is a modifier tail
	// applied to an empty parameter, as in `${(%):-default}` where
	// the `(%)` flag group is followed directly by `:-`. Without
	// this guard, parseExpression tries to find a prefix for `:` and
	// errors out. If the peek is a modifier punctuator, skip straight
	// to the opaque modifier-tail scanner below.
	if p.peekTokenIs(token.COLON) || p.peekTokenIs(token.HASH) ||
		p.peekTokenIs(token.PERCENT) || p.peekTokenIs(token.SLASH) {
		// Nothing to parse for the subject; the modifier tail loop
		// will consume the rest of the body.
		exp.Left = nil
	} else {
		p.nextToken() // move to subject
		// Parse the subject narrowly. Using parseExpression(LOWEST)
		// pulls modifier operators (`%`, `#`, `/`) into an infix
		// chain, which then misreads patterns like `${a%%[[:space:]]*}`
		// (PERCENT then `[` → LBRACKET prefix calls parseSingleCommand
		// on the bracket class). The modifier-tail scanner below is
		// the right home for that body. Limit the subject to the
		// minimal shapes it can be — IDENT (with optional adjacent
		// subscript), VARIABLE, INT — and let the opaque scanner
		// consume the rest.
		switch {
		case p.curTokenIs(token.IDENT) && strings.Contains(p.curToken.Literal, "/"):
			// Pattern-substitution head like `line//` already
			// absorbed the slashes; bracket that follows is a
			// glob class, not a subscript.
			exp.Left = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		case p.curTokenIs(token.IDENT):
			id := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			exp.Left = id
			// Allow a single adjacent subscript so `${arr[i]}` keeps
			// its IndexExpression shape that katas walk for array
			// access. The subscript closes at the matching `]`; the
			// modifier tail (if any) starts after it.
			if p.peekTokenIs(token.LBRACKET) && !p.peekToken.HasPrecedingSpace {
				p.nextToken() // onto [
				if idx, ok := p.parseIndexExpression(id).(*ast.IndexExpression); ok {
					exp.Left = idx.Left
					exp.Index = idx.Index
				}
			}
		case p.curTokenIs(token.VARIABLE), p.curTokenIs(token.INT),
			p.curTokenIs(token.ASTERISK), p.curTokenIs(token.QUESTION),
			p.curTokenIs(token.MINUS), p.curTokenIs(token.BANG):
			// Special positional / array-style subject names like
			// `${#*}`, `${?}`, `${-}`, `${!}`. Treat the
			// punctuation as a literal subject so the modifier tail
			// scanner takes over.
			exp.Left = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		default:
			expr := p.parseExpression(LOWEST)
			if idxExpr, ok := expr.(*ast.IndexExpression); ok {
				exp.Left = idxExpr.Left
				exp.Index = idxExpr.Index
			} else {
				exp.Left = expr
			}
		}
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
	if p.curTokenIs(token.RBRACE) && !p.peekTokenIs(token.RBRACE) {
		return exp
	}
	if !p.peekTokenIs(token.RBRACE) {
		depth := 0
		for !p.peekTokenIs(token.EOF) {
			switch {
			case p.peekTokenIs(token.DollarLbrace) || p.peekTokenIs(token.LBRACE):
				depth++
				p.nextToken()
			case p.peekTokenIs(token.RBRACE):
				if depth == 0 {
					goto done
				}
				depth--
				p.nextToken()
			default:
				p.nextToken()
			}
		}
	done:
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return exp
}

func (p *Parser) parseInvalidArrayAccessPrefix() ast.Expression {
	dollarToken := p.curToken

	// A bare `$` followed by a command terminator or EOF is a
	// literal dollar character. Real code in the oh-my-zsh corpus
	// writes `echo $` (print a literal `$`) or splits long
	// expressions with `= $` at end of line. Return a $ identifier
	// so downstream walkers see a well-formed expression rather
	// than the "expected next token to be IDENT" path below.
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF) ||
		p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AMPERSAND) ||
		p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) ||
		p.peekTokenIs(token.RPAREN) || p.peekTokenIs(token.RBRACE) ||
		p.peekTokenIs(token.RDBRACKET) || p.peekTokenIs(token.RBRACKET) {
		return &ast.Identifier{Token: dollarToken, Value: "$"}
	}

	// `$[expr]` — Zsh's older / deprecated arithmetic expansion
	// form, equivalent to `$((expr))`. Consume the body opaquely
	// to the matching `]` so callers see a single Identifier and
	// the rest of the line keeps parsing.
	if p.peekTokenIs(token.LBRACKET) {
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

	if p.peekTokenIs(token.HASH) || p.peekTokenIs(token.INT) || p.peekTokenIs(token.ASTERISK) || p.peekTokenIs(token.BANG) || p.peekTokenIs(token.MINUS) {
		p.nextToken()
		opToken := p.curToken
		// `$#name` is Zsh's length-of operator. When the special char
		// is followed by an identifier, the identifier names the
		// parameter being measured and belongs to the same expression
		// — don't leak it back into the caller's token stream.
		if opToken.Type == token.HASH && p.peekTokenIs(token.IDENT) {
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
		ident := &ast.Identifier{Token: opToken, Value: opToken.Literal}
		return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: ident}
	}

	// `$+name` / `$+name[key]` — parameter-existence test, equivalent to
	// `${+name}` / `${+name[key]}`. Commonly used inside `(( ... ))`.
	if p.peekTokenIs(token.PLUS) {
		p.nextToken() // move to '+'
		plusToken := p.curToken
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		plus := &ast.PrefixExpression{Token: plusToken, Operator: "+", Right: ident}
		if !p.peekTokenIs(token.LBRACKET) {
			return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: plus}
		}
		p.nextToken() // consume [
		exp := &ast.InvalidArrayAccess{Token: dollarToken, Left: plus}
		p.nextToken()
		exp.Index = p.parseExpression(LOWEST)
		// Subscript body may carry tokens the arithmetic expression
		// parser didn't consume — e.g. `_$cmd` lexes as IDENT +
		// VARIABLE; the first IDENT returns from parseExpression
		// and the VARIABLE falls out. Drain opaquely to the
		// matching `]` so `$+name[_$cmd]` parses cleanly.
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	// Variable / expansion as the function name:
	// `function ${=X} { … }` declares functions named by the
	// split words of `$X`. Skip past the matching `}` and treat
	// the expansion as the (opaque) name.
	if p.peekTokenIs(token.DollarLbrace) {
		nameTok := p.peekToken
		p.nextToken() // onto ${
		depth := 1
		for depth > 0 && !p.peekTokenIs(token.EOF) {
			p.nextToken()
			switch {
			case p.curTokenIs(token.DollarLbrace), p.curTokenIs(token.LBRACE):
				depth++
			case p.curTokenIs(token.RBRACE):
				depth--
			}
		}
		lit.Name = &ast.Identifier{Token: nameTok, Value: nameTok.Literal}
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		// Composite function names like `function name"${1:-}"suffix() {}`
		// appear in gitstatus and other Zsh modules that scope the
		// function to a caller-provided suffix. Absorb any adjacent
		// (no preceding whitespace) word-forming tokens into the name
		// so the trailing `()` and `{` position correctly.
		for !p.peekToken.HasPrecedingSpace {
			switch {
			case p.peekTokenIs(token.IDENT),
				p.peekTokenIs(token.STRING),
				p.peekTokenIs(token.VARIABLE):
				p.nextToken()
			case p.peekTokenIs(token.DollarLbrace):
				p.nextToken() // onto ${
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
			default:
				goto fnNameDone
			}
		}
	fnNameDone:
	}

	// Multi-name definitions: `function a b c { ... }` declares the
	// same body for each of a/b/c. oh-my-zsh's prompt_info_functions
	// uses this pattern to stub out optional prompt hooks. Swallow
	// any additional space-separated identifiers before the body so
	// the parser reaches the opening `{` (or `(`) correctly; the AST
	// keeps only the first name, which is enough for kata detection.
	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
	}

	// Zsh/Bash allows `function name { ... }` without parens.
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
	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken() // onto ;
		p.nextToken() // onto next stmt head
		_ = p.parseStatement()
	}
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
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseDoubleParenExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	return exp
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

	// Zsh subscript flags: `arr[(R)value]`, `arr[(r)pat]`, `arr[(I)i]`,
	// `arr[(ri)pat]`, etc. The `(flags)` tuple precedes the actual
	// index subject and modifies how the match is performed. Consume
	// the tuple opaquely before handing the remainder to the generic
	// expression parser. Without this guard the `(…)` was parsed as a
	// grouped expression, after which the subject IDENT had nowhere
	// to land and `expectPeek(RBRACKET)` fired on that token.
	if p.curTokenIs(token.LPAREN) {
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
		// Advance onto the subject after the closing paren.
		p.nextToken()
		// When a flag tuple was present the subject is a glob
		// pattern (`arr[(r)*.zsh]`, `${1[(wr)^(*=*|sudo)]}`), not
		// an arithmetic expression. Consume the remainder of the
		// body opaquely so mixed glob alternations, negations, and
		// nested classes don't crash the arithmetic parser. The
		// AST keeps Index set to a placeholder; detection katas
		// that need the raw subscript text read it from source.
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

	prevInArithmetic := p.inArithmetic
	p.inArithmetic = true
	exp.Index = p.parseExpression(LOWEST)
	p.inArithmetic = prevInArithmetic

	// Array slices: `${arr[1,8]}`, `${arr[$a,$b]}`. The comma is
	// the Zsh range separator and the second index is the slice
	// endpoint. Skip it opaquely and consume the rest of the
	// subscript body so expectPeek(RBRACKET) lands on the closing
	// bracket. The AST keeps the first index in Index; detection
	// katas that need slice info can walk source directly.
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // onto ,
		p.nextToken() // onto second index token
		_ = p.parseExpression(LOWEST)
	}

	// Zsh associative-array keys accept arbitrary tokens — keys
	// starting with a digit (`emoji[1st_place_medal]`) tokenise as
	// INT + IDENT, and keys with punctuation (`arr[foo-bar]`,
	// `arr[x.y]`) split across multiple tokens. The arithmetic
	// parse above may have already landed curToken on the closing
	// `]` when a prefix expression's failed RHS swallowed it
	// (e.g. `${#y[*]}`: parsePrefixExpression on `*` advanced into
	// the RBRACKET). Only short-circuit when we're clearly at the
	// outermost subscript's close — detected by peek being the
	// enclosing `${…}`'s RBRACE. In nested forms like
	// `FG[$colors[color+1]]` peek is another RBRACKET (the outer
	// subscript's close) and we must continue to expectPeek it.
	if p.curTokenIs(token.RBRACKET) && p.peekTokenIs(token.RBRACE) {
		return exp
	}
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
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
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

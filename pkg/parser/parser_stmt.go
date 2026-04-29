// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package parser

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
		return nil
	}
	if stmt, ok := p.parseSimpleStatement(); ok {
		return stmt
	}
	if stmt, ok := p.parsePipelineHeadStatement(); ok {
		return stmt
	}
	return p.parseStatementBranch()
}

// parseSimpleStatement covers the cases whose dispatch is just a token
// match plus a single helper call.
func (p *Parser) parseSimpleStatement() (ast.Statement, bool) {
	switch p.curToken.Type {
	case token.RETURN:
		return p.parseReturnStatement(), true
	case token.LET:
		return p.parseLetStatement(), true
	case token.SHEBANG:
		return p.parseShebangStatement(), true
	case token.HASH:
		return nil, true
	case token.COPROC:
		return p.parseCoprocStatement(), true
	case token.TYPESET, token.DECLARE:
		return p.parseDeclarationStatement(), true
	case token.LPAREN:
		return p.parseSubshellStatement(), true
	}
	return nil, false
}

// parsePipelineHeadStatement covers block-shaped statements that may
// head a trailing pipeline tail.
func (p *Parser) parsePipelineHeadStatement() (ast.Statement, bool) {
	switch p.curToken.Type {
	case token.If:
		stmt := p.parseIfStatement()
		p.consumePipelineTail()
		return stmt, true
	case token.FOR:
		stmt := p.parseForLoopStatement()
		p.consumePipelineTail()
		return stmt, true
	case token.WHILE:
		stmt := p.parseWhileLoopStatement()
		p.consumePipelineTail()
		return stmt, true
	case token.SELECT:
		stmt := p.parseSelectStatement()
		p.consumePipelineTail()
		return stmt, true
	case token.CASE:
		stmt := p.parseCaseStatement()
		p.consumePipelineTail()
		return stmt, true
	}
	return nil, false
}

func (p *Parser) parseStatementBranch() ast.Statement {
	switch p.curToken.Type {
	case token.LBRACE:
		return p.parseBraceGroupStatement()
	case token.DoubleLparen:
		return p.parseDoubleLparenStatement()
	case token.LDBRACKET:
		return p.parseLDBracketStatement()
	case token.COLON, token.DOT, token.LBRACKET,
		token.GT, token.LT, token.GTGT, token.LTLT,
		token.GTAMP, token.LTAMP, token.AMPERSAND, token.SLASH:
		return p.parseSimpleCommandStatement()
	case token.BANG:
		return p.parseBangStatement()
	case token.BACKTICK, token.DOLLAR_LPAREN, token.VARIABLE, token.DollarLbrace:
		return p.parsePipelineStartingWithExpression()
	case token.IDENT:
		return p.parseIdentStatement()
	default:
		return p.parseExpressionOrFunctionDefinition()
	}
}

// parsePipelineHead dispatches the head expression of a pipeline.
// Returns (expr, true) when the head is a keyword-compound command
// whose parser already chained logical / redirection tails; in that
// case parseCommandPipeline returns immediately. Otherwise returns
// (expr, false) for the standard redirection / pipe-tail follow-up.
func (p *Parser) parsePipelineHead() (ast.Expression, bool) {
	switch p.curToken.Type {
	case token.WHILE:
		return p.parseWhileLoopStatement(), false
	case token.LPAREN:
		return p.parseGroupedExpression(), false
	case token.LBRACE:
		// Brace-group: `{ cmd1; cmd2 } 2>&1` appears inside `$(…)`
		// and as a pipeline head. Parse as a brace block so the
		// generic parseSingleCommand path doesn't read `{` as a
		// command name and crash on the closing `}`.
		left := keywordStmtToExpression(p.parseBraceGroupStatement())
		p.drainFDPrefixedRedirections()
		return left, false
	case token.LDBRACKET:
		return p.parseDoubleBracketExpression(), false
	case token.DoubleLparen:
		return p.parseArithmeticCommand(), false
	case token.If, token.FOR, token.CASE, token.SELECT:
		// Keyword-headed compound commands have their own parsers
		// that already chain redirection / logical tails. Return
		// (expr, true) so parseCommandPipeline exits early.
		stmt := p.parseStatement()
		return keywordStmtToExpression(stmt), true
	}
	return p.parseSingleCommand(), false
}

// drainFDPrefixedRedirections consumes trailing `N>...` / `N<...`
// redirections after a brace-group / subshell body. The default
// redirection loop in parseCommandPipeline only matches bare
// GT/GTAMP openers so an explicit FD number prefix would orphan
// without this drain.
func (p *Parser) drainFDPrefixedRedirections() {
	for p.peekTokenIs(token.INT) {
		p.nextToken() // FD number
		if !p.peekTokenIs(token.GT) && !p.peekTokenIs(token.GTGT) &&
			!p.peekTokenIs(token.GTAMP) && !p.peekTokenIs(token.LT) &&
			!p.peekTokenIs(token.LTAMP) {
			return
		}
		p.nextToken() // operator
		p.nextToken() // target
		_ = p.parseCommandWord()
	}
}

func (p *Parser) parseBraceGroupStatement() ast.Statement {
	tok := p.curToken
	p.nextToken()
	block := p.parseBlockStatement(token.RBRACE)
	block.Token = tok
	for p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
		p.nextToken()
		p.nextToken()
		_ = p.parseCommandPipeline()
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return block
}

func (p *Parser) parseDoubleLparenStatement() ast.Statement {
	cmd := p.parseArithmeticCommand()
	if cmd == nil {
		return nil
	}
	if chained := p.chainLogical(cmd, cmd.Token); chained != nil {
		return chained
	}
	return cmd
}

func (p *Parser) parseLDBracketStatement() ast.Statement {
	startTok := p.curToken
	expr := p.parseDoubleBracketExpression()
	if expr == nil {
		return nil
	}
	if chained := p.chainLogical(expr, startTok); chained != nil {
		return chained
	}
	stmt := &ast.ExpressionStatement{Token: startTok, Expression: expr}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBangStatement() ast.Statement {
	if p.peekStartsCommand() {
		return p.parseSimpleCommandStatement()
	}
	return p.parseExpressionOrFunctionDefinition()
}

func (p *Parser) peekStartsCommand() bool {
	switch {
	case p.peekTokenIs(token.IDENT), p.peekTokenIs(token.LPAREN):
		return true
	case p.peekTokenIs(token.LBRACKET), p.peekTokenIs(token.LDBRACKET):
		return true
	case p.peekTokenIs(token.DoubleLparen), p.peekTokenIs(token.VARIABLE):
		return true
	case p.peekTokenIs(token.DollarLbrace), p.peekTokenIs(token.BACKTICK):
		return true
	case p.peekTokenIs(token.DOLLAR_LPAREN):
		return true
	}
	return false
}

func (p *Parser) parseIdentStatement() ast.Statement {
	if p.curToken.Literal == "test" {
		return p.parseSimpleCommandStatement()
	}
	if p.peekStartsSimpleCommand() {
		return p.parseSimpleCommandStatement()
	}
	return p.parseExpressionOrFunctionDefinition()
}

func (p *Parser) peekStartsSimpleCommand() bool {
	if p.peekStartsArgPrefix() {
		return true
	}
	if p.peekTokenIs(token.DEC) || p.peekTokenIs(token.INC) {
		return true
	}
	return p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR)
}

func (p *Parser) peekStartsArgPrefix() bool {
	switch {
	case p.peekTokenIs(token.IDENT), p.peekTokenIs(token.STRING), p.peekTokenIs(token.INT):
		return true
	case p.peekTokenIs(token.MINUS), p.peekTokenIs(token.DOT), p.peekTokenIs(token.VARIABLE):
		return true
	case p.peekTokenIs(token.DOLLAR), p.peekTokenIs(token.DollarLbrace):
		return true
	case p.peekTokenIs(token.DOLLAR_LPAREN), p.peekTokenIs(token.SLASH):
		return true
	case p.peekTokenIs(token.TILDE), p.peekTokenIs(token.ASTERISK):
		return true
	case p.peekTokenIs(token.BANG), p.peekTokenIs(token.LBRACE):
		return true
	}
	return false
}

// chainLogical threads `&&` / `||` continuations onto an arbitrary
// left-hand expression, returning a wrapped ExpressionStatement. The
// helper exists because `(( … ))` and `[[ … ]]` are both legitimate
// starts of a logical chain but live on different parse paths; both
// now funnel through here. Returns nil when the peek is not a logical
// operator so the caller can emit its native shape untouched.
func (p *Parser) chainLogical(left ast.Expression, startTok token.Token) ast.Statement {
	if !p.peekTokenIs(token.AND) && !p.peekTokenIs(token.OR) {
		return nil
	}
	for p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
		p.nextToken()
		op := p.curToken
		p.nextToken() // move to start of right-hand command
		right := p.parseCommandPipeline()
		left = &ast.InfixExpression{
			Token:    op,
			Operator: op.Literal,
			Left:     left,
			Right:    right,
		}
	}
	stmt := &ast.ExpressionStatement{Token: startTok, Expression: left}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionOrFunctionDefinition() ast.Statement {
	stmt := p.parseExpressionStatement()

	// `IDENT= cmd | other` — Zsh inline env-prefix assignment
	// followed by a pipeline. The expression path successfully
	// parses the `IDENT=value` infix but leaves the trailing `|`
	// for the next iteration which crashes on PIPE. Drain
	// pipeline / logical continuations onto the statement so they
	// stay attached.
	p.consumePipelineTail()

	// Check if it matches function definition pattern: name()
	if call, ok := stmt.Expression.(*ast.CallExpression); ok {
		if len(call.Arguments) == 0 {
			if ident, ok := call.Function.(*ast.Identifier); ok {
				// It is `name()`. In Zsh this must be followed by a body to be a valid definition.
				// If we are here, we consumed `name` and `()`.
				// Parse the next statement as body.
				funcDef := &ast.FunctionDefinition{
					Token: ident.Token,
					Name:  ident,
				}
				// We expect a statement now.
				// If we are at semicolon, skip it? `func(); body` is valid? No. `func() body`.
				// But lexer might have produced semicolon if newline?
				// If next is `{`, `(` (subshell body), or command.

				// If we are at EOF or semicolon without body, it's just a CallExpression (incomplete func def).
				if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.EOF) {
					return stmt
				}

				// Parse body
				p.nextToken() // Move to start of body statement
				funcDef.Body = p.parseStatement()
				return funcDef
			}
		}
	}
	return stmt
}

// keywordStmtToExpression wraps a Statement returned by parseStatement
// (typically an IfStatement / ForLoopStatement / CaseStatement /
// SelectStatement) as an Expression so it can flow through pipeline
// chaining where parseCommandPipeline expects an ast.Expression.
// Detection katas that walk the statement type still see the original
// node when traversing the wrapper.
func keywordStmtToExpression(stmt ast.Statement) ast.Expression {
	if stmt == nil {
		return nil
	}
	if es, ok := stmt.(*ast.ExpressionStatement); ok {
		return es.Expression
	}
	// Wrap in a stub Identifier so callers see a non-nil expression.
	// The Token preserves the head keyword for kata-side walks of
	// containing CallExpression / DollarParenExpression bodies.
	return &ast.Identifier{Token: stmt.TokenLiteralNode(), Value: stmt.TokenLiteral()}
}

// consumePipelineTail drains trailing `| cmd` / `&& cmd` / `|| cmd`
// continuations that follow a block-shaped statement (if/for/while/
// case). These structures can head pipelines in Zsh
// (`for f in *; do …; done | column -t`) but have no AST node for a
// pipeline with a block left-hand side, so the continuation is
// consumed opaquely. Detection katas that need the full pipeline
// walk source directly.
func (p *Parser) consumePipelineTail() {
	for p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
		p.nextToken() // onto op
		p.nextToken() // onto RHS head
		_ = p.parseCommandPipeline()
	}
}

// parsePipelineStartingWithExpression parses a statement whose head
// is a command-producing expression (backtick or `$(…)`) and then
// folds any trailing pipeline / logical chain onto it. The generic
// parseSingleCommand path can't handle this because it expects the
// head to be an IDENT; doing the expression parse first and grafting
// the pipeline on top keeps the AST shape identical to what the
// IDENT path would produce for `cmd | other`.
func (p *Parser) parsePipelineStartingWithExpression() ast.Statement {
	tok := p.curToken
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}
	for p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
		p.nextToken()
		op := p.curToken
		p.nextToken()
		right := p.parseCommandPipeline()
		expr = &ast.InfixExpression{Token: op, Operator: op.Literal, Left: expr, Right: right}
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return &ast.ExpressionStatement{Token: tok, Expression: expr}
}

func (p *Parser) parseSimpleCommandStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseCommandList()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseCommandList() ast.Expression {
	left := p.parseCommandPipeline()

	for p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
		p.nextToken()
		op := p.curToken
		p.nextToken() // move to start of right command
		right := p.parseCommandPipeline()
		left = &ast.InfixExpression{
			Token:    op,
			Operator: op.Literal,
			Left:     left,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseCommandPipeline() ast.Expression {
	if p.curTokenIs(token.BANG) {
		tok := p.curToken
		p.nextToken()
		right := p.parseCommandPipeline()
		return &ast.PrefixExpression{Token: tok, Operator: "!", Right: right}
	}

	left, returned := p.parsePipelineHead()
	if returned {
		return left
	}

	// Parse redirections
	for p.peekTokenIs(token.GT) || p.peekTokenIs(token.GTGT) ||
		p.peekTokenIs(token.LT) || p.peekTokenIs(token.LTLT) ||
		p.peekTokenIs(token.GTAMP) || p.peekTokenIs(token.LTAMP) {

		p.nextToken()
		op := p.curToken
		p.nextToken() // consume op
		// Redirection target is file/expression. Use parseCommandWord to handle paths/strings correctly.
		right := p.parseCommandWord()

		left = &ast.Redirection{
			Token:    op,
			Left:     left,
			Operator: op.Literal,
			Right:    right,
		}
	}

	if p.peekTokenIs(token.PIPE) && p.peekPrecedence() == LOWEST+1 {
		p.nextToken() // consume '|'
		op := p.curToken
		p.nextToken() // move to the start of the next command
		right := p.parseCommandPipeline()
		left = &ast.InfixExpression{
			Token:    op,
			Operator: op.Literal,
			Left:     left,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseSingleCommand() ast.Expression {
	// When the head is a command-producing expression (`$(cmd)`,
	// `` `cmd` ``, `$VAR`, `${name}`), let the prefix parser run
	// so DollarParenExpression / CommandSubstitution / Identifier
	// is captured properly. Without this, the head was forced into
	// an Identifier whose Value was the literal `$(` and the
	// trailing `)` of the substitution leaked back to the dispatch
	// loop. parseSimpleCommand still wraps the result so downstream
	// argument gathering and pipeline chaining keep working.
	if p.curTokenIs(token.DOLLAR_LPAREN) || p.curTokenIs(token.BACKTICK) ||
		p.curTokenIs(token.VARIABLE) || p.curTokenIs(token.DollarLbrace) {
		startTok := p.curToken
		head := p.parseExpression(LOWEST)
		cmd := &ast.SimpleCommand{Token: startTok, Name: head, Arguments: []ast.Expression{}}
		for !p.isCommandDelimiter(p.peekToken) && p.peekOnSameLogicalLine() {
			p.nextToken()
			cmd.Arguments = append(cmd.Arguments, p.parseCommandWord())
		}
		return cmd
	}
	cmd := &ast.SimpleCommand{
		Token: p.curToken,
		Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	// Check for function definition syntax: name()
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // curToken is (
		if p.peekTokenIs(token.RPAREN) {
			p.nextToken() // curToken is )
			// Function Definition!
			funcDef := &ast.FunctionDefinition{
				Token: cmd.Token,
				Name:  cmd.Name.(*ast.Identifier),
			}

			p.nextToken() // Move to start of body
			funcDef.Body = p.parseStatement()
			return funcDef
		} else {
			// It was not (), it was `name ( ...`
			arg := p.parseCommandWord()
			cmd.Arguments = append(cmd.Arguments, arg)
		}
	} else {
		cmd.Arguments = []ast.Expression{}
	}

	// Continue parsing arguments
	for !p.isCommandDelimiter(p.peekToken) && p.peekOnSameLogicalLine() {
		p.nextToken()
		arg := p.parseCommandWord()
		cmd.Arguments = append(cmd.Arguments, arg)
	}

	return cmd
}

// commandWordLiteralTokens is the set of token types that appear inside
// a command-argument word but should be emitted as StringLiterals,
// not parsed as expressions. Reserved words (do/done/then/etc.) live
// in the same set because they routinely show up as literal arguments
// (`print -l function`, `arr=(do done)`).
var commandWordLiteralTokens = map[token.Type]struct{}{
	token.ASTERISK: {}, token.QUESTION: {}, token.PLUS: {},
	token.MINUS: {}, token.CARET: {}, token.TILDE: {}, token.DOT: {},
	token.GT: {}, token.LT: {}, token.AMPERSAND: {},
	token.LBRACKET: {}, token.RBRACKET: {}, token.HASH: {},
	token.COMMA: {}, token.COLON: {}, token.GTGT: {}, token.LTLT: {},
	token.GTAMP: {}, token.LTAMP: {},
	token.DEC: {}, token.INC: {},
	token.ASSIGN: {}, token.PLUSEQ: {},
	token.FUNCTION: {}, token.SELECT: {}, token.COPROC: {},
	token.DO: {}, token.DONE: {}, token.ESAC: {},
	token.THEN: {}, token.ELSE: {}, token.ELIF: {}, token.Fi: {},
	token.If: {}, token.FOR: {}, token.WHILE: {}, token.CASE: {},
	token.IN: {}, token.LET: {}, token.RETURN: {},
	token.TYPESET: {}, token.DECLARE: {},
}

func (p *Parser) commandWordIsExpression(t token.Type) bool {
	if _, hit := commandWordLiteralTokens[t]; hit {
		return false
	}
	return p.prefixParseFns[t] != nil
}

func (p *Parser) parseCommandWord() ast.Expression {
	firstToken := p.curToken
	// Seed braceDepth with the first token so a word that opens with
	// `{` (brace expansion or literal) keeps its closing `}` glued
	// to the same word. Without this, `{a}` as a command arg ended
	// at the first `}` (RBRACE is a command delimiter), splitting
	// the word and confusing the surrounding parse.
	braceDepth := updateCommandWordBraceDepth(0, p.curToken.Type)
	parts := []ast.Expression{p.parseCommandWordPart()}
	for p.commandWordContinues(braceDepth) {
		p.nextToken()
		braceDepth = updateCommandWordBraceDepth(braceDepth, p.curToken.Type)
		parts = append(parts, p.parseCommandWordPart())
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return &ast.ConcatenatedExpression{Token: firstToken, Parts: parts}
}

// parseCommandWordPart consumes the current token as either a literal
// string part or a sub-expression, depending on commandWordIsExpression.
func (p *Parser) parseCommandWordPart() ast.Expression {
	if !p.commandWordIsExpression(p.curToken.Type) {
		return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	}
	return p.parseExpression(CALL)
}

// commandWordContinues reports whether the next token belongs to the
// current word. The brace-depth context lets `@{upstream}` and similar
// mid-word brace runs survive without splitting at the inner `}`.
func (p *Parser) commandWordContinues(braceDepth int) bool {
	if p.peekToken.HasPrecedingSpace || !p.peekOnSameLogicalLine() {
		return false
	}
	// Zsh glob qualifiers `#` / `##` attach to the preceding pattern
	// character (`]`, `*`, `?`, `+`, `)`, an IDENT) without a space,
	// e.g. `[[:space:]]##` or `(a|b)#`. Without this exception
	// isCommandDelimiter would split the word at the HASH because
	// HASH starts a comment in command position.
	if p.peekTokenIs(token.HASH) && p.curIsGlobQualifierLeft() {
		return true
	}
	if braceDepth == 0 && p.isCommandDelimiter(p.peekToken) {
		return false
	}
	if braceDepth > 0 && p.peekTokenIs(token.EOF) {
		return false
	}
	return true
}

// curIsGlobQualifierLeft reports whether curToken is a token type that
// can carry a trailing `#` / `##` glob-qualifier without a space.
// HASH itself is included so the second `#` of a `##` doubled
// qualifier glues onto the first.
func (p *Parser) curIsGlobQualifierLeft() bool {
	switch p.curToken.Type {
	case token.RBRACKET, token.RPAREN, token.ASTERISK, token.QUESTION,
		token.PLUS, token.IDENT, token.STRING, token.HASH:
		return true
	}
	return false
}

func updateCommandWordBraceDepth(depth int, t token.Type) int {
	switch t {
	case token.LBRACE:
		return depth + 1
	case token.RBRACE:
		if depth > 0 {
			return depth - 1
		}
	}
	return depth
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Zsh `let` accepts `local` / `-i` / `-x` etc. modifier words
	// between the keyword and the assignment target, e.g.
	// `let local elapsed=1`. Skip leading IDENT modifiers (but not
	// the final name, which is the one paired with `=`).
	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken() // consume '='

	prevInArithmetic := p.inArithmetic
	p.inArithmetic = true
	stmt.Value = p.parseExpression(LOWEST)
	p.inArithmetic = prevInArithmetic

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.curToken}
	p.nextToken()
	// RPAREN / RBRACE join the cond terminator set so the Zsh
	// shortcut `if (( cond )) cmd` inside `=( … )` / `( … )` /
	// function bodies hands control back to the enclosing
	// construct once the `then`-less form ends. Fi / DONE / ESAC /
	// ELSE / ELIF cover the case where the shortcut sits inside an
	// outer if/loop/case body — without them parseBlockStatement
	// would absorb the outer construct's terminator into the cond.
	stmt.Condition = p.parseBlockStatement(token.THEN, token.LBRACE,
		token.RPAREN, token.RBRACE,
		token.Fi, token.DONE, token.ESAC, token.ELSE, token.ELIF)

	// Zsh short form `if cond { body } [elif cond { body }]…
	// [else { body }]` uses brace blocks instead of `then … fi`.
	// Detect by curToken being LBRACE after the condition.
	if p.curTokenIs(token.LBRACE) {
		p.nextToken() // into body
		stmt.Consequence = p.parseBlockStatement(token.RBRACE)
		if alt := p.parseBraceFormElifChain(); alt != nil {
			stmt.Alternative = alt
		}
		// Step past the closing RBRACE so an enclosing brace-scoped
		// body (for/while/subshell brace body) does not mistake the
		// if's terminator for its own. The consumedBraceTerminator
		// flag tells parseBlockStatement to skip its usual post-
		// statement nextToken — otherwise we'd overshoot the next
		// statement's head.
		if p.curTokenIs(token.RBRACE) {
			p.nextToken()
			p.consumedBraceTerminator = true
		}
		return stmt
	}

	if !p.curTokenIs(token.THEN) {
		if p.tryDegradeNoThenShortcut(stmt) {
			return stmt
		}
		return nil
	}

	p.nextToken() // consume "then"
	stmt.Consequence = p.parseBlockStatement(token.ELSE, token.ELIF, token.Fi)

	// Collapse any chain of `elif CONDITION; then BODY` clauses into a
	// right-nested IfStatement stored on the outer `Alternative`. We
	// thread the latest elif so the next one can attach to it.
	var tailElif *ast.IfStatement
	for p.curTokenIs(token.ELIF) {
		elifToken := p.curToken
		p.nextToken() // consume "elif"
		elif := &ast.IfStatement{Token: elifToken}
		elif.Condition = p.parseBlockStatement(token.THEN)
		if !p.curTokenIs(token.THEN) {
			return nil
		}
		p.nextToken() // consume "then"
		elif.Consequence = p.parseBlockStatement(token.ELSE, token.ELIF, token.Fi)

		if tailElif == nil {
			stmt.Alternative = elif
		} else {
			tailElif.Alternative = elif
		}
		tailElif = elif
	}

	if p.curTokenIs(token.ELSE) {
		p.nextToken() // consume "else"
		tail := p.parseBlockStatement(token.Fi)
		if tailElif == nil {
			stmt.Alternative = tail
		} else {
			tailElif.Alternative = tail
		}
	}
	if !p.curTokenIs(token.Fi) {
		p.peekError(token.Fi)
		return nil
	}
	return stmt
}

// tryDegradeNoThenShortcut handles the Zsh `if (( cond )) cmd` /
// `if [[ cond ]] cmd` shortcut. parseBlockStatement absorbs the
// trailing cmd into the cond block; once cur lands on the enclosing
// terminator (`)`, `}`, EOF, or an outer keyword like `fi`/`done`/
// `esac`/`else`/`elif`) we hand control back so the surrounding
// construct closes cleanly. Promotes the absorbed last cond statement
// to Consequence so the AST still records the body. Returns true
// when the shortcut shape applies.
func (p *Parser) tryDegradeNoThenShortcut(stmt *ast.IfStatement) bool {
	switch p.curToken.Type {
	case token.RPAREN, token.RBRACE, token.EOF,
		token.Fi, token.DONE, token.ESAC, token.ELSE, token.ELIF:
	default:
		return false
	}
	if cond, ok := stmt.Condition.(*ast.BlockStatement); ok && len(cond.Statements) >= 2 {
		last := cond.Statements[len(cond.Statements)-1]
		cond.Statements = cond.Statements[:len(cond.Statements)-1]
		stmt.Consequence = &ast.BlockStatement{Statements: []ast.Statement{last}}
	}
	return true
}

// parseBraceFormElifChain walks any `} elif COND { BODY }` chain plus
// an optional trailing `} else { BODY }` for the Zsh brace-form `if`.
// Caller has just parsed the consequence body and curToken is RBRACE.
// Returns the head Alternative (right-nested IfStatement chain), or
// nil when no chain is present.
func (p *Parser) parseBraceFormElifChain() ast.Statement {
	var head ast.Statement
	var tail *ast.IfStatement
	for p.peekTokenIs(token.ELIF) {
		p.nextToken() // onto elif
		elifTok := p.curToken
		p.nextToken() // into condition
		elif := &ast.IfStatement{Token: elifTok}
		elif.Condition = p.parseBlockStatement(token.LBRACE)
		if !p.curTokenIs(token.LBRACE) {
			return head
		}
		p.nextToken() // into body
		elif.Consequence = p.parseBlockStatement(token.RBRACE)
		if head == nil {
			head = elif
		} else {
			tail.Alternative = elif
		}
		tail = elif
	}
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // onto else
		p.nextToken() // expect {
		if p.curTokenIs(token.LBRACE) {
			p.nextToken() // into else body
			elseBody := p.parseBlockStatement(token.RBRACE)
			if head == nil {
				return elseBody
			}
			tail.Alternative = elseBody
		}
	}
	return head
}

func (p *Parser) parseBlockStatement(terminators ...token.Type) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		isTerm := false
		for _, t := range terminators {
			if p.curTokenIs(t) {
				isTerm = true
				break
			}
		}
		if isTerm {
			break
		}

		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}

		// Clear the consumedParenTerminator signal before each
		// statement so it reflects only the statement we're about
		// to parse. Without this, an inner `$(cmd)` inside a
		// different construct (e.g. `(( x = $(cmd) ))`) leaves the
		// flag set and a later iteration misfires the RPAREN skip.
		p.consumedParenTerminator = false

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		// An inner expression consumed its own RPAREN (e.g.
		// `HOST=$(cmd)` ends with curToken on the `$(…)`'s `)`).
		// That RPAREN is not this block's terminator even when
		// RPAREN is in our terminator set. Skip the RPAREN and the
		// usual nextToken so the block keeps parsing the real next
		// statement. When the flag is set but curToken is not a
		// bare `)` (the inner `$()` was wrapped by `((…))` or a
		// larger construct), the outer wrapper already accounted
		// for the `)`, so clear the flag and fall through to the
		// normal advance.
		if p.consumedParenTerminator {
			p.consumedParenTerminator = false
			if p.curTokenIs(token.RPAREN) {
				p.nextToken()
				continue
			}
		}

		// A statement whose natural terminator is the block
		// terminator itself (e.g. bare `return` / `break` right
		// before `fi`, `done`, `esac`) leaves curToken sitting on
		// that terminator because parseExpression/LOWEST yields nil
		// on block keywords. Advancing unconditionally here would
		// step past it and the outer if/loop/case would then see
		// EOF and report "expected FI got EOF".
		curIsTerm := false
		for _, t := range terminators {
			if p.curTokenIs(t) {
				// An RBRACE that closed a `${…}` is not the block's
				// terminator — `cmd ${X} }` leaves curToken on the
				// `${X}`'s `}` while the brace-block close is the
				// next RBRACE. The lexer flags the inner one.
				if t == token.RBRACE && p.curToken.ClosesDollarBrace {
					break
				}
				curIsTerm = true
				break
			}
		}
		if curIsTerm {
			break
		}

		// Brace-form statements (e.g. `if cond { body }`) advance
		// past their own closing RBRACE and set this flag so we do
		// not double-advance past the following statement's head.
		if p.consumedBraceTerminator {
			p.consumedBraceTerminator = false
			continue
		}

		p.nextToken()
	}
	return block
}

func (p *Parser) parseSubshellStatement() ast.Statement {
	subshellToken := p.curToken
	p.nextToken()
	// Anonymous-function form `() { body }` — empty parens
	// followed by an opening brace. parseStatement routes LPAREN
	// here, so detect that pattern and parse the body as a brace
	// block rather than a subshell.
	if p.curTokenIs(token.RPAREN) && p.peekTokenIs(token.LBRACE) {
		p.nextToken() // onto {
		p.nextToken() // into body
		body := p.parseBlockStatement(token.RBRACE)
		return &ast.FunctionDefinition{Token: subshellToken, Body: body}
	}
	block := p.parseBlockStatement(token.RPAREN)
	if !p.curTokenIs(token.RPAREN) {
		p.peekError(token.RPAREN)
		return nil
	}
	p.nextToken()
	// Return a Subshell node instead of BlockStatement
	return &ast.Subshell{Token: subshellToken, Command: block}
}

func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	stmt := &ast.CaseStatement{Token: p.curToken}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken()
	for !p.curTokenIs(token.ESAC) && !p.curTokenIs(token.EOF) {
		for p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		if p.curTokenIs(token.ESAC) {
			break
		}
		clause := p.parseCaseClause()
		if clause == nil {
			return nil
		}
		stmt.Clauses = append(stmt.Clauses, clause)
		if p.curTokenIs(token.DSEMI) {
			p.nextToken()
		}
	}
	// Consume the closing ESAC so a caller parsing a nested case
	// (`case x in a) case y in …;; esac ;; esac`) doesn't see the
	// inner `esac` as the outer's terminator. parseBlockStatement
	// includes ESAC in its terminator set so a clause body whose
	// final `;;` is replaced by `esac` stops cleanly; without the
	// advance below, an inner `esac` would close the outer body too.
	// Skip the advance when the successor is a pipeline / logical
	// continuation — consumePipelineTail expects peek=PIPE/AND/OR
	// — or EOF, where there is nothing to advance onto.
	if p.curTokenIs(token.ESAC) && p.shouldAdvancePastEsac() {
		p.nextToken()
		p.consumedBraceTerminator = true
	}
	return stmt
}

// shouldAdvancePastEsac reports whether parseCaseStatement should step
// past ESAC. It must NOT advance when peek is a pipeline / logical
// continuation (those need peek=PIPE/AND/OR for consumePipelineTail)
// or EOF (nothing to advance onto).
func (p *Parser) shouldAdvancePastEsac() bool {
	switch p.peekToken.Type {
	case token.EOF, token.PIPE, token.AND, token.OR:
		return false
	}
	return true
}

func (p *Parser) parseShebangStatement() *ast.Shebang {
	return &ast.Shebang{Token: p.curToken, Path: p.curToken.Literal}
}

func (p *Parser) parseCaseClause() *ast.CaseClause {
	clause := &ast.CaseClause{Token: p.curToken}
	if p.curTokenIs(token.LPAREN) && p.lookaheadCaseLabelOpener() {
		p.nextToken()
	}
	// Zsh allows `((alt|alt))` glob patterns as a case label (e.g.
	// `((add-|)fpath)`). The lexer fuses the leading `((` into
	// DoubleLparen which the regular pattern parser would route into
	// arithmetic. Re-enter the normal pattern path with curToken on
	// the first `(` of the alternation by leaving the DoubleLparen
	// in place and letting parseCommandWord drive — but first detect
	// the case-label-opener form and drop the fused token.
	if p.curTokenIs(token.DoubleLparen) {
		// Treat the fused `((` as a leading `(` opener (case-label
		// open paren) followed by an arithmetic-style `(...)` group
		// for the pattern. The pattern parsing loop in parseCommandWord
		// handles plain `(` via parseGroupedExpression; emit a
		// synthetic LPAREN by overwriting the token type.
		p.curToken.Type = token.LPAREN
		p.curToken.Literal = "("
	}
	clause.Patterns = p.parseCaseClausePatterns()
	if !p.alignToCaseClauseClose() {
		return nil
	}
	p.nextToken()
	clause.Body = p.parseBlockStatement(token.DSEMI, token.ESAC)
	return clause
}

func (p *Parser) parseCaseClausePatterns() []ast.Expression {
	var patterns []ast.Expression
	for {
		patterns = append(patterns, p.parseCommandWord())
		if !p.peekTokenIs(token.PIPE) {
			return patterns
		}
		p.nextToken()
		p.nextToken()
	}
}

// alignToCaseClauseClose advances curToken to the `)` that closes the
// label, handling the inline-glob-group case where parseCommandWord
// already consumed an inner `)`.
func (p *Parser) alignToCaseClauseClose() bool {
	if p.curTokenIs(token.RPAREN) && p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return true
	}
	if p.curTokenIs(token.RPAREN) {
		return true
	}
	return p.expectPeek(token.RPAREN)
}

func (p *Parser) parseForLoopStatement() *ast.ForLoopStatement {
	stmt := &ast.ForLoopStatement{Token: p.curToken}
	if p.peekTokenIs(token.DoubleLparen) {
		return p.parseArithmeticForLoop(stmt)
	}
	if !p.consumeForLoopName(stmt) {
		return nil
	}
	if p.peekTokenIs(token.LPAREN) {
		return p.parseShortFormForLoop(stmt)
	}
	if p.peekTokenIs(token.IN) {
		p.consumeForLoopInItems(stmt)
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return p.consumeForLoopBody(stmt)
}

func (p *Parser) parseArithmeticForLoop(stmt *ast.ForLoopStatement) *ast.ForLoopStatement {
	p.nextToken() // consume ((
	if !p.parseArithSlot(&stmt.Init, token.SEMICOLON) {
		return nil
	}
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	if !p.parseArithSlot(&stmt.Condition, token.SEMICOLON) {
		return nil
	}
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	if !p.parseArithSlot(&stmt.Post, token.DoubleRparen) {
		return nil
	}
	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	// Zsh short-form arithmetic for: `for ((..)) { body }`. Detect a
	// LBRACE peek and treat it as the body opener instead of `do`.
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken() // onto {
		p.nextToken() // into body
		stmt.Body = p.parseBlockStatement(token.RBRACE)
		if p.curTokenIs(token.RBRACE) {
			p.nextToken()
			p.consumedBraceTerminator = true
		}
		return stmt
	}
	if !p.expectPeek(token.DO) {
		return nil
	}
	p.nextToken()
	stmt.Body = p.parseBlockStatement(token.DONE)
	return stmt
}

// parseArithSlot fills the optional init / cond / post slot of an
// arithmetic for loop. Empty slot = peek already at terminator.
func (p *Parser) parseArithSlot(target *ast.Expression, terminator token.Type) bool {
	if p.peekTokenIs(terminator) {
		return true
	}
	p.nextToken()
	if p.prefixParseFns[p.curToken.Type] == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return false
	}
	*target = p.parseExpression(LOWEST)
	return true
}

func (p *Parser) consumeForLoopName(stmt *ast.ForLoopStatement) bool {
	if !p.peekTokenIs(token.IDENT) && !p.peekTokenIs(token.INT) {
		p.peekError(token.IDENT)
		return false
	}
	p.nextToken()
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	for p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.INT) {
		p.nextToken()
	}
	return true
}

func (p *Parser) parseShortFormForLoop(stmt *ast.ForLoopStatement) *ast.ForLoopStatement {
	p.nextToken() // consume (
	stmt.Items = []ast.Expression{}
	for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		if p.curTokenIs(token.RPAREN) {
			break
		}
		stmt.Items = append(stmt.Items, p.parseCommandWord())
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if p.peekTokenIs(token.DO) {
		p.nextToken()
		p.nextToken()
		stmt.Body = p.parseBlockStatement(token.DONE)
		return stmt
	}
	p.nextToken()
	stmt.Body = wrapForLoopBody(stmt.Token, p.parseStatement())
	return stmt
}

func (p *Parser) consumeForLoopInItems(stmt *ast.ForLoopStatement) {
	p.nextToken()
	stmt.Items = []ast.Expression{}
	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.DO) &&
		!p.peekTokenIs(token.EOF) && p.peekOnSameLogicalLine() {
		p.nextToken()
		stmt.Items = append(stmt.Items, p.parseCommandWord())
	}
}

func (p *Parser) consumeForLoopBody(stmt *ast.ForLoopStatement) *ast.ForLoopStatement {
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		p.nextToken()
		stmt.Body = p.parseBlockStatement(token.RBRACE)
		return stmt
	}
	if !p.peekTokenIs(token.DO) && !p.peekTokenIs(token.EOF) && !p.peekOnSameLogicalLine() {
		p.nextToken()
		stmt.Body = wrapForLoopBody(stmt.Token, p.parseStatement())
		return stmt
	}
	if !p.expectPeek(token.DO) {
		return nil
	}
	p.nextToken()
	stmt.Body = p.parseBlockStatement(token.DONE)
	return stmt
}

// wrapForLoopBody normalises a single-statement body into a
// BlockStatement so ForLoopStatement.Body stays homogeneous.
func wrapForLoopBody(tok token.Token, body ast.Statement) *ast.BlockStatement {
	if body == nil {
		return nil
	}
	if block, ok := body.(*ast.BlockStatement); ok {
		return block
	}
	return &ast.BlockStatement{Token: tok, Statements: []ast.Statement{body}}
}

func (p *Parser) parseWhileLoopStatement() *ast.WhileLoopStatement {
	stmt := &ast.WhileLoopStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseBlockStatement(token.DO, token.LBRACE)
	// Zsh short-form `while cond { body }`.
	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
		stmt.Body = p.parseBlockStatement(token.RBRACE)
		if p.curTokenIs(token.RBRACE) {
			p.nextToken()
			p.consumedBraceTerminator = true
		}
		return stmt
	}
	if !p.curTokenIs(token.DO) {
		return nil
	}
	p.nextToken()
	stmt.Body = p.parseBlockStatement(token.DONE)
	return stmt
}

func (p *Parser) parseSelectStatement() *ast.SelectStatement {
	stmt := &ast.SelectStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.IN) {
		p.nextToken()
		stmt.Items = []ast.Expression{}
		for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.DO) && !p.peekTokenIs(token.EOF) &&
			p.peekOnSameLogicalLine() {
			p.nextToken()
			arg := p.parseCommandWord()
			stmt.Items = append(stmt.Items, arg)
		}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken()
	stmt.Body = p.parseBlockStatement(token.DONE)
	return stmt
}

func (p *Parser) parseCoprocStatement() *ast.CoprocStatement {
	stmt := &ast.CoprocStatement{Token: p.curToken}
	p.nextToken()
	// Handle optional name (Bash style: coproc name { ... })
	// If next is IDENT and next-next is LBRACE?
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.LBRACE) {
		stmt.Name = p.curToken.Literal
		p.nextToken()
	}

	// Parse the command/statement
	stmt.Command = p.parseStatement()
	return stmt
}

func (p *Parser) parseDeclarationStatement() *ast.DeclarationStatement {
	stmt := &ast.DeclarationStatement{Token: p.curToken, Command: p.curToken.Literal}
	startLine := stmt.Token.Line
	p.nextToken()
	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) && p.curToken.Line == startLine {
		if p.declarationConsumeFlag(stmt, startLine) {
			continue
		}
		if p.declarationConsumeAssignment(stmt, startLine) {
			continue
		}
		// Unknown token inside a declaration — stop the loop so we do
		// not skip tokens belonging to the next statement.
		break
	}
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseArithmeticCommand() *ast.ArithmeticCommand {
	cmd := &ast.ArithmeticCommand{Token: p.curToken}
	p.nextToken()

	prevInArithmetic := p.inArithmetic
	p.inArithmetic = true
	p.consumeArithmeticRadixPrefix()
	cmd.Expression = p.parseExpression(LOWEST)
	p.inArithmetic = prevInArithmetic

	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	return cmd
}

// declarationAdvanceOrStop consumes one token within a declaration
// statement. Returns false when the next token belongs to a following
// statement (different line, terminator, or EOF) so the caller breaks
// out of the loop without over-consuming.
func (p *Parser) declarationAdvanceOrStop(startLine int) bool {
	if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF) {
		return false
	}
	if p.peekToken.Line != startLine {
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) declarationConsumeFlag(stmt *ast.DeclarationStatement, startLine int) bool {
	if !p.declarationCurIsFlag() {
		return false
	}
	stmt.Flags = append(stmt.Flags, p.curToken.Literal)
	if !p.declarationAdvanceOrStop(startLine) {
		return false
	}
	return true
}

func (p *Parser) declarationCurIsFlag() bool {
	if p.curTokenIs(token.MINUS) {
		return true
	}
	if !p.curTokenIs(token.IDENT) {
		return false
	}
	return len(p.curToken.Literal) > 0 && p.curToken.Literal[0] == '-'
}

func (p *Parser) declarationConsumeAssignment(stmt *ast.DeclarationStatement, startLine int) bool {
	if !p.declarationCurIsName() {
		return false
	}
	assign := &ast.DeclarationAssignment{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	p.consumeCompositeName(startLine)
	p.consumeDeclarationValueTail(assign, startLine)
	stmt.Assignments = append(stmt.Assignments, assign)
	if !p.declarationAdvanceOrStop(startLine) {
		return false
	}
	return true
}

func (p *Parser) declarationCurIsName() bool {
	return p.curTokenIs(token.IDENT) || p.curTokenIs(token.STRING) ||
		p.curTokenIs(token.VARIABLE) || p.curTokenIs(token.DollarLbrace)
}

// consumeCompositeName stitches adjacent (no preceding space) word
// parts onto the declaration name so `prefix"${1:-}"_suffix` and
// `${arr[$i]}` survive as one logical name.
func (p *Parser) consumeCompositeName(startLine int) {
	for !p.peekToken.HasPrecedingSpace && p.peekToken.Line == startLine {
		switch {
		case p.peekTokenIs(token.LBRACE):
			p.nextToken()
			p.consumeBraceTail(token.LBRACE, token.RBRACE)
		case p.peekTokenIs(token.DollarLbrace):
			p.nextToken()
			p.consumeBraceTail(token.LBRACE, token.RBRACE)
		case p.peekTokenIs(token.STRING),
			p.peekTokenIs(token.IDENT),
			p.peekTokenIs(token.VARIABLE):
			p.nextToken()
		default:
			return
		}
	}
}

// consumeBraceTail walks past the matching close brace of a brace
// or `${ … }` group, treating both forms as the same nesting kind.
func (p *Parser) consumeBraceTail(_ token.Type, closeT token.Type) {
	depth := 1
	for depth > 0 && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		switch {
		case p.curTokenIs(token.DollarLbrace) || p.curTokenIs(token.LBRACE):
			depth++
		case p.curTokenIs(closeT):
			depth--
		}
	}
}

func (p *Parser) consumeDeclarationValueTail(assign *ast.DeclarationAssignment, startLine int) {
	switch {
	case p.peekTokenIs(token.PLUSEQ):
		p.nextToken()
		assign.IsAppend = true
		p.consumeAssignedValue(assign, startLine)
	case p.peekTokenIs(token.ASSIGN):
		p.nextToken()
		p.consumeAssignedValue(assign, startLine)
	}
}

func (p *Parser) consumeAssignedValue(assign *ast.DeclarationAssignment, startLine int) {
	if p.peekToken.Line != startLine || p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF) {
		return
	}
	p.nextToken()
	assign.Value = p.parseDeclarationValue()
}

func (p *Parser) parseDeclarationValue() ast.Expression {
	// Check for Array literal `( ... )`
	if p.curTokenIs(token.LPAREN) {
		paren := p.curToken
		p.nextToken() // consume (

		// Track paren and brace depth so nested `$(...)`,
		// `${...}`, `$((...))` inside the array literal don't
		// terminate the scan prematurely. Without this,
		// `x=($(cmd))` fell off the end looking for the outer `)`
		// because the lexer's `$(` + `)` pair consumed the `)` we
		// expected.
		val := "("
		depth := 0
		for !p.curTokenIs(token.EOF) {
			switch {
			case p.curTokenIs(token.RPAREN):
				if depth == 0 {
					goto arrDone
				}
				depth--
			case p.curTokenIs(token.LPAREN),
				p.curTokenIs(token.DOLLAR_LPAREN),
				p.curTokenIs(token.DoubleLparen),
				p.curTokenIs(token.LBRACE),
				p.curTokenIs(token.DollarLbrace):
				depth++
			case p.curTokenIs(token.DoubleRparen):
				if depth > 0 {
					depth--
				}
			case p.curTokenIs(token.RBRACE):
				if depth > 0 {
					depth--
				}
			}
			val += " " + p.curToken.Literal
			p.nextToken()
		}
	arrDone:
		val += " )"
		// Leave curToken on the closing `)` rather than advancing
		// past it. The caller's declaration loop checks
		// `curToken.Line == startLine`; if we step past `)` onto
		// the next line's first token, the loop drops out and
		// parseProgram then re-advances, double-skipping the next
		// statement (e.g. `typeset -g X=(a b)\ntypeset -g Y=…`
		// dropped the second TYPESET).
		return &ast.StringLiteral{Token: paren, Value: val}
	}

	// Normal expression
	return p.parseExpression(LOWEST)
}

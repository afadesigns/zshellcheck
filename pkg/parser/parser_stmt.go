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
	switch p.curToken.Type {
	case token.RETURN:
		return p.parseReturnStatement()
	case token.LET:
		return p.parseLetStatement()
	case token.If:
		stmt := p.parseIfStatement()
		p.consumePipelineTail()
		return stmt
	case token.SHEBANG:
		return p.parseShebangStatement()
	case token.HASH:
		// Skip comments for now
		return nil
	case token.FOR:
		stmt := p.parseForLoopStatement()
		p.consumePipelineTail()
		return stmt
	case token.WHILE:
		stmt := p.parseWhileLoopStatement()
		p.consumePipelineTail()
		return stmt
	case token.SELECT:
		stmt := p.parseSelectStatement()
		p.consumePipelineTail()
		return stmt
	case token.COPROC:
		return p.parseCoprocStatement()
	case token.TYPESET, token.DECLARE:
		return p.parseDeclarationStatement()
	case token.LBRACE:
		tok := p.curToken
		p.nextToken()
		block := p.parseBlockStatement(token.RBRACE)
		block.Token = tok
		// A brace group can head a pipeline or logical chain:
		// `{ cmd1; cmd2 } | sort`, `{ a || b } | awk`. Consume any
		// trailing pipeline / logical continuations opaquely so
		// parseStatement doesn't choke on the leading `|` / `&&`
		// / `||` as an unknown prefix. The AST keeps the block as
		// the statement; detection katas that care about the full
		// pipeline can walk source.
		for p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
			p.nextToken() // onto op
			p.nextToken() // onto RHS head
			_ = p.parseCommandPipeline()
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return block
	case token.LPAREN:
		return p.parseSubshellStatement()
	case token.DoubleLparen:
		cmd := p.parseArithmeticCommand()
		if cmd == nil {
			return nil
		}
		if chained := p.chainLogical(cmd, cmd.Token); chained != nil {
			return chained
		}
		return cmd
	case token.LDBRACKET:
		// `[[ … ]]` is a prefix expression by default. As a statement
		// we need to capture the bracketed expression AND the `&&` /
		// `||` continuations without letting the generic
		// parseExpression loop pick OR/AND up as internal infix
		// operators — that swallows the continuation's right-hand
		// command (e.g. `|| return 0`) into a single expression
		// whose RHS starts at `return`, which has no prefix parse
		// entry and errors out.
		//
		// Call the prefix function directly so the expression stops
		// exactly at `]]`, then route post-`]]` logical chains
		// through chainLogical, which uses parseCommandPipeline for
		// the RHS — the command-aware path that knows how to handle
		// `return`, builtins, simple commands, and so on.
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
	case token.COLON, token.DOT, token.LBRACKET,
		token.GT, token.LT, token.GTGT, token.LTLT, token.GTAMP, token.LTAMP, token.AMPERSAND, token.SLASH:
		return p.parseSimpleCommandStatement()
	case token.BANG:
		// Shell `!` negates the exit status of the following
		// pipeline: `! cmd 2>/dev/null | grep`. Route through the
		// command-pipeline path so redirects and pipes on the
		// right chain correctly. Keep the expression-level prefix
		// behaviour for C-style inputs like `!5` / `!true` by
		// checking peek: IDENT / LPAREN / LBRACKET / LDBRACKET /
		// DoubleLparen / VARIABLE / DollarLbrace / BACKTICK /
		// DOLLAR_LPAREN are command starts; anything else falls
		// back to the expression parser.
		if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.LPAREN) ||
			p.peekTokenIs(token.LBRACKET) || p.peekTokenIs(token.LDBRACKET) ||
			p.peekTokenIs(token.DoubleLparen) || p.peekTokenIs(token.VARIABLE) ||
			p.peekTokenIs(token.DollarLbrace) || p.peekTokenIs(token.BACKTICK) ||
			p.peekTokenIs(token.DOLLAR_LPAREN) {
			return p.parseSimpleCommandStatement()
		}
		return p.parseExpressionOrFunctionDefinition()
	case token.BACKTICK, token.DOLLAR_LPAREN, token.VARIABLE, token.DollarLbrace:
		// A command-producing expression (`cmd`, $(cmd), $name,
		// ${name}) can stand on its own as a statement, but can
		// also head a pipeline or a logical chain:
		// `` `_cmd` | sed … ``, `$(date) && ...`, `$VAR | awk`.
		// Parse the expression via the normal prefix path, then
		// fold any pipeline / AND / OR continuations into an infix
		// tree so the trailing `|` / `&&` / `||` do not leak back
		// into parseStatement's next-iteration dispatch.
		return p.parsePipelineStartingWithExpression()
	case token.CASE:
		stmt := p.parseCaseStatement()
		p.consumePipelineTail()
		return stmt
	case token.IDENT:
		if p.curToken.Literal == "test" {
			return p.parseSimpleCommandStatement()
		}
		if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.STRING) || p.peekTokenIs(token.INT) ||
			p.peekTokenIs(token.MINUS) || p.peekTokenIs(token.DOT) || p.peekTokenIs(token.VARIABLE) ||
			p.peekTokenIs(token.DOLLAR) || p.peekTokenIs(token.DollarLbrace) ||
			p.peekTokenIs(token.DOLLAR_LPAREN) || p.peekTokenIs(token.SLASH) ||
			p.peekTokenIs(token.TILDE) || p.peekTokenIs(token.ASTERISK) ||
			p.peekTokenIs(token.BANG) || p.peekTokenIs(token.LBRACE) ||
			// Zero-arg commands followed by a pipe / logical chain
			// must route through parseSimpleCommandStatement so the
			// pipeline / AND / OR chain is parsed at the command
			// layer. Without this `cmd1 |\n cmd2` left `cmd1` as a
			// bare Identifier expression, and the block loop then
			// tried to start a new statement at `|`.
			p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.AND) || p.peekTokenIs(token.OR) {
			return p.parseSimpleCommandStatement()
		}
		return p.parseExpressionOrFunctionDefinition()
	default:
		return p.parseExpressionOrFunctionDefinition()
	}
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

	var left ast.Expression
	switch p.curToken.Type {
	case token.WHILE:
		left = p.parseWhileLoopStatement()
	case token.LPAREN:
		// Subshell group: `( cmd1; cmd2 )` appears as the RHS of
		// logical chains like `[[ … ]] && ( … )`. Parse the group
		// as a grouped expression so parseSingleCommand doesn't
		// treat `(` as a command name.
		left = p.parseGroupedExpression()
	case token.LDBRACKET:
		// `[[ … ]]` condition as a pipeline term (RHS of `&&`/`||`
		// or head of a pipe). Call the prefix directly so the
		// caller doesn't try to use it as a simple-command name.
		left = p.parseDoubleBracketExpression()
	case token.DoubleLparen:
		left = p.parseArithmeticCommand()
	default:
		left = p.parseSingleCommand()
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

func (p *Parser) parseCommandWord() ast.Expression {
	firstToken := p.curToken
	parts := []ast.Expression{}

	// Helper to determine if we should parse as expression
	isExpression := func(t token.Type) bool {
		// Treat these as literals in command args, even if they have prefix fns
		if t == token.ASTERISK || t == token.QUESTION || t == token.PLUS ||
			t == token.MINUS || t == token.CARET || t == token.TILDE || t == token.DOT ||
			t == token.GT || t == token.LT || t == token.AMPERSAND || t == token.LBRACKET ||
			t == token.COMMA || t == token.COLON || t == token.GTGT || t == token.LTLT ||
			t == token.GTAMP || t == token.LTAMP ||
			// `=` is an assignment operator in expression context but a
			// literal in command arguments (e.g. `alias -- -='cd -'`,
			// or `env FOO=bar cmd`). Treat it as a literal word part
			// when it appears mid-command. The declaration parser has
			// its own dedicated handling for the IDENT=VALUE form.
			t == token.ASSIGN || t == token.PLUSEQ {
			return false
		}
		return p.prefixParseFns[t] != nil
	}

	// Track LBRACE depth opened MID-WORD so a matching `}` isn't
	// mistaken for a delimiter. Zsh git refspecs like
	// `@{upstream}` appear bare on command lines; without this the
	// RBRACE closed the arg at `{upstream` and the outer `$(` lost
	// its closing `)`. When the word STARTS with `{` we leave
	// braceDepth at zero so brace expansions `{1..10}` still
	// terminate at `}` (tests like ZC1083 expect `{1..10}$var`
	// to parse as two concatenated words: `{1..10}` then `$var`).
	braceDepth := 0

	// Parse the first part
	if !isExpression(p.curToken.Type) {
		parts = append(parts, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})
	} else {
		parts = append(parts, p.parseExpression(CALL))
	}

	// Continue parsing while the next token is adjacent (no preceding space)
	for !p.peekToken.HasPrecedingSpace && p.peekOnSameLogicalLine() {
		if braceDepth == 0 && p.isCommandDelimiter(p.peekToken) {
			break
		}
		if braceDepth > 0 && p.peekTokenIs(token.EOF) {
			break
		}
		p.nextToken()
		switch p.curToken.Type {
		case token.LBRACE:
			braceDepth++
		case token.RBRACE:
			if braceDepth > 0 {
				braceDepth--
			}
		}
		if !isExpression(p.curToken.Type) {
			// Treat as literal string part
			parts = append(parts, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})
		} else {
			parts = append(parts, p.parseExpression(CALL))
		}
	}

	if len(parts) == 1 {
		return parts[0]
	}

	return &ast.ConcatenatedExpression{
		Token: firstToken,
		Parts: parts,
	}
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
	stmt.Condition = p.parseBlockStatement(token.THEN)

	if !p.curTokenIs(token.THEN) {
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

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
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
				curIsTerm = true
				break
			}
		}
		if curIsTerm {
			break
		}

		p.nextToken()
	}
	return block
}

func (p *Parser) parseSubshellStatement() ast.Statement {
	subshellToken := p.curToken
	p.nextToken()
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
		clause := &ast.CaseClause{Token: p.curToken}
		if p.curTokenIs(token.LPAREN) {
			p.nextToken()
		}
		for {
			pat := p.parseCommandWord()
			clause.Patterns = append(clause.Patterns, pat)
			if p.peekTokenIs(token.PIPE) {
				p.nextToken()
				p.nextToken()
			} else {
				break
			}
		}
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		p.nextToken()
		clause.Body = p.parseBlockStatement(token.DSEMI, token.ESAC)
		stmt.Clauses = append(stmt.Clauses, clause)
		if p.curTokenIs(token.DSEMI) {
			p.nextToken()
		}
	}
	// Consume the closing ESAC so a caller parsing a nested case
	// (`case x in a) case y in …;; esac ;; esac`) doesn't see the
	// inner `esac` as the outer's terminator. Without this advance
	// parseBlockStatement's terminator check would fire on the
	// inner ESAC and collapse the outer case body.
	if p.curTokenIs(token.ESAC) {
		// Leave at-ESAC if we're the outermost call (no advance
		// necessary — parseBlockStatement will re-check peek);
		// but always leave curToken pointing at ESAC's successor
		// when there is one so the outer DSEMI check works.
		// Peek past ESAC only if there's a trailing `;;` pattern
		// terminator or ;/newline.
		if p.peekTokenIs(token.DSEMI) || p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	return stmt
}

func (p *Parser) parseShebangStatement() *ast.Shebang {
	return &ast.Shebang{Token: p.curToken, Path: p.curToken.Literal}
}

func (p *Parser) parseForLoopStatement() *ast.ForLoopStatement {
	stmt := &ast.ForLoopStatement{Token: p.curToken}

	if p.peekTokenIs(token.DoubleLparen) {
		// Arithmetic for loop: for (( init; cond; post ))
		p.nextToken() // consume ((

		// Init (optional)
		if !p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
			if p.prefixParseFns[p.curToken.Type] != nil {
				stmt.Init = p.parseExpression(LOWEST)
			} else {
				p.noPrefixParseFnError(p.curToken.Type)
				return nil
			}
		}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		// Condition (optional)
		if !p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
			if p.prefixParseFns[p.curToken.Type] != nil {
				stmt.Condition = p.parseExpression(LOWEST)
			} else {
				p.noPrefixParseFnError(p.curToken.Type)
				return nil
			}
		}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		// Post (optional)
		if !p.peekTokenIs(token.DoubleRparen) {
			p.nextToken()
			if p.prefixParseFns[p.curToken.Type] != nil {
				stmt.Post = p.parseExpression(LOWEST)
			} else {
				p.noPrefixParseFnError(p.curToken.Type)
				return nil
			}
		}

		if !p.expectPeek(token.DoubleRparen) {
			return nil
		}

		// Optional semicolon before DO
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		if !p.expectPeek(token.DO) {
			return nil
		}
		p.nextToken() // consume DO
		stmt.Body = p.parseBlockStatement(token.DONE)
		return stmt
	}

	// For-each loop: for name [in words]; do
	// Zsh accepts numeric positional names like `for 1 in "$@"; do`
	// (shorthand to iterate over positionals). Allow INT as the
	// binding name alongside IDENT.
	if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.INT) {
		p.nextToken()
	} else {
		p.peekError(token.IDENT)
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Zsh multi-variable for loop: `for k v in …` / `for a b c in …`
	// pairs each element of the item list against the declared
	// variables in turn. The AST currently only models a single
	// Name, so skip extra names forward until we hit IN / LPAREN /
	// SEMICOLON / DO. Detection katas that need the full name list
	// can read source directly.
	for p.peekTokenIs(token.IDENT) {
		p.nextToken()
	}

	// Zsh short form: `for NAME ( items ) body`. The item list is
	// wrapped in parentheses and the body is a single command (or
	// block) with no `do`/`done`. Real-world example in prezto
	// init.zsh: `for zmodule ("$zmodules[@]") zmodload "zsh/…"`.
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // consume (
		stmt.Items = []ast.Expression{}
		for !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.EOF) {
			p.nextToken()
			if p.curTokenIs(token.RPAREN) {
				break
			}
			arg := p.parseCommandWord()
			stmt.Items = append(stmt.Items, arg)
		}
		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		// Body form varies. Some Zsh code uses the pure short form
		// (`for x (items) body`); other code mixes short-form items
		// with a classic `do/done` body (`for x (items); do … done`).
		// A leading `;` and `do` indicates the latter.
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		if p.peekTokenIs(token.DO) {
			p.nextToken() // onto DO
			p.nextToken() // into body
			stmt.Body = p.parseBlockStatement(token.DONE)
			return stmt
		}

		// Body is a single statement on the same line (usually a
		// command) or a braced block. Wrap non-block statements in
		// a BlockStatement so the Body field stays homogeneous.
		p.nextToken()
		body := p.parseStatement()
		if block, ok := body.(*ast.BlockStatement); ok {
			stmt.Body = block
		} else if body != nil {
			stmt.Body = &ast.BlockStatement{
				Token:      stmt.Token,
				Statements: []ast.Statement{body},
			}
		}
		return stmt
	}

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

	// Zsh short body form: `for x in items; { body }` replaces
	// `do … done` with a brace block. Accept LBRACE here alongside
	// the classic DO keyword.
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken() // onto {
		p.nextToken() // into body
		stmt.Body = p.parseBlockStatement(token.RBRACE)
		return stmt
	}

	if !p.expectPeek(token.DO) {
		return nil
	}

	p.nextToken()
	stmt.Body = p.parseBlockStatement(token.DONE)
	return stmt
}

func (p *Parser) parseWhileLoopStatement() *ast.WhileLoopStatement {
	stmt := &ast.WhileLoopStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseBlockStatement(token.DO)
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

	// Consume flags and assignments until the statement's line ends or a
	// terminator fires. advanceOrStop is the key helper: it only moves
	// past the current token when the next token is still part of this
	// declaration (same line, not a terminator). When the next token is
	// on a new line the declaration ends with curToken on its last real
	// token so the outer block's unconditional nextToken() advances to
	// the following statement's first token — without this guard, a
	// declaration immediately followed by an `if` (or any statement) on
	// the next line caused the parser to swallow the statement's leading
	// token. Reported against oh-my-zsh / zsh-autosuggestions bodies.
	advanceOrStop := func() bool {
		if p.peekTokenIs(token.SEMICOLON) || p.peekTokenIs(token.EOF) {
			return false
		}
		if p.peekToken.Line != startLine {
			return false
		}
		p.nextToken()
		return true
	}

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) && p.curToken.Line == startLine {
		// Flags (e.g. -g, -A, -r, --).
		if p.curTokenIs(token.MINUS) || (p.curTokenIs(token.IDENT) && len(p.curToken.Literal) > 0 && p.curToken.Literal[0] == '-') {
			stmt.Flags = append(stmt.Flags, p.curToken.Literal)
			if !advanceOrStop() {
				break
			}
			continue
		}

		// Identifier (optionally followed by = or += value).
		if p.curTokenIs(token.IDENT) {
			assign := &ast.DeclarationAssignment{
				Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			}
			// Non-trivial composite names: the parser should treat
			// `X_{A,B}_Y`, `prefix"${1:-}"_suffix`, and
			// `${arr[$i]}=value` as one logical name. Consume any
			// sequence of adjacent (no preceding space) tokens that
			// can participate in a shell-word: LBRACE expansions,
			// STRING literals, DollarLbrace expansions, VARIABLE
			// references, and plain IDENT suffixes.
			for !p.peekToken.HasPrecedingSpace && p.peekToken.Line == startLine {
				switch {
				case p.peekTokenIs(token.LBRACE):
					p.nextToken() // onto {
					depth := 1
					for depth > 0 && !p.peekTokenIs(token.EOF) {
						p.nextToken()
						switch {
						case p.curTokenIs(token.LBRACE):
							depth++
						case p.curTokenIs(token.RBRACE):
							depth--
						}
					}
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
				case p.peekTokenIs(token.STRING):
					p.nextToken()
				case p.peekTokenIs(token.IDENT):
					p.nextToken()
				case p.peekTokenIs(token.VARIABLE):
					p.nextToken()
				default:
					goto nameDone
				}
			}
		nameDone:
			// Peek the =/+= before consuming the name so we can decide
			// whether to stay on the name token (bare declaration) or
			// move onto the operator (value follows). An empty RHS
			// (`typeset -g VAR=` at end-of-line) is valid Zsh and sets
			// the variable to the empty string — do NOT try to parse
			// the next-line token as the value, that over-consumes
			// into the following statement exactly like the pre-
			// declaration fix handled for bare declarations.
			if p.peekTokenIs(token.PLUSEQ) {
				p.nextToken() // onto +=
				assign.IsAppend = true
				if p.peekToken.Line == startLine && !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) {
					p.nextToken() // onto value token
					assign.Value = p.parseDeclarationValue()
				}
			} else if p.peekTokenIs(token.ASSIGN) {
				p.nextToken() // onto =
				if p.peekToken.Line == startLine && !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) {
					p.nextToken() // onto value token
					assign.Value = p.parseDeclarationValue()
				}
			}
			stmt.Assignments = append(stmt.Assignments, assign)

			if !advanceOrStop() {
				break
			}
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
	cmd.Expression = p.parseExpression(LOWEST)
	p.inArithmetic = prevInArithmetic

	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	return cmd
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
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
		return &ast.StringLiteral{Token: paren, Value: val}
	}

	// Normal expression
	return p.parseExpression(LOWEST)
}

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
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.If:
		return p.parseIfStatement()
	case token.SHEBANG:
		return p.parseShebangStatement()
	case token.HASH:
		// Skip comments for now
		return nil
	case token.FOR:
		return p.parseForLoopStatement()
	case token.WHILE:
		return p.parseWhileLoopStatement()
	case token.SELECT:
		return p.parseSelectStatement()
	case token.COPROC:
		return p.parseCoprocStatement()
	case token.TYPESET, token.DECLARE:
		return p.parseDeclarationStatement()
	case token.LBRACE:
		tok := p.curToken
		p.nextToken()
		block := p.parseBlockStatement(token.RBRACE)
		block.Token = tok
		return block
	case token.LPAREN:
		return p.parseSubshellStatement()
	case token.DoubleLparen:
		cmd := p.parseArithmeticCommand()
		if cmd == nil {
			return nil
		}
		return cmd
	case token.COLON, token.DOT, token.LBRACKET:
		return p.parseSimpleCommandStatement()
	case token.CASE:
		return p.parseCaseStatement()
	case token.IDENT:
		if p.curToken.Literal == "test" {
			return p.parseSimpleCommandStatement()
		}
		if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.STRING) || p.peekTokenIs(token.INT) ||
			p.peekTokenIs(token.MINUS) || p.peekTokenIs(token.DOT) || p.peekTokenIs(token.VARIABLE) ||
			p.peekTokenIs(token.DOLLAR) || p.peekTokenIs(token.DollarLbrace) ||
			p.peekTokenIs(token.DOLLAR_LPAREN) || p.peekTokenIs(token.SLASH) ||
			p.peekTokenIs(token.TILDE) || p.peekTokenIs(token.ASTERISK) ||
			p.peekTokenIs(token.BANG) || p.peekTokenIs(token.LBRACE) {
			return p.parseSimpleCommandStatement()
		}
		return p.parseExpressionOrFunctionDefinition()
	default:
		return p.parseExpressionOrFunctionDefinition()
	}
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

	left := p.parseSingleCommand()

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

	if p.peekTokenIs(token.PIPE) {
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
	for !p.isCommandDelimiter(p.peekToken) && p.peekToken.Line == p.curToken.Line {
		p.nextToken()
		arg := p.parseCommandWord()
		cmd.Arguments = append(cmd.Arguments, arg)
	}

	return cmd
}

func (p *Parser) parseCommandWord() ast.Expression {
	firstToken := p.curToken
	parts := []ast.Expression{}

	// Parse the first part
	if p.prefixParseFns[p.curToken.Type] == nil {
		parts = append(parts, &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal})
	} else {
		parts = append(parts, p.parseExpression(CALL))
	}

	// Continue parsing while the next token is adjacent (no preceding space)
	for !p.peekToken.HasPrecedingSpace && !p.isCommandDelimiter(p.peekToken) &&
		p.peekToken.Line == p.curToken.Line {

		p.nextToken()

		if p.prefixParseFns[p.curToken.Type] == nil {
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
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
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
	stmt.Consequence = p.parseBlockStatement(token.ELSE, token.Fi)

	if p.curTokenIs(token.ELSE) {
		p.nextToken() // consume "else"
		stmt.Alternative = p.parseBlockStatement(token.Fi)
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
	return &ast.Subshell{Token: subshellToken, Block: block}
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
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Init = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.IN) {
		p.nextToken()
		stmt.Items = []ast.Expression{}
		for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.DO) && !p.peekTokenIs(token.EOF) &&
			p.peekToken.Line == p.curToken.Line {
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
			p.peekToken.Line == p.curToken.Line {
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
	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) && p.curToken.Line == stmt.Token.Line {
		// Check for flags
		if p.curTokenIs(token.MINUS) || (p.curTokenIs(token.IDENT) && len(p.curToken.Literal) > 0 && p.curToken.Literal[0] == '-') {
			// It's a flag (e.g., -A, -r, --)
			stmt.Flags = append(stmt.Flags, p.curToken.Literal)
			p.nextToken()
			continue
		} else if p.curTokenIs(token.MINUS) {
			// Standalone minus?
			stmt.Flags = append(stmt.Flags, "-")
			p.nextToken()
			continue
		}

		// Expect Identifier
		if p.curTokenIs(token.IDENT) {
			assign := &ast.DeclarationAssignment{
				Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			}
			p.nextToken() // consume Name

			// Check for = or +=
			if p.curTokenIs(token.PLUSEQ) {
				assign.IsAppend = true
				p.nextToken() // consume +=
				assign.Value = p.parseDeclarationValue()
			} else if p.curTokenIs(token.ASSIGN) {
				p.nextToken() // consume =
				assign.Value = p.parseDeclarationValue()
			}

			stmt.Assignments = append(stmt.Assignments, assign)
		} else {
			// Unexpected token in declaration
			p.nextToken()
		}
	}
	
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseArithmeticCommand() *ast.ArithmeticCommand {
	cmd := &ast.ArithmeticCommand{Token: p.curToken}
	p.nextToken()
	cmd.Expression = p.parseExpression(LOWEST)

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
		
		val := "("
		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			val += " " + p.curToken.Literal // Very rough
			p.nextToken()
		}
		val += " )"
		if p.curTokenIs(token.RPAREN) {
			p.nextToken()
		}
		return &ast.StringLiteral{Token: paren, Value: val}
	}
	
	// Normal expression
	return p.parseExpression(LOWEST)
}

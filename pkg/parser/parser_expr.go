package parser

import (
	"fmt"
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		if !p.inArithmetic && p.peekTokenIs(token.LBRACKET) && p.peekToken.HasPrecedingSpace {
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
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	tok := p.curToken
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return &ast.GroupedExpression{Token: tok, Expression: exp}
}

func (p *Parser) parseArrayAccess() ast.Expression {
	exp := &ast.ArrayAccess{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.Left = p.parseIdentifier()

	// check for optional index
	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken() // consume [
		p.nextToken() // move to start of index expression
		exp.Index = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return exp
}

func (p *Parser) parseInvalidArrayAccessPrefix() ast.Expression {
	dollarToken := p.curToken
	if p.peekTokenIs(token.HASH) || p.peekTokenIs(token.INT) || p.peekTokenIs(token.ASTERISK) || p.peekTokenIs(token.BANG) || p.peekTokenIs(token.MINUS) {
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		return &ast.PrefixExpression{Token: dollarToken, Operator: "$", Right: ident}
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
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
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
	exp.Command = p.parseCommandList()
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
	exp.Command = p.parseCommandList()
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
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

	prevInArithmetic := p.inArithmetic
	p.inArithmetic = true
	exp.Index = p.parseExpression(LOWEST)
	p.inArithmetic = prevInArithmetic

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseProcessSubstitution() ast.Expression {
	exp := &ast.ProcessSubstitution{Token: p.curToken}
	p.nextToken()

	// Process substitution contains a command list
	exp.Command = p.parseCommandList()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

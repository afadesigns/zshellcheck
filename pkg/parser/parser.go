package parser

import (
	"fmt"
	"strconv"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

const (
	_ int = iota
	LOWEST
	LOGICAL     // || or &&
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	POSTFIX     // i++
)

var precedences = map[token.Type]int{
	token.AND:      LOGICAL,
	token.OR:       LOGICAL,
	token.EQ:       EQUALS,
	token.NotEq:    EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.PIPE:     CALL,
	token.ASSIGN:   EQUALS,
	token.INC:      POSTFIX,
	token.DEC:      POSTFIX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseSingleCommand)
	p.registerPrefix(token.LDBRACKET, p.parseDoubleBracketExpression)
	p.registerPrefix(token.DollarLbrace, p.parseArrayAccess)
	p.registerPrefix(token.DOLLAR, p.parseInvalidArrayAccessPrefix)
	p.registerPrefix(token.VARIABLE, p.parseIdentifier)
	p.registerPrefix(token.DOLLAR_LPAREN, p.parseDollarParenExpression)
	p.registerPrefix(token.DoubleLparen, p.parseDoubleParenExpression)
	p.registerPrefix(token.BACKTICK, p.parseCommandSubstitution)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PIPE, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

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
	case token.LBRACE:
		return p.parseBlockStatement(token.RBRACE)
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
			p.peekTokenIs(token.DOLLAR) || p.peekTokenIs(token.DollarLbrace) || p.peekTokenIs(token.DOLLAR_LPAREN) {
			return p.parseSimpleCommandStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
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
		// Command pipeline follows (parseCommandPipeline handles simple commands and pipes)
		// But we need to move to start of next command. 
		// nextToken was called. curToken is AND.
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
	left := p.parseSingleCommand()

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
	cmd.Arguments = []ast.Expression{}

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.PIPE) &&
		!p.peekTokenIs(token.AND) && !p.peekTokenIs(token.OR) &&
		!p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.LPAREN) && !p.peekTokenIs(token.RBRACE) &&
		!p.peekTokenIs(token.HASH) &&
		!p.peekTokenIs(token.THEN) && !p.peekTokenIs(token.ELSE) && !p.peekTokenIs(token.ELIF) && !p.peekTokenIs(token.Fi) &&
		!p.peekTokenIs(token.DO) && !p.peekTokenIs(token.DONE) &&
		!p.peekTokenIs(token.ESAC) && !p.peekTokenIs(token.DSEMI) &&
		p.peekToken.Line == p.curToken.Line {
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
	for !p.peekToken.HasPrecedingSpace && !p.peekTokenIs(token.EOF) &&
		!p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.PIPE) &&
		!p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.HASH) &&
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
		// If we stopped for EOF or other reason
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

	isTerminator := func(tok token.Token) bool {
		for _, t := range terminators {
			if tok.Type == t {
				return true
			}
		}
		return false
	}

	loopCount := 0
	for !isTerminator(p.curToken) && !p.curTokenIs(token.EOF) {
		loopCount++
		if loopCount > 10000 {
			break
		}
		s := p.parseStatement()
		if s != nil {
			block.Statements = append(block.Statements, s)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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

func (p *Parser) parseBracketExpression() ast.Expression {
	bracketToken := p.curToken
	p.nextToken()
	expressions := []ast.Expression{}
	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		exp := p.parseExpression(LOWEST)
		if exp != nil {
			expressions = append(expressions, exp)
		}
		if !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
			p.nextToken()
		}
	}
	if !p.curTokenIs(token.RBRACKET) {
		p.peekError(token.RBRACKET)
		return nil
	}
	return &ast.BracketExpression{Token: bracketToken, Expressions: expressions}
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
	return &ast.DoubleBracketExpression{Token: bracketToken, Expressions: expressions}
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
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseArrayAccess() ast.Expression {
	exp := &ast.ArrayAccess{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	exp.Left = p.parseIdentifier()
	if !p.expectPeek(token.LBRACKET) {
		return nil
	}
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return exp
}

func (p *Parser) parseInvalidArrayAccessPrefix() ast.Expression {
	dollarToken := p.curToken
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

func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	stmt := &ast.CaseStatement{Token: p.curToken}

	p.nextToken() // consume CASE
	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.IN) {
		return nil
	}
	p.nextToken() // consume IN

	for !p.curTokenIs(token.ESAC) && !p.curTokenIs(token.EOF) {
		// Allow optional newlines before patterns
		for p.curTokenIs(token.SEMICOLON) { // Lexer might emit semicolons for newlines if I configured it so?
			// Or parser skips them? parseBlockStatement skips them.
			// But here we are in Case body.
			p.nextToken()
		}
		if p.curTokenIs(token.ESAC) {
			break
		}

		clause := &ast.CaseClause{Token: p.curToken}
		
		// Parse patterns
		// Optional leading (
		if p.curTokenIs(token.LPAREN) {
			p.nextToken()
		}
		
		// Patterns separated by |
		for {
			// Parse pattern word. 
			// Should be parseCommandWord? Or parseExpression?
			// Patterns are words. Glob patterns.
			// parseCommandWord matches "word" logic.
			pat := p.parseCommandWord()
			clause.Patterns = append(clause.Patterns, pat)
			
			if p.peekTokenIs(token.PIPE) {
				p.nextToken() // consume arg
				p.nextToken() // consume |
			} else {
				break
			}
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		p.nextToken() // consume )

		// Parse body
		clause.Body = p.parseBlockStatement(token.DSEMI, token.ESAC)
		
		stmt.Clauses = append(stmt.Clauses, clause)

		if p.curTokenIs(token.DSEMI) {
			p.nextToken() // consume ;;
		}
	}

	return stmt
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Zsh functions can be anonymous `function() {}` or named `function my_func() {}`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken() // consume function name
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	p.nextToken() // consume "{"
	lit.Body = p.parseBlockStatement(token.RBRACE)

	return lit
}

func (p *Parser) parseCommandSubstitution() ast.Expression {
	exp := &ast.CommandSubstitution{Token: p.curToken}
	p.nextToken()
	exp.Command = p.parseExpression(LOWEST)
	if !p.expectPeek(token.BACKTICK) {
		return nil
	}
	return exp
}

func (p *Parser) parseDollarParenExpression() ast.Expression {
	exp := &ast.DollarParenExpression{Token: p.curToken}
	p.nextToken()
	exp.Command = p.parseExpression(LOWEST)
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
	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
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

func (p *Parser) parseShebangStatement() *ast.Shebang {
	return &ast.Shebang{Token: p.curToken, Path: p.curToken.Literal}
}

func (p *Parser) parseForLoopStatement() *ast.ForLoopStatement {
	stmt := &ast.ForLoopStatement{Token: p.curToken}

	if p.peekTokenIs(token.DoubleLparen) {
		// Arithmetic for loop: for (( init; cond; post ))
		p.nextToken() // consume ((

		if !p.peekTokenIs(token.SEMICOLON) {
			stmt.Init = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		// Consumed semicolon.
		// Now Cond.
		if !p.peekTokenIs(token.SEMICOLON) {
			stmt.Condition = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		// Consumed semicolon.
		// Now Post.
		if !p.peekTokenIs(token.DoubleRparen) {
			stmt.Post = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.DoubleRparen) {
			return nil
		}
	} else {
		// For-each loop: for name [in words]; do
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		// Use Init field to store the variable name (as Identifier)
		stmt.Init = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if p.peekTokenIs(token.IN) {
			p.nextToken() // consume "in"
			
			stmt.Items = []ast.Expression{}
			for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.DO) && !p.peekTokenIs(token.EOF) &&
				p.peekToken.Line == p.curToken.Line {
				p.nextToken()
				arg := p.parseCommandWord()
				stmt.Items = append(stmt.Items, arg)
			}
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

	p.nextToken() // consume DO
	stmt.Body = p.parseBlockStatement(token.DONE)

	return stmt
}

func (p *Parser) parseDoubleParenExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.DoubleRparen) {
		return nil
	}
	return exp
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	if t == token.RPAREN && p.peekTokenIs(token.DoubleRparen) {
		p.curToken = p.peekToken
		p.curToken.Type = token.RPAREN
		p.curToken.Literal = ")"

		p.peekToken.Type = token.RPAREN
		p.peekToken.Literal = ")"
		p.peekToken.Column++
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

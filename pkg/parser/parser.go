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
	p.registerPrefix(token.LBRACKET, p.parseBracketExpression)
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
	case token.FOR:
		return p.parseForLoopStatement()
	case token.IDENT:
		if p.curToken.Literal == "test" {
			return p.parseSimpleCommandStatement()
		}
		if p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.STRING) || p.peekTokenIs(token.INT) ||
			p.peekTokenIs(token.MINUS) || p.peekTokenIs(token.DOT) {
			return p.parseSimpleCommandStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseSimpleCommandStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseCommandPipeline()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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
		!p.peekTokenIs(token.RPAREN) && p.peekToken.Line == p.curToken.Line {
		p.nextToken()
		arg := p.parseExpression(CALL)
		cmd.Arguments = append(cmd.Arguments, arg)
	}

	return cmd
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
	stmt.Condition = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if !p.expectPeek(token.THEN) {
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

	for !isTerminator(p.curToken) && !p.curTokenIs(token.EOF) {
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

	if !p.expectPeek(token.DoubleLparen) { // curToken is 'for', peekToken should be '(('. After expectPeek, curToken is '(('
		return nil
	}

	// Parse Init
	p.nextToken() // curToken is now 'i' (start of init expression)
	stmt.Init = p.parseExpression(LOWEST) // Parses 'i=0'. curToken is '0'. peekToken is ';'

	if !p.expectPeek(token.SEMICOLON) { // expects ';'
		return nil
	}
	// After expectPeek, curToken is ';'. peekToken is 'i' (start of condition expression)

	// Parse Condition
	p.nextToken() // curToken becomes 'i' (start of condition expression)
	stmt.Condition = p.parseExpression(LOWEST) // Parses 'i<10'. curToken is '10'. peekToken is ';'

	if !p.expectPeek(token.SEMICOLON) { // expects ';'
		return nil
	}
	// After expectPeek, curToken is ';'. peekToken is 'i' (start of post expression)

	// Parse Post
	p.nextToken() // curToken becomes 'i' (start of post expression)
	stmt.Post = p.parseExpression(LOWEST) // Parses 'i++'. curToken is '++'. peekToken is '))'

	if !p.expectPeek(token.DoubleRparen) { // expects '))'
		return nil
	}
	// After expectPeek, curToken is '))'. peekToken is ';' (before 'do')

	// Now we are at the token after '))', which should be the SEMICOLON before 'do'
	if !p.curTokenIs(token.SEMICOLON) { // current token should be the SEMICOLON
		return nil
	}
	p.nextToken() // Consume the SEMICOLON. curToken is ';'. peekToken is 'do'

	if !p.expectPeek(token.IDENT) || p.curToken.Literal != "do" { // expects 'do'
		return nil
	}
	// After expectPeek, curToken is 'do'. peekToken is 'echo'

	p.nextToken() // Consume "do". curToken is 'echo'.
	stmt.Body = p.parseBlockStatement(token.DONE)

	return stmt
}

func (p *Parser) parseDoubleParenExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.RPAREN) {
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

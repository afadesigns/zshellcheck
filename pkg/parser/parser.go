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
	POSTFIX     // X++
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.Type]int{
	token.EQ:            EQUALS,
	token.NotEq:         EQUALS,
	token.LT:            LESSGREATER,
	token.GT:            LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.ASTERISK:      PRODUCT,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.DollarLbrace:  INDEX,
	token.DOLLAR_LPAREN: CALL,
	token.DoubleLparen:  CALL,
	token.ASSIGN:        EQUALS,
	token.INC:           POSTFIX,
	token.DEC:           POSTFIX,
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

	inBackticks int
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
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)
	p.registerInfix(token.GTGT, p.parseRedirection)
	p.registerInfix(token.LTLT, p.parseRedirection)
	p.registerInfix(token.GTAMP, p.parseRedirection)
	p.registerInfix(token.LTAMP, p.parseRedirection)

	p.nextToken() // Initialize curToken
	p.nextToken() // Initialize peekToken

	return p
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
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
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
		tok := p.curToken
		p.nextToken()
		block := p.parseBlockStatement(token.RBRACE)
		block.Token = tok
		return block
	case token.LPAREN:
		return p.parseSubshellStatement()
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
			p.peekTokenIs(token.DOLLAR) || p.peekTokenIs(token.DollarLbrace) || p.peekTokenIs(token.DOLLAR_LPAREN) ||
			p.peekTokenIs(token.SLASH) || p.peekTokenIs(token.TILDE) || p.peekTokenIs(token.ASTERISK) || p.peekTokenIs(token.BANG) ||
			p.peekTokenIs(token.LPAREN) {
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

func (p *Parser) isCommandDelimiter(t token.Token) bool {
	if t.Type == token.BACKTICK && p.inBackticks > 0 {
		return true // Terminating backtick is a delimiter
	}
	// If NOT in backticks, start backtick is NOT a delimiter (it starts a substitution)
	if t.Type == token.BACKTICK && p.inBackticks == 0 {
		return false
	}
	
	return t.Type == token.EOF || t.Type == token.SEMICOLON || t.Type == token.PIPE ||
		t.Type == token.AND || t.Type == token.OR ||
		t.Type == token.RPAREN || t.Type == token.RBRACE ||
		t.Type == token.HASH ||
		t.Type == token.THEN || t.Type == token.ELSE || t.Type == token.ELIF || t.Type == token.Fi ||
		t.Type == token.DO || t.Type == token.DONE ||
		t.Type == token.ESAC || t.Type == token.DSEMI ||
		t.Type == token.GT || t.Type == token.LT ||
		t.Type == token.GTGT || t.Type == token.LTLT ||
		t.Type == token.GTAMP || t.Type == token.LTAMP
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
	tok := p.curToken
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return &ast.GroupedExpression{Token: tok, Exp: exp}
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

func (p *Parser) parseSubshellStatement() ast.Statement {
	subshellToken := p.curToken
	p.nextToken()
	block := p.parseBlockStatement(token.RPAREN)
	if !p.curTokenIs(token.RPAREN) {
		p.peekError(token.RPAREN)
		return nil
	}
	p.nextToken()
	block.Token = subshellToken
	return block
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()
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
		cmd := p.parseExpression(LOWEST)
		
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

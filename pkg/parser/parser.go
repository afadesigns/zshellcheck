package parser

import (
	"fmt"

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
	token.EQTILDE:       EQUALS,
	token.EQ_NUM:        EQUALS,
	token.NE_NUM:        EQUALS,
	token.LT_NUM:        LESSGREATER,
	token.LE_NUM:        LESSGREATER,
	token.GT_NUM:        LESSGREATER,
	token.GE_NUM:        LESSGREATER,
	token.INC:           POSTFIX,
	token.DEC:           POSTFIX,
	token.PIPE:          LOWEST + 1,
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

	inBackticks  int
	inArithmetic bool
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.CARET, p.parsePrefixExpression)
	p.registerPrefix(token.ASTERISK, p.parsePrefixExpression)
	p.registerPrefix(token.QUESTION, p.parsePrefixExpression)
	p.registerPrefix(token.TILDE, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACE, p.parseStringLiteral)
	p.registerPrefix(token.RBRACE, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseSingleCommand)
	p.registerPrefix(token.LDBRACKET, p.parseDoubleBracketExpression)
	p.registerPrefix(token.DollarLbrace, p.parseArrayAccess)
	p.registerPrefix(token.DOLLAR, p.parseInvalidArrayAccessPrefix)
	p.registerPrefix(token.VARIABLE, p.parseIdentifier)
	p.registerPrefix(token.DOLLAR_LPAREN, p.parseDollarParenExpression)
	p.registerPrefix(token.DoubleLparen, p.parseDoubleParenExpression)
	p.registerPrefix(token.BACKTICK, p.parseCommandSubstitution)
	p.registerPrefix(token.LT_LPAREN, p.parseProcessSubstitution)
	p.registerPrefix(token.GT_LPAREN, p.parseProcessSubstitution)
	p.registerPrefix(token.EQ_LPAREN, p.parseProcessSubstitution)

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
	p.registerInfix(token.EQTILDE, p.parseInfixExpression)
	p.registerInfix(token.EQ_NUM, p.parseInfixExpression)
	p.registerInfix(token.NE_NUM, p.parseInfixExpression)
	p.registerInfix(token.LT_NUM, p.parseInfixExpression)
	p.registerInfix(token.LE_NUM, p.parseInfixExpression)
	p.registerInfix(token.GT_NUM, p.parseInfixExpression)
	p.registerInfix(token.GE_NUM, p.parseInfixExpression)
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
		t.Type == token.ESAC || t.Type == token.DSEMI
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

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
	LOGICAL     // && or ||
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
	token.AND:           LOGICAL,
	token.OR:            LOGICAL,
	token.EQ:            EQUALS,
	token.NotEq:         EQUALS,
	token.LT:            LESSGREATER,
	token.GT:            LESSGREATER,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.SLASH:         PRODUCT,
	token.ASTERISK:      PRODUCT,
	token.PERCENT:       PRODUCT,
	token.LPAREN:        CALL,
	token.LBRACKET:      INDEX,
	token.DollarLbrace:  INDEX,
	token.DOLLAR_LPAREN: CALL,
	token.DoubleLparen:  CALL,
	token.ASSIGN:        EQUALS,
	token.PLUSEQ:        EQUALS,
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
	// Zsh arithmetic ternary `cond ? a : b`. QUESTION uses LOGICAL
	// precedence; COLON is consumed inside parseInfixExpression's
	// right-hand parse so it doesn't need its own infix entry (and
	// adding one regressed several files outside arithmetic).
	token.QUESTION: LOGICAL,
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
	// inDoubleBracket is set while parsing the body of a `[[ … ]]`
	// conditional. Inside that context `(pat|pat)` is a glob
	// alternation, not a call expression, and adjacent groups like
	// `(a|b)(c|d)` should concatenate as a pattern — not become
	// `(a|b) called with (c|d)`. The flag gates the LPAREN infix
	// so parseCallExpression doesn't fire on pattern groups.
	inDoubleBracket bool

	// consumedBraceTerminator signals that the most recently
	// returned statement already advanced past its own closing
	// RBRACE (e.g. Zsh brace-form `if cond { body }`). Tells the
	// outer parseBlockStatement to skip its own post-statement
	// nextToken so we don't overshoot the following statement's
	// head. parseBlockStatement clears the flag after honouring it.
	consumedBraceTerminator bool
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// Zsh `return` can legitimately appear as the right-hand side
	// of a logical chain (`cmd || return`, `[[ … ]] && return 0`).
	// Top-level statement parsing of the RETURN keyword still wins
	// via the parseStatement switch because that runs before
	// expression dispatch. Registering it as a prefix only matters
	// when the parser is already mid-expression (OR/AND folded as
	// infix into an expression chain with `return` on the RHS).
	p.registerPrefix(token.RETURN, p.parseKeywordAsCommand)
	// Zsh `=cmd` substitutes the absolute path of `cmd` (equivalent
	// to `$(which cmd)`). Lexer emits ASSIGN+IDENT; parseEqualsForm
	// fuses them when the `=` has no preceding space. Only relevant
	// as a statement head — infix ASSIGN still handles `x=y`.
	p.registerPrefix(token.ASSIGN, p.parseEqualsForm)
	// DOT as a prefix models literal-word contexts like `*.zsh`
	// inside a glob, `.*` as a Zsh conditional pattern, or `./path`
	// inside an argument list. Wrap the dot in an Identifier and
	// let parseCommandWord / the bracket scanner fold surrounding
	// tokens in; without this, every dot in a conditional or
	// subscript expression fired "no prefix parse function for .".
	p.registerPrefix(token.DOT, p.parseIdentifier)
	// PERCENT as a prefix handles the prompt-escape words `%F{…}`,
	// `%B`, `%~`, `%n`, `%m` etc. that appear as bare argument
	// tokens in theme files across oh-my-zsh. Treat the percent as
	// a literal word; surrounding tokens concatenate via
	// parseCommandWord. Without this, every prompt-style argument
	// produced "no prefix parse function for %".
	p.registerPrefix(token.PERCENT, p.parseIdentifier)
	// SLASH as a prefix covers path-literal arguments like `/`,
	// `/tmp`, `/usr/bin/*`, where the leading slash starts a
	// command-word. Without this the test `[[ "$dir" != / ]]`
	// fired "no prefix parse function for /". SLASH has no infix
	// entry so this only widens prefix acceptance.
	p.registerPrefix(token.SLASH, p.parseIdentifier)
	// RBRACKET as a prefix covers literal `]` arguments — e.g.
	// Zsh `alias ]=cat` defines an alias whose name is the
	// single-character `]`. Without a prefix the lexer's RBRACKET
	// token had nowhere to land at statement / argument position.
	p.registerPrefix(token.RBRACKET, p.parseIdentifier)
	// COLON as prefix: `dir=:$X` (literal `:` value), `${(j::)`
	// flag forms, and `:` (POSIX null-command) all need the
	// token to fold into an Identifier when it appears as the
	// start of an expression.
	p.registerPrefix(token.COLON, p.parseIdentifier)
	// COMMA as prefix folds bare commas (`function {a,b,c}_x()`,
	// brace-expansion bodies that the parser stumbles into as
	// individual statements) into Identifier tokens so the
	// surrounding block parses cleanly.
	p.registerPrefix(token.COMMA, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	p.registerPrefix(token.CARET, p.parsePrefixExpression)
	p.registerPrefix(token.ASTERISK, p.parsePrefixExpression)
	p.registerPrefix(token.QUESTION, p.parsePrefixExpression)
	p.registerPrefix(token.TILDE, p.parsePrefixExpression)
	// `++x` / `--x` are pre-increment / pre-decrement in Zsh
	// arithmetic (`(( ++x ))`, `(( --x ))`). Register both as
	// prefix; the statement layer routes `cmd --flag arg` through
	// parseSimpleCommandStatement before the expression path sees
	// the DEC token, so prefix DEC is only reached inside actual
	// arithmetic or where a long flag survives as an argument to
	// the expression path (which the simple-command layer drains).
	p.registerPrefix(token.INC, p.parsePrefixExpression)
	p.registerPrefix(token.DEC, p.parsePrefixExpression)
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
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.QUESTION, p.parseTernaryExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.PLUSEQ, p.parseInfixExpression)
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

// lookaheadCaseLabelOpener reports true when the leading `(` at
// curToken is the optional case-label opener (classic `( pat ) body
// ;;` form) rather than the start of a glob-alternation pattern
// (`(a|b)*) body ;;`). The reliable signal is whether peek after
// `(` has space before it — `( pat )` always has whitespace between
// `(` and the pattern, while glob `(pat|pat)` does not. Default to
// TRUE; the curToken=`)` recovery branch in parseCaseStatement
// handles classification mistakes.
func (p *Parser) lookaheadCaseLabelOpener() bool {
	// `( pat …` (with space) is the case-label opener.
	// `(pat|…` (no space) is glob alternation.
	return p.peekToken.HasPrecedingSpace
}

// peekOnSameLogicalLine reports whether the peek token is part of the
// same logical command as the current one. A Zsh `\<NL>` pair joins
// two physical lines into one command; the lexer marks the first token
// after such a pair with HasPrecedingContinuation so argument-gathering
// loops don't terminate at the newline. Multi-line tokens (strings,
// heredocs) report their closing line via EndLine; falling back to
// Line when EndLine is zero.
func (p *Parser) peekOnSameLogicalLine() bool {
	curEnd := p.curToken.EndLine
	if curEnd == 0 {
		curEnd = p.curToken.Line
	}
	return p.peekToken.Line == curEnd || p.peekToken.HasPrecedingContinuation
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

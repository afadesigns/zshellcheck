package ast

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node

	statementNode()
}

type Expression interface {
	Node

	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out []byte
	for _, s := range p.Statements {
		out = append(out, s.String()...)
	}
	return string(out)
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out []byte

	out = append(out, ls.TokenLiteral()...)
	out = append(out, " ", ls.Name.String(), " = ")

	if ls.Value != nil {
		out = append(out, ls.Value.String()...)
	}

	out = append(out, ";")

	return string(out)
}

type ReturnStatement struct {
	Token token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out []byte

	out = append(out, rs.TokenLiteral()...)
	out = append(out, " ")

	if rs.ReturnValue != nil {
		out = append(out, rs.ReturnValue.String()...)
	}

	out = append(out, ";")

	return string(out)
}

type ExpressionStatement struct {
	Token token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token // the token.TRUE or token.FALSE token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The operator token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out []byte

	out = append(out, "(")
	out = append(out, pe.Operator...)
	out = append(out, pe.Right.String()...)
	out = append(out, ")")

	return string(out)
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out []byte

	out = append(out, "(")
	out = append(out, ie.Left.String()...)
	out = append(out, " ", ie.Operator, " ")
	out = append(out, ie.Right.String()...)
	out = append(out, ")")

	return string(out)
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out []byte

	out = append(out, "if")
	out = append(out, ie.Condition.String()...)
	out = append(out, " ")
	out = append(out, ie.Consequence.String()...)

	if ie.Alternative != nil {
		out = append(out, "else ")
		out = append(out, ie.Alternative.String()...)
	}

	return string(out)
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out []byte

	for _, s := range bs.Statements {
		out = append(out, s.String()...)
	}

	return string(out)
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out []byte

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out = append(out, fl.TokenLiteral()...)
	out = append(out, "(")
	out = append(out, strings.Join(params, ", ")...)
	out = append(out, "){")
	out = append(out, fl.Body.String()...)
	out = append(out, "}")

	return string(out)
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out []byte

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out = append(out, ce.Function.String()...)
	out = append(out, "(")
	out = append(out, strings.Join(args, ", ")...)
	out = append(out, ")")

	return string(out)
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
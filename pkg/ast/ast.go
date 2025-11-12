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
		out = append(out, []byte(s.String())...)
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

	out = append(out, []byte(ls.TokenLiteral())...)
	out = append(out, []byte(" ")...)
	out = append(out, []byte(ls.Name.String())...)
	out = append(out, []byte(" = ")...)

	if ls.Value != nil {
		out = append(out, []byte(ls.Value.String())...)
	}

	out = append(out, []byte(";")...)

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

	out = append(out, []byte(rs.TokenLiteral())...)
	out = append(out, []byte(" ")...)

	if rs.ReturnValue != nil {
		out = append(out, []byte(rs.ReturnValue.String())...)
	}

	out = append(out, []byte(";")...)

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

	out = append(out, []byte("(")...)
	out = append(out, []byte(pe.Operator)...)
	out = append(out, []byte(pe.Right.String())...)
	out = append(out, []byte(")")...)

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

	out = append(out, []byte("(")...)
	out = append(out, []byte(ie.Left.String())...)
	out = append(out, []byte(" ")...)
	out = append(out, []byte(ie.Operator)...)
	out = append(out, []byte(" ")...)
	out = append(out, []byte(ie.Right.String())...)
	out = append(out, []byte(")")...)

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

	out = append(out, []byte("if")...)
	out = append(out, []byte(ie.Condition.String())...)
	out = append(out, []byte(" ")...)
	out = append(out, []byte(ie.Consequence.String())...)

	if ie.Alternative != nil {
		out = append(out, []byte("else ")...)
		out = append(out, []byte(ie.Alternative.String())...)
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
		out = append(out, []byte(s.String())...)
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

	out = append(out, []byte(fl.TokenLiteral())...)
	out = append(out, []byte("(")...)
	out = append(out, []byte(strings.Join(params, ", "))...)
	out = append(out, []byte("){")...)
	out = append(out, []byte(fl.Body.String())...)
	out = append(out, []byte("}")...)

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

	out = append(out, []byte(ce.Function.String())...)
	out = append(out, []byte("(")...)
	out = append(out, []byte(strings.Join(args, ", "))...)
	out = append(out, []byte(")")...)

	return string(out)
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// WalkFn is the type of the function called for each node visited by Walk.
// The return value of WalkFn controls how Walk proceeds.
// If true, the children of the node are visited.
// If false, the children of the node are skipped.
type WalkFn func(node Node) bool

// Walk traverses an AST in depth-first order: it starts by calling f(node);
// if f returns true, it then calls Walk on each of the children of node.
func Walk(node Node, f WalkFn) {
	if node == nil {
		return
	}
	if !f(node) {
		return
	}

	switch n := node.(type) {
	case *Program:
		for _, stmt := range n.Statements {
			Walk(stmt, f)
		}

	case *LetStatement:
		Walk(n.Name, f)
		Walk(n.Value, f)

	case *ReturnStatement:
		Walk(n.ReturnValue, f)

	case *ExpressionStatement:
		Walk(n.Expression, f)

	case *BlockStatement:
		for _, stmt := range n.Statements {
			Walk(stmt, f)
		}

	case *Identifier:
		// Leaf node, nothing to walk further

	case *IntegerLiteral:
		// Leaf node

	case *Boolean:
		// Leaf node

	case *PrefixExpression:
		Walk(n.Right, f)

	case *InfixExpression:
		Walk(n.Left, f)
		Walk(n.Right, f)

	case *IfExpression:
		Walk(n.Condition, f)
		Walk(n.Consequence, f)
		Walk(n.Alternative, f)

	case *FunctionLiteral:
		for _, param := range n.Parameters {
			Walk(param, f)
		}
		Walk(n.Body, f)

	case *CallExpression:
		Walk(n.Function, f)
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}

	case *StringLiteral:
		// Leaf node
	}
}
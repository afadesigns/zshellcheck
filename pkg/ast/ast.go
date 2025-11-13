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
	Token       token.Token // the token.RETURN token
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
	Token      token.Token // the first token of the expression
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
	if pe.Right != nil {
		out = append(out, []byte(pe.Right.String())...)
	}
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
	if ie.Left != nil {
		out = append(out, []byte(ie.Left.String())...)
	}
	out = append(out, []byte(" ")...)
	out = append(out, []byte(ie.Operator)...)
	out = append(out, []byte(" ")...)
	if ie.Right != nil {
		out = append(out, []byte(ie.Right.String())...)
	}
	out = append(out, []byte(")")...)
	return string(out)
}

type BlockStatement struct {
	Token      token.Token // the { token or then token
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

type IfStatement struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out []byte
	out = append(out, []byte("if ")...)
	if is.Condition != nil {
		out = append(out, []byte(is.Condition.String())...)
	}
	out = append(out, []byte(" then ")...)
	if is.Consequence != nil {
		out = append(out, []byte(is.Consequence.String())...)
	}
	if is.Alternative != nil {
		out = append(out, []byte(" else ")...)
		out = append(out, []byte(is.Alternative.String())...)
	}
	out = append(out, []byte(" fi")...)
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

type BracketExpression struct {
	Token       token.Token // The '[' token
	Expressions []Expression
}

func (be *BracketExpression) expressionNode()      {}
func (be *BracketExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BracketExpression) String() string {
	var out []byte
	out = append(out, []byte("[")...)
	args := []string{}
	for _, e := range be.Expressions {
		args = append(args, e.String())
	}
	out = append(out, []byte(strings.Join(args, " "))...)
	out = append(out, []byte("]")...)
	return string(out)
}

type DoubleBracketExpression struct {
	Token       token.Token // The '[[' token
	Expressions []Expression
}

func (dbe *DoubleBracketExpression) expressionNode()      {}
func (dbe *DoubleBracketExpression) TokenLiteral() string { return dbe.Token.Literal }
func (dbe *DoubleBracketExpression) String() string {
	var out []byte
	out = append(out, []byte("[[")...)
	args := []string{}
	for _, e := range dbe.Expressions {
		args = append(args, e.String())
	}
	out = append(out, []byte(strings.Join(args, " "))...)
	out = append(out, []byte("]]")...)
	return string(out)
}


type ArrayAccess struct {
	Token token.Token // The '${' token
	Left  Expression
	Index Expression
}

func (aa *ArrayAccess) expressionNode()      {}
func (aa *ArrayAccess) TokenLiteral() string { return aa.Token.Literal }
func (aa *ArrayAccess) String() string {
	var out []byte
	out = append(out, []byte("${")...)
	out = append(out, []byte(aa.Left.String())...)
	out = append(out, []byte("[")...)
	out = append(out, []byte(aa.Index.String())...)
	out = append(out, []byte("]}")...)
	return string(out)
}

type CommandSubstitution struct {
	Token   token.Token // The '`' token
	Command Expression
}

func (cs *CommandSubstitution) expressionNode()      {}
func (cs *CommandSubstitution) TokenLiteral() string { return cs.Token.Literal }
func (cs *CommandSubstitution) String() string {
	var out []byte
	out = append(out, []byte("`")...)
	out = append(out, []byte(cs.Command.String())...)
	out = append(out, []byte("`")...)
	return string(out)
}

type SimpleCommand struct {
	Token     token.Token // The first token of the command
	Name      Expression
	Arguments []Expression
}

func (sc *SimpleCommand) expressionNode()      {}
func (sc *SimpleCommand) TokenLiteral() string { return sc.Token.Literal }
func (sc *SimpleCommand) String() string {
	var out []byte
	args := []string{}
	for _, a := range sc.Arguments {
		args = append(args, a.String())
	}
	out = append(out, []byte(sc.Name.String())...)
	out = append(out, []byte(" ")...)
	out = append(out, []byte(strings.Join(args, " "))...)
	return string(out)
}

type WalkFn func(node Node) bool

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
	case *IntegerLiteral:
	case *Boolean:
	case *PrefixExpression:
		if n.Right != nil {
			Walk(n.Right, f)
		}
	case *InfixExpression:
		if n.Left != nil {
			Walk(n.Left, f)
		}
		if n.Right != nil {
			Walk(n.Right, f)
		}
	case *IfStatement:
		if n.Condition != nil {
			Walk(n.Condition, f)
		}
		if n.Consequence != nil {
			Walk(n.Consequence, f)
		}
		if n.Alternative != nil {
			Walk(n.Alternative, f)
		}
	case *FunctionLiteral:
		for _, param := range n.Parameters {
			Walk(param, f)
		}
		if n.Body != nil {
			Walk(n.Body, f)
		}
	case *CallExpression:
		if n.Function != nil {
			Walk(n.Function, f)
		}
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}
	case *StringLiteral:
	case *BracketExpression:
		for _, exp := range n.Expressions {
			Walk(exp, f)
		}
	case *DoubleBracketExpression:
		for _, exp := range n.Expressions {
			Walk(exp, f)
		}
	case *ArrayAccess:
		Walk(n.Left, f)
		Walk(n.Index, f)
	case *CommandSubstitution:
	case *SimpleCommand:
		Walk(n.Name, f)
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}
	}
}
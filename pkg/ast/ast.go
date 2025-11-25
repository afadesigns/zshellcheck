package ast

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

type NodeType int

const (
	ProgramNode NodeType = iota
	LetStatementNode
	ReturnStatementNode
	ExpressionStatementNode
	IdentifierNode
	IntegerLiteralNode
	BooleanNode
	PrefixExpressionNode
	PostfixExpressionNode
	InfixExpressionNode
	BlockStatementNode
	IfStatementNode
	ForLoopStatementNode
	WhileLoopStatementNode
	FunctionLiteralNode
	CallExpressionNode
	StringLiteralNode
	BracketExpressionNode
	DoubleBracketExpressionNode
	ArrayAccessNode
	InvalidArrayAccessNode
	CommandSubstitutionNode
	ShebangNode
	DollarParenExpressionNode
	SimpleCommandNode
	IndexExpressionNode
	ConcatenatedExpressionNode
	CaseStatementNode
	RedirectionNode
	FunctionDefinitionNode
	GroupedExpressionNode
	ArithmeticCommandNode
	SubshellNode
)

type Subshell struct {
	Token token.Token // The '(' token
	Block *BlockStatement
}

func (s *Subshell) Type() NodeType       { return SubshellNode }
func (s *Subshell) statementNode()       {}
func (s *Subshell) TokenLiteral() string { return s.Token.Literal }
func (s *Subshell) String() string {
	return "( " + s.Block.String() + " )"
}

type ArithmeticCommand struct {
	Token      token.Token // The '((' token
	Expression Expression
}

func (ac *ArithmeticCommand) Type() NodeType                { return ArithmeticCommandNode }
func (ac *ArithmeticCommand) statementNode()                {}
func (ac *ArithmeticCommand) TokenLiteral() string          { return ac.Token.Literal }
func (ac *ArithmeticCommand) TokenLiteralNode() token.Token { return ac.Token }
func (ac *ArithmeticCommand) String() string {
	return "((" + ac.Expression.String() + "))"
}

type GroupedExpression struct {
	Token token.Token // The '(' token
	Exp   Expression
}

func (ge *GroupedExpression) Type() NodeType                { return GroupedExpressionNode }
func (ge *GroupedExpression) expressionNode()               {}
func (ge *GroupedExpression) TokenLiteral() string          { return ge.Token.Literal }
func (ge *GroupedExpression) TokenLiteralNode() token.Token { return ge.Token }
func (ge *GroupedExpression) String() string {
	return "(" + ge.Exp.String() + ")"
}

type FunctionDefinition struct {
	Token token.Token // The name token
	Name  *Identifier
	Body  Statement // The function body (usually BlockStatement)
}

func (fd *FunctionDefinition) Type() NodeType       { return FunctionDefinitionNode }
func (fd *FunctionDefinition) expressionNode()      {}
func (fd *FunctionDefinition) statementNode()       {}
func (fd *FunctionDefinition) TokenLiteral() string { return fd.Token.Literal }
func (fd *FunctionDefinition) String() string {
	var out []byte
	if fd.Name != nil {
		out = append(out, []byte(fd.Name.String())...)
	}
	out = append(out, []byte("() ")...)
	if fd.Body != nil {
		out = append(out, []byte(fd.Body.String())...)
	}
	return string(out)
}

type Redirection struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (r *Redirection) Type() NodeType       { return RedirectionNode }
func (r *Redirection) expressionNode()      {}
func (r *Redirection) TokenLiteral() string { return r.Token.Literal }
func (r *Redirection) String() string {
	var out []byte
	if r.Left != nil {
		out = append(out, []byte(r.Left.String())...)
	}
	out = append(out, []byte(" ")...)
	out = append(out, []byte(r.Operator)...)
	out = append(out, []byte(" ")...)
	if r.Right != nil {
		out = append(out, []byte(r.Right.String())...)
	}
	return string(out)
}

type Node interface {
	TokenLiteral() string
	String() string
	Type() NodeType
	TokenLiteralNode() token.Token
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

func (p *Program) Type() NodeType { return ProgramNode }
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) TokenLiteralNode() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteralNode()
	}
	return token.Token{}
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

func (ls *LetStatement) Type() NodeType       { return LetStatementNode }
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

func (rs *ReturnStatement) Type() NodeType       { return ReturnStatementNode }
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

func (es *ExpressionStatement) Type() NodeType       { return ExpressionStatementNode }
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

func (i *Identifier) Type() NodeType       { return IdentifierNode }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (il *IntegerLiteral) Type() NodeType       { return IntegerLiteralNode }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token // the token.TRUE or token.FALSE token
	Value bool
}

func (b *Boolean) Type() NodeType       { return BooleanNode }
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The operator token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Type() NodeType       { return PrefixExpressionNode }
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

type PostfixExpression struct {
	Token    token.Token // The operator token, e.g. ++
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) Type() NodeType       { return PostfixExpressionNode }
func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out []byte
	out = append(out, []byte("(")...)
	if pe.Left != nil {
		out = append(out, []byte(pe.Left.String())...)
	}
	out = append(out, []byte(pe.Operator)...)
	out = append(out, []byte(")")...)
	return string(out)
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) Type() NodeType       { return InfixExpressionNode }
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

func (bs *BlockStatement) Type() NodeType       { return BlockStatementNode }
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
	Condition   *BlockStatement
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) Type() NodeType       { return IfStatementNode }
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

type ForLoopStatement struct {
	Token     token.Token  // The 'for' token
	Init      Expression   // Variable name (for-each) or Init expr (arithmetic)
	Condition Expression   // Arithmetic condition
	Post      Expression   // Arithmetic post
	Items     []Expression // Items to iterate over (for-each)
	Body      *BlockStatement
}

func (fls *ForLoopStatement) Type() NodeType       { return ForLoopStatementNode }
func (fls *ForLoopStatement) statementNode()       {}
func (fls *ForLoopStatement) TokenLiteral() string { return fls.Token.Literal }
func (fls *ForLoopStatement) String() string {
	var out []byte
	// Heuristic: if Condition or Post is present, or Items is nil (and not implicit?), it's arithmetic?
	// Actually, explicit `Items` makes it for-each.
	// But `for i` (implicit in) has Items=nil.
	// Arithmetic `for ((...))` usually has params. `for ((;;))` is possible.
	// I'll assume if Items is non-nil (even empty) it's for-each.
	// Or if Init is Identifier and others are nil?

	if fls.Items != nil {
		out = append(out, []byte("for ")...)
		if fls.Init != nil {
			out = append(out, []byte(fls.Init.String())...)
		}
		out = append(out, []byte(" in ")...)
		for _, item := range fls.Items {
			out = append(out, []byte(item.String())...)
			out = append(out, []byte(" ")...)
		}
		out = append(out, []byte("; do ")...)
	} else {
		out = append(out, []byte("for ((")...)
		if fls.Init != nil {
			out = append(out, []byte(fls.Init.String())...)
		}
		out = append(out, []byte("; ")...)
		if fls.Condition != nil {
			out = append(out, []byte(fls.Condition.String())...)
		}
		out = append(out, []byte("; ")...)
		if fls.Post != nil {
			out = append(out, []byte(fls.Post.String())...)
		}
		out = append(out, []byte(")); do ")...)
	}

	if fls.Body != nil {
		out = append(out, []byte(fls.Body.String())...)
	}
	out = append(out, []byte("done")...)
	return string(out)
}

type WhileLoopStatement struct {
	Token     token.Token // The 'while' token
	Condition *BlockStatement
	Body      *BlockStatement
}

func (wls *WhileLoopStatement) Type() NodeType       { return WhileLoopStatementNode }
func (wls *WhileLoopStatement) statementNode()       {}
func (wls *WhileLoopStatement) TokenLiteral() string { return wls.Token.Literal }
func (wls *WhileLoopStatement) String() string {
	var out []byte
	out = append(out, []byte("while ")...)
	if wls.Condition != nil {
		out = append(out, []byte(wls.Condition.String())...)
	}
	out = append(out, []byte("; do ")...)
	if wls.Body != nil {
		out = append(out, []byte(wls.Body.String())...)
	}
	out = append(out, []byte("done")...)
	return string(out)
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Name       string      // The function name (optional)
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Type() NodeType       { return FunctionLiteralNode }
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

func (ce *CallExpression) Type() NodeType       { return CallExpressionNode }
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

func (sl *StringLiteral) Type() NodeType       { return StringLiteralNode }
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type BracketExpression struct {
	Token       token.Token // The '[' token
	Expressions []Expression
}

func (be *BracketExpression) Type() NodeType       { return BracketExpressionNode }
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

func (dbe *DoubleBracketExpression) Type() NodeType       { return DoubleBracketExpressionNode }
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

func (aa *ArrayAccess) Type() NodeType       { return ArrayAccessNode }
func (aa *ArrayAccess) expressionNode()      {}
func (aa *ArrayAccess) TokenLiteral() string { return aa.Token.Literal }
func (aa *ArrayAccess) String() string {
	var out []byte
	out = append(out, []byte("${")...)
	if aa.Left != nil {
		out = append(out, []byte(aa.Left.String())...)
	}
	if aa.Index != nil {
		out = append(out, []byte("[")...)
		out = append(out, []byte(aa.Index.String())...)
		out = append(out, []byte("]")...)
	}
	out = append(out, []byte("}")...)
	return string(out)
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) Type() NodeType       { return IndexExpressionNode }
func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out []byte
	out = append(out, []byte("(")...)
	if ie.Left != nil {
		out = append(out, []byte(ie.Left.String())...)
	}
	out = append(out, []byte("[")...)
	if ie.Index != nil {
		out = append(out, []byte(ie.Index.String())...)
	}
	out = append(out, []byte("])")...)
	return string(out)
}

type InvalidArrayAccess struct {
	Token token.Token // The '$' token
	Left  Expression
	Index Expression
}

func (iaa *InvalidArrayAccess) Type() NodeType       { return InvalidArrayAccessNode }
func (iaa *InvalidArrayAccess) expressionNode()      {}
func (iaa *InvalidArrayAccess) TokenLiteral() string { return iaa.Token.Literal }
func (iaa *InvalidArrayAccess) String() string {
	var out []byte
	out = append(out, []byte("$")...)
	out = append(out, []byte(iaa.Left.String())...)
	out = append(out, []byte("[")...)
	out = append(out, []byte(iaa.Index.String())...)
	out = append(out, []byte("]")...)
	return string(out)
}

type CommandSubstitution struct {
	Token   token.Token // The ` or $() token
	Command Expression
}

func (cs *CommandSubstitution) Type() NodeType       { return CommandSubstitutionNode }
func (cs *CommandSubstitution) expressionNode()      {}
func (cs *CommandSubstitution) TokenLiteral() string { return cs.Token.Literal }
func (cs *CommandSubstitution) String() string {
	return "`" + cs.Command.String() + "`"
}

type Shebang struct {
	Token token.Token // The #! token
	Path  string
}

func (s *Shebang) Type() NodeType       { return ShebangNode }
func (s *Shebang) statementNode()       {}
func (s *Shebang) TokenLiteral() string { return s.Token.Literal }
func (s *Shebang) String() string {
	return s.Token.Literal
}

type DollarParenExpression struct {
	Token   token.Token // The '$(' token
	Command Expression
}

func (dpe *DollarParenExpression) Type() NodeType       { return DollarParenExpressionNode }
func (dpe *DollarParenExpression) expressionNode()      {}
func (dpe *DollarParenExpression) TokenLiteral() string { return dpe.Token.Literal }
func (dpe *DollarParenExpression) String() string {
	var out []byte
	out = append(out, []byte("$(")...)
	out = append(out, []byte(dpe.Command.String())...)
	out = append(out, []byte(")")...)
	return string(out)
}

type SimpleCommand struct {
	Token     token.Token // The first token of the command
	Name      Expression
	Arguments []Expression
}

func (sc *SimpleCommand) Type() NodeType       { return SimpleCommandNode }
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

type ConcatenatedExpression struct {
	Token token.Token
	Parts []Expression
}

func (ce *ConcatenatedExpression) Type() NodeType       { return ConcatenatedExpressionNode }
func (ce *ConcatenatedExpression) expressionNode()      {}
func (ce *ConcatenatedExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ConcatenatedExpression) String() string {
	var out []byte
	for _, p := range ce.Parts {
		if p != nil {
			out = append(out, []byte(p.String())...)
		}
	}
	return string(out)
}

type CaseStatement struct {
	Token   token.Token // The 'case' token
	Value   Expression
	Clauses []*CaseClause
}

func (cs *CaseStatement) Type() NodeType       { return CaseStatementNode }
func (cs *CaseStatement) statementNode()       {}
func (cs *CaseStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CaseStatement) String() string {
	var out []byte
	out = append(out, []byte("case ")...)
	if cs.Value != nil {
		out = append(out, []byte(cs.Value.String())...)
	}
	out = append(out, []byte(" in ")...)
	for _, c := range cs.Clauses {
		out = append(out, []byte(c.String())...)
	}
	out = append(out, []byte("esac")...)
	return string(out)
}

type CaseClause struct {
	Token    token.Token // The first token of pattern
	Patterns []Expression
	Body     *BlockStatement
}

func (cc *CaseClause) String() string {
	var out []byte
	for i, p := range cc.Patterns {
		out = append(out, []byte(p.String())...)
		if i < len(cc.Patterns)-1 {
			out = append(out, []byte(" | ")...)
		}
	}
	out = append(out, []byte(") ")...)
	if cc.Body != nil {
		out = append(out, []byte(cc.Body.String())...)
	}
	out = append(out, []byte(" ;; ")...)
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
	case *PostfixExpression:
		if n.Left != nil {
			Walk(n.Left, f)
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
	case *ForLoopStatement:
		Walk(n.Init, f)
		Walk(n.Condition, f)
		Walk(n.Post, f)
		for _, item := range n.Items {
			Walk(item, f)
		}
		Walk(n.Body, f)
	case *WhileLoopStatement:
		Walk(n.Condition, f)
		Walk(n.Body, f)
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
	case *IndexExpression:
		Walk(n.Left, f)
		Walk(n.Index, f)
	case *InvalidArrayAccess:
		Walk(n.Left, f)
		Walk(n.Index, f)
	case *CommandSubstitution:
		Walk(n.Command, f)
	case *DollarParenExpression:
		Walk(n.Command, f)
	case *SimpleCommand:
		Walk(n.Name, f)
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}
	case *ConcatenatedExpression:
		for _, p := range n.Parts {
			Walk(p, f)
		}
	case *CaseStatement:
		Walk(n.Value, f)
		for _, c := range n.Clauses {
			for _, p := range c.Patterns {
				Walk(p, f)
			}
			Walk(c.Body, f)
		}
	case *Redirection:
		Walk(n.Left, f)
		Walk(n.Right, f)
	case *FunctionDefinition:
		Walk(n.Name, f)
		Walk(n.Body, f)
	case *GroupedExpression:
		Walk(n.Exp, f)
	case *ArithmeticCommand:
		Walk(n.Expression, f)
	case *Subshell:
		Walk(n.Block, f)
	}
}

func (n *Subshell) TokenLiteralNode() token.Token                { return n.Token }
func (n *FunctionDefinition) TokenLiteralNode() token.Token      { return n.Token }
func (n *Redirection) TokenLiteralNode() token.Token             { return n.Token }
func (n *LetStatement) TokenLiteralNode() token.Token            { return n.Token }
func (n *ReturnStatement) TokenLiteralNode() token.Token         { return n.Token }
func (n *ExpressionStatement) TokenLiteralNode() token.Token     { return n.Token }
func (n *Identifier) TokenLiteralNode() token.Token              { return n.Token }
func (n *IntegerLiteral) TokenLiteralNode() token.Token          { return n.Token }
func (n *Boolean) TokenLiteralNode() token.Token                 { return n.Token }
func (n *PrefixExpression) TokenLiteralNode() token.Token        { return n.Token }
func (n *PostfixExpression) TokenLiteralNode() token.Token       { return n.Token }
func (n *InfixExpression) TokenLiteralNode() token.Token         { return n.Token }
func (n *BlockStatement) TokenLiteralNode() token.Token          { return n.Token }
func (n *IfStatement) TokenLiteralNode() token.Token             { return n.Token }
func (n *ForLoopStatement) TokenLiteralNode() token.Token        { return n.Token }
func (n *WhileLoopStatement) TokenLiteralNode() token.Token      { return n.Token }
func (n *FunctionLiteral) TokenLiteralNode() token.Token         { return n.Token }
func (n *CallExpression) TokenLiteralNode() token.Token          { return n.Token }
func (n *StringLiteral) TokenLiteralNode() token.Token           { return n.Token }
func (n *BracketExpression) TokenLiteralNode() token.Token       { return n.Token }
func (n *DoubleBracketExpression) TokenLiteralNode() token.Token { return n.Token }
func (n *ArrayAccess) TokenLiteralNode() token.Token             { return n.Token }
func (n *IndexExpression) TokenLiteralNode() token.Token         { return n.Token }
func (n *InvalidArrayAccess) TokenLiteralNode() token.Token      { return n.Token }
func (n *CommandSubstitution) TokenLiteralNode() token.Token     { return n.Token }
func (n *Shebang) TokenLiteralNode() token.Token                 { return n.Token }
func (n *DollarParenExpression) TokenLiteralNode() token.Token   { return n.Token }
func (n *SimpleCommand) TokenLiteralNode() token.Token           { return n.Token }
func (n *ConcatenatedExpression) TokenLiteralNode() token.Token  { return n.Token }
func (n *CaseStatement) TokenLiteralNode() token.Token           { return n.Token }
func (n *CaseClause) TokenLiteralNode() token.Token              { return n.Token }

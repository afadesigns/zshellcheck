package ast

import (
	"fmt"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

// Node represents a node in the AST.
type Node interface {
	TokenLiteral() string // literal value of the node
	TokenLiteralNode() token.Token
	String() string // string representation of the node
}

// Statement represents a statement in the AST.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the AST.
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string          { return "" }
func (p *Program) TokenLiteralNode() token.Token { return token.Token{} }

// String returns a string representation of the Program.
func (p *Program) String() string {
	var sb strings.Builder
	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

// LetStatement represents a let statement.
type LetStatement struct {
	Token token.Token // the 'let' token
	Name  *Identifier // identifier node
	Value Expression  // expression node
}

func (ls *LetStatement) statementNode()                {}
func (ls *LetStatement) expressionNode()               {}
func (ls *LetStatement) TokenLiteral() string          { return ls.Token.Literal }
func (ls *LetStatement) TokenLiteralNode() token.Token { return ls.Token }

// String returns a string representation of the LetStatement.
func (ls *LetStatement) String() string {
	return ls.Token.Literal + " " + ls.Name.String() + " = " + ls.Value.String() + ";"
}

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression  // expression node
}

func (rs *ReturnStatement) statementNode()                {}
func (rs *ReturnStatement) expressionNode()               {}
func (rs *ReturnStatement) TokenLiteral() string          { return rs.Token.Literal }
func (rs *ReturnStatement) TokenLiteralNode() token.Token { return rs.Token }

// String returns a string representation of the ReturnStatement.
func (rs *ReturnStatement) String() string {
	if rs.ReturnValue != nil {
		return rs.Token.Literal + " " + rs.ReturnValue.String()
	}
	return rs.Token.Literal
}

// ExpressionStatement represents an expression statement.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression  // expression node
}

func (es *ExpressionStatement) statementNode()                {}
func (es *ExpressionStatement) expressionNode()               {}
func (es *ExpressionStatement) TokenLiteral() string          { return es.Token.Literal }
func (es *ExpressionStatement) TokenLiteralNode() token.Token { return es.Token }

// String returns a string representation of the ExpressionStatement.
func (es *ExpressionStatement) String() string {
	return es.Expression.String()
}

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token // the IDENT token
	Value string      // identifier name
}

func (i *Identifier) expressionNode()               {}
func (i *Identifier) TokenLiteral() string          { return i.Token.Literal }
func (i *Identifier) TokenLiteralNode() token.Token { return i.Token }

// String returns a string representation of the Identifier.
func (i *Identifier) String() string { return i.Value }

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token // the INT token
	Value int64       // integer value
}

func (il *IntegerLiteral) expressionNode()               {}
func (il *IntegerLiteral) TokenLiteral() string          { return il.Token.Literal }
func (il *IntegerLiteral) TokenLiteralNode() token.Token { return il.Token }

// String returns a string representation of the IntegerLiteral.
func (il *IntegerLiteral) String() string { return fmt.Sprintf("%d", il.Value) }

// Boolean represents a boolean literal.
type Boolean struct {
	Token token.Token // the TRUE or FALSE token
	Value bool        // boolean value
}

func (b *Boolean) expressionNode()               {}
func (b *Boolean) TokenLiteral() string          { return b.Token.Literal }
func (b *Boolean) TokenLiteralNode() token.Token { return b.Token }

// String returns a string representation of the Boolean.
func (b *Boolean) String() string { return fmt.Sprintf("%t", b.Value) }

// PrefixExpression represents a prefix expression (e.g., -5, !true).
type PrefixExpression struct {
	Token    token.Token // The prefix operator token (! or -)
	Operator string      // Prefix operator
	Right    Expression  // Right operand
}

func (pe *PrefixExpression) expressionNode()               {}
func (pe *PrefixExpression) TokenLiteral() string          { return pe.Token.Literal }
func (pe *PrefixExpression) TokenLiteralNode() token.Token { return pe.Token }

// String returns a string representation of the PrefixExpression.
func (pe *PrefixExpression) String() string {
	var right string
	if pe.Right != nil {
		right = pe.Right.String()
	}
	return "(" + pe.Operator + right + ")"
}

// PostfixExpression represents a postfix expression (e.g. --i, ++i).
type PostfixExpression struct {
	Token    token.Token // The postfix operator token (-- or ++)
	Operator string      // Postfix operator
	Left     Expression  // Left operand
}

func (pe *PostfixExpression) expressionNode()               {}
func (pe *PostfixExpression) TokenLiteral() string          { return pe.Token.Literal }
func (pe *PostfixExpression) TokenLiteralNode() token.Token { return pe.Token }

// String returns a string representation of the PostfixExpression.
func (pe *PostfixExpression) String() string {
	return "(" + pe.Left.String() + pe.Operator + ")"
}

// InfixExpression represents an infix expression (e.g. 5 + 5).
type InfixExpression struct {
	Token    token.Token // The operator token
	Left     Expression  // Left operand
	Operator string      // Infix operator
	Right    Expression  // Right operand
}

func (ie *InfixExpression) expressionNode()               {}
func (ie *InfixExpression) TokenLiteral() string          { return ie.Token.Literal }
func (ie *InfixExpression) TokenLiteralNode() token.Token { return ie.Token }

// String returns a string representation of the InfixExpression.
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// BlockStatement represents a block of statements (e.g., in if or function bodies).
type BlockStatement struct {
	Token      token.Token // the '{' token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()                {}
func (bs *BlockStatement) expressionNode()               {}
func (bs *BlockStatement) TokenLiteral() string          { return bs.Token.Literal }
func (bs *BlockStatement) TokenLiteralNode() token.Token { return bs.Token }

// String returns a string representation of the BlockStatement.
func (bs *BlockStatement) String() string {
	var sb strings.Builder
	for _, s := range bs.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

// IfStatement represents an if statement.
type IfStatement struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // Statement, could be *BlockStatement or *ExpressionStatement
}

func (is *IfStatement) statementNode()                {}
func (is *IfStatement) expressionNode()               {}
func (is *IfStatement) TokenLiteral() string          { return is.Token.Literal }
func (is *IfStatement) TokenLiteralNode() token.Token { return is.Token }

// String returns a string representation of the IfStatement.
func (is *IfStatement) String() string {
	var sb strings.Builder
	sb.WriteString(is.Token.Literal + "(" + is.Condition.String() + ") " + is.Consequence.String())
	if is.Alternative != nil {
		sb.WriteString("else " + is.Alternative.String())
	}
	return sb.String()
}

// ForLoopStatement represents a for loop statement.
type ForLoopStatement struct {
	Token     token.Token // The 'for' token
	Init      Expression
	Condition Expression
	Post      Expression
	Name      *Identifier
	Items     []Expression
	Body      *BlockStatement
}

func (fls *ForLoopStatement) statementNode()                {}
func (fls *ForLoopStatement) expressionNode()               {}
func (fls *ForLoopStatement) TokenLiteral() string          { return fls.Token.Literal }
func (fls *ForLoopStatement) TokenLiteralNode() token.Token { return fls.Token }

// String returns a string representation of the ForLoopStatement.
func (fls *ForLoopStatement) String() string {
	var sb strings.Builder
	sb.WriteString(fls.Token.Literal + " ")
	if fls.Name != nil {
		sb.WriteString(fls.Name.String() + " in")
		for _, item := range fls.Items {
			sb.WriteString(" " + item.String())
		}
	} else {
		sb.WriteString("(")
		if fls.Init != nil {
			sb.WriteString(fls.Init.String())
		}
		sb.WriteString("; ")
		if fls.Condition != nil {
			sb.WriteString(fls.Condition.String())
		}
		sb.WriteString("; ")
		if fls.Post != nil {
			sb.WriteString(fls.Post.String())
		}
		sb.WriteString(")")
	}
	sb.WriteString(" ")
	sb.WriteString(fls.Body.String())
	return sb.String()
}

// WhileLoopStatement represents a while loop statement.
type WhileLoopStatement struct {
	Token     token.Token // The 'while' token
	Condition Expression
	Body      *BlockStatement
}

func (wls *WhileLoopStatement) statementNode()                {}
func (wls *WhileLoopStatement) expressionNode()               {}
func (wls *WhileLoopStatement) TokenLiteral() string          { return wls.Token.Literal }
func (wls *WhileLoopStatement) TokenLiteralNode() token.Token { return wls.Token }

// String returns a string representation of the WhileLoopStatement.
func (wls *WhileLoopStatement) String() string {
	var sb strings.Builder
	sb.WriteString(wls.Token.Literal + "(" + wls.Condition.String() + ") ")
	sb.WriteString(wls.Body.String())
	return sb.String()
}

// FunctionLiteral represents a function literal.
type FunctionLiteral struct {
	Token  token.Token // The 'fn' token
	Name   *Identifier
	Params []*Identifier
	Body   *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()               {}
func (fl *FunctionLiteral) TokenLiteral() string          { return fl.Token.Literal }
func (fl *FunctionLiteral) TokenLiteralNode() token.Token { return fl.Token }

// String returns a string representation of the FunctionLiteral.
func (fl *FunctionLiteral) String() string {
	var sb strings.Builder
	sb.WriteString(fl.TokenLiteral())
	sb.WriteString("(")
	params := []string{}
	for _, p := range fl.Params {
		params = append(params, p.String())
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") ")
	sb.WriteString(fl.Body.String())
	return sb.String()
}

// CallExpression represents a function call expression.
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()               {}
func (ce *CallExpression) TokenLiteral() string          { return ce.Token.Literal }
func (ce *CallExpression) TokenLiteralNode() token.Token { return ce.Token }

// String returns a string representation of the CallExpression.
func (ce *CallExpression) String() string {
	var sb strings.Builder
	sb.WriteString(ce.Function.String())
	sb.WriteString("(")
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")")
	return sb.String()
}

// IndexExpression represents an index expression (e.g., arr[0]).
type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()               {}
func (ie *IndexExpression) TokenLiteral() string          { return ie.Token.Literal }
func (ie *IndexExpression) TokenLiteralNode() token.Token { return ie.Token }

// String returns a string representation of the IndexExpression.
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

// BracketExpression represents a bracket expression (e.g. {a, b, c}).
type BracketExpression struct {
	Token    token.Token // The '{' token
	Elements []Expression
}

func (be *BracketExpression) expressionNode()               {}
func (be *BracketExpression) TokenLiteral() string          { return be.Token.Literal }
func (be *BracketExpression) TokenLiteralNode() token.Token { return be.Token }

// String returns a string representation of the BracketExpression.
func (be *BracketExpression) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	elements := []string{}
	for _, el := range be.Elements {
		elements = append(elements, el.String())
	}
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteString("}")
	return sb.String()
}

// DoubleBracketExpression represents a double bracket expression (e.g. [[ a ]]).
type DoubleBracketExpression struct {
	Token    token.Token // The '[[' token
	Elements []Expression
}

func (dbe *DoubleBracketExpression) expressionNode()               {}
func (dbe *DoubleBracketExpression) TokenLiteral() string          { return dbe.Token.Literal }
func (dbe *DoubleBracketExpression) TokenLiteralNode() token.Token { return dbe.Token }

// String returns a string representation of the DoubleBracketExpression.
func (dbe *DoubleBracketExpression) String() string {
	var sb strings.Builder
	sb.WriteString("[[")
	elements := []string{}
	for _, el := range dbe.Elements {
		elements = append(elements, el.String())
	}
	sb.WriteString(strings.Join(elements, ", "))
	sb.WriteString("]]")
	return sb.String()
}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()               {}
func (sl *StringLiteral) TokenLiteral() string          { return sl.Token.Literal }
func (sl *StringLiteral) TokenLiteralNode() token.Token { return sl.Token }
func (sl *StringLiteral) String() string                { return sl.Token.Literal }

// GroupedExpression represents a grouped expression (e.g. (1 + 2)).
type GroupedExpression struct {
	Token      token.Token // The '(' token
	Expression Expression
}

func (ge *GroupedExpression) expressionNode()               {}
func (ge *GroupedExpression) TokenLiteral() string          { return ge.Token.Literal }
func (ge *GroupedExpression) TokenLiteralNode() token.Token { return ge.Token }
func (ge *GroupedExpression) String() string {
	return "(" + ge.Expression.String() + ")"
}

// ArrayAccess represents an array access expression (e.g. ${arr[0]}).
type ArrayAccess struct {
	Token token.Token // The '$' token
	Left  Expression
	Index Expression
}

func (aa *ArrayAccess) expressionNode()               {}
func (aa *ArrayAccess) TokenLiteral() string          { return aa.Token.Literal }
func (aa *ArrayAccess) TokenLiteralNode() token.Token { return aa.Token }

// String returns a string representation of the ArrayAccess.
func (aa *ArrayAccess) String() string {
	var sb strings.Builder
	sb.WriteString("${ ")
	if aa.Left != nil {
		sb.WriteString(aa.Left.String())
	}
	if aa.Index != nil {
		sb.WriteString("[" + aa.Index.String() + "]")
	}
	sb.WriteString("}")
	return sb.String()
}

// CommandSubstitution represents a command substitution (e.g. $(...)).
type CommandSubstitution struct {
	Token   token.Token // The '$(' token
	Command Node
}

func (cs *CommandSubstitution) expressionNode()               {}
func (cs *CommandSubstitution) TokenLiteral() string          { return cs.Token.Literal }
func (cs *CommandSubstitution) TokenLiteralNode() token.Token { return cs.Token }
func (cs *CommandSubstitution) String() string {
	return "$(" + cs.Command.String() + ")"
}

// InvalidArrayAccess represents an invalid array access.
type InvalidArrayAccess struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ia *InvalidArrayAccess) expressionNode()               {}
func (ia *InvalidArrayAccess) TokenLiteral() string          { return ia.Token.Literal }
func (ia *InvalidArrayAccess) TokenLiteralNode() token.Token { return ia.Token }
func (ia *InvalidArrayAccess) String() string {
	return ia.Left.String() + "[" + ia.Index.String() + "]"
}

// ArrayLiteral represents an array literal (e.g., (val1 val2)).
type ArrayLiteral struct {
	Token    token.Token // The '(' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()               {}
func (al *ArrayLiteral) TokenLiteral() string          { return al.Token.Literal }
func (al *ArrayLiteral) TokenLiteralNode() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, el := range al.Elements {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(el.String())
	}
	sb.WriteString(")")
	return sb.String()
}

// Shebang represents a shebang comment (e.g., #!/bin/zsh).
type Shebang struct {
	Token token.Token // The '#!' token
	Path  string
}

func (s *Shebang) statementNode()                {}
func (s *Shebang) expressionNode()               {}
func (s *Shebang) TokenLiteral() string          { return s.Token.Literal }
func (s *Shebang) TokenLiteralNode() token.Token { return s.Token }

// String returns a string representation of the Shebang.
func (s *Shebang) String() string { return s.Path }

// DollarParenExpression represents a dollar paren expression (e.g., $(command)).
type DollarParenExpression struct {
	Token   token.Token // The '$(' token
	Command Node
}

func (dpe *DollarParenExpression) expressionNode()               {}
func (dpe *DollarParenExpression) TokenLiteral() string          { return dpe.Token.Literal }
func (dpe *DollarParenExpression) TokenLiteralNode() token.Token { return dpe.Token }

// String returns a string representation of the DollarParenExpression.
func (dpe *DollarParenExpression) String() string { return "$(" + dpe.Command.String() + ")" }

// SimpleCommand represents a simple command (e.g., ls -l).
type SimpleCommand struct {
	Token     token.Token
	Name      Expression
	Arguments []Expression
}

func (sc *SimpleCommand) statementNode()                {}
func (sc *SimpleCommand) expressionNode()               {}
func (sc *SimpleCommand) TokenLiteral() string          { return sc.Token.Literal }
func (sc *SimpleCommand) TokenLiteralNode() token.Token { return sc.Token }

// String returns a string representation of the SimpleCommand.
func (sc *SimpleCommand) String() string {
	var sb strings.Builder
	sb.WriteString(sc.Name.String())
	for _, arg := range sc.Arguments {
		sb.WriteString(" " + arg.String())
	}
	return sb.String()
}

// ConcatenatedExpression represents a concatenated expression.
type ConcatenatedExpression struct {
	Token token.Token
	Parts []Expression
}

func (ce *ConcatenatedExpression) expressionNode()               {}
func (ce *ConcatenatedExpression) TokenLiteral() string          { return ce.Token.Literal }
func (ce *ConcatenatedExpression) TokenLiteralNode() token.Token { return ce.Token }

// String returns a string representation of the ConcatenatedExpression.
func (ce *ConcatenatedExpression) String() string {
	var sb strings.Builder
	for _, part := range ce.Parts {
		sb.WriteString(part.String())
	}
	return sb.String()
}

// CaseStatement represents a case statement.
type CaseStatement struct {
	Token   token.Token
	Value   Expression
	Clauses []*CaseClause
}

func (cs *CaseStatement) statementNode()                {}
func (cs *CaseStatement) expressionNode()               {}
func (cs *CaseStatement) TokenLiteral() string          { return cs.Token.Literal }
func (cs *CaseStatement) TokenLiteralNode() token.Token { return cs.Token }

// String returns a string representation of the CaseStatement.
func (cs *CaseStatement) String() string {
	var sb strings.Builder
	sb.WriteString("case ")
	sb.WriteString(cs.Value.String())
	sb.WriteString(" in\n")
	for _, clause := range cs.Clauses {
		sb.WriteString(clause.String())
	}
	sb.WriteString("esac")
	return sb.String()
}

// CaseClause represents a branch in a case statement.
type CaseClause struct {
	Token    token.Token // The pattern token
	Patterns []Expression
	Body     *BlockStatement
}

func (cc *CaseClause) statementNode()                {}
func (cc *CaseClause) expressionNode()               {}
func (cc *CaseClause) TokenLiteral() string          { return cc.Token.Literal }
func (cc *CaseClause) TokenLiteralNode() token.Token { return cc.Token }

func (cc *CaseClause) String() string {
	var sb strings.Builder
	sb.WriteString("  ")
	for i, p := range cc.Patterns {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(p.String())
	}
	sb.WriteString(") ")
	sb.WriteString(cc.Body.String())
	sb.WriteString("\n")
	return sb.String()
}

// WalkFn is a function that will be called for each node in the AST.
// It returns true if the children of the node should be visited.
type WalkFn func(node Node) bool

// Walk traverses the AST starting from the given node.
// It calls the walkFn for each node.
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
	case *IfStatement:
		Walk(n.Condition, f)
		Walk(n.Consequence, f)
		Walk(n.Alternative, f)
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
		Walk(n.Name, f)
		for _, p := range n.Params {
			Walk(p, f)
		}
		Walk(n.Body, f)
	case *CallExpression:
		Walk(n.Function, f)
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}
	case *IndexExpression:
		Walk(n.Left, f)
		Walk(n.Index, f)
	case *ArrayAccess:
		Walk(n.Index, f)
	case *BracketExpression:
		for _, el := range n.Elements {
			Walk(el, f)
		}
	case *DoubleBracketExpression:
		for _, el := range n.Elements {
			Walk(el, f)
		}
	case *CommandSubstitution:
		Walk(n.Command, f)
	case *Shebang:
		// No nested AST nodes for Shebang
	case *DollarParenExpression:
		Walk(n.Command, f)
	case *ProcessSubstitution:
		Walk(n.Command, f)
	case *Subshell:
		Walk(n.Command, f)
	case *SimpleCommand:
		Walk(n.Name, f)
		for _, arg := range n.Arguments {
			Walk(arg, f)
		}
	case *SelectStatement:
		Walk(n.Name, f)
		for _, item := range n.Items {
			Walk(item, f)
		}
		Walk(n.Body, f)
	case *CoprocStatement:
		Walk(n.Command, f)
	case *DeclarationStatement:
		for _, assign := range n.Assignments {
			Walk(assign.Name, f)
			if assign.Value != nil {
				Walk(assign.Value, f)
			}
		}
	case *ArithmeticCommand:
		Walk(n.Expression, f)
	case *Redirection:
		Walk(n.Left, f)
		Walk(n.Right, f)
	case *ConcatenatedExpression:
		for _, part := range n.Parts {
			Walk(part, f)
		}
	case *CaseStatement:
		Walk(n.Value, f)
		for _, clause := range n.Clauses {
			Walk(clause, f)
		}
	case *CaseClause:
		Walk(n.Body, f)
	case *Identifier:
		// No nested AST nodes for Identifier
	case *IntegerLiteral:
		// No nested AST nodes for IntegerLiteral
	case *Boolean:
		// No nested AST nodes for Boolean
	case *PrefixExpression:
		Walk(n.Right, f)
	case *PostfixExpression:
		Walk(n.Left, f)
	case *InfixExpression:
		Walk(n.Left, f)
		Walk(n.Right, f)
	case *InvalidArrayAccess:
		// No nested AST nodes for InvalidArrayAccess
	case *ArrayLiteral:
		for _, el := range n.Elements {
			Walk(el, f)
		}
	}
}

// SelectStatement represents a select loop.
type SelectStatement struct {
	Token token.Token // The 'select' token
	Name  *Identifier
	Items []Expression
	Body  *BlockStatement
}

func (ss *SelectStatement) statementNode()                {}
func (ss *SelectStatement) expressionNode()               {}
func (ss *SelectStatement) TokenLiteral() string          { return ss.Token.Literal }
func (ss *SelectStatement) TokenLiteralNode() token.Token { return ss.Token }

func (ss *SelectStatement) String() string {
	var sb strings.Builder
	sb.WriteString("select ")
	sb.WriteString(ss.Name.String())
	if len(ss.Items) > 0 {
		sb.WriteString(" in")
		for _, item := range ss.Items {
			sb.WriteString(" ")
			sb.WriteString(item.String())
		}
	}
	sb.WriteString("; do ")
	sb.WriteString(ss.Body.String())
	sb.WriteString("done")
	return sb.String()
}

// CoprocStatement represents a coproc statement.
type CoprocStatement struct {
	Token   token.Token // The 'coproc' token
	Name    string
	Command Statement
}

func (cs *CoprocStatement) statementNode()                {}
func (cs *CoprocStatement) expressionNode()               {}
func (cs *CoprocStatement) TokenLiteral() string          { return cs.Token.Literal }
func (cs *CoprocStatement) TokenLiteralNode() token.Token { return cs.Token }

func (cs *CoprocStatement) String() string {
	var sb strings.Builder
	sb.WriteString("coproc ")
	if cs.Name != "" {
		sb.WriteString(cs.Name + " ")
	}
	sb.WriteString(cs.Command.String())
	return sb.String()
}

// DeclarationAssignment represents an assignment in a declaration.
type DeclarationAssignment struct {
	Name     *Identifier
	IsAppend bool
	Value    Expression
}

func (da *DeclarationAssignment) String() string {
	var sb strings.Builder
	sb.WriteString(da.Name.String())
	if da.Value != nil {
		if da.IsAppend {
			sb.WriteString("+=")
		} else {
			sb.WriteString("=")
		}
		sb.WriteString(da.Value.String())
	}
	return sb.String()
}

// DeclarationStatement represents a typeset/declare statement.
type DeclarationStatement struct {
	Token       token.Token // The 'typeset'/'declare' token
	Command     string      // "typeset", "declare", "local", "export", "readonly"
	Flags       []string
	Assignments []*DeclarationAssignment
}

func (ds *DeclarationStatement) statementNode()                {}
func (ds *DeclarationStatement) expressionNode()               {}
func (ds *DeclarationStatement) TokenLiteral() string          { return ds.Token.Literal }
func (ds *DeclarationStatement) TokenLiteralNode() token.Token { return ds.Token }

func (ds *DeclarationStatement) String() string {
	var sb strings.Builder
	sb.WriteString(ds.Command)
	for _, flag := range ds.Flags {
		sb.WriteString(" " + flag)
	}
	for _, assign := range ds.Assignments {
		sb.WriteString(" ")
		sb.WriteString(assign.String())
	}
	return sb.String()
}

// ArithmeticCommand represents an arithmetic command (( ... )).
type ArithmeticCommand struct {
	Token      token.Token // The '((' token
	Expression Expression
}

func (ac *ArithmeticCommand) statementNode()                {}
func (ac *ArithmeticCommand) expressionNode()               {}
func (ac *ArithmeticCommand) TokenLiteral() string          { return ac.Token.Literal }
func (ac *ArithmeticCommand) TokenLiteralNode() token.Token { return ac.Token }

func (ac *ArithmeticCommand) String() string {
	return "((" + ac.Expression.String() + "))"
}

// Redirection represents a redirection.
type Redirection struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (r *Redirection) expressionNode()               {}
func (r *Redirection) TokenLiteral() string          { return r.Token.Literal }
func (r *Redirection) TokenLiteralNode() token.Token { return r.Token }
func (r *Redirection) String() string {
	return r.Left.String() + " " + r.Operator + " " + r.Right.String()
}

// ProcessSubstitution represents <(...) or >(...).
type ProcessSubstitution struct {
	Token   token.Token
	Command Node
}

func (ps *ProcessSubstitution) expressionNode()               {}
func (ps *ProcessSubstitution) TokenLiteral() string          { return ps.Token.Literal }
func (ps *ProcessSubstitution) TokenLiteralNode() token.Token { return ps.Token }
func (ps *ProcessSubstitution) String() string                { return ps.Token.Literal + ps.Command.String() + ")" }

// Subshell represents ( ... ).
type Subshell struct {
	Token   token.Token
	Command Node
}

func (s *Subshell) statementNode()                {}
func (s *Subshell) expressionNode()               {}
func (s *Subshell) TokenLiteral() string          { return s.Token.Literal }
func (s *Subshell) TokenLiteralNode() token.Token { return s.Token }
func (s *Subshell) String() string                { return "(" + s.Command.String() + ")" }

// FunctionDefinition represents a function definition statement.
type FunctionDefinition struct {
	Token token.Token
	Name  *Identifier
	Body  Statement
}

func (fd *FunctionDefinition) statementNode()                {}
func (fd *FunctionDefinition) expressionNode()               {}
func (fd *FunctionDefinition) TokenLiteral() string          { return fd.Token.Literal }
func (fd *FunctionDefinition) TokenLiteralNode() token.Token { return fd.Token }
func (fd *FunctionDefinition) String() string {
	return "function " + fd.Name.String() + " " + fd.Body.String()
}

var (
	ProgramNode                 = &Program{}
	LetStatementNode            = &LetStatement{}
	ReturnStatementNode         = &ReturnStatement{}
	ExpressionStatementNode     = &ExpressionStatement{}
	BlockStatementNode          = &BlockStatement{}
	IfStatementNode             = &IfStatement{}
	ForLoopStatementNode        = &ForLoopStatement{}
	WhileLoopStatementNode      = &WhileLoopStatement{}
	IdentifierNode              = &Identifier{}
	IntegerLiteralNode          = &IntegerLiteral{}
	BooleanNode                 = &Boolean{}
	PrefixExpressionNode        = &PrefixExpression{}
	PostfixExpressionNode       = &PostfixExpression{}
	InfixExpressionNode         = &InfixExpression{}
	CallExpressionNode          = &CallExpression{}
	IndexExpressionNode         = &IndexExpression{}
	ArrayAccessNode             = &ArrayAccess{}
	BracketExpressionNode       = &BracketExpression{}
	DoubleBracketExpressionNode = &DoubleBracketExpression{}
	CommandSubstitutionNode     = &CommandSubstitution{}
	DollarParenExpressionNode   = &DollarParenExpression{}
	SimpleCommandNode           = &SimpleCommand{}
	ConcatenatedExpressionNode  = &ConcatenatedExpression{}
	InvalidArrayAccessNode      = &InvalidArrayAccess{}
	ArrayLiteralNode            = &ArrayLiteral{}
	StringLiteralNode           = &StringLiteral{}
	GroupedExpressionNode       = &GroupedExpression{}
	SelectStatementNode         = &SelectStatement{}
	CoprocStatementNode         = &CoprocStatement{}
	DeclarationStatementNode    = &DeclarationStatement{}
	ArithmeticCommandNode       = &ArithmeticCommand{}
	RedirectionNode             = &Redirection{}
	ProcessSubstitutionNode     = &ProcessSubstitution{}
	SubshellNode                = &Subshell{}
	FunctionDefinitionNode      = &FunctionDefinition{}
	FunctionLiteralNode         = &FunctionLiteral{}
	CaseStatementNode           = &CaseStatement{}
	ShebangNode                 = &Shebang{}
)
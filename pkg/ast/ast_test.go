package ast

import (
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/token"
)

func TestWalk(t *testing.T) {
	tests := []struct {
		name          string
		node          Node
		expectedCount int
	}{
		{
			"Let Statement",
			&Program{
				Statements: []Statement{
					&LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
						Value: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			4, // Program, LetStatement, Identifier, IntegerLiteral
		},
		{
			"If Statement",
			&IfStatement{
				Condition: &InfixExpression{
					Left:     &IntegerLiteral{Value: 1},
					Operator: "<",
					Right:    &IntegerLiteral{Value: 2},
				},
				Consequence: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: &Identifier{Value: "a"},
						},
					},
				},
			},
			7, // If, Infix, Int, Int, Block, ExprStmt, Ident
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nodes []Node
			Walk(tt.node, func(node Node) bool {
				nodes = append(nodes, node)
				return true
			})

			if len(nodes) != tt.expectedCount {
				t.Errorf("Walk did not visit all nodes. got=%d, want=%d", len(nodes), tt.expectedCount)
			}
		})
	}
}

func TestProgram_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	expected := "let myVar = anotherVar;"
	if program.String() != expected {
		t.Errorf("program.String() wrong. expected=%q, got=%q", expected, program.String())
	}
}

// helper to verify all Node interface methods work on a given node.
func verifyNode(t *testing.T, name string, n Node) {
	t.Helper()
	_ = n.TokenLiteral()
	_ = n.TokenLiteralNode()
	_ = n.String()
}

func TestAllNodeTypes_InterfaceMethods(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "test"}
	ident := &Identifier{Token: tok, Value: "x"}
	intLit := &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "42"}, Value: 42}
	block := &BlockStatement{Token: tok, Statements: []Statement{}}

	nodes := map[string]Node{
		"Program":              &Program{Statements: []Statement{}},
		"LetStatement":         &LetStatement{Token: tok, Name: ident, Value: intLit},
		"ReturnStatement":      &ReturnStatement{Token: tok, ReturnValue: intLit},
		"ReturnStatement_nil":  &ReturnStatement{Token: tok, ReturnValue: nil},
		"ExpressionStatement":  &ExpressionStatement{Token: tok, Expression: ident},
		"Identifier":           ident,
		"IntegerLiteral":       intLit,
		"Boolean":              &Boolean{Token: tok, Value: true},
		"PrefixExpression":     &PrefixExpression{Token: tok, Operator: "!", Right: intLit},
		"PrefixExpression_nil": &PrefixExpression{Token: tok, Operator: "-", Right: nil},
		"PostfixExpression":    &PostfixExpression{Token: tok, Operator: "++", Left: ident},
		"InfixExpression":      &InfixExpression{Token: tok, Left: intLit, Operator: "+", Right: intLit},
		"BlockStatement":       block,
		"IfStatement":          &IfStatement{Token: tok, Condition: ident, Consequence: block},
		"IfStatement_alt":      &IfStatement{Token: tok, Condition: ident, Consequence: block, Alternative: block},
		"ForLoopStatement_named": &ForLoopStatement{
			Token: token.Token{Type: token.FOR, Literal: "for"},
			Name:  ident,
			Items: []Expression{ident},
			Body:  block,
		},
		"ForLoopStatement_c": &ForLoopStatement{
			Token:     token.Token{Type: token.FOR, Literal: "for"},
			Init:      ident,
			Condition: ident,
			Post:      ident,
			Body:      block,
		},
		"ForLoopStatement_empty": &ForLoopStatement{
			Token: token.Token{Type: token.FOR, Literal: "for"},
			Body:  block,
		},
		"WhileLoopStatement": &WhileLoopStatement{
			Token:     token.Token{Type: token.WHILE, Literal: "while"},
			Condition: ident,
			Body:      block,
		},
		"FunctionLiteral": &FunctionLiteral{
			Token:  tok,
			Name:   ident,
			Params: []*Identifier{ident},
			Body:   block,
		},
		"CallExpression": &CallExpression{
			Token:     tok,
			Function:  ident,
			Arguments: []Expression{intLit},
		},
		"IndexExpression": &IndexExpression{Token: tok, Left: ident, Index: intLit},
		"BracketExpression": &BracketExpression{
			Token:    tok,
			Elements: []Expression{ident},
		},
		"DoubleBracketExpression": &DoubleBracketExpression{
			Token:    tok,
			Elements: []Expression{ident},
		},
		"StringLiteral":       &StringLiteral{Token: tok, Value: "hello"},
		"GroupedExpression":   &GroupedExpression{Token: tok, Expression: ident},
		"ArrayAccess":         &ArrayAccess{Token: tok, Left: ident, Index: intLit},
		"ArrayAccess_nil":     &ArrayAccess{Token: tok, Left: nil, Index: nil},
		"CommandSubstitution": &CommandSubstitution{Token: tok, Command: ident},
		"InvalidArrayAccess":  &InvalidArrayAccess{Token: tok, Left: ident, Index: intLit},
		"ArrayLiteral": &ArrayLiteral{
			Token:    tok,
			Elements: []Expression{ident, intLit},
		},
		"Shebang":               &Shebang{Token: tok, Path: "#!/bin/zsh"},
		"DollarParenExpression": &DollarParenExpression{Token: tok, Command: ident},
		"SimpleCommand": &SimpleCommand{
			Token:     tok,
			Name:      ident,
			Arguments: []Expression{ident},
		},
		"ConcatenatedExpression": &ConcatenatedExpression{
			Token: tok,
			Parts: []Expression{ident, intLit},
		},
		"CaseStatement": &CaseStatement{
			Token:   tok,
			Value:   ident,
			Clauses: []*CaseClause{},
		},
		"CaseClause": &CaseClause{
			Token:    tok,
			Patterns: []Expression{ident, intLit},
			Body:     block,
		},
		"SelectStatement": &SelectStatement{
			Token: tok,
			Name:  ident,
			Items: []Expression{ident},
			Body:  block,
		},
		"SelectStatement_noitems": &SelectStatement{
			Token: tok,
			Name:  ident,
			Items: []Expression{},
			Body:  block,
		},
		"CoprocStatement": &CoprocStatement{
			Token:   tok,
			Name:    "myproc",
			Command: &ExpressionStatement{Token: tok, Expression: ident},
		},
		"CoprocStatement_noname": &CoprocStatement{
			Token:   tok,
			Command: &ExpressionStatement{Token: tok, Expression: ident},
		},
		"DeclarationStatement": &DeclarationStatement{
			Token:   tok,
			Command: "declare",
			Flags:   []string{"-a", "-g"},
			Assignments: []*DeclarationAssignment{
				{Name: ident, Value: intLit, IsAppend: false},
				{Name: ident, Value: intLit, IsAppend: true},
				{Name: ident, Value: nil},
			},
		},
		"ArithmeticCommand": &ArithmeticCommand{
			Token:      tok,
			Expression: intLit,
		},
		"Redirection": &Redirection{
			Token:    tok,
			Operator: ">",
			Left:     ident,
			Right:    ident,
		},
		"ProcessSubstitution": &ProcessSubstitution{
			Token:   token.Token{Type: token.LT_LPAREN, Literal: "<("},
			Command: ident,
		},
		"Subshell": &Subshell{
			Token:   tok,
			Command: ident,
		},
		"FunctionDefinition": &FunctionDefinition{
			Token: tok,
			Name:  ident,
			Body:  block,
		},
	}

	for name, n := range nodes {
		t.Run(name, func(t *testing.T) {
			verifyNode(t, name, n)
		})
	}
}

func TestWalk_AllNodeTypes(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "t"}
	ident := &Identifier{Token: tok, Value: "v"}
	intLit := &IntegerLiteral{Token: tok, Value: 1}
	block := &BlockStatement{Token: tok, Statements: []Statement{}}

	// Walk each node type to ensure all switch cases are covered
	walkCases := []Node{
		// Walk nil
		nil,
		&Program{Statements: []Statement{}},
		&LetStatement{Token: tok, Name: ident, Value: intLit},
		&ReturnStatement{Token: tok, ReturnValue: intLit},
		&ExpressionStatement{Token: tok, Expression: ident},
		&BlockStatement{Token: tok, Statements: []Statement{
			&ExpressionStatement{Token: tok, Expression: ident},
		}},
		&IfStatement{Token: tok, Condition: ident, Consequence: block, Alternative: block},
		&ForLoopStatement{Token: tok, Name: ident, Items: []Expression{ident}, Body: block, Init: ident, Condition: ident, Post: ident},
		&WhileLoopStatement{Token: tok, Condition: ident, Body: block},
		&FunctionLiteral{Token: tok, Name: ident, Params: []*Identifier{ident}, Body: block},
		&CallExpression{Token: tok, Function: ident, Arguments: []Expression{intLit}},
		&IndexExpression{Token: tok, Left: ident, Index: intLit},
		&ArrayAccess{Token: tok, Index: intLit},
		&BracketExpression{Token: tok, Elements: []Expression{ident}},
		&DoubleBracketExpression{Token: tok, Elements: []Expression{ident}},
		&CommandSubstitution{Token: tok, Command: ident},
		&Shebang{Token: tok, Path: "#!/bin/zsh"},
		&DollarParenExpression{Token: tok, Command: ident},
		&ProcessSubstitution{Token: tok, Command: ident},
		&Subshell{Token: tok, Command: ident},
		&SimpleCommand{Token: tok, Name: ident, Arguments: []Expression{intLit}},
		&SelectStatement{Token: tok, Name: ident, Items: []Expression{ident}, Body: block},
		&CoprocStatement{Token: tok, Command: &ExpressionStatement{Token: tok, Expression: ident}},
		&DeclarationStatement{Token: tok, Assignments: []*DeclarationAssignment{
			{Name: ident, Value: intLit},
		}},
		&ArithmeticCommand{Token: tok, Expression: intLit},
		&Redirection{Token: tok, Left: ident, Right: ident},
		&ConcatenatedExpression{Token: tok, Parts: []Expression{ident}},
		&CaseStatement{Token: tok, Value: ident, Clauses: []*CaseClause{
			{Token: tok, Body: block},
		}},
		&CaseClause{Token: tok, Body: block},
		ident,
		intLit,
		&Boolean{Token: tok, Value: false},
		&PrefixExpression{Token: tok, Operator: "-", Right: intLit},
		&PostfixExpression{Token: tok, Operator: "++", Left: ident},
		&InfixExpression{Token: tok, Left: intLit, Operator: "+", Right: intLit},
		&InvalidArrayAccess{Token: tok, Left: ident, Index: intLit},
		&ArrayLiteral{Token: tok, Elements: []Expression{ident}},
		&StringLiteral{Token: tok, Value: "s"},
		&GroupedExpression{Token: tok, Expression: ident},
		&FunctionDefinition{Token: tok, Name: ident, Body: block},
	}

	for _, n := range walkCases {
		Walk(n, func(node Node) bool { return true })
	}

	// Test Walk with walkFn returning false (stop traversal)
	count := 0
	Walk(&Program{
		Statements: []Statement{
			&ExpressionStatement{Token: tok, Expression: ident},
		},
	}, func(node Node) bool {
		count++
		return false // stop after first node
	})
	if count != 1 {
		t.Errorf("Walk with stop: expected 1 node visited, got %d", count)
	}
}

func TestCaseStatement_String(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "case"}
	ident := &Identifier{Token: tok, Value: "x"}
	block := &BlockStatement{Token: tok, Statements: []Statement{}}

	cs := &CaseStatement{
		Token: tok,
		Value: ident,
		Clauses: []*CaseClause{
			{
				Token:    tok,
				Patterns: []Expression{ident},
				Body:     block,
			},
		},
	}

	str := cs.String()
	if !strings.Contains(str, "case") {
		t.Error("CaseStatement.String() missing 'case'")
	}
	if !strings.Contains(str, "esac") {
		t.Error("CaseStatement.String() missing 'esac'")
	}
}

func TestSelectStatement_String(t *testing.T) {
	tok := token.Token{Type: token.SELECT, Literal: "select"}
	ident := &Identifier{Token: tok, Value: "item"}
	block := &BlockStatement{Token: tok, Statements: []Statement{}}

	ss := &SelectStatement{
		Token: tok,
		Name:  ident,
		Items: []Expression{ident},
		Body:  block,
	}

	str := ss.String()
	if !strings.Contains(str, "select") {
		t.Error("SelectStatement.String() missing 'select'")
	}
	if !strings.Contains(str, "done") {
		t.Error("SelectStatement.String() missing 'done'")
	}
}

// TestMarkerMethods explicitly calls statementNode() and expressionNode() on
// all AST types to cover those empty interface-conformance markers.
func TestMarkerMethods(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "test"}
	ident := &Identifier{Token: tok, Value: "x"}
	intLit := &IntegerLiteral{Token: tok, Value: 42}
	block := &BlockStatement{Token: tok, Statements: []Statement{}}

	// Statement + Expression markers for types that implement both
	type stmtExpr interface {
		statementNode()
		expressionNode()
	}
	bothNodes := []stmtExpr{
		&LetStatement{Token: tok, Name: ident, Value: intLit},
		&ReturnStatement{Token: tok, ReturnValue: intLit},
		&ExpressionStatement{Token: tok, Expression: ident},
		&BlockStatement{Token: tok, Statements: []Statement{}},
		&IfStatement{Token: tok, Condition: ident, Consequence: block},
		&ForLoopStatement{Token: tok, Body: block},
		&WhileLoopStatement{Token: tok, Condition: ident, Body: block},
		&SimpleCommand{Token: tok, Name: ident},
		&CaseStatement{Token: tok, Value: ident},
		&CaseClause{Token: tok, Patterns: []Expression{ident}, Body: block},
		&SelectStatement{Token: tok, Name: ident, Body: block},
		&CoprocStatement{Token: tok, Command: &ExpressionStatement{Token: tok, Expression: ident}},
		&DeclarationStatement{Token: tok, Command: "declare"},
		&ArithmeticCommand{Token: tok, Expression: intLit},
		&Shebang{Token: tok, Path: "#!/bin/zsh"},
		&Subshell{Token: tok, Command: ident},
		&FunctionDefinition{Token: tok, Name: ident, Body: block},
	}
	for _, n := range bothNodes {
		n.statementNode()
		n.expressionNode()
	}

	// Expression-only markers
	type exprOnly interface {
		expressionNode()
	}
	exprNodes := []exprOnly{
		&Identifier{Token: tok, Value: "x"},
		&IntegerLiteral{Token: tok, Value: 1},
		&Boolean{Token: tok, Value: true},
		&PrefixExpression{Token: tok, Operator: "-", Right: intLit},
		&PostfixExpression{Token: tok, Operator: "++", Left: ident},
		&InfixExpression{Token: tok, Left: intLit, Operator: "+", Right: intLit},
		&FunctionLiteral{Token: tok, Name: ident, Params: []*Identifier{}, Body: block},
		&CallExpression{Token: tok, Function: ident, Arguments: []Expression{}},
		&IndexExpression{Token: tok, Left: ident, Index: intLit},
		&BracketExpression{Token: tok, Elements: []Expression{}},
		&DoubleBracketExpression{Token: tok, Elements: []Expression{}},
		&StringLiteral{Token: tok, Value: "s"},
		&GroupedExpression{Token: tok, Expression: ident},
		&ArrayAccess{Token: tok, Left: ident, Index: intLit},
		&CommandSubstitution{Token: tok, Command: ident},
		&InvalidArrayAccess{Token: tok, Left: ident, Index: intLit},
		&ArrayLiteral{Token: tok, Elements: []Expression{}},
		&DollarParenExpression{Token: tok, Command: ident},
		&ConcatenatedExpression{Token: tok, Parts: []Expression{}},
		&Redirection{Token: tok, Left: ident, Right: ident},
		&ProcessSubstitution{Token: tok, Command: ident},
	}
	for _, n := range exprNodes {
		n.expressionNode()
	}
}

func TestDeclarationAssignment_String(t *testing.T) {
	tok := token.Token{Type: token.IDENT, Literal: "x"}
	ident := &Identifier{Token: tok, Value: "x"}
	intLit := &IntegerLiteral{Token: tok, Value: 5}

	// Regular assignment
	da := &DeclarationAssignment{Name: ident, Value: intLit, IsAppend: false}
	if !strings.Contains(da.String(), "=") {
		t.Error("expected = in assignment string")
	}

	// Append assignment
	da2 := &DeclarationAssignment{Name: ident, Value: intLit, IsAppend: true}
	if !strings.Contains(da2.String(), "+=") {
		t.Error("expected += in append assignment string")
	}

	// No value
	da3 := &DeclarationAssignment{Name: ident}
	s := da3.String()
	if s != "x" {
		t.Errorf("expected just name, got %q", s)
	}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package ast

import "testing"

// TestMarkerMethodsCoverage ensures every node's statementNode() /
// expressionNode() marker is invoked. The methods are no-op interface
// witnesses but must still appear in coverage reports.
func TestMarkerMethodsCoverage(t *testing.T) {
	statements := []Statement{
		&LetStatement{},
		&ReturnStatement{},
		&ExpressionStatement{},
		&BlockStatement{},
		&IfStatement{},
		&ForLoopStatement{},
		&WhileLoopStatement{},
		&Shebang{},
		&SimpleCommand{},
		&CaseStatement{},
		&CaseClause{},
		&SelectStatement{},
		&CoprocStatement{},
		&DeclarationStatement{},
		&ArithmeticCommand{},
		&Subshell{},
		&FunctionDefinition{},
	}
	for _, s := range statements {
		s.statementNode()
	}

	expressions := []Expression{
		&LetStatement{},
		&ReturnStatement{},
		&ExpressionStatement{},
		&Identifier{},
		&IntegerLiteral{},
		&Boolean{},
		&PrefixExpression{},
		&PostfixExpression{},
		&InfixExpression{},
		&BlockStatement{},
		&IfStatement{},
		&ForLoopStatement{},
		&WhileLoopStatement{},
		&FunctionLiteral{},
		&CallExpression{},
		&IndexExpression{},
		&BracketExpression{},
		&DoubleBracketExpression{},
		&StringLiteral{},
		&GroupedExpression{},
		&ArrayAccess{},
		&CommandSubstitution{},
		&InvalidArrayAccess{},
		&ArrayLiteral{},
		&Shebang{},
		&DollarParenExpression{},
		&SimpleCommand{},
		&ConcatenatedExpression{},
		&CaseStatement{},
		&CaseClause{},
		&SelectStatement{},
		&CoprocStatement{},
		&DeclarationStatement{},
		&ArithmeticCommand{},
		&Redirection{},
		&ProcessSubstitution{},
		&Subshell{},
		&FunctionDefinition{},
	}
	for _, e := range expressions {
		e.expressionNode()
	}
}

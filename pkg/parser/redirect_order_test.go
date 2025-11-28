package parser

import (
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/lexer"
)

func TestRedirectOrder(t *testing.T) {
	input := "cmd 2>&1 > out"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	t.Logf("Root: %T %s", stmt.Expression, stmt.Expression.String())

	if infix, ok := stmt.Expression.(*ast.InfixExpression); ok {
		t.Logf("Infix Op: %s", infix.Operator)
		t.Logf("Left: %T", infix.Left)
	}
	if redir, ok := stmt.Expression.(*ast.Redirection); ok {
		t.Logf("Redir Op: %s", redir.Operator)
		t.Logf("Left: %T", redir.Left)
		if subRedir, ok := redir.Left.(*ast.Redirection); ok {
			t.Logf("  SubRedir Op: %s Right: %s", subRedir.Operator, subRedir.Right.String())
		}
	}
}

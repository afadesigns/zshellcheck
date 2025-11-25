package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.RedirectionNode, Kata{
		ID:    "ZC1089",
		Title: "Redirection order matters (`2>&1 > file`)",
		Description: "Redirecting stderr to stdout (`2>&1`) before redirecting stdout to a file (`> file`) " +
			"means stderr goes to the *original* stdout (usually tty), not the file. " +
			"Use `> file 2>&1` or `&> file` to redirect both.",
		Check: checkZC1089,
	})
}

func checkZC1089(node ast.Node) []Violation {
	// Structure: Redirection(Op=">", Left=Redirection(Op=">&", Right=1))
	
	redir, ok := node.(*ast.Redirection)
	if !ok {
		return nil
	}
	
	// Check outer operator > or >>
	if redir.Operator != ">" && redir.Operator != ">>" {
		return nil
	}
	
	// Check inner redirection
	inner, ok := redir.Left.(*ast.Redirection)
	if !ok {
		return nil
	}
	
	// Check inner operator >&
	if inner.Operator != ">&" {
		return nil
	}
	
	// Check inner right is 1
	if intLit, ok := inner.Right.(*ast.IntegerLiteral); ok {
		if intLit.Value == 1 {
			return []Violation{
				{
					KataID:  "ZC1089",
					Message: "Redirection order matters. `2>&1 > file` does not redirect stderr to file. Use `> file 2>&1` instead.",
					Line:    redir.TokenLiteralNode().Line,
					Column:  redir.TokenLiteralNode().Column,
				},
			}
		}
	}
	
	return nil
}

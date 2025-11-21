package katas

import (

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var plumbingCommands = map[string]string{
	"rev-parse":    "git-rev-parse",
	"update-ref":   "git-update-ref",
	"symbolic-ref": "git-symbolic-ref",
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1011",	
		Title: "Use `git` porcelain commands instead of plumbing commands",
		Description: "Plumbing commands in `git` are designed for scripting and can be unstable. " +
			"Porcelain commands are designed for interactive use and are more stable.",
		Check: checkZC1011,
	})
}

func checkZC1011(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if ident, ok := cmd.Name.(*ast.Identifier); ok {
			if ident.Value == "git" {
				for _, arg := range cmd.Arguments {
					val := getArgValue(arg)
					if _, ok := plumbingCommands[val]; ok {
						violations = append(violations, Violation{
							KataID:  "ZC1011",
							Message: "Avoid using `git` plumbing commands in scripts. They are not guaranteed to be stable.",
							Line:    ident.Token.Line,
							Column:  ident.Token.Column,
						})
					}
				}
			}
		}
	}

	return violations
}

func getArgValue(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Value
	case *ast.StringLiteral:
		return e.Value
	case *ast.InfixExpression:
		return getArgValue(e.Left) + e.Operator + getArgValue(e.Right)
	case *ast.ConcatenatedExpression:
		var out string
		for _, p := range e.Parts {
			out += getArgValue(p)
		}
		return out
	case *ast.PrefixExpression:
		// Special handling for prefix expressions in command args to avoid parens
		return e.Operator + getArgValue(e.Right)
	}
	return ""
}

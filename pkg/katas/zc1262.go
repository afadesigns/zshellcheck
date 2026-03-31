package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1262",
		Title:    "Avoid `chmod -R 777` — recursive world-writable is critical",
		Severity: SeverityError,
		Description: "`chmod -R 777` makes every file and directory world-writable and executable. " +
			"Use specific permissions like `755` for directories and `644` for files.",
		Check: checkZC1262,
	})
}

func checkZC1262(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "chmod" {
		return nil
	}

	hasRecursive := false
	has777 := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-R" {
			hasRecursive = true
		}
		if val == "777" {
			has777 = true
		}
	}

	if hasRecursive && has777 {
		return []Violation{{
			KataID: "ZC1262",
			Message: "Never use `chmod -R 777` — it makes everything world-writable. " +
				"Use `find -type d -exec chmod 755` and `find -type f -exec chmod 644` instead.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

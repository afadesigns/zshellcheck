package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1388",
		Title:    "Use Zsh lowercase `$mailpath` array instead of colon-separated `$MAILPATH`",
		Severity: SeverityWarning,
		Description: "Bash uses `$MAILPATH` — a colon-separated string of mail files with " +
			"optional `?message` suffixes. Zsh uses lowercase `$mailpath` as an array (each " +
			"element: `file?message`), which is typed and parseable. Setting the uppercase " +
			"name in Zsh is ignored.",
		Check: checkZC1388,
	})
}

func checkZC1388(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" && ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "$MAILPATH") || strings.Contains(v, "${MAILPATH}") ||
			strings.Contains(v, "MAILPATH=") {
			return []Violation{{
				KataID: "ZC1388",
				Message: "Use Zsh lowercase `$mailpath` (array) instead of Bash uppercase " +
					"`$MAILPATH` (colon-separated string).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

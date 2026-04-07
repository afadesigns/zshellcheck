package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1279",
		Title:    "Use `realpath` instead of `readlink -f` for canonical paths",
		Severity: SeverityInfo,
		Description: "`readlink -f` is not portable across all platforms (notably macOS). " +
			"Use `realpath` which is POSIX-standard and available on modern systems.",
		Check: checkZC1279,
	})
}

func checkZC1279(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "readlink" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			return []Violation{{
				KataID:  "ZC1279",
				Message: "Use `realpath` instead of `readlink -f`. `realpath` is more portable, especially on macOS.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityInfo,
			}}
		}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1151",
		Title:    "Avoid `cat -A` — use `print -v` or od for non-printable characters",
		Severity: SeverityStyle,
		Description: "`cat -A` shows non-printable characters but varies across platforms. " +
			"Use Zsh `print -v` or `od -c` for reliable non-printable character inspection.",
		Check: checkZC1151,
	})
}

func checkZC1151(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cat" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-A" || val == "-v" || val == "-e" {
			return []Violation{{
				KataID: "ZC1151",
				Message: "Avoid `cat " + val + "` for inspecting non-printable characters. " +
					"Use `od -c` or `hexdump -C` for reliable cross-platform output.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

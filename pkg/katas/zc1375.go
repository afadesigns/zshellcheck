package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1375",
		Title:    "Use `[[ -t fd ]]` instead of `tty -s` for tty-check",
		Severity: SeverityStyle,
		Description: "`tty -s` exits 0 if stdin is a terminal. Zsh's `[[ -t 0 ]]` (or `[[ -t 1 ]]` " +
			"for stdout, `[[ -t 2 ]]` for stderr) does the same check without spawning `tty`.",
		Check: checkZC1375,
	})
}

func checkZC1375(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tty" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-s" {
			return []Violation{{
				KataID: "ZC1375",
				Message: "Use `[[ -t 0 ]]` (stdin), `[[ -t 1 ]]` (stdout), or `[[ -t 2 ]]` (stderr) " +
					"instead of `tty -s`. In-shell file-descriptor test, no external process.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

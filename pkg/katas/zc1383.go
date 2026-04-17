package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1383",
		Title:    "Avoid `$TIMEFORMAT` — Zsh uses `$TIMEFMT`",
		Severity: SeverityWarning,
		Description: "Bash's `$TIMEFORMAT` controls the output of the `time` builtin. Zsh uses a " +
			"shorter name, `$TIMEFMT`, for the same purpose. Setting `TIMEFORMAT` in a Zsh script " +
			"has no effect; the Zsh `time` builtin reads `$TIMEFMT`.",
		Check: checkZC1383,
	})
}

func checkZC1383(node ast.Node) []Violation {
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
		if strings.Contains(v, "TIMEFORMAT") {
			return []Violation{{
				KataID: "ZC1383",
				Message: "`$TIMEFORMAT` is Bash-only. Zsh reads `$TIMEFMT` (shorter name) for the " +
					"`time` builtin's output format.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

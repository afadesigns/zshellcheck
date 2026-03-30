package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1169",
		Title:    "Avoid `install` for simple copy+chmod — use `cp` then `chmod`",
		Severity: SeverityStyle,
		Description: "`install` command is less common and may confuse readers. " +
			"For clarity, use separate `cp` and `chmod` commands or `install` only in Makefiles.",
		Check: checkZC1169,
	})
}

func checkZC1169(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "install" {
		return nil
	}

	// Only flag install with -m (mode) flag in scripts
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-m" {
			return []Violation{{
				KataID: "ZC1169",
				Message: "Consider using `cp` + `chmod` instead of `install -m`. " +
					"Separate commands are clearer in shell scripts.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

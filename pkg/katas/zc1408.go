package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1408",
		Title:    "Avoid `$BASH_FUNC_...%%` — Bash-specific exported-function envvar",
		Severity: SeverityError,
		Description: "Bash exports functions into environment variables named `BASH_FUNC_NAME%%`. " +
			"These are consumed only by other Bash shells. Zsh does not recognize the format " +
			"and will neither inherit the function nor clean these envvars.",
		Check: checkZC1408,
	})
}

func checkZC1408(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_FUNC_") {
			return []Violation{{
				KataID: "ZC1408",
				Message: "`BASH_FUNC_*` exported-function envvars are Bash-only. Zsh does not " +
					"consume them; export function definitions via `autoload` + `$FPATH` instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

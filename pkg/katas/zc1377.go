package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1377",
		Title:    "Avoid `$BASH_ALIASES` — use Zsh `$aliases` associative array",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_ALIASES` is an associative array of alias→value mappings. Zsh " +
			"exposes the same information via `$aliases` (also an assoc array). `$BASH_ALIASES` " +
			"is unset in Zsh; reading it yields nothing.",
		Check: checkZC1377,
	})
}

func checkZC1377(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "echo" && ident.Value != "print" && ident.Value != "printf" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.Contains(v, "BASH_ALIASES") {
			return []Violation{{
				KataID: "ZC1377",
				Message: "`$BASH_ALIASES` is Bash-only. In Zsh use `$aliases` (assoc array) — " +
					"same structure, e.g. `print -l ${(kv)aliases}`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

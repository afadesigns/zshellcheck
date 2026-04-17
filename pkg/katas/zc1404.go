package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1404",
		Title:    "Avoid `$BASH_CMDS` — Bash-specific hash-table mirror, use Zsh `$commands`",
		Severity: SeverityWarning,
		Description: "Bash's `$BASH_CMDS` associative array mirrors the hash-table of command " +
			"names→paths. Zsh exposes the same via `$commands` (assoc array from " +
			"`zsh/parameter`). `$BASH_CMDS` is unset in Zsh.",
		Check: checkZC1404,
	})
}

func checkZC1404(node ast.Node) []Violation {
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
		if strings.Contains(v, "BASH_CMDS") {
			return []Violation{{
				KataID: "ZC1404",
				Message: "`$BASH_CMDS` is Bash-only. In Zsh use `$commands` (assoc array, " +
					"names→paths) via `zsh/parameter`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

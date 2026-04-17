package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1384",
		Title:    "Avoid `$EXECIGNORE` — Bash-only; Zsh uses completion-system ignore patterns",
		Severity: SeverityWarning,
		Description: "Bash's `$EXECIGNORE` excludes matching commands from PATH hashing. Zsh does " +
			"not honor this variable; use the compsys tag-based filters " +
			"(`zstyle ':completion:*' ignored-patterns ...`) for a similar effect on completion.",
		Check: checkZC1384,
	})
}

func checkZC1384(node ast.Node) []Violation {
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
		if strings.Contains(v, "EXECIGNORE") {
			return []Violation{{
				KataID: "ZC1384",
				Message: "`$EXECIGNORE` is Bash-only. For completion filtering in Zsh use " +
					"`zstyle ':completion:*' ignored-patterns 'pattern'`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

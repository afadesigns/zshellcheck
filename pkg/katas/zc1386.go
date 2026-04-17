package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1386",
		Title:    "Avoid `$FIGNORE` — Bash-only; Zsh uses compsys tag patterns",
		Severity: SeverityWarning,
		Description: "Bash's `$FIGNORE` hides filenames matching listed suffixes from completion. " +
			"Zsh does not honor this variable; use `zstyle ':completion:*' ignored-patterns '*.o *.pyc'` " +
			"or the file-patterns tag for equivalent filtering.",
		Check: checkZC1386,
	})
}

func checkZC1386(node ast.Node) []Violation {
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
		if strings.Contains(v, "FIGNORE") {
			return []Violation{{
				KataID: "ZC1386",
				Message: "`$FIGNORE` is Bash-only. In Zsh use " +
					"`zstyle ':completion:*' ignored-patterns '*.o *.pyc'` for completion filtering.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

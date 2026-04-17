package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1403",
		Title:    "Setting `$HISTFILESIZE` alone is incomplete in Zsh — pair with `$SAVEHIST`",
		Severity: SeverityWarning,
		Description: "Bash uses `$HISTSIZE` (in-memory) and `$HISTFILESIZE` (on disk). Zsh uses " +
			"`$HISTSIZE` (in-memory) and `$SAVEHIST` (on disk). Setting only `$HISTFILESIZE` in " +
			"Zsh has no effect on disk — `$SAVEHIST` must be set. Mixing both names leaves " +
			"disk-history behavior undefined.",
		Check: checkZC1403,
	})
}

func checkZC1403(node ast.Node) []Violation {
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
		if strings.Contains(v, "HISTFILESIZE") {
			return []Violation{{
				KataID: "ZC1403",
				Message: "`$HISTFILESIZE` is Bash-only. Zsh uses `$SAVEHIST` for on-disk history " +
					"size. Setting `HISTFILESIZE` in Zsh has no effect.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

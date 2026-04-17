package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1380",
		Title:    "Avoid `$HISTIGNORE` — use Zsh `$HISTORY_IGNORE`",
		Severity: SeverityWarning,
		Description: "Bash filters history entries matching `$HISTIGNORE` patterns. Zsh uses a " +
			"parameter named `$HISTORY_IGNORE` (underscore in the middle). Setting `HISTIGNORE` " +
			"in Zsh is a no-op.",
		Check: checkZC1380,
	})
}

func checkZC1380(node ast.Node) []Violation {
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
		if strings.Contains(v, "HISTIGNORE") && !strings.Contains(v, "HISTORY_IGNORE") {
			return []Violation{{
				KataID: "ZC1380",
				Message: "`$HISTIGNORE` is Bash-only. In Zsh use `$HISTORY_IGNORE` (underscored) " +
					"for the same history-pattern filter.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

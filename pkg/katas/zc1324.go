package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1324",
		Title:    "Avoid `$PROMPT_COMMAND` — use `precmd` hook in Zsh",
		Severity: SeverityWarning,
		Description: "`$PROMPT_COMMAND` is a Bash variable that executes a command before " +
			"each prompt. Zsh uses the `precmd` hook function for the same purpose.",
		Check: checkZC1324,
	})
}

func checkZC1324(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$PROMPT_COMMAND" && ident.Value != "PROMPT_COMMAND" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1324",
		Message: "Avoid `$PROMPT_COMMAND` in Zsh — use the `precmd` hook function instead. `PROMPT_COMMAND` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

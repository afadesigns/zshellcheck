package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1320",
		Title:    "Avoid `$BASH_ARGV` — use `$argv` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_ARGV` is a Bash array containing arguments in reverse order. " +
			"Zsh provides `$argv` (or `$@`) for positional parameters.",
		Check: checkZC1320,
	})
}

func checkZC1320(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_ARGV" && ident.Value != "BASH_ARGV" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1320",
		Message: "Avoid `$BASH_ARGV` in Zsh — use `$argv` or `$@` for positional parameters. `BASH_ARGV` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

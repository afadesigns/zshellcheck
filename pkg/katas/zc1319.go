package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1319",
		Title:    "Avoid `$BASH_ARGC` — use `$#` in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_ARGC` is a Bash array tracking argument counts per stack frame. " +
			"Zsh uses `$#` for argument count and `$argv` for the argument array.",
		Check: checkZC1319,
	})
}

func checkZC1319(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_ARGC" && ident.Value != "BASH_ARGC" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1319",
		Message: "Avoid `$BASH_ARGC` in Zsh — use `$#` for argument count. `BASH_ARGC` is Bash-specific.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

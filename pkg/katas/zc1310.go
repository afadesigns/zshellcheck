package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1310",
		Title:    "Avoid `$BASH_EXECUTION_STRING` — not available in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_EXECUTION_STRING` contains the argument to `bash -c`. " +
			"Zsh does not provide this variable. Access the script argument directly.",
		Check: checkZC1310,
	})
}

func checkZC1310(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$BASH_EXECUTION_STRING" && ident.Value != "BASH_EXECUTION_STRING" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1310",
		Message: "Avoid `$BASH_EXECUTION_STRING` in Zsh — it is undefined. Access command arguments directly instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

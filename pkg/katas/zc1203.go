package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1203",
		Title:    "Avoid `netstat` — use `ss` for socket statistics",
		Severity: SeverityInfo,
		Description: "`netstat` is deprecated on modern Linux in favor of `ss` from iproute2. " +
			"`ss` is faster and provides more detailed socket information.",
		Check: checkZC1203,
	})
}

func checkZC1203(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "netstat" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1203",
		Message: "Avoid `netstat` — it is deprecated on modern Linux. " +
			"Use `ss` from iproute2 for faster, more detailed socket statistics.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

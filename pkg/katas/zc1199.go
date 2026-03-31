package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1199",
		Title:    "Avoid `telnet` in scripts — use `curl` or `zsh/net/tcp`",
		Severity: SeverityWarning,
		Description: "`telnet` is interactive and sends data in plain text. " +
			"Use `curl` for HTTP or `zmodload zsh/net/tcp` for port checks in scripts.",
		Check: checkZC1199,
	})
}

func checkZC1199(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "telnet" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1199",
		Message: "Avoid `telnet` in scripts — it is interactive and insecure. " +
			"Use `curl` for HTTP checks or `zmodload zsh/net/tcp` for port testing.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

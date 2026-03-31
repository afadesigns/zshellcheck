package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1200",
		Title:    "Avoid `ftp` — use `sftp` or `curl` for secure transfers",
		Severity: SeverityWarning,
		Description: "`ftp` transmits credentials and data in plain text. " +
			"Use `sftp`, `scp`, or `curl` with HTTPS/SFTP for secure file transfers.",
		Check: checkZC1200,
	})
}

func checkZC1200(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ftp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1200",
		Message: "Avoid `ftp` — it transmits credentials in plain text. " +
			"Use `sftp`, `scp`, or `curl` with HTTPS for secure file transfers.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

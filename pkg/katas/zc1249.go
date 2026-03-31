package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1249",
		Title:    "Use `ssh-keygen -f` to specify key file in scripts",
		Severity: SeverityWarning,
		Description: "`ssh-keygen` without `-f` prompts for a file path interactively. " +
			"Use `-f /path/to/key` and `-N ''` for non-interactive key generation.",
		Check: checkZC1249,
	})
}

func checkZC1249(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keygen" {
		return nil
	}

	hasFile := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" {
			hasFile = true
		}
	}

	if !hasFile {
		return []Violation{{
			KataID: "ZC1249",
			Message: "Use `ssh-keygen -f /path/to/key -N ''` in scripts. " +
				"Without `-f`, ssh-keygen prompts interactively for the file path.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

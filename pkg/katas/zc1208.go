package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1208",
		Title:    "Avoid `visudo` in scripts — use sudoers.d drop-in files",
		Severity: SeverityWarning,
		Description: "`visudo` opens an interactive editor. For programmatic sudoers changes, " +
			"write to `/etc/sudoers.d/` drop-in files with `visudo -c` for validation.",
		Check: checkZC1208,
	})
}

func checkZC1208(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "visudo" {
		return nil
	}

	// visudo -c (check) is fine — it's non-interactive validation
	for _, arg := range cmd.Arguments {
		if arg.String() == "-c" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1208",
		Message: "Avoid `visudo` in scripts — it opens an interactive editor. " +
			"Write to `/etc/sudoers.d/` drop-in files and validate with `visudo -c`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

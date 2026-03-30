package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1129",
		Title: "Use Zsh `stat` module instead of `wc -c` for file size",
		Description: "Zsh's `zstat` (via `zmodload zsh/stat`) provides file size without " +
			"spawning `wc`. Use `zstat +size file` for efficient file size queries.",
		Severity: SeverityStyle,
		Check:    checkZC1129,
	})
}

func checkZC1129(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasCharFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" {
			hasCharFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if !hasCharFlag || !hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1129",
		Message: "Use `zstat +size file` (via `zmodload zsh/stat`) instead of `wc -c file`. " +
			"Avoids reading the entire file for a simple size query.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

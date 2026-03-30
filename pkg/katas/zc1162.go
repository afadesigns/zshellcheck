package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1162",
		Title:    "Use `cp -a` instead of `cp -r` to preserve attributes",
		Severity: SeverityInfo,
		Description: "`cp -r` copies recursively but may not preserve permissions, timestamps, " +
			"or symlinks. Use `cp -a` (archive mode) to preserve all attributes.",
		Check: checkZC1162,
	})
}

func checkZC1162(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cp" {
		return nil
	}

	hasRecursive := false
	hasArchive := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-r" || val == "-R" {
			hasRecursive = true
		}
		if val == "-a" || val == "-rp" || val == "-Rp" {
			hasArchive = true
		}
	}

	if hasRecursive && !hasArchive {
		return []Violation{{
			KataID: "ZC1162",
			Message: "Use `cp -a` instead of `cp -r` to preserve permissions, timestamps, and symlinks. " +
				"Archive mode ensures a faithful copy.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

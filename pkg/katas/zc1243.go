package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1243",
		Title:    "Use `grep -lZ` with `xargs -0` for safe file lists",
		Severity: SeverityWarning,
		Description: "`grep -l` outputs one filename per line, breaking on names with newlines. " +
			"Use `grep -lZ` (null-terminated) paired with `xargs -0` for safe processing.",
		Check: checkZC1243,
	})
}

func checkZC1243(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasListFiles := false
	hasNullTerm := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-l" || val == "-rl" || val == "-lr" {
			hasListFiles = true
		}
		if val == "-Z" || val == "-lZ" || val == "-Zl" {
			hasNullTerm = true
		}
	}

	if hasListFiles && !hasNullTerm {
		return []Violation{{
			KataID: "ZC1243",
			Message: "Use `grep -lZ` instead of `grep -l` for null-terminated file lists. " +
				"Pair with `xargs -0` to safely handle filenames with special characters.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1244",
		Title:    "Consider `mv -n` to prevent overwriting existing files",
		Severity: SeverityInfo,
		Description: "`mv` overwrites existing files without warning by default. " +
			"Use `-n` (no-clobber) to prevent accidental overwrites in scripts.",
		Check: checkZC1244,
	})
}

func checkZC1244(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mv" {
		return nil
	}

	hasSafe := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-n" || val == "-i" || val == "-f" {
			hasSafe = true
		}
	}

	if !hasSafe && len(cmd.Arguments) >= 2 {
		return []Violation{{
			KataID: "ZC1244",
			Message: "Consider `mv -n` to prevent overwriting existing files. " +
				"Without `-n`, `mv` silently overwrites the target.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

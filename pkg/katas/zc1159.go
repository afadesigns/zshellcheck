package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1159",
		Title:    "Avoid `tar` without explicit compression flag",
		Severity: SeverityInfo,
		Description: "Use explicit compression flags (`-z` for gzip, `-j` for bzip2, `-J` for xz) " +
			"instead of relying on `tar` auto-detection for clarity and portability.",
		Check: checkZC1159,
	})
}

func checkZC1159(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tar" {
		return nil
	}

	hasCreate := false
	hasCompression := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-cf" || val == "cf" {
			hasCreate = true
		}
		if val == "-z" || val == "-j" || val == "-J" || val == "--gzip" || val == "--bzip2" || val == "--xz" {
			hasCompression = true
		}
		// Combined flags like czf
		if len(val) > 1 && val[0] != '-' {
			for _, ch := range val {
				if ch == 'c' {
					hasCreate = true
				}
				if ch == 'z' || ch == 'j' || ch == 'J' {
					hasCompression = true
				}
			}
		}
	}

	if hasCreate && !hasCompression {
		return []Violation{{
			KataID: "ZC1159",
			Message: "Specify an explicit compression flag (`-z`, `-j`, `-J`) when creating tar archives. " +
				"Relying on auto-detection reduces clarity and portability.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

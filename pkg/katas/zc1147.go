package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1147",
		Title:    "Avoid `mkdir` without `-p` for nested paths",
		Severity: SeverityInfo,
		Description: "Using `mkdir` without `-p` fails if parent directories don't exist. " +
			"Use `mkdir -p` to create the full path safely.",
		Check: checkZC1147,
	})
}

func checkZC1147(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "mkdir" {
		return nil
	}

	hasParentFlag := false
	hasNestedPath := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-p" {
			hasParentFlag = true
		}
		// Check for paths with multiple slashes (nested)
		if len(val) > 0 && val[0] != '-' {
			slashCount := 0
			for _, ch := range val {
				if ch == '/' {
					slashCount++
				}
			}
			if slashCount >= 2 {
				hasNestedPath = true
			}
		}
	}

	if hasParentFlag || !hasNestedPath {
		return nil
	}

	return []Violation{{
		KataID: "ZC1147",
		Message: "Use `mkdir -p` when creating nested directories. " +
			"Without `-p`, `mkdir` fails if parent directories don't exist.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1196",
		Title:    "Use `grep -F` for fixed string matching",
		Severity: SeverityStyle,
		Description: "When searching for a literal string (no regex metacharacters), " +
			"`grep -F` is faster because it skips regex compilation.",
		Check: checkZC1196,
	})
}

func checkZC1196(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasFixed := false
	hasExtended := false
	pattern := ""

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-F" {
			hasFixed = true
		}
		if val == "-E" || val == "-P" {
			hasExtended = true
		}
		if len(val) > 0 && val[0] != '-' && pattern == "" {
			pattern = strings.Trim(val, "'\"")
		}
	}

	if hasFixed || hasExtended || pattern == "" {
		return nil
	}

	// Check if pattern contains regex metacharacters
	metacharacters := ".[](){}*+?|^$\\"
	hasRegex := false
	for _, ch := range pattern {
		if strings.ContainsRune(metacharacters, ch) {
			hasRegex = true
			break
		}
	}

	if !hasRegex && len(pattern) > 2 {
		return []Violation{{
			KataID: "ZC1196",
			Message: "Use `grep -F` for literal string matching. " +
				"It skips regex compilation and is faster for fixed patterns.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

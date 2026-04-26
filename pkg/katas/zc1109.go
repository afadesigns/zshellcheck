// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1109",
		Title: "Use parameter expansion instead of `cut` for field extraction",
		Description: "For simple field extraction from variables, use Zsh parameter expansion " +
			"like `${var%%:*}` or `${(s.:.)var}` instead of piping through `cut`.",
		Severity: SeverityStyle,
		Check:    checkZC1109,
	})
}

func checkZC1109(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cut" {
		return nil
	}

	// Only flag simple cut with -d and -f flags and no file argument
	hasDelimiter := false
	hasField := false
	hasFileArg := false

	for _, arg := range cmd.Arguments {
		val := strings.Trim(arg.String(), "'\"")
		switch {
		case strings.HasPrefix(val, "-d"), strings.HasPrefix(val, "--delimiter"):
			hasDelimiter = true
		case strings.HasPrefix(val, "-f"), strings.HasPrefix(val, "--fields"):
			hasField = true
		case len(val) > 0 && val[0] != '-':
			hasFileArg = true
		}
	}

	if hasFileArg || !hasDelimiter || !hasField {
		return nil
	}

	return []Violation{{
		KataID: "ZC1109",
		Message: "Use Zsh parameter expansion for field extraction instead of `cut`. " +
			"`${var%%delim*}` or `${(s.delim.)var}` avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

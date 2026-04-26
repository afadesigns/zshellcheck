// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1333",
		Title:    "Avoid `$TIMEFORMAT` — use `$TIMEFMT` in Zsh",
		Severity: SeverityInfo,
		Description: "`$TIMEFORMAT` is the Bash variable for customizing `time` output. " +
			"Zsh uses `$TIMEFMT` for the same purpose, with different format specifiers.",
		Check: checkZC1333,
		Fix:   fixZC1333,
	})
}

// fixZC1333 renames the Bash `$TIMEFORMAT` identifier to the Zsh
// `$TIMEFMT` variable. Format specifiers differ between the two
// shells; the rename preserves the identifier itself but authors
// should still review the format string after conversion.
func fixZC1333(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}
	switch ident.Value {
	case "$TIMEFORMAT":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$TIMEFORMAT"),
			Replace: "$TIMEFMT",
		}}
	case "TIMEFORMAT":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("TIMEFORMAT"),
			Replace: "TIMEFMT",
		}}
	}
	return nil
}

func checkZC1333(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident == nil {
		return nil
	}

	if ident.Value != "$TIMEFORMAT" && ident.Value != "TIMEFORMAT" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1333",
		Message: "Avoid `$TIMEFORMAT` in Zsh — use `$TIMEFMT` instead. Format specifiers differ between Bash and Zsh.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityInfo,
	}}
}

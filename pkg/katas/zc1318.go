package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IdentifierNode, Kata{
		ID:       "ZC1318",
		Title:    "Avoid `$BASH_CMDS` — use `$commands` hash in Zsh",
		Severity: SeverityWarning,
		Description: "`$BASH_CMDS` is a Bash associative array caching command lookups. " +
			"Zsh provides the `$commands` hash for the same purpose, mapping " +
			"command names to their full paths.",
		Check: checkZC1318,
		Fix:   fixZC1318,
	})
}

// fixZC1318 renames the Bash `$BASH_CMDS` identifier to the Zsh
// `$commands` hash.
func fixZC1318(node ast.Node, v Violation, source []byte) []FixEdit {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "$BASH_CMDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("$BASH_CMDS"),
			Replace: "$commands",
		}}
	case "BASH_CMDS":
		return []FixEdit{{
			Line:    v.Line,
			Column:  v.Column,
			Length:  len("BASH_CMDS"),
			Replace: "commands",
		}}
	}
	return nil
}

func checkZC1318(node ast.Node) []Violation {
	ident, ok := node.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "$BASH_CMDS" && ident.Value != "BASH_CMDS" {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1318",
		Message: "Avoid `$BASH_CMDS` in Zsh — use the `$commands` hash for command path lookups instead.",
		Line:    ident.Token.Line,
		Column:  ident.Token.Column,
		Level:   SeverityWarning,
	}}
}

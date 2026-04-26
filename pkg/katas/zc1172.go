// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1172",
		Title:    "Use `read -A` instead of Bash `read -a` for arrays",
		Severity: SeverityInfo,
		Description: "Bash uses `read -a` to read into an array, but Zsh uses `read -A`. " +
			"Using `-a` in Zsh reads into a scalar, not an array.",
		Check: checkZC1172,
		Fix:   fixZC1172,
	})
}

// fixZC1172 swaps the lowercase `-a` flag for the uppercase `-A` Zsh
// equivalent. Single-byte replacement at the argument's column.
// Idempotent: a re-run sees `-A`, not `-a`, so the detector won't
// fire. Defensive byte-match guard refuses to insert unless the
// source at the offset is literally `-a`.
func fixZC1172(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if arg.String() != "-a" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-a" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "-A",
		}}
	}
	return nil
}

func checkZC1172(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "read" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-a" {
			return []Violation{{
				KataID: "ZC1172",
				Message: "Use `read -A` instead of `read -a` in Zsh. " +
					"The `-a` flag is Bash syntax; Zsh uses `-A` to read into arrays.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1260",
		Title:    "Use `git branch -d` instead of `-D` for safe deletion",
		Severity: SeverityWarning,
		Description: "`git branch -D` force-deletes branches even if unmerged. " +
			"Use `-d` which refuses to delete unmerged branches, preventing data loss.",
		Check: checkZC1260,
		Fix:   fixZC1260,
	})
}

// fixZC1260 rewrites `git branch -D` to `git branch -d`. Single-character
// flag swap at the `-D` argument position; surrounding subcommand and
// branch-name arguments stay in place.
func fixZC1260(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}
	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "branch" {
		return nil
	}
	for _, arg := range cmd.Arguments[1:] {
		if arg.String() != "-D" {
			continue
		}
		tok := arg.TokenLiteralNode()
		off := LineColToByteOffset(source, tok.Line, tok.Column)
		if off < 0 || off+2 > len(source) {
			return nil
		}
		if string(source[off:off+2]) != "-D" {
			return nil
		}
		return []FixEdit{{
			Line:    tok.Line,
			Column:  tok.Column,
			Length:  2,
			Replace: "-d",
		}}
	}
	return nil
}

func checkZC1260(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "git" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "branch" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		if arg.String() == "-D" {
			return []Violation{{
				KataID: "ZC1260",
				Message: "Use `git branch -d` instead of `-D`. The lowercase `-d` refuses to " +
					"delete unmerged branches, preventing accidental data loss.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

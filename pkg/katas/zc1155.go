// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1155",
		Title:    "Use `whence -a` instead of `which -a`",
		Severity: SeverityInfo,
		Description: "`which -a` may be an external command on some systems. " +
			"Zsh builtin `whence -a` reliably lists all command locations.",
		Check: checkZC1155,
		Fix:   fixZC1155,
	})
}

// fixZC1155 rewrites the `which` command name to `whence`, leaving the
// `-a` flag and any other arguments in place. Detector already
// guarantees the shape (which + -a anywhere in argv).
func fixZC1155(node ast.Node, v Violation, _ []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("which"),
		Replace: "whence",
	}}
}

func checkZC1155(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "which" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-a" {
			return []Violation{{
				KataID: "ZC1155",
				Message: "Use `whence -a` instead of `which -a`. " +
					"Zsh `whence` is a reliable builtin for listing all command locations.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityInfo,
			}}
		}
	}

	return nil
}

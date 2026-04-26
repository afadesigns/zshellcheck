// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1161",
		Title:    "Avoid `openssl` for simple hashing — use Zsh modules",
		Severity: SeverityStyle,
		Description: "For simple SHA/MD5 hashing, Zsh provides `zmodload zsh/sha256` and " +
			"`zmodload zsh/md5`. Avoid spawning `openssl` or `sha256sum` for basic hash operations.",
		Check: checkZC1161,
	})
}

func checkZC1161(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "sha256sum" && name != "sha1sum" && name != "md5sum" && name != "md5" {
		return nil
	}

	// Only flag when used without file arguments (pipeline usage)
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] != '-' {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1161",
		Message: "Consider `zmodload zsh/sha256` or `zmodload zsh/md5` for hash operations. " +
			"Zsh modules avoid spawning external hashing processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

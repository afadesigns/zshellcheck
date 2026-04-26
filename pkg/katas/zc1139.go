// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1139",
		Title: "Avoid `source` with URL — use local files",
		Description: "Sourcing scripts from URLs (curl | source) is a security risk. " +
			"Download, verify, then source local files.",
		Severity: SeverityWarning,
		Check:    checkZC1139,
	})
}

func checkZC1139(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "source" && ident.Value != "." {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 8 && (val[:8] == "https://" || val[:7] == "http://") {
			return []Violation{{
				KataID: "ZC1139",
				Message: "Avoid sourcing scripts from URLs. Download, verify integrity, " +
					"then source from a local path to prevent supply-chain attacks.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

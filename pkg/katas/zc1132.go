// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1132",
		Title: "Use Zsh pattern extraction instead of `grep -o`",
		Description: "For extracting matching parts from variables, use Zsh `${(M)var:#pattern}` " +
			"or `${match[1]}` with `=~` instead of piping through `grep -o`.",
		Severity: SeverityStyle,
		Check:    checkZC1132,
	})
}

func checkZC1132(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasOnlyMatching := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-o" {
				hasOnlyMatching = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasOnlyMatching || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1132",
		Message: "Use Zsh pattern extraction `${(M)var:#pattern}` or `[[ $var =~ regex ]] && echo $match[1]` " +
			"instead of piping through `grep -o`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

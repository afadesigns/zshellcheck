// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1125",
		Title: "Avoid `echo | grep` for string matching",
		Description: "Using `echo $var | grep pattern` spawns two unnecessary processes. " +
			"Use Zsh `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` for string matching.",
		Severity: SeverityStyle,
		Check:    checkZC1125,
	})
}

func checkZC1125(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	// Only flag grep with -q (quiet) and no file argument
	// grep -q is typically used for string matching in conditionals
	hasQuiet := false
	hasFile := false
	patternSeen := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			if val == "-q" {
				hasQuiet = true
			}
		} else {
			if patternSeen {
				hasFile = true
				break
			}
			patternSeen = true
		}
	}

	if !hasQuiet || hasFile {
		return nil
	}

	return []Violation{{
		KataID: "ZC1125",
		Message: "Use `[[ $var =~ pattern ]]` or `[[ $var == *pattern* ]]` instead of piping " +
			"through `grep -q`. Zsh pattern matching avoids spawning external processes.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

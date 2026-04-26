// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1111",
		Title: "Avoid `xargs` for simple command invocation",
		Description: "Zsh can iterate arrays directly with `for` loops or use `${(f)...}` to split " +
			"command output by newlines. Avoid `xargs` when processing lines one at a time.",
		Severity: SeverityStyle,
		Check:    checkZC1111,
	})
}

func checkZC1111(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "xargs" {
		return nil
	}

	// Only flag simple xargs without complex flags
	// -0, -P (parallel), -I (replace string), -L are complex uses — skip them
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 1 && val[0] == '-' {
			switch {
			case val == "-0", val == "--null":
				return nil
			case val == "-P", val == "--max-procs":
				return nil
			case val == "-I", val == "--replace":
				return nil
			case val == "-L", val == "--max-lines":
				return nil
			case val == "-p", val == "--interactive":
				return nil
			}
		}
	}

	return []Violation{{
		KataID: "ZC1111",
		Message: "Consider using Zsh array iteration instead of `xargs`. " +
			"`for item in ${(f)$(cmd)}` splits output by newlines without spawning xargs.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

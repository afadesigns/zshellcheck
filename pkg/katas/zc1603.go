// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1603",
		Title:    "Warn on `gdb -p PID` / `ltrace -p PID` — live attach reads target memory",
		Severity: SeverityWarning,
		Description: "`gdb -p PID` and `ltrace -p PID` attach via ptrace and hand the caller " +
			"full read / write access to the target process: registers, heap, stack, open file " +
			"descriptors, and every environment variable. Credentials in `$AWS_SECRET_ACCESS_" +
			"KEY`, session tokens on the stack, TLS keys in memory — all readable. A root-run " +
			"script that attaches to another user's process extracts everything that user has. " +
			"Keep production scripts out of the debugger; if post-mortem diagnostics are " +
			"needed, use `coredumpctl` against a captured core file instead.",
		Check: checkZC1603,
	})
}

func checkZC1603(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gdb" && ident.Value != "ltrace" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-p" {
			return []Violation{{
				KataID: "ZC1603",
				Message: "`" + ident.Value + " -p PID` attaches via ptrace — memory, " +
					"registers, env, and stack of the target are readable. Use " +
					"`coredumpctl` on a captured core, not a live attach from a script.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

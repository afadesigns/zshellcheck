// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1746",
		Title:    "Error on `sysctl -w kernel.randomize_va_space=0|1` — weakens or disables ASLR",
		Severity: SeverityError,
		Description: "`kernel.randomize_va_space` controls Address Space Layout Randomization. " +
			"Value `2` (default) randomizes stack, heap, VDSO, and mmap regions; value `1` " +
			"omits the heap; value `0` disables ASLR entirely, making every memory layout " +
			"deterministic. Exploits that rely on absolute addresses — stack overflows, ROP " +
			"chains, kernel gadgets — become one-shot instead of brute-forceable. Never " +
			"lower this below `2` outside a sandboxed kernel-debug context.",
		Check: checkZC1746,
	})
}

func checkZC1746(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sysctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "kernel.randomize_va_space=0" || v == "kernel.randomize_va_space=1" {
			return []Violation{{
				KataID: "ZC1746",
				Message: "`sysctl " + v + "` weakens ASLR — absolute-address exploits " +
					"become deterministic (stack overflows, ROP). Keep " +
					"`kernel.randomize_va_space=2` outside a sandboxed debug context.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

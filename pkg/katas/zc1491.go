// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1491",
		Title:    "Warn on `export LD_PRELOAD=...` / `LD_LIBRARY_PATH=...` — library injection",
		Severity: SeverityWarning,
		Description: "Setting `LD_PRELOAD` in a script forces every subsequent dynamically-linked " +
			"command to load the specified shared object first, a classic post-compromise " +
			"privesc and persistence technique. Setting `LD_LIBRARY_PATH` to a writable path is " +
			"a gentler variant of the same class. Legitimate uses exist (perf profiling, " +
			"asan instrumentation) but should be scoped to a single invocation and the path " +
			"pinned to a read-only location.",
		Check: checkZC1491,
	})
}

func checkZC1491(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "export" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if strings.HasPrefix(v, "LD_PRELOAD=") && len(v) > len("LD_PRELOAD=") {
			return zc1491Violation(cmd, "LD_PRELOAD")
		}
		if strings.HasPrefix(v, "LD_LIBRARY_PATH=") && len(v) > len("LD_LIBRARY_PATH=") {
			return zc1491Violation(cmd, "LD_LIBRARY_PATH")
		}
		if strings.HasPrefix(v, "LD_AUDIT=") && len(v) > len("LD_AUDIT=") {
			return zc1491Violation(cmd, "LD_AUDIT")
		}
	}
	return nil
}

func zc1491Violation(cmd *ast.SimpleCommand, varName string) []Violation {
	return []Violation{{
		KataID: "ZC1491",
		Message: "`export " + varName + "=...` forces every subsequent binary to load a custom " +
			"library — classic privesc/persistence. Scope to a single invocation if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

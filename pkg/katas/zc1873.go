// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1873",
		Title:    "Warn on `setopt ERR_RETURN` — functions silently bail out on the first non-zero exit",
		Severity: SeverityWarning,
		Description: "`ERR_RETURN` is the function-scoped cousin of `ERR_EXIT` and is off by " +
			"default in Zsh. Turning it on script-wide makes every function `return` at " +
			"the first command whose status is non-zero, which in practice means helpers " +
			"that deliberately probe the environment (`test -f /some/file`, `grep -q " +
			"PATTERN`, `id -u user`) will bail before they reach the branch that was meant " +
			"to run when the probe failed. Callers see a success-or-nothing return and " +
			"no stderr. Keep the option off at script level; inside one function that " +
			"really wants fail-fast semantics, scope with `setopt LOCAL_OPTIONS; setopt " +
			"ERR_RETURN` so the behaviour cannot leak to the rest of the shell.",
		Check: checkZC1873,
	})
}

func checkZC1873(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1873IsErrReturn(arg.String()) {
				return zc1873Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOERRRETURN" {
				return zc1873Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1873IsErrReturn(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "ERRRETURN"
}

func zc1873Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1873",
		Message: "`" + where + "` returns from every function on first non-zero " +
			"exit — probing helpers (`test -f`, `grep -q`) bail before the " +
			"fallback branch. Scope inside a `LOCAL_OPTIONS` function if " +
			"needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

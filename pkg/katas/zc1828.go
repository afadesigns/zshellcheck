// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1828",
		Title:    "Warn on `gcore PID` / `strace -p PID` — live ptrace attach dumps target memory",
		Severity: SeverityWarning,
		Description: "`gcore PID` writes a core dump of the running process to disk; `strace -p " +
			"PID` streams every syscall the process makes. Both attach via ptrace and expose " +
			"the target's memory, stack, environment variables, and argument buffers — " +
			"credentials, TLS session keys, and `$AWS_SECRET_ACCESS_KEY`-style env vars are " +
			"all readable. A root-run script that attaches to another user's process extracts " +
			"whatever that user has. Keep production scripts off ptrace; reach for " +
			"`coredumpctl` with a captured core or vendor-specific `perf` counters when you " +
			"only need syscall statistics.",
		Check: checkZC1828,
	})
}

func checkZC1828(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "gcore":
		if len(cmd.Arguments) > 0 && !zc1828IsHelp(cmd.Arguments) {
			return zc1828Hit(cmd, "gcore")
		}
	case "strace":
		for _, arg := range cmd.Arguments {
			if arg.String() == "-p" {
				return zc1828Hit(cmd, ident.Value+" -p")
			}
		}
	}
	return nil
}

func zc1828IsHelp(args []ast.Expression) bool {
	for _, a := range args {
		v := a.String()
		if v == "-h" || v == "--help" || v == "-?" || v == "--version" {
			return true
		}
	}
	return false
}

func zc1828Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1828",
		Message: "`" + what + " PID` attaches via ptrace — target memory, env, and " +
			"syscall args are exposed. Production scripts should not run ptrace; " +
			"use `coredumpctl` on a captured core instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

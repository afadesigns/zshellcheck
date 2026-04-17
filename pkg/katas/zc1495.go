package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1495",
		Title:    "Warn on `ulimit -c unlimited` — enables core dumps from setuid binaries",
		Severity: SeverityWarning,
		Description: "`ulimit -c unlimited` enables unbounded core dumps for the current shell " +
			"and its children. On a system with `fs.suid_dumpable=2` and a world-readable " +
			"coredump directory, a setuid process that segfaults leaks its memory into a file " +
			"any user can read — Dirty COW-class keys, TLS session material, kerberos tickets. " +
			"Leave core dumps at the distro default (usually 0) and use systemd-coredump with " +
			"access controls if you genuinely need post-mortems.",
		Check: checkZC1495,
	})
}

func checkZC1495(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ulimit" {
		return nil
	}

	var coreFlag bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-c" {
			coreFlag = true
			continue
		}
		if coreFlag && v == "unlimited" {
			return []Violation{{
				KataID: "ZC1495",
				Message: "`ulimit -c unlimited` exposes setuid-process memory via core dumps. " +
					"Leave the distro default and use systemd-coredump if you need post-mortems.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
		coreFlag = false
	}
	return nil
}

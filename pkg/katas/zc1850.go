// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1850",
		Title:    "Warn on `ssh -o LogLevel=QUIET` — silences security-relevant ssh diagnostics",
		Severity: SeverityWarning,
		Description: "`LogLevel=QUIET` (aliased to the `-q` short flag) suppresses every " +
			"informational or warning message ssh would otherwise print: host-key " +
			"changes, key-exchange downgrades, agent-forwarding permission denials, " +
			"canonical-hostname rewrites. In a script, that means the output looks clean " +
			"even when ssh is shouting about a MITM on the other end. Keep the default " +
			"`INFO` level (or raise to `VERBOSE` during debugging), capture stderr to a " +
			"log if the noise bothers you, and never pair `LogLevel=QUIET` with " +
			"`StrictHostKeyChecking=no` in the same call — that combination actively " +
			"hides known-bad-key events.",
		Check: checkZC1850,
	})
}

func checkZC1850(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ssh" && ident.Value != "scp" && ident.Value != "sftp" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		var kv string
		switch {
		case v == "-o" && i+1 < len(args):
			kv = args[i+1].String()
		case strings.HasPrefix(v, "-o"):
			kv = strings.TrimPrefix(v, "-o")
		default:
			continue
		}
		if zc1850IsLogLevelQuiet(kv) {
			return []Violation{{
				KataID: "ZC1850",
				Message: "`" + ident.Value + " -o LogLevel=QUIET` silences host-key, " +
					"agent-forward, and canonical-hostname warnings — a MITM " +
					"event produces no stderr. Keep the default level; capture " +
					"stderr to a log if you need it clean.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1850IsLogLevelQuiet(kv string) bool {
	norm := strings.ToLower(strings.Trim(kv, "\"' \t"))
	if !strings.HasPrefix(norm, "loglevel") {
		return false
	}
	rest := strings.TrimSpace(strings.TrimPrefix(norm, "loglevel"))
	if !strings.HasPrefix(rest, "=") {
		return false
	}
	val := strings.TrimSpace(strings.TrimPrefix(rest, "="))
	return val == "quiet" || val == "fatal" || val == "error"
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1845",
		Title:    "Warn on `setopt PATH_DIRS` — slash-bearing command names fall back to `$PATH` lookup",
		Severity: SeverityWarning,
		Description: "`PATH_DIRS` (off by default) changes how Zsh resolves a command that " +
			"contains a `/`: instead of treating `./foo/bar` or `subdir/cmd` as a direct " +
			"path, Zsh walks `$path` and retries `${path[i]}/subdir/cmd` until one is " +
			"executable. The surface intent — run a local binary — is silently replaced by " +
			"`/usr/local/bin/subdir/cmd` or any other same-shaped subtree that exists on " +
			"`$PATH`. This gets even worse on shared build hosts where `$PATH` contains " +
			"user-owned directories. Leave the option off and call local binaries with an " +
			"explicit leading `./`, or hand the full absolute path to the shell.",
		Check: checkZC1845,
	})
}

func checkZC1845(node ast.Node) []Violation {
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
			v := arg.String()
			if zc1845IsPathDirs(v) {
				return zc1845Hit(cmd, "setopt "+v)
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOPATHDIRS" {
				return zc1845Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1845IsPathDirs(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "PATHDIRS"
}

func zc1845Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1845",
		Message: "`" + where + "` lets `subdir/cmd` fall back to a `$PATH` lookup — " +
			"a missing local binary silently runs a same-named subtree elsewhere on " +
			"`$PATH`. Leave the option off; call locals as `./subdir/cmd`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

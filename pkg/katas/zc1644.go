// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1644",
		Title:    "Error on `unzip -P SECRET` / `zip -P SECRET` — archive password in process list",
		Severity: SeverityError,
		Description: "`unzip -P PASSWORD` / `zip -P PASSWORD` (or the concatenated `-PPASSWORD` " +
			"form) places the archive password in argv. The expanded value shows up in `ps`, " +
			"`/proc/<pid>/cmdline`, shell history, and audit logs for every local user who " +
			"can list processes. Both tools prompt interactively if `-P` is absent — use that " +
			"for human workflows. For automation prefer an archive format with a real key-" +
			"derivation story (for example `7z -p` piped over stdin, or `age` / `gpg` " +
			"envelope encryption that reads keys from a protected file).",
		Check: checkZC1644,
	})
}

func checkZC1644(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "unzip" && ident.Value != "zip" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" && i+1 < len(cmd.Arguments) {
			return zc1644Hit(cmd, ident.Value)
		}
		if strings.HasPrefix(v, "-P") && len(v) > 2 {
			return zc1644Hit(cmd, ident.Value)
		}
	}
	return nil
}

func zc1644Hit(cmd *ast.SimpleCommand, name string) []Violation {
	return []Violation{{
		KataID: "ZC1644",
		Message: "`" + name + " -P` places the archive password in argv — visible via " +
			"`ps`. Drop `-P` for interactive prompt, or switch to `7z -p` (reads from " +
			"stdin) / `age` / `gpg` with keys in a protected file.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

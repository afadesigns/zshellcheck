// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1853",
		Title:    "Warn on `setopt MARK_DIRS` — glob-matched directories gain a silent trailing `/`",
		Severity: SeverityWarning,
		Description: "With `MARK_DIRS` on, every filename produced by a glob that resolves to a " +
			"directory picks up a trailing `/`. Inside a shell it looks harmless, but " +
			"scripts that pass the glob result to other tools break in quiet ways: " +
			"`[[ -f \"$f\" ]]` rejects `dir/` because it is not a regular file, `rm -f *` " +
			"sees `dir/` and silently skips it (GNU rm refuses to remove directories " +
			"without `-r`), and downstream hash maps indexed on basenames suddenly carry " +
			"two keys for what the user thinks is one entry. Keep the option off at the " +
			"script level and request the trailing slash per-glob with the `(/)` qualifier " +
			"(`dirs=( *(/) )`) when you really need directories only.",
		Check: checkZC1853,
	})
}

func checkZC1853(node ast.Node) []Violation {
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
			if zc1853IsMarkDirs(arg.String()) {
				return zc1853Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMARKDIRS" {
				return zc1853Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1853IsMarkDirs(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MARKDIRS"
}

func zc1853Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1853",
		Message: "`" + where + "` appends a trailing `/` to every glob-matched " +
			"directory — `[[ -f \"$f\" ]]` and `rm -f *` start skipping, hash " +
			"maps keyed on basenames double up. Keep the option off; use " +
			"`*(/)` when you need dirs only.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

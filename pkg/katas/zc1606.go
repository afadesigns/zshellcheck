// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1606",
		Title:    "Warn on `mkdir -m NNN` / `install -m NNN` with world-write bit (no sticky)",
		Severity: SeverityWarning,
		Description: "`mkdir -m 777 /path` and `install -m 777 src /dest` create a path that " +
			"every local user can write and rename inside. If the script later creates files " +
			"there, classic TOCTOU symlink attacks become trivial — the attacker drops a " +
			"symlink named like the expected output file, redirecting the write wherever they " +
			"choose. A sticky-bit mode (`1777`) mitigates this for shared temp dirs. Prefer " +
			"`mkdir -m 700` (or 750), and scope access by group or ACL rather than everyone.",
		Check: checkZC1606,
	})
}

var zc1606Names = map[string]struct{}{"mkdir": {}, "install": {}}

func checkZC1606(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	name := CommandIdentifier(cmd)
	if _, hit := zc1606Names[name]; !hit {
		return nil
	}
	mode := zc1606WorldWriteMode(cmd)
	if mode == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1606",
		Message: "`" + name + " -m " + mode + "` creates a world-writable path " +
			"without the sticky bit — TOCTOU symlink-attack ground. Use `-m 700` / " +
			"`-m 750`, or `-m 1777` if a shared sticky dir is actually needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

func zc1606WorldWriteMode(cmd *ast.SimpleCommand) string {
	for i := 0; i+1 < len(cmd.Arguments); i++ {
		if cmd.Arguments[i].String() != "-m" {
			continue
		}
		mode := cmd.Arguments[i+1].String()
		if zc1606IsWorldWritable(mode) {
			return mode
		}
	}
	return ""
}

func zc1606IsWorldWritable(mode string) bool {
	if len(mode) != 3 {
		return false
	}
	for _, c := range mode {
		if c < '0' || c > '7' {
			return false
		}
	}
	switch mode[2] {
	case '2', '3', '6', '7':
		return true
	}
	return false
}

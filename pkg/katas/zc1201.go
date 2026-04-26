// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1201",
		Title:    "Avoid `rsh`/`rlogin`/`rcp` — use `ssh`/`scp`",
		Severity: SeverityWarning,
		Description: "`rsh`, `rlogin`, and `rcp` are insecure legacy protocols. " +
			"Use `ssh`, `scp`, or `rsync` over SSH for encrypted remote operations.",
		Check: checkZC1201,
		Fix:   fixZC1201,
	})
}

// fixZC1201 rewrites the legacy `rsh` / `rlogin` / `rcp` command
// names to `ssh` / `ssh` / `scp` respectively. Single-edit
// replacement at the violation column. Argument syntax is
// compatible (host + optional command for rsh/rlogin/ssh; src dst
// for rcp/scp). Idempotent — a re-run sees `ssh` or `scp`, not
// the legacy names.
func fixZC1201(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	var replacement string
	switch ident.Value {
	case "rsh", "rlogin":
		replacement = "ssh"
	case "rcp":
		replacement = "scp"
	default:
		return nil
	}
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len(ident.Value) > len(source) {
		return nil
	}
	if string(source[off:off+len(ident.Value)]) != ident.Value {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len(ident.Value),
		Replace: replacement,
	}}
}

func checkZC1201(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "rsh" && name != "rlogin" && name != "rcp" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1201",
		Message: "Avoid `" + name + "` — it is an insecure legacy protocol. " +
			"Use `ssh`/`scp`/`rsync` for encrypted remote operations.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

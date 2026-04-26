// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1502",
		Title:    "Warn on `grep \"$var\" file` without `--` — flag injection when `$var` starts with `-`",
		Severity: SeverityWarning,
		Description: "Without a `--` end-of-flags marker, `grep` (and most POSIX tools) treats " +
			"any argument that starts with `-` as a flag. If `$var` comes from user input or a " +
			"fuzzed filename, an attacker can pass `--include=*secret*` or `-f /etc/shadow` " +
			"and get grep to read paths the script author never intended. Always write " +
			"`grep -- \"$var\" file` or use a grep-compatible library with explicit pattern API.",
		Check: checkZC1502,
		Fix:   fixZC1502,
	})
}

// fixZC1502 inserts `-- ` before the first variable-shaped argument
// of a grep / egrep / fgrep / rg / ag invocation that lacks the
// `--` end-of-options marker. Idempotent — the detector gates on
// the absence of `--`, so once `-- ` is present a re-run won't
// re-insert.
func fixZC1502(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "grep" && ident.Value != "egrep" && ident.Value != "fgrep" &&
		ident.Value != "rg" && ident.Value != "ag" {
		return nil
	}
	var firstVar ast.Expression
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--" {
			return nil
		}
		if firstVar == nil && (strings.HasPrefix(v, "\"$") || strings.HasPrefix(v, "$")) {
			firstVar = arg
		}
	}
	if firstVar == nil {
		return nil
	}
	tok := firstVar.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return nil
	}
	insLine, insCol := offsetLineColZC1502(source, off)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: "-- ",
	}}
}

func offsetLineColZC1502(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1502(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "grep" && ident.Value != "egrep" && ident.Value != "fgrep" &&
		ident.Value != "rg" && ident.Value != "ag" {
		return nil
	}

	hasDashDash := false
	firstVarArg := ""
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--" {
			hasDashDash = true
			break
		}
		// First variable-like pattern argument (no surrounding single quote means interpolation)
		if firstVarArg == "" && (strings.HasPrefix(v, "\"$") || strings.HasPrefix(v, "$")) {
			firstVarArg = v
		}
	}

	if firstVarArg == "" || hasDashDash {
		return nil
	}
	return []Violation{{
		KataID: "ZC1502",
		Message: "Variable `" + firstVarArg + "` used as pattern without `--` end-of-flags " +
			"marker — attacker-controlled leading `-` becomes a flag. Write `grep -- \"$var\"`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

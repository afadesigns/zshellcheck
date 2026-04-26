// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1273",
		Title:    "Use `grep -q` instead of redirecting grep output to `/dev/null`",
		Severity: SeverityStyle,
		Description: "`grep -q` suppresses output and exits on first match, which is faster and more " +
			"idiomatic than piping or redirecting to `/dev/null`.",
		Check: checkZC1273,
		Fix:   fixZC1273,
	})
}

// fixZC1273 inserts ` -q` after `grep` and strips the trailing
// `/dev/null` argument (including its leading whitespace). Two edits;
// the detector already gates on the absence of `-q`, so the rewrite
// is idempotent on a re-run.
func fixZC1273(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}
	var devNull ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "/dev/null" {
			devNull = arg
			break
		}
	}
	if devNull == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || IdentLenAt(source, nameOff) != len("grep") {
		return nil
	}
	insertAt := nameOff + len("grep")
	insLine, insCol := offsetLineColZC1273(source, insertAt)
	if insLine < 0 {
		return nil
	}
	stripEdits := zc1238StripFlag(source, devNull, "/dev/null")
	if stripEdits == nil {
		return nil
	}
	return append([]FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -q",
	}}, stripEdits...)
}

func offsetLineColZC1273(source []byte, offset int) (int, int) {
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

func checkZC1273(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasQuiet := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" || val == "--quiet" || val == "--silent" {
			hasQuiet = true
			break
		}
	}

	if hasQuiet {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/null" {
			return []Violation{{
				KataID:  "ZC1273",
				Message: "Use `grep -q` instead of redirecting to `/dev/null`. It is faster and more idiomatic.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

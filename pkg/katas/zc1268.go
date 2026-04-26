// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1268",
		Title:    "Use `du -sh --` to handle filenames starting with dash",
		Severity: SeverityInfo,
		Description: "`du -sh *` breaks if a filename starts with `-`. " +
			"Use `--` to signal end of options and safely handle all filenames.",
		Check: checkZC1268,
		Fix:   fixZC1268,
	})
}

// fixZC1268 inserts `-- ` before the first positional argument of a
// `du …` invocation that lacks the `--` end-of-options marker. The
// detector already gates on a glob (`*` / `.`) being present, and on
// the absence of `--`, so the insertion is idempotent.
func fixZC1268(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "du" {
		return nil
	}
	// Find the first positional (non-flag) argument.
	var positional ast.Expression
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if len(v) > 0 && v[0] != '-' {
			positional = arg
			break
		}
	}
	if positional == nil {
		return nil
	}
	tok := positional.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 {
		return nil
	}
	insLine, insCol := offsetLineColZC1268(source, off)
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

func offsetLineColZC1268(source []byte, offset int) (int, int) {
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

func checkZC1268(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "du" {
		return nil
	}

	hasEndOfOpts := false
	hasGlob := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--" {
			hasEndOfOpts = true
		}
		if val == "*" || val == "." {
			hasGlob = true
		}
	}

	if hasGlob && !hasEndOfOpts {
		return []Violation{{
			KataID: "ZC1268",
			Message: "Use `du -sh -- *` instead of `du -sh *`. The `--` prevents " +
				"filenames starting with `-` from being interpreted as options.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

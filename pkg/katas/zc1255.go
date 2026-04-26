// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1255",
		Title:    "Use `curl -L` to follow HTTP redirects",
		Severity: SeverityInfo,
		Description: "`curl` without `-L` does not follow redirects, returning 301/302 responses " +
			"instead of the actual content. Use `-L` to follow redirects automatically.",
		Check: checkZC1255,
		Fix:   fixZC1255,
	})
}

// fixZC1255 inserts ` -L` after the `curl` command name. Detector
// already guards against any existing follow-redirect flag, so the
// insertion is idempotent on a re-run.
func fixZC1255(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	if IdentLenAt(source, nameOff) != len("curl") {
		return nil
	}
	insertAt := nameOff + len("curl")
	insLine, insCol := offsetLineColZC1255(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -L",
	}}
}

func offsetLineColZC1255(source []byte, offset int) (int, int) {
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

func checkZC1255(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFollow := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-L" || val == "-fsSL" || val == "-fSL" || val == "-sL" {
			hasFollow = true
		}
		if len(val) > 7 && val[:5] == "https" {
			hasURL = true
		}
	}

	if hasURL && !hasFollow {
		return []Violation{{
			KataID: "ZC1255",
			Message: "Use `curl -L` to follow HTTP redirects. Without `-L`, curl returns " +
				"redirect responses (301/302) instead of the actual content.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

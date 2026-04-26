// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1153",
		Title:    "Use `cmp -s` instead of `diff` for equality check",
		Severity: SeverityStyle,
		Description: "When only checking if two files are identical (not viewing differences), " +
			"`cmp -s` is faster than `diff` as it stops at the first difference.",
		Check: checkZC1153,
		Fix:   fixZC1153,
	})
}

// fixZC1153 rewrites `diff -q FILE1 FILE2` into `cmp -s FILE1 FILE2`.
// Two non-overlapping edits: the command name (`diff` → `cmp`) and the
// quiet flag (`-q` → `-s`). Other arguments stay byte-identical.
// Idempotent because the detector gates on `diff -q` literal presence.
func fixZC1153(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}
	var dashQ ast.Expression
	for _, arg := range cmd.Arguments {
		if arg.String() == "-q" {
			dashQ = arg
			break
		}
	}
	if dashQ == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("diff") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("diff")]) != "diff" {
		return nil
	}
	dashTok := dashQ.TokenLiteralNode()
	dashOff := LineColToByteOffset(source, dashTok.Line, dashTok.Column)
	if dashOff < 0 || dashOff+2 > len(source) {
		return nil
	}
	if string(source[dashOff:dashOff+2]) != "-q" {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: len("diff"), Replace: "cmp"},
		{Line: dashTok.Line, Column: dashTok.Column, Length: 2, Replace: "-s"},
	}
}

func checkZC1153(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "diff" {
		return nil
	}

	// Only flag diff -q (quiet) which is used for equality checks
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" {
			return []Violation{{
				KataID: "ZC1153",
				Message: "Use `cmp -s file1 file2` instead of `diff -q`. " +
					"`cmp -s` is faster for equality checks as it stops at the first difference.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1219",
		Title:    "Use `curl -fsSL` instead of `wget -O -` for piped downloads",
		Severity: SeverityStyle,
		Description: "`wget -O -` outputs to stdout but lacks `curl`'s error handling. " +
			"`curl -fsSL` fails on HTTP errors, is silent, follows redirects, and is more portable.",
		Check: checkZC1219,
		Fix:   fixZC1219,
	})
}

// fixZC1219 collapses `wget -O- URL` / `wget -qO- URL` into
// `curl -fsSL URL`. The span covers the `wget` command name and the
// `-O-`/`-qO-` flag in a single edit so the rewrite stays deterministic
// even if a separate kata also fires on the `wget` name; trailing URL
// argument(s) stay in place.
func fixZC1219(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}
	var flag ast.Expression
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-O-" || val == "-qO-" {
			flag = arg
			break
		}
	}
	if flag == nil {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("wget") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("wget")]) != "wget" {
		return nil
	}
	flagTok := flag.TokenLiteralNode()
	flagOff := LineColToByteOffset(source, flagTok.Line, flagTok.Column)
	flagLen := len(flag.String())
	if flagOff < 0 || flagOff+flagLen > len(source) {
		return nil
	}
	if string(source[flagOff:flagOff+flagLen]) != flag.String() {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  flagOff + flagLen - nameOff,
		Replace: "curl -fsSL",
	}}
}

func checkZC1219(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wget" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-O-" || val == "-qO-" {
			return []Violation{{
				KataID: "ZC1219",
				Message: "Use `curl -fsSL` instead of `wget -O -` for piped downloads. " +
					"`curl` fails on HTTP errors and is available on more platforms.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityStyle,
			}}
		}
	}

	return nil
}

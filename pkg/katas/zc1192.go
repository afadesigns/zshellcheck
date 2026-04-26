// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1192",
		Title:    "Avoid `sleep 0` — it is a no-op external process",
		Severity: SeverityInfo,
		Description: "`sleep 0` spawns an external process that does nothing. " +
			"Remove it or use `:` if an explicit no-op is needed.",
		Check: checkZC1192,
		Fix:   fixZC1192,
	})
}

// fixZC1192 rewrites the no-op `sleep 0` invocation into `:`, the
// builtin no-op. Span covers the command name through the `0` arg.
func fixZC1192(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) != 1 {
		return nil
	}
	zeroArg := cmd.Arguments[0]
	if zeroArg.String() != "0" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || nameOff+len("sleep") > len(source) {
		return nil
	}
	if string(source[nameOff:nameOff+len("sleep")]) != "sleep" {
		return nil
	}
	argTok := zeroArg.TokenLiteralNode()
	argOff := LineColToByteOffset(source, argTok.Line, argTok.Column)
	if argOff < 0 || argOff+1 > len(source) || source[argOff] != '0' {
		return nil
	}
	end := argOff + 1
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - nameOff,
		Replace: ":",
	}}
}

func checkZC1192(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sleep" {
		return nil
	}

	if len(cmd.Arguments) == 1 && cmd.Arguments[0].String() == "0" {
		return []Violation{{
			KataID: "ZC1192",
			Message: "Remove `sleep 0` — it spawns a process that does nothing. " +
				"Use `:` if an explicit no-op is needed.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityInfo,
		}}
	}

	return nil
}

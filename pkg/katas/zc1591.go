// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1591",
		Title:    "Style: use Zsh `print -l` / `${(F)array}` instead of `printf '%s\\n' \"${array[@]}\"`",
		Severity: SeverityStyle,
		Description: "`printf '%s\\n' \"${array[@]}\"` is the Bash-idiomatic way to print one " +
			"element per line. Zsh has `print -l -r -- \"${array[@]}\"` (one element per line, " +
			"raw, sentinel-safe) and the parameter-expansion flag `${(F)array}` (newline-join, " +
			"fine for `$(...)`). Both are shorter than the printf incantation and avoid format-" +
			"string surprises if the array ever contains a literal `%`.",
		Check: checkZC1591,
		Fix:   fixZC1591,
	})
}

// fixZC1591 rewrites `printf '%s\n' "${array[@]}"` (or `printf '%s'
// "${array[@]}"`) to `print -l -r -- "${array[@]}"`. Single span
// replacement covers the `printf` command name and the format
// argument; subsequent args (the array expansion) pass through
// unchanged. Idempotent — a re-run sees `print`, not `printf`.
func fixZC1591(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "printf" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}
	formatArg := cmd.Arguments[0]
	formatVal := formatArg.String()
	trimmed := strings.Trim(formatVal, "'\"")
	if trimmed != `%s\n` && trimmed != "%s" {
		return nil
	}
	cmdOff := LineColToByteOffset(source, v.Line, v.Column)
	if cmdOff < 0 || cmdOff+len("printf") > len(source) {
		return nil
	}
	if string(source[cmdOff:cmdOff+len("printf")]) != "printf" {
		return nil
	}
	formatTok := formatArg.TokenLiteralNode()
	formatOff := LineColToByteOffset(source, formatTok.Line, formatTok.Column)
	if formatOff < 0 || formatOff+len(formatVal) > len(source) {
		return nil
	}
	if string(source[formatOff:formatOff+len(formatVal)]) != formatVal {
		return nil
	}
	end := formatOff + len(formatVal)
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - cmdOff,
		Replace: "print -l -r --",
	}}
}

func checkZC1591(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "printf" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	format := strings.Trim(cmd.Arguments[0].String(), "'\"")
	if format != `%s\n` && format != "%s" {
		return nil
	}

	second := cmd.Arguments[1].String()
	if !strings.Contains(second, "[@]") && !strings.Contains(second, "[*]") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1591",
		Message: "`printf '%s\\n' \"${array[@]}\"` — use Zsh `print -l -r -- \"${array[@]}\"` " +
			"or `${(F)array}` for newline-joined output.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

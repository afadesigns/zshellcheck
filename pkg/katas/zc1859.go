// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1859",
		Title:    "Warn on `unsetopt MULTIOS` — `cmd >a >b` silently keeps only the last redirection",
		Severity: SeverityWarning,
		Description: "`MULTIOS` is on by default in Zsh: `cmd >out.log >>archive.log` sends stdout " +
			"to both files via an implicit `tee`, and `cmd <a <b` concatenates the two " +
			"inputs in order. Disabling it reverts to POSIX-sh semantics — Zsh opens each " +
			"earlier redirection, closes it immediately, and only the last one in the " +
			"direction wins. Any script that was written for Zsh suddenly starts dropping " +
			"the `archive.log` tail, and log collectors that opened `archive.log` keep " +
			"the fd but never receive new lines. Keep the option on at the script level; " +
			"if one specific line really needs POSIX behaviour, wrap it in a function with " +
			"`setopt LOCAL_OPTIONS; unsetopt MULTIOS`.",
		Check: checkZC1859,
	})
}

func checkZC1859(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1859IsMultios(arg.String()) {
				return zc1859Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMULTIOS" {
				return zc1859Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1859IsMultios(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MULTIOS"
}

func zc1859Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1859",
		Message: "`" + where + "` reverts to POSIX single-output redirection — " +
			"`cmd >a >b` silently drops `a`, log collectors stop receiving new " +
			"lines. Keep the option on; scope inside a `LOCAL_OPTIONS` function " +
			"if one line really needs POSIX.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

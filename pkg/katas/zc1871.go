// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1871",
		Title:    "Warn on `setopt IGNORE_BRACES` — brace expansion stops working script-wide",
		Severity: SeverityWarning,
		Description: "`IGNORE_BRACES` is off by default in Zsh, which means `{1..10}`, " +
			"`file.{log,bak}`, and nested combinations like `{a..z}{1..3}` all expand " +
			"exactly as they do in Bash with `brace_expand` on. Turning it on disables " +
			"every one of those — `for i in {1..10}` iterates over the single literal " +
			"token `{1..10}`, and `cp app.{conf,conf.bak}` tries to copy a file literally " +
			"called `app.{conf,conf.bak}`. Scripts that depend on either numeric or " +
			"comma-list expansion silently become no-ops or fail with ENOENT. Keep the " +
			"option off; if you really need a literal brace string, quote the specific " +
			"argument (`'app.{conf,bak}'`) instead of flipping the shell-wide behaviour.",
		Check: checkZC1871,
	})
}

func checkZC1871(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1871IsIgnoreBraces(arg.String()) {
				return zc1871Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOIGNOREBRACES" {
				return zc1871Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1871IsIgnoreBraces(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "IGNOREBRACES"
}

func zc1871Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1871",
		Message: "`" + where + "` disables brace expansion — `for i in {1..10}` " +
			"loops once over the literal token, `cp app.{conf,bak}` fails ENOENT. " +
			"Keep the option off; quote the specific argument if you need a " +
			"literal brace string.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

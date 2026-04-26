// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1867",
		Title:    "Warn on `unsetopt GLOB` — pattern expansion turned off, `rm *.log` tries the literal filename",
		Severity: SeverityWarning,
		Description: "`GLOB` is on by default in Zsh: `*`, `?`, `[...]`, and `**/` expand against " +
			"the filesystem before the command runs. Turning the option off script-wide " +
			"(via `unsetopt GLOB` or the equivalent `setopt NO_GLOB`, same as POSIX " +
			"`set -f`) means every later pattern is handed to the command verbatim, so " +
			"`rm *.log` tries to remove a file literally named `*.log`, `for f in *.txt` " +
			"iterates over the single literal string, and expected-array-length checks " +
			"always return 1. Keep the option on at the script level; if one specific " +
			"line needs the pattern as a literal, quote the argument (`'*.log'`) or scope " +
			"with `setopt LOCAL_OPTIONS; setopt NO_GLOB` inside a function.",
		Check: checkZC1867,
	})
}

func checkZC1867(node ast.Node) []Violation {
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
			if zc1867IsGlob(arg.String()) {
				return zc1867Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOGLOB" {
				return zc1867Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1867IsGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "GLOB"
}

func zc1867Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1867",
		Message: "`" + where + "` disables glob expansion — `rm *.log` chases the " +
			"literal `*.log`, `for f in *.txt` loops once. Quote specific args or " +
			"scope with `LOCAL_OPTIONS` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

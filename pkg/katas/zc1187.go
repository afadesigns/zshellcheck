// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1187",
		Title:    "Avoid `notify-send` without fallback — check availability first",
		Severity: SeverityInfo,
		Description: "`notify-send` is Linux-only (libnotify). For portable notifications, " +
			"check `$OSTYPE` and fall back to `osascript` on macOS or `print` as default.",
		Check: checkZC1187,
	})
}

func checkZC1187(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "notify-send" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1187",
		Message: "Wrap `notify-send` with an `$OSTYPE` check or `command -v` guard. " +
			"It is Linux-only and will fail silently on macOS.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

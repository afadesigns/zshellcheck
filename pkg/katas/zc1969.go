// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1969",
		Title:    "Warn on `zsh -f` / `zsh -d` — skips `/etc/zsh*` and `~/.zsh*` startup files",
		Severity: SeverityWarning,
		Description: "`zsh -f` is the short form of `--no-rcs`, which skips every personal " +
			"and system-wide startup file: `/etc/zshenv`, `/etc/zprofile`, `/etc/zshrc`, " +
			"`/etc/zlogin`, `~/.zshenv`, `~/.zshrc`, `~/.zlogin`. `zsh -d` (`--no-" +
			"globalrcs`) drops only the `/etc/zsh*` set but keeps per-user ones. Either " +
			"form strips corp-mandated settings — proxy/hosts overrides, audit hooks, " +
			"umask, `HISTFILE` redirection, `PATH` hardening — silently. Use it " +
			"deliberately only for a pristine test harness or a minimal repro; never as " +
			"the shebang of a production script. When isolation is required, prefer " +
			"`env -i zsh` with an explicit allow-list of variables.",
		Check: checkZC1969,
	})
}

func checkZC1969(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "zsh" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-f" || v == "-d" {
			return zc1969Hit(cmd, "zsh "+v)
		}
	}
	return nil
}

func zc1969Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1969",
		Message: "`" + form + "` skips `/etc/zsh*` and `~/.zsh*` startup files — " +
			"corp proxy/audit/`PATH` hardening silently dropped. For a pristine " +
			"shell use `env -i zsh` with an explicit allow-list.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1991",
		Title:    "Warn on `setopt CSH_NULLCMD` — bare `> file` raises an error instead of running `$NULLCMD`",
		Severity: SeverityWarning,
		Description: "Default Zsh executes `$NULLCMD` (initially `cat`) when a line has " +
			"redirections but no command, so `> file < input` copies input to file " +
			"and `< file` pages through it with `$READNULLCMD` (initially `more`). " +
			"`setopt CSH_NULLCMD` drops the Zsh convention and follows csh — any " +
			"command line without an explicit command is a parse error, regardless of " +
			"redirections. Scripts that rely on the bare-redirect idiom (log " +
			"truncation via `> $LOG`, drop-in includes via `< file`, piped filters " +
			"built from aliases) stop working with a confusing `parse error near '<'`. " +
			"Keep the option off; write `: > file` (or `true > file`) explicitly when " +
			"you mean to truncate.",
		Check: checkZC1991,
	})
}

func checkZC1991(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var enabling bool
	switch ident.Value {
	case "setopt":
		enabling = true
	case "unsetopt":
		enabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1991Canonical(arg.String())
		switch v {
		case "CSHNULLCMD":
			if enabling {
				return zc1991Hit(cmd, "setopt CSH_NULLCMD")
			}
		case "NOCSHNULLCMD":
			if !enabling {
				return zc1991Hit(cmd, "unsetopt NO_CSH_NULLCMD")
			}
		}
	}
	return nil
}

func zc1991Canonical(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		out = append(out, c)
	}
	return string(out)
}

func zc1991Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1991",
		Message: "`" + form + "` makes `> file` / `< file` (no command) a parse error " +
			"— log truncation and bare-redirect idioms stop working. Write `: > " +
			"file` explicitly for truncation.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

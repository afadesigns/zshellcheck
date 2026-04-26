// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1930",
		Title:    "Warn on `unsetopt HASH_CMDS` — every command invocation re-walks `$PATH`",
		Severity: SeverityWarning,
		Description: "`HASH_CMDS` (on by default) caches the resolved absolute path of every " +
			"command after its first successful lookup. `unsetopt HASH_CMDS` disables the " +
			"cache, so each invocation re-walks every `$PATH` entry and re-runs `stat()` on " +
			"every candidate. On a slow filesystem (NFS home, encrypted volume, large `$PATH`) " +
			"this adds tens to hundreds of milliseconds per command and can double the runtime " +
			"of a long pipeline. Keep the option on; if you are changing a binary and want the " +
			"cache invalidated, `rehash` (one-shot) or `hash -r` is the scoped fix.",
		Check: checkZC1930,
	})
}

func checkZC1930(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	var disabling bool
	switch ident.Value {
	case "unsetopt":
		disabling = true
	case "setopt":
		disabling = false
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := zc1930Canonical(arg.String())
		switch v {
		case "HASHCMDS":
			if disabling {
				return zc1930Hit(cmd, "unsetopt HASH_CMDS")
			}
		case "NOHASHCMDS":
			if !disabling {
				return zc1930Hit(cmd, "setopt NO_HASH_CMDS")
			}
		}
	}
	return nil
}

func zc1930Canonical(s string) string {
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

func zc1930Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1930",
		Message: "`" + form + "` re-walks `$PATH` on every call — tens to hundreds of ms " +
			"per command on slow filesystems. Keep it on; use `rehash` or `hash -r` to " +
			"invalidate the cache after a targeted binary swap.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

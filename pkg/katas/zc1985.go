package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1985",
		Title:    "Warn on `setopt SH_FILE_EXPANSION` — expansion order flips from Zsh-native to sh/bash, `~` leaks",
		Severity: SeverityWarning,
		Description: "Default Zsh runs parameter expansion first, then filename/`~` " +
			"expansion — so a `VAR='~/cache'` keeps the tilde literal when you do " +
			"`mkdir -p -- $VAR` because the `~` never leaves the value. `setopt " +
			"SH_FILE_EXPANSION` (POSIX/sh ordering) flips the pass: filename expansion " +
			"runs first on the raw text, then parameter expansion happens, so the " +
			"same line suddenly makes the tilde resolve to `$HOME`, paths pointing at " +
			"`~evil/.cache` resolve into another user's home, and `=cmd` spellings " +
			"look up `$PATH` silently. Keep the option off; when a specific helper " +
			"needs POSIX ordering use `emulate -LR sh` inside that function.",
		Check: checkZC1985,
	})
}

func checkZC1985(node ast.Node) []Violation {
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
		v := zc1985Canonical(arg.String())
		switch v {
		case "SHFILEEXPANSION":
			if enabling {
				return zc1985Hit(cmd, "setopt SH_FILE_EXPANSION")
			}
		case "NOSHFILEEXPANSION":
			if !enabling {
				return zc1985Hit(cmd, "unsetopt NO_SH_FILE_EXPANSION")
			}
		}
	}
	return nil
}

func zc1985Canonical(s string) string {
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

func zc1985Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1985",
		Message: "`" + form + "` flips expansion order to POSIX — a `~` or `=cmd` " +
			"sitting inside a `$VAR` value suddenly resolves, so a user-typed " +
			"`~other/.cache` escapes into another home. Scope with `emulate -LR " +
			"sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

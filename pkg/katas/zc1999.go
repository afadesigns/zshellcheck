package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1999",
		Title:    "Warn on `setopt AUTO_NAMED_DIRS` — every scalar holding a directory path becomes `~name`",
		Severity: SeverityWarning,
		Description: "Off by default, Zsh only treats `~USER` / explicit `hash -d` entries as " +
			"named directories. `setopt AUTO_NAMED_DIRS` auto-registers any scalar " +
			"whose value is an existing directory — so `release=/srv/app/releases` " +
			"suddenly makes `~release/config` a valid path, and `ls ~release` lists " +
			"`/srv/app/releases`. That silently collides with real usernames " +
			"(`alice` in `/etc/passwd` vs. an `alice=$HOME/stage` scalar the script " +
			"happens to set) and turns every unquoted `~$var` inside a heredoc or " +
			"`cd` arg into a parameter that the prompt expander happily replaces " +
			"with the wrong path. Keep the option off; when a script legitimately " +
			"wants a named dir, register it explicitly with `hash -d NAME=PATH`.",
		Check: checkZC1999,
	})
}

func checkZC1999(node ast.Node) []Violation {
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
		v := zc1999Canonical(arg.String())
		switch v {
		case "AUTONAMEDDIRS":
			if enabling {
				return zc1999Hit(cmd, "setopt AUTO_NAMED_DIRS")
			}
		case "NOAUTONAMEDDIRS":
			if !enabling {
				return zc1999Hit(cmd, "unsetopt NO_AUTO_NAMED_DIRS")
			}
		}
	}
	return nil
}

func zc1999Canonical(s string) string {
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

func zc1999Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1999",
		Message: "`" + form + "` auto-registers every dir-valued scalar as `~name` — " +
			"collisions with real usernames and stray `~$var` expansions. " +
			"Register named dirs explicitly with `hash -d NAME=PATH`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

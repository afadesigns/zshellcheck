package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1942",
		Title:    "Warn on `setopt CLOBBER_EMPTY` — `>file` still overwrites zero-length files under `NO_CLOBBER`",
		Severity: SeverityWarning,
		Description: "`setopt CLOBBER_EMPTY` relaxes `NO_CLOBBER`: a bare `>file` redirect still " +
			"succeeds when the target is zero bytes. Scripts that rely on `setopt NO_CLOBBER` " +
			"as a guard against accidental overwrite lose their safety net for every " +
			"freshly-`touch`ed lock file, sentinel, or `install -D`-created placeholder — the " +
			"next stray `>sentinel` quietly overwrites it. Keep the option off; use `>|file` " +
			"explicitly when you do want to bypass the `NO_CLOBBER` guard for a specific write.",
		Check: checkZC1942,
	})
}

func checkZC1942(node ast.Node) []Violation {
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
		v := zc1942Canonical(arg.String())
		switch v {
		case "CLOBBEREMPTY":
			if enabling {
				return zc1942Hit(cmd, "setopt CLOBBER_EMPTY")
			}
		case "NOCLOBBEREMPTY":
			if !enabling {
				return zc1942Hit(cmd, "unsetopt NO_CLOBBER_EMPTY")
			}
		}
	}
	return nil
}

func zc1942Canonical(s string) string {
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

func zc1942Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1942",
		Message: "`" + form + "` lets `>file` overwrite zero-length files even under " +
			"`NO_CLOBBER` — `touch`ed lock / sentinel files lose their safety net. Keep " +
			"off; use explicit `>|file` to bypass `NO_CLOBBER` for a specific write.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

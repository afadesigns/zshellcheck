package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1918",
		Title:    "Warn on `setopt HIST_SUBST_PATTERN` — `!:s/old/new/` silently switches to pattern matching",
		Severity: SeverityWarning,
		Description: "`HIST_SUBST_PATTERN` makes the `:s` and `:&` history modifiers, as well as " +
			"the identically-named parameter-expansion modifier `${foo:s/pat/rep/}`, match on " +
			"patterns rather than literal strings. Text that looked safe as a constant " +
			"(`#` comments, `^` anchors, `?`, `*`) suddenly gets interpreted as glob " +
			"metacharacters, and replacements that always returned the original string now " +
			"edit it in surprising ways. Keep the option off and use `${var//pat/rep}` " +
			"explicitly when you do want glob substitution — that form declares the intent " +
			"at the call site instead of via a shell-wide flag.",
		Check: checkZC1918,
	})
}

func checkZC1918(node ast.Node) []Violation {
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
		v := zc1918Canonical(arg.String())
		switch v {
		case "HISTSUBSTPATTERN":
			if enabling {
				return zc1918Hit(cmd, "setopt HIST_SUBST_PATTERN")
			}
		case "NOHISTSUBSTPATTERN":
			if !enabling {
				return zc1918Hit(cmd, "unsetopt NO_HIST_SUBST_PATTERN")
			}
		}
	}
	return nil
}

func zc1918Canonical(s string) string {
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

func zc1918Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1918",
		Message: "`" + form + "` switches `:s` history/param modifiers to pattern " +
			"matching — literal `*`/`?`/`^` suddenly act as glob metacharacters. Keep it off; " +
			"use `${var//pat/rep}` when you actually want pattern substitution.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

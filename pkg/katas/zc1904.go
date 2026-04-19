package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1904",
		Title:    "Warn on `setopt KSH_GLOB` — reinterprets `*(pattern)` and breaks Zsh glob qualifiers",
		Severity: SeverityWarning,
		Description: "`setopt KSH_GLOB` turns `@(a|b)`, `*(x)`, `+(x)`, `?(x)`, `!(x)` into " +
			"Korn-shell extended glob operators. The side effect is that `*(N)`, `*(D)`, " +
			"`*(.)`, and every other Zsh glob qualifier stop working — `*(N)` becomes " +
			"\"zero or more `N` characters\", silently shattering null-glob idioms across the " +
			"script. If you need Korn-style patterns, prefer `setopt EXTENDED_GLOB` and its " +
			"`(^...)` / `(#...)` forms, which coexist with the qualifier syntax. Otherwise " +
			"scope the switch inside a function with `setopt LOCAL_OPTIONS KSH_GLOB`.",
		Check: checkZC1904,
	})
}

func checkZC1904(node ast.Node) []Violation {
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
		v := zc1904Canonical(arg.String())
		switch v {
		case "KSHGLOB":
			if enabling {
				return zc1904Hit(cmd, "setopt KSH_GLOB")
			}
		case "NOKSHGLOB":
			if !enabling {
				return zc1904Hit(cmd, "unsetopt NO_KSH_GLOB")
			}
		}
	}
	return nil
}

func zc1904Canonical(s string) string {
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

func zc1904Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1904",
		Message: "`" + form + "` reinterprets `*(...)` as a ksh-style operator — every " +
			"Zsh glob qualifier (`*(N)`, `*(D)`, `*(.)`) silently stops working. Prefer " +
			"`setopt EXTENDED_GLOB`, or scope inside a function with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

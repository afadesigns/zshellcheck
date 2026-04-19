package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1925",
		Title:    "Warn on `unsetopt EQUALS` — disables `=cmd` path expansion and tilde-after-colon",
		Severity: SeverityWarning,
		Description: "Zsh's `EQUALS` option (on by default) is what makes `=python`, `=ls`, and " +
			"`=vim` expand to the absolute path of the command via `$PATH` lookup. It also " +
			"drives the `PATH=~/bin:$PATH` tilde-after-colon expansion. `unsetopt EQUALS` " +
			"turns both off: `=cmd` becomes a literal argument (breaking any idiom that " +
			"relies on the short-path), and `PATH=~/bin:$PATH` stops expanding the tilde " +
			"inside the colon-separated list. Keep the option on; if one function needs " +
			"literal `=` arguments, scope via `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside it.",
		Check: checkZC1925,
	})
}

func checkZC1925(node ast.Node) []Violation {
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
		v := zc1925Canonical(arg.String())
		switch v {
		case "EQUALS":
			if disabling {
				return zc1925Hit(cmd, "unsetopt EQUALS")
			}
		case "NOEQUALS":
			if !disabling {
				return zc1925Hit(cmd, "setopt NO_EQUALS")
			}
		}
	}
	return nil
}

func zc1925Canonical(s string) string {
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

func zc1925Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1925",
		Message: "`" + form + "` turns off `=cmd` path expansion and tilde-after-colon — " +
			"`=python`/`=ls` become literals and `PATH=~/bin:$PATH` stops tilde-expanding. " +
			"Keep on; scope with `setopt LOCAL_OPTIONS; unsetopt EQUALS` inside a function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

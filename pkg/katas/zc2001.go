package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC2001",
		Title:    "Warn on `unsetopt EVAL_LINENO` — `$LINENO` inside `eval` stops tracking source, stack traces go blank",
		Severity: SeverityWarning,
		Description: "On by default, Zsh's `EVAL_LINENO` keeps `$LINENO`, `$funcfiletrace`, and " +
			"`$funcstack` pointing at the line inside the `eval`ed string where the " +
			"error actually happened. Turning the option off (`unsetopt EVAL_LINENO` " +
			"or `setopt NO_EVAL_LINENO`) reverts to pre-Zsh-4.3 behaviour: `$LINENO` " +
			"collapses to the line that launched the `eval`, so every runtime error " +
			"inside a generated config, a lazy-loaded function, or a `compile`d string " +
			"reports the same line number and the stack trace loses every frame past " +
			"the eval. Keep the option on; if strict POSIX-matching line numbers are " +
			"needed inside one helper, scope with `emulate -LR sh` in that function.",
		Check: checkZC2001,
	})
}

func checkZC2001(node ast.Node) []Violation {
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
		v := zc2001Canonical(arg.String())
		switch v {
		case "EVALLINENO":
			if !enabling {
				return zc2001Hit(cmd, "unsetopt EVAL_LINENO")
			}
		case "NOEVALLINENO":
			if enabling {
				return zc2001Hit(cmd, "setopt NO_EVAL_LINENO")
			}
		}
	}
	return nil
}

func zc2001Canonical(s string) string {
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

func zc2001Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC2001",
		Message: "`" + form + "` reverts `$LINENO` inside `eval` to the outer line — " +
			"errors in generated configs collapse to a single source line and " +
			"stack frames past `eval` vanish. Keep on; scope via `emulate -LR sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

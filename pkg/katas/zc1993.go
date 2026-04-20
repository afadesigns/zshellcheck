package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1993",
		Title:    "Warn on `setopt KSH_TYPESET` — `typeset var=$val` starts word-splitting the RHS",
		Severity: SeverityWarning,
		Description: "Off by default, Zsh treats every `typeset`/`declare` assignment like a " +
			"shell assignment: the whole RHS after `=` is one token, so `typeset " +
			"msg=\"a b c\"` produces a single-element string. `setopt KSH_TYPESET` " +
			"follows ksh instead — each word on the `typeset` line is its own " +
			"assignment or name, and the shell re-splits the RHS on whitespace. " +
			"Functions that used to accept `typeset path=$HOME/My Files` suddenly " +
			"treat `Files` as a second variable name, and `local` (an alias for " +
			"`typeset` inside functions) inherits the same change. Keep the option " +
			"off; if ksh compatibility is genuinely needed, scope with `emulate -LR " +
			"ksh` inside the helper that needs it.",
		Check: checkZC1993,
	})
}

func checkZC1993(node ast.Node) []Violation {
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
		v := zc1993Canonical(arg.String())
		switch v {
		case "KSHTYPESET":
			if enabling {
				return zc1993Hit(cmd, "setopt KSH_TYPESET")
			}
		case "NOKSHTYPESET":
			if !enabling {
				return zc1993Hit(cmd, "unsetopt NO_KSH_TYPESET")
			}
		}
	}
	return nil
}

func zc1993Canonical(s string) string {
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

func zc1993Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1993",
		Message: "`" + form + "` re-splits the RHS of every later `typeset`/`local` — " +
			"`typeset path=$HOME/My Files` now treats `Files` as a second name. " +
			"Scope with `emulate -LR ksh` inside the one helper that needs it.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1989",
		Title:    "Warn on `setopt REMATCH_PCRE` — `[[ =~ ]]` regex flips from POSIX ERE to PCRE, changes semantics",
		Severity: SeverityWarning,
		Description: "By default Zsh's `[[ $str =~ pattern ]]` uses POSIX extended regex " +
			"(ERE). `setopt REMATCH_PCRE` (after `zmodload zsh/pcre`) swaps the engine " +
			"to PCRE for every later match. Patterns that pass through both engines " +
			"change meaning subtly: `\\b` is a word boundary in PCRE but a literal `b` " +
			"in ERE, `\\d`/`\\s`/`\\w` work in PCRE but not ERE, lookahead/lookbehind " +
			"(`(?=…)`) parse in PCRE but error in ERE, and inline flags `(?i)` only " +
			"exist in PCRE. Flipping the option globally silently rewrites the " +
			"meaning of every existing regex — prefer an explicit `pcre_match`/`pcre_compile` " +
			"call when PCRE is needed, or a `setopt LOCAL_OPTIONS REMATCH_PCRE` inside " +
			"the single function that uses PCRE syntax.",
		Check: checkZC1989,
	})
}

func checkZC1989(node ast.Node) []Violation {
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
		v := zc1989Canonical(arg.String())
		switch v {
		case "REMATCHPCRE":
			if enabling {
				return zc1989Hit(cmd, "setopt REMATCH_PCRE")
			}
		case "NOREMATCHPCRE":
			if !enabling {
				return zc1989Hit(cmd, "unsetopt NO_REMATCH_PCRE")
			}
		}
	}
	return nil
}

func zc1989Canonical(s string) string {
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

func zc1989Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1989",
		Message: "`" + form + "` swaps `[[ =~ ]]` from POSIX ERE to PCRE — `\\b`, " +
			"`\\d`, lookahead, `(?i)` change meaning across every later match. " +
			"Prefer `pcre_match` where PCRE is needed or scope with `LOCAL_OPTIONS`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

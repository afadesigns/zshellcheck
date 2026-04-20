package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1987",
		Title:    "Warn on `setopt BRACE_CCL` — `{a-z}` expands to each character instead of staying literal",
		Severity: SeverityWarning,
		Description: "`BRACE_CCL` is off by default: `echo {a-z}` stays literal `a-z` in Zsh, " +
			"which is what most scripts that only want the numeric range form " +
			"`{1..10}` actually expect. `setopt BRACE_CCL` promotes single-character " +
			"ranges and enumerations inside braces to csh-style character-class " +
			"expansion, so `echo {a-z}` suddenly prints every letter from `a` to `z` " +
			"and `echo {ABC}` becomes `A B C`. Any later command line that embeds " +
			"single-character ranges — regex fragments, hex masks, CI job names with " +
			"stage suffixes — expands unexpectedly. Leave the option off; use `{a..z}` " +
			"when a real range is wanted and quote literals that contain `{…}`.",
		Check: checkZC1987,
	})
}

func checkZC1987(node ast.Node) []Violation {
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
		v := zc1987Canonical(arg.String())
		switch v {
		case "BRACECCL":
			if enabling {
				return zc1987Hit(cmd, "setopt BRACE_CCL")
			}
		case "NOBRACECCL":
			if !enabling {
				return zc1987Hit(cmd, "unsetopt NO_BRACE_CCL")
			}
		}
	}
	return nil
}

func zc1987Canonical(s string) string {
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

func zc1987Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1987",
		Message: "`" + form + "` promotes single-character braces to csh-style classes " +
			"— `{a-z}` now expands to every letter, `{ABC}` to `A B C`, breaking " +
			"regex/hex/CI-name literals. Use `{a..z}` when a real range is wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

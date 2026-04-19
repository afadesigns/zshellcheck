package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1906",
		Title:    "Warn on `setopt POSIX_CD` — changes when `cd` / `pushd` consult `CDPATH`",
		Severity: SeverityWarning,
		Description: "`setopt POSIX_CD` makes `cd`, `chdir`, and `pushd` skip `CDPATH` for any " +
			"argument that starts with `/`, `.`, or `..`. Zsh's default — consulting `CDPATH` " +
			"for anything that does not start with `/` — was exactly what made `cd foo` resolve " +
			"the \"project\" dir via `CDPATH` even when a local `./foo` existed. Flipping " +
			"the option globally makes scripts that relied on the Zsh behaviour silently enter " +
			"different directories. Keep the option off; if POSIX parity is needed, wrap a " +
			"single function with `emulate -LR sh`.",
		Check: checkZC1906,
	})
}

func checkZC1906(node ast.Node) []Violation {
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
		v := zc1906Canonical(arg.String())
		switch v {
		case "POSIXCD":
			if enabling {
				return zc1906Hit(cmd, "setopt POSIX_CD")
			}
		case "NOPOSIXCD":
			if !enabling {
				return zc1906Hit(cmd, "unsetopt NO_POSIX_CD")
			}
		}
	}
	return nil
}

func zc1906Canonical(s string) string {
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

func zc1906Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1906",
		Message: "`" + form + "` changes when `cd`/`pushd` read `CDPATH` — scripts that " +
			"relied on Zsh's default silently enter different directories. Keep it off; " +
			"wrap POSIX-specific code with `emulate -LR sh`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

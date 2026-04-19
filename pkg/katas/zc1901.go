package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1901",
		Title:    "Warn on `setopt POSIX_BUILTINS` — flips `command`/special-builtin semantics",
		Severity: SeverityWarning,
		Description: "`setopt POSIX_BUILTINS` switches Zsh to the POSIX rules for special " +
			"builtins: assignments before `export`, `readonly`, `eval`, `.`, `trap`, `set`, " +
			"etc. stay in the caller's scope, and `command builtin` can now resolve shell " +
			"builtins. Mid-script Zsh code written against native semantics — where those " +
			"assignments are local — silently leaks state. Leave the option off; scope any " +
			"POSIX-specific block with `emulate -LR sh` instead of toggling globally.",
		Check: checkZC1901,
	})
}

func checkZC1901(node ast.Node) []Violation {
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
		v := zc1901Canonical(arg.String())
		switch v {
		case "POSIXBUILTINS":
			if enabling {
				return zc1901Hit(cmd, "setopt POSIX_BUILTINS")
			}
		case "NOPOSIXBUILTINS":
			if !enabling {
				return zc1901Hit(cmd, "unsetopt NO_POSIX_BUILTINS")
			}
		}
	}
	return nil
}

func zc1901Canonical(s string) string {
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

func zc1901Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1901",
		Message: "`" + form + "` switches Zsh to POSIX special-builtin rules — " +
			"assignments before `export`/`readonly`/`eval` stop being local, silently " +
			"leaking state. Scope any POSIX block with `emulate -LR sh` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

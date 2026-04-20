package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1975",
		Title:    "Warn on `unsetopt EXEC` / `setopt NO_EXEC` — parser keeps scanning, commands stop running",
		Severity: SeverityWarning,
		Description: "`EXEC` is on by default; the shell both parses and runs each command. " +
			"Turning it off (`unsetopt EXEC` or `setopt NO_EXEC`) tells Zsh to parse " +
			"everything but silently skip the execution step — nothing fires, yet " +
			"parameter assignments on the same line don't either, `$?` stays frozen, " +
			"and functions that follow look defined but never run. That is the " +
			"semantics behind `zsh -n script.zsh` for a pure syntax check; flipping " +
			"the option in the middle of a production script converts every later " +
			"line into a no-op without a visible error. Run syntax checks via `zsh " +
			"-n` from the outside, never by flipping `EXEC` in-line.",
		Check: checkZC1975,
	})
}

func checkZC1975(node ast.Node) []Violation {
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
		v := zc1975Canonical(arg.String())
		switch v {
		case "EXEC":
			if !enabling {
				return zc1975Hit(cmd, "unsetopt EXEC")
			}
		case "NOEXEC":
			if enabling {
				return zc1975Hit(cmd, "setopt NO_EXEC")
			}
		}
	}
	return nil
}

func zc1975Canonical(s string) string {
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

func zc1975Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1975",
		Message: "`" + form + "` stops running commands but keeps parsing — every " +
			"later line becomes a silent no-op. For syntax checks run `zsh -n " +
			"script.zsh` from outside the script.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1997",
		Title:    "Warn on `setopt HIST_NO_FUNCTIONS` — function definitions skipped from `$HISTFILE`, breaks forensic trail",
		Severity: SeverityWarning,
		Description: "Default Zsh writes every command you type, including function " +
			"definitions, to `$HISTFILE`. `setopt HIST_NO_FUNCTIONS` suppresses " +
			"storage of commands that define a function. On a multi-admin box or a " +
			"shared root account this breaks the forensic trail — the function the " +
			"attacker just defined (or that an operator typed before running the " +
			"destructive bit) vanishes from history while the invocation that used " +
			"it still shows, leaving responders with a command that references a " +
			"name that no longer exists on disk or in any log. Keep the option off " +
			"and scope any hiding needs with the Zsh hook `zshaddhistory { return " +
			"1 }` inside a function where the secret actually lives.",
		Check: checkZC1997,
	})
}

func checkZC1997(node ast.Node) []Violation {
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
		v := zc1997Canonical(arg.String())
		switch v {
		case "HISTNOFUNCTIONS":
			if enabling {
				return zc1997Hit(cmd, "setopt HIST_NO_FUNCTIONS")
			}
		case "NOHISTNOFUNCTIONS":
			if !enabling {
				return zc1997Hit(cmd, "unsetopt NO_HIST_NO_FUNCTIONS")
			}
		}
	}
	return nil
}

func zc1997Canonical(s string) string {
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

func zc1997Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1997",
		Message: "`" + form + "` drops function-definition commands from `$HISTFILE` " +
			"— forensic trail loses the definition while the call that used it " +
			"still shows. Scope hiding via `zshaddhistory` hook instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

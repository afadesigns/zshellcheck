// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1995",
		Title:    "Warn on `unsetopt BGNICE` — background jobs run at full interactive priority, starve the foreground",
		Severity: SeverityWarning,
		Description: "Default Zsh applies `nice +5` to every backgrounded job so long-running " +
			"work does not starve the interactive session. `unsetopt BGNICE` (or " +
			"`setopt NO_BGNICE`) turns that off and bg jobs compete at the same " +
			"priority as the foreground shell — SSH keystroke handling, editor " +
			"redraws, and `cmd &` batch fan-out all feel laggy, and a single CPU-" +
			"bound bg job can peg every core of a container it shares with a human " +
			"operator. Keep the option on; when a background job legitimately needs " +
			"full priority (audio pipeline, realtime simulator), wrap just that one " +
			"with `nice -n 0 -- cmd &` or a systemd unit with `Nice=` instead of " +
			"flipping globally.",
		Check: checkZC1995,
	})
}

func checkZC1995(node ast.Node) []Violation {
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
		v := zc1995Canonical(arg.String())
		switch v {
		case "BGNICE":
			if !enabling {
				return zc1995Hit(cmd, "unsetopt BG_NICE")
			}
		case "NOBGNICE":
			if enabling {
				return zc1995Hit(cmd, "setopt NO_BG_NICE")
			}
		}
	}
	return nil
}

func zc1995Canonical(s string) string {
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

func zc1995Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1995",
		Message: "`" + form + "` drops the `nice +5` that bg jobs get by default — a " +
			"CPU-bound `cmd &` now competes with SSH/editor work. Wrap specific " +
			"jobs with `nice -n 0` or a systemd `Nice=` unit instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

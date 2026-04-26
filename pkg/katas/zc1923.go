// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1923",
		Title:    "Warn on `setopt PRINT_EXIT_VALUE` — every non-zero exit leaks a status line to stderr",
		Severity: SeverityWarning,
		Description: "`PRINT_EXIT_VALUE` makes Zsh emit `zsh: exit N` on stderr after every " +
			"foreground command that returns a non-zero status. In a script the stream is " +
			"typically captured by a supervisor or shipped to a log aggregator, and the " +
			"extra line reveals which tool returned what — including grep / test / curl " +
			"probes that were supposed to stay silent. Worse, tools that parse stderr for " +
			"diagnostics (`git`, `ssh`, `rsync`) now see interleaved shell chatter. Remove " +
			"the `setopt` call; if you actually want a per-command post-mortem, rely on " +
			"`precmd`/`preexec` hooks or an explicit `|| printf …`.",
		Check: checkZC1923,
	})
}

func checkZC1923(node ast.Node) []Violation {
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
		v := zc1923Canonical(arg.String())
		switch v {
		case "PRINTEXITVALUE":
			if enabling {
				return zc1923Hit(cmd, "setopt PRINT_EXIT_VALUE")
			}
		case "NOPRINTEXITVALUE":
			if !enabling {
				return zc1923Hit(cmd, "unsetopt NO_PRINT_EXIT_VALUE")
			}
		}
	}
	return nil
}

func zc1923Canonical(s string) string {
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

func zc1923Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1923",
		Message: "`" + form + "` prints `zsh: exit N` on stderr for every non-zero " +
			"exit — silent grep/test/curl probes suddenly leak status, and tools parsing " +
			"stderr see interleaved shell chatter. Remove; use `|| printf …` per call.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

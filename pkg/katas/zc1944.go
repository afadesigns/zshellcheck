// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1944",
		Title:    "Warn on `setopt IGNORE_EOF` — Ctrl-D no longer exits the shell, masking runaway pipelines",
		Severity: SeverityWarning,
		Description: "`IGNORE_EOF` tells the interactive shell to treat an end-of-file on stdin as " +
			"if it were nothing, so `Ctrl-D` stops terminating a login. In an unattended `zsh " +
			"-i -c` launch, or a sourced rc, this keeps a subshell alive that was supposed to " +
			"wind down when the controlling terminal went away — sudo sessions, SSH tunnels, " +
			"port-forwards, and build supervisors then linger long after the parent left. Keep " +
			"the option off; if a stale-tty guard is truly wanted, set `TMOUT=NN` for a timed " +
			"exit instead.",
		Check: checkZC1944,
	})
}

func checkZC1944(node ast.Node) []Violation {
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
		v := zc1944Canonical(arg.String())
		switch v {
		case "IGNOREEOF":
			if enabling {
				return zc1944Hit(cmd, "setopt IGNORE_EOF")
			}
		case "NOIGNOREEOF":
			if !enabling {
				return zc1944Hit(cmd, "unsetopt NO_IGNORE_EOF")
			}
		}
	}
	return nil
}

func zc1944Canonical(s string) string {
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

func zc1944Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1944",
		Message: "`" + form + "` makes `Ctrl-D` stop terminating the shell — subshells, " +
			"sudo holds, SSH tunnels linger after the parent left. Keep off; use " +
			"`TMOUT=NN` for a timed stale-tty exit if needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

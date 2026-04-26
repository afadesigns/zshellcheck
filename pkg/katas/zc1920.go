// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1920",
		Title:    "Warn on `setopt VERBOSE` — every executed command is echoed to stderr",
		Severity: SeverityWarning,
		Description: "`setopt VERBOSE` is Zsh's name for the POSIX `set -v` flag: the shell prints " +
			"each command line to stderr immediately after reading it. In a script that " +
			"processes secrets the stderr stream then carries every command that mentions them, " +
			"including `mysql -pSECRET`, `curl -u user:pass`, `export DB_PASS=…`. Unlike " +
			"`set -x` (which already has dedicated detectors) the `VERBOSE` flag is easy to " +
			"leave on by accident because the output looks like normal command echo. Remove " +
			"the call and rely on `printf` / a proper logger; if a debug trace is required, " +
			"scope it in a function with `setopt LOCAL_OPTIONS VERBOSE` then `unsetopt VERBOSE`.",
		Check: checkZC1920,
	})
}

func checkZC1920(node ast.Node) []Violation {
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
		v := zc1920Canonical(arg.String())
		switch v {
		case "VERBOSE":
			if enabling {
				return zc1920Hit(cmd, "setopt VERBOSE")
			}
		case "NOVERBOSE":
			if !enabling {
				return zc1920Hit(cmd, "unsetopt NO_VERBOSE")
			}
		}
	}
	return nil
}

func zc1920Canonical(s string) string {
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

func zc1920Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1920",
		Message: "`" + form + "` echoes every executed command to stderr — any line that " +
			"mentions a password, token, or API key leaks with the trace. Remove and use " +
			"`printf` / a logger, or scope via `setopt LOCAL_OPTIONS VERBOSE` in a helper.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

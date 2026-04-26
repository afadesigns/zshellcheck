// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1973",
		Title:    "Warn on `setopt POSIX_IDENTIFIERS` — restricts parameter names to ASCII, breaks Unicode `$var`",
		Severity: SeverityWarning,
		Description: "Zsh accepts Unicode parameter names by default: `$café`, `$π`, `$данные` " +
			"all parse. `setopt POSIX_IDENTIFIERS` tightens that to the POSIX subset — " +
			"ASCII letters, digits, underscore, not starting with a digit. Once the " +
			"option is on, every later `${café}` or `café=1` is a parse error, and " +
			"scripts/libraries that expose i18n-named vars stop loading. If you need " +
			"POSIX identifiers for a specific helper, scope it inside a function with " +
			"`emulate -LR sh`; leave the global option off so the rest of the shell " +
			"keeps the Zsh behaviour the user expects.",
		Check: checkZC1973,
	})
}

func checkZC1973(node ast.Node) []Violation {
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
		v := zc1973Canonical(arg.String())
		switch v {
		case "POSIXIDENTIFIERS":
			if enabling {
				return zc1973Hit(cmd, "setopt POSIX_IDENTIFIERS")
			}
		case "NOPOSIXIDENTIFIERS":
			if !enabling {
				return zc1973Hit(cmd, "unsetopt NO_POSIX_IDENTIFIERS")
			}
		}
	}
	return nil
}

func zc1973Canonical(s string) string {
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

func zc1973Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1973",
		Message: "`" + form + "` restricts parameter names to ASCII; later " +
			"`${café}`/`${π}` fail to parse and i18n-named libs stop loading. " +
			"Scope with `emulate -LR sh` inside the helper instead of flipping " +
			"globally.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

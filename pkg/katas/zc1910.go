// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1910",
		Title:    "Warn on `setopt GLOB_STAR_SHORT` — makes bare `**` recurse instead of matching literal",
		Severity: SeverityWarning,
		Description: "`GLOB_STAR_SHORT` teaches Zsh to expand bare `**` (not followed by `/`) as if " +
			"it were `**/*` — suddenly `rm **` wipes every file under the current directory " +
			"instead of erroring or matching the two-star literal. Scripts that pass `**` as a " +
			"literal argument to `grep`, `sed`, or a logger call silently turn into deep " +
			"directory recursions. Keep the option off; when you really need recursive globs, " +
			"spell `**/*` explicitly so reviewers can see the intent.",
		Check: checkZC1910,
	})
}

func checkZC1910(node ast.Node) []Violation {
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
		v := zc1910Canonical(arg.String())
		switch v {
		case "GLOBSTARSHORT":
			if enabling {
				return zc1910Hit(cmd, "setopt GLOB_STAR_SHORT")
			}
		case "NOGLOBSTARSHORT":
			if !enabling {
				return zc1910Hit(cmd, "unsetopt NO_GLOB_STAR_SHORT")
			}
		}
	}
	return nil
}

func zc1910Canonical(s string) string {
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

func zc1910Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1910",
		Message: "`" + form + "` turns bare `**` into `**/*` — `rm **` now wipes the tree. " +
			"Keep the option off and spell `**/*` when recursion is actually wanted.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

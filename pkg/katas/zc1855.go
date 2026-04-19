package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1855",
		Title:    "Avoid `$GROUPS` — Bash-only array; Zsh exposes supplementary groups as `$groups`",
		Severity: SeverityWarning,
		Description: "`$GROUPS` is a Bash magic parameter that holds the caller's supplementary " +
			"GIDs as a numeric array. Zsh does not populate `$GROUPS`; it has " +
			"`$groups`, a lowercase associative array keyed by group *name* with the GID " +
			"as value (`${(k)groups}` for names, `${(v)groups}` for IDs). Scripts ported " +
			"from Bash that iterate `${GROUPS[@]}` therefore see an empty list under " +
			"Zsh and silently skip group-membership checks. Use `${(k)groups}` for names " +
			"or `${(v)groups}` for numeric GIDs; the Zsh `id -Gn` fallback keeps the " +
			"script portable across shells.",
		Check: checkZC1855,
	})
}

func checkZC1855(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil {
		return nil
	}
	for _, arg := range cmd.Arguments {
		if zc1855RefersToGROUPS(arg.String()) {
			return []Violation{{
				KataID: "ZC1855",
				Message: "`$GROUPS` is a Bash-only array — Zsh populates `$groups` " +
					"(associative name→GID) instead. Iterate `${(k)groups}` for " +
					"names or `${(v)groups}` for GIDs, or fall back to `id -Gn`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

func zc1855RefersToGROUPS(v string) bool {
	// Walk the arg looking for `$GROUPS` or `${GROUPS...}` as a distinct
	// token. Accept trailing `[`, `}`, or end-of-string so callers like
	// `${GROUPS[@]}` still match but `$GROUPSIZE` does not.
	i := 0
	for {
		idx := strings.Index(v[i:], "GROUPS")
		if idx < 0 {
			return false
		}
		idx += i
		// Require `$` or `${` immediately before.
		prefixOK := false
		switch {
		case idx >= 2 && v[idx-2:idx] == "${":
			prefixOK = true
		case idx >= 1 && v[idx-1] == '$':
			prefixOK = true
		}
		if prefixOK {
			end := idx + len("GROUPS")
			if end == len(v) {
				return true
			}
			next := v[end]
			if next == '[' || next == '}' || next == '"' || next == ' ' || next == '\t' {
				return true
			}
		}
		i = idx + 1
		if i >= len(v) {
			return false
		}
	}
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1098",
		Title: "Use `(q)` flag for quoting variables in eval",
		Description: "When constructing a command string for `eval`, use the `(q)` flag (or `(qq)`, `(q-)`) to safely quote variables " +
			"and prevent command injection.",
		Check: checkZC1098,
	})
}

func checkZC1098(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "eval" {
		for _, arg := range cmd.Arguments {
			// Check if the argument string contains '$' and NOT '(q)'
			argStr := arg.String()
			// Very rough heuristic. Real parsing inside the string would be better but complex.
			// If arg contains '$' and not `(q)`, warn.
			// Also skip if it contains `(qq)` or `(q-)`.

			// We need to handle the case where user wrote `${(q)var}`.
			// arg.String() would be `${(q)var}`.

			// If we find `$` but no `(q`, warn.

			// Check for variable usage
			hasVar := false
			for i := 0; i < len(argStr); i++ {
				if argStr[i] == '$' {
					hasVar = true
					break
				}
			}

			if hasVar {
				// Check for quoting flags
				if !containsFlag(argStr) {
					return []Violation{{
						KataID:  "ZC1098",
						Message: "Use the `(q)` flag (or `(qq)`, `(q-)`) when using variables in `eval` to prevent injection.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
					}}
				}
			}
		}
	}

	return nil
}

func containsFlag(s string) bool {
	// Simple check for (q), (qq), (q-)
	// This is not perfect (could be inside a string literal), but good enough for a linter warning.
	// We look for `(q` pattern after `$`.
	// e.g. `${(q)var}` or `$var[(q)...]`? No, flags are at start of expansion.
	// `${(q)...}` or `$(...)` (command subst is also dangerous in eval without q).

	// Let's just check if the string contains "(q".
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '(' && s[i+1] == 'q' {
			return true
		}
	}
	return false
}

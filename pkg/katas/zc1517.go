package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1517",
		Title:    "Warn on `print -P \"$var\"` — prompt-escape injection via user-controlled string",
		Severity: SeverityWarning,
		Description: "`print -P` enables prompt-escape expansion (`%F`, `%K`, `%B`, `%S`, plus " +
			"arbitrary command substitution via `%{...%}`). Interpolating a shell variable " +
			"means any of those sequences inside the variable are expanded — at best messing " +
			"up terminal state, at worst running the attacker's command via `%(e:...)` or " +
			"similar. Either drop `-P` or wrap the variable with `${(q-)var}` / `${(V)var}` " +
			"to neutralize metacharacters before printing.",
		Check: checkZC1517,
	})
}

func checkZC1517(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "print" {
		return nil
	}

	hasP := false
	var varArg string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-P" {
			hasP = true
			continue
		}
		if !hasP || varArg != "" {
			continue
		}
		// Look for interpolation:  "$x" or $x (double-quoted or unquoted).
		raw := v
		if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
			raw = raw[1 : len(raw)-1]
		}
		if strings.Contains(raw, "$") && !(len(v) >= 2 && v[0] == '\'' && v[len(v)-1] == '\'') {
			varArg = v
		}
	}
	if !hasP || varArg == "" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1517",
		Message: "`print -P " + varArg + "` expands prompt escapes inside the variable — use " +
			"`${(V)var}` / `${(q-)var}` to neutralize metacharacters, or drop -P.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

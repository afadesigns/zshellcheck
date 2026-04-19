package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1865",
		Title:    "Warn on `unsetopt CASE_MATCH` — `[[ =~ ]]` and pattern tests quietly fold case",
		Severity: SeverityWarning,
		Description: "`CASE_MATCH` on is Zsh's default: `[[ $x =~ ^FOO ]]`, `[[ $x == Foo* ]]`, " +
			"and the subst-in-conditional forms honour letter case exactly as written. " +
			"Turning the option off flips every later test to case-insensitive — " +
			"`[[ $user == Admin ]]` also matches `admin`/`ADMIN`, regex dispatchers stop " +
			"distinguishing `README` from `readme`, and log-pattern filters over-collect. " +
			"Keep the option on at script level; if one specific regex really needs " +
			"case-folding, request it per-pattern with the Zsh `(#i)` flag " +
			"(e.g. `[[ $x =~ (#i)foo ]]`).",
		Check: checkZC1865,
	})
}

func checkZC1865(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			if zc1865IsCaseMatch(arg.String()) {
				return zc1865Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOCASEMATCH" {
				return zc1865Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1865IsCaseMatch(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "CASEMATCH"
}

func zc1865Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1865",
		Message: "`" + where + "` flips every `[[ =~ ]]` / `[[ == pat ]]` to " +
			"case-insensitive — `Admin` matches `ADMIN`, dispatchers collide. " +
			"Keep it on; scope per-line with `(#i)pattern`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1897",
		Title:    "Warn on `setopt SH_GLOB` — Zsh-specific glob patterns (`*(N)`, `<1-10>`, alternation) stop parsing",
		Severity: SeverityWarning,
		Description: "`SH_GLOB` is off by default in Zsh. With it off, the shell recognises Zsh's " +
			"extended patterns: `*(N)` null-glob qualifier, `<1-10>` numeric range globs, " +
			"`(alt1|alt2)` in-glob alternation, and the whole `(#i)`/`(#c,m)` flag " +
			"family. Turning the option on forces strict POSIX-sh parsing, so the parser " +
			"re-interprets `(...)` as command grouping and the null-glob / range idioms " +
			"raise parse errors. Every kata recommending `*(N)` (see ZC1830, ZC1893) " +
			"silently breaks, and downstream helpers sourced after the setopt inherit the " +
			"restricted pattern syntax. Keep the option off; scope inside a function " +
			"with `setopt LOCAL_OPTIONS; setopt SH_GLOB` if a specific block genuinely " +
			"needs POSIX patterns.",
		Check: checkZC1897,
	})
}

func checkZC1897(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	switch ident.Value {
	case "setopt":
		for _, arg := range cmd.Arguments {
			if zc1897IsShGlob(arg.String()) {
				return zc1897Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOSHGLOB" {
				return zc1897Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1897IsShGlob(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "SHGLOB"
}

func zc1897Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1897",
		Message: "`" + where + "` disables Zsh-extended glob patterns — `*(N)` " +
			"qualifiers, `<1-10>` ranges, and `(alt1|alt2)` alternation raise " +
			"parse errors. Keep the option off; scope with `LOCAL_OPTIONS` if " +
			"strict POSIX is needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

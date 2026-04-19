package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1893",
		Title:    "Warn on `unsetopt BARE_GLOB_QUAL` — `*(N)` null-glob qualifier stops being special",
		Severity: SeverityWarning,
		Description: "`BARE_GLOB_QUAL` is on by default in Zsh — that is what makes the " +
			"per-glob qualifier syntax (`*(N)` for null-glob, `*(.x)` for " +
			"executable, `*(Om)` for sort-by-mtime) work. Turning it off reverts " +
			"to ksh-style parsing where `(...)` inside a glob is a pattern " +
			"alternation, so `*(N)` stops being a null-glob and turns into " +
			"\"match zero-or-one N\" — a completely different pattern. Scripts that " +
			"relied on `for f in *.log(N)` to cope with empty directories then " +
			"silently iterate the literal string or fail under NOMATCH. Keep the " +
			"option on; if you really want ksh-style qualifiers, use " +
			"`setopt LOCAL_OPTIONS; unsetopt BARE_GLOB_QUAL` inside a function.",
		Check: checkZC1893,
	})
}

func checkZC1893(node ast.Node) []Violation {
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
			if zc1893IsBareGlobQual(arg.String()) {
				return zc1893Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOBAREGLOBQUAL" {
				return zc1893Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1893IsBareGlobQual(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "BAREGLOBQUAL"
}

func zc1893Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1893",
		Message: "`" + where + "` disables `*(qualifier)` syntax — `*(N)` stops being " +
			"null-glob and becomes an alternation, so null-glob idioms silently " +
			"break. Keep the option on; scope inside a `LOCAL_OPTIONS` function.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

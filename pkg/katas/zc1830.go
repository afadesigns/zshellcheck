package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1830",
		Title:    "Warn on `unsetopt NOMATCH` — unmatched glob becomes the literal pattern, silent bugs",
		Severity: SeverityWarning,
		Description: "`NOMATCH` is on by default in Zsh — an unmatched glob (`*.log` with no matching " +
			"files) errors out instead of silently passing through. Disabling it " +
			"(`unsetopt NOMATCH` or the equivalent `setopt NO_NOMATCH`) reverts to POSIX-sh " +
			"behaviour: the pattern is handed to the command verbatim, so `rm *.log` with no " +
			"matches runs `rm '*.log'` — which fails noisily for `rm` but, for commands that " +
			"accept arbitrary strings, silently processes the literal `*.log` instead of " +
			"files. Prefer scoped `*(N)` null-glob qualifier or `setopt LOCAL_OPTIONS; setopt " +
			"NULL_GLOB` inside a function, so the rest of the script keeps the default safety.",
		Check: checkZC1830,
	})
}

func checkZC1830(node ast.Node) []Violation {
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
			if zc1830IsNomatch(arg.String()) {
				return zc1830Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NONOMATCH" {
				return zc1830Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1830IsNomatch(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "NOMATCH"
}

func zc1830Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1830",
		Message: "`" + where + "` silences Zsh's unmatched-glob error — typos pass " +
			"through literally. Use `*(N)` per-glob or scope inside a function " +
			"with `setopt LOCAL_OPTIONS; setopt NULL_GLOB`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

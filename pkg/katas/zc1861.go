package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1861",
		Title:    "Warn on `setopt OCTAL_ZEROES` — leading-zero integers silently reinterpret as octal",
		Severity: SeverityWarning,
		Description: "`OCTAL_ZEROES` is off in Zsh by default: arithmetic treats `0100` as the " +
			"decimal integer one hundred, matching what every other scripting language " +
			"does. Setting it on reverts to POSIX-shell semantics where the leading `0` " +
			"flags the literal as octal — `(( n = 0100 ))` assigns 64, not 100. Scripts " +
			"that read timestamps padded to `00:59`, CSVs of phone-number prefixes " +
			"(`0049`), or file modes formatted as `0700` silently return the wrong " +
			"integer. Keep the option off at script level; if you really want C-style " +
			"octal literals, stay explicit with `(( n = 8#100 ))` or `$(( 8#$val ))` " +
			"so the intent is obvious.",
		Check: checkZC1861,
	})
}

func checkZC1861(node ast.Node) []Violation {
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
			if zc1861IsOctalZeroes(arg.String()) {
				return zc1861Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOOCTALZEROES" {
				return zc1861Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1861IsOctalZeroes(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "OCTALZEROES"
}

func zc1861Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1861",
		Message: "`" + where + "` reinterprets leading-zero integers as octal — " +
			"`(( n = 0100 ))` assigns 64 instead of 100, breaking timestamp, " +
			"phone-prefix, and mode parsing. Keep the option off; use `8#100` " +
			"when you want explicit octal.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

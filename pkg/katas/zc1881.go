package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1881",
		Title:    "Warn on `unsetopt MULTIBYTE` — `${#str}`, substring, and `[[ =~ ]]` stop counting characters",
		Severity: SeverityWarning,
		Description: "`MULTIBYTE` is on in Zsh by default: `${#str}` returns character count, " +
			"`${str:0:3}` extracts the first three characters, and `[[ $str =~ ... ]]` " +
			"matches whole UTF-8 codepoints. Turning it off reverts every string " +
			"operation to per-byte math, so an emoji that encodes to four bytes counts " +
			"as four, a substring spanning a multi-byte character slices mid-codepoint " +
			"and produces invalid UTF-8, and `[[ =~ ]]` regex ranges no longer cover " +
			"Unicode blocks. Filenames containing non-ASCII, i18n log strings, and JSON " +
			"snippets silently drift from their assumed layout. Keep the option on; if " +
			"you truly need byte-level counting, use `${#${(%)str}}` or `wc -c <<< $str`.",
		Check: checkZC1881,
	})
}

func checkZC1881(node ast.Node) []Violation {
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
			if zc1881IsMultibyte(arg.String()) {
				return zc1881Hit(cmd, "unsetopt "+arg.String())
			}
		}
	case "setopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NOMULTIBYTE" {
				return zc1881Hit(cmd, "setopt "+v)
			}
		}
	}
	return nil
}

func zc1881IsMultibyte(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "MULTIBYTE"
}

func zc1881Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1881",
		Message: "`" + where + "` flips every string op to per-byte math — `${#str}` " +
			"on an emoji returns 4, substrings slice mid-codepoint, `[[ =~ ]]` " +
			"Unicode ranges break. Keep the option on; byte-count with " +
			"`wc -c <<< $str` when truly needed.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

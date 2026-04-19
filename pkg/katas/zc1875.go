package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1875",
		Title:    "Warn on `setopt RC_QUOTES` — `''` inside single quotes flips from empty-concat to literal apostrophe",
		Severity: SeverityWarning,
		Description: "`RC_QUOTES` is off by default in Zsh: inside a single-quoted string `'it''s'` " +
			"parses as two adjacent single-quoted regions with an empty middle, producing " +
			"the literal `its`. Turning the option on reinterprets the doubled apostrophe " +
			"as one escaped quote, so `'it''s'` suddenly becomes `it's`. That is a " +
			"source-level change to every already-written string literal in the file — " +
			"password strings, SQL fragments, display text — so log lines, stored tokens, " +
			"and API payloads silently diverge. Keep the option off; write a literal " +
			"apostrophe with `\\'` outside the quotes or with double-quoted wrapping.",
		Check: checkZC1875,
	})
}

func checkZC1875(node ast.Node) []Violation {
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
			if zc1875IsRCQuotes(arg.String()) {
				return zc1875Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NORCQUOTES" {
				return zc1875Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1875IsRCQuotes(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "RCQUOTES"
}

func zc1875Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1875",
		Message: "`" + where + "` reinterprets `''` inside single quotes as a " +
			"literal apostrophe — `'it''s'` flips from `its` to `it's`, " +
			"breaking tokens and SQL. Use double quotes or `\\'` instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

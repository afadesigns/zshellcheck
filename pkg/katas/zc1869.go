package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1869",
		Title:    "Warn on `setopt RC_EXPAND_PARAM` — brace-adjacent array expansion silently distributes",
		Severity: SeverityWarning,
		Description: "`RC_EXPAND_PARAM` is off in Zsh by default: `echo x${arr[@]}y` concatenates " +
			"once, producing `xay xby xcy` only if you wrote the template carefully. " +
			"Turning it on changes the rule — every adjacent literal is distributed " +
			"across each array element, so `cp src/${files[@]}.bak /tmp` suddenly " +
			"rewrites as `cp src/a.bak src/b.bak src/c.bak /tmp`. That is exactly what " +
			"you want when you want it, and a nasty surprise anywhere else because the " +
			"same syntax keeps working silently. Leave the option off at script level; " +
			"if one specific line needs distributive expansion, request it per-use with " +
			"`${^arr}` (the `^` flag scopes the behaviour to that parameter only).",
		Check: checkZC1869,
	})
}

func checkZC1869(node ast.Node) []Violation {
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
			if zc1869IsRcExpandParam(arg.String()) {
				return zc1869Hit(cmd, "setopt "+arg.String())
			}
		}
	case "unsetopt":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
			if norm == "NORCEXPANDPARAM" {
				return zc1869Hit(cmd, "unsetopt "+v)
			}
		}
	}
	return nil
}

func zc1869IsRcExpandParam(v string) bool {
	norm := strings.ToUpper(strings.ReplaceAll(v, "_", ""))
	return norm == "RCEXPANDPARAM"
}

func zc1869Hit(cmd *ast.SimpleCommand, where string) []Violation {
	return []Violation{{
		KataID: "ZC1869",
		Message: "`" + where + "` distributes literal prefix/suffix across every " +
			"array element — `cp src/${arr[@]}.bak dst` silently rewrites as " +
			"`cp src/a.bak src/b.bak dst`. Keep it off; opt in per-use with " +
			"`${^arr}`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

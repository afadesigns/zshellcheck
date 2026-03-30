package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1108",
		Title: "Use Zsh case conversion instead of `tr`",
		Description: "Zsh provides `${(U)var}` for uppercase and `${(L)var}` for lowercase. " +
			"Avoid piping through `tr '[:lower:]' '[:upper:]'` for simple case conversion.",
		Check: checkZC1108,
	})
}

func checkZC1108(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tr" {
		return nil
	}

	if len(cmd.Arguments) != 2 {
		return nil
	}

	arg1 := strings.Trim(cmd.Arguments[0].String(), "'\"")
	arg2 := strings.Trim(cmd.Arguments[1].String(), "'\"")

	isLowerToUpper := (arg1 == "[:lower:]" && arg2 == "[:upper:]") ||
		(arg1 == "a-z" && arg2 == "A-Z")
	isUpperToLower := (arg1 == "[:upper:]" && arg2 == "[:lower:]") ||
		(arg1 == "A-Z" && arg2 == "a-z")

	if !isLowerToUpper && !isUpperToLower {
		return nil
	}

	var suggestion string
	if isLowerToUpper {
		suggestion = "`${(U)var}`"
	} else {
		suggestion = "`${(L)var}`"
	}

	return []Violation{{
		KataID: "ZC1108",
		Message: "Use " + suggestion + " for case conversion instead of `tr`. " +
			"Zsh parameter expansion flags avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
	}}
}

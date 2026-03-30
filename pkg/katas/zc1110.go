package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1110",
		Title: "Use Zsh subscripts instead of `head -1` or `tail -1`",
		Description: "Zsh array subscripts `${lines[1]}` and `${lines[-1]}` can extract the first or last " +
			"element without spawning `head` or `tail` as external processes.",
		Severity: SeverityStyle,
		Check:    checkZC1110,
	})
}

func checkZC1110(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "head" && name != "tail" {
		return nil
	}

	// Only flag `head -1` / `head -n 1` / `tail -1` / `tail -n 1` without file args
	hasFileArg := false
	isSingleLine := false
	skipNext := false

	for i, arg := range cmd.Arguments {
		if skipNext {
			skipNext = false
			continue
		}
		val := arg.String()
		switch {
		case val == "-1" || val == "-n1":
			isSingleLine = true
		case val == "-n" && i+1 < len(cmd.Arguments) && cmd.Arguments[i+1].String() == "1":
			isSingleLine = true
			skipNext = true
		case len(val) > 0 && val[0] != '-':
			hasFileArg = true
		}
	}

	if hasFileArg || !isSingleLine {
		return nil
	}

	var suggestion string
	if name == "head" {
		suggestion = "`${lines[1]}`"
	} else {
		suggestion = "`${lines[-1]}`"
	}

	return []Violation{{
		KataID: "ZC1110",
		Message: "Use " + suggestion + " instead of `" + name + " -1`. " +
			"Zsh array subscripts avoid spawning an external process.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

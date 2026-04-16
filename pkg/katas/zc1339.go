package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1339",
		Title:    "Use Zsh `${#${(f)var}}` instead of `wc -l` for line count",
		Severity: SeverityStyle,
		Description: "Zsh `${(f)var}` splits a string into lines and `${#...}` counts them. " +
			"Avoid piping through `wc -l` for simple line counting from variables.",
		Check: checkZC1339,
	})
}

func checkZC1339(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "wc" {
		return nil
	}

	hasLineFlag := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-l" {
			hasLineFlag = true
		} else if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasLineFlag && !hasFile {
		return []Violation{{
			KataID: "ZC1339",
			Message: "Use Zsh `${#${(f)var}}` for line counting instead of piping through `wc -l`. " +
				"Parameter expansion avoids spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

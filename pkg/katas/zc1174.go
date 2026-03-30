package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1174",
		Title:    "Use Zsh `${(j:delim:)}` instead of `paste -sd`",
		Severity: SeverityStyle,
		Description: "Zsh `${(j:delim:)array}` joins array elements with a delimiter. " +
			"Avoid spawning `paste` for simple field joining from variables.",
		Check: checkZC1174,
	})
}

func checkZC1174(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "paste" {
		return nil
	}

	hasSD := false
	hasFile := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-sd" || val == "-s" {
			hasSD = true
		}
		if len(val) > 0 && val[0] != '-' {
			hasFile = true
		}
	}

	if hasSD && !hasFile {
		return []Violation{{
			KataID: "ZC1174",
			Message: "Use Zsh `${(j:delim:)array}` to join array elements instead of `paste -sd`. " +
				"Parameter expansion avoids spawning an external process.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

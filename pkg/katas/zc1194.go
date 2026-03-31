package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1194",
		Title:    "Avoid `sed` with multiple `-e` — use a single script",
		Severity: SeverityStyle,
		Description: "Multiple `sed -e 's/a/b/' -e 's/c/d/'` can be combined into " +
			"`sed 's/a/b/; s/c/d/'` for cleaner syntax and fewer shell word splits.",
		Check: checkZC1194,
	})
}

func checkZC1194(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sed" {
		return nil
	}

	eCount := 0
	for _, arg := range cmd.Arguments {
		if arg.String() == "-e" {
			eCount++
		}
	}

	if eCount >= 2 {
		return []Violation{{
			KataID: "ZC1194",
			Message: "Combine multiple `sed -e` expressions into a single script: " +
				"`sed 's/a/b/; s/c/d/'` is cleaner than multiple `-e` flags.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

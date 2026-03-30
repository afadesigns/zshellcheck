package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1168",
		Title:    "Use `${(f)...}` instead of `readarray`/`mapfile`",
		Severity: SeverityStyle,
		Description: "`readarray` and `mapfile` are Bash builtins not available in Zsh. " +
			"Use Zsh `${(f)...}` parameter expansion flag to split output into an array by newlines.",
		Check: checkZC1168,
	})
}

func checkZC1168(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "readarray" && ident.Value != "mapfile" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1168",
		Message: "Use Zsh `${(f)$(cmd)}` instead of `" + ident.Value + "`. " +
			"`readarray`/`mapfile` are Bash builtins not available in Zsh.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}

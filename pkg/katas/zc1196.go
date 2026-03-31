package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1196",
		Title:    "Avoid `cat` for reading single file into variable",
		Severity: SeverityStyle,
		Description: "Use Zsh `$(<file)` instead of `$(cat file)` to read file contents. " +
			"`$(<file)` is a Zsh builtin that avoids spawning cat.",
		Check: checkZC1196,
	})
}

func checkZC1196(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "less" {
		return nil
	}

	// less without flags in a script is likely a mistake
	hasFlags := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			hasFlags = true
		}
	}

	if !hasFlags && len(cmd.Arguments) > 0 {
		return []Violation{{
			KataID: "ZC1196",
			Message: "Avoid `less` in scripts — it requires interactive terminal input. " +
				"Use `cat` or redirect output to a pager only when `$TERM` is available.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

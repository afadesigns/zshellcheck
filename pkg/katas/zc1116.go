package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1116",
		Title: "Use Zsh multios instead of `tee`",
		Description: "Zsh `setopt multios` allows redirecting output to multiple files with " +
			"`cmd > file1 > file2`. Avoid spawning `tee` for simple output duplication.",
		Check: checkZC1116,
	})
}

func checkZC1116(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "tee" {
		return nil
	}

	// Only flag simple tee without -a (append) or -i (ignore interrupt)
	// tee -a is append mode which multios handles differently
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if len(val) > 0 && val[0] == '-' {
			return nil
		}
	}

	// Must have at least one file argument
	if len(cmd.Arguments) == 0 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1116",
		Message: "Use Zsh multios (`setopt multios`) instead of `tee`. " +
			"With multios, `cmd > file1 > file2` writes to both files without spawning tee.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
	}}
}

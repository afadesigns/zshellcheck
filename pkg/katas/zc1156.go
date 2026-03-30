package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1156",
		Title:    "Avoid `ln` without `-s` for symlinks",
		Severity: SeverityInfo,
		Description: "Hard links (`ln` without `-s`) share inodes and can cause confusion. " +
			"Prefer symbolic links (`ln -s`) unless you specifically need hard links.",
		Check: checkZC1156,
	})
}

func checkZC1156(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ln" {
		return nil
	}

	hasSymlink := false
	hasForce := false
	fileCount := 0

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-s" || val == "-sf" || val == "-snf" {
			hasSymlink = true
		}
		if val == "-f" {
			hasForce = true
		}
		if len(val) > 0 && val[0] != '-' {
			fileCount++
		}
	}

	_ = hasForce

	if hasSymlink || fileCount < 2 {
		return nil
	}

	return []Violation{{
		KataID: "ZC1156",
		Message: "Use `ln -s` for symbolic links instead of hard links. " +
			"Hard links share inodes and don't work across filesystems.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityInfo,
	}}
}

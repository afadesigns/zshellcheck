package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1270",
		Title:    "Use `mktemp` instead of hardcoded `/tmp` paths",
		Severity: SeverityWarning,
		Description: "Hardcoding `/tmp/filename` is vulnerable to symlink attacks and race conditions. " +
			"Use `mktemp` to create unique temporary files safely.",
		Check: checkZC1270,
	})
}

func checkZC1270(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "touch" && name != "cat" && name != "echo" && name != "cp" && name != "mv" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "/tmp/") && !strings.Contains(val, "$") {
			return []Violation{{
				KataID:  "ZC1270",
				Message: "Use `mktemp` instead of hardcoded `" + val + "`. Hardcoded `/tmp` paths are vulnerable to symlink attacks and race conditions.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityWarning,
			}}
		}
	}

	return nil
}

package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1246",
		Title:    "Avoid hardcoded passwords in command arguments",
		Severity: SeverityError,
		Description: "Passing passwords as command arguments exposes them in process lists " +
			"and shell history. Use environment variables or credential files instead.",
		Check: checkZC1246,
	})
}

func checkZC1246(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	name := ident.Value
	if name != "mysql" && name != "psql" && name != "mongosh" && name != "redis-cli" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if strings.HasPrefix(val, "-p") && len(val) > 2 && val != "-p" {
			return []Violation{{
				KataID: "ZC1246",
				Message: "Avoid passing passwords as command arguments — they appear in process lists. " +
					"Use environment variables (e.g., `MYSQL_PWD`) or credential files instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}

	return nil
}

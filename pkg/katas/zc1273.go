package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1273",
		Title:    "Use `grep -q` instead of redirecting grep output to `/dev/null`",
		Severity: SeverityStyle,
		Description: "`grep -q` suppresses output and exits on first match, which is faster and more " +
			"idiomatic than piping or redirecting to `/dev/null`.",
		Check: checkZC1273,
	})
}

func checkZC1273(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "grep" {
		return nil
	}

	hasQuiet := false
	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-q" || val == "--quiet" || val == "--silent" {
			hasQuiet = true
			break
		}
	}

	if hasQuiet {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "/dev/null" {
			return []Violation{{
				KataID:  "ZC1273",
				Message: "Use `grep -q` instead of redirecting to `/dev/null`. It is faster and more idiomatic.",
				Line:    cmd.Token.Line,
				Column:  cmd.Token.Column,
				Level:   SeverityStyle,
			}}
		}
	}

	return nil
}

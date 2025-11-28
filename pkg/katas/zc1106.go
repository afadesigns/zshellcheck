package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1106",
		Title: "Avoid `set -x` in production scripts for sensitive data exposure",
		Description: "Using `set -x` (xtrace) in production environments can expose sensitive " +
			"information, such as API keys or passwords, in logs. While useful for debugging, " +
			"it should be avoided in production. Consider using targeted debugging or secure logging.",
		Check: checkZC1106,
	})
}

func checkZC1106(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name != nil && cmd.Name.String() == "set" {
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			argStr = strings.Trim(argStr, "\"'")
			if strings.HasPrefix(argStr, "-") {
				// Check for -x flag explicitly or combined flags like -eux
				if strings.Contains(argStr, "x") {
					return []Violation{{
						KataID:  "ZC1106",
						Message: "Avoid `set -x` in production scripts to prevent sensitive data exposure.",
						Line:    cmd.Token.Line,
						Column:  cmd.Token.Column,
					}}
				}
			}
		}
	}

	return nil
}

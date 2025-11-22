package katas

import (
    "strings"
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1012",
		Title:       "Use `read -r` to prevent backslash escaping",
		Description: "By default, `read` interprets backslashes as escape characters. " +
			"Use `read -r` to treat backslashes literally, which is usually what you want.",
		Check: checkZC1012,
	})
}

func checkZC1012(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if cmd.Name.String() == "read" {
			hasR := false
			for _, arg := range cmd.Arguments {
				s := arg.String()
                
                // Handle PrefixExpression String() format: "(-r)" -> "-r"
                s = strings.Trim(s, "()")
                
                if len(s) > 0 && s[0] == '-' {
                    if strings.Contains(s, "r") {
                        hasR = true
                        break
                    }
                }
			}

			if !hasR {
				violations = append(violations, Violation{
					KataID:  "ZC1012",
					Message: "Use `read -r` to read input without interpreting backslashes.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
				})
			}
		}
	}

	return violations
}

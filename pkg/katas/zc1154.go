package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1154",
		Title:    "Use `find -exec {} +` instead of `find -exec {} \\;`",
		Severity: SeverityStyle,
		Description: "`find -exec cmd {} \\;` runs cmd once per file. " +
			"`find -exec cmd {} +` batches files into fewer invocations, improving performance.",
		Check: checkZC1154,
	})
}

func checkZC1154(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "find" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-exec" {
			// Check if the exec block ends with \; (not +)
			for j := i + 1; j < len(cmd.Arguments); j++ {
				endVal := cmd.Arguments[j].String()
				if endVal == ";" {
					return []Violation{{
						KataID: "ZC1154",
						Message: "Use `find -exec cmd {} +` instead of `find -exec cmd {} \\;`. " +
							"The `+` form batches files for fewer process invocations.",
						Line:   cmd.Token.Line,
						Column: cmd.Token.Column,
						Level:  SeverityStyle,
					}}
				}
				if endVal == "+" {
					break
				}
			}
		}
	}

	return nil
}

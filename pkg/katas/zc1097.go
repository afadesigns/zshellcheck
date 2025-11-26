package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:    "ZC1097",
		Title: "Declare loop variables as `local` in functions",
		Description: "Loop variables in `for` loops are global by default in Zsh functions. " +
			"Use `local` to scope them to the function before the loop.",
		Check: checkZC1097,
	})
}

func checkZC1097(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}

	violations := []Violation{}
	locals := make(map[string]bool)

	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		// Stop walking into nested function definitions
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}

		// Track local declarations
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			nameStr := cmd.Name.String()
			if nameStr == "local" || nameStr == "typeset" || nameStr == "declare" ||
				nameStr == "integer" || nameStr == "float" || nameStr == "readonly" {
				for _, arg := range cmd.Arguments {
					// Arg can be "x" or "x=1" or "-r"
					argStr := arg.String()
					if len(argStr) > 0 && argStr[0] == '-' {
						continue // Skip options
					}
					// Extract name before '='
					varName := argStr
					for i, c := range argStr {
						if c == '=' {
							varName = argStr[:i]
							break
						}
					}
					locals[varName] = true
				}
			}
		}

		// Check ForLoopStatement
		if forLoop, ok := n.(*ast.ForLoopStatement); ok {
			// Check if Init is an Identifier (loop variable)
			if ident, ok := forLoop.Init.(*ast.Identifier); ok {
				if !locals[ident.Value] {
					violations = append(violations, Violation{
						KataID: "ZC1097",
						Message: "Loop variable '" + ident.Value + "' is used without 'local'. It will be global. " +
							"Use `local " + ident.Value + "` before the loop.",
						Line:   ident.Token.Line,
						Column: ident.Token.Column,
					})
				}
			}
		}

		return true
	})

	return violations
}

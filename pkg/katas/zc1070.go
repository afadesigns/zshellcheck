package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.FunctionDefinitionNode, Kata{
		ID:          "ZC1070",
		Title:       "Use `builtin` or `command` to avoid infinite recursion in wrapper functions",
		Description: "When defining a wrapper function with the same name as a builtin or command (e.g., `cd`), " +
			"calling the command directly inside the function causes infinite recursion. " +
			"Use `builtin cd` or `command cd`.",
		Check: checkZC1070,
	})
}

func checkZC1070(node ast.Node) []Violation {
	funcDef, ok := node.(*ast.FunctionDefinition)
	if !ok {
		return nil
	}

	name := funcDef.Name.String()

	// Only check for common builtins/commands to avoid flagging valid recursive algorithms
	targets := map[string]bool{
		"cd": true, "echo": true, "printf": true, "read": true, "source": true, ".": true,
		"eval": true, "exec": true, "exit": true, "export": true, "kill": true,
		"local": true, "pwd": true, "return": true, "set": true, "shift": true,
		"test": true, "trap": true, "typeset": true, "umask": true, "unset": true, "wait": true,
		"ls": true, "grep": true, "mkdir": true, "rm": true, "mv": true, "cp": true, "git": true,
		"dirs": true, "popd": true, "pushd": true,
	}

	if !targets[name] {
		return nil
	}

	violations := []Violation{}

	// Walk body to find self-calls
	ast.Walk(funcDef.Body, func(n ast.Node) bool {
		// Don't recurse into nested functions (they mask the name)
		if _, ok := n.(*ast.FunctionDefinition); ok && n != funcDef {
			return false
		}
		if _, ok := n.(*ast.FunctionLiteral); ok {
			return false
		}

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			cmdName := cmd.Name.String()
			if cmdName == name {
				// Found self-call.
				// Check if it is "builtin name" or "command name" is handled?
				// SimpleCommand "builtin" with arg "name".
				// But here `cmd.Name` IS `name`.
				// So `builtin cd` -> Name="builtin", Args=["cd"]
				// `cd` -> Name="cd".

				// If Name == function name, it IS a recursive call.
				// Unless it is `builtin` or `command`?
				// If I write `builtin cd`, the parser sees Name="builtin".
				// So if Name matches `name`, it is NOT `builtin` or `command`.

				// Exception: `command` might not be a keyword in parser?
				// `command -v cd` -> Name="command".

				// So if `cmdName == name`, it is a direct call.

				violations = append(violations, Violation{
					KataID: "ZC1070",
					Message: "Recursive call to `" + name + "` inside `" + name + "`. " +
						"Use `builtin " + name + "` or `command " + name + "` to invoke the underlying command.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
				})
			}
		}
		return true
	})

	return violations
}

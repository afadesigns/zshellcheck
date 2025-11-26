package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1103",
		Title: "Suggest `path` array instead of `$PATH` string manipulation",
		Description: "Zsh automatically maps the `$PATH` environment variable to the `$path` array. " +
			"Modifying `$path` is cleaner and less error-prone than manipulating the colon-separated `$PATH` string.",
		Check: checkZC1103,
	})
}

func checkZC1103(node ast.Node) []Violation {
	// This is tricky because we need to check assignments.
	// SimpleCommand handles assignments if they are part of the command (e.g. PATH=... cmd)
	// But typically assignments `PATH=...` are parsed as SimpleCommand with Name="PATH=...".
	// Let's check if the node represents an assignment to PATH.
	
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	
	// In our parser, `PATH=$PATH:/bin` might be parsed as Name="PATH=$PATH:/bin" (if no spaces).
	// Or if it's `export PATH=...`, it's different.
	
	name := cmd.Name.String()
	// Check for assignment `PATH=...`
	// This is a heuristic.
	if len(name) > 5 && name[:5] == "PATH=" {
		return []Violation{{
			KataID:  "ZC1103",
			Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
		}}
	}
	
	// Also check for `export PATH=...`
	if name == "export" {
		for _, arg := range cmd.Arguments {
			argStr := arg.String()
			if len(argStr) > 5 && argStr[:5] == "PATH=" {
				return []Violation{{
					KataID:  "ZC1103",
					Message: "Use the `path` array instead of manually manipulating the `$PATH` string.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
				}}
			}
		}
	}

	return nil
}

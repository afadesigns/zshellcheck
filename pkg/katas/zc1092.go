package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1092",
		Title: "Prefer `print` or `printf` over `echo` in Zsh",
		Description: "In Zsh, `echo` behavior can vary significantly based on options like `BSD_ECHO`. " +
			"`print` is a builtin with consistent behavior and more features. " +
			"For formatted output, `printf` is preferred.",
		Check: checkZC1092,
	})
}

func checkZC1092(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name == nil {
		return nil
	}

	if cmd.Name.String() == "echo" {
		// Check if it's just a simple echo or if flags are involved
		// If flags are used (like -n, -e), print is definitely better.
		// Even without flags, print is idiomatic Zsh.
		
		// We can be slightly lenient and only warn if flags are present OR if it contains backslashes?
		// The prompt suggests "Prefer 'print' over 'echo'". Let's be strict for now as it's "Platinum Standard".
		
		msg := "Prefer `print` over `echo`. `echo` behavior varies. `print` is the Zsh builtin. Especially with flags, `print -n` or `print -r` is more reliable."

		return []Violation{{
			KataID:  "ZC1092",
			Message: msg,
			Line:    cmd.Token.Line,
			Column:  cmd.Token.Column,
		}}
	}

	return nil
}

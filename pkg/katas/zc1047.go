package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:          "ZC1047",
		Title:       "Avoid `sudo` in scripts",
		Description: "Using `sudo` in scripts is generally discouraged. It makes the script interactive and less portable. Run the script as root or use `sudo` to invoke the script.",
		Check:       checkZC1047,
	})
}

func checkZC1047(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is sudo
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "sudo" {
		return []Violation{{
			KataID:  "ZC1047",
			Message: "Avoid `sudo` in scripts. Run the entire script as root if privileges are required.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		}}
	}

	return nil
}

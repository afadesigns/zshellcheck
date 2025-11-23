package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	kata := Kata{
		ID:          "ZC1068",
		Title:       "Use `add-zsh-hook` instead of defining hook functions directly",
		Description: "Defining special functions like `precmd`, `preexec`, `chpwd`, etc. directly overwrites any previously defined hooks. Use `autoload -Uz add-zsh-hook; add-zsh-hook <hook> <function>` to append to the hook list safely.",
		Check:       checkZC1068,
	}
	RegisterKata(ast.FunctionDefinitionNode, kata)
	RegisterKata(ast.FunctionLiteralNode, kata)
}

func checkZC1068(node ast.Node) []Violation {
	var name string
	var tokenLine, tokenCol int

	switch n := node.(type) {
	case *ast.FunctionDefinition:
		name = n.Name.String()
		tokenLine = n.Token.Line
		tokenCol = n.Token.Column
	case *ast.FunctionLiteral:
		name = n.Name
		tokenLine = n.Token.Line
		tokenCol = n.Token.Column
	default:
		return nil
	}

	// List of special hook functions in Zsh
	specialHooks := map[string]bool{
		"precmd":             true,
		"preexec":            true,
		"chpwd":              true,
		"periodic":           true,
		"zshaddhistory":      true,
		"zshexit":            true,
		"zsh_directory_name": true,
	}

	if specialHooks[name] {
		return []Violation{
			{
				KataID:  "ZC1068",
				Message: "Defining `" + name + "` directly overwrites existing hooks. Use `autoload -Uz add-zsh-hook; add-zsh-hook " + name + " my_func` instead.",
				Line:    tokenLine,
				Column:  tokenCol,
			},
		}
	}

	return nil
}

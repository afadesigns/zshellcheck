package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	// Register for Redirection (> file)
	RegisterKata(ast.RedirectionNode, Kata{
		ID:    "ZC1087",
		Title: "Output redirection overwrites input file",
		Description: "Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. " +
			"Use a temporary file or `sponge`.",
		Check: checkZC1087,
	})
	// Register for Pipeline (|) to detect clobbering across pipe
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:    "ZC1087",
		Title: "Output redirection overwrites input file",
		Description: "Redirecting output to a file that is also being read as input causes the file to be truncated before it is read. " +
			"Use a temporary file or `sponge`.",
		Check: checkZC1087,
	})
}

func checkZC1087(node ast.Node) []Violation {
	// Case 1: Redirection (> file)
	if redir, ok := node.(*ast.Redirection); ok {
		// Only care about output truncation: > (and maybe >| in Zsh)
		if redir.Operator != ">" && redir.Operator != ">|" {
			return nil
		}

		outputFile := redir.Right.String()
		inputs := collectInputs(redir.Left)
		
		for _, input := range inputs {
			if input == outputFile {
				return []Violation{
					{
						KataID:  "ZC1087",
						Message: "Output redirection overwrites input file `" + outputFile + "`. The file is truncated before reading.",
						Line:    redir.TokenLiteralNode().Line,
						Column:  redir.TokenLiteralNode().Column,
					},
				}
			}
		}
		return nil
	}

	// Case 2: Pipeline (cmd1 | cmd2)
	if infix, ok := node.(*ast.InfixExpression); ok {
		if infix.Operator != "|" {
			return nil
		}
		
		// Left side inputs
		inputs := collectInputs(infix.Left)
		// Right side outputs
		outputs := collectOutputs(infix.Right)
		
		for _, output := range outputs {
			for _, input := range inputs {
				if input == output {
					return []Violation{
						{
							KataID:  "ZC1087",
							Message: "Output redirection overwrites input file `" + output + "`. The file is truncated before reading.",
							Line:    infix.TokenLiteralNode().Line,
							Column:  infix.TokenLiteralNode().Column,
						},
					}
				}
			}
		}
	}

	return nil
}

func collectInputs(node ast.Node) []string {
	var inputs []string
	
	ast.Walk(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		
		// 1. SimpleCommand arguments
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			for _, arg := range cmd.Arguments {
				inputs = append(inputs, arg.String())
			}
		}
		
		// 2. Input Redirection (<)
		// It seems < might be parsed as Redirection or InfixExpression depending on version/flags?
		// Debug showed *ast.Redirection for <.
		if redir, ok := n.(*ast.Redirection); ok {
			if redir.Operator == "<" {
				inputs = append(inputs, redir.Right.String())
			}
		}

		if infix, ok := n.(*ast.InfixExpression); ok {
			if infix.Operator == "<" {
				inputs = append(inputs, infix.Right.String())
			}
		}
		
		return true
	})
	
	return inputs
}

func collectOutputs(node ast.Node) []string {
	var outputs []string
	ast.Walk(node, func(n ast.Node) bool {
		if n == nil { return true }
		if redir, ok := n.(*ast.Redirection); ok {
			if redir.Operator == ">" || redir.Operator == ">|" {
				outputs = append(outputs, redir.Right.String())
			}
		}
		// Also check if > is Infix?
		if infix, ok := n.(*ast.InfixExpression); ok {
			if infix.Operator == ">" || infix.Operator == ">|" {
				outputs = append(outputs, infix.Right.String())
			}
		}
		return true
	})
	return outputs
}

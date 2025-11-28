package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	// Register for SimpleCommand (to check args)
	RegisterKata(ast.SimpleCommandNode, Kata{
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
	// Case 1: SimpleCommand (checking args for > file)
	if cmd, ok := node.(*ast.SimpleCommand); ok {
		inputs := collectInputs(cmd)
		outputs := collectOutputs(cmd)

		for _, output := range outputs {
			for _, input := range inputs {
				if input == output {
					return []Violation{
						{
							KataID:  "ZC1087",
							Message: "Output redirection overwrites input file `" + output + "`. The file is truncated before reading.",
							Line:    cmd.TokenLiteralNode().Line,
							Column:  cmd.TokenLiteralNode().Column,
						},
					}
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

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			for i := 0; i < len(cmd.Arguments); i++ {
				arg := cmd.Arguments[i].String()
				if arg == "<" {
					if i+1 < len(cmd.Arguments) {
						inputs = append(inputs, cmd.Arguments[i+1].String())
						i++
					}
				} else if arg != ">" && arg != ">>" && arg != ">|" && arg != "&>" {
					// Assume args are inputs unless they are flags
					if len(arg) > 0 && arg[0] != '-' {
						inputs = append(inputs, arg)
					}
				} else {
					// Skip next (output file)
					i++
				}
			}
		}
		return true
	})

	return inputs
}

func collectOutputs(node ast.Node) []string {
	var outputs []string
	ast.Walk(node, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			for i := 0; i < len(cmd.Arguments); i++ {
				arg := cmd.Arguments[i].String()
				// Only output redirection that truncates: > or >|
				// Ignore append >>, &>>, etc. unless we want to catch clobbering there too?
				// Kata description says "truncated". >> does not truncate.
				if arg == ">" || arg == ">|" {
					if i+1 < len(cmd.Arguments) {
						outputs = append(outputs, cmd.Arguments[i+1].String())
						i++
					}
				} else if arg == ">>" || arg == "&>" || arg == "&>>" {
					// Skip operator and file
					i++
				}
			}
		}
		return true
	})
	return outputs
}

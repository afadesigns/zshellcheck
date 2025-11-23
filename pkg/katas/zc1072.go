package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.InfixExpressionNode, Kata{
		ID:          "ZC1072",
		Title:       "Use `awk` instead of `grep | awk`",
		Description: "`grep pattern | awk '{...}'` is inefficient. Use `awk '/pattern/ {...}'` to combine matching and processing in a single process.",
		Check:       checkZC1072,
	})
}

func checkZC1072(node ast.Node) []Violation {
	pipe, ok := node.(*ast.InfixExpression)
	if !ok || pipe.Operator != "|" {
		return nil
	}

	// Check left command is grep
	grepCmd, ok := pipe.Left.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(grepCmd, "grep") {
		return nil
	}

	// Check right command is awk/gawk/mawk
	awkCmd, ok := pipe.Right.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if !isCommandName(awkCmd, "awk") && !isCommandName(awkCmd, "gawk") && !isCommandName(awkCmd, "mawk") {
		return nil
	}

	// Check grep flags. If flags are complex (like -r, -v, -l), we might skip warning.
	// But `grep | awk` is almost always replaceable.
	// Only if grep does something awk can't easily do (like -r recursive search) should we allow it?
	// Awk doesn't do recursive directory search by default.
	// So if grep has `-r` or `-R`, it's valid.
	
	if hasRecursiveFlag(grepCmd) {
		return nil
	}

	return []Violation{{
		KataID:  "ZC1072",
		Message: "Use `awk '/pattern/ {...}'` instead of `grep pattern | awk '{...}'` to avoid a pipeline.",
		Line:    pipe.TokenLiteralNode().Line,
		Column:  pipe.TokenLiteralNode().Column,
	}}
}

func hasRecursiveFlag(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		// Check PrefixExpression
		if prefix, ok := arg.(*ast.PrefixExpression); ok && prefix.Operator == "-" {
			if ident, ok := prefix.Right.(*ast.Identifier); ok {
				val := ident.Value
				if val == "r" || val == "R" || val == "-recursive" {
					return true
				}
				// Combined flags e.g. -rn
				if strings.Contains(val, "r") || strings.Contains(val, "R") {
					return true
				}
			}
		}
		// Check StringLiteral
		if str, ok := arg.(*ast.StringLiteral); ok {
			val := str.Value
			// Remove quotes
			if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') {
				val = val[1 : len(val)-1]
			}
			if val == "-r" || val == "-R" || val == "--recursive" {
				return true
			}
		}
	}
	return false
}

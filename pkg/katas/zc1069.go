// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1069",
		Title: "Avoid `local` outside of functions",
		Description: "The `local` builtin can only be used inside functions. " +
			"Using it in the global scope causes an error.",
		Severity: SeverityInfo,
		Check:    checkZC1069,
		Fix:      fixZC1069,
	})
}

// fixZC1069 rewrites `local` to `typeset` when used at file scope.
// `typeset` works in both function and global contexts, so the
// rewrite is safe wherever the detector fires. Single-edit name
// swap at the violation column. Idempotent — a re-run sees
// `typeset`, not `local`. Defensive byte-match guard.
func fixZC1069(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 || off+len("local") > len(source) {
		return nil
	}
	if string(source[off:off+len("local")]) != "local" {
		return nil
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("local"),
		Replace: "typeset",
	}}
}

func checkZC1069(node ast.Node) []Violation {
	program, ok := node.(*ast.Program)
	if !ok {
		return nil
	}

	violations := []Violation{}

	// Helper to walk and track scope
	var walk func(n ast.Node, inFunction bool)
	walk = func(n ast.Node, inFunction bool) {
		if n == nil {
			return
		}

		// Check for local usage
		if cmd, ok := n.(*ast.SimpleCommand); ok {
			if name, ok := cmd.Name.(*ast.Identifier); ok && (name.Value == "local" || name.Value == "typeset") {
				if name.Value == "local" && !inFunction {
					violations = append(violations, Violation{
						KataID: "ZC1069",
						Message: "`local` can only be used inside functions. " +
							"Use `typeset`, `declare`, or just assignment for global variables.",
						Line:   name.Token.Line,
						Column: name.Token.Column,
						Level:  SeverityInfo,
					})
				}
			}
		}

		// Determine if we are entering a function
		// We handle scope change in the switch below explicitly for definition nodes.

		switch t := n.(type) {
		case *ast.Program:
			for _, s := range t.Statements {
				walk(s, inFunction)
			}
		case *ast.BlockStatement:
			for _, s := range t.Statements {
				walk(s, inFunction)
			}
		case *ast.IfStatement:
			if t.Condition != nil {
				walk(t.Condition, inFunction)
			}
			if t.Consequence != nil {
				walk(t.Consequence, inFunction)
			}
			if t.Alternative != nil {
				walk(t.Alternative, inFunction)
			}
		case *ast.ForLoopStatement:
			if t.Init != nil {
				walk(t.Init, inFunction)
			}
			if t.Condition != nil {
				walk(t.Condition, inFunction)
			}
			if t.Post != nil {
				walk(t.Post, inFunction)
			}
			for _, item := range t.Items {
				walk(item, inFunction)
			}
			if t.Body != nil {
				walk(t.Body, inFunction)
			}
		case *ast.WhileLoopStatement:
			if t.Condition != nil {
				walk(t.Condition, inFunction)
			}
			if t.Body != nil {
				walk(t.Body, inFunction)
			}
		case *ast.FunctionDefinition:
			if t.Name != nil {
				walk(t.Name, inFunction)
			}
			if t.Body != nil {
				walk(t.Body, true)
			}
		case *ast.FunctionLiteral:
			for _, p := range t.Params {
				walk(p, inFunction)
			}
			if t.Body != nil {
				walk(t.Body, true)
			}
		case *ast.SimpleCommand:
			walk(t.Name, inFunction)
			for _, arg := range t.Arguments {
				walk(arg, inFunction)
			}
		case *ast.ExpressionStatement:
			walk(t.Expression, inFunction)
		case *ast.InfixExpression:
			walk(t.Left, inFunction)
			walk(t.Right, inFunction)
		case *ast.PrefixExpression:
			walk(t.Right, inFunction)
		case *ast.PostfixExpression:
			walk(t.Left, inFunction)
		case *ast.GroupedExpression:
			walk(t.Expression, inFunction)
		case *ast.CaseStatement:
			walk(t.Value, inFunction)
			for _, clause := range t.Clauses {
				for _, p := range clause.Patterns {
					walk(p, inFunction)
				}
				walk(clause.Body, inFunction)
			}
		case *ast.ConcatenatedExpression:
			for _, p := range t.Parts {
				walk(p, inFunction)
			}
		case *ast.CommandSubstitution:
			walk(t.Command, inFunction)
		case *ast.DollarParenExpression:
			walk(t.Command, inFunction)
		case *ast.Subshell:
			walk(t.Command, inFunction)
		}
	}

	walk(program, false)

	return violations
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1088",
		Title: "Subshell isolates state changes",
		Description: "Commands inside `( ... )` run in a subshell. " +
			"State changes like `cd`, `export`, or variable assignments are lost when the subshell exits. " +
			"Use `{ ... }` for grouping if you want to preserve state changes.",
		Check: checkZC1088,
	})
}

func checkZC1088(node ast.Node) []Violation {
	// We perform a context-aware traversal from the root
	v := &zc1088Visitor{violations: []Violation{}}
	v.traverse(node, false)
	return v.violations
}

type zc1088Visitor struct {
	violations []Violation
}

func (v *zc1088Visitor) traverse(node ast.Node, expectsStatus bool) {
	if node == nil {
		return
	}
	
	// Handle typed nil interfaces
	if t, ok := node.(*ast.BlockStatement); ok && t == nil { return }
	if t, ok := node.(*ast.IfStatement); ok && t == nil { return }
	if t, ok := node.(*ast.WhileLoopStatement); ok && t == nil { return }
	if t, ok := node.(*ast.ExpressionStatement); ok && t == nil { return }
	if t, ok := node.(*ast.InfixExpression); ok && t == nil { return }
	if t, ok := node.(*ast.PrefixExpression); ok && t == nil { return }
	if t, ok := node.(*ast.GroupedExpression); ok && t == nil { return }
	if t, ok := node.(*ast.Program); ok && t == nil { return }
	if t, ok := node.(*ast.Subshell); ok && t == nil { return }

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			v.traverse(stmt, false)
		}
	case *ast.Subshell:
		// ( ... ) in statement context
		if !expectsStatus {
			v.checkSubshell(n)
		}
		// Traverse content (BlockStatement) with expectsStatus=false for children
		v.traverse(n.Block, false)

	case *ast.BlockStatement:
		// Normal block { ... } or implicit (e.g. If Condition)
		// If this block expects status, only the LAST statement expects status.
		for i, stmt := range n.Statements {
			isLast := i == len(n.Statements)-1
			v.traverse(stmt, expectsStatus && isLast)
		}
	case *ast.IfStatement:
		v.traverse(n.Condition, true)
		v.traverse(n.Consequence, false)
		v.traverse(n.Alternative, false)
	case *ast.WhileLoopStatement:
		v.traverse(n.Condition, true)
		v.traverse(n.Body, false)
	case *ast.ExpressionStatement:
		v.traverse(n.Expression, expectsStatus)
	case *ast.InfixExpression:
		if n.Operator == "&&" || n.Operator == "||" {
			v.traverse(n.Left, true)
			v.traverse(n.Right, true)
		} else {
			v.traverse(n.Left, false)
			v.traverse(n.Right, false)
		}
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			v.traverse(n.Right, true)
		} else {
			v.traverse(n.Right, false)
		}
	case *ast.GroupedExpression:
		// ( ... ) expression context
		if !expectsStatus {
			v.checkGroupedExpression(n)
		}
		v.traverse(n.Exp, false)
	default:
		// Do nothing for other nodes
	}
}

func (v *zc1088Visitor) checkSubshell(sub *ast.Subshell) {
	if v.isStateChangeOnly(sub.Block) {
		v.violations = append(v.violations, Violation{
			KataID:  "ZC1088",
			Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
			Line:    sub.TokenLiteralNode().Line,
			Column:  sub.TokenLiteralNode().Column,
		})
	}
}

func (v *zc1088Visitor) checkGroupedExpression(group *ast.GroupedExpression) {
	// Similar logic for GroupedExpression if it wraps state changes
	// But GroupedExpression wraps Expression.
	// SimpleCommand is Expression.
	if v.isStateChangeOnly(group.Exp) {
		v.violations = append(v.violations, Violation{
			KataID:  "ZC1088",
			Message: "Subshell `( ... )` isolates state changes. The changes (e.g. `cd`, variable assignment) will be lost. Use `{ ... }` to preserve them, or add commands that use the changed state.",
			Line:    group.TokenLiteralNode().Line,
			Column:  group.TokenLiteralNode().Column,
		})
	}
}

func (v *zc1088Visitor) isStateChangeOnly(node ast.Node) bool {
	hasStateChange := false
	hasSideEffect := false
	
	ast.Walk(node, func(n ast.Node) bool {
		if n == nil || n == node { return true }
		
		if hasSideEffect { return false } // Optimization

		if cmd, ok := n.(*ast.SimpleCommand); ok {
			name := cmd.Name.String()
			if isStateChanger(name) {
				hasStateChange = true
			} else {
				hasSideEffect = true
			}
		} else if infix, ok := n.(*ast.InfixExpression); ok {
			if infix.Operator == "=" {
				hasStateChange = true
			} else {
				// e.g. a && b. If a is state change, b might be side effect.
				// We just traverse.
			}
		} else if _, ok := n.(*ast.IfStatement); ok {
			hasSideEffect = true
		} else if _, ok := n.(*ast.ForLoopStatement); ok {
			hasSideEffect = true
		} else if _, ok := n.(*ast.WhileLoopStatement); ok {
			hasSideEffect = true
		}
		return true
	})
	
	return hasStateChange && !hasSideEffect
}

func isStateChanger(name string) bool {
	switch name {
	case "cd", "export", "unset", "alias", "unalias", "declare", "typeset", "local", "shift":
		return true
	}
	return false
}

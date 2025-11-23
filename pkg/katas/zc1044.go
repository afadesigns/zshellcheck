package katas

import (
	"reflect"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.ProgramNode, Kata{
		ID:    "ZC1044",
		Title: "Check for unchecked `cd` commands",
		Description: "`cd` failures should be handled to avoid executing commands in the wrong directory. " +
			"Use `cd ... || return` (or `exit`).",
		Check: checkZC1044,
	})
}

func checkZC1044(node ast.Node) []Violation {
	// We only run on ProgramNode, but we do a full context-aware traversal.
	// Since main.go calls Check on every node, we might get called for Program.
	// But we might also want to ensure we don't double check if we traverse children.
	// Actually, if we register ONLY for ProgramNode, we are called once per file.
	// BUT, standard ast.Walk visits all nodes and calls Check.
	// If Check returns violations, they are added.
	// So if we implement a full walker here, it works.

	violations := []Violation{}

	walkZC1044(node, false, &violations)

	return violations
}

func walkZC1044(node ast.Node, isChecked bool, violations *[]Violation) {
	if node == nil {
		return
	}

	// Handle nil pointers inside interface
	if v := reflect.ValueOf(node); v.Kind() == reflect.Ptr && v.IsNil() {
		return
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			walkZC1044(stmt, false, violations)
		}
	case *ast.BlockStatement:
		for i, stmt := range n.Statements {
			check := false
			// Only the last statement in a block inherits the "checked" status of the block itself
			// (e.g. the last command in an `if` condition determines the result).
			// Previous statements are unchecked unless they have their own handling.
			if isChecked && i == len(n.Statements)-1 {
				check = true
			}
			walkZC1044(stmt, check, violations)
		}
	case *ast.GroupedExpression:
		walkZC1044(n.Exp, isChecked, violations)
	case *ast.IfStatement:
		walkZC1044(n.Condition, true, violations) // Condition is checked
		walkZC1044(n.Consequence, false, violations)
		walkZC1044(n.Alternative, false, violations)
	case *ast.WhileLoopStatement:
		walkZC1044(n.Condition, true, violations) // Condition is checked
		walkZC1044(n.Body, false, violations)
	case *ast.ForLoopStatement:
		walkZC1044(n.Init, false, violations)
		walkZC1044(n.Condition, true, violations) // Loop condition checked (arithmetic)
		walkZC1044(n.Post, false, violations)
		for _, item := range n.Items {
			walkZC1044(item, false, violations)
		}
		walkZC1044(n.Body, false, violations)
	case *ast.ExpressionStatement:
		walkZC1044(n.Expression, isChecked, violations)
	case *ast.InfixExpression:
		switch n.Operator {
		case "||":
			walkZC1044(n.Left, true, violations) // Left checked by Right
			walkZC1044(n.Right, isChecked, violations)
		case "&&":
			walkZC1044(n.Left, isChecked, violations) // Left inherits check
			walkZC1044(n.Right, isChecked, violations)
		default:
			walkZC1044(n.Left, false, violations)
			walkZC1044(n.Right, false, violations)
		}
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			walkZC1044(n.Right, true, violations) // Negation checks it
		} else {
			walkZC1044(n.Right, false, violations)
		}
	case *ast.FunctionDefinition:
		walkZC1044(n.Body, false, violations)
	case *ast.SimpleCommand:
		checkCommandZC1044(n, isChecked, violations)
		// Also walk args? Arguments might contain subshells etc.
		for _, arg := range n.Arguments {
			walkZC1044(arg, false, violations)
		}
	case *ast.CommandSubstitution:
		walkZC1044(n.Command, false, violations) // Inner command starts unchecked

	// Recursion for other nodes
	default:
		// Use generic walk or manual?
		// Generic walk doesn't pass state.
		// We must implement all nodes that contain statements/expressions.
		// ...
	}
}

func checkCommandZC1044(cmd *ast.SimpleCommand, isChecked bool, violations *[]Violation) {
	if isChecked {
		return
	}
	if name, ok := cmd.Name.(*ast.Identifier); ok && name.Value == "cd" {
		*violations = append(*violations, Violation{
			KataID:  "ZC1044",
			Message: "Use `cd ... || return` (or `exit`) in case cd fails.",
			Line:    name.Token.Line,
			Column:  name.Token.Column,
		})
	}
}

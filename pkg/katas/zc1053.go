package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.IfStatementNode, Kata{
		ID:    "ZC1053",
		Title: "Silence `grep` output in conditions",
		Description: "Using `grep` in a condition prints matches to stdout. " +
			"Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.",
		Check: checkZC1053,
	})
	RegisterKata(ast.WhileLoopStatementNode, Kata{
		ID:    "ZC1053",
		Title: "Silence `grep` output in conditions",
		Description: "Using `grep` in a condition prints matches to stdout. " +
			"Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.",
		Check: checkZC1053,
	})
}

func checkZC1053(node ast.Node) []Violation {
	violations := []Violation{}

	var condition ast.Node

	switch n := node.(type) {
	case *ast.IfStatement:
		condition = n.Condition
	case *ast.WhileLoopStatement:
		condition = n.Condition
	default:
		return nil
	}

	if condition == nil {
		return nil
	}

	walkZC1053(condition, false, &violations)

	return violations
}

func walkZC1053(node ast.Node, isSilenced bool, violations *[]Violation) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			// In a block (condition), usually all statements execute, but only the last one's exit code
			// matters for the condition?
			// No, `if cmd1; cmd2; then`. Both execute. `cmd1` prints. `cmd2` prints.
			// Should we check ALL commands in condition?
			// Yes, because `grep` printing in a condition is usually unwanted noise.
			walkZC1053(stmt, isSilenced, violations)
		}
	case *ast.ExpressionStatement:
		walkZC1053(n.Expression, isSilenced, violations)
	case *ast.InfixExpression:
		if n.Operator == "|" {
			// Left side of pipe is silenced (stdout goes to pipe)
			walkZC1053(n.Left, true, violations)
			// Right side inherits current silence state
			walkZC1053(n.Right, isSilenced, violations)
		} else {
			// &&, || etc. inherit state
			walkZC1053(n.Left, isSilenced, violations)
			walkZC1053(n.Right, isSilenced, violations)
		}
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			walkZC1053(n.Right, isSilenced, violations)
		}
	case *ast.Redirection:
		// Check if redirection silences stdout
		newSilenced := isSilenced
		if n.Operator == ">" || n.Operator == ">>" || n.Operator == "&>" {
			// Check if Right is /dev/null
			if isDevNull(n.Right) {
				newSilenced = true
			}
		}
		walkZC1053(n.Left, newSilenced, violations)
	case *ast.SimpleCommand:
		checkCommandZC1053(n, isSilenced, violations)
	case *ast.GroupedExpression:
		walkZC1053(n.Exp, isSilenced, violations)
		// case *ast.Subshell:
		// Handled by BlockStatement
	}
}

func checkCommandZC1053(cmd *ast.SimpleCommand, isSilenced bool, violations *[]Violation) {
	if isSilenced {
		return
	}

	if name, ok := cmd.Name.(*ast.Identifier); ok {
		if name.Value == "grep" || name.Value == "egrep" || name.Value == "fgrep" || name.Value == "zgrep" {
			// Check args for -q, --quiet, --silent
			hasQuiet := false
			for _, arg := range cmd.Arguments {
				if str, ok := arg.(*ast.StringLiteral); ok {
					if str.Value == "-q" || str.Value == "--quiet" || str.Value == "--silent" {
						hasQuiet = true
						break
					}
				} else if prefix, ok := arg.(*ast.PrefixExpression); ok && prefix.Operator == "-" {
					if ident, ok := prefix.Right.(*ast.Identifier); ok {
						if strings.Contains(ident.Value, "q") { // e.g. -rq
							hasQuiet = true
							break
						}
					}
				}
			}

			if !hasQuiet {
				*violations = append(*violations, Violation{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
				})
			}
		}
	}
}

func isDevNull(node ast.Node) bool {
	val := getStringValueZC1053(node)
	// Remove quotes
	if len(val) >= 2 && (val[0] == '"' || val[0] == '\'') {
		val = val[1 : len(val)-1]
	}
	return val == "/dev/null"
}

func getStringValueZC1053(node ast.Node) string {
	switch n := node.(type) {
	case *ast.StringLiteral:
		return n.Value
	case *ast.ConcatenatedExpression:
		var sb strings.Builder
		for _, p := range n.Parts {
			sb.WriteString(getStringValueZC1053(p))
		}
		return sb.String()
	}
	return ""
}

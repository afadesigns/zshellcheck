package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
	"github.com/afadesigns/zshellcheck/pkg/token"
)

func init() {
	kata := Kata{
		ID:    "ZC1004",
		Title: "Use `return` instead of `exit` in functions",
		Description: "Using `exit` in a function terminates the entire shell, which is often unintended " +
			"in interactive sessions or sourced scripts. Use `return` to exit the function.",
		Severity: SeverityWarning,
		Check:    checkZC1004,
		Fix:      fixZC1004,
	}
	RegisterKata(ast.FunctionDefinitionNode, kata)
	RegisterKata(ast.FunctionLiteralNode, kata)
}

// fixZC1004 rewrites `exit` to `return` at the command-name position
// inside a function body. Arguments (the exit/return code) stay
// unchanged — `exit 1` becomes `return 1` with the `1` byte-identical.
// The violation's Line/Column already point at the command name.
func fixZC1004(node ast.Node, v Violation, source []byte) []FixEdit {
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  len("exit"),
		Replace: "return",
	}}
}

func checkZC1004(node ast.Node) []Violation {
	var body ast.Statement

	switch n := node.(type) {
	case *ast.FunctionDefinition:
		body = n.Body
	case *ast.FunctionLiteral:
		body = n.Body
	default:
		return nil
	}

	violations := []Violation{}

	ast.Walk(body, func(n ast.Node) bool {
		// Stop traversal at subshell boundaries where exit is safe/scoped
		switch t := n.(type) {
		case *ast.GroupedExpression: // ( ... )
			return false
		case *ast.Subshell: // ( ... ) as subshell
			return false
		case *ast.CommandSubstitution: // ` ... `
			return false
		case *ast.DollarParenExpression: // $( ... )
			return false
		case *ast.BlockStatement:
			if t.Token.Type == token.LPAREN { // ( ... ) as a statement block
				return false
			}
		}

		// Match both SimpleCommand (`exit 1`) and bare Identifier
		// wrapped in an ExpressionStatement (`exit` with no args —
		// the parser folds zero-arg command invocations into a plain
		// identifier expression rather than a SimpleCommand).
		switch sn := n.(type) {
		case *ast.SimpleCommand:
			if sn.Name != nil && sn.Name.String() == "exit" {
				violations = append(violations, Violation{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    sn.Token.Line,
					Column:  sn.Token.Column,
					Level:   SeverityWarning,
				})
			}
			// Don't descend — the Name Identifier would otherwise
			// double-count as a bare-`exit` hit below.
			return false
		case *ast.Identifier:
			if sn.Value == "exit" {
				violations = append(violations, Violation{
					KataID:  "ZC1004",
					Message: "Use `return` instead of `exit` in functions to avoid killing the shell.",
					Line:    sn.Token.Line,
					Column:  sn.Token.Column,
					Level:   SeverityWarning,
				})
			}
		}
		return true
	})

	return violations
}

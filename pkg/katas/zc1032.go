package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.LetStatementNode, Kata{
		ID:    "ZC1032",
		Title: "Use `((...))` for C-style incrementing",
		Description: "Instead of `let i=i+1` or `let i=i-1`, you can use the more concise and idiomatic " +
			"C-style increment `(( i++ ))` / decrement `(( i-- ))` in Zsh.",
		Severity: SeverityStyle,
		Check:    checkZC1032,
		Fix:      fixZC1032,
	})
}

// zc1032Op classifies the let-statement value as either an increment
// (`name + 1`) or decrement (`name - 1`). Returns the C-style suffix
// (`++` / `--`) and true on match; empty string and false otherwise.
//
// The Zsh lexer treats `-` as an identifier byte, so `let i=i-1`
// parses as a single Identifier value `i-1` rather than an
// InfixExpression. The detector handles both shapes:
//   - InfixExpression: `i + 1` (operator `+`, integer literal 1).
//   - Identifier with `NAME-1` literal: synthetic decrement form.
func zc1032Op(stmt *ast.LetStatement) (string, bool) {
	if stmt == nil || stmt.Name == nil {
		return "", false
	}
	if infix, ok := stmt.Value.(*ast.InfixExpression); ok {
		leftIdent, ok := infix.Left.(*ast.Identifier)
		if !ok {
			return "", false
		}
		rightInt, ok := infix.Right.(*ast.IntegerLiteral)
		if !ok {
			return "", false
		}
		if stmt.Name.Value != leftIdent.Value || rightInt.Value != 1 {
			return "", false
		}
		switch infix.Operator {
		case "+":
			return "++", true
		case "-":
			return "--", true
		}
		return "", false
	}
	if ident, ok := stmt.Value.(*ast.Identifier); ok {
		want := stmt.Name.Value + "-1"
		if ident.Value == want {
			return "--", true
		}
	}
	return "", false
}

// fixZC1032 rewrites `let NAME=NAME+1` / `let NAME=NAME-1` into the
// C-style `(( NAME++ ))` / `(( NAME-- ))` arithmetic form. The
// replacement spans from the `let` keyword to the end of the logical
// line (first `;`, `\n`, or EOF). Bails when the AST does not match
// the increment/decrement shape so the fix stays unambiguous and
// idempotent on a re-run (the rewritten form is no longer a
// LetStatement).
func fixZC1032(node ast.Node, v Violation, source []byte) []FixEdit {
	stmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	suffix, ok := zc1032Op(stmt)
	if !ok {
		return nil
	}
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 {
		return nil
	}
	end := start
	for end < len(source) {
		b := source[end]
		if b == '\n' || b == ';' {
			break
		}
		end++
	}
	return []FixEdit{{
		Line:    v.Line,
		Column:  v.Column,
		Length:  end - start,
		Replace: "(( " + stmt.Name.Value + suffix + " ))",
	}}
}

func checkZC1032(node ast.Node) []Violation {
	letStmt, ok := node.(*ast.LetStatement)
	if !ok {
		return nil
	}
	suffix, ok := zc1032Op(letStmt)
	if !ok {
		return nil
	}
	msg := "Use `(( " + letStmt.Name.Value + suffix + " ))` for C-style "
	if suffix == "++" {
		msg += "incrementing instead of `let " + letStmt.Name.Value + "=" + letStmt.Name.Value + "+1`."
	} else {
		msg += "decrementing instead of `let " + letStmt.Name.Value + "=" + letStmt.Name.Value + "-1`."
	}
	return []Violation{{
		KataID:  "ZC1032",
		Message: msg,
		Line:    letStmt.Token.Line,
		Column:  letStmt.Token.Column,
		Level:   SeverityStyle,
	}}
}

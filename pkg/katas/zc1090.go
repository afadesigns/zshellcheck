package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.DoubleBracketExpressionNode, Kata{
		ID:    "ZC1090",
		Title: "Quoted regex pattern in `=~`",
		Description: "Quoting the pattern on the right side of `=~` forces literal string matching in Zsh/Bash. " +
			"Regex metacharacters inside quotes will be matched literally. " +
			"Remove quotes to enable regex matching, or use `==` for literal string comparison.",
		Check: checkZC1090,
	})
}

func checkZC1090(node ast.Node) []Violation {
	expr, ok := node.(*ast.DoubleBracketExpression)
	if !ok {
		return nil
	}

	var violations []Violation

	for _, e := range expr.Elements {
		infix, ok := e.(*ast.InfixExpression)
		if !ok {
			continue
		}

		if infix.Operator != "=~" {
			continue
		}

		// Check Right operand
		checkOperand(infix.Right, infix, &violations)
	}

	return violations
}

func checkOperand(node ast.Expression, infix *ast.InfixExpression, violations *[]Violation) {
	switch n := node.(type) {
	case *ast.StringLiteral:
		if containsRegexMeta(n.Value) {
			*violations = append(*violations, Violation{
				KataID:  "ZC1090",
				Message: "Quoted regex pattern matches literally. Remove quotes to enable regex matching.",
				Line:    n.TokenLiteralNode().Line,
				Column:  n.TokenLiteralNode().Column,
			})
		} else {
			// If no metachars, it's a literal match. Recommend == if it looks like they wanted string match?
			// But `=~ "foo"` works. Just weird.
			// We focus on BROKEN regex.
		}
	case *ast.ConcatenatedExpression:
		for _, part := range n.Parts {
			if sl, ok := part.(*ast.StringLiteral); ok {
				if containsRegexMeta(sl.Value) {
					*violations = append(*violations, Violation{
						KataID:  "ZC1090",
						Message: "Quoted regex pattern matches literally. Remove quotes from the regex part.",
						Line:    sl.TokenLiteralNode().Line,
						Column:  sl.TokenLiteralNode().Column,
					})
					return // One violation per expression is enough
				}
			}
		}
	}
}

func containsRegexMeta(s string) bool {
	// Check for regex metacharacters that are likely intended as regex but broken by quotes.
	// ^ $ * + ? [ ( |
	// We exclude . because it's common in text.
	// We exclude $ because it's used for variables (and my parser keeps it in StringLiteral?).
	// Wait, "$var" literal is "$var" or "value"?
	// Parser stores raw literal including quotes usually?
	// Lexer `readString` returns content WITH quotes.
	// So `s` includes quotes!
	// `containsRegexMeta` should check INSIDE quotes.

	if len(s) < 2 {
		return false
	}
	// Strip quotes
	content := s[1 : len(s)-1]

	for _, char := range content {
		switch char {
		case '^', '*', '+', '?', '[', '(', '|':
			return true
			// case '$':
			// 	// $ might be variable. Don't flag.
			// case '.':
			//  // . is common. Don't flag "file.txt".
		}
	}
	// Check for $ at end? `foo$` -> regex end anchor.
	if strings.HasSuffix(content, "$") {
		// But `price$` might be text.
		// `^` is stronger indicator.
		// `.*` is strong.
	}
	return false
}

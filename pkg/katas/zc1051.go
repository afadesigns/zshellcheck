package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1051",
		Title: "Quote variables in `rm` to avoid globbing",
		Description: "`rm $VAR` is dangerous if `$VAR` contains spaces or glob characters. " +
			"Quote the variable (`rm \"$VAR\"`) to ensure safe deletion.",
		Severity: SeverityWarning,
		Check:    checkZC1051,
		Fix:      fixZC1051,
	})
}

// fixZC1051 wraps an unquoted `$VAR` argument in double-quotes.
// Two edits: one `"` before the arg and one after. Arg span is
// measured from source — we scan forward from the arg's token
// position until the first unescaped whitespace / delimiter,
// honouring `{…}` / `[…]` / `(…)` nesting so expansions like
// `${var[1]}`, `$(cmd)`, `${arr[@]}` stay whole.
func fixZC1051(_ ast.Node, v Violation, source []byte) []FixEdit {
	start := LineColToByteOffset(source, v.Line, v.Column)
	if start < 0 || start >= len(source) {
		return nil
	}
	argLen := unquotedArgLen(source, start)
	if argLen == 0 {
		return nil
	}
	endLine, endCol := offsetLineColZC1051(source, start+argLen)
	if endLine < 0 {
		return nil
	}
	return []FixEdit{
		{Line: v.Line, Column: v.Column, Length: 0, Replace: `"`},
		{Line: endLine, Column: endCol, Length: 0, Replace: `"`},
	}
}

// unquotedArgLen returns the byte length of a shell-word starting
// at offset. Honours brace / paren / bracket nesting so
// `${arr[$i]}` and `$(cmd (sub))` stay whole, and stops on the
// first top-level delimiter.
func unquotedArgLen(source []byte, offset int) int {
	if offset < 0 || offset >= len(source) {
		return 0
	}
	n := 0
	braceDepth := 0
	parenDepth := 0
	bracketDepth := 0
	for offset+n < len(source) {
		c := source[offset+n]
		if braceDepth == 0 && parenDepth == 0 && bracketDepth == 0 {
			switch c {
			case ' ', '\t', '\n', ';', '|', '&', '>', '<':
				return n
			}
			if c == ')' {
				return n
			}
		}
		switch c {
		case '{':
			braceDepth++
		case '}':
			if braceDepth > 0 {
				braceDepth--
			}
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		}
		n++
	}
	return n
}

func offsetLineColZC1051(source []byte, offset int) (int, int) {
	if offset < 0 || offset > len(source) {
		return -1, -1
	}
	line := 1
	col := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
			col = 1
			continue
		}
		col++
	}
	return line, col
}

func checkZC1051(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	// Check if command is rm
	if name, ok := cmd.Name.(*ast.Identifier); !ok || name.Value != "rm" {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		isUnquoted := false

		switch n := arg.(type) {
		case *ast.Identifier:
			// $VAR
			if len(n.Value) > 0 && n.Value[0] == '$' {
				isUnquoted = true
			}
		case *ast.PrefixExpression:
			// $var (if parsed as prefix)
			if n.Operator == "$" {
				isUnquoted = true
			}
		case *ast.ArrayAccess:
			// ${var[...]} unquoted
			// Zsh DOES NOT split unquoted variable expansions by default!
			// BUT it DOES glob them.
			// `rm $var`. If var="a b", it deletes "a b" (one file).
			// If var="*", it expands to all files.
			// So checking for globbing safety is key.
			// `rm \"$var\"` prevents globbing.
			isUnquoted = true
		case *ast.DollarParenExpression:
			// $(...)
			isUnquoted = true
		}

		if isUnquoted {
			violations = append(violations, Violation{
				KataID:  "ZC1051",
				Message: "Unquoted variable in `rm`. Quote it to prevent globbing (e.g. `rm \"$VAR\"`).",
				Line:    arg.TokenLiteralNode().Line,
				Column:  arg.TokenLiteralNode().Column,
				Level:   SeverityWarning,
			})
		}
	}

	return violations
}

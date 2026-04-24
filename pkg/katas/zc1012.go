package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1012",
		Title: "Use `read -r` to prevent backslash escaping",
		Description: "By default, `read` interprets backslashes as escape characters. " +
			"Use `read -r` to treat backslashes literally, which is usually what you want.",
		Severity: SeverityStyle,
		Check:    checkZC1012,
		Fix:      fixZC1012,
	})
}

// fixZC1012 inserts ` -r` directly after the `read` command name.
// Existing flags are left untouched (`read -p "x" VAR` becomes
// `read -r -p "x" VAR`) so the fix is order-preserving and idempotent
// on a second pass (the re-parse will see `-r` and detection won't
// fire).
func fixZC1012(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name.String() != "read" {
		return nil
	}
	nameOffset := LineColToByteOffset(source, v.Line, v.Column)
	if nameOffset < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOffset)
	if nameLen != len("read") {
		return nil
	}
	insertAt := nameOffset + nameLen
	insertLine, insertCol := byteOffsetToLineColZC1012(source, insertAt)
	if insertLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insertLine,
		Column:  insertCol,
		Length:  0,
		Replace: " -r",
	}}
}

// byteOffsetToLineColZC1012 converts a byte offset to a 1-based
// (line, column). Kept kata-local to avoid exposing a shared helper
// that the rest of the package does not yet need.
func byteOffsetToLineColZC1012(source []byte, offset int) (int, int) {
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

func checkZC1012(node ast.Node) []Violation {
	violations := []Violation{}

	if cmd, ok := node.(*ast.SimpleCommand); ok {
		if cmd.Name.String() == "read" {
			hasR := false
			for _, arg := range cmd.Arguments {
				s := arg.String()

				// Handle PrefixExpression String() format: "(-r)" -> "-r"
				s = strings.Trim(s, "()")

				if len(s) > 0 && s[0] == '-' {
					if strings.Contains(s, "r") {
						hasR = true
						break
					}
				}
			}

			if !hasR {
				violations = append(violations, Violation{
					KataID:  "ZC1012",
					Message: "Use `read -r` to read input without interpreting backslashes.",
					Line:    cmd.Token.Line,
					Column:  cmd.Token.Column,
					Level:   SeverityStyle,
				})
			}
		}
	}

	return violations
}

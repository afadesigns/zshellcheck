// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
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
		Severity: SeverityStyle,
		Check:    checkZC1053,
		Fix:      fixZC1053,
	})
	RegisterKata(ast.WhileLoopStatementNode, Kata{
		ID:    "ZC1053",
		Title: "Silence `grep` output in conditions",
		Description: "Using `grep` in a condition prints matches to stdout. " +
			"Use `grep -q` (or `> /dev/null`) to silence output if you only care about the exit code.",
		Severity: SeverityStyle,
		Check:    checkZC1053,
		Fix:      fixZC1053,
	})
}

// fixZC1053 inserts ` -q` directly after the grep / egrep / fgrep /
// zgrep command name reported at the violation column. Idempotent —
// once `-q` is present the detector's hasQuiet check short-circuits
// and the kata no longer fires. Defensive byte-match guard refuses
// to insert unless the source at the offset is one of the recognised
// grep variants followed by whitespace.
func fixZC1053(_ ast.Node, v Violation, source []byte) []FixEdit {
	off := LineColToByteOffset(source, v.Line, v.Column)
	if off < 0 {
		return nil
	}
	var name string
	for _, n := range []string{"grep", "egrep", "fgrep", "zgrep"} {
		end := off + len(n)
		if end > len(source) {
			continue
		}
		if string(source[off:end]) != n {
			continue
		}
		// Boundary: next byte must be whitespace, newline, or end of file.
		if end < len(source) {
			c := source[end]
			if c != ' ' && c != '\t' && c != '\n' {
				continue
			}
		}
		name = n
		break
	}
	if name == "" {
		return nil
	}
	insertAt := off + len(name)
	line, col := offsetLineColZC1053(source, insertAt)
	if line < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    line,
		Column:  col,
		Length:  0,
		Replace: " -q",
	}}
}

func offsetLineColZC1053(source []byte, offset int) (int, int) {
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
			walkZC1053(stmt, isSilenced, violations)
		}
	case *ast.ExpressionStatement:
		walkZC1053(n.Expression, isSilenced, violations)
	case *ast.InfixExpression:
		zc1053WalkInfix(n, isSilenced, violations)
	case *ast.PrefixExpression:
		if n.Operator == "!" {
			walkZC1053(n.Right, isSilenced, violations)
		}
	case *ast.Redirection:
		walkZC1053(n.Left, isSilenced || zc1053SilencesStdout(n), violations)
	case *ast.SimpleCommand:
		checkCommandZC1053(n, isSilenced, violations)
	case *ast.GroupedExpression:
		walkZC1053(n.Expression, isSilenced, violations)
	}
}

func zc1053WalkInfix(n *ast.InfixExpression, isSilenced bool, violations *[]Violation) {
	if n.Operator == "|" {
		// Left side of pipe is silenced (stdout goes to pipe).
		walkZC1053(n.Left, true, violations)
		walkZC1053(n.Right, isSilenced, violations)
		return
	}
	walkZC1053(n.Left, isSilenced, violations)
	walkZC1053(n.Right, isSilenced, violations)
}

func zc1053SilencesStdout(n *ast.Redirection) bool {
	switch n.Operator {
	case ">", ">>", "&>":
		return isDevNull(n.Right)
	}
	return false
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
				argStr := arg.String()
				argStr = strings.Trim(argStr, "\"'")
				if strings.HasPrefix(argStr, "-") {
					if argStr == "-q" || argStr == "--quiet" || argStr == "--silent" {
						hasQuiet = true
						break
					}
					// Check for combined flags e.g. -rq
					if !strings.HasPrefix(argStr, "--") && strings.Contains(argStr, "q") {
						hasQuiet = true
						break
					}
				}
			}

			if !hasQuiet {
				*violations = append(*violations, Violation{
					KataID:  "ZC1053",
					Message: "Silence `grep` output in conditions. Use `grep -q` or redirect to `/dev/null`.",
					Line:    name.Token.Line,
					Column:  name.Token.Column,
					Level:   SeverityStyle,
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

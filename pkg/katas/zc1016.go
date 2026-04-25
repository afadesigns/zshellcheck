package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:    "ZC1016",
		Title: "Use `read -s` when reading sensitive information",
		Description: "When asking for passwords or secrets, use `read -s` to prevent " +
			"the input from being echoed to the terminal.",
		Severity: SeverityStyle,
		Check:    checkZC1016,
		Fix:      fixZC1016,
	})
}

// fixZC1016 inserts ` -s` after the `read` command name. The detector
// gates on the absence of `-s` in any flag bundle, so the insertion
// is idempotent on a re-run.
func fixZC1016(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if cmd.Name == nil || cmd.Name.String() != "read" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 || IdentLenAt(source, nameOff) != len("read") {
		return nil
	}
	insertAt := nameOff + len("read")
	insLine, insCol := offsetLineColZC1016(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -s",
	}}
}

func offsetLineColZC1016(source []byte, offset int) (int, int) {
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

func checkZC1016(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	if cmd.Name.String() != "read" {
		return nil
	}

	hasS := false
	sensitiveVars := []string{"password", "passwd", "pwd", "secret", "token", "key", "api_key"}

	// Check flags
	for _, arg := range cmd.Arguments {
		argStr := arg.String()
		// Remove quotes if present
		argStr = strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStr, "-") {
			if strings.Contains(argStr, "s") {
				hasS = true
			}
		}
	}

	if hasS {
		return nil
	}

	violations := []Violation{}

	for _, arg := range cmd.Arguments {
		// Skip flags
		argStr := arg.String()
		argStrClean := strings.Trim(argStr, "\"'")
		if strings.HasPrefix(argStrClean, "-") {
			continue
		}

		// Handle Zsh read syntax: variable?prompt
		parts := strings.Split(argStr, "?")
		varName := strings.TrimSpace(parts[0])
		varName = strings.Trim(varName, "'\"")

		varLower := strings.ToLower(varName)
		isSensitive := false
		for _, s := range sensitiveVars {
			if strings.Contains(varLower, s) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			violations = append(violations, Violation{
				KataID:  "ZC1016",
				Message: "Use `read -s` to hide input when reading sensitive variable '" + varName + "'.",
				Line:    cmd.TokenLiteralNode().Line,
				Column:  cmd.TokenLiteralNode().Column,
				Level:   SeverityStyle,
			})
		}
	}

	return violations
}

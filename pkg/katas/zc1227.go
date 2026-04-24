package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1227",
		Title:    "Use `curl -f` to fail on HTTP errors",
		Severity: SeverityWarning,
		Description: "`curl` without `-f` silently returns error pages (404, 500) as success. " +
			"Use `-f` or `--fail` to return exit code 22 on HTTP errors.",
		Check: checkZC1227,
		Fix:   fixZC1227,
	})
}

// fixZC1227 inserts ` -f` after the `curl` command name so HTTP
// errors translate into a non-zero exit code. Detector guards the
// shape (URL arg present, no existing -f/-fsSL/etc.).
func fixZC1227(node ast.Node, v Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}
	nameOff := LineColToByteOffset(source, v.Line, v.Column)
	if nameOff < 0 {
		return nil
	}
	nameLen := IdentLenAt(source, nameOff)
	if nameLen != len("curl") {
		return nil
	}
	insertAt := nameOff + nameLen
	insLine, insCol := offsetLineColZC1227(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " -f",
	}}
}

func offsetLineColZC1227(source []byte, offset int) (int, int) {
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

func checkZC1227(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "curl" {
		return nil
	}

	hasFail := false
	hasURL := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-f" || val == "-fsSL" || val == "-fSL" || val == "-fsS" {
			hasFail = true
		}
		if len(val) > 4 && (val[:4] == "http" || val[:5] == "https") {
			hasURL = true
		}
	}

	if hasURL && !hasFail {
		return []Violation{{
			KataID: "ZC1227",
			Message: "Use `curl -f` to fail on HTTP errors. Without `-f`, curl silently " +
				"returns error pages (404, 500) as if they were successful.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1234",
		Title:    "Use `docker run --rm` to auto-remove containers",
		Severity: SeverityStyle,
		Description: "`docker run` without `--rm` leaves stopped containers behind. " +
			"Use `--rm` in scripts to automatically clean up after execution.",
		Check: checkZC1234,
		Fix:   fixZC1234,
	})
}

// fixZC1234 inserts ` --rm` after the `run` subcommand in a
// `docker run …` invocation. Detector has already verified the shape
// (docker + run + no --rm + no -d).
func fixZC1234(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	runArg := cmd.Arguments[0]
	if runArg.String() != "run" {
		return nil
	}
	tok := runArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+3 > len(source) {
		return nil
	}
	if string(source[off:off+3]) != "run" {
		return nil
	}
	insertAt := off + 3
	insLine, insCol := offsetLineColZC1234(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --rm",
	}}
}

func offsetLineColZC1234(source []byte, offset int) (int, int) {
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

func checkZC1234(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "run" {
		return nil
	}

	hasRM := false
	hasDetach := false

	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--rm" {
			hasRM = true
		}
		if val == "-d" {
			hasDetach = true
		}
	}

	if !hasRM && !hasDetach {
		return []Violation{{
			KataID: "ZC1234",
			Message: "Use `docker run --rm` to auto-remove containers after exit. " +
				"Without `--rm`, stopped containers accumulate on disk.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

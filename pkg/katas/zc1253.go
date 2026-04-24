package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1253",
		Title:    "Use `docker build --no-cache` in CI for reproducible builds",
		Severity: SeverityStyle,
		Description: "`docker build` uses layer caching which can mask dependency changes. " +
			"Use `--no-cache` in CI pipelines to ensure fully reproducible builds.",
		Check: checkZC1253,
		Fix:   fixZC1253,
	})
}

// fixZC1253 inserts ` --no-cache` after the `build` subcommand in
// `docker build …`. Mirrors ZC1234's subcommand-level insertion.
func fixZC1253(node ast.Node, _ Violation, source []byte) []FixEdit {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	if len(cmd.Arguments) == 0 {
		return nil
	}
	subArg := cmd.Arguments[0]
	if subArg.String() != "build" {
		return nil
	}
	tok := subArg.TokenLiteralNode()
	off := LineColToByteOffset(source, tok.Line, tok.Column)
	if off < 0 || off+5 > len(source) {
		return nil
	}
	if string(source[off:off+5]) != "build" {
		return nil
	}
	insertAt := off + 5
	insLine, insCol := offsetLineColZC1253(source, insertAt)
	if insLine < 0 {
		return nil
	}
	return []FixEdit{{
		Line:    insLine,
		Column:  insCol,
		Length:  0,
		Replace: " --no-cache",
	}}
}

func offsetLineColZC1253(source []byte, offset int) (int, int) {
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

func checkZC1253(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "docker" {
		return nil
	}

	if len(cmd.Arguments) < 1 || cmd.Arguments[0].String() != "build" {
		return nil
	}

	hasNoCache := false
	for _, arg := range cmd.Arguments[1:] {
		val := arg.String()
		if val == "--no-cache" {
			hasNoCache = true
		}
	}

	if !hasNoCache {
		return []Violation{{
			KataID: "ZC1253",
			Message: "Consider `docker build --no-cache` in CI for reproducible builds. " +
				"Layer caching can mask changed dependencies.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}

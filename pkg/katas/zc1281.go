package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1281",
		Title:    "Use `sort -u` instead of `sort | uniq` for deduplication",
		Severity: SeverityStyle,
		Description: "`sort -u` combines sorting and deduplication in a single pass, " +
			"which is more efficient than piping `sort` into `uniq` as a separate process.",
		Check: checkZC1281,
	})
}

func checkZC1281(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "uniq" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "-c" || val == "-d" || val == "-D" || val == "-u" {
			return nil
		}
	}

	return []Violation{{
		KataID:  "ZC1281",
		Message: "Use `sort -u` instead of `sort | uniq`. The `-u` flag deduplicates in a single pass.",
		Line:    cmd.Token.Line,
		Column:  cmd.Token.Column,
		Level:   SeverityStyle,
	}}
}

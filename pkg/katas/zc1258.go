package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1258",
		Title:    "Consider `rsync --delete` for directory sync",
		Severity: SeverityWarning,
		Description: "`rsync` without `--delete` keeps files on the destination that were " +
			"removed from the source. Use `--delete` for true directory mirroring.",
		Check: checkZC1258,
	})
}

func checkZC1258(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "rsync" {
		return nil
	}

	hasDelete := false
	hasTrailingSlash := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "--delete" {
			hasDelete = true
		}
		if strings.HasSuffix(val, "/") && !strings.HasPrefix(val, "-") {
			hasTrailingSlash = true
		}
	}

	if hasTrailingSlash && !hasDelete {
		return []Violation{{
			KataID: "ZC1258",
			Message: "Consider `rsync --delete` for directory sync. Without `--delete`, " +
				"files removed from source remain on the destination.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

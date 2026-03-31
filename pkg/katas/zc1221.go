package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1221",
		Title:    "Avoid `fdisk` in scripts — use `parted` or `sfdisk`",
		Severity: SeverityWarning,
		Description: "`fdisk` is interactive and not scriptable. " +
			"Use `parted -s` or `sfdisk` for non-interactive disk partitioning.",
		Check: checkZC1221,
	})
}

func checkZC1221(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "fdisk" {
		return nil
	}

	// fdisk -l (list) is non-interactive and fine
	for _, arg := range cmd.Arguments {
		if arg.String() == "-l" {
			return nil
		}
	}

	return []Violation{{
		KataID: "ZC1221",
		Message: "Avoid `fdisk` in scripts — it is interactive. " +
			"Use `parted -s` or `sfdisk` for scriptable disk partitioning.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

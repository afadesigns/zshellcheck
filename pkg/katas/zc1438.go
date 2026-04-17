package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1438",
		Title:    "`systemctl mask` permanently prevents service start — document the unmask path",
		Severity: SeverityWarning,
		Description: "`systemctl mask unit` symlinks the unit to `/dev/null`, preventing any " +
			"start (manual, dependency, or at boot). Even `systemctl start` fails with 'Unit is " +
			"masked.'. The reverse `systemctl unmask` is easy to forget. Document the unmask in " +
			"provisioning scripts or use `disable` (which still allows manual start).",
		Check: checkZC1438,
	})
}

func checkZC1438(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "systemctl" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "mask" {
			return []Violation{{
				KataID: "ZC1438",
				Message: "`systemctl mask` permanently blocks service start. If this is a " +
					"policy choice, document the `unmask` path. For a softer block, use `disable`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1616",
		Title:    "Warn on `fsfreeze -f MOUNTPOINT` — filesystem stays frozen until `-u` runs",
		Severity: SeverityWarning,
		Description: "`fsfreeze -f` blocks every write on the mountpoint until `fsfreeze -u` " +
			"thaws it. The intended use is a short window around a hypervisor or LVM snapshot. " +
			"If the script errors between the freeze and the unfreeze (or is killed), the " +
			"filesystem stays frozen — every subsequent write hangs forever until the admin " +
			"manually thaws it, and a reboot may be the only way out on the root fs. Pair " +
			"every freeze with `trap 'fsfreeze -u MOUNTPOINT' EXIT` and keep the window under " +
			"a few seconds.",
		Check: checkZC1616,
	})
}

func checkZC1616(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "fsfreeze" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "-f" {
			return []Violation{{
				KataID: "ZC1616",
				Message: "`fsfreeze -f` freezes the mountpoint — every write hangs until " +
					"`fsfreeze -u` runs. Wrap the call in `trap 'fsfreeze -u PATH' EXIT` " +
					"so the thaw fires even on failure.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

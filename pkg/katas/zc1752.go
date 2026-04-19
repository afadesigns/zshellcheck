package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1752ForceFlags = map[string]bool{
	"-f": true, "-ff": true, "-fff": true,
	"--force": true,
	"-y":      true, "--yes": true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1752",
		Title:    "Error on `pvcreate/vgcreate/lvcreate -ff|--yes` — force-init LVM over existing data",
		Severity: SeverityError,
		Description: "LVM prompts before overwriting existing filesystem, RAID, or LVM signatures " +
			"on a device — that prompt is the only thing saving you from a typo'd target " +
			"destroying someone else's data. `pvcreate -ff`, `pvcreate --yes`, and the same " +
			"flags on `vgcreate` / `lvcreate` skip the prompt. Drop the flag, inspect with " +
			"`wipefs -n` and `lsblk -f` first, then confirm the target before re-running " +
			"the create command.",
		Check: checkZC1752,
	})
}

func checkZC1752(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	switch ident.Value {
	case "pvcreate", "vgcreate", "lvcreate":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if zc1752ForceFlags[v] {
			return []Violation{{
				KataID: "ZC1752",
				Message: "`" + ident.Value + " " + v + "` skips the LVM confirmation — a " +
					"wrong device gets its filesystem / RAID / LVM signatures wiped. " +
					"Inspect with `wipefs -n` + `lsblk -f` first, drop the flag, re-" +
					"run after checking the target.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

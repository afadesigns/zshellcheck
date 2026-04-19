package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1813",
		Title:    "Warn on `cryptsetup luksFormat` / `reencrypt` — destructive LUKS header write",
		Severity: SeverityWarning,
		Description: "`cryptsetup luksFormat DEV` writes a new LUKS2 header at the start of DEV " +
			"and marks the remaining space as fresh ciphertext — any pre-existing filesystem " +
			"or LUKS metadata is gone. `cryptsetup reencrypt DEV` rewrites the entire device " +
			"in place, and an interruption mid-write leaves the volume partially re-encrypted " +
			"and dependent on the `--resume-only` recovery path. Pair `luksFormat` with " +
			"`--batch-mode` only after verifying DEV via `lsblk -o NAME,MODEL,SERIAL`, always " +
			"back up the header (`cryptsetup luksHeaderBackup`) before touching it, and run " +
			"`reencrypt` on an unmounted volume with UPS-backed power.",
		Check: checkZC1813,
	})
}

func checkZC1813(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cryptsetup" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "luksFormat", "reencrypt", "luks-format":
			return []Violation{{
				KataID: "ZC1813",
				Message: "`cryptsetup " + v + "` rewrites the LUKS header / device. " +
					"Verify the target (`lsblk`), back up with " +
					"`luksHeaderBackup`, and run on an unmounted volume with UPS.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1540",
		Title:    "Error on `cryptsetup erase` / `luksErase` — destroys LUKS header, data unrecoverable",
		Severity: SeverityError,
		Description: "`cryptsetup erase` (alias `luksErase`) overwrites the LUKS header and " +
			"every key slot. Without the header the ciphertext on the device is unrecoverable " +
			"— even the original passphrase cannot unlock it. Keep a `cryptsetup " +
			"luksHeaderBackup` image somewhere safe before running erase, and prefer " +
			"`luksRemoveKey`/`luksKillSlot` when only rotating one slot.",
		Check: checkZC1540,
	})
}

func checkZC1540(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "cryptsetup" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "erase" || v == "luksErase" {
			return []Violation{{
				KataID: "ZC1540",
				Message: "`cryptsetup " + v + "` wipes the LUKS header — ciphertext becomes " +
					"unrecoverable. Back up the header first, or use luksRemoveKey/" +
					"luksKillSlot for single-slot rotation.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

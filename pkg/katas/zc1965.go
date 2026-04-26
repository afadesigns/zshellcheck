// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1965",
		Title:    "Error on `systemd-cryptenroll --wipe-slot=all` — wipes every LUKS key slot",
		Severity: SeverityError,
		Description: "`systemd-cryptenroll --wipe-slot=all $DEV` removes every key slot on the " +
			"LUKS volume — passphrase, recovery key, TPM2, FIDO2, PKCS#11 — in one call. " +
			"`--wipe-slot=recovery` / `--wipe-slot=empty` are scoped; the `all` form is a " +
			"one-shot brick with no confirmation. Either enrol the new slot first and then " +
			"wipe the specific index you are retiring (`--wipe-slot=<n>`), or back up the " +
			"header with `cryptsetup luksHeaderBackup` before the call so recovery is " +
			"possible.",
		Check: checkZC1965,
	})
}

func checkZC1965(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `systemd-cryptenroll --wipe-slot=all $DEV` mangles the
	// command name to `wipe-slot=all`.
	if strings.HasPrefix(ident.Value, "wipe-slot=") {
		if ident.Value == "wipe-slot=all" {
			return zc1965Hit(cmd)
		}
	}
	if ident.Value != "systemd-cryptenroll" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--wipe-slot=all" || v == "--wipe-slot" {
			if v == "--wipe-slot=all" {
				return zc1965Hit(cmd)
			}
		}
	}
	return nil
}

func zc1965Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1965",
		Message: "`systemd-cryptenroll --wipe-slot=all` wipes every LUKS key slot " +
			"(passphrase/recovery/TPM2/FIDO2) in one call. Enrol the new slot first, " +
			"wipe a specific index, back up the header with `cryptsetup luksHeaderBackup`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

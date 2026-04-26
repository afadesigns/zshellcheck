// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1998",
		Title:    "Error on `tpm2_clear` / `tpm2 clear` — wipes TPM storage hierarchy, kills every sealed key",
		Severity: SeverityError,
		Description: "`tpm2_clear -c p` (or `tpm2 clear -c p`) invokes the TPM 2.0 `TPM2_Clear` " +
			"command, which invalidates every object sealed against the storage " +
			"hierarchy — LUKS-TPM2 keyslots, systemd-cryptenroll's `--tpm2-device` " +
			"slot, sshd TPM-backed host keys, and SecureBoot measured-boot state. " +
			"The machine can still boot but any disk that unlocked through the TPM " +
			"now needs a recovery passphrase, and every TLS cert issued from a " +
			"TPM-sealed CA loses its anchor. There is no undo. Run `tpm2_clear` only " +
			"under a documented recovery runbook with the recovery material in hand; " +
			"never put it in an automated scheduled script.",
		Check: checkZC1998,
	})
}

func checkZC1998(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value == "tpm2_clear" {
		return zc1998Hit(cmd, "tpm2_clear")
	}
	if ident.Value == "tpm2" && len(cmd.Arguments) > 0 &&
		cmd.Arguments[0].String() == "clear" {
		return zc1998Hit(cmd, "tpm2 clear")
	}
	return nil
}

func zc1998Hit(cmd *ast.SimpleCommand, form string) []Violation {
	return []Violation{{
		KataID: "ZC1998",
		Message: "`" + form + "` wipes the TPM storage hierarchy — every LUKS-TPM2 " +
			"keyslot, `systemd-cryptenroll --tpm2-device` slot, and TPM-sealed " +
			"TLS/sshd key is destroyed. No undo. Gate behind a recovery runbook.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

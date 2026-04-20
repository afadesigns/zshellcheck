package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1978",
		Title:    "Warn on `ftp` / `tftp` — cleartext transfer, credentials and payload exposed on the wire",
		Severity: SeverityWarning,
		Description: "`ftp HOST` negotiates USER/PASS and moves the payload over plaintext TCP " +
			"(port 21) — every credential and byte is visible to anything between the " +
			"caller and the server. `tftp` has no auth at all and runs over UDP/69, so " +
			"any packet capture recovers the full transfer. Both are also routinely " +
			"mishandled by NAT/firewall gear because of their dual-channel design. " +
			"Replace with `curl -u USER: https://…` / `sftp` / `scp` / `rsync -e ssh` " +
			"for authenticated transfers, and with a signed-payload pull over HTTPS for " +
			"PXE-style provisioning that used to rely on `tftp`.",
		Check: checkZC1978,
	})
}

func checkZC1978(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "ftp" && ident.Value != "tftp" {
		return nil
	}
	// Require at least one arg so bare `ftp` at a prompt isn't flagged.
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1978",
		Message: "`" + ident.Value + "` transfers in plaintext — creds and payload " +
			"visible on the wire. Use `sftp`/`scp`/`rsync -e ssh` or a signed-" +
			"payload `curl` over HTTPS instead.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

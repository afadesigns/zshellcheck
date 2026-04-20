package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1978",
		Title:    "Warn on `tftp` — cleartext, unauthenticated UDP transfer",
		Severity: SeverityWarning,
		Description: "`tftp` has no authentication at all and moves the payload in plaintext " +
			"over UDP/69 — any packet capture on the path recovers the full transfer " +
			"and an attacker at the server can push an arbitrary file under the " +
			"expected name without noticing a lack of credentials. The dual-channel " +
			"design is also routinely mishandled by NAT/firewall gear. For PXE-style " +
			"provisioning that historically used `tftp`, fetch a signed payload over " +
			"HTTPS with `curl` and verify the signature locally before use. (See " +
			"ZC1200 for `ftp`, the authenticated-but-plaintext sibling.)",
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
	// `ftp` is owned by ZC1200; ZC1978 narrows to tftp (no auth, UDP).
	if ident.Value != "tftp" {
		return nil
	}
	// Require at least one arg so bare `tftp` at a prompt isn't flagged.
	if len(cmd.Arguments) == 0 {
		return nil
	}
	return []Violation{{
		KataID: "ZC1978",
		Message: "`tftp` transfers over plaintext UDP/69 with no authentication — " +
			"capture the payload, or push a crafted file under the expected " +
			"name. Use a signed-payload `curl` over HTTPS and verify the " +
			"signature before use.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

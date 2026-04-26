// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1862",
		Title:    "Warn on `ssh-keygen -R HOST` — deletes a known-hosts entry, next `ssh` re-trusts silently",
		Severity: SeverityWarning,
		Description: "`ssh-keygen -R HOST` scrubs the entry for `HOST` from `~/.ssh/known_hosts`. " +
			"The legitimate trigger is a real key rotation (server reinstall, HSM " +
			"replacement), but the flag is frequently dropped into automation to " +
			"silence the REMOTE HOST IDENTIFICATION HAS CHANGED banner without ever " +
			"confirming the new fingerprint. The very next `ssh` call then prompts " +
			"once (or not at all under `StrictHostKeyChecking=no`) and blindly accepts " +
			"whatever the network hands back — a MITM attacker who was waiting for a " +
			"rebuild slips in without a trace. Fetch the new key out-of-band and " +
			"`ssh-keyscan -t rsa,ed25519 HOST | ssh-keygen -lf -` before adding it, or " +
			"pin fingerprints in a managed `known_hosts` file.",
		Check: checkZC1862,
	})
}

func checkZC1862(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "ssh-keygen" {
		return nil
	}
	args := cmd.Arguments
	for i, arg := range args {
		v := arg.String()
		if v != "-R" {
			continue
		}
		if i+1 >= len(args) {
			return nil
		}
		host := args[i+1].String()
		return []Violation{{
			KataID: "ZC1862",
			Message: "`ssh-keygen -R " + host + "` deletes a known-hosts entry — the " +
				"next `ssh` silently re-trusts whatever key the network returns. " +
				"Fetch the new fingerprint out-of-band and verify before " +
				"re-adding.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}
	return nil
}

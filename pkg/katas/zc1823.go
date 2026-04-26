// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1823",
		Title:    "Warn on `keytool -import -noprompt` — Java trust store imports without fingerprint check",
		Severity: SeverityWarning,
		Description: "`keytool -import -noprompt -trustcacerts -alias X -file CERT -keystore KS` " +
			"adds CERT to the Java trust store without showing its SHA-256 fingerprint or " +
			"asking the operator to confirm. If CERT came from an HTTP download, an attacker " +
			"wrote it in a shared temp dir, or a provisioning step fetched the wrong file, the " +
			"JVM will happily pin the attacker's CA as trusted and verify everything signed " +
			"against it. Drop `-noprompt`, or pre-verify with `keytool -printcert -file CERT` " +
			"and keep the alias+fingerprint pair in a versioned inventory before adding to any " +
			"trust store.",
		Check: checkZC1823,
	})
}

func checkZC1823(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "keytool" {
		return nil
	}

	hasImport := false
	hasNoPrompt := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-import" || v == "-importcert" || v == "-importkeystore" {
			hasImport = true
		}
		if v == "-noprompt" {
			hasNoPrompt = true
		}
	}
	if !hasImport || !hasNoPrompt {
		return nil
	}
	return []Violation{{
		KataID: "ZC1823",
		Message: "`keytool -import -noprompt` pins a cert to the Java trust store " +
			"without a fingerprint check. Drop `-noprompt`, verify with " +
			"`keytool -printcert -file CERT`, and store (alias, SHA-256) pairs in " +
			"an audited inventory.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

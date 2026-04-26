// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1794CosignSkipFlags = map[string]bool{
	"--insecure-ignore-tlog":    true,
	"--insecure-ignore-sct":     true,
	"--insecure-skip-verify":    true,
	"--allow-insecure-registry": true,
	"--allow-http-registry":     true,
	"--allow-insecure-bundle":   true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1794",
		Title:    "Error on `cosign verify --insecure-ignore-tlog` / `--allow-insecure-registry` — signature chain disabled",
		Severity: SeverityError,
		Description: "`cosign verify` with `--insecure-ignore-tlog` skips Rekor transparency-log " +
			"verification, `--insecure-ignore-sct` skips Fulcio SCT verification, and " +
			"`--insecure-skip-verify` turns off TLS certificate validation for the registry / " +
			"Rekor / Fulcio endpoints. `cosign sign --allow-insecure-registry` and " +
			"`--allow-http-registry` push signatures over plain HTTP. Each flag removes a " +
			"distinct rung of the signature chain that `cosign` was built to enforce — a " +
			"malicious registry or on-path attacker now passes verification without detection. " +
			"Drop the flag, fix the underlying trust anchor (CA bundle, Rekor URL, " +
			"Fulcio OIDC), and keep signature verification strict.",
		Check: checkZC1794,
	})
}

func checkZC1794(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "cosign" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		flag := v
		if idx := strings.Index(flag, "="); idx >= 0 {
			flag = flag[:idx]
		}
		if zc1794CosignSkipFlags[flag] {
			return []Violation{{
				KataID: "ZC1794",
				Message: "`cosign " + v + "` removes a rung of the signature chain " +
					"(transparency log / SCT / TLS / HTTPS-only registry). Drop " +
					"the flag and fix the underlying trust anchor.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

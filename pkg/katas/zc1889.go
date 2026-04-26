// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1889",
		Title:    "Error on `skopeo copy --src-tls-verify=false` / `--dest-tls-verify=false` — MITM on image copy",
		Severity: SeverityError,
		Description: "`skopeo copy` is the glue for promoting container images between registries in " +
			"CI, mirroring upstream images into internal caches, and rehydrating images " +
			"to an air-gapped registry. `--src-tls-verify=false` and " +
			"`--dest-tls-verify=false` drop certificate verification on the respective " +
			"leg, which means any on-path attacker can substitute a malicious manifest or " +
			"layer and the copy completes without a warning. Use `--src-cert-dir`/" +
			"`--dest-cert-dir` to pin a private CA if you are mirroring to or from an " +
			"internal registry with self-signed certs, or fix the upstream's cert.",
		Check: checkZC1889,
	})
}

func checkZC1889(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "skopeo" {
		return nil
	}
	for _, arg := range cmd.Arguments {
		v := arg.String()
		lower := strings.ToLower(v)
		for _, prefix := range []string{"--src-tls-verify=", "--dest-tls-verify=", "--tls-verify="} {
			if !strings.HasPrefix(lower, prefix) {
				continue
			}
			val := strings.TrimPrefix(lower, prefix)
			val = strings.Trim(val, "\"'")
			if val == "false" || val == "0" || val == "no" || val == "off" {
				return []Violation{{
					KataID: "ZC1889",
					Message: "`skopeo " + v + "` disables TLS verification on image " +
						"copy — on-path attacker can substitute a malicious manifest. " +
						"Pin a private CA with `--src-cert-dir`/`--dest-cert-dir` " +
						"instead.",
					Line:   cmd.Token.Line,
					Column: cmd.Token.Column,
					Level:  SeverityError,
				}}
			}
		}
	}
	return nil
}

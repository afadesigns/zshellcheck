// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1822",
		Title:    "Error on `csrutil disable` / `spctl --master-disable` — disables macOS system integrity / Gatekeeper",
		Severity: SeverityError,
		Description: "`csrutil disable` turns off System Integrity Protection: the kernel stops " +
			"blocking writes under `/System`, `/bin`, `/sbin`, runtime attachment to " +
			"protected processes becomes possible, and unsigned kexts can load. `spctl " +
			"--master-disable` (and `--global-disable`, `kext-consent disable`) removes " +
			"Gatekeeper / kext-consent enforcement, so any downloaded binary or kernel " +
			"extension runs without the user being prompted. Neither has a legitimate " +
			"provisioning use; both belong to ad-hoc developer workflows and are high-value " +
			"persistence steps for malware. Re-enable with `csrutil enable` in recovery mode " +
			"and `spctl --master-enable`.",
		Check: checkZC1822,
	})
}

func checkZC1822(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	// Parser caveat: `spctl --master-disable` mangles to name=`master-disable`,
	// `spctl --global-disable` to `global-disable`.
	switch ident.Value {
	case "master-disable", "global-disable":
		return zc1822Hit(cmd, "spctl --"+ident.Value)
	}

	switch ident.Value {
	case "csrutil":
		if len(cmd.Arguments) > 0 && cmd.Arguments[0].String() == "disable" {
			return zc1822Hit(cmd, "csrutil disable")
		}
	case "spctl":
		for _, arg := range cmd.Arguments {
			v := arg.String()
			if v == "--master-disable" || v == "--global-disable" {
				return zc1822Hit(cmd, "spctl "+v)
			}
		}
		if len(cmd.Arguments) >= 2 &&
			cmd.Arguments[0].String() == "kext-consent" &&
			cmd.Arguments[1].String() == "disable" {
			return zc1822Hit(cmd, "spctl kext-consent disable")
		}
	}
	return nil
}

func zc1822Hit(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1822",
		Message: "`" + what + "` disables macOS SIP / Gatekeeper / kext-consent — " +
			"every malware analyst's favorite persistence primitive. Re-enable " +
			"(`csrutil enable` in recovery, `spctl --master-enable`) and keep " +
			"the default policy on.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

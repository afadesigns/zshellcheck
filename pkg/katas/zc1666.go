// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1666",
		Title:    "Warn on `kubectl patch --type=json` — bypasses strategic-merge defaults",
		Severity: SeverityWarning,
		Description: "`kubectl patch --type=json` applies a raw RFC-6902 JSON patch: `remove`, " +
			"`replace`, `add /spec/containers/0`, and `move` land verbatim on the resource. " +
			"Unlike strategic-merge or merge-patch, Kubernetes does not reconcile the " +
			"patch against field ownership or default values — so a mistyped `path` or an " +
			"index that no longer exists fails silently or drops the wrong field. From a " +
			"script this is a foot-gun for drift and supply-chain compromise: an attacker " +
			"with write access to the patch file can slip `privileged: true` or `hostPath` " +
			"mounts in. Prefer `--type=strategic` (the default) and hold JSON patches " +
			"behind code review.",
		Check: checkZC1666,
	})
}

func checkZC1666(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" {
		return nil
	}

	if len(cmd.Arguments) == 0 || cmd.Arguments[0].String() != "patch" {
		return nil
	}

	for i, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--type=json" {
			return zc1666Hit(cmd)
		}
		if v == "--type" && i+1 < len(cmd.Arguments) &&
			cmd.Arguments[i+1].String() == "json" {
			return zc1666Hit(cmd)
		}
	}
	return nil
}

func zc1666Hit(cmd *ast.SimpleCommand) []Violation {
	return []Violation{{
		KataID: "ZC1666",
		Message: "`kubectl patch --type=json` applies a raw RFC-6902 patch that bypasses " +
			"strategic-merge reconciliation — prefer `--type=strategic` and hold JSON " +
			"patches behind code review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

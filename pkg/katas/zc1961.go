// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1961",
		Title:    "Warn on `gcloud iam service-accounts keys create` — mints a long-lived service-account JSON key",
		Severity: SeverityWarning,
		Description: "`gcloud iam service-accounts keys create key.json --iam-account=SA@PROJECT` " +
			"exports an RSA key pair wrapped in a JSON file. Once written it is effectively a " +
			"forever-valid bearer credential: no automatic rotation, no refresh, and a single " +
			"\"leaked by a `cat key.json`\" is game-over. Prefer Workload Identity Federation " +
			"(`gcloud iam workload-identity-pools …`), short-lived impersonation via " +
			"`gcloud auth print-access-token --impersonate-service-account=SA`, or the " +
			"key-less GCE/GKE attached service account. Reserve static JSON keys for provably " +
			"off-platform callers.",
		Check: checkZC1961,
	})
}

func checkZC1961(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "gcloud" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "iam" ||
		cmd.Arguments[1].String() != "service-accounts" ||
		cmd.Arguments[2].String() != "keys" ||
		cmd.Arguments[3].String() != "create" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1961",
		Message: "`gcloud iam service-accounts keys create` mints a long-lived JSON key — " +
			"no auto-rotate, no refresh. Prefer Workload Identity Federation, " +
			"`--impersonate-service-account`, or the attached service account.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

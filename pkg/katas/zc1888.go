// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1888",
		Title:    "Warn on `aws iam create-access-key` — mints long-lived static AWS credentials",
		Severity: SeverityWarning,
		Description: "`aws iam create-access-key` hands out a static `AKIA.../secret` pair that is " +
			"valid forever until someone rotates it; whoever gets the pair speaks for " +
			"the IAM user on every API call AWS accepts. Most modern deploys no longer " +
			"need these: EC2 instance profiles, EKS/IRSA, Lambda roles, GitHub OIDC, " +
			"and IAM Identity Center all hand out short-lived session credentials on " +
			"demand. Prefer those; if a static key is genuinely required (legacy third-" +
			"party tooling), store it in AWS Secrets Manager, scope the user to the " +
			"narrowest policy possible, and rotate on a schedule with `aws iam update-" +
			"access-key --status Inactive` / `delete-access-key`.",
		Check: checkZC1888,
	})
}

func checkZC1888(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "aws" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "iam" || args[1].String() != "create-access-key" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1888",
		Message: "`aws iam create-access-key` mints a long-lived `AKIA.../secret` — " +
			"prefer short-lived creds via instance profiles, IRSA, Lambda roles, " +
			"or OIDC federation. If static keys are unavoidable, store in Secrets " +
			"Manager and rotate.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

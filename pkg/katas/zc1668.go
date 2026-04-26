// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1668",
		Title:    "Error on `aws iam attach-*-policy ... AdministratorAccess` — grants full AWS admin",
		Severity: SeverityError,
		Description: "Attaching the AWS-managed `AdministratorAccess` (or `PowerUserAccess`) " +
			"policy gives the target principal `*:*` — create/delete IAM users, mutate KMS " +
			"keys, rotate root passwords, exfiltrate every S3 bucket. Scripts rarely need " +
			"full admin; the pattern usually means someone hit a permissions error and " +
			"replaced the scoped policy with the blanket one. Write a least-privilege inline " +
			"policy (`iam put-user-policy --policy-document`), or reference a customer-" +
			"managed policy with only the `Action`/`Resource` pairs the workload needs. Admin " +
			"attachment should land via change-reviewed Terraform, not a shell loop.",
		Check: checkZC1668,
	})
}

func checkZC1668(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "aws" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "iam" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "attach-user-policy" && sub != "attach-role-policy" &&
		sub != "attach-group-policy" {
		return nil
	}

	for _, arg := range cmd.Arguments[2:] {
		v := arg.String()
		if strings.HasSuffix(v, "/AdministratorAccess") ||
			strings.HasSuffix(v, "/PowerUserAccess") ||
			strings.HasSuffix(v, "/IAMFullAccess") {
			return []Violation{{
				KataID: "ZC1668",
				Message: "`aws iam " + sub + " ... " + v + "` grants sweeping admin — " +
					"use a scoped inline policy (`put-user-policy`) or a customer-managed " +
					"policy with the minimum `Action`/`Resource` set.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

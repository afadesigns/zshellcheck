// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1480",
		Title:    "Warn on `terraform apply -auto-approve` / `destroy -auto-approve` in scripts",
		Severity: SeverityWarning,
		Description: "Running `terraform apply -auto-approve` or `destroy -auto-approve` from a " +
			"shell script skips the plan-review step that exists to catch schema drift, " +
			"accidental `-replace`, and resources being deleted. Fine for throwaway CI against " +
			"a PR environment, but dangerous against shared state. Prefer running `plan` + " +
			"`apply` with an out-file and human approval, or scope the auto-apply to specific " +
			"branches/environments.",
		Check: checkZC1480,
	})
}

func checkZC1480(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "terraform" && ident.Value != "terragrunt" && ident.Value != "tofu" {
		return nil
	}

	if len(cmd.Arguments) == 0 {
		return nil
	}
	sub := cmd.Arguments[0].String()
	if sub != "apply" && sub != "destroy" {
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "-auto-approve" || v == "--auto-approve" ||
			v == "-auto-approve=true" || v == "--auto-approve=true" {
			return []Violation{{
				KataID: "ZC1480",
				Message: "`" + ident.Value + " " + sub + " " + v + "` skips plan review. " +
					"Gate behind a branch/env check or require manual approval.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

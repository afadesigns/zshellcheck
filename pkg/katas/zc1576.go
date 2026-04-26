// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1576",
		Title:    "Warn on `terraform apply -target=...` — cherry-pick apply bypasses dependencies",
		Severity: SeverityWarning,
		Description: "`-target` restricts `terraform apply` to a specific resource / module and " +
			"everything it depends on. In theory that is a surgical fix; in practice it " +
			"routinely skips changes the targeted resource actually depends on, leading to " +
			"drift between state and configuration. HashiCorp documents `-target` as a tool " +
			"for incident response, not routine operations. Re-run without `-target` or " +
			"split the configuration into separate root modules.",
		Check: checkZC1576,
	})
}

func checkZC1576(node ast.Node) []Violation {
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
	if sub != "apply" && sub != "destroy" && sub != "plan" {
		return nil
	}

	var prevTarget bool
	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if prevTarget {
			return zc1576Violation(cmd, "-target "+v)
		}
		if v == "-target" {
			prevTarget = true
			continue
		}
		if strings.HasPrefix(v, "-target=") {
			return zc1576Violation(cmd, v)
		}
	}
	return nil
}

func zc1576Violation(cmd *ast.SimpleCommand, what string) []Violation {
	return []Violation{{
		KataID: "ZC1576",
		Message: "`terraform " + what + "` bypasses dependency order — documented as incident " +
			"response tool only. Re-run without -target or split root modules.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

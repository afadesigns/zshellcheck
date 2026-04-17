package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1443",
		Title:    "Dangerous: `terraform destroy` / `apply -destroy` without `-target`",
		Severity: SeverityWarning,
		Description: "`terraform destroy` (or `terraform apply -destroy`) without a `-target` " +
			"removes every resource in state — entire environments, databases, volumes, DNS, " +
			"everything. Always prefer targeted destroy or scope via workspaces. Consider " +
			"guarding state-destroying commands behind an interactive confirmation.",
		Check: checkZC1443,
	})
}

func checkZC1443(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || (ident.Value != "terraform" && ident.Value != "tofu") {
		return nil
	}

	hasDestroy := false
	hasTarget := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "destroy" || v == "-destroy" {
			hasDestroy = true
		}
		if strings.HasPrefix(v, "-target") || v == "-target" {
			hasTarget = true
		}
	}
	if hasDestroy && !hasTarget {
		return []Violation{{
			KataID: "ZC1443",
			Message: "`terraform destroy` without `-target` removes every resource in state. " +
				"Scope with `-target=...` or gate behind interactive confirmation.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}

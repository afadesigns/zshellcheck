package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1695",
		Title:    "Warn on `terraform state rm` / `state push` — surgery on shared state outside plan/apply",
		Severity: SeverityWarning,
		Description: "`terraform state rm RESOURCE` drops the resource from Terraform's " +
			"tracking without touching the real cloud object — the next `terraform apply` " +
			"sees it as newly-created and tries to re-provision, often hitting name-" +
			"collision errors. `terraform state push FILE` replaces the entire remote " +
			"state with a local file, bypassing locking and overwriting any concurrent " +
			"changes. Both commands skirt the usual plan/apply audit trail. Reach for " +
			"`terraform import` / `terraform apply -replace=ADDR` instead, and only run " +
			"`state rm|push` from a reviewed fix-up PR with state backup in hand.",
		Check: checkZC1695,
	})
}

func checkZC1695(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "terraform" && ident.Value != "tofu" && ident.Value != "terragrunt" {
		return nil
	}

	if len(cmd.Arguments) < 2 {
		return nil
	}
	if cmd.Arguments[0].String() != "state" {
		return nil
	}
	sub := cmd.Arguments[1].String()
	if sub != "rm" && sub != "push" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1695",
		Message: "`" + ident.Value + " state " + sub + "` mutates shared state outside " +
			"plan/apply — use `terraform import` or `apply -replace=ADDR` instead, and " +
			"review / back up first.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

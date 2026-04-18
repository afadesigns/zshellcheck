package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1726",
		Title:    "Error on `gcloud ... delete --quiet` — silent destruction of GCP resources",
		Severity: SeverityError,
		Description: "`gcloud` accepts `--quiet` (`-q`) globally to suppress every confirmation " +
			"prompt. Combined with `delete` on projects, SQL instances, GKE clusters, " +
			"compute VMs, secrets, or storage buckets, a single misresolved variable wipes " +
			"the resource with no human-in-the-loop. Project deletion has a 30-day soft " +
			"window but compute disks, secrets, and BigQuery tables are gone immediately. " +
			"Drop `--quiet` from delete commands or route the bulk-destroy through a " +
			"Terraform plan that surfaces the diff for review.",
		Check: checkZC1726,
	})
}

func checkZC1726(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}

	hasDelete, hasQuiet := false, false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "delete":
			hasDelete = true
		case "--quiet", "-q":
			hasQuiet = true
		}
	}
	if !hasDelete || !hasQuiet {
		return nil
	}

	return []Violation{{
		KataID: "ZC1726",
		Message: "`gcloud ... delete --quiet` skips confirmation — a wrong argument " +
			"wipes the resource (compute disks, secrets, BigQuery tables have no soft-" +
			"delete). Drop `--quiet` or destroy through a Terraform plan with review.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

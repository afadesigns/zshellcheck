package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1442",
		Title:    "Dangerous: `kubectl delete --all` / `--all-namespaces` deletes cluster resources",
		Severity: SeverityError,
		Description: "`kubectl delete --all pods` (in the current namespace) or " +
			"`-A`/`--all-namespaces` scopes delete operations across the whole cluster. A typo " +
			"on the resource type can wipe deployments, services, secrets, or even CRDs. " +
			"Always use `--dry-run=client` first, then apply with `-n` explicit namespace.",
		Check: checkZC1442,
	})
}

func checkZC1442(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "kubectl" && ident.Value != "oc" {
		return nil
	}

	hasDelete := false
	hasAll := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "delete" {
			hasDelete = true
		}
		if v == "--all" || v == "-A" || v == "--all-namespaces" {
			hasAll = true
		}
	}
	if hasDelete && hasAll {
		return []Violation{{
			KataID: "ZC1442",
			Message: "`kubectl delete --all` (or `-A`) deletes resources cluster-wide. Dry-run " +
				"with `--dry-run=client -o yaml` first, and scope with `-n` namespace.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}

	return nil
}

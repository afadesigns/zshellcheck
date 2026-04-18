package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1679BroadRoles = map[string]struct{}{
	"roles/owner":                             {},
	"roles/editor":                            {},
	"roles/iam.securityAdmin":                 {},
	"roles/iam.serviceAccountTokenCreator":    {},
	"roles/iam.serviceAccountKeyAdmin":        {},
	"roles/iam.workloadIdentityUser":          {},
	"roles/resourcemanager.organizationAdmin": {},
	"roles/resourcemanager.projectIamAdmin":   {},
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1679",
		Title:    "Error on `gcloud ... add-iam-policy-binding ... --role=roles/owner` — GCP primitive admin",
		Severity: SeverityError,
		Description: "`gcloud projects|folders|organizations add-iam-policy-binding` with the " +
			"primitive roles `roles/owner` or `roles/editor`, or with the IAM-escalation " +
			"roles (`roles/iam.securityAdmin`, `roles/iam.serviceAccountTokenCreator`, " +
			"`roles/iam.serviceAccountKeyAdmin`, `roles/resourcemanager.organizationAdmin`), " +
			"hands the principal the ability to grant themselves any other permission. " +
			"Scripts rarely need that scope; the pattern signals someone papering over a " +
			"permissions error. Grant a specific predefined role (e.g. `roles/compute." +
			"viewer`) or build a custom role with only the `Action`s the workload needs, " +
			"and apply admin changes via Terraform under change review.",
		Check: checkZC1679,
	})
}

func checkZC1679(node ast.Node) []Violation {
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

	hasAdd := false
	var hit string
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "add-iam-policy-binding" {
			hasAdd = true
			continue
		}
		if strings.HasPrefix(v, "--role=") {
			if _, broad := zc1679BroadRoles[strings.TrimPrefix(v, "--role=")]; broad {
				hit = v
			}
		}
	}
	if !hasAdd || hit == "" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1679",
		Message: "`gcloud ... add-iam-policy-binding " + hit + "` grants primitive / IAM-" +
			"admin — use a predefined role with the minimum scope or a custom role, and " +
			"apply admin changes via Terraform.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

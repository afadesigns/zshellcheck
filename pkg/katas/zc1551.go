package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1551",
		Title:    "Warn on `helm install/upgrade --skip-crds` — chart CRs land before their CRDs",
		Severity: SeverityWarning,
		Description: "`--skip-crds` tells Helm to install only the `.Release` objects and skip " +
			"the CustomResourceDefinition manifests under `crds/`. Without the CRDs present, " +
			"any `.Release` object that references a custom resource is rejected by the API " +
			"server at validation time, or — worse — fails later when a reconciler tries to " +
			"watch a type that does not exist. Use the default (install CRDs) on first roll- " +
			"out; if you need split lifecycle, install CRDs manually (`kubectl apply -f " +
			"chart/crds/`) before the `helm install`.",
		Check: checkZC1551,
	})
}

func checkZC1551(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "helm" {
		return nil
	}

	var sawVerb bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "upgrade" {
			sawVerb = true
			continue
		}
		if !sawVerb {
			continue
		}
		if v == "--skip-crds" {
			return []Violation{{
				KataID: "ZC1551",
				Message: "`helm --skip-crds` installs .Release objects without their CRDs — " +
					"custom resources fail validation. Install CRDs first.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}

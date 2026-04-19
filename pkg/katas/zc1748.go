package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1748",
		Title:    "Error on `helm repo add NAME http://...` — plaintext chart repo allows MITM",
		Severity: SeverityError,
		Description: "`helm repo add NAME http://URL` registers a chart repository reached over " +
			"plaintext HTTP. Any network-position attacker can swap `index.yaml` or a " +
			"chart tarball in flight, and subsequent `helm install` pulls container images " +
			"and Kubernetes manifests straight from the substituted content — fast path to " +
			"cluster-wide code execution. Use `https://`, and pair it with chart provenance " +
			"(`helm install --verify` or OCI signatures) to pin the digest.",
		Check: checkZC1748,
	})
}

func checkZC1748(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "helm" {
		return nil
	}
	if len(cmd.Arguments) < 4 {
		return nil
	}
	if cmd.Arguments[0].String() != "repo" || cmd.Arguments[1].String() != "add" {
		return nil
	}

	url := cmd.Arguments[3].String()
	if !strings.HasPrefix(url, "http://") {
		return nil
	}

	return []Violation{{
		KataID: "ZC1748",
		Message: "`helm repo add " + cmd.Arguments[2].String() + " " + url + "` fetches " +
			"charts over plaintext HTTP — any MITM swaps the chart and its referenced " +
			"images. Use `https://` and `helm install --verify` for provenance.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityError,
	}}
}

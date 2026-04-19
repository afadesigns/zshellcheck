package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1793",
		Title:    "Warn on `kubectl certificate approve CSR` — signs the identity baked into the CSR",
		Severity: SeverityWarning,
		Description: "`kubectl certificate approve NAME` tells the cluster signer to sign the " +
			"pending CSR unchanged. The signer respects the Subject (CN, O) and the " +
			"SubjectAltName extensions the caller put in the CSR — approve one that requests " +
			"`system:masters` and you have handed the requester full admin on the cluster. " +
			"In automation, review the CSR body first (`kubectl get csr NAME -o " +
			"jsonpath='{.spec.request}' | base64 -d | openssl req -text`) and reject (`kubectl " +
			"certificate deny`) any request that names a privileged group, kube-system service " +
			"account, or hostname outside the intended scope.",
		Check: checkZC1793,
	})
}

func checkZC1793(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	if len(cmd.Arguments) < 3 {
		return nil
	}
	if cmd.Arguments[0].String() != "certificate" {
		return nil
	}
	if cmd.Arguments[1].String() != "approve" {
		return nil
	}
	return []Violation{{
		KataID: "ZC1793",
		Message: "`kubectl certificate approve` signs the identity embedded in the CSR " +
			"— a `system:masters` request becomes cluster admin. Decode with " +
			"`openssl req -text` first; use `kubectl certificate deny` otherwise.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}

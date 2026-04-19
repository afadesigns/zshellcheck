package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1891",
		Title:    "Error on `kubectl config view --raw` — prints the full kubeconfig with client keys",
		Severity: SeverityError,
		Description: "`kubectl config view` by default redacts secrets: `client-certificate-data`, " +
			"`client-key-data`, `token`, and `password` fields are replaced with `REDACTED`. " +
			"Adding `--raw` (or the synonym `-R`) undoes every redaction and prints the " +
			"client's base64-encoded private key, bearer tokens, and any embedded user " +
			"password to stdout. In a script where stdout lands in CI log storage, a " +
			"`journalctl` ring buffer, or a Slack paste, the entire kubeconfig walks out. " +
			"Emit only the specific field you need (e.g. `kubectl config view -o " +
			"jsonpath='{.current-context}'`) or decrypt once into a temp file and `shred` it.",
		Check: checkZC1891,
	})
}

func checkZC1891(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "kubectl" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 2 {
		return nil
	}
	if args[0].String() != "config" || args[1].String() != "view" {
		return nil
	}
	for _, arg := range args[2:] {
		v := arg.String()
		if v == "--raw" || v == "-R" {
			return []Violation{{
				KataID: "ZC1891",
				Message: "`kubectl config view --raw` prints the full kubeconfig " +
					"including client-certificate/key-data and bearer tokens — " +
					"any script-captured stdout exfiltrates the creds. Emit " +
					"the specific field with `-o jsonpath='…'`.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}

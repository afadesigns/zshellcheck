package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1868",
		Title:    "Error on `gcloud config set auth/disable_ssl_validation true` — disables TLS on every later gcloud call",
		Severity: SeverityError,
		Description: "`gcloud config set auth/disable_ssl_validation true` writes the flag into " +
			"the active configuration file, so every subsequent `gcloud` invocation on " +
			"that machine stops verifying the Google API certificate until someone " +
			"reverses it. A MITM holding a self-signed cert can then intercept service " +
			"account tokens, project-level credentials, and every deploy that runs under " +
			"the same user. Remove the setting (`gcloud config unset " +
			"auth/disable_ssl_validation`), and if a corporate proxy really needs a custom " +
			"CA use `core/custom_ca_certs_file` to pin it rather than disabling the check.",
		Check: checkZC1868,
	})
}

func checkZC1868(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "gcloud" {
		return nil
	}
	args := cmd.Arguments
	if len(args) < 4 {
		return nil
	}
	if args[0].String() != "config" || args[1].String() != "set" {
		return nil
	}
	key := args[2].String()
	val := strings.ToLower(args[3].String())
	if key == "auth/disable_ssl_validation" && (val == "true" || val == "1" || val == "on") {
		return []Violation{{
			KataID: "ZC1868",
			Message: "`gcloud config set auth/disable_ssl_validation " + args[3].String() +
				"` turns off TLS for every later `gcloud` call — service-account " +
				"tokens and deploys become interceptable. Unset it; pin custom " +
				"CAs via `core/custom_ca_certs_file`.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityError,
		}}
	}
	return nil
}
